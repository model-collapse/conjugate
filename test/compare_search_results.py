#!/usr/bin/env python3
"""Compare search results between CONJUGATE and OpenSearch for all 42 Big5 queries.

Both engines have independently generated Big5 data (same schema definition,
different synthetic values). Schema discovery runs per-engine so each query
uses field names/values that actually exist in that engine's index.

Compares:
- Total hit counts (structural match)
- Top document samples (_id, _score, sort values, _source snippet)
- Aggregation bucket counts, keys, and sub-agg values
"""

import json
import sys
import time

try:
    import requests
except ImportError:
    print("ERROR: requests library required. Install with: pip3 install requests")
    sys.exit(1)


CONJ_URL = "http://localhost:9201"
OS_URL = "http://localhost:9200"
INDEX = "big5-benchmark"


# ------- schema discovery (same as benchmark script) -------

def discover_schema(session, url, index):
    info = {'ts_min': None, 'ts_max': None, 'keywords': [], 'text_fields': [],
            'numeric_fields': [], 'doc_count': 0}
    try:
        r = session.get(f'{url}/{index}/_count', timeout=10)
        info['doc_count'] = r.json().get('count', 0) if r.status_code == 200 else 0
    except Exception:
        pass

    # Timestamp range — use date_histogram day to avoid slow min/max on 116M docs
    try:
        r = session.post(f'{url}/{index}/_search', json={
            "size": 0,
            "aggs": {"by_day": {"date_histogram": {"field": "@timestamp", "calendar_interval": "day"}}}
        }, timeout=60)
        if r.status_code == 200:
            buckets = r.json().get('aggregations', {}).get('by_day', {}).get('buckets', [])
            if buckets:
                info['ts_min'] = buckets[0].get('key_as_string', '2023-01-01T00:00:00')
                info['ts_max'] = buckets[-1].get('key_as_string', '2023-01-14T00:00:00')
    except Exception:
        pass

    # Get mapping
    try:
        r = session.get(f'{url}/{index}/_mapping', timeout=10)
        if r.status_code == 200:
            mapping = r.json()
            for idx_name, idx_data in mapping.items():
                props = idx_data.get('mappings', {}).get('properties', {})
                _extract_fields(props, '', info)
    except Exception:
        pass
    return info


def _extract_fields(props, prefix, info):
    for name, field_def in props.items():
        full_name = f'{prefix}{name}' if not prefix else f'{prefix}.{name}'
        ftype = field_def.get('type', '')
        if ftype == 'keyword':
            info['keywords'].append(full_name)
        elif ftype == 'text':
            info['text_fields'].append(full_name)
            if 'fields' in field_def:
                for sub_name, sub_def in field_def['fields'].items():
                    if sub_def.get('type') == 'keyword':
                        info['keywords'].append(f'{full_name}.{sub_name}')
        elif ftype in ('long', 'integer', 'float', 'double', 'short'):
            info['numeric_fields'].append(full_name)
        if 'properties' in field_def:
            _extract_fields(field_def['properties'], full_name, info)


