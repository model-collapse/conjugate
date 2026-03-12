#!/usr/bin/env python3
"""Fixed Big5 benchmark — identical DSL queries for both OpenSearch and Conjugate.

Every query is byte-for-byte identical regardless of engine. No auto-discovery,
no per-engine adaptation. Uses hardcoded field names and values from the shared
Big5 dataset (big5_shared_data_gen.py).

Usage:
  python3 test/big5_fixed_benchmark.py --url http://localhost:9200 --label "OpenSearch" --output /tmp/os_fixed.json
  python3 test/big5_fixed_benchmark.py --url http://localhost:9201 --label "Conjugate" --output /tmp/conj_fixed.json
"""

import argparse
import hashlib
import json
import statistics
import sys
import time

try:
    import requests
except ImportError:
    print("ERROR: requests library required. Install with: pip3 install requests")
    sys.exit(1)


# --- Fixed constants from shared dataset ---
KW1 = "process.name"
KW1_VAL = "udev"
KW2 = "cloud.region"
KW2_VAL = "us-east-1"
KW3 = "host.name"
KW3_VAL = "host-001"
TEXT_FIELD = "message"
SEARCH_TERMS = "scale dog hero"
NUM1 = "metrics.size"
NUM2 = "metrics.tmin"

# Fixed time windows (14-day range: 2024-01-01 to 2024-01-14)
TW_FULL_START = "2024-01-01T00:00:00"
TW_FULL_END = "2024-01-14T23:59:59"
TW_SMALL_START = "2024-01-02T00:00:00"
TW_SMALL_END = "2024-01-03T00:00:00"
TW_MEDIUM_START = "2024-01-02T00:00:00"
TW_MEDIUM_END = "2024-01-08T00:00:00"
TW_LARGE_START = "2024-01-01T12:00:00"
TW_LARGE_END = "2024-01-12T12:00:00"
TW_MIDPOINT = "2024-01-07T12:00:00.000Z"


