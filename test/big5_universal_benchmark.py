#!/usr/bin/env python3
"""Universal Big5 benchmark that auto-discovers schema and timestamp range.

Generates 42 equivalent queries adapted to the actual data, ensuring fair
apples-to-apples comparison between engines with different datasets.
"""

import argparse
import json
import sys
import time
import statistics

try:
    import requests
except ImportError:
    print("ERROR: requests library required. Install with: pip3 install requests")
    sys.exit(1)


def percentile(sorted_data, p):
    if not sorted_data:
        return 0.0
    n = len(sorted_data)
    k = (n - 1) * (p / 100.0)
    f = int(k)
    c = min(f + 1, n - 1)
    return sorted_data[f] + (k - f) * (sorted_data[c] - sorted_data[f]) if f != c else sorted_data[f]


def discover_schema(session, url, index):
    """Auto-discover index schema: timestamp range, keyword/text/numeric fields."""
    info = {'ts_min': None, 'ts_max': None, 'keywords': [], 'text_fields': [],
            'numeric_fields': [], 'doc_count': 0, 'segments': 0}

    # Doc count
    try:
        r = session.get(f'{url}/{index}/_count', timeout=10)
        info['doc_count'] = r.json().get('count', 0) if r.status_code == 200 else 0
    except Exception:
        pass

    # Segments
    try:
        r = session.get(f'{url}/{index}/_stats/segments', timeout=10)
        if r.status_code == 200:
            info['segments'] = r.json().get('_all', {}).get('primaries', {}).get('segments', {}).get('count', 0)
    except Exception:
        pass

    # Timestamp range
    try:
        r = session.post(f'{url}/{index}/_search', json={
            "size": 0,
            "aggs": {
                "min_ts": {"min": {"field": "@timestamp"}},
                "max_ts": {"max": {"field": "@timestamp"}}
            }
        }, timeout=30)
        if r.status_code == 200:
            aggs = r.json().get('aggregations', {})
            info['ts_min'] = aggs.get('min_ts', {}).get('value_as_string')
            info['ts_max'] = aggs.get('max_ts', {}).get('value_as_string')
    except Exception:
        pass

    # Get mapping to discover field types
    try:
        r = session.get(f'{url}/{index}/_mapping', timeout=10)
        if r.status_code == 200:
            mapping = r.json()
            # Navigate to properties
            for idx_name, idx_data in mapping.items():
                props = idx_data.get('mappings', {}).get('properties', {})
                _extract_fields(props, '', info)
    except Exception:
        pass

    return info


def _extract_fields(props, prefix, info):
    """Recursively extract field types from mapping."""
    for name, field_def in props.items():
        full_name = f'{prefix}{name}' if not prefix else f'{prefix}.{name}'
        ftype = field_def.get('type', '')

        if ftype == 'keyword':
            info['keywords'].append(full_name)
        elif ftype == 'text':
            info['text_fields'].append(full_name)
            # Check for .keyword subfield (OpenSearch dynamic mapping)
            if 'fields' in field_def:
                for sub_name, sub_def in field_def['fields'].items():
                    if sub_def.get('type') == 'keyword':
                        info['keywords'].append(f'{full_name}.{sub_name}')
        elif ftype in ('long', 'integer', 'float', 'double', 'short'):
            info['numeric_fields'].append(full_name)

        # Recurse into nested properties
        if 'properties' in field_def:
            _extract_fields(field_def['properties'], full_name, info)