def compute_time_windows(ts_min_str, ts_max_str):
    from datetime import datetime, timedelta
    for fmt in ('%Y-%m-%dT%H:%M:%S.%fZ', '%Y-%m-%dT%H:%M:%SZ', '%Y-%m-%dT%H:%M:%S',
                '%Y-%m-%d %H:%M:%S', '%Y-%m-%dT%H:%M:%S.%f'):
        try:
            ts_min = datetime.strptime(ts_min_str.split('.')[0], '%Y-%m-%dT%H:%M:%S')
            break
        except ValueError:
            continue
    else:
        ts_min = datetime(2023, 1, 1)
    for fmt in ('%Y-%m-%dT%H:%M:%S.%fZ', '%Y-%m-%dT%H:%M:%SZ', '%Y-%m-%dT%H:%M:%S',
                '%Y-%m-%d %H:%M:%S', '%Y-%m-%dT%H:%M:%S.%f'):
        try:
            ts_max = datetime.strptime(ts_max_str.split('.')[0], '%Y-%m-%dT%H:%M:%S')
            break
        except ValueError:
            continue
    else:
        ts_max = datetime(2023, 1, 14)
    total_span = ts_max - ts_min
    if total_span.total_seconds() <= 0:
        total_span = timedelta(hours=1)
    return {
        'full_start': ts_min.strftime('%Y-%m-%dT%H:%M:%S'),
        'full_end': ts_max.strftime('%Y-%m-%dT%H:%M:%S'),
        'small_start': (ts_min + total_span * 0.1).strftime('%Y-%m-%dT%H:%M:%S'),
        'small_end': (ts_min + total_span * 0.2).strftime('%Y-%m-%dT%H:%M:%S'),
        'medium_start': (ts_min + total_span * 0.1).strftime('%Y-%m-%dT%H:%M:%S'),
        'medium_end': (ts_min + total_span * 0.5).strftime('%Y-%m-%dT%H:%M:%S'),
        'large_start': (ts_min + total_span * 0.05).strftime('%Y-%m-%dT%H:%M:%S'),
        'large_end': (ts_min + total_span * 0.85).strftime('%Y-%m-%dT%H:%M:%S'),
        'midpoint': (ts_min + total_span * 0.5).strftime('%Y-%m-%dT%H:%M:%S.000Z'),
    }


def discover_keyword_value(session, url, index, field):
    try:
        r = session.post(f'{url}/{index}/_search', json={
            "size": 0, "aggs": {"vals": {"terms": {"field": field, "size": 1}}}
        }, timeout=30)
        if r.status_code == 200:
            buckets = r.json().get('aggregations', {}).get('vals', {}).get('buckets', [])
            if buckets:
                return buckets[0]['key']
    except Exception:
        pass
    return None


def discover_text_terms(session, url, index, field):
    try:
        r = session.post(f'{url}/{index}/_search', json={"size": 1, "_source": [field]}, timeout=10)
        if r.status_code == 200:
            hits = r.json().get('hits', {}).get('hits', [])
            if hits:
                src = hits[0].get('_source', {})
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
                    import random
                    random.seed(42)  # Reproducible
                    if len(words) >= 3:
                        return ' '.join(random.sample(words, 3))
                    return val
    except Exception:
        pass
    return "error timeout connection"


# ------- query builder -------

