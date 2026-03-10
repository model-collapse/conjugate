package diagon

/*
#cgo CFLAGS: -I${SRCDIR}/../../../src/3rdparty/diagon/src/core/include
#cgo LDFLAGS: -L${SRCDIR}/build -ldiagon -L${SRCDIR}/../../../src/3rdparty/diagon/build/src/core -ldiagon_core -lz -lzstd -llz4 -Wl,-rpath,${SRCDIR}/build -Wl,-rpath,${SRCDIR}/../../../src/3rdparty/diagon/build/src/core
#include <stdlib.h>
#include <string.h>
#include "diagon/c_api/diagon_c_api.h"

// TermBucketC holds a terms aggregation result from C-level computation
typedef struct {
    char key[256];
    int doc_count;
} TermBucketC;

// compute_terms_agg_stored samples documents uniformly across the reader,
// reads a single stored field, and computes terms aggregation entirely in C.
// max_docs_to_scan caps how many documents are sampled (0 = all).
// Uses strided sampling (stride = max_doc / max_docs_to_scan) to ensure coverage
// across all segments — critical when early segments lack stored fields.
// Returns number of unique terms found. Results written to out_buckets (sorted by doc_count desc).
static int compute_terms_agg_stored(
    DiagonIndexReader reader,
    const char* field_name,
    TermBucketC* out_buckets,
    int max_buckets,
    int max_docs_to_scan)
{
    int max_doc = (int)diagon_reader_max_doc(reader);
    if (max_doc <= 0) return 0;

    // Use strided sampling for large indices to cover all segments uniformly
    int stride = 1;
    int scan_limit = max_doc;
    if (max_docs_to_scan > 0 && max_docs_to_scan < max_doc) {
        stride = max_doc / max_docs_to_scan;
        if (stride < 1) stride = 1;
        scan_limit = max_doc; // iterate full range with stride
    }

    // Simple hash map: use linear probing with 4x overallocation
    // Max unique terms we track
    int map_cap = 8192;
    // Inline hash map entries
    typedef struct {
        char key[256];
        int count;
        int used;
    } Entry;

    Entry* map = (Entry*)calloc(map_cap, sizeof(Entry));
    if (!map) return 0;

    int unique_count = 0;
    char tmp[4096];

    int sampled = 0;
    for (int doc_id = 0; doc_id < scan_limit; doc_id += stride) {
        sampled++;
        DiagonDocument doc = diagon_reader_get_document(reader, doc_id);
        if (!doc) continue;

        if (diagon_document_get_field_value(doc, field_name, tmp, sizeof(tmp))) {
            // Hash the field value
            unsigned int hash = 5381;
            for (const char* p = tmp; *p; p++) {
                hash = ((hash << 5) + hash) + (unsigned char)*p;
            }
            int idx = hash % map_cap;

            // Linear probing
            int found = 0;
            for (int probe = 0; probe < map_cap; probe++) {
                int slot = (idx + probe) % map_cap;
                if (!map[slot].used) {
                    // New entry
                    strncpy(map[slot].key, tmp, 255);
                    map[slot].key[255] = '\0';
                    map[slot].count = 1;
                    map[slot].used = 1;
                    unique_count++;
                    found = 1;
                    break;
                }
                if (strcmp(map[slot].key, tmp) == 0) {
                    // Existing entry
                    map[slot].count++;
                    found = 1;
                    break;
                }
            }
        }
        diagon_free_document(doc);
    }

    // Collect all entries into output, sorted by count desc
    // Scale counts by stride to estimate full-index frequencies
    TermBucketC* all = (TermBucketC*)calloc(unique_count, sizeof(TermBucketC));
    if (!all) { free(map); return 0; }

    int out_idx = 0;
    for (int i = 0; i < map_cap && out_idx < unique_count; i++) {
        if (map[i].used) {
            strncpy(all[out_idx].key, map[i].key, 255);
            all[out_idx].key[255] = '\0';
            all[out_idx].doc_count = map[i].count * stride; // scale to full index
            out_idx++;
        }
    }
    free(map);

    // Simple selection sort for top-N (efficient when max_buckets << unique_count)
    int n = unique_count < max_buckets ? unique_count : max_buckets;
    for (int i = 0; i < n; i++) {
        int max_idx = i;
        for (int j = i + 1; j < unique_count; j++) {
            if (all[j].doc_count > all[max_idx].doc_count) {
                max_idx = j;
            }
        }
        if (max_idx != i) {
            TermBucketC tmp_bucket = all[i];
            all[i] = all[max_idx];
            all[max_idx] = tmp_bucket;
        }
        out_buckets[i] = all[i];
    }
    free(all);

    return n;
}

// compute_cardinality_sampled counts approximate unique values for a field
// by sampling uniformly-spaced documents. Uses a hash set with early termination:
// stops sampling when no new unique values found in last `patience` docs.
// sample_size controls max docs to read (stride = maxDoc / sample_size).
// Returns estimated cardinality.
static int64_t compute_cardinality_sampled(
    DiagonIndexReader reader,
    const char* field_name,
    int sample_size)
{
    int max_doc = (int)diagon_reader_max_doc(reader);
    if (max_doc <= 0) return 0;

    int stride = max_doc / sample_size;
    if (stride < 1) stride = 1;
    int is_full_scan = (stride == 1);

    // Hash set for unique counting (8K entries, ~260KB)
    int map_cap = 8192;
    typedef struct { char key[256]; int used; } Entry;
    Entry* map = (Entry*)calloc(map_cap, sizeof(Entry));
    if (!map) return 0;

    int unique_count = 0;
    int sampled = 0;
    int since_last_new = 0; // early termination counter
    int patience = 200;     // stop if 200 consecutive samples yield no new value
    char tmp[4096];

    for (int doc_id = 0; doc_id < max_doc; doc_id += stride) {
        DiagonDocument doc = diagon_reader_get_document(reader, doc_id);
        if (!doc) continue;
        sampled++;

        int is_new = 0;
        if (diagon_document_get_field_value(doc, field_name, tmp, sizeof(tmp))) {
            unsigned int hash = 5381;
            for (const char* p = tmp; *p; p++) {
                hash = ((hash << 5) + hash) + (unsigned char)*p;
            }
            int idx = hash % map_cap;

            for (int probe = 0; probe < 512; probe++) {
                int slot = (idx + probe) % map_cap;
                if (!map[slot].used) {
                    strncpy(map[slot].key, tmp, 255);
                    map[slot].key[255] = '\0';
                    map[slot].used = 1;
                    unique_count++;
                    is_new = 1;
                    break;
                }
                if (strcmp(map[slot].key, tmp) == 0) {
                    break;
                }
            }
        }
        diagon_free_document(doc);

        if (is_new) {
            since_last_new = 0;
        } else {
            since_last_new++;
            // Early termination: if no new values in `patience` samples,
            // we've likely seen all unique values (low/medium cardinality).
            if (!is_full_scan && since_last_new >= patience) {
                break;
            }
        }
    }
    free(map);

    if (is_full_scan || unique_count == 0) {
        return (int64_t)unique_count;
    }

    // For sampled data: if unique_count < 0.5 * sampled, we likely found most
    // unique values (low cardinality). Return exact count.
    if (unique_count < sampled / 2) {
        return (int64_t)unique_count;
    }
    // High cardinality estimate: scale by sampling ratio
    int64_t estimated = (int64_t)unique_count * max_doc / sampled;
    return estimated;
}

// batch_extract_field_values extracts a single stored field value from ALL documents
// in the TopDocs results in a single CGO call.
// Output: values concatenated in out_buf with null separators, lengths in out_lengths.
// Returns total bytes written to out_buf.
static int batch_extract_field_values(
    DiagonIndexReader reader,
    DiagonTopDocs top_docs,
    const char* field_name,
    char* out_buf,
    int buf_size,
    int* out_lengths,
    int max_docs)
{
    int num_results = diagon_top_docs_score_docs_length(top_docs);
    if (num_results > max_docs) num_results = max_docs;

    int offset = 0;
    char tmp[4096];

    for (int i = 0; i < num_results; i++) {
        out_lengths[i] = 0;

        DiagonScoreDoc sd = diagon_top_docs_score_doc_at(top_docs, i);
        if (!sd) continue;

        int doc_id = diagon_score_doc_get_doc(sd);
        DiagonDocument doc = diagon_reader_get_document(reader, doc_id);
        if (!doc) continue;

        if (diagon_document_get_field_value(doc, field_name, tmp, sizeof(tmp))) {
            int len = (int)strlen(tmp);
            if (offset + len + 1 <= buf_size) {
                memcpy(out_buf + offset, tmp, len + 1);
                out_lengths[i] = len;
                offset += len + 1;
            }
        }
        diagon_free_document(doc);
    }
    return offset;
}

// batch_extract_two_fields extracts two stored field values from ALL documents
// in a single CGO call. Used for aggregations with sub-aggs.
static int batch_extract_two_fields(
    DiagonIndexReader reader,
    DiagonTopDocs top_docs,
    const char* field1_name,
    const char* field2_name,
    char* out_buf1, int buf1_size, int* out_lengths1,
    char* out_buf2, int buf2_size, int* out_lengths2,
    int max_docs)
{
    int num_results = diagon_top_docs_score_docs_length(top_docs);
    if (num_results > max_docs) num_results = max_docs;

    int offset1 = 0, offset2 = 0;
    char tmp[4096];

    for (int i = 0; i < num_results; i++) {
        out_lengths1[i] = 0;
        out_lengths2[i] = 0;

        DiagonScoreDoc sd = diagon_top_docs_score_doc_at(top_docs, i);
        if (!sd) continue;

        int doc_id = diagon_score_doc_get_doc(sd);
        DiagonDocument doc = diagon_reader_get_document(reader, doc_id);
        if (!doc) continue;

        if (diagon_document_get_field_value(doc, field1_name, tmp, sizeof(tmp))) {
            int len = (int)strlen(tmp);
            if (offset1 + len + 1 <= buf1_size) {
                memcpy(out_buf1 + offset1, tmp, len + 1);
                out_lengths1[i] = len;
                offset1 += len + 1;
            }
        }
        if (diagon_document_get_field_value(doc, field2_name, tmp, sizeof(tmp))) {
            int len = (int)strlen(tmp);
            if (offset2 + len + 1 <= buf2_size) {
                memcpy(out_buf2 + offset2, tmp, len + 1);
                out_lengths2[i] = len;
                offset2 += len + 1;
            }
        }
        diagon_free_document(doc);
    }
    return num_results;
}

// batch_extract_multi_fields extracts N stored field values from ALL documents
// in the TopDocs results. Each document is loaded exactly once from stored fields,
// then all N field values are extracted. This is critical for composite aggs with
// 3+ keys where loading each doc once (vs N times) cuts I/O by Nx.
// field_names: array of N field name pointers
// out_bufs: array of N output buffer pointers (one per field)
// buf_sizes: array of N buffer sizes
// out_lengths: array of N int-array pointers (lengths for each field per doc)
// num_fields: number of fields
// max_docs: maximum documents to process
static int batch_extract_multi_fields(
    DiagonIndexReader reader,
    DiagonTopDocs top_docs,
    const char** field_names,
    char** out_bufs,
    int* buf_sizes,
    int** out_lengths,
    int num_fields,
    int max_docs)
{
    int num_results = diagon_top_docs_score_docs_length(top_docs);
    if (num_results > max_docs) num_results = max_docs;

    int offsets[16]; // up to 16 fields
    if (num_fields > 16) num_fields = 16;
    for (int f = 0; f < num_fields; f++) {
        offsets[f] = 0;
    }

    char tmp[4096];

    for (int i = 0; i < num_results; i++) {
        for (int f = 0; f < num_fields; f++) {
            out_lengths[f][i] = 0;
        }

        DiagonScoreDoc sd = diagon_top_docs_score_doc_at(top_docs, i);
        if (!sd) continue;

        int doc_id = diagon_score_doc_get_doc(sd);
        DiagonDocument doc = diagon_reader_get_document(reader, doc_id);
        if (!doc) continue;

        for (int f = 0; f < num_fields; f++) {
            if (diagon_document_get_field_value(doc, field_names[f], tmp, sizeof(tmp))) {
                int len = (int)strlen(tmp);
                if (offsets[f] + len + 1 <= buf_sizes[f]) {
                    memcpy(out_bufs[f] + offsets[f], tmp, len + 1);
                    out_lengths[f][i] = len;
                    offsets[f] += len + 1;
                }
            }
        }
        diagon_free_document(doc);
    }
    return num_results;
}

// batch_scan_field_values scans sequential doc IDs [0, max_docs) and extracts
// a single stored field value from each document. Unlike batch_extract_field_values
// which uses TopDocs from a search, this iterates doc IDs directly — ideal for
// match_all aggregation queries that don't need search scoring.
// Output: values concatenated in out_buf with null separators, lengths in out_lengths.
// Returns number of valid documents scanned (skipping deleted docs).
static int batch_scan_field_values(
    DiagonIndexReader reader,
    const char* field_name,
    char* out_buf,
    int buf_size,
    int* out_lengths,
    int max_docs)
{
    int max_doc_id = (int)diagon_reader_max_doc(reader);
    if (max_doc_id <= 0) return 0;

    int offset = 0;
    char tmp[4096];
    int valid_docs = 0;

    for (int doc_id = 0; doc_id < max_doc_id && valid_docs < max_docs; doc_id++) {
        DiagonDocument doc = diagon_reader_get_document(reader, doc_id);
        if (!doc) continue;

        out_lengths[valid_docs] = 0;
        if (diagon_document_get_field_value(doc, field_name, tmp, sizeof(tmp))) {
            int len = (int)strlen(tmp);
            if (offset + len + 1 <= buf_size) {
                memcpy(out_buf + offset, tmp, len + 1);
                out_lengths[valid_docs] = len;
                offset += len + 1;
            }
        }
        diagon_free_document(doc);
        valid_docs++;
    }

    return valid_docs;
}

// batch_scan_field_values_from scans sequential doc IDs starting from start_doc_id,
// extracting up to max_docs stored field values. Returns number of valid docs scanned.
// Sets *next_doc_id to the doc ID to resume from on the next call (-1 if done).
static int batch_scan_field_values_from(
    DiagonIndexReader reader,
    const char* field_name,
    char* out_buf,
    int buf_size,
    int* out_lengths,
    int max_docs,
    int start_doc_id,
    int* next_doc_id)
{
    int max_doc_id = (int)diagon_reader_max_doc(reader);
    if (max_doc_id <= 0) { *next_doc_id = -1; return 0; }

    int offset = 0;
    char tmp[4096];
    int valid_docs = 0;
    int doc_id;

    for (doc_id = start_doc_id; doc_id < max_doc_id && valid_docs < max_docs; doc_id++) {
        DiagonDocument doc = diagon_reader_get_document(reader, doc_id);
        if (!doc) continue;

        out_lengths[valid_docs] = 0;
        if (diagon_document_get_field_value(doc, field_name, tmp, sizeof(tmp))) {
            int len = (int)strlen(tmp);
            if (offset + len + 1 <= buf_size) {
                memcpy(out_buf + offset, tmp, len + 1);
                out_lengths[valid_docs] = len;
                offset += len + 1;
            }
        }
        diagon_free_document(doc);
        valid_docs++;
    }

    *next_doc_id = (doc_id < max_doc_id) ? doc_id : -1;
    return valid_docs;
}

// (batch_scan_numeric_values removed: stored-field scan is O(N*7µs) = ~800s for 116M docs,
// making it slower than the C++ range query path which uses NumericDocValues column store)

// ---------- BKD-based aggregation functions ----------
// Range agg uses individual BKD range queries (few buckets: 3-10 typical).
// Date histogram uses diagon_compute_histogram C API (single-pass BKD traversal).

typedef struct {
    double from_val;
    double to_val;
    int has_from;
    int has_to;
    int64_t doc_count;
} RangeAggBucketC;

// compute_range_agg_bkd counts docs per numeric range bucket using BKD range queries.
// Ranges are specified as (from, to) pairs with optional bounds.
// Returns num_ranges (all ranges are filled with doc_count).
static int compute_range_agg_bkd(
    DiagonIndexSearcher searcher,
    const char* field_name,
    RangeAggBucketC* ranges,
    int num_ranges,
    DiagonQuery filter_query)
{
    for (int i = 0; i < num_ranges; i++) {
        double lower = ranges[i].has_from ? ranges[i].from_val : -1e308;
        double upper = ranges[i].has_to ? ranges[i].to_val : 1e308;

        DiagonQuery rangeQ = diagon_create_double_range_query(
            field_name, lower, upper, true, false);
        if (!rangeQ) { ranges[i].doc_count = 0; continue; }

        DiagonQuery searchQ = NULL;
        if (filter_query) {
            DiagonQuery boolQ = diagon_create_bool_query();
            if (!boolQ) { diagon_free_query(rangeQ); ranges[i].doc_count = 0; continue; }
            diagon_bool_query_add_must(boolQ, filter_query);
            diagon_bool_query_add_filter(boolQ, rangeQ);
            searchQ = diagon_bool_query_build(boolQ);
            diagon_free_query(rangeQ);
        } else {
            DiagonQuery matchAll = diagon_create_match_all_query();
            if (!matchAll) { diagon_free_query(rangeQ); ranges[i].doc_count = 0; continue; }
            DiagonQuery boolQ = diagon_create_bool_query();
            if (!boolQ) { diagon_free_query(rangeQ); diagon_free_query(matchAll); ranges[i].doc_count = 0; continue; }
            diagon_bool_query_add_must(boolQ, matchAll);
            diagon_bool_query_add_filter(boolQ, rangeQ);
            searchQ = diagon_bool_query_build(boolQ);
            diagon_free_query(matchAll);
            diagon_free_query(rangeQ);
        }

        if (!searchQ) { ranges[i].doc_count = 0; continue; }

        DiagonTopDocs td = diagon_search(searcher, searchQ, 1);
        diagon_free_query(searchQ);

        if (td) {
            ranges[i].doc_count = diagon_top_docs_total_hits(td);
            diagon_free_top_docs(td);
        } else {
            ranges[i].doc_count = 0;
        }
    }
    return num_ranges;
}

// CStringArena: memory arena for batched CString allocations.
// Reduces malloc/free overhead from ~80,000 calls to 1 per 5K-doc batch.
typedef struct {
    char* buf;
    size_t pos;
    size_t cap;
} CStringArena;

static CStringArena* arena_create(size_t cap) {
    CStringArena* a = (CStringArena*)malloc(sizeof(CStringArena));
    if (!a) return NULL;
    a->buf = (char*)malloc(cap);
    if (!a->buf) { free(a); return NULL; }
    a->pos = 0;
    a->cap = cap;
    return a;
}

static char* arena_cstring_bytes(CStringArena* a, const void* src, size_t len) {
    if (!a || a->pos + len + 1 > a->cap) return NULL;
    char* ptr = a->buf + a->pos;
    memcpy(ptr, src, len);
    ptr[len] = '\0';
    a->pos += len + 1;
    return ptr;
}

static void arena_free(CStringArena* a) {
    if (a) { free(a->buf); free(a); }
}
*/
import "C"

