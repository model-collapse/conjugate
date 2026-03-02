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

// compute_terms_agg_stored iterates ALL documents in the reader (0..maxDoc-1),
// reads a single stored field, and computes terms aggregation entirely in C.
// This avoids going through TopDocs / search results.
// Returns number of unique terms found. Results written to out_buckets (sorted by doc_count desc).
// max_buckets limits the output size.
static int compute_terms_agg_stored(
    DiagonIndexReader reader,
    const char* field_name,
    TermBucketC* out_buckets,
    int max_buckets)
{
    int max_doc = (int)diagon_reader_max_doc(reader);
    if (max_doc <= 0) return 0;

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

    for (int doc_id = 0; doc_id < max_doc; doc_id++) {
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
    // First pass: copy to temp array
    TermBucketC* all = (TermBucketC*)calloc(unique_count, sizeof(TermBucketC));
    if (!all) { free(map); return 0; }

    int out_idx = 0;
    for (int i = 0; i < map_cap && out_idx < unique_count; i++) {
        if (map[i].used) {
            strncpy(all[out_idx].key, map[i].key, 255);
            all[out_idx].key[255] = '\0';
            all[out_idx].doc_count = map[i].count;
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
*/
import "C"

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unsafe"

	json "github.com/goccy/go-json"

	"go.uber.org/zap"
)

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
	C.diagon_config_set_ram_buffer_size(config, 64.0)                   // 64MB buffer
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
				// Index as double for range queries
				doubleField := C.diagon_create_indexed_double_field(cFieldName, C.double(float64(epochMs)))
				C.diagon_document_add_field(diagonDoc, doubleField)
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


// BulkIndexDocuments indexes multiple documents in a batch to reduce CGO overhead.
// This is significantly faster than calling IndexDocument repeatedly.
func (s *Shard) BulkIndexDocuments(docs []struct {
	ID         string
	Doc        map[string]interface{}
	SourceJSON []byte
}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Debug("BulkIndexDocuments",
		zap.Int("num_docs", len(docs)))

	// Helper function to create a Diagon document from Go map
	createDiagonDoc := func(docID string, doc map[string]interface{}, rawJSON []byte) C.DiagonDocument {
		diagonDoc := C.diagon_create_document()

		// Add ID field - both indexed (for searching) and stored (for retrieval)
		cDocID := C.CString(docID)
		defer C.free(unsafe.Pointer(cDocID))
		cIDFieldName := C.CString("_id")
		defer C.free(unsafe.Pointer(cIDFieldName))

		// Add as StringField for exact-match searching
		idField := C.diagon_create_string_field(cIDFieldName, cDocID)
		C.diagon_document_add_field(diagonDoc, idField)

		// Add as StoredField for retrieval
		storedIDField := C.diagon_create_stored_field(cIDFieldName, cDocID)
		C.diagon_document_add_field(diagonDoc, storedIDField)

		// Store full _source as JSON for reliable document retrieval
		// Use pre-existing raw JSON bytes when available to avoid re-marshaling
		sourceBytes := rawJSON
		if sourceBytes == nil {
			sourceBytes, _ = json.Marshal(doc)
		}
		if sourceBytes != nil {
			cSourceFieldName := C.CString("_source")
			defer C.free(unsafe.Pointer(cSourceFieldName))
			cSourceValue := C.CString(string(sourceBytes))
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

			switch v := value.(type) {
			case string:
				if epochMs, isDate := tryParseDateToEpochMs(v); isDate {
					doubleField := C.diagon_create_indexed_double_field(cFieldName, C.double(float64(epochMs)))
					C.diagon_document_add_field(diagonDoc, doubleField)
					cValue := C.CString(v)
					defer C.free(unsafe.Pointer(cValue))
					storedField := C.diagon_create_stored_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, storedField)
				} else if isKeywordLike(v) {
					// Keyword-like: index as StringField (exact, not analyzed)
					// Enables fast terms aggregation via inverted index
					cValue := C.CString(v)
					defer C.free(unsafe.Pointer(cValue))
					stringField := C.diagon_create_string_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, stringField)
					storedField2 := C.diagon_create_stored_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, storedField2)
				} else {
					// Text field (analyzed) + StringField for terms aggregation
					cValue := C.CString(v)
					defer C.free(unsafe.Pointer(cValue))
					field := C.diagon_create_text_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, field)
					stringField := C.diagon_create_string_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, stringField)
				}

			case int, int32, int64:
				val := int64(0)
				switch n := v.(type) {
				case int:
					val = int64(n)
				case int32:
					val = int64(n)
				case int64:
					val = n
				}
				field := C.diagon_create_indexed_long_field(cFieldName, C.int64_t(val))
				C.diagon_document_add_field(diagonDoc, field)

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

				cValueStr := C.CString(fmt.Sprintf("%f", val))
				defer C.free(unsafe.Pointer(cValueStr))
				storedField := C.diagon_create_stored_field(cFieldName, cValueStr)
				C.diagon_document_add_field(diagonDoc, storedField)

			default:
				// Convert to JSON string for complex types
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					continue
				}
				cValue := C.CString(string(jsonBytes))
				defer C.free(unsafe.Pointer(cValue))
				field := C.diagon_create_stored_field(cFieldName, cValue)
				C.diagon_document_add_field(diagonDoc, field)
			}
		}

		return diagonDoc
	}

	// Index all documents in batch
	for i, item := range docs {
		diagonDoc := createDiagonDoc(item.ID, item.Doc, item.SourceJSON)
		defer C.diagon_free_document(diagonDoc)

		result := C.diagon_add_document(s.writer, diagonDoc)
		if !result {
			errMsg := C.GoString(C.diagon_last_error())
			return fmt.Errorf("failed to add document %d (ID: %s): %s", i, item.ID, errMsg)
		}
	}

	s.readerDirty = true
	s.logger.Debug("Bulk indexed documents to RAM buffer",
		zap.Int("count", len(docs)))

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

			s.logger.Debug(" Creating Diagon double range query",
				zap.String("field", field),
				zap.Float64("lower", lowerValue),
				zap.Float64("upper", upperValue),
				zap.Bool("include_lower", includeLower),
				zap.Bool("include_upper", includeUpper))

			// Use double range query for correct comparison semantics.
			// diagon_create_numeric_range_query uses bit_cast<int64_t>(double)
			// which breaks for negative doubles (bit patterns don't sort as integers).
			// diagon_create_double_range_query compares doubles directly.
			rangeQ := C.diagon_create_double_range_query(
				cField,
				C.double(lowerValue),
				C.double(upperValue),
				C.bool(includeLower),
				C.bool(includeUpper),
			)

			if rangeQ == nil {
				errMsg := C.GoString(C.diagon_last_error())
				s.logger.Error("DEBUG: Failed to create Diagon double range query", zap.String("error", errMsg))
				return nil, fmt.Errorf("failed to create double range query: %s", errMsg)
			}

			// Wrap range query in bool(must: match_all, filter: range) to exclude
			// phantom documents. Diagon's doc values iteration returns default 0.0
			// for documents without the field, causing false matches when the range
			// includes 0.0. match_all restricts results to real documents only.
			matchAllQ := C.diagon_create_match_all_query()
			if matchAllQ == nil {
				errMsg := C.GoString(C.diagon_last_error())
				return nil, fmt.Errorf("failed to create match_all for range wrapper: %s", errMsg)
			}

			boolBuilder := C.diagon_create_bool_query()
			if boolBuilder == nil {
				errMsg := C.GoString(C.diagon_last_error())
				return nil, fmt.Errorf("failed to create bool query for range wrapper: %s", errMsg)
			}

			C.diagon_bool_query_add_must(boolBuilder, matchAllQ)
			C.diagon_bool_query_add_filter(boolBuilder, rangeQ)

			diagonQuery = C.diagon_bool_query_build(boolBuilder)
			if diagonQuery == nil {
				errMsg := C.GoString(C.diagon_last_error())
				return nil, fmt.Errorf("failed to build range wrapper bool query: %s", errMsg)
			}
			s.logger.Debug(" Diagon double range query created successfully (wrapped with match_all)")
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
	return s.searchInternal(query, filterExpression, maxResults, fields)
}