def build_queries():
    """Build all 42 fixed queries."""
    return {
        # Category 1: TEXT QUERYING (6 queries)
        "1_text_01_match_all": {
            "category": "Text Querying", "name": "match-all",
            "query": {"query": {"match_all": {}}}
        },
        "1_text_02_term": {
            "category": "Text Querying", "name": "term",
            "query": {"query": {"term": {KW1: {"value": KW1_VAL}}}}
        },
        "1_text_03_query_string": {
            "category": "Text Querying", "name": "query-string-on-message",
            "query": {"query": {"query_string": {"query": "message: scale dog hero"}}}
        },
        "1_text_04_query_string_filtered": {
            "category": "Text Querying", "name": "query-string-on-message-filtered",
            "query": {"query": {"bool": {"must": [
                {"range": {"@timestamp": {"gte": TW_SMALL_START, "lt": TW_SMALL_END}}},
                {"query_string": {"query": "message: scale dog hero"}}
            ]}}}
        },
        "1_text_05_query_string_filtered_sorted": {
            "category": "Text Querying", "name": "query-string-filtered-sorted-num",
            "query": {
                "query": {"bool": {"must": [
                    {"range": {"@timestamp": {"gte": TW_SMALL_START, "lt": TW_SMALL_END}}},
                    {"query_string": {"query": "message: scale dog hero"}}
                ]}},
                "sort": [{"@timestamp": {"order": "asc"}}]
            }
        },
        "1_text_06_keyword_in_range": {
            "category": "Text Querying", "name": "keyword-in-range",
            "query": {"query": {"bool": {"must": [
                {"range": {"@timestamp": {"gte": TW_SMALL_START, "lt": TW_SMALL_END}}},
                {"term": {KW1: KW1_VAL}}
            ]}}}
        },

        # Category 2: SORTING (13 queries)
        "2_sort_01_desc_timestamp": {
            "category": "Sorting", "name": "desc-sort-timestamp",
            "query": {"query": {"match_all": {}}, "sort": [{"@timestamp": "desc"}]}
        },
        "2_sort_02_desc_with_after": {
            "category": "Sorting", "name": "desc-sort-with-after-timestamp",
            "query": {"track_total_hits": False, "query": {"match_all": {}},
                      "sort": [{"@timestamp": "desc"}], "search_after": [TW_MIDPOINT]}
        },
        "2_sort_03_asc_timestamp": {
            "category": "Sorting", "name": "asc-sort-timestamp",
            "query": {"query": {"match_all": {}}, "sort": [{"@timestamp": "asc"}]}
        },
        "2_sort_04_asc_with_after": {
            "category": "Sorting", "name": "asc-sort-with-after-timestamp",
            "query": {"track_total_hits": False, "query": {"match_all": {}},
                      "sort": [{"@timestamp": "asc"}], "search_after": [TW_MIDPOINT]}
        },
        "2_sort_05_desc_can_match": {
            "category": "Sorting", "name": "desc-sort-timestamp-can-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {KW1: KW1_VAL}},
                      "sort": [{"@timestamp": "desc"}]}
        },
        "2_sort_06_asc_can_match": {
            "category": "Sorting", "name": "asc-sort-timestamp-can-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {KW1: KW1_VAL}},
                      "sort": [{"@timestamp": "asc"}]}
        },
        "2_sort_07_keyword_sort": {
            "category": "Sorting", "name": "sort-keyword-can-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {KW1: KW1_VAL}},
                      "sort": [{KW2: "asc"}]}
        },
        "2_sort_08_numeric_desc": {
            "category": "Sorting", "name": "sort-numeric-desc",
            "query": {"track_total_hits": False, "query": {"match_all": {}},
                      "sort": [{NUM1: "desc"}]}
        },
        "2_sort_09_numeric_asc": {
            "category": "Sorting", "name": "sort-numeric-asc",
            "query": {"track_total_hits": False, "query": {"match_all": {}},
                      "sort": [{NUM1: "asc"}]}
        },
        "2_sort_10_numeric_desc_with_match": {
            "category": "Sorting", "name": "sort-numeric-desc-with-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {KW1: KW1_VAL}},
                      "sort": [{NUM1: "desc"}]}
        },
        "2_sort_11_numeric_asc_with_match": {
            "category": "Sorting", "name": "sort-numeric-asc-with-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {KW1: KW1_VAL}},
                      "sort": [{NUM1: "asc"}]}
        },
        "2_sort_12_range_asc": {
            "category": "Sorting", "name": "range-with-asc-sort",
            "query": {"query": {"range": {"@timestamp": {"gte": TW_LARGE_START, "lte": TW_LARGE_END}}},
                      "sort": [{"@timestamp": "asc"}]}
        },
        "2_sort_13_range_desc": {
            "category": "Sorting", "name": "range-with-desc-sort",
            "query": {"query": {"range": {"@timestamp": {"gte": TW_LARGE_START, "lte": TW_LARGE_END}}},
                      "sort": [{"@timestamp": "desc"}]}
        },

        # Category 3: DATE HISTOGRAM (4 queries)
        "3_hist_01_hourly": {
            "category": "Date Histogram", "name": "date-histogram-hourly",
            "query": {"size": 0, "aggs": {"by_hour": {"date_histogram": {
                "field": "@timestamp", "calendar_interval": "hour"}}}}
        },
        "3_hist_02_hourly_with_filter": {
            "category": "Date Histogram", "name": "date-histogram-hourly-with-filter",
            "query": {"size": 0, "query": {"term": {KW1: KW1_VAL}},
                      "aggs": {"by_hour": {"date_histogram": {
                          "field": "@timestamp", "calendar_interval": "hour"}}}}
        },
        "3_hist_03_minute": {
            "category": "Date Histogram", "name": "date-histogram-minute",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": TW_SMALL_START, "lt": TW_SMALL_END}}},
                      "aggs": {"by_hour": {"date_histogram": {
                          "field": "@timestamp", "calendar_interval": "minute"}}}}
        },
        "3_hist_04_composite_daily": {
            "category": "Date Histogram", "name": "composite-date-histogram-daily",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": TW_FULL_START, "lt": TW_FULL_END}}},
                      "aggs": {"logs": {"composite": {"sources": [
                          {"date": {"date_histogram": {"field": "@timestamp", "calendar_interval": "day"}}}
                      ]}}}}
        },

        # Category 4: RANGE QUERIES (10 queries)
        "4_range_01_date": {
            "category": "Range Queries", "name": "range-date",
            "query": {"query": {"range": {"@timestamp": {"gte": TW_SMALL_START, "lt": TW_SMALL_END}}}}
        },
        "4_range_02_numeric": {
            "category": "Range Queries", "name": "range-numeric",
            "query": {"query": {"range": {NUM1: {"gte": 20, "lte": 200}}}}
        },
        "4_range_03_conj_big_range_big_term": {
            "category": "Range Queries", "name": "range-conjunction-big-range-big-term",
            "query": {"query": {"bool": {"must": [
                {"term": {KW1: KW1_VAL}},
                {"range": {NUM1: {"gte": 1, "lte": 100}}}
            ]}}}
        },
        "4_range_04_disj_big_range_small_term": {
            "category": "Range Queries", "name": "range-disjunction-big-range-small-term",
            "query": {"query": {"bool": {"should": [
                {"term": {KW2: KW2_VAL}},
                {"range": {NUM1: {"gte": 1, "lte": 100}}}
            ]}}}
        },
        "4_range_05_small_range_only": {
            "category": "Range Queries", "name": "range-small-range-big-term",
            "query": {"query": {"bool": {"must": [{"range": {NUM1: {"gte": 20, "lte": 30}}}]}}}
        },
        "4_range_06_agg": {
            "category": "Range Queries", "name": "range-agg-1",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": NUM1,
                      "ranges": [{"to": -10}, {"from": -10, "to": 10}, {"from": 10, "to": 100},
                                 {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]}}}}
        },
        "4_range_07_agg_2": {
            "category": "Range Queries", "name": "range-agg-2",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": NUM1,
                      "ranges": [{"to": 100}, {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]}}}}
        },
        "4_range_08_with_metrics": {
            "category": "Range Queries", "name": "range-with-metrics",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": NUM1,
                      "ranges": [{"to": -10}, {"from": -10, "to": 10}, {"from": 10, "to": 100},
                                 {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]},
                      "aggs": {"tsum": {"sum": {"field": NUM1}}, "tmin": {"min": {"field": NUM1}},
                               "tavg": {"avg": {"field": NUM1}}, "tmax_m": {"max": {"field": NUM1}},
                               "tstats": {"stats": {"field": NUM1}}}}}}
        },
        "4_range_09_auto_date_histo": {
            "category": "Range Queries", "name": "range-auto-date-histo",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": NUM1,
                      "ranges": [{"to": -10}, {"from": -10, "to": 10}, {"from": 10, "to": 100},
                                 {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]},
                      "aggs": {"date": {"auto_date_histogram": {"field": "@timestamp", "buckets": 20}}}}}}
        },
        "4_range_10_auto_date_histo_with_metrics": {
            "category": "Range Queries", "name": "range-auto-date-histo-with-metrics",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": NUM1,
                      "ranges": [{"to": 100}, {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]},
                      "aggs": {"date": {"auto_date_histogram": {"field": "@timestamp", "buckets": 10},
                               "aggs": {"tmin": {"min": {"field": NUM1}},
                                        "tavg": {"avg": {"field": NUM1}},
                                        "tmax_m": {"max": {"field": NUM1}}}}}}}}
        },

        # Category 5: TERMS AGGREGATION (9 queries)
        "5_terms_01_keyword_terms_500": {
            "category": "Terms Aggregation", "name": "keyword-terms-500",
            "query": {"size": 0, "aggs": {"station": {"terms": {"field": KW1, "size": 500}}}}
        },
        "5_terms_02_keyword_terms_50": {
            "category": "Terms Aggregation", "name": "keyword-terms-50",
            "query": {"size": 0, "aggs": {"country": {"terms": {"field": KW1, "size": 50}}}}
        },
        "5_terms_03_significant_1": {
            "category": "Terms Aggregation", "name": "terms-significant-1",
            "query": {"track_total_hits": False, "size": 0,
                      "query": {"range": {"@timestamp": {"gte": TW_SMALL_START, "lt": TW_SMALL_END}}},
                      "aggs": {"terms": {"terms": {"field": KW1, "size": 10},
                               "aggs": {"significant_ips": {"significant_terms": {"field": KW2}}}}}}
        },
        "5_terms_04_significant_2": {
            "category": "Terms Aggregation", "name": "terms-significant-2",
            "query": {"track_total_hits": False, "size": 0,
                      "query": {"range": {"@timestamp": {"gte": TW_SMALL_START, "lt": TW_SMALL_END}}},
                      "aggs": {"terms": {"terms": {"field": KW2, "size": 10},
                               "aggs": {"significant_ips": {"significant_terms": {"field": KW1}}}}}}
        },
        "5_terms_05_multi_terms": {
            "category": "Terms Aggregation", "name": "multi-terms-keyword",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": TW_MEDIUM_START, "lt": TW_MEDIUM_END}}},
                      "aggs": {"important_terms": {"multi_terms": {"terms": [{"field": KW1}, {"field": KW2}]}}}}
        },
        "5_terms_06_composite": {
            "category": "Terms Aggregation", "name": "composite-terms",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": TW_SMALL_START, "lt": TW_SMALL_END}}},
                      "aggs": {"logs": {"composite": {"sources": [
                          {"process_name": {"terms": {"field": KW1, "order": "desc"}}},
                          {"cloud_region": {"terms": {"field": KW2, "order": "asc"}}}
                      ]}}}}
        },
        "5_terms_07_composite_3key": {
            "category": "Terms Aggregation", "name": "composite-terms-3key",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": TW_SMALL_START, "lt": TW_SMALL_END}}},
                      "aggs": {"logs": {"composite": {"sources": [
                          {"process_name": {"terms": {"field": KW1, "order": "desc"}}},
                          {"cloud_region": {"terms": {"field": KW2, "order": "asc"}}},
                          {"host_name": {"terms": {"field": KW3, "order": "asc"}}}
                      ]}}}}
        },
        "5_terms_08_cardinality_low": {
            "category": "Terms Aggregation", "name": "cardinality-low",
            "query": {"size": 0, "aggs": {"region": {"cardinality": {"field": KW1}}}}
        },
        "5_terms_09_cardinality_high": {
            "category": "Terms Aggregation", "name": "cardinality-high",
            "query": {"size": 0, "aggs": {"agent": {"cardinality": {"field": KW2}}}}
        },
    }