import (
	"fmt"
	"math"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"

	json "github.com/goccy/go-json"

	"go.uber.org/zap"
)

// countCFSFiles counts compound segment files (.cfs) in a directory.
func countCFSFiles(dir string) int {
	// Segments exist as either compound (.cfs) or non-compound (separate .tim files).
	// Total segment count = compound + non-compound.
	cfsMatches, _ := filepath.Glob(filepath.Join(dir, "*.cfs"))
	timMatches, _ := filepath.Glob(filepath.Join(dir, "*.tim"))
	return len(cfsMatches) + len(timMatches)
}

// tryParseDateToEpochMs attempts to parse a string as an ISO date and returns epoch millis.
// Returns (epochMs, true) if successful, (0, false) otherwise.
func tryParseDateToEpochMs(s string) (int64, bool) {
	// Fast pre-filter: reject strings that can't possibly be dates.
	// All supported date formats start with "YYYY-MM" (digit at [0], '-' at [4]).
	// Minimum length is 10 ("2006-01-02").
	n := len(s)
	if n < 10 || s[0] < '0' || s[0] > '9' || s[4] != '-' || s[7] != '-' {
		return 0, false
	}

	// Try formats in order of likelihood for log/observability data.
	// Most timestamps are RFC3339 with milliseconds.
	var t time.Time
	var err error

	if n >= 24 && s[10] == 'T' && (s[n-1] == 'Z' || s[n-1] == '0') {
		// Fast path: ISO8601 with timezone (covers ~95% of log timestamps)
		t, err = time.Parse(time.RFC3339Nano, s)
		if err == nil {
			return t.UnixMilli(), true
		}
		t, err = time.Parse("2006-01-02T15:04:05.000Z", s)
		if err == nil {
			return t.UnixMilli(), true
		}
	}

	if n >= 19 && s[10] == 'T' {
		t, err = time.Parse(time.RFC3339, s)
		if err == nil {
			return t.UnixMilli(), true
		}
		t, err = time.Parse("2006-01-02T15:04:05", s)
		if err == nil {
			return t.UnixMilli(), true
		}
		t, err = time.Parse("2006-01-02T15:04:05.000", s)
		if err == nil {
			return t.UnixMilli(), true
		}
	}

	if n == 10 {
		t, err = time.Parse("2006-01-02", s)
		if err == nil {
			return t.UnixMilli(), true
		}
	}

	return 0, false
}

// parseDateMathToEpochMs parses OpenSearch date math expressions like "now-7d", "now-1h".
func parseDateMathToEpochMs(expr string) int64 {
	now := time.Now()
	if expr == "now" {
		return now.UnixMilli()
	}
	// Parse "now-Xd", "now-Xh", "now-Xm", "now+Xd" etc.
	rest := strings.TrimPrefix(expr, "now")
	if rest == "" {
		return now.UnixMilli()
	}
	sign := 1
	if rest[0] == '-' {
		sign = -1
		rest = rest[1:]
	} else if rest[0] == '+' {
		rest = rest[1:]
	}
	// Parse number
	numStr := ""
	for i, c := range rest {
		if c >= '0' && c <= '9' {
			numStr += string(c)
		} else {
			rest = rest[i:]
			break
		}
	}
	if numStr == "" {
		return now.UnixMilli()
	}
	num := 0
	for _, c := range numStr {
		num = num*10 + int(c-'0')
	}
	var dur time.Duration
	switch {
	case strings.HasPrefix(rest, "d"):
		dur = time.Duration(num) * 24 * time.Hour
	case strings.HasPrefix(rest, "h"):
		dur = time.Duration(num) * time.Hour
	case strings.HasPrefix(rest, "m"):
		dur = time.Duration(num) * time.Minute
	case strings.HasPrefix(rest, "s"):
		dur = time.Duration(num) * time.Second
	case strings.HasPrefix(rest, "w"):
		dur = time.Duration(num) * 7 * 24 * time.Hour
	case strings.HasPrefix(rest, "M"):
		return now.AddDate(0, sign*num, 0).UnixMilli()
	case strings.HasPrefix(rest, "y"):
		return now.AddDate(sign*num, 0, 0).UnixMilli()
	default:
		dur = time.Duration(num) * 24 * time.Hour // default to days
	}
	return now.Add(time.Duration(sign) * dur).UnixMilli()
}

// DiagonBridge provides a Go interface to the real Diagon C++ search engine
type DiagonBridge struct {
	config     *Config
	logger     *zap.Logger
	shards     map[string]*Shard
	mu         sync.RWMutex
}

// Config holds Diagon configuration
type Config struct {
	DataDir     string
	SIMDEnabled bool
	Logger      *zap.Logger
}

// NewDiagonBridge creates a new Diagon bridge
func NewDiagonBridge(cfg *Config) (*DiagonBridge, error) {
	if cfg.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	bridge := &DiagonBridge{
		config: cfg,
		logger: cfg.Logger,
		shards: make(map[string]*Shard),
	}

	return bridge, nil
}

// Start starts the Diagon engine
func (db *DiagonBridge) Start() error {
	db.logger.Info("Starting real Diagon C++ search engine",
		zap.String("data_dir", db.config.DataDir),
		zap.Bool("simd_enabled", db.config.SIMDEnabled))

	return nil
}

// Stop stops the Diagon engine
func (db *DiagonBridge) Stop() error {
	db.logger.Info("Stopping Diagon engine")

	db.mu.Lock()
	defer db.mu.Unlock()

	// Close all shards
	for path, shard := range db.shards {
		db.logger.Info("Closing Diagon shard", zap.String("path", path))
		if err := shard.Close(); err != nil {
			db.logger.Error("Error closing shard", zap.String("path", path), zap.Error(err))
		}
	}

	return nil
}

// CreateShard creates a new shard using real Diagon IndexWriter
func (db *DiagonBridge) CreateShard(path string) (*Shard, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if shard already exists
	if _, exists := db.shards[path]; exists {
		return nil, fmt.Errorf("shard at path %s already exists", path)
	}

	// Open directory
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	dir := C.diagon_open_mmap_directory(cPath) // Use MMapDirectory for performance
	if dir == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return nil, fmt.Errorf("failed to open directory: %s", errMsg)
	}

	// Create IndexWriter config
	config := C.diagon_create_index_writer_config()
	C.diagon_config_set_ram_buffer_size(config, 256.0)                  // 256MB buffer (tuned for bulk throughput)
	C.diagon_config_set_max_buffered_docs(config, 100000)              // Let RAM buffer control flushing (default 1000 caused segment explosion)
	C.diagon_config_set_open_mode(config, 2)                            // CREATE_OR_APPEND
	C.diagon_config_set_commit_on_close(config, true)

	// Create IndexWriter
	writer := C.diagon_create_index_writer(dir, config)
	C.diagon_free_index_writer_config(config)

	if writer == nil {
		C.diagon_close_directory(dir)
		errMsg := C.GoString(C.diagon_last_error())
		return nil, fmt.Errorf("failed to create IndexWriter: %s", errMsg)
	}

	shard := &Shard{
		path:      path,
		bridge:    db,
		directory: dir,
		writer:    writer,
		reader:    nil, // Will be opened when needed
		logger:    db.logger.With(zap.String("shard_path", path)),
	}

	db.shards[path] = shard

	shard.logger.Info("Created real Diagon shard with IndexWriter")

	return shard, nil
}

// GetShard retrieves an existing shard
func (db *DiagonBridge) GetShard(path string) (*Shard, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	shard, exists := db.shards[path]
	if !exists {
		return nil, fmt.Errorf("shard at path %s not found", path)
	}

	return shard, nil
}

// Shard represents a real Diagon shard with IndexWriter/IndexReader
type Shard struct {
	path         string
	bridge       *DiagonBridge
	directory    C.DiagonDirectory
	writer       C.DiagonIndexWriter
	reader       C.DiagonIndexReader
	searcher     C.DiagonIndexSearcher
	readerDirty  bool // true when writes occurred since last reader open
	logger       *zap.Logger
	mu           sync.RWMutex
	// termsAggCache caches terms aggregation results keyed by "field:size".
	// Invalidated when reader is reopened (readerDirty).
	termsAggCache   map[string][]TermBucket
	termsAggCacheMu sync.RWMutex
	// columnCache stores extracted field values as Go string slices.
	// After the first batch extraction, subsequent agg queries on the same
	// fields skip CGO entirely and operate on pure Go arrays.
	// Invalidated on commit/refresh (readerDirty).
	columnCache   map[string][]string // field -> values for docs 0..N
	columnCacheMu sync.RWMutex
	columnCacheN  int // number of docs in cache (all columns same length)

}

// isKeywordLike returns true if a string value looks like a keyword (enum/identifier)
// rather than full text. Keywords are: short (<256 chars), no whitespace, typically
// identifiers, enum values, or short labels.
func isKeywordLike(s string) bool {
	if len(s) > 256 || len(s) == 0 {
		return false
	}
	for _, c := range s {
		if c == ' ' || c == '\t' || c == '\n' {
			return false
		}
	}
	return true
}

// flattenMap recursively flattens nested maps into dotted field paths.
// e.g., {"cloud": {"region": "us-west-2"}} becomes {"cloud.region": "us-west-2"}
// Arrays of primitives are joined with spaces for text indexing.
func flattenMap(prefix string, m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	flattenMapInto(prefix, m, result)
	return result
}

// flattenMapInto writes flattened key-value pairs directly into dst,
// avoiding intermediate map allocations at each recursion level.
func flattenMapInto(prefix string, m map[string]interface{}, dst map[string]interface{}) {
	for key, value := range m {
		var fullKey string
		if prefix == "" {
			fullKey = key
		} else {
			fullKey = prefix + "." + key
		}
		switch v := value.(type) {
		case map[string]interface{}:
			flattenMapInto(fullKey, v, dst)
		case []interface{}:
			parts := make([]string, 0, len(v))
			for _, elem := range v {
				parts = append(parts, fmt.Sprintf("%v", elem))
			}
			if len(parts) > 0 {
				dst[fullKey] = strings.Join(parts, " ")
			}
		default:
			dst[fullKey] = value
		}
	}
}