def build_queries(schema, session, url, index):
    tw = compute_time_windows(schema['ts_min'], schema['ts_max'])
    kw_fields = schema['keywords'][:5]
    text_field = schema['text_fields'][0] if schema['text_fields'] else 'message'
    num_fields = schema['numeric_fields'][:3]
    kw1 = kw_fields[0] if kw_fields else 'host.name'
    kw2 = kw_fields[1] if len(kw_fields) > 1 else kw1
    kw3 = kw_fields[2] if len(kw_fields) > 2 else kw1
    num1 = num_fields[0] if num_fields else 'network.bytes'
    kw1_val = discover_keyword_value(session, url, index, kw1) or "unknown"
    kw2_val = discover_keyword_value(session, url, index, kw2) or "unknown"
    search_terms = discover_text_terms(session, url, index, text_field)

    meta = {'kw1': kw1, 'kw2': kw2, 'kw3': kw3, 'kw1_val': kw1_val, 'kw2_val': kw2_val,
            'text_field': text_field, 'num1': num1, 'search_terms': search_terms, 'tw': tw}

    queries = {
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
    return queries, meta


# ------- result extraction -------

def extract_total_hits(resp):
    hits = resp.get('hits', {}).get('total', {})
    if isinstance(hits, dict):
        return hits.get('value', 0)
    return int(hits) if hits else 0


def extract_hit_samples(resp, n=3):
    hits = resp.get('hits', {}).get('hits', [])[:n]
    samples = []
    for h in hits:
        sample = {'_id': h.get('_id', ''), '_score': h.get('_score')}
        if 'sort' in h:
            sample['sort'] = h['sort']
        src = h.get('_source', {})
        if src:
            keys = sorted(src.keys())[:6]
            sample['_source_fields'] = {k: _truncate(src[k]) for k in keys}
        else:
            sample['_source_fields'] = {}
        samples.append(sample)
    return samples


def _truncate(val, maxlen=80):
    if isinstance(val, str) and len(val) > maxlen:
        return val[:maxlen] + '...'
    if isinstance(val, dict):
        return {k: _truncate(v, 40) for k, v in list(val.items())[:4]}
    if isinstance(val, list) and len(val) > 3:
        return val[:3] + ['...']
    return val


def extract_agg_summary(resp):
    aggs = resp.get('aggregations', {})
    if not aggs:
        return None
    return _summarize_agg(aggs)


def _summarize_agg(agg_dict, depth=0):
    if depth > 3:
        return "..."
    result = {}
    for key, val in agg_dict.items():
        if isinstance(val, dict):
            if 'buckets' in val:
                buckets = val['buckets']
                if isinstance(buckets, list):
                    result[key] = {
                        'bucket_count': len(buckets),
                        'top_buckets': _summarize_buckets(buckets, depth),
                    }
                    if 'after_key' in val:
                        result[key]['after_key'] = val['after_key']
                else:
                    result[key] = {'named_buckets': {k: v.get('doc_count', 0) for k, v in buckets.items()}}
            elif 'value' in val:
                result[key] = val['value']
            elif any(k in val for k in ('count', 'min', 'max', 'avg', 'sum')):
                result[key] = {k: val.get(k) for k in ('count', 'min', 'max', 'avg', 'sum') if k in val}
            else:
                result[key] = _summarize_agg(val, depth + 1)
    return result


def _summarize_buckets(buckets, depth):
    result = []
    show = buckets[:5] if len(buckets) > 8 else buckets[:8]
    for b in show:
        entry = {}
        if 'key_as_string' in b:
            entry['key'] = b['key_as_string']
        elif 'key' in b:
            entry['key'] = b['key']
        entry['doc_count'] = b.get('doc_count', 0)
        # Sub-aggs
        for k, v in b.items():
            if k in ('key', 'key_as_string', 'doc_count', 'from', 'to', 'from_as_string', 'to_as_string'):
                continue
            if isinstance(v, dict):
                if 'value' in v:
                    entry[k] = v['value']
                elif 'buckets' in v:
                    sub = v['buckets']
                    entry[k] = f"{len(sub)} buckets" if isinstance(sub, list) else f"{len(sub)} named"
                elif any(mk in v for mk in ('count', 'min', 'max', 'avg', 'sum')):
                    entry[k] = {mk: round(v[mk], 2) if isinstance(v.get(mk), float) else v.get(mk)
                                for mk in ('count', 'min', 'max', 'avg', 'sum') if mk in v and v[mk] is not None}
        result.append(entry)
    if len(buckets) > len(show):
        result.append(f"... and {len(buckets) - len(show)} more buckets")
    return result


# ------- main -------

def main():
    session = requests.Session()
    session.headers.update({'Content-Type': 'application/json'})

    print("=" * 70)
    print("  Big5 Result Comparison: CONJUGATE vs OpenSearch")
    print("=" * 70)

    # Verify both engines
    for name, url in [("CONJUGATE", CONJ_URL), ("OpenSearch", OS_URL)]:
        try:
            r = session.get(f'{url}/{INDEX}/_count', timeout=10)
            count = r.json().get('count', 0)
            print(f"  {name} ({url}): {count:,} docs")
        except Exception as e:
            print(f"  ERROR: {name} ({url}) not reachable: {e}")
            sys.exit(1)

    # Schema discovery per engine
    print(f"\n  Discovering CONJUGATE schema...")
    conj_schema = discover_schema(session, CONJ_URL, INDEX)
    if not conj_schema['keywords']:
        conj_schema['keywords'] = ['process.name', 'cloud.region', 'aws.cloudwatch.log_stream',
                                   'meta.file', 'agent.name', 'host.name']
        conj_schema['text_fields'] = ['message', 'log.file.path']
        conj_schema['numeric_fields'] = ['metrics.size', 'metrics.tmin']
    if not conj_schema['ts_min']:
        conj_schema['ts_min'] = '2023-01-01T00:00:00'
        conj_schema['ts_max'] = '2023-01-14T00:00:00'

    print(f"  Discovering OpenSearch schema...")
    os_schema = discover_schema(session, OS_URL, INDEX)
    if not os_schema['keywords']:
        os_schema['keywords'] = ['agent.type.keyword', 'agent.version.keyword']
        os_schema['text_fields'] = ['agent.type', 'message']
        os_schema['numeric_fields'] = ['event_duration_ms']
    if not os_schema['ts_min']:
        os_schema['ts_min'] = '2026-03-10T08:05:31'
        os_schema['ts_max'] = '2026-03-10T09:05:31'

    print(f"\n  Building queries per-engine (field discovery)...")
    conj_queries, conj_meta = build_queries(conj_schema, session, CONJ_URL, INDEX)
    os_queries, os_meta = build_queries(os_schema, session, OS_URL, INDEX)

    all_results = {
        'timestamp': time.strftime('%Y-%m-%d %H:%M:%S'),
        'dataset': f'{INDEX} — 115.9M docs per engine (independently generated)',
        'conj_meta': conj_meta,
        'os_meta': os_meta,
        'queries': {}
    }

    print(f"\n  Running {len(conj_queries)} queries against both engines...\n")

    for key in sorted(conj_queries.keys()):
        cdef = conj_queries[key]
        odef = os_queries[key]
        name = cdef['name']
        cat = cdef['category']
        print(f"  [{cat}] {name} ... ", end='', flush=True)

        # Query CONJUGATE
        try:
            r = session.post(f'{CONJ_URL}/{INDEX}/_search', json=cdef['query'], timeout=120)
            conj_resp = r.json() if r.status_code == 200 else {'error': r.text[:200], 'status': r.status_code}
            conj_ms = round(r.elapsed.total_seconds() * 1000, 1)
        except Exception as e:
            conj_resp = {'error': str(e)[:200]}
            conj_ms = -1

        # Query OpenSearch
        try:
            r = session.post(f'{OS_URL}/{INDEX}/_search', json=odef['query'], timeout=120)
            os_resp = r.json() if r.status_code == 200 else {'error': r.text[:200], 'status': r.status_code}
            os_ms = round(r.elapsed.total_seconds() * 1000, 1)
        except Exception as e:
            os_resp = {'error': str(e)[:200]}
            os_ms = -1

        conj_total = extract_total_hits(conj_resp)
        os_total = extract_total_hits(os_resp)

        entry = {
            'name': name,
            'category': cat,
            'conj_query': cdef['query'],
            'os_query': odef['query'],
            'latency_ms': {'conjugate': conj_ms, 'opensearch': os_ms},
            'total_hits': {'conjugate': conj_total, 'opensearch': os_total},
        }

        # Hit samples
        conj_samples = extract_hit_samples(conj_resp, n=3)
        os_samples = extract_hit_samples(os_resp, n=3)
        if conj_samples or os_samples:
            entry['hit_samples'] = {'conjugate': conj_samples, 'opensearch': os_samples}

        # Aggregation results
        conj_aggs = extract_agg_summary(conj_resp)
        os_aggs = extract_agg_summary(os_resp)
        if conj_aggs or os_aggs:
            entry['aggregations'] = {'conjugate': conj_aggs, 'opensearch': os_aggs}

        all_results['queries'][key] = entry

        print(f"CONJ={conj_total:,}({conj_ms:.0f}ms) OS={os_total:,}({os_ms:.0f}ms)")

    # Save
    output_path = '/tmp/big5_result_comparison.json'
    with open(output_path, 'w') as f:
        json.dump(all_results, f, indent=2, default=str)

    print(f"\n{'='*70}")
    print(f"  Results saved to: {output_path}")
    print(f"{'='*70}")
    session.close()


if __name__ == '__main__':
    main()