def compute_time_windows(ts_min_str, ts_max_str):
    """Compute time windows for queries based on actual data range."""
    from datetime import datetime, timedelta

    # Parse ISO timestamps (handle various formats)
    for fmt in ('%Y-%m-%dT%H:%M:%S.%fZ', '%Y-%m-%dT%H:%M:%SZ', '%Y-%m-%dT%H:%M:%S'):
        try:
            ts_min = datetime.strptime(ts_min_str.rstrip('Z') + 'Z', fmt if fmt.endswith('Z') else fmt + 'Z')
            break
        except ValueError:
            continue
    else:
        ts_min = datetime(2023, 1, 1)

    for fmt in ('%Y-%m-%dT%H:%M:%S.%fZ', '%Y-%m-%dT%H:%M:%SZ', '%Y-%m-%dT%H:%M:%S'):
        try:
            ts_max = datetime.strptime(ts_max_str.rstrip('Z') + 'Z', fmt if fmt.endswith('Z') else fmt + 'Z')
            break
        except ValueError:
            continue
    else:
        ts_max = datetime(2023, 1, 14)

    total_span = ts_max - ts_min
    if total_span.total_seconds() <= 0:
        total_span = timedelta(hours=1)

    # Window sizes as fractions of total span
    return {
        'full_start': ts_min.strftime('%Y-%m-%dT%H:%M:%S'),
        'full_end': ts_max.strftime('%Y-%m-%dT%H:%M:%S'),
        # ~10% of data
        'small_start': (ts_min + total_span * 0.1).strftime('%Y-%m-%dT%H:%M:%S'),
        'small_end': (ts_min + total_span * 0.2).strftime('%Y-%m-%dT%H:%M:%S'),
        # ~40% of data
        'medium_start': (ts_min + total_span * 0.1).strftime('%Y-%m-%dT%H:%M:%S'),
        'medium_end': (ts_min + total_span * 0.5).strftime('%Y-%m-%dT%H:%M:%S'),
        # ~80% of data
        'large_start': (ts_min + total_span * 0.05).strftime('%Y-%m-%dT%H:%M:%S'),
        'large_end': (ts_min + total_span * 0.85).strftime('%Y-%m-%dT%H:%M:%S'),
        # midpoint for search_after
        'midpoint': (ts_min + total_span * 0.5).strftime('%Y-%m-%dT%H:%M:%S.000Z'),
    }


def discover_keyword_value(session, url, index, field):
    """Find an actual value for a keyword field using terms agg."""
    try:
        r = session.post(f'{url}/{index}/_search', json={
            "size": 0, "aggs": {"vals": {"terms": {"field": field, "size": 1}}}
        }, timeout=10)
        if r.status_code == 200:
            buckets = r.json().get('aggregations', {}).get('vals', {}).get('buckets', [])
            if buckets:
                return buckets[0]['key']
    except Exception:
        pass
    return None


def discover_text_terms(session, url, index, field):
    """Find searchable terms by sampling a document's text field."""
    try:
        r = session.post(f'{url}/{index}/_search', json={"size": 1, "_source": [field]}, timeout=10)
        if r.status_code == 200:
            hits = r.json().get('hits', {}).get('hits', [])
            if hits:
                src = hits[0].get('_source', {})
                # Navigate nested field
                parts = field.split('.')
                val = src
                for p in parts:
                    if isinstance(val, dict):
                        val = val.get(p, '')
                    else:
                        val = ''
                        break
                if isinstance(val, str) and len(val) > 3:
                    words = val.split()
                    # Pick 3 words that are likely to not match everything
                    import random
                    if len(words) >= 3:
                        return ' '.join(random.sample(words, 3))
                    return val
    except Exception:
        pass
    return "error timeout connection"