// IndexDocument indexes a document using real Diagon IndexWriter
func (s *Shard) IndexDocument(docID string, doc map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debug("IndexDocument",
		zap.String("doc_id", docID),
		zap.Int("num_fields", len(doc)))

	// Create Diagon document
	diagonDoc := C.diagon_create_document()
	defer C.diagon_free_document(diagonDoc)

	s.logger.Debug("Created Diagon document object", zap.String("doc_id", docID))

	// Add ID field - both indexed (for searching) and stored (for retrieval)
	cDocID := C.CString(docID)
	defer C.free(unsafe.Pointer(cDocID))
	cIDFieldName := C.CString("_id")
	defer C.free(unsafe.Pointer(cIDFieldName))

	// Add as StringField for exact-match searching (indexed, not analyzed)
	idField := C.diagon_create_string_field(cIDFieldName, cDocID)
	C.diagon_document_add_field(diagonDoc, idField)

	// ALSO add as StoredField so we can retrieve it
	storedIDField := C.diagon_create_stored_field(cIDFieldName, cDocID)
	C.diagon_document_add_field(diagonDoc, storedIDField)

	// Store full _source as JSON for reliable document retrieval
	sourceJSON, err := json.Marshal(doc)
	if err == nil {
		cSourceFieldName := C.CString("_source")
		defer C.free(unsafe.Pointer(cSourceFieldName))
		cSourceValue := C.CString(string(sourceJSON))
		defer C.free(unsafe.Pointer(cSourceValue))
		sourceField := C.diagon_create_stored_field(cSourceFieldName, cSourceValue)
		C.diagon_document_add_field(diagonDoc, sourceField)
	}

	// Flatten nested objects into dotted field paths for indexing
	flatDoc := flattenMap("", doc)

	// Add other fields
	for key, value := range flatDoc {
		cFieldName := C.CString(key)
		defer C.free(unsafe.Pointer(cFieldName))

		s.logger.Debug("Indexing field",
			zap.String("field", key),
			zap.String("type", fmt.Sprintf("%T", value)))

		switch v := value.(type) {
		case string:
			// Check if this is a date string - if so, index as double (epoch millis)
			// We use double instead of long because the C API's numeric range query
			// uses bit_cast<int64>(double) which corrupts the value. Double range
			// query compares doubles directly, which works for epoch millis.
			if epochMs, isDate := tryParseDateToEpochMs(v); isDate {
				// Index as double for range queries (doc values for sort/agg)
				doubleField := C.diagon_create_indexed_double_field(cFieldName, C.double(float64(epochMs)))
				C.diagon_document_add_field(diagonDoc, doubleField)
				// Also add BKD point field for O(log N) range queries
				pointField := C.diagon_create_double_point_field(cFieldName, C.double(float64(epochMs)))
				C.diagon_document_add_field(diagonDoc, pointField)
				// Also store the original string value
				cValue := C.CString(v)
				defer C.free(unsafe.Pointer(cValue))
				storedField := C.diagon_create_stored_field(cFieldName, cValue)
				C.diagon_document_add_field(diagonDoc, storedField)
			} else if isKeywordLike(v) {
				// Keyword-like string: index as StringField (exact, not analyzed)
				// This enables fast terms aggregation via inverted index O(unique_terms)
				// instead of O(total_docs) document extraction.
				cValue := C.CString(v)
				defer C.free(unsafe.Pointer(cValue))
				stringField := C.diagon_create_string_field(cFieldName, cValue)
				C.diagon_document_add_field(diagonDoc, stringField)
				storedField := C.diagon_create_stored_field(cFieldName, cValue)
				C.diagon_document_add_field(diagonDoc, storedField)
			} else {
				// Regular text field for strings (analyzed, indexed, stored)
				cValue := C.CString(v)
				defer C.free(unsafe.Pointer(cValue))
				field := C.diagon_create_text_field(cFieldName, cValue)
				C.diagon_document_add_field(diagonDoc, field)
				// Also add as StringField so terms aggregation works on the exact value
				stringField := C.diagon_create_string_field(cFieldName, cValue)
				C.diagon_document_add_field(diagonDoc, stringField)
			}

		case int, int32, int64:
			// Index integers as doubles for consistent range query support.
			// The double_range_query compares doubles directly, whereas
			// numeric_range_query uses bit_cast which corrupts values.
			val := int64(0)
			switch n := v.(type) {
			case int:
				val = int64(n)
			case int32:
				val = int64(n)
			case int64:
				val = n
			}
			field := C.diagon_create_indexed_double_field(cFieldName, C.double(float64(val)))
			C.diagon_document_add_field(diagonDoc, field)
			pointField := C.diagon_create_double_point_field(cFieldName, C.double(float64(val)))
			C.diagon_document_add_field(diagonDoc, pointField)

			cValueStr := C.CString(fmt.Sprintf("%d", val))
			defer C.free(unsafe.Pointer(cValueStr))
			storedField := C.diagon_create_stored_field(cFieldName, cValueStr)
			C.diagon_document_add_field(diagonDoc, storedField)

		case float32, float64:
			val := float64(0)
			switch f := v.(type) {
			case float32:
				val = float64(f)
			case float64:
				val = f
			}
			field := C.diagon_create_indexed_double_field(cFieldName, C.double(val))
			C.diagon_document_add_field(diagonDoc, field)
			pointField := C.diagon_create_double_point_field(cFieldName, C.double(val))
			C.diagon_document_add_field(diagonDoc, pointField)

			cValueStr := C.CString(fmt.Sprintf("%f", val))
			defer C.free(unsafe.Pointer(cValueStr))
			storedField := C.diagon_create_stored_field(cFieldName, cValueStr)
			C.diagon_document_add_field(diagonDoc, storedField)

		default:
			// Convert to JSON string for complex types
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				s.logger.Warn("Failed to marshal field, skipping",
					zap.String("field", key),
					zap.Error(err))
				continue
			}
			cValue := C.CString(string(jsonBytes))
			defer C.free(unsafe.Pointer(cValue))
			field := C.diagon_create_stored_field(cFieldName, cValue)
			C.diagon_document_add_field(diagonDoc, field)
		}
	}

	// Add document to IndexWriter
	result := C.diagon_add_document(s.writer, diagonDoc)

	if !result {
		errMsg := C.GoString(C.diagon_last_error())
		s.logger.Error("C.diagon_add_document FAILED",
			zap.String("doc_id", docID),
			zap.String("error", errMsg))
		return fmt.Errorf("failed to add document: %s", errMsg)
	}

	s.readerDirty = true
	s.logger.Debug("Document added to RAM buffer",
		zap.String("doc_id", docID))

	return nil
}


// Cached C strings for field names that are identical across all documents.
// Allocated once (process-lifetime), never freed.
var (
	cachedIDFieldName     *C.char
	cachedSourceFieldName *C.char
	cachedFieldNamesOnce  sync.Once
)

func initCachedFieldNames() {
	cachedFieldNamesOnce.Do(func() {
		cachedIDFieldName = C.CString("_id")
		cachedSourceFieldName = C.CString("_source")
	})
}

// cStringFromBytes allocates a null-terminated C string directly from a Go
// byte slice, avoiding the intermediate string copy that C.CString(string(b))
// would incur. The caller must C.free the returned pointer.
func cStringFromBytes(b []byte) *C.char {
	n := C.size_t(len(b))
	cstr := (*C.char)(C.malloc(n + 1))
	if len(b) > 0 {
		C.memcpy(unsafe.Pointer(cstr), unsafe.Pointer(&b[0]), n)
	}
	// Null-terminate
	*(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cstr)) + uintptr(n))) = 0
	return cstr
}

// arenaAllocCString allocates a null-terminated C string from the arena.
// Returns (cstring, true) if arena succeeded; (cstring, false) if fallback to C.CString.
func arenaAllocCString(arena *C.CStringArena, s string) (*C.char, bool) {
	b := []byte(s)
	if len(b) == 0 {
		b = []byte{0}
		cs := C.arena_cstring_bytes(arena, unsafe.Pointer(&b[0]), 0)
		if cs == nil {
			return C.CString(""), false
		}
		return cs, true
	}
	cs := C.arena_cstring_bytes(arena, unsafe.Pointer(&b[0]), C.size_t(len(b)))
	if cs == nil {
		return C.CString(s), false
	}
	return cs, true
}

// arenaAllocCStringFromBytes allocates from arena using []byte directly.
func arenaAllocCStringFromBytes(arena *C.CStringArena, b []byte) (*C.char, bool) {
	if len(b) == 0 {
		empty := []byte{0}
		cs := C.arena_cstring_bytes(arena, unsafe.Pointer(&empty[0]), 0)
		if cs == nil {
			return cStringFromBytes(nil), false
		}
		return cs, true
	}
	cs := C.arena_cstring_bytes(arena, unsafe.Pointer(&b[0]), C.size_t(len(b)))
	if cs == nil {
		return cStringFromBytes(b), false
	}
	return cs, true
}

// processDocChunk processes a range of documents [start, end) for BulkIndexDocuments.
// Each call uses its own arena and caches — safe for concurrent use.
func processDocChunk(docs []struct {
	ID         string
	Doc        map[string]interface{}
	SourceJSON []byte
}, start, end int, cDocs []C.DiagonDocument) error {
	arenaSize := C.size_t(2 * 1024 * 1024) // 2MB per worker
	arena := C.arena_create(arenaSize)
	if arena == nil {
		return fmt.Errorf("failed to create arena for chunk [%d,%d)", start, end)
	}
	defer C.arena_free(arena)

	fieldNameCache := make(map[string]*C.char, 32)
	valueCache := make(map[string]*C.char, 256)
	fallbackPtrs := make([]unsafe.Pointer, 0, 8)

	for i := start; i < end; i++ {
		item := docs[i]
		diagonDoc := C.diagon_create_document()

		// _id field
		cDocID, inArena := arenaAllocCString(arena, item.ID)
		if !inArena {
			fallbackPtrs = append(fallbackPtrs, unsafe.Pointer(cDocID))
		}
		idField := C.diagon_create_string_field(cachedIDFieldName, cDocID)
		C.diagon_document_add_field(diagonDoc, idField)
		storedIDField := C.diagon_create_stored_field(cachedIDFieldName, cDocID)
		C.diagon_document_add_field(diagonDoc, storedIDField)

		// _source field
		sourceBytes := item.SourceJSON
		if sourceBytes == nil {
			sourceBytes, _ = json.Marshal(item.Doc)
		}
		if len(sourceBytes) > 0 {
			cSourceValue, inArena := arenaAllocCStringFromBytes(arena, sourceBytes)
			if !inArena {
				fallbackPtrs = append(fallbackPtrs, unsafe.Pointer(cSourceValue))
			}
			sourceField := C.diagon_create_stored_field(cachedSourceFieldName, cSourceValue)
			C.diagon_document_add_field(diagonDoc, sourceField)
		}

		// Flatten and index user fields
		flatDoc := flattenMap("", item.Doc)
		for key, value := range flatDoc {
			cFieldName, ok := fieldNameCache[key]
			if !ok {
				cFieldName, inArena = arenaAllocCString(arena, key)
				if !inArena {
					fallbackPtrs = append(fallbackPtrs, unsafe.Pointer(cFieldName))
				}
				fieldNameCache[key] = cFieldName
			}

			switch v := value.(type) {
			case string:
				if epochMs, isDate := tryParseDateToEpochMs(v); isDate {
					doubleField := C.diagon_create_indexed_double_field(cFieldName, C.double(float64(epochMs)))
					C.diagon_document_add_field(diagonDoc, doubleField)
					pointField := C.diagon_create_double_point_field(cFieldName, C.double(float64(epochMs)))
					C.diagon_document_add_field(diagonDoc, pointField)
					cValue, inArena := arenaAllocCString(arena, v)
					if !inArena {
						fallbackPtrs = append(fallbackPtrs, unsafe.Pointer(cValue))
					}
					storedField := C.diagon_create_stored_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, storedField)
				} else if isKeywordLike(v) {
					cValue, cached := valueCache[v]
					if !cached {
						cValue, inArena = arenaAllocCString(arena, v)
						if !inArena {
							fallbackPtrs = append(fallbackPtrs, unsafe.Pointer(cValue))
						}
						valueCache[v] = cValue
					}
					stringField := C.diagon_create_string_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, stringField)
					storedField := C.diagon_create_stored_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, storedField)
				} else {
					cValue, inArena := arenaAllocCString(arena, v)
					if !inArena {
						fallbackPtrs = append(fallbackPtrs, unsafe.Pointer(cValue))
					}
					field := C.diagon_create_text_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, field)
					stringField := C.diagon_create_string_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, stringField)
				}

			case int, int32, int64:
				val := int64(0)
				switch iv := v.(type) {
				case int:
					val = int64(iv)
				case int32:
					val = int64(iv)
				case int64:
					val = iv
				}
				field := C.diagon_create_indexed_double_field(cFieldName, C.double(float64(val)))
				C.diagon_document_add_field(diagonDoc, field)
				pointField := C.diagon_create_double_point_field(cFieldName, C.double(float64(val)))
				C.diagon_document_add_field(diagonDoc, pointField)
				cValueStr, inArena := arenaAllocCString(arena, strconv.FormatInt(val, 10))
				if !inArena {
					fallbackPtrs = append(fallbackPtrs, unsafe.Pointer(cValueStr))
				}
				storedField := C.diagon_create_stored_field(cFieldName, cValueStr)
				C.diagon_document_add_field(diagonDoc, storedField)

			case float32, float64:
				val := float64(0)
				switch fv := v.(type) {
				case float32:
					val = float64(fv)
				case float64:
					val = fv
				}
				field := C.diagon_create_indexed_double_field(cFieldName, C.double(val))
				C.diagon_document_add_field(diagonDoc, field)
				pointField := C.diagon_create_double_point_field(cFieldName, C.double(val))
				C.diagon_document_add_field(diagonDoc, pointField)
				cValueStr, inArena := arenaAllocCString(arena, strconv.FormatFloat(val, 'f', 6, 64))
				if !inArena {
					fallbackPtrs = append(fallbackPtrs, unsafe.Pointer(cValueStr))
				}
				storedField := C.diagon_create_stored_field(cFieldName, cValueStr)
				C.diagon_document_add_field(diagonDoc, storedField)

			default:
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					continue
				}
				cValue, inArena := arenaAllocCStringFromBytes(arena, jsonBytes)
				if !inArena {
					fallbackPtrs = append(fallbackPtrs, unsafe.Pointer(cValue))
				}
				field := C.diagon_create_stored_field(cFieldName, cValue)
				C.diagon_document_add_field(diagonDoc, field)
			}
		}

		cDocs[i] = diagonDoc
	}

	// Free fallback allocations (arena strings freed with arena)
	for _, ptr := range fallbackPtrs {
		C.free(ptr)
	}
	// Field name cache strings are in arena — freed with arena
	return nil
}