// SearchWithLimit executes a search with a specified maximum number of results.
// For non-aggregation queries, use a small limit (e.g., from+size).
// For aggregation queries that need all matching docs, use a large limit.
func (s *Shard) SearchWithLimit(query []byte, filterExpression []byte, maxResults int) (*SearchResult, error) {
	return s.searchInternal(query, filterExpression, maxResults, nil)
}

func (s *Shard) searchInternal(query []byte, filterExpression []byte, maxResults int, fieldsOnly []string) (*SearchResult, error) {
	totalStart := time.Now()
	if maxResults <= 0 {
		maxResults = 100
	}

	reopenStart := time.Now()
	s.mu.Lock()

	// Only commit and reopen reader if there have been writes since last open
	needReopen := s.readerDirty || s.reader == nil
	if needReopen {
		// Commit pending changes to make them visible
		if s.writer != nil {
			if !C.diagon_commit(s.writer) {
				s.mu.Unlock()
				errMsg := C.GoString(C.diagon_last_error())
				s.logger.Warn("Failed to commit before search", zap.String("error", errMsg))
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
	reopenTime := time.Since(reopenStart)

	// Parse query JSON
	parseStart := time.Now()
	var queryObj map[string]interface{}
	if err := json.Unmarshal(query, &queryObj); err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	// Convert to Diagon query
	diagonQuery, err := s.convertQueryToDiagon(queryObj)
	if err != nil {
		return nil, err
	}
	defer C.diagon_free_query(diagonQuery)
	parseTime := time.Since(parseStart)

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
	if maxResults <= 0 {
		maxResults = 200000
	}

	// Ensure reader is open
	s.mu.Lock()
	needReopen := s.readerDirty || s.reader == nil
	if needReopen {
		if s.writer != nil {
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

	// Execute search
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
		// 3+ fields: one batch call per field
		for _, field := range fields {
			buf := make([]byte, bufSize)
			lengths := make([]C.int, numResults)
			cField := C.CString(field)

			C.batch_extract_field_values(
				s.reader, topDocs, cField,
				(*C.char)(unsafe.Pointer(&buf[0])), C.int(bufSize),
				&lengths[0], C.int(numResults))

			C.free(unsafe.Pointer(cField))

			offset := 0
			for i := 0; i < numResults; i++ {
				l := int(lengths[i])
				if l > 0 && offset+l <= len(buf) {
					docs[i].Fields[field] = AggFieldValue{StringVal: string(buf[offset : offset+l])}
					offset += l + 1
				}
			}
		}
	}
	s.mu.RUnlock()

	return totalHits, docs, nil
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

// ComputeTermsAgg computes a terms aggregation entirely in C, scanning all documents
// in the reader and building a hash map of field values → counts.
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

	// Ensure reader is open
	s.mu.Lock()
	if s.reader == nil || s.readerDirty {
		// Invalidate cache on reader reopen
		s.termsAggCacheMu.Lock()
		s.termsAggCache = nil
		s.termsAggCacheMu.Unlock()

		if s.writer != nil {
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

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.reader == nil {
		return nil, fmt.Errorf("reader not initialized")
	}

	cField := C.CString(field)
	defer C.free(unsafe.Pointer(cField))

	// Allocate output buckets
	outBuckets := make([]C.TermBucketC, size)

	n := C.compute_terms_agg_stored(
		s.reader,
		cField,
		&outBuckets[0],
		C.int(size),
	)

	if n <= 0 {
		return nil, nil
	}

	buckets := make([]TermBucket, int(n))
	for i := 0; i < int(n); i++ {
		buckets[i] = TermBucket{
			Key:      C.GoString(&outBuckets[i].key[0]),
			DocCount: int64(outBuckets[i].doc_count),
		}
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

// ForceMerge optimizes the index by merging segments
func (s *Shard) ForceMerge(maxSegments int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.writer == nil {
		return fmt.Errorf("writer not initialized")
	}

	s.logger.Info("Starting force merge", zap.Int("max_segments", maxSegments))

	if !C.diagon_force_merge(s.writer, C.int(maxSegments)) {
		errMsg := C.GoString(C.diagon_last_error())
		return fmt.Errorf("force merge failed: %s", errMsg)
	}

	s.logger.Info("Force merge completed", zap.Int("max_segments", maxSegments))
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