def build_queries(schema, session, url, index):
    """Build 42 equivalent queries adapted to the discovered schema."""
    tw = compute_time_windows(schema['ts_min'], schema['ts_max'])

    # Pick fields
    kw_fields = schema['keywords'][:5]  # up to 5 keyword fields
    text_field = schema['text_fields'][0] if schema['text_fields'] else 'message'
    num_fields = schema['numeric_fields'][:3]
    kw1 = kw_fields[0] if kw_fields else 'host.name'
    kw2 = kw_fields[1] if len(kw_fields) > 1 else kw1
    kw3 = kw_fields[2] if len(kw_fields) > 2 else kw1
    num1 = num_fields[0] if num_fields else 'network.bytes'

    # Discover actual values
    kw1_val = discover_keyword_value(session, url, index, kw1) or "unknown"
    kw2_val = discover_keyword_value(session, url, index, kw2) or "unknown"
    search_terms = discover_text_terms(session, url, index, text_field)

    print(f"  Schema discovery:")
    print(f"    Keyword fields: {kw_fields[:5]}")
    print(f"    Text field: {text_field}")
    print(f"    Numeric fields: {num_fields[:3]}")
    print(f"    Keyword values: {kw1}={kw1_val}, {kw2}={kw2_val}")
    print(f"    Text search: '{search_terms}'")
    print(f"    Timestamp range: {tw['full_start']} → {tw['full_end']}")
    print()

    queries = {
        # Category 1: TEXT QUERYING (6 queries)
        "1_text_01_match_all": {
            "category": "Text Querying", "name": "match-all",
            "query": {"query": {"match_all": {}}}
        },
        "1_text_02_term": {
            "category": "Text Querying", "name": "term",
            "query": {"query": {"term": {kw1: {"value": kw1_val}}}}
        },
        "1_text_03_query_string": {
            "category": "Text Querying", "name": "query-string-on-message",
            "query": {"query": {"query_string": {"query": f"{text_field}: {search_terms}"}}}
        },
        "1_text_04_query_string_filtered": {
            "category": "Text Querying", "name": "query-string-on-message-filtered",
            "query": {"query": {"bool": {"must": [
                {"range": {"@timestamp": {"gte": tw['small_start'], "lt": tw['small_end']}}},
                {"query_string": {"query": f"{text_field}: {search_terms}"}}
            ]}}}
        },
        "1_text_05_query_string_filtered_sorted": {
            "category": "Text Querying", "name": "query-string-filtered-sorted-num",
            "query": {
                "query": {"bool": {"must": [
                    {"range": {"@timestamp": {"gte": tw['small_start'], "lt": tw['small_end']}}},
                    {"query_string": {"query": f"{text_field}: {search_terms}"}}
                ]}},
                "sort": [{"@timestamp": {"order": "asc"}}]
            }
        },
        "1_text_06_keyword_in_range": {
            "category": "Text Querying", "name": "keyword-in-range",
            "query": {"query": {"bool": {"must": [
                {"range": {"@timestamp": {"gte": tw['small_start'], "lt": tw['small_end']}}},
                {"term": {kw1: kw1_val}}
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
                      "sort": [{"@timestamp": "desc"}], "search_after": [tw['midpoint']]}
        },
        "2_sort_03_asc_timestamp": {
            "category": "Sorting", "name": "asc-sort-timestamp",
            "query": {"query": {"match_all": {}}, "sort": [{"@timestamp": "asc"}]}
        },
        "2_sort_04_asc_with_after": {
            "category": "Sorting", "name": "asc-sort-with-after-timestamp",
            "query": {"track_total_hits": False, "query": {"match_all": {}},
                      "sort": [{"@timestamp": "asc"}], "search_after": [tw['midpoint']]}
        },
        "2_sort_05_desc_can_match": {
            "category": "Sorting", "name": "desc-sort-timestamp-can-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {kw1: kw1_val}},
                      "sort": [{"@timestamp": "desc"}]}
        },
        "2_sort_06_asc_can_match": {
            "category": "Sorting", "name": "asc-sort-timestamp-can-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {kw1: kw1_val}},
                      "sort": [{"@timestamp": "asc"}]}
        },
        "2_sort_07_keyword_sort": {
            "category": "Sorting", "name": "sort-keyword-can-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {kw1: kw1_val}},
                      "sort": [{kw2: "asc"}]}
        },
        "2_sort_08_numeric_desc": {
            "category": "Sorting", "name": "sort-numeric-desc",
            "query": {"track_total_hits": False, "query": {"match_all": {}},
                      "sort": [{num1: "desc"}]}
        },
        "2_sort_09_numeric_asc": {
            "category": "Sorting", "name": "sort-numeric-asc",
            "query": {"track_total_hits": False, "query": {"match_all": {}},
                      "sort": [{num1: "asc"}]}
        },
        "2_sort_10_numeric_desc_with_match": {
            "category": "Sorting", "name": "sort-numeric-desc-with-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {kw1: kw1_val}},
                      "sort": [{num1: "desc"}]}
        },
        "2_sort_11_numeric_asc_with_match": {
            "category": "Sorting", "name": "sort-numeric-asc-with-match",
            "query": {"track_total_hits": False,
                      "query": {"term": {kw1: kw1_val}},
                      "sort": [{num1: "asc"}]}
        },
        "2_sort_12_range_asc": {
            "category": "Sorting", "name": "range-with-asc-sort",
            "query": {"query": {"range": {"@timestamp": {"gte": tw['large_start'], "lte": tw['large_end']}}},
                      "sort": [{"@timestamp": "asc"}]}
        },
        "2_sort_13_range_desc": {
            "category": "Sorting", "name": "range-with-desc-sort",
            "query": {"query": {"range": {"@timestamp": {"gte": tw['large_start'], "lte": tw['large_end']}}},
                      "sort": [{"@timestamp": "desc"}]}
        },

        # Category 3: DATE HISTOGRAM (4 queries)
        "3_hist_01_hourly": {
            "category": "Date Histogram", "name": "date-histogram-hourly",
            "query": {"size": 0, "aggs": {"by_hour": {"date_histogram": {"field": "@timestamp", "calendar_interval": "hour"}}}}
        },
        "3_hist_02_hourly_with_filter": {
            "category": "Date Histogram", "name": "date-histogram-hourly-with-filter",
            "query": {"size": 0, "query": {"term": {kw1: kw1_val}},
                      "aggs": {"by_hour": {"date_histogram": {"field": "@timestamp", "calendar_interval": "hour"}}}}
        },
        "3_hist_03_minute": {
            "category": "Date Histogram", "name": "date-histogram-minute",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": tw['small_start'], "lt": tw['small_end']}}},
                      "aggs": {"by_hour": {"date_histogram": {"field": "@timestamp", "calendar_interval": "minute"}}}}
        },
        "3_hist_04_composite_daily": {
            "category": "Date Histogram", "name": "composite-date-histogram-daily",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": tw['full_start'], "lt": tw['full_end']}}},
                      "aggs": {"logs": {"composite": {"sources": [{"date": {"date_histogram": {"field": "@timestamp", "calendar_interval": "day"}}}]}}}}
        },

        # Category 4: RANGE QUERIES (10 queries)
        "4_range_01_date": {
            "category": "Range Queries", "name": "range-date",
            "query": {"query": {"range": {"@timestamp": {"gte": tw['small_start'], "lt": tw['small_end']}}}}
        },
        "4_range_02_numeric": {
            "category": "Range Queries", "name": "range-numeric",
            "query": {"query": {"range": {num1: {"gte": 20, "lte": 200}}}}
        },
        "4_range_03_conj_big_range_big_term": {
            "category": "Range Queries", "name": "range-conjunction-big-range-big-term",
            "query": {"query": {"bool": {"must": [
                {"term": {kw1: kw1_val}},
                {"range": {num1: {"gte": 1, "lte": 100}}}
            ]}}}
        },
        "4_range_04_disj_big_range_small_term": {
            "category": "Range Queries", "name": "range-disjunction-big-range-small-term",
            "query": {"query": {"bool": {"should": [
                {"term": {kw2: kw2_val}},
                {"range": {num1: {"gte": 1, "lte": 100}}}
            ]}}}
        },
        "4_range_05_small_range_only": {
            "category": "Range Queries", "name": "range-small-range-big-term",
            "query": {"query": {"bool": {"must": [{"range": {num1: {"gte": 20, "lte": 30}}}]}}}
        },
        "4_range_06_agg": {
            "category": "Range Queries", "name": "range-agg-1",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": num1,
                      "ranges": [{"to": -10}, {"from": -10, "to": 10}, {"from": 10, "to": 100},
                                 {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]}}}}
        },
        "4_range_07_agg_2": {
            "category": "Range Queries", "name": "range-agg-2",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": num1,
                      "ranges": [{"to": 100}, {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]}}}}
        },
        "4_range_08_with_metrics": {
            "category": "Range Queries", "name": "range-with-metrics",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": num1,
                      "ranges": [{"to": -10}, {"from": -10, "to": 10}, {"from": 10, "to": 100},
                                 {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]},
                      "aggs": {"tsum": {"sum": {"field": num1}}, "tmin": {"min": {"field": num1}},
                               "tavg": {"avg": {"field": num1}}, "tmax_m": {"max": {"field": num1}},
                               "tstats": {"stats": {"field": num1}}}}}}
        },
        "4_range_09_auto_date_histo": {
            "category": "Range Queries", "name": "range-auto-date-histo",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": num1,
                      "ranges": [{"to": -10}, {"from": -10, "to": 10}, {"from": 10, "to": 100},
                                 {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]},
                      "aggs": {"date": {"auto_date_histogram": {"field": "@timestamp", "buckets": 20}}}}}}
        },
        "4_range_10_auto_date_histo_with_metrics": {
            "category": "Range Queries", "name": "range-auto-date-histo-with-metrics",
            "query": {"size": 0, "aggs": {"tmax": {"range": {"field": num1,
                      "ranges": [{"to": 100}, {"from": 100, "to": 1000}, {"from": 1000, "to": 2000}, {"from": 2000}]},
                      "aggs": {"date": {"auto_date_histogram": {"field": "@timestamp", "buckets": 10},
                               "aggs": {"tmin": {"min": {"field": num1}},
                                        "tavg": {"avg": {"field": num1}},
                                        "tmax_m": {"max": {"field": num1}}}}}}}}
        },

        # Category 5: TERMS AGGREGATION (9 queries)
        "5_terms_01_keyword_terms_500": {
            "category": "Terms Aggregation", "name": "keyword-terms-500",
            "query": {"size": 0, "aggs": {"station": {"terms": {"field": kw1, "size": 500}}}}
        },
        "5_terms_02_keyword_terms_50": {
            "category": "Terms Aggregation", "name": "keyword-terms-50",
            "query": {"size": 0, "aggs": {"country": {"terms": {"field": kw1, "size": 50}}}}
        },
        "5_terms_03_significant_1": {
            "category": "Terms Aggregation", "name": "terms-significant-1",
            "query": {"track_total_hits": False, "size": 0,
                      "query": {"range": {"@timestamp": {"gte": tw['small_start'], "lt": tw['small_end']}}},
                      "aggs": {"terms": {"terms": {"field": kw1, "size": 10},
                               "aggs": {"significant_ips": {"significant_terms": {"field": kw2}}}}}}
        },
        "5_terms_04_significant_2": {
            "category": "Terms Aggregation", "name": "terms-significant-2",
            "query": {"track_total_hits": False, "size": 0,
                      "query": {"range": {"@timestamp": {"gte": tw['small_start'], "lt": tw['small_end']}}},
                      "aggs": {"terms": {"terms": {"field": kw2, "size": 10},
                               "aggs": {"significant_ips": {"significant_terms": {"field": kw1}}}}}}
        },
        "5_terms_05_multi_terms": {
            "category": "Terms Aggregation", "name": "multi-terms-keyword",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": tw['medium_start'], "lt": tw['medium_end']}}},
                      "aggs": {"important_terms": {"multi_terms": {"terms": [{"field": kw1}, {"field": kw2}]}}}}
        },
        "5_terms_06_composite": {
            "category": "Terms Aggregation", "name": "composite-terms",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": tw['small_start'], "lt": tw['small_end']}}},
                      "aggs": {"logs": {"composite": {"sources": [
                          {"kw1": {"terms": {"field": kw1, "order": "desc"}}},
                          {"kw2": {"terms": {"field": kw2, "order": "asc"}}}
                      ]}}}}
        },
        "5_terms_07_composite_3key": {
            "category": "Terms Aggregation", "name": "composite-terms-3key",
            "query": {"size": 0,
                      "query": {"range": {"@timestamp": {"gte": tw['small_start'], "lt": tw['small_end']}}},
                      "aggs": {"logs": {"composite": {"sources": [
                          {"kw1": {"terms": {"field": kw1, "order": "desc"}}},
                          {"kw2": {"terms": {"field": kw2, "order": "asc"}}},
                          {"kw3": {"terms": {"field": kw3, "order": "asc"}}}
                      ]}}}}
        },
        "5_terms_08_cardinality_low": {
            "category": "Terms Aggregation", "name": "cardinality-low",
            "query": {"size": 0, "aggs": {"region": {"cardinality": {"field": kw1}}}}
        },
        "5_terms_09_cardinality_high": {
            "category": "Terms Aggregation", "name": "cardinality-high",
            "query": {"size": 0, "aggs": {"agent": {"cardinality": {"field": kw2}}}}
        },
    }

    return queries