// BulkIndexDocuments indexes multiple documents in a batch to reduce CGO overhead.
// Optimizations over per-document IndexDocument:
//   - Uses diagon_add_documents batch API (single mutex acquisition in C++)
//   - Caches "_id" and "_source" field name CStrings at package level
//   - Caches user field name CStrings per batch (all docs share schema)
//   - Avoids per-field defer; collects pointers for batch cleanup
//   - Uses cStringFromBytes to skip []byte→string copy for _source
//   - Uses strconv instead of fmt.Sprintf for numeric formatting
func (s *Shard) BulkIndexDocuments(docs []struct {
	ID         string
	Doc        map[string]interface{}
	SourceJSON []byte
}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	numDocs := len(docs)
	if numDocs == 0 {
		return nil
	}

	initCachedFieldNames()

	s.logger.Debug("BulkIndexDocuments",
		zap.Int("num_docs", numDocs))

	// Allocate C array for batch diagon_add_documents call.
	cDocsRaw := C.malloc(C.size_t(numDocs) * C.size_t(unsafe.Sizeof(uintptr(0))))
	defer C.free(cDocsRaw)
	cDocs := unsafe.Slice((*C.DiagonDocument)(cDocsRaw), numDocs)

	// Parallel document preparation: split across workers.
	// Each worker gets its own arena, field name cache, and value cache.
	numWorkers := 8
	if maxProcs := runtime.GOMAXPROCS(0); maxProcs < numWorkers {
		numWorkers = maxProcs
	}
	if numDocs < numWorkers*10 {
		numWorkers = 1 // sequential for small batches
	}

	if numWorkers == 1 {
		// Sequential path
		if err := processDocChunk(docs, 0, numDocs, cDocs); err != nil {
			return err
		}
	} else {
		// Parallel path
		chunkSize := (numDocs + numWorkers - 1) / numWorkers
		var wg sync.WaitGroup
		workerErrors := make([]error, numWorkers)

		for w := 0; w < numWorkers; w++ {
			start := w * chunkSize
			end := start + chunkSize
			if start >= numDocs {
				break
			}
			if end > numDocs {
				end = numDocs
			}
			wg.Add(1)
			go func(workerID, s, e int) {
				defer wg.Done()
				workerErrors[workerID] = processDocChunk(docs, s, e, cDocs)
			}(w, start, end)
		}
		wg.Wait()

		for _, err := range workerErrors {
			if err != nil {
				// Cleanup any created documents on error
				for i := 0; i < numDocs; i++ {
					if cDocs[i] != nil {
						C.diagon_free_document(cDocs[i])
					}
				}
				return err
			}
		}
	}

	// Single CGO call for entire batch (one mutex acquisition in C++)
	result := C.diagon_add_documents(s.writer, &cDocs[0], C.int(numDocs))

	// Free all document handles
	for i := 0; i < numDocs; i++ {
		C.diagon_free_document(cDocs[i])
	}

	if int(result) < 0 {
		errMsg := C.GoString(C.diagon_last_error())
		return fmt.Errorf("batch add_documents failed: %s", errMsg)
	}
	if int(result) < numDocs {
		return fmt.Errorf("batch add_documents: only %d of %d documents added", int(result), numDocs)
	}

	s.readerDirty = true
	s.logger.Debug("Bulk indexed documents to RAM buffer",
		zap.Int("count", numDocs))

	return nil
}

// Commit commits all pending changes
func (s *Shard) Commit() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !C.diagon_commit(s.writer) {
		errMsg := C.GoString(C.diagon_last_error())
		return fmt.Errorf("commit failed: %s", errMsg)
	}

	s.logger.Debug("Committed changes")
	return nil
}

// Flush flushes buffered documents to disk
func (s *Shard) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !C.diagon_flush(s.writer) {
		errMsg := C.GoString(C.diagon_last_error())
		return fmt.Errorf("flush failed: %s", errMsg)
	}

	s.logger.Debug("Flushed buffered documents")
	return nil
}

// Refresh reopens the reader to see recent changes
func (s *Shard) Refresh() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Commit first to ensure changes are visible
	if !C.diagon_commit(s.writer) {
		errMsg := C.GoString(C.diagon_last_error())
		return fmt.Errorf("commit failed during refresh: %s", errMsg)
	}

	// Close old reader and searcher
	if s.searcher != nil {
		C.diagon_free_index_searcher(s.searcher)
		s.searcher = nil
	}
	if s.reader != nil {
		C.diagon_close_index_reader(s.reader)
		s.reader = nil
	}

	// Open new reader
	s.reader = C.diagon_open_index_reader(s.directory)
	if s.reader == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return fmt.Errorf("failed to reopen reader: %s", errMsg)
	}

	// Create new searcher
	s.searcher = C.diagon_create_index_searcher(s.reader)
	if s.searcher == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return fmt.Errorf("failed to create searcher: %s", errMsg)
	}

	s.readerDirty = false
	s.logger.Debug("Refreshed shard (reopened reader)")
	return nil
}

// WarmReader proactively opens the index reader without committing.
// Call this at startup to avoid the cold-start penalty on the first search.
func (s *Shard) WarmReader() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.reader != nil {
		return nil // Already open
	}

	s.reader = C.diagon_open_index_reader(s.directory)
	if s.reader == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return fmt.Errorf("failed to open reader: %s", errMsg)
	}

	s.searcher = C.diagon_create_index_searcher(s.reader)
	if s.searcher == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return fmt.Errorf("failed to create searcher: %s", errMsg)
	}

	numDocs := int(C.diagon_reader_num_docs(s.reader))
	s.logger.Info("Reader warmed", zap.Int("num_docs", numDocs))
	return nil
}

// convertQueryToDiagon converts a query object to a Diagon query
// This is a helper function used by Search and for recursive bool query parsing
// Caller is responsible for freeing the returned query
func (s *Shard) convertQueryToDiagon(queryObj map[string]interface{}) (C.DiagonQuery, error) {
	var diagonQuery C.DiagonQuery

	// Handle different query types
	if termQuery, ok := queryObj["term"].(map[string]interface{}); ok {
		// Term query: {"term": {"field_name": "term_value"}} or {"term": {"field_name": {"value": "term_value"}}}
		for field, value := range termQuery {
			cField := C.CString(field)
			defer C.free(unsafe.Pointer(cField))

			// Handle both simple and complex term query formats
			var termValue string
			switch v := value.(type) {
			case string:
				termValue = v
			case map[string]interface{}:
				if val, ok := v["value"]; ok {
					termValue = fmt.Sprintf("%v", val)
				}
			default:
				termValue = fmt.Sprintf("%v", v)
			}

			cValue := C.CString(termValue)
			defer C.free(unsafe.Pointer(cValue))

			term := C.diagon_create_term(cField, cValue)
			defer C.diagon_free_term(term)

			diagonQuery = C.diagon_create_term_query(term)
			if diagonQuery == nil {
				errMsg := C.GoString(C.diagon_last_error())
				return nil, fmt.Errorf("failed to create term query: %s", errMsg)
			}
			break // Only support single term for now
		}
	} else if matchQuery, ok := queryObj["match"].(map[string]interface{}); ok {
		// Match query: {"match": {"field_name": "query_text"}} or {"match": {"field_name": {"query": "text"}}}
		// For now, treat match query as term query (no text analysis in Diagon Phase 4)
		for field, value := range matchQuery {
			cField := C.CString(field)
			defer C.free(unsafe.Pointer(cField))

			// Handle both simple and complex match query formats
			var matchText string
			switch v := value.(type) {
			case string:
				matchText = v
			case map[string]interface{}:
				if q, ok := v["query"].(string); ok {
					matchText = q
				}
			default:
				matchText = fmt.Sprintf("%v", v)
			}

			cValue := C.CString(matchText)
			defer C.free(unsafe.Pointer(cValue))

			term := C.diagon_create_term(cField, cValue)
			defer C.diagon_free_term(term)

			diagonQuery = C.diagon_create_term_query(term)
			if diagonQuery == nil {
				errMsg := C.GoString(C.diagon_last_error())
				return nil, fmt.Errorf("failed to create match query: %s", errMsg)
			}
			break // Only support single field for now
		}
	} else if _, ok := queryObj["match_all"]; ok {
		// Match all query: {"match_all": {}}
		// Use proper MatchAllDocsQuery from Diagon C API
		s.logger.Debug(" Creating match_all query")
		diagonQuery = C.diagon_create_match_all_query()
		if diagonQuery == nil {
			errMsg := C.GoString(C.diagon_last_error())
			s.logger.Error("Failed to create match_all query", zap.String("error", errMsg))
			return nil, fmt.Errorf("failed to create match_all query: %s", errMsg)
		}
		s.logger.Debug(" match_all query created successfully")
	} else if rangeQuery, ok := queryObj["range"].(map[string]interface{}); ok {
		// Range query: {"range": {"field_name": {"gte": 100, "lte": 1000}}}
		for field, rangeParams := range rangeQuery {
			params := rangeParams.(map[string]interface{})

			s.logger.Debug(" Range query params",
				zap.String("field", field),
				zap.Any("params", params))

			var lowerValue, upperValue float64
			var includeLower, includeUpper bool

			// Parse lower bound (supports float64, date strings, and "now" expressions)
			if gte, ok := params["gte"].(float64); ok {
				lowerValue = gte
				includeLower = true
			} else if gteStr, ok := params["gte"].(string); ok {
				if epochMs, isDate := tryParseDateToEpochMs(gteStr); isDate {
					lowerValue = float64(epochMs)
					includeLower = true
				} else if strings.HasPrefix(gteStr, "now") {
					// Handle "now-7d", "now-1h" etc.
					lowerValue = float64(parseDateMathToEpochMs(gteStr))
					includeLower = true
				} else {
					lowerValue = -9007199254740992
					includeLower = true
				}
			} else if gt, ok := params["gt"].(float64); ok {
				lowerValue = gt
				includeLower = false
			} else if gtStr, ok := params["gt"].(string); ok {
				if epochMs, isDate := tryParseDateToEpochMs(gtStr); isDate {
					lowerValue = float64(epochMs)
					includeLower = false
				} else if strings.HasPrefix(gtStr, "now") {
					lowerValue = float64(parseDateMathToEpochMs(gtStr))
					includeLower = false
				} else {
					lowerValue = -9007199254740992
					includeLower = true
				}
			} else {
				lowerValue = -9007199254740992
				includeLower = true
			}

			// Parse upper bound (supports float64, date strings, and "now" expressions)
			if lte, ok := params["lte"].(float64); ok {
				upperValue = lte
				includeUpper = true
			} else if lteStr, ok := params["lte"].(string); ok {
				if epochMs, isDate := tryParseDateToEpochMs(lteStr); isDate {
					upperValue = float64(epochMs)
					includeUpper = true
				} else if strings.HasPrefix(lteStr, "now") {
					upperValue = float64(parseDateMathToEpochMs(lteStr))
					includeUpper = true
				} else {
					upperValue = 9007199254740992
					includeUpper = true
				}
			} else if lt, ok := params["lt"].(float64); ok {
				upperValue = lt
				includeUpper = false
			} else if ltStr, ok := params["lt"].(string); ok {
				if epochMs, isDate := tryParseDateToEpochMs(ltStr); isDate {
					upperValue = float64(epochMs)
					includeUpper = false
				} else if strings.HasPrefix(ltStr, "now") {
					upperValue = float64(parseDateMathToEpochMs(ltStr))
					includeUpper = false
				} else {
					upperValue = 9007199254740992
					includeUpper = true
				}
			} else {
				upperValue = 9007199254740992
				includeUpper = true
			}

			cField := C.CString(field)
			defer C.free(unsafe.Pointer(cField))

			s.logger.Debug(" Creating Diagon BKD point range query",
				zap.String("field", field),
				zap.Float64("lower", lowerValue),
				zap.Float64("upper", upperValue),
				zap.Bool("include_lower", includeLower),
				zap.Bool("include_upper", includeUpper))

			// Adjust bounds for exclusive endpoints.
			// PointRangeQuery is always inclusive [lower, upper].
			// For exclusive bounds, nudge the value by the smallest increment.
			adjLower := lowerValue
			adjUpper := upperValue
			if !includeLower {
				adjLower = nextDouble(lowerValue)
			}
			if !includeUpper {
				adjUpper = prevDouble(upperValue)
			}

			// Use BKD tree-based PointRangeQuery — O(log N) per segment.
			// Requires fields indexed with diagon_create_double_point_field.
			// This replaces the old DoubleRangeQuery which did O(N) doc-values scan.
			diagonQuery = C.diagon_create_double_point_range_query(
				cField,
				C.double(adjLower),
				C.double(adjUpper),
			)

			if diagonQuery == nil {
				errMsg := C.GoString(C.diagon_last_error())
				s.logger.Error("Failed to create BKD point range query", zap.String("error", errMsg))
				return nil, fmt.Errorf("failed to create point range query: %s", errMsg)
			}
			s.logger.Debug(" Diagon BKD point range query created successfully")
			break // Only support single field for now
		}
	} else if boolQuery, ok := queryObj["bool"].(map[string]interface{}); ok {
		// Bool query: {"bool": {"must": [...], "should": [...], "filter": [...], "must_not": [...]}}
		boolQueryBuilder := C.diagon_create_bool_query()
		if boolQueryBuilder == nil {
			errMsg := C.GoString(C.diagon_last_error())
			return nil, fmt.Errorf("failed to create bool query: %s", errMsg)
		}

		// Add MUST clauses
		if mustClauses, ok := boolQuery["must"]; ok {
			clauseArray, isArray := mustClauses.([]interface{})
			if !isArray {
				return nil, fmt.Errorf("must clauses must be an array")
			}

			for _, clause := range clauseArray {
				clauseMap, ok := clause.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("clause must be an object")
				}

				// Recursively parse sub-query
				subQuery, err := s.convertQueryToDiagon(clauseMap)
				if err != nil {
					return nil, fmt.Errorf("failed to convert must sub-query: %w", err)
				}

				C.diagon_bool_query_add_must(boolQueryBuilder, subQuery)
			}
		}

		// Add SHOULD clauses
		if shouldClauses, ok := boolQuery["should"]; ok {
			clauseArray, isArray := shouldClauses.([]interface{})
			if !isArray {
				return nil, fmt.Errorf("should clauses must be an array")
			}

			for _, clause := range clauseArray {
				clauseMap, ok := clause.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("clause must be an object")
				}

				subQuery, err := s.convertQueryToDiagon(clauseMap)
				if err != nil {
					return nil, fmt.Errorf("failed to convert should sub-query: %w", err)
				}

				C.diagon_bool_query_add_should(boolQueryBuilder, subQuery)
			}
		}

		// Add FILTER clauses
		if filterClauses, ok := boolQuery["filter"]; ok {
			clauseArray, isArray := filterClauses.([]interface{})
			if !isArray {
				return nil, fmt.Errorf("filter clauses must be an array")
			}

			for _, clause := range clauseArray {
				clauseMap, ok := clause.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("clause must be an object")
				}

				subQuery, err := s.convertQueryToDiagon(clauseMap)
				if err != nil {
					return nil, fmt.Errorf("failed to convert filter sub-query: %w", err)
				}

				C.diagon_bool_query_add_filter(boolQueryBuilder, subQuery)
			}
		}

		// Add MUST_NOT clauses
		if mustNotClauses, ok := boolQuery["must_not"]; ok {
			clauseArray, isArray := mustNotClauses.([]interface{})
			if !isArray {
				return nil, fmt.Errorf("must_not clauses must be an array")
			}

			for _, clause := range clauseArray {
				clauseMap, ok := clause.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("clause must be an object")
				}

				subQuery, err := s.convertQueryToDiagon(clauseMap)
				if err != nil {
					return nil, fmt.Errorf("failed to convert must_not sub-query: %w", err)
				}

				C.diagon_bool_query_add_must_not(boolQueryBuilder, subQuery)
			}
		}

		// Set minimum_should_match if specified
		if minShould, ok := boolQuery["minimum_should_match"].(float64); ok {
			C.diagon_bool_query_set_minimum_should_match(boolQueryBuilder, C.int(minShould))
		}

		// Build the final query
		diagonQuery = C.diagon_bool_query_build(boolQueryBuilder)
		if diagonQuery == nil {
			errMsg := C.GoString(C.diagon_last_error())
			return nil, fmt.Errorf("failed to build bool query: %s", errMsg)
		}
	} else {
		// Extract query type for better error message
		queryTypes := make([]string, 0, len(queryObj))
		for k := range queryObj {
			queryTypes = append(queryTypes, k)
		}
		return nil, fmt.Errorf("unsupported query type: %v (currently supported: 'term', 'match', 'match_all', 'range', 'bool')", queryTypes)
	}

	return diagonQuery, nil
}