def percentile(sorted_data, p):
    if not sorted_data:
        return 0.0
    n = len(sorted_data)
    k = (n - 1) * (p / 100.0)
    f = int(k)
    c = min(f + 1, n - 1)
    return sorted_data[f] + (k - f) * (sorted_data[c] - sorted_data[f]) if f != c else sorted_data[f]


def run_query(session, url, index, query_body, timeout=120):
    start = time.perf_counter()
    try:
        resp = session.post(f'{url}/{index}/_search', json=query_body, timeout=timeout)
        elapsed_ms = (time.perf_counter() - start) * 1000
        if resp.status_code != 200:
            return elapsed_ms, 0, resp.status_code, None
        result = resp.json()
        hits = result.get('hits', {}).get('total', {})
        total = hits.get('value', 0) if isinstance(hits, dict) else (int(hits) if hits else 0)
        return elapsed_ms, total, 200, result
    except requests.exceptions.Timeout:
        return (time.perf_counter() - start) * 1000, 0, -1, None
    except Exception:
        return (time.perf_counter() - start) * 1000, 0, -2, None


def capture_parity(result):
    """Extract parity-checkable snapshot from a search response.

    Returns dict with:
      hits_total: total hit count
      agg_summary: {agg_name: {type, bucket_count, first_3_buckets}} for bucket aggs
                   {agg_name: {type, value}} for metric aggs
    """
    if not result:
        return None
    snapshot = {}

    # Hit count
    hits = result.get('hits', {}).get('total', {})
    if isinstance(hits, dict):
        snapshot['hits_total'] = hits.get('value', 0)
    else:
        snapshot['hits_total'] = int(hits) if hits else 0

    # Aggregation summary
    aggs = result.get('aggregations', {})
    if aggs:
        agg_summary = {}
        for agg_name, agg_data in aggs.items():
            if not isinstance(agg_data, dict):
                continue
            if 'buckets' in agg_data:
                buckets = agg_data['buckets']
                if isinstance(buckets, list):
                    first_3 = []
                    for b in buckets[:3]:
                        first_3.append({
                            "key": b.get("key_as_string", b.get("key")),
                            "doc_count": b.get("doc_count", 0),
                        })
                    agg_summary[agg_name] = {
                        "type": "bucketed",
                        "bucket_count": len(buckets),
                        "first_3_buckets": first_3,
                    }
                elif isinstance(buckets, dict):
                    first_3 = []
                    for k in list(buckets.keys())[:3]:
                        first_3.append({
                            "key": k,
                            "doc_count": buckets[k].get("doc_count", 0),
                        })
                    agg_summary[agg_name] = {
                        "type": "named_buckets",
                        "bucket_count": len(buckets),
                        "first_3_buckets": first_3,
                    }
            elif 'value' in agg_data:
                agg_summary[agg_name] = {"type": "metric", "value": agg_data['value']}
        if agg_summary:
            snapshot['agg_summary'] = agg_summary

    return snapshot