def run_query(session, url, index, query_body, timeout=120):
    start = time.perf_counter()
    try:
        resp = session.post(f'{url}/{index}/_search', json=query_body, timeout=timeout)
        elapsed_ms = (time.perf_counter() - start) * 1000
        if resp.status_code != 200:
            return elapsed_ms, 0, resp.status_code
        result = resp.json()
        hits = result.get('hits', {}).get('total', {})
        total = hits.get('value', 0) if isinstance(hits, dict) else (int(hits) if hits else 0)
        return elapsed_ms, total, 200
    except requests.exceptions.Timeout:
        return (time.perf_counter() - start) * 1000, 0, -1
    except Exception:
        return (time.perf_counter() - start) * 1000, 0, -2


def probe_query(session, url, index, query_body):
    times = []
    last_status = 200
    for _ in range(3):
        lat, _, status = run_query(session, url, index, query_body, timeout=120)
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
            print(f"    TIMEOUT — skipping")
            results[key] = {'name': qdef['name'], 'category': qdef['category'],
                            'timeout': True, 'tier': tier}
            continue

        if tier == 'error':
            print(f"    ERROR (status {probe_status}) — skipping")
            results[key] = {'name': qdef['name'], 'category': qdef['category'],
                            'error': True, 'status_code': probe_status, 'tier': tier}
            continue

        print(f"    Probe: {avg_probe:.0f}ms ({tier}) → {warmup}w + {iterations}i")

        print(f"    Warmup ({warmup})...", end='', flush=True)
        for _ in range(warmup):
            run_query(session, url, index, qdef['query'])
        print(" done")

        latencies = []
        hit_count = 0
        errors = 0
        timeouts = 0
        print(f"    Measuring ({iterations})", end='', flush=True)
        for i in range(iterations):
            lat, hits, status = run_query(session, url, index, qdef['query'])
            if status == 200:
                latencies.append(lat)
                hit_count = hits
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
        results[key] = stats
        print(f"    Hits: {hit_count:,} | Avg: {stats['avg']:.1f}ms | "
              f"P50: {stats['p50']:.1f}ms | P90: {stats['p90']:.1f}ms | "
              f"P95: {stats['p95']:.1f}ms | P99: {stats['p99']:.1f}ms")

    session.close()
    return results