// Search executes a search query using real Diagon IndexSearcher
func (s *Shard) Search(query []byte, filterExpression []byte) (*SearchResult, error) {
	return s.SearchWithLimit(query, filterExpression, 100000)
}

// SearchFieldsOnly executes a search extracting only specific stored fields (no _source parsing).
// Much faster for aggregation queries that only need a few field values per document.
func (s *Shard) SearchFieldsOnly(query []byte, filterExpression []byte, maxResults int, fields []string) (*SearchResult, error) {
	return s.searchInternal(query, filterExpression, maxResults, fields, nil)
}

// SearchWithLimit executes a search with a specified maximum number of results.
// For non-aggregation queries, use a small limit (e.g., from+size).
// For aggregation queries that need all matching docs, use a large limit.
func (s *Shard) SearchWithLimit(query []byte, filterExpression []byte, maxResults int) (*SearchResult, error) {
	return s.searchInternal(query, filterExpression, maxResults, nil, nil)
}

// SearchWithSort executes a search with sort push-down.
// sort entries are "field:asc" or "field:desc". For match_all queries with desc sort
// on a time field, reads docs from the back of the index (newest docs have highest
// internal doc IDs in chronologically-indexed data).
func (s *Shard) SearchWithSort(query []byte, filterExpression []byte, maxResults int, sort []string) (*SearchResult, error) {
	return s.searchInternal(query, filterExpression, maxResults, nil, sort)
}

// nextDouble returns the smallest double strictly greater than v.
// Used to convert exclusive lower bounds to inclusive for BKD PointRangeQuery.
func nextDouble(v float64) float64 {
	if v != v { // NaN
		return v
	}
	if v == 0 {
		return math.SmallestNonzeroFloat64
	}
	bits := math.Float64bits(v)
	if v > 0 {
		bits++
	} else {
		bits--
	}
	return math.Float64frombits(bits)
}

// prevDouble returns the largest double strictly less than v.
// Used to convert exclusive upper bounds to inclusive for BKD PointRangeQuery.
func prevDouble(v float64) float64 {
	if v != v { // NaN
		return v
	}
	if v == 0 {
		return -math.SmallestNonzeroFloat64
	}
	bits := math.Float64bits(v)
	if v > 0 {
		bits--
	} else {
		bits++
	}
	return math.Float64frombits(bits)
}

// isMatchAllQuery returns true if queryObj is a pure {"match_all": {}} with no
// other clauses (no bool wrapper, no filters).
func isMatchAllQuery(queryObj map[string]interface{}) bool {
	if len(queryObj) != 1 {
		return false
	}
	_, ok := queryObj["match_all"]
	return ok
}

// hasDescSort returns true if any sort entry contains ":desc".
func hasDescSort(sort []string) bool {
	for _, s := range sort {
		if strings.HasSuffix(s, ":desc") {
			return true
		}
	}
	return false
}

// matchAllShortcut reads N documents directly by internal doc ID, bypassing the
// C++ search entirely. For match_all queries on 116M docs this reduces latency
// from ~618ms (full TopK collection) to <2ms.
//
// When sort contains a desc field, reads from the BACK of the index (highest
// doc IDs first). For chronologically-indexed time-series data, the newest
// documents have the highest internal doc IDs, so reverse reading produces
// correct desc-sorted results without scanning all documents.
func (s *Shard) matchAllShortcut(totalStart time.Time, reopenTime, parseTime time.Duration, maxResults int, fieldsOnly []string, sort []string) (*SearchResult, error) {
	readStart := time.Now()
	s.mu.RLock()

	totalDocs := int64(C.diagon_reader_num_docs(s.reader))
	maxDocID := int(C.diagon_reader_max_doc(s.reader))
	s.mu.RUnlock()

	numToRead := maxResults
	if int64(numToRead) > totalDocs {
		numToRead = int(totalDocs)
	}
	if numToRead > maxDocID {
		numToRead = maxDocID
	}

	// Determine read direction: reverse for desc sort (newest docs = highest doc IDs)
	reverseRead := hasDescSort(sort)

	hits := make([]*Hit, 0, numToRead)

	// Pre-allocate C strings and reusable buffers
	cIDField := C.CString("_id")
	defer C.free(unsafe.Pointer(cIDField))
	idBuf := make([]byte, 1024)

	var cSourceField *C.char
	var sourceBuf []byte
	var cFieldNames []*C.char

	if len(fieldsOnly) > 0 {
		cFieldNames = make([]*C.char, len(fieldsOnly))
		for i, f := range fieldsOnly {
			cFieldNames[i] = C.CString(f)
			defer C.free(unsafe.Pointer(cFieldNames[i]))
		}
	} else {
		cSourceField = C.CString("_source")
		defer C.free(unsafe.Pointer(cSourceField))
		sourceBuf = make([]byte, 65536)
	}

	fieldBuf := make([]byte, 4096)

	s.mu.RLock()

	// Choose iteration direction based on sort
	var docIDIter func() (int, bool)
	if reverseRead {
		cur := maxDocID - 1
		docIDIter = func() (int, bool) {
			if cur < 0 || len(hits) >= numToRead {
				return 0, false
			}
			id := cur
			cur--
			return id, true
		}
	} else {
		cur := 0
		docIDIter = func() (int, bool) {
			if cur >= maxDocID || len(hits) >= numToRead {
				return 0, false
			}
			id := cur
			cur++
			return id, true
		}
	}

	for {
		docID, ok := docIDIter()
		if !ok {
			break
		}

		diagonDoc := C.diagon_reader_get_document(s.reader, C.int(docID))
		if diagonDoc == nil {
			continue // deleted or invalid doc
		}

		// Get _id
		docIDString := fmt.Sprintf("doc_%d", docID)
		if C.diagon_document_get_field_value(diagonDoc, cIDField,
			(*C.char)(unsafe.Pointer(&idBuf[0])), C.size_t(len(idBuf))) {
			for j := 0; j < len(idBuf); j++ {
				if idBuf[j] == 0 {
					docIDString = string(idBuf[:j])
					break
				}
			}
		}

		var doc map[string]interface{}
		if len(fieldsOnly) > 0 {
			doc = make(map[string]interface{}, len(fieldsOnly))
			for fi, cFN := range cFieldNames {
				if C.diagon_document_get_field_value(diagonDoc, cFN,
					(*C.char)(unsafe.Pointer(&fieldBuf[0])), C.size_t(len(fieldBuf))) {
					for j := 0; j < len(fieldBuf); j++ {
						if fieldBuf[j] == 0 {
							if j > 0 {
								doc[fieldsOnly[fi]] = string(fieldBuf[:j])
							}
							break
						}
					}
				}
			}
		} else {
			if C.diagon_document_get_field_value(diagonDoc, cSourceField,
				(*C.char)(unsafe.Pointer(&sourceBuf[0])), C.size_t(len(sourceBuf))) {
				for j := 0; j < len(sourceBuf); j++ {
					if sourceBuf[j] == 0 {
						if j > 0 {
							json.Unmarshal(sourceBuf[:j], &doc)
						}
						break
					}
				}
			}
		}

		C.diagon_free_document(diagonDoc)

		if doc == nil {
			doc = make(map[string]interface{})
		}

		hits = append(hits, &Hit{
			ID:     docIDString,
			Score:  1.0, // match_all: constant score
			Source: doc,
		})
	}
	s.mu.RUnlock()
	readTime := time.Since(readStart)

	totalTime := time.Since(totalStart)
	result := &SearchResult{
		Took:      totalTime.Milliseconds(),
		TotalHits: totalDocs,
		MaxScore:  1.0,
		Hits:      hits,
	}

	s.logger.Info("match_all shortcut: direct doc read",
		zap.Duration("reopen_reader", reopenTime),
		zap.Duration("parse_query", parseTime),
		zap.Duration("direct_read", readTime),
		zap.Duration("total", totalTime),
		zap.Int("hits_returned", len(hits)),
		zap.Int64("total_hits", totalDocs),
		zap.Bool("reverse", reverseRead))

	return result, nil
}

func (s *Shard) searchInternal(query []byte, filterExpression []byte, maxResults int, fieldsOnly []string, sort []string) (*SearchResult, error) {
	totalStart := time.Now()
	if maxResults <= 0 {
		maxResults = 100
	}

	reopenStart := time.Now()

	// Fast path (RLock): check if reader is already open. This allows
	// concurrent searches without serializing on the write lock.
	s.mu.RLock()
	needReopen := s.readerDirty || s.reader == nil
	s.mu.RUnlock()

	if needReopen {
		// Slow path (WLock): reopen reader. Double-check after acquiring
		// write lock since another goroutine may have already reopened.
		s.mu.Lock()
		needReopen = s.readerDirty || s.reader == nil
		if needReopen {
			// Only commit if there are actual dirty writes. When reader is nil
			// after restart (readerDirty=false), the index is already committed
			// from the previous run. Skipping the commit avoids a potentially
			// minutes-long waitForMerges on large indices (116M+ docs).
			if s.readerDirty && s.writer != nil {
				if !C.diagon_commit(s.writer) {
					errMsg := C.GoString(C.diagon_last_error())
					s.logger.Warn("Failed to commit before search", zap.String("error", errMsg))
					// Continue anyway - try to open reader with whatever is committed
				}
			}

			// Close existing reader/searcher
			if s.searcher != nil {
				C.diagon_free_index_searcher(s.searcher)
				s.searcher = nil
			}
			if s.reader != nil {
				C.diagon_close_index_reader(s.reader)
				s.reader = nil
			}

			// Open fresh reader
			s.reader = C.diagon_open_index_reader(s.directory)
			if s.reader == nil {
				s.mu.Unlock()
				errMsg := C.GoString(C.diagon_last_error())
				return nil, fmt.Errorf("failed to open reader: %s", errMsg)
			}

			// Create fresh searcher
			s.searcher = C.diagon_create_index_searcher(s.reader)
			if s.searcher == nil {
				s.mu.Unlock()
				errMsg := C.GoString(C.diagon_last_error())
				return nil, fmt.Errorf("failed to create searcher: %s", errMsg)
			}

			s.readerDirty = false
		}
		s.mu.Unlock()
	}

	// Hold RLock during the actual search to prevent Refresh() from
	// closing the reader/searcher while we're using them.
	s.mu.RLock()
	defer s.mu.RUnlock()

	reopenTime := time.Since(reopenStart)

	// Parse query JSON
	parseStart := time.Now()
	var queryObj map[string]interface{}
	if err := json.Unmarshal(query, &queryObj); err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}
	parseTime := time.Since(parseStart)

	// match_all shortcut: for small result sets, skip the full C++ search and
	// directly read documents by internal doc ID. diagon_search with match_all
	// iterates ALL documents (O(N) scoring) even when only a few results are
	// needed. Direct doc reads are O(maxResults) instead.
	if isMatchAllQuery(queryObj) && maxResults <= 10000 {
		return s.matchAllShortcut(totalStart, reopenTime, parseTime, maxResults, fieldsOnly, sort)
	}

	// Convert to Diagon query
	diagonQuery, err := s.convertQueryToDiagon(queryObj)
	if err != nil {
		return nil, err
	}
	defer C.diagon_free_query(diagonQuery)

	// Execute search with the specified limit
	searchStart := time.Now()
	s.mu.RLock()
	topDocs := C.diagon_search(s.searcher, diagonQuery, C.int(maxResults))
	s.mu.RUnlock()
	searchTime := time.Since(searchStart)

	if topDocs == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return nil, fmt.Errorf("search failed: %s", errMsg)
	}
	defer C.diagon_free_top_docs(topDocs)

	// Extract results
	extractStart := time.Now()
	totalHits := int64(C.diagon_top_docs_total_hits(topDocs))
	maxScore := float64(C.diagon_top_docs_max_score(topDocs))
	numResults := int(C.diagon_top_docs_score_docs_length(topDocs))

	hits := make([]*Hit, 0, numResults)

	// Pre-allocate C strings and reusable buffers for batch extraction
	cIDField := C.CString("_id")
	defer C.free(unsafe.Pointer(cIDField))
	idBuf := make([]byte, 1024)
	fieldBuf := make([]byte, 4096) // Reusable buffer for individual field reads

	// Pre-allocate C strings for fields-only mode OR _source mode
	var cSourceField *C.char
	var sourceBuf []byte
	var cFieldNames []*C.char
	if len(fieldsOnly) > 0 {
		// Fields-only mode: pre-allocate C strings for each field
		cFieldNames = make([]*C.char, len(fieldsOnly))
		for i, f := range fieldsOnly {
			cFieldNames[i] = C.CString(f)
			defer C.free(unsafe.Pointer(cFieldNames[i]))
		}
	} else {
		// Full _source mode
		cSourceField = C.CString("_source")
		defer C.free(unsafe.Pointer(cSourceField))
		sourceBuf = make([]byte, 65536)
	}

	s.mu.RLock()
	for i := 0; i < numResults; i++ {
		scoreDoc := C.diagon_top_docs_score_doc_at(topDocs, C.int(i))
		if scoreDoc == nil {
			continue
		}

		internalDocID := int(C.diagon_score_doc_get_doc(scoreDoc))
		score := float64(C.diagon_score_doc_get_score(scoreDoc))

		diagonDoc := C.diagon_reader_get_document(s.reader, C.int(internalDocID))
		if diagonDoc == nil {
			hits = append(hits, &Hit{
				ID:     fmt.Sprintf("doc_%d", internalDocID),
				Score:  score,
				Source: map[string]interface{}{"_internal_doc_id": internalDocID},
			})
			continue
		}

		// Get _id
		docIDString := fmt.Sprintf("doc_%d", internalDocID)
		if C.diagon_document_get_field_value(diagonDoc, cIDField,
			(*C.char)(unsafe.Pointer(&idBuf[0])), C.size_t(len(idBuf))) {
			for j := 0; j < len(idBuf); j++ {
				if idBuf[j] == 0 {
					docIDString = string(idBuf[:j])
					break
				}
			}
		}

		var doc map[string]interface{}
		if len(fieldsOnly) > 0 {
			// Fields-only mode: read individual stored fields (skip _source JSON parsing)
			doc = make(map[string]interface{}, len(fieldsOnly))
			for fi, cFN := range cFieldNames {
				if C.diagon_document_get_field_value(diagonDoc, cFN,
					(*C.char)(unsafe.Pointer(&fieldBuf[0])), C.size_t(len(fieldBuf))) {
					for j := 0; j < len(fieldBuf); j++ {
						if fieldBuf[j] == 0 {
							if j > 0 {
								doc[fieldsOnly[fi]] = string(fieldBuf[:j])
							}
							break
						}
					}
				}
			}
		} else {
			// Full _source mode: read and parse entire document JSON
			if C.diagon_document_get_field_value(diagonDoc, cSourceField,
				(*C.char)(unsafe.Pointer(&sourceBuf[0])), C.size_t(len(sourceBuf))) {
				for j := 0; j < len(sourceBuf); j++ {
					if sourceBuf[j] == 0 {
						if j > 0 {
							json.Unmarshal(sourceBuf[:j], &doc)
						}
						break
					}
				}
			}
		}

		C.diagon_free_document(diagonDoc)

		if doc == nil {
			doc = make(map[string]interface{})
		}

		hits = append(hits, &Hit{
			ID:     docIDString,
			Score:  score,
			Source: doc,
		})
	}
	s.mu.RUnlock()
	extractTime := time.Since(extractStart)

	totalTime := time.Since(totalStart)
	result := &SearchResult{
		Took:      totalTime.Milliseconds(),
		TotalHits: totalHits,
		MaxScore:  maxScore,
		Hits:      hits,
	}

	s.logger.Info("Diagon SearchWithLimit timing",
		zap.Duration("reopen_reader", reopenTime),
		zap.Duration("parse_query", parseTime),
		zap.Duration("cgo_search", searchTime),
		zap.Duration("extract_docs", extractTime),
		zap.Duration("total", totalTime),
		zap.Int("max_results", maxResults),
		zap.Int("num_results", numResults),
		zap.Int64("total_hits", totalHits))

	return result, nil
}