def probe_query(session, url, index, query_body):
    times = []
    last_status = 200
    for _ in range(3):
        lat, _, status, _ = run_query(session, url, index, query_body, timeout=120)
        last_status = status
        if status == 200:
            times.append(lat)
        elif status == -1:
            return 120000, 'timeout', 0, 3, last_status
    if not times:
        return 0, 'error', 0, 3, last_status
    avg = statistics.mean(times)
    if avg < 100:
        return avg, 'fast', 100, 200, 200
    elif avg < 2000:
        return avg, 'medium', 30, 100, 200
    elif avg < 30000:
        return avg, 'slow', 5, 20, 200
    else:
        return avg, 'very_slow', 2, 5, 200


def benchmark_queries(url, index, queries, max_warmup, max_iterations):
    session = requests.Session()
    session.headers.update({'Content-Type': 'application/json'})
    results = {}

    for key in sorted(queries.keys()):
        qdef = queries[key]
        print(f"\n  [{qdef['category']}] {qdef['name']}")

        avg_probe, tier, warmup, iterations, probe_status = probe_query(
            session, url, index, qdef['query'])

        warmup = min(warmup, max_warmup)
        iterations = min(iterations, max_iterations)

        if tier == 'timeout':
            print(f"    TIMEOUT -- skipping")
            results[key] = {'name': qdef['name'], 'category': qdef['category'],
                            'timeout': True, 'tier': tier}
            continue

        if tier == 'error':
            print(f"    ERROR (status {probe_status}) -- skipping")
            results[key] = {'name': qdef['name'], 'category': qdef['category'],
                            'error': True, 'status_code': probe_status, 'tier': tier}
            continue

        print(f"    Probe: {avg_probe:.0f}ms ({tier}) -> {warmup}w + {iterations}i")

        print(f"    Warmup ({warmup})...", end='', flush=True)
        for _ in range(warmup):
            run_query(session, url, index, qdef['query'])
        print(" done")

        latencies = []
        hit_count = 0
        errors = 0
        timeouts = 0
        parity = None
        print(f"    Measuring ({iterations})", end='', flush=True)
        for i in range(iterations):
            lat, hits, status, result = run_query(session, url, index, qdef['query'])
            if status == 200:
                latencies.append(lat)
                hit_count = hits
                if parity is None:
                    parity = capture_parity(result)
            elif status == -1:
                timeouts += 1
            else:
                errors += 1
            if (i + 1) % max(1, iterations // 4) == 0:
                print(".", end='', flush=True)
        err_str = f" ({errors} errors)" if errors else ""
        err_str += f" ({timeouts} timeouts)" if timeouts else ""
        print(f" done{err_str}")

        if not latencies:
            results[key] = {'name': qdef['name'], 'category': qdef['category'],
                            'error': True, 'errors': errors, 'timeouts': timeouts, 'tier': tier}
            continue

        latencies.sort()
        stats = {
            'name': qdef['name'], 'category': qdef['category'], 'tier': tier,
            'iterations': len(latencies), 'errors': errors, 'timeouts': timeouts,
            'hits': hit_count,
            'avg': round(statistics.mean(latencies), 2),
            'min': round(latencies[0], 2), 'max': round(latencies[-1], 2),
            'p50': round(percentile(latencies, 50), 2),
            'p90': round(percentile(latencies, 90), 2),
            'p95': round(percentile(latencies, 95), 2),
            'p99': round(percentile(latencies, 99), 2),
            'stddev': round(statistics.stdev(latencies), 2) if len(latencies) > 1 else 0,
        }
        if parity:
            stats['parity'] = parity
        results[key] = stats
        print(f"    Hits: {hit_count:,} | Avg: {stats['avg']:.1f}ms | "
              f"P50: {stats['p50']:.1f}ms | P90: {stats['p90']:.1f}ms | "
              f"P95: {stats['p95']:.1f}ms | P99: {stats['p99']:.1f}ms")

    session.close()
    return results


def verify_schema(url, index):
    """Quick schema verification: check doc count and sample a doc."""
    session = requests.Session()
    session.headers.update({'Content-Type': 'application/json'})
    try:
        r = session.get(f'{url}/{index}/_count', timeout=10)
        count = r.json().get('count', 0) if r.status_code == 200 else 0
    except Exception:
        count = 0

    # Verify key fields exist
    missing = []
    try:
        r = session.post(f'{url}/{index}/_search', json={"size": 1}, timeout=10)
        if r.status_code == 200:
            hits = r.json().get('hits', {}).get('hits', [])
            if hits:
                src = hits[0].get('_source', {})
                for field_path in ['@timestamp', 'process.name', 'cloud.region',
                                   'host.name', 'message', 'metrics.size']:
                    parts = field_path.split('.')
                    val = src
                    for p in parts:
                        if isinstance(val, dict):
                            val = val.get(p)
                        else:
                            val = None
                            break
                    if val is None:
                        missing.append(field_path)
    except Exception:
        pass

    session.close()
    return count, missing


def main():
    parser = argparse.ArgumentParser(description='Fixed Big5 Benchmark (identical DSL)')
    parser.add_argument('--url', default='http://localhost:9200')
    parser.add_argument('--index', default='big5-benchmark')
    parser.add_argument('--label', default='Engine')
    parser.add_argument('--warmup', type=int, default=3)
    parser.add_argument('--iterations', type=int, default=10)
    parser.add_argument('--output', default=None,
                        help='JSON output file (default: /tmp/<label>_fixed_results.json)')
    args = parser.parse_args()

    if args.output is None:
        safe_label = args.label.lower().replace(' ', '_')
        args.output = f'/tmp/{safe_label}_fixed_results.json'

    print("=" * 70)
    print(f"  Fixed Big5 Benchmark -- {args.label}")
    print(f"  Target: {args.url}/{args.index}")
    print("=" * 70)

    # Verify schema
    count, missing = verify_schema(args.url, args.index)
    print(f"\n  Documents: {count:,}")
    if missing:
        print(f"  WARNING: Missing fields: {missing}")
        print(f"  This index may not have the shared Big5 data. Results will be unreliable.")
    if count == 0:
        print(f"  ERROR: No documents in index")
        sys.exit(1)

    queries = build_queries()
    print(f"  Queries: {len(queries)} operations across 5 categories")
    print(f"  Config: up to {args.warmup} warmup + {args.iterations} measured per query")

    # Verify identical DSL
    query_json = json.dumps({k: v['query'] for k, v in sorted(queries.items())}, sort_keys=True)
    dsl_hash = hashlib.sha256(query_json.encode()).hexdigest()[:12]
    print(f"  DSL hash: {dsl_hash} (must match across engines)")

    # Run benchmark
    print(f"\n{'='*70}")
    query_start = time.time()
    results = benchmark_queries(args.url, args.index, queries, args.warmup, args.iterations)
    query_elapsed = time.time() - query_start

    # Category summary
    categories = {}
    cat_totals = {}
    for key, r in sorted(results.items()):
        cat = r['category']
        cat_totals[cat] = cat_totals.get(cat, 0) + 1
        if not r.get('error') and not r.get('timeout'):
            categories.setdefault(cat, []).append(r)

    all_cats = ["Text Querying", "Sorting", "Date Histogram", "Range Queries", "Terms Aggregation"]

    print(f"\n{'='*70}")
    print(f"  CATEGORY SUMMARY -- {args.label}")
    print(f"{'='*70}\n")
    print(f"  {'Category':<22} {'OK':>3}/{'Tot':>3} {'Avg P50':>10} {'Avg P90':>10} {'Avg P99':>10}")
    print("  " + "-" * 60)

    for cat in all_cats:
        total = cat_totals.get(cat, 0)
        if cat in categories:
            qs = categories[cat]
            n = len(qs)
            avg_p50 = statistics.mean(q['p50'] for q in qs)
            avg_p90 = statistics.mean(q['p90'] for q in qs)
            avg_p99 = statistics.mean(q['p99'] for q in qs)
            print(f"  {cat:<22} {n:>3}/{total:>3} {avg_p50:>9.1f}ms {avg_p90:>9.1f}ms {avg_p99:>9.1f}ms")
        else:
            print(f"  {cat:<22} {0:>3}/{total:>3}   (all failed)")

    # Per-query results
    print(f"\n{'='*70}")
    print("  PER-QUERY RESULTS")
    print(f"{'='*70}\n")
    print(f"  {'Query':<42} {'Avg':>9} {'P50':>9} {'P90':>9} {'P99':>9}  {'Hits':>12}")
    print("  " + "-" * 95)
    current_cat = None
    for key in sorted(results.keys()):
        r = results[key]
        cat = r.get('category', '')
        if cat != current_cat:
            current_cat = cat
            print(f"\n  --- {cat} ---")
        if r.get('timeout'):
            print(f"  {r['name']:<42} {'TIMEOUT':>9}")
            continue
        if r.get('error'):
            print(f"  {r['name']:<42} {'ERROR':>9}")
            continue
        print(f"  {r['name']:<42} {r['avg']:>8.1f}ms {r['p50']:>8.1f}ms "
              f"{r['p90']:>8.1f}ms {r['p99']:>8.1f}ms  {r['hits']:>12,}")

    ok = sum(1 for r in results.values() if not r.get('error') and not r.get('timeout'))
    print(f"\n  Completed: {ok}/{len(results)} | Total time: {query_elapsed:.1f}s")

    # Save
    output = {
        'timestamp': time.strftime('%Y-%m-%d %H:%M:%S'),
        'label': args.label,
        'doc_count': count,
        'dsl_hash': dsl_hash,
        'config': {'url': args.url, 'index': args.index,
                   'warmup': args.warmup, 'iterations': args.iterations},
        'constants': {
            'kw1': KW1, 'kw1_val': KW1_VAL,
            'kw2': KW2, 'kw2_val': KW2_VAL,
            'kw3': KW3, 'kw3_val': KW3_VAL,
            'text_field': TEXT_FIELD, 'search_terms': SEARCH_TERMS,
            'num1': NUM1, 'num2': NUM2,
            'full_start': TW_FULL_START, 'full_end': TW_FULL_END,
            'small_start': TW_SMALL_START, 'small_end': TW_SMALL_END,
            'medium_start': TW_MEDIUM_START, 'medium_end': TW_MEDIUM_END,
            'large_start': TW_LARGE_START, 'large_end': TW_LARGE_END,
            'midpoint': TW_MIDPOINT,
        },
        'total_time_sec': round(query_elapsed, 2),
        'queries': results,
    }
    with open(args.output, 'w') as f:
        json.dump(output, f, indent=2)
    print(f"  Results saved to: {args.output}")
    print(f"{'='*70}")


if __name__ == '__main__':
    main()