def main():
    parser = argparse.ArgumentParser(description='Universal Big5 Benchmark')
    parser.add_argument('--url', default='http://localhost:9200')
    parser.add_argument('--index', default='big5-benchmark')
    parser.add_argument('--label', default='Engine')
    parser.add_argument('--warmup', type=int, default=100)
    parser.add_argument('--iterations', type=int, default=200)
    parser.add_argument('--output', default='/tmp/universal_benchmark_results.json')
    args = parser.parse_args()

    print("=" * 70)
    print(f"  Universal Big5 Benchmark — {args.label}")
    print("=" * 70)

    session = requests.Session()
    session.headers.update({'Content-Type': 'application/json'})

    # Discover schema
    print(f"\n  Discovering schema for {args.url}/{args.index}...")
    schema = discover_schema(session, args.url, args.index)

    if not schema['doc_count']:
        print(f"\n  ERROR: Index '{args.index}' has no documents")
        sys.exit(1)

    # Fallback: if no fields discovered, use known Big5 real data schema
    if not schema['keywords'] and not schema['text_fields']:
        print("  (No mapping returned — using known Big5 schema)")
        schema['keywords'] = ['process.name', 'cloud.region', 'aws.cloudwatch.log_stream',
                              'meta.file', 'agent.name', 'host.name']
        schema['text_fields'] = ['message', 'log.file.path']
        schema['numeric_fields'] = ['metrics.size', 'metrics.tmin']
    if not schema['ts_min'] or not schema['ts_max']:
        # Try fetching a sample doc to detect timestamp range
        try:
            r = session.post(f'{args.url}/{args.index}/_search', json={
                "size": 0, "aggs": {"by_hour": {"date_histogram": {"field": "@timestamp", "calendar_interval": "day"}}}
            }, timeout=30)
            if r.status_code == 200:
                buckets = r.json().get('aggregations', {}).get('by_hour', {}).get('buckets', [])
                if buckets:
                    schema['ts_min'] = buckets[0].get('key_as_string', '2023-01-01T00:00:00')
                    schema['ts_max'] = buckets[-1].get('key_as_string', '2023-01-14T00:00:00')
        except Exception:
            pass
    if not schema['ts_min']:
        schema['ts_min'] = '2023-01-01T00:00:00'
        schema['ts_max'] = '2023-01-14T00:00:00'
        print("  (Using default Big5 timestamp range: 2023-01-01 to 2023-01-14)")

    print(f"  Documents: {schema['doc_count']:,}")
    print(f"  Segments: {schema['segments']}")
    print(f"  Timestamp range: {schema['ts_min']} → {schema['ts_max']}")
    print(f"  Keywords: {len(schema['keywords'])}, Text: {len(schema['text_fields'])}, Numeric: {len(schema['numeric_fields'])}")

    # Build queries
    queries = build_queries(schema, session, args.url, args.index)
    session.close()

    print(f"\n  Config: up to {args.warmup} warmup + {args.iterations} measured per query")
    print(f"  Queries: {len(queries)} operations across 5 categories")

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
    print(f"  CATEGORY SUMMARY — {args.label}")
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

    # Per-query
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
        'schema': {
            'doc_count': schema['doc_count'], 'segments': schema['segments'],
            'ts_min': schema['ts_min'], 'ts_max': schema['ts_max'],
            'keyword_count': len(schema['keywords']),
            'text_count': len(schema['text_fields']),
            'numeric_count': len(schema['numeric_fields']),
        },
        'config': {'url': args.url, 'index': args.index,
                   'warmup': args.warmup, 'iterations': args.iterations},
        'total_time_sec': round(query_elapsed, 2),
        'queries': results,
    }
    with open(args.output, 'w') as f:
        json.dump(output, f, indent=2)
    print(f"  Results saved to: {args.output}")
    print(f"{'='*70}")


if __name__ == '__main__':
    main()