// AggFieldValue holds a single extracted field value for a document
type AggFieldValue struct {
	StringVal  string
	NumericVal float64
	IsNumeric  bool
}

// AggDocValues holds extracted field values for aggregation computation
type AggDocValues struct {
	Fields map[string]AggFieldValue
}

// SearchAndAggregate performs search then extracts only the fields needed for aggregation.
// Returns (totalHits, fieldValues-per-doc) without building full Hit objects.
// Uses batch C extraction (single CGO call for all docs) to eliminate per-doc CGO overhead.
func (s *Shard) SearchAndAggregate(query []byte, maxResults int, fields []string) (int64, []AggDocValues, error) {
	// maxResults <= 0 means "all matching docs" — resolved after reader open

	// Fast path: check under RLock if reader is already open
	s.mu.RLock()
	needReopen := s.readerDirty || s.reader == nil
	s.mu.RUnlock()

	if needReopen {
		s.mu.Lock()
		needReopen = s.readerDirty || s.reader == nil
		if needReopen {
			if s.readerDirty && s.writer != nil {
				C.diagon_commit(s.writer)
			}
			if s.searcher != nil {
				C.diagon_free_index_searcher(s.searcher)
				s.searcher = nil
			}
			if s.reader != nil {
				C.diagon_close_index_reader(s.reader)
				s.reader = nil
			}
			s.reader = C.diagon_open_index_reader(s.directory)
			if s.reader != nil {
				s.searcher = C.diagon_create_index_searcher(s.reader)
			}
			s.readerDirty = false
		}
		s.mu.Unlock()
	}

	// Resolve maxResults: 0 means "all docs"
	if maxResults <= 0 {
		s.mu.RLock()
		if s.reader != nil {
			maxResults = int(C.diagon_reader_max_doc(s.reader))
		} else {
			maxResults = 1000000 // fallback
		}
		s.mu.RUnlock()
	}

	// Parse query
	var queryObj map[string]interface{}
	if err := json.Unmarshal(query, &queryObj); err != nil {
		return 0, nil, fmt.Errorf("failed to parse query: %w", err)
	}
	diagonQuery, err := s.convertQueryToDiagon(queryObj)
	if err != nil {
		return 0, nil, err
	}
	defer C.diagon_free_query(diagonQuery)

	// Execute search — collect all matching docs for aggregation
	s.mu.RLock()
	topDocs := C.diagon_search(s.searcher, diagonQuery, C.int(maxResults))
	s.mu.RUnlock()

	if topDocs == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return 0, nil, fmt.Errorf("search failed: %s", errMsg)
	}
	defer C.diagon_free_top_docs(topDocs)

	totalHits := int64(C.diagon_top_docs_total_hits(topDocs))
	numResults := int(C.diagon_top_docs_score_docs_length(topDocs))

	if numResults == 0 || len(fields) == 0 {
		return totalHits, nil, nil
	}

	// Use batch extraction: single CGO call for all documents per field.
	// For 1 or 2 fields, use specialized batch functions.
	// This eliminates ~55K CGO round-trips (the main bottleneck).

	// Allocate buffers: ~100 bytes per value * numResults = ~5.5MB for 55K docs
	bufSize := numResults * 128 // average 128 bytes per field value
	if bufSize < 65536 {
		bufSize = 65536
	}

	docs := make([]AggDocValues, numResults)
	for i := range docs {
		docs[i].Fields = make(map[string]AggFieldValue, len(fields))
	}

	s.mu.RLock()
	if len(fields) == 1 {
		// Single field: one batch extraction call
		buf := make([]byte, bufSize)
		lengths := make([]C.int, numResults)
		cField := C.CString(fields[0])
		defer C.free(unsafe.Pointer(cField))

		C.batch_extract_field_values(
			s.reader, topDocs, cField,
			(*C.char)(unsafe.Pointer(&buf[0])), C.int(bufSize),
			&lengths[0], C.int(numResults))

		// Parse the concatenated buffer into per-doc values
		offset := 0
		for i := 0; i < numResults; i++ {
			l := int(lengths[i])
			if l > 0 && offset+l <= len(buf) {
				docs[i].Fields[fields[0]] = AggFieldValue{StringVal: string(buf[offset : offset+l])}
				offset += l + 1 // +1 for null terminator
			}
		}
	} else if len(fields) == 2 {
		// Two fields: one batch extraction call
		buf1 := make([]byte, bufSize)
		buf2 := make([]byte, bufSize)
		lengths1 := make([]C.int, numResults)
		lengths2 := make([]C.int, numResults)
		cField1 := C.CString(fields[0])
		cField2 := C.CString(fields[1])
		defer C.free(unsafe.Pointer(cField1))
		defer C.free(unsafe.Pointer(cField2))

		C.batch_extract_two_fields(
			s.reader, topDocs, cField1, cField2,
			(*C.char)(unsafe.Pointer(&buf1[0])), C.int(bufSize), &lengths1[0],
			(*C.char)(unsafe.Pointer(&buf2[0])), C.int(bufSize), &lengths2[0],
			C.int(numResults))

		offset1, offset2 := 0, 0
		for i := 0; i < numResults; i++ {
			l1 := int(lengths1[i])
			if l1 > 0 && offset1+l1 <= len(buf1) {
				docs[i].Fields[fields[0]] = AggFieldValue{StringVal: string(buf1[offset1 : offset1+l1])}
				offset1 += l1 + 1
			}
			l2 := int(lengths2[i])
			if l2 > 0 && offset2+l2 <= len(buf2) {
				docs[i].Fields[fields[1]] = AggFieldValue{StringVal: string(buf2[offset2 : offset2+l2])}
				offset2 += l2 + 1
			}
		}
	} else {
		// 3+ fields: single pass - load each document once, extract all fields.
		// Allocate pointer arrays in C memory to satisfy CGO pointer rules.
		nFields := len(fields)
		ptrSize := C.size_t(unsafe.Sizeof((*C.char)(nil)))

		// C-allocated arrays of pointers
		cFieldNamesArr := (*[16]*C.char)(C.malloc(C.size_t(nFields) * ptrSize))
		cBufsArr := (*[16]*C.char)(C.malloc(C.size_t(nFields) * ptrSize))
		cBufSizesArr := (*[16]C.int)(C.malloc(C.size_t(nFields) * C.size_t(unsafe.Sizeof(C.int(0)))))
		cLengthPtrsArr := (*[16]*C.int)(C.malloc(C.size_t(nFields) * ptrSize))

		// Go-side slices for reading results back
		bufs := make([][]byte, nFields)
		lengthArrays := make([][]C.int, nFields)

		for i, field := range fields {
			cFieldNamesArr[i] = C.CString(field)
			bufs[i] = make([]byte, bufSize)
			cBufsArr[i] = (*C.char)(unsafe.Pointer(&bufs[i][0]))
			cBufSizesArr[i] = C.int(bufSize)
			lengthArrays[i] = make([]C.int, numResults)
			cLengthPtrsArr[i] = &lengthArrays[i][0]
		}

		C.batch_extract_multi_fields(
			s.reader, topDocs,
			&cFieldNamesArr[0],
			&cBufsArr[0],
			&cBufSizesArr[0],
			&cLengthPtrsArr[0],
			C.int(nFields),
			C.int(numResults))

		for i := 0; i < nFields; i++ {
			C.free(unsafe.Pointer(cFieldNamesArr[i]))
		}
		C.free(unsafe.Pointer(cFieldNamesArr))
		C.free(unsafe.Pointer(cBufsArr))
		C.free(unsafe.Pointer(cBufSizesArr))
		C.free(unsafe.Pointer(cLengthPtrsArr))

		// Parse concatenated buffers into per-doc values
		for fi, field := range fields {
			offset := 0
			for i := 0; i < numResults; i++ {
				l := int(lengthArrays[fi][i])
				if l > 0 && offset+l <= len(bufs[fi]) {
					docs[i].Fields[field] = AggFieldValue{StringVal: string(bufs[fi][offset : offset+l])}
					offset += l + 1
				}
			}
		}
	}
	s.mu.RUnlock()

	return totalHits, docs, nil
}

// DirectAggScan iterates documents by internal doc ID and extracts field values
// for aggregation. Uses batch C extraction (1 CGO call per field instead of
// 3 CGO calls per doc per field) and an in-memory column cache to avoid
// re-scanning on repeated queries. For match_all agg-only queries.
func (s *Shard) DirectAggScan(fields []string, maxDocs int) (int64, []AggDocValues, error) {
	// maxDocs <= 0 means "scan all docs in shard" — resolved after reader open

	// Ensure reader is open
	s.mu.RLock()
	needReopen := s.readerDirty || s.reader == nil
	s.mu.RUnlock()

	if needReopen {
		s.mu.Lock()
		needReopen = s.readerDirty || s.reader == nil
		if needReopen {
			if s.readerDirty && s.writer != nil {
				C.diagon_commit(s.writer)
			}
			if s.searcher != nil {
				C.diagon_free_index_searcher(s.searcher)
				s.searcher = nil
			}
			if s.reader != nil {
				C.diagon_close_index_reader(s.reader)
				s.reader = nil
			}
			s.reader = C.diagon_open_index_reader(s.directory)
			if s.reader != nil {
				s.searcher = C.diagon_create_index_searcher(s.reader)
			}
			s.readerDirty = false
			// Invalidate column cache on reader reopen
			s.columnCacheMu.Lock()
			s.columnCache = nil
			s.columnCacheN = 0
			s.columnCacheMu.Unlock()
		}
		s.mu.Unlock()
	}

	s.mu.RLock()
	if s.reader == nil {
		s.mu.RUnlock()
		return 0, nil, fmt.Errorf("reader not initialized")
	}
	totalDocs := int64(C.diagon_reader_num_docs(s.reader))
	maxDocID := int(C.diagon_reader_max_doc(s.reader))
	s.mu.RUnlock()

	// maxDocs <= 0 means "all docs" — resolve to actual doc count
	if maxDocs <= 0 || maxDocs > maxDocID {
		maxDocs = maxDocID
	}

	if len(fields) == 0 {
		return totalDocs, nil, nil
	}

	// Check column cache for all requested fields
	s.columnCacheMu.RLock()
	allCached := s.columnCache != nil && s.columnCacheN > 0
	if allCached {
		for _, f := range fields {
			if _, ok := s.columnCache[f]; !ok {
				allCached = false
				break
			}
		}
	}
	if allCached {
		// All fields cached — build AggDocValues from pure Go arrays (no CGO)
		n := s.columnCacheN
		if n > maxDocs {
			n = maxDocs
		}
		docs := make([]AggDocValues, n)
		for i := 0; i < n; i++ {
			docs[i].Fields = make(map[string]AggFieldValue, len(fields))
			for _, f := range fields {
				v := s.columnCache[f][i]
				if v != "" {
					docs[i].Fields[f] = AggFieldValue{StringVal: v}
				}
			}
		}
		s.columnCacheMu.RUnlock()
		return totalDocs, docs, nil
	}
	s.columnCacheMu.RUnlock()

	// Chunked batch extraction: scan docs in chunks to avoid allocating huge buffers.
	// Each chunk uses a fixed-size buffer (~64MB), then results are accumulated into
	// pre-allocated column slices. For 116M docs this uses ~3GB for column cache
	// (24 bytes/timestamp * 116M) vs ~15GB for a single monolithic buffer.
	const chunkSize = 500000 // docs per chunk
	chunkBufSize := chunkSize * 128
	if chunkBufSize < 65536 {
		chunkBufSize = 65536
	}

	// Pre-allocate column slices to full capacity
	columnData := make(map[string][]string, len(fields))
	for _, f := range fields {
		columnData[f] = make([]string, 0, maxDocs)
	}

	s.mu.RLock()
	for _, field := range fields {
		buf := make([]byte, chunkBufSize)
		lengths := make([]C.int, chunkSize)
		cField := C.CString(field)
		nextDocID := C.int(0)

		for nextDocID >= 0 {
			var cNextDocID C.int
			n := int(C.batch_scan_field_values_from(
				s.reader, cField,
				(*C.char)(unsafe.Pointer(&buf[0])), C.int(chunkBufSize),
				&lengths[0], C.int(chunkSize),
				nextDocID, &cNextDocID))

			// Parse chunk into column slice
			offset := 0
			for i := 0; i < n; i++ {
				l := int(lengths[i])
				if l > 0 && offset+l <= len(buf) {
					columnData[field] = append(columnData[field], string(buf[offset:offset+l]))
					offset += l + 1
				} else {
					columnData[field] = append(columnData[field], "")
				}
			}
			nextDocID = cNextDocID
		}

		C.free(unsafe.Pointer(cField))
	}
	s.mu.RUnlock()

	// Determine actual doc count from first field
	numDocs := 0
	for _, vals := range columnData {
		numDocs = len(vals)
		break
	}

	// Store in column cache
	s.columnCacheMu.Lock()
	if s.columnCache == nil {
		s.columnCache = make(map[string][]string, len(fields))
	}
	for f, vals := range columnData {
		s.columnCache[f] = vals
	}
	s.columnCacheN = numDocs
	s.columnCacheMu.Unlock()

	// Build AggDocValues from extracted columns
	n := numDocs
	if n > maxDocs {
		n = maxDocs
	}
	docs := make([]AggDocValues, n)
	for i := 0; i < n; i++ {
		docs[i].Fields = make(map[string]AggFieldValue, len(fields))
		for _, f := range fields {
			if i < len(columnData[f]) && columnData[f][i] != "" {
				docs[i].Fields[f] = AggFieldValue{StringVal: columnData[f][i]}
			}
		}
	}

	return totalDocs, docs, nil
}

// DirectAggColumns returns raw column data (map[field][]string) and total doc count
// for computing aggregations directly without building per-doc AggDocValues structs.
// This reduces memory from ~90GB (AggDocValues for 116M docs) to ~3GB (column cache).
func (s *Shard) DirectAggColumns(fields []string) (int64, map[string][]string, int, error) {
	// Ensure reader is open (same as DirectAggScan)
	s.mu.RLock()
	needReopen := s.readerDirty || s.reader == nil
	s.mu.RUnlock()

	if needReopen {
		s.mu.Lock()
		needReopen = s.readerDirty || s.reader == nil
		if needReopen {
			if s.readerDirty && s.writer != nil {
				C.diagon_commit(s.writer)
			}
			if s.searcher != nil {
				C.diagon_free_index_searcher(s.searcher)
				s.searcher = nil
			}
			if s.reader != nil {
				C.diagon_close_index_reader(s.reader)
				s.reader = nil
			}
			s.reader = C.diagon_open_index_reader(s.directory)
			if s.reader != nil {
				s.searcher = C.diagon_create_index_searcher(s.reader)
			}
			s.readerDirty = false
			s.columnCacheMu.Lock()
			s.columnCache = nil
			s.columnCacheN = 0
			s.columnCacheMu.Unlock()
		}
		s.mu.Unlock()
	}

	s.mu.RLock()
	if s.reader == nil {
		s.mu.RUnlock()
		return 0, nil, 0, fmt.Errorf("reader not initialized")
	}
	totalDocs := int64(C.diagon_reader_num_docs(s.reader))
	maxDocID := int(C.diagon_reader_max_doc(s.reader))
	s.mu.RUnlock()

	if len(fields) == 0 {
		return totalDocs, nil, 0, nil
	}

	// Check column cache
	s.columnCacheMu.RLock()
	allCached := s.columnCache != nil && s.columnCacheN > 0
	if allCached {
		for _, f := range fields {
			if _, ok := s.columnCache[f]; !ok {
				allCached = false
				break
			}
		}
	}
	if allCached {
		columnData := make(map[string][]string, len(fields))
		for _, f := range fields {
			columnData[f] = s.columnCache[f]
		}
		n := s.columnCacheN
		s.columnCacheMu.RUnlock()
		return totalDocs, columnData, n, nil
	}
	s.columnCacheMu.RUnlock()

	// Chunked batch extraction (same as DirectAggScan)
	const chunkSize = 500000
	chunkBufSize := chunkSize * 128

	columnData := make(map[string][]string, len(fields))
	for _, f := range fields {
		columnData[f] = make([]string, 0, maxDocID)
	}

	s.mu.RLock()
	for _, field := range fields {
		buf := make([]byte, chunkBufSize)
		lengths := make([]C.int, chunkSize)
		cField := C.CString(field)
		nextDocID := C.int(0)

		for nextDocID >= 0 {
			var cNextDocID C.int
			n := int(C.batch_scan_field_values_from(
				s.reader, cField,
				(*C.char)(unsafe.Pointer(&buf[0])), C.int(chunkBufSize),
				&lengths[0], C.int(chunkSize),
				nextDocID, &cNextDocID))

			offset := 0
			for i := 0; i < n; i++ {
				l := int(lengths[i])
				if l > 0 && offset+l <= len(buf) {
					columnData[field] = append(columnData[field], string(buf[offset:offset+l]))
					offset += l + 1
				} else {
					columnData[field] = append(columnData[field], "")
				}
			}
			nextDocID = cNextDocID
		}
		C.free(unsafe.Pointer(cField))
	}
	s.mu.RUnlock()

	numDocs := 0
	for _, vals := range columnData {
		numDocs = len(vals)
		break
	}

	// Store in column cache
	s.columnCacheMu.Lock()
	if s.columnCache == nil {
		s.columnCache = make(map[string][]string, len(fields))
	}
	for f, vals := range columnData {
		s.columnCache[f] = vals
	}
	s.columnCacheN = numDocs
	s.columnCacheMu.Unlock()

	return totalDocs, columnData, numDocs, nil
}

// DocCount returns the number of documents in the shard index.
func (s *Shard) DocCount() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.reader != nil {
		return int64(C.diagon_reader_num_docs(s.reader))
	}
	return 0
}

// TermBucket represents a single bucket in a terms aggregation result
type TermBucket struct {
	Key      string
	DocCount int64
}

// ComputeTermsAgg computes a terms aggregation in C by scanning documents
// and building a hash map of field values → counts. Capped at 500K docs to
// avoid multi-minute scans on large indices; results are approximate above that.
// Results are cached until the reader is refreshed.
func (s *Shard) ComputeTermsAgg(field string, size int) ([]TermBucket, error) {
	if size <= 0 {
		size = 10
	}

	cacheKey := field + ":" + fmt.Sprintf("%d", size)

	// Check cache first
	s.termsAggCacheMu.RLock()
	if cached, ok := s.termsAggCache[cacheKey]; ok {
		s.termsAggCacheMu.RUnlock()
		return cached, nil
	}
	s.termsAggCacheMu.RUnlock()

	// Fast path: check under RLock if reader is already open
	s.mu.RLock()
	needReopen := s.reader == nil || s.readerDirty
	s.mu.RUnlock()

	if needReopen {
		s.mu.Lock()
		needReopen = s.reader == nil || s.readerDirty
		if needReopen {
			// Invalidate caches on reader reopen
			s.termsAggCacheMu.Lock()
			s.termsAggCache = nil
			s.termsAggCacheMu.Unlock()
			s.columnCacheMu.Lock()
			s.columnCache = nil
			s.columnCacheN = 0
			s.columnCacheMu.Unlock()

			if s.readerDirty && s.writer != nil {
				C.diagon_commit(s.writer)
			}
			if s.searcher != nil {
				C.diagon_free_index_searcher(s.searcher)
				s.searcher = nil
			}
			if s.reader != nil {
				C.diagon_close_index_reader(s.reader)
			}
			s.reader = C.diagon_open_index_reader(s.directory)
			if s.reader != nil {
				s.searcher = C.diagon_create_index_searcher(s.reader)
			}
			s.readerDirty = false
		}
		s.mu.Unlock()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.reader == nil {
		return nil, fmt.Errorf("reader not initialized")
	}

	cField := C.CString(field)
	defer C.free(unsafe.Pointer(cField))

	// Allocate output buckets
	outBuckets := make([]C.TermBucketC, size)

	// Cap scan at 50K docs. On cold disk (no OS file cache), stored field reads
	// are ~240µs/doc, so 50K docs = ~12s. At 7µs/doc warm = 0.35s.
	// 50K samples is sufficient for accurate term frequency estimation
	// (Big5 has ~100 unique log_stream values → ~500 samples per term).
	const maxDocsToScan = 50000

	n := C.compute_terms_agg_stored(
		s.reader,
		cField,
		&outBuckets[0],
		C.int(size),
		C.int(maxDocsToScan),
	)

	if n <= 0 {
		return nil, nil
	}

	buckets := make([]TermBucket, 0, int(n))
	for i := 0; i < int(n); i++ {
		key := C.GoString(&outBuckets[i].key[0])
		// Sanitize non-UTF-8 strings to prevent gRPC protobuf marshaling failures.
		// Stored field values may contain arbitrary bytes from the C++ engine.
		if !utf8.ValidString(key) {
			continue
		}
		buckets = append(buckets, TermBucket{
			Key:      key,
			DocCount: int64(outBuckets[i].doc_count),
		})
	}

	// Cache the result
	s.termsAggCacheMu.Lock()
	if s.termsAggCache == nil {
		s.termsAggCache = make(map[string][]TermBucket)
	}
	s.termsAggCache[cacheKey] = buckets
	s.termsAggCacheMu.Unlock()

	return buckets, nil
}

// ComputeCardinality returns the approximate number of unique values for a field.
// Uses sampled reading (50K docs uniformly spaced) to avoid full-index scan on large indices.
// For low-cardinality fields, result is exact. For high-cardinality, approximate.
func (s *Shard) ComputeCardinality(field string) (int64, error) {
	// Ensure reader is open (RLock fast path)
	s.mu.RLock()
	needReopen := s.reader == nil || s.readerDirty
	s.mu.RUnlock()

	if needReopen {
		s.mu.Lock()
		needReopen = s.reader == nil || s.readerDirty
		if needReopen {
			if s.readerDirty && s.writer != nil {
				C.diagon_commit(s.writer)
			}
			if s.searcher != nil {
				C.diagon_free_index_searcher(s.searcher)
				s.searcher = nil
			}
			if s.reader != nil {
				C.diagon_close_index_reader(s.reader)
			}
			s.reader = C.diagon_open_index_reader(s.directory)
			if s.reader != nil {
				s.searcher = C.diagon_create_index_searcher(s.reader)
			}
			s.readerDirty = false
		}
		s.mu.Unlock()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.reader == nil {
		return 0, fmt.Errorf("reader not initialized")
	}

	cField := C.CString(field)
	defer C.free(unsafe.Pointer(cField))

	result := C.compute_cardinality_sampled(s.reader, cField, C.int(5000))
	return int64(result), nil
}

// getDocumentByInternalID retrieves a document's stored fields given its internal Diagon doc ID
// Returns the document fields map and the document's _id string
func (s *Shard) getDocumentByInternalID(internalDocID int) (map[string]interface{}, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	maxDoc := int(C.diagon_reader_max_doc(s.reader))
	if internalDocID >= maxDoc {
		return nil, "", fmt.Errorf("internal docID %d >= maxDoc %d", internalDocID, maxDoc)
	}

	diagonDoc := C.diagon_reader_get_document(s.reader, C.int(internalDocID))
	if diagonDoc == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return nil, "", fmt.Errorf("failed to retrieve document: %s", errMsg)
	}
	defer C.diagon_free_document(diagonDoc)

	var docIDString string

	// Get _id field
	idBuf := make([]byte, 1024)
	cIDFieldName := C.CString("_id")
	defer C.free(unsafe.Pointer(cIDFieldName))
	if C.diagon_document_get_field_value(diagonDoc, cIDFieldName,
		(*C.char)(unsafe.Pointer(&idBuf[0])), C.size_t(len(idBuf))) {
		nullIdx := 0
		for i, b := range idBuf {
			if b == 0 {
				nullIdx = i
				break
			}
		}
		docIDString = string(idBuf[:nullIdx])
	} else {
		docIDString = fmt.Sprintf("doc_%d", internalDocID)
	}

	// Retrieve _source field (full document JSON stored during indexing)
	sourceBuf := make([]byte, 65536) // 64KB buffer
	cSourceFieldName := C.CString("_source")
	defer C.free(unsafe.Pointer(cSourceFieldName))
	if C.diagon_document_get_field_value(diagonDoc, cSourceFieldName,
		(*C.char)(unsafe.Pointer(&sourceBuf[0])), C.size_t(len(sourceBuf))) {
		nullIdx := 0
		for i, b := range sourceBuf {
			if b == 0 {
				nullIdx = i
				break
			}
		}
		if nullIdx > 0 {
			var doc map[string]interface{}
			if err := json.Unmarshal(sourceBuf[:nullIdx], &doc); err == nil {
				return doc, docIDString, nil
			}
		}
	}

	// Fallback: return minimal document
	return map[string]interface{}{"_id": docIDString}, docIDString, nil
}

// GetDocument retrieves a document by ID
func (s *Shard) GetDocument(docID string) (map[string]interface{}, error) {
	s.logger.Debug("GetDocument", zap.String("doc_id", docID))
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debug("GetDocument called", zap.String("doc_id", docID))

	// Ensure reader and searcher are initialized
	if s.reader == nil || s.searcher == nil {
		s.logger.Debug("Reader not initialized, opening now", zap.String("doc_id", docID))

		// Commit first to ensure changes are visible
		if !C.diagon_commit(s.writer) {
			errMsg := C.GoString(C.diagon_last_error())
			return nil, fmt.Errorf("commit failed: %s", errMsg)
		}

		// Open reader
		s.reader = C.diagon_open_index_reader(s.directory)
		if s.reader == nil {
			errMsg := C.GoString(C.diagon_last_error())
			return nil, fmt.Errorf("failed to open reader: %s", errMsg)
		}

		// Create searcher
		s.searcher = C.diagon_create_index_searcher(s.reader)
		if s.searcher == nil {
			errMsg := C.GoString(C.diagon_last_error())
			return nil, fmt.Errorf("failed to create searcher: %s", errMsg)
		}

		s.logger.Debug("Reader and searcher initialized successfully")
	}

	// Search for the document by _id field to get internal doc ID
	s.logger.Debug("STEP1: Creating term for _id search")
	cIDField := C.CString("_id")
	defer C.free(unsafe.Pointer(cIDField))

	cDocID := C.CString(docID)
	defer C.free(unsafe.Pointer(cDocID))

	term := C.diagon_create_term(cIDField, cDocID)
	if term == nil {
		errMsg := C.GoString(C.diagon_last_error())
		s.logger.Error("FAILED at create term", zap.String("error", errMsg))
		return nil, fmt.Errorf("failed to create term: %s", errMsg)
	}
	defer C.diagon_free_term(term)

	s.logger.Debug("STEP2: Creating term query")
	query := C.diagon_create_term_query(term)
	if query == nil {
		errMsg := C.GoString(C.diagon_last_error())
		s.logger.Error("FAILED at create query", zap.String("error", errMsg))
		return nil, fmt.Errorf("failed to create query: %s", errMsg)
	}
	defer C.diagon_free_query(query)

	s.logger.Debug("STEP3: Executing search", zap.String("doc_id", docID))

	// Search to find the internal doc ID
	topDocs := C.diagon_search(s.searcher, query, 1)
	if topDocs == nil {
		errMsg := C.GoString(C.diagon_last_error())
		s.logger.Error("FAILED at search", zap.String("error", errMsg))
		return nil, fmt.Errorf("search failed: %s", errMsg)
	}
	defer C.diagon_free_top_docs(topDocs)

	totalHits := int64(C.diagon_top_docs_total_hits(topDocs))
	s.logger.Debug("Search completed", zap.Int64("total_hits", totalHits))

	if totalHits == 0 {
		return nil, fmt.Errorf("document not found")
	}

	// Get internal doc ID from search result
	scoreDoc := C.diagon_top_docs_score_doc_at(topDocs, 0)
	if scoreDoc == nil {
		return nil, fmt.Errorf("failed to get score doc")
	}

	internalDocID := int(C.diagon_score_doc_get_doc(scoreDoc))
	s.logger.Debug("Found document", zap.Int("internal_doc_id", internalDocID))

	// Retrieve stored fields using reader
	diagonDoc := C.diagon_reader_get_document(s.reader, C.int(internalDocID))
	if diagonDoc == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return nil, fmt.Errorf("failed to retrieve document: %s", errMsg)
	}
	defer C.diagon_free_document(diagonDoc)

	// Retrieve _source field (full document JSON)
	sourceBuf := make([]byte, 65536)
	cSourceFieldName := C.CString("_source")
	defer C.free(unsafe.Pointer(cSourceFieldName))
	if C.diagon_document_get_field_value(diagonDoc, cSourceFieldName,
		(*C.char)(unsafe.Pointer(&sourceBuf[0])), C.size_t(len(sourceBuf))) {
		nullIdx := 0
		for i, b := range sourceBuf {
			if b == 0 {
				nullIdx = i
				break
			}
		}
		if nullIdx > 0 {
			var doc map[string]interface{}
			if err := json.Unmarshal(sourceBuf[:nullIdx], &doc); err == nil {
				return doc, nil
			}
		}
	}

	// Fallback
	return map[string]interface{}{"_id": docID}, nil
}

// DeleteDocument deletes a document (not yet implemented in Phase 4)
func (s *Shard) DeleteDocument(docID string) error {
	// TODO: Implement when document deletion is available in Diagon
	s.logger.Warn("Document deletion not yet implemented in Diagon Phase 4", zap.String("doc_id", docID))
	return fmt.Errorf("document deletion not yet implemented in Diagon Phase 4")
}

// ForceMerge optimizes the index by cascading forced merges with progressively
// lower segment targets. Each step uses diagon_force_merge(N) which creates groups
// of ~(current/N) segments and merges each group separately via SegmentMerger.
// Key insight: SegmentMerger's MergedTermsEnum uses O(N) linear scan per term,
// so keeping group sizes small (~10) is critical. Cascading 10x reductions
// (10000→1000→100→10→1) keeps each merge manageable.
func (s *Shard) ForceMerge(maxSegments int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.writer == nil {
		return fmt.Errorf("writer not initialized")
	}

	// NOTE: Do NOT call cleanupMergedFiles here! After a successful force merge,
	// committed segments may have _merged_* prefixes. Deleting them destroys data.
	// Diagon's IndexWriter handles cleanup of unreferenced files on open.

	currentSegs := countCFSFiles(s.path)
	s.logger.Info("Starting cascading force merge",
		zap.Int("current_segments", currentSegs),
		zap.Int("target_segments", maxSegments))

	if currentSegs <= maxSegments {
		s.logger.Info("Already at target segment count")
		return nil
	}

	// Build cascade targets: 10x reduction each step
	// e.g., 3372 segments → targets: [337, 33, 3, 1]
	var targets []int
	target := currentSegs / 10
	for target > maxSegments {
		targets = append(targets, target)
		target = target / 10
	}
	targets = append(targets, maxSegments)

	for step, t := range targets {
		segs := countCFSFiles(s.path)
		if segs <= t {
			s.logger.Info("Step skipped, already at target",
				zap.Int("step", step+1),
				zap.Int("segments", segs),
				zap.Int("target", t))
			continue
		}

		groupSize := segs / t
		if groupSize < 2 {
			groupSize = 2
		}

		s.logger.Info("Force merge step starting",
			zap.Int("step", step+1),
			zap.Int("total_steps", len(targets)),
			zap.Int("from_segments", segs),
			zap.Int("to_segments", t),
			zap.Int("group_size", groupSize))

		startTime := time.Now()

		if !C.diagon_force_merge(s.writer, C.int(t)) {
			errMsg := C.GoString(C.diagon_last_error())
			return fmt.Errorf("force merge step %d (target=%d) failed: %s", step+1, t, errMsg)
		}

		afterSegs := countCFSFiles(s.path)
		elapsed := time.Since(startTime)
		s.logger.Info("Force merge step completed",
			zap.Int("step", step+1),
			zap.Int("segments_before", segs),
			zap.Int("segments_after", afterSegs),
			zap.Duration("elapsed", elapsed))
	}

	finalSegs := countCFSFiles(s.path)
	s.logger.Info("Cascading force merge completed",
		zap.Int("final_segments", finalSegs))
	return nil
}

// ShardStats contains statistics about the shard
type ShardStats struct {
	NumDocs      int64 `json:"num_docs"`       // Number of documents
	MaxDoc       int64 `json:"max_doc"`        // Maximum document ID
	SegmentCount int   `json:"segment_count"`  // Number of segments
	SizeBytes    int64 `json:"size_bytes"`     // Index size in bytes
}

// GetStats returns statistics about the shard
func (s *Shard) GetStats() (*ShardStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.directory == nil {
		return nil, fmt.Errorf("directory not initialized")
	}

	// Open a fresh reader to get current stats
	reader := C.diagon_open_index_reader(s.directory)
	if reader == nil {
		errMsg := C.GoString(C.diagon_last_error())
		return nil, fmt.Errorf("failed to open reader for stats: %s", errMsg)
	}
	defer C.diagon_close_index_reader(reader)

	stats := &ShardStats{
		NumDocs:      int64(C.diagon_reader_num_docs(reader)),
		MaxDoc:       int64(C.diagon_reader_max_doc(reader)),
		SegmentCount: int(C.diagon_reader_get_segment_count(reader)),
		SizeBytes:    int64(C.diagon_directory_get_size(s.directory)),
	}

	s.logger.Debug("Retrieved shard stats",
		zap.Int64("num_docs", stats.NumDocs),
		zap.Int64("max_doc", stats.MaxDoc),
		zap.Int("segment_count", stats.SegmentCount),
		zap.Int64("size_bytes", stats.SizeBytes))

	return stats, nil
}

// Close closes the shard and frees all resources
func (s *Shard) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close searcher
	if s.searcher != nil {
		C.diagon_free_index_searcher(s.searcher)
		s.searcher = nil
	}

	// Close reader
	if s.reader != nil {
		C.diagon_close_index_reader(s.reader)
		s.reader = nil
	}

	// Close writer
	if s.writer != nil {
		C.diagon_close_index_writer(s.writer)
		s.writer = nil
	}

	// Close directory
	if s.directory != nil {
		C.diagon_close_directory(s.directory)
		s.directory = nil
	}

	s.logger.Info("Closed real Diagon shard")

	return nil
}

// SearchResult represents search results
type SearchResult struct {
	Took         int64                        `json:"took"`
	TotalHits    int64                        `json:"total_hits"`
	MaxScore     float64                      `json:"max_score"`
	Hits         []*Hit                       `json:"hits"`
	Aggregations map[string]AggregationResult `json:"aggregations,omitempty"`
}

// Hit represents a search hit
type Hit struct {
	ID     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

// AggregationResult represents an aggregation result
type AggregationResult struct {
	Type    string                   `json:"type"`
	Buckets []map[string]interface{} `json:"buckets,omitempty"`
	Count   int64                    `json:"count,omitempty"`
	Min     float64                  `json:"min,omitempty"`
	Max     float64                  `json:"max,omitempty"`
	Avg     float64                  `json:"avg,omitempty"`
	Sum     float64                  `json:"sum,omitempty"`
	Value   int64                    `json:"value,omitempty"`
	Values  map[string]float64       `json:"values,omitempty"`
}

// ---------- BKD-based aggregation Go wrappers ----------

// DateHistogramBucket holds one bucket from BKD-based date histogram.
type DateHistogramBucket struct {
	KeyMs    float64 // epoch ms of bucket start
	DocCount int64
}

// RangeSpec defines a single range for range aggregation.
type RangeSpec struct {
	From *float64
	To   *float64
	Key  string // optional user-provided key
}

// RangeAggBucket holds one bucket from BKD-based range aggregation.
type RangeAggBucket struct {
	From     *float64
	To       *float64
	Key      string
	DocCount int64
}

// ensureReaderOpen ensures the index reader and searcher are open and current.
// Returns true if reader is ready, false on error.
func (s *Shard) ensureReaderOpen() bool {
	s.mu.RLock()
	needReopen := s.reader == nil || s.readerDirty
	s.mu.RUnlock()

	if needReopen {
		s.mu.Lock()
		needReopen = s.reader == nil || s.readerDirty
		if needReopen {
			if s.readerDirty && s.writer != nil {
				C.diagon_commit(s.writer)
			}
			if s.searcher != nil {
				C.diagon_free_index_searcher(s.searcher)
				s.searcher = nil
			}
			if s.reader != nil {
				C.diagon_close_index_reader(s.reader)
				s.reader = nil
			}
			s.reader = C.diagon_open_index_reader(s.directory)
			if s.reader != nil {
				s.searcher = C.diagon_create_index_searcher(s.reader)
			}
			s.readerDirty = false
			// Invalidate caches
			s.termsAggCacheMu.Lock()
			s.termsAggCache = nil
			s.termsAggCacheMu.Unlock()
			s.columnCacheMu.Lock()
			s.columnCache = nil
			s.columnCacheN = 0
			s.columnCacheMu.Unlock()
		}
		s.mu.Unlock()
	}

	s.mu.RLock()
	ready := s.reader != nil && s.searcher != nil
	s.mu.RUnlock()
	return ready
}

// ConvertFilterQuery converts a JSON query object to a Diagon query handle.
// The returned handle must be freed by the caller with FreeQuery.
func (s *Shard) ConvertFilterQuery(queryObj map[string]interface{}) (unsafe.Pointer, error) {
	q, err := s.convertQueryToDiagon(queryObj)
	if err != nil {
		return nil, err
	}
	return unsafe.Pointer(q), nil
}

// FreeQuery frees a Diagon query handle returned by ConvertFilterQuery.
func FreeQuery(q unsafe.Pointer) {
	if q != nil {
		C.diagon_free_query(C.DiagonQuery(q))
	}
}

// ComputeDateHistogramBKD computes a date histogram using single-pass BKD tree traversal.
// Uses diagon_compute_histogram C API: O(N) once through BKD tree, bucketing all points.
// For unfiltered (match_all) queries, this replaces 10K+ individual searches with 1 call.
// filterQuery is currently ignored — filtered histograms should use SearchAndAggregate.
func (s *Shard) ComputeDateHistogramBKD(field string, intervalMs float64, minEpochMs float64, maxEpochMs float64, filterQuery map[string]interface{}) ([]DateHistogramBucket, error) {
	if !s.ensureReaderOpen() {
		return nil, fmt.Errorf("reader not initialized")
	}

	numBuckets := int((maxEpochMs-minEpochMs)/intervalMs) + 1
	if numBuckets <= 0 {
		return nil, nil
	}
	if numBuckets > 1000000 {
		numBuckets = 1000000
	}

	// For filtered histograms, return nil to let caller use SearchAndAggregate
	if filterQuery != nil {
		return nil, nil
	}

	cField := C.CString(field)
	defer C.free(unsafe.Pointer(cField))

	bucketCounts := make([]C.int64_t, numBuckets)

	s.mu.RLock()
	totalCounted := C.diagon_compute_histogram(
		s.reader,
		cField,
		C.double(minEpochMs),
		C.double(intervalMs),
		C.int(numBuckets),
		&bucketCounts[0],
	)
	s.mu.RUnlock()

	if int64(totalCounted) < 0 {
		return nil, fmt.Errorf("histogram computation failed")
	}

	// Collect non-empty buckets
	var result []DateHistogramBucket
	for i := 0; i < numBuckets; i++ {
		count := int64(bucketCounts[i])
		if count > 0 {
			result = append(result, DateHistogramBucket{
				KeyMs:    minEpochMs + float64(i)*intervalMs,
				DocCount: count,
			})
		}
	}
	return result, nil
}

// ComputeRangeAggBKD computes a range aggregation using BKD range queries.
// Each range is counted via a single BKD range query — O(log N) per range.
func (s *Shard) ComputeRangeAggBKD(field string, ranges []RangeSpec, filterQuery map[string]interface{}) ([]RangeAggBucket, error) {
	if !s.ensureReaderOpen() {
		return nil, fmt.Errorf("reader not initialized")
	}
	if len(ranges) == 0 {
		return nil, nil
	}

	var filterQ C.DiagonQuery
	if filterQuery != nil {
		var err error
		filterQ, err = s.convertQueryToDiagon(filterQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to convert filter query: %w", err)
		}
		defer C.diagon_free_query(filterQ)
	}

	cField := C.CString(field)
	defer C.free(unsafe.Pointer(cField))

	cRanges := make([]C.RangeAggBucketC, len(ranges))
	for i, r := range ranges {
		if r.From != nil {
			cRanges[i].from_val = C.double(*r.From)
			cRanges[i].has_from = 1
		}
		if r.To != nil {
			cRanges[i].to_val = C.double(*r.To)
			cRanges[i].has_to = 1
		}
	}

	s.mu.RLock()
	C.compute_range_agg_bkd(
		s.searcher,
		cField,
		&cRanges[0],
		C.int(len(ranges)),
		filterQ,
	)
	s.mu.RUnlock()

	result := make([]RangeAggBucket, len(ranges))
	for i, r := range ranges {
		result[i] = RangeAggBucket{
			From:     r.From,
			To:       r.To,
			Key:      r.Key,
			DocCount: int64(cRanges[i].doc_count),
		}
	}
	return result, nil
}

// CountQuery counts total hits for a query using BKD. Returns 0 on error.
func (s *Shard) CountQuery(queryObj map[string]interface{}) int64 {
	if !s.ensureReaderOpen() {
		return 0
	}
	q, err := s.convertQueryToDiagon(queryObj)
	if err != nil {
		return 0
	}
	defer C.diagon_free_query(q)

	s.mu.RLock()
	td := C.diagon_search(s.searcher, q, 1)
	s.mu.RUnlock()
	if td == nil {
		return 0
	}
	hits := int64(C.diagon_top_docs_total_hits(td))
	C.diagon_free_top_docs(td)
	return hits
}
