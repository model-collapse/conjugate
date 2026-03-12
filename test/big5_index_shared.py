#!/usr/bin/env python3
"""Index shared Big5 NDJSON data into OpenSearch or Conjugate.

Creates explicit mapping first (so keyword fields are native keyword type,
not text+.keyword), then bulk indexes from NDJSON file.

Usage:
  # Generate data first:
  python3 test/big5_shared_data_gen.py --docs 10000000 --output /tmp/big5_shared.ndjson

  # Index into OpenSearch:
  python3 test/big5_index_shared.py --url http://localhost:9200 --data /tmp/big5_shared.ndjson

  # Index into Conjugate:
  python3 test/big5_index_shared.py --url http://localhost:9201 --data /tmp/big5_shared.ndjson
"""

import argparse
import json
import os
import sys
import time

try:
    import requests
except ImportError:
    print("ERROR: requests library required. Install with: pip3 install requests")
    sys.exit(1)

# Explicit mapping — keyword fields are keyword type (not text+.keyword).
# This ensures both engines use identical field paths in queries.
INDEX_MAPPING = {
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0,
        "refresh_interval": "30s"
    },
    "mappings": {
        "properties": {
            "@timestamp": {"type": "date"},
            "process": {
                "properties": {
                    "name": {"type": "keyword"}
                }
            },
            "cloud": {
                "properties": {
                    "region": {"type": "keyword"}
                }
            },
            "host": {
                "properties": {
                    "name": {"type": "keyword"}
                }
            },
            "agent": {
                "properties": {
                    "name": {"type": "keyword"},
                    "type": {"type": "keyword"},
                    "version": {"type": "keyword"}
                }
            },
            "aws": {
                "properties": {
                    "cloudwatch": {
                        "properties": {
                            "log_stream": {"type": "keyword"},
                            "log_group": {"type": "keyword"}
                        }
                    }
                }
            },
            "message": {"type": "text"},
            "log": {
                "properties": {
                    "file": {
                        "properties": {
                            "path": {"type": "text"}
                        }
                    }
                }
            },
            "metrics": {
                "properties": {
                    "size": {"type": "long"},
                    "tmin": {"type": "long"}
                }
            }
        }
    }
}


def delete_index(session, url, index):
    """Delete existing index if present."""
    resp = session.delete(f'{url}/{index}', timeout=30)
    if resp.status_code == 200:
        print(f"  Deleted existing index '{index}'")
    elif resp.status_code == 404:
        print(f"  Index '{index}' does not exist (OK)")
    else:
        print(f"  Delete response: {resp.status_code} {resp.text[:200]}")


def create_index(session, url, index):
    """Create index with explicit mapping."""
    resp = session.put(f'{url}/{index}', json=INDEX_MAPPING, timeout=30)
    if resp.status_code == 200:
        print(f"  Created index '{index}' with explicit mapping")
    else:
        print(f"  Create index failed: {resp.status_code} {resp.text[:300]}")
        sys.exit(1)


def bulk_index(session, url, index, data_path, batch_size):
    """Bulk index documents from NDJSON file."""
    file_size = os.path.getsize(data_path)
    print(f"  Data file: {data_path} ({file_size / (1024**3):.2f} GB)")

    start_time = time.time()
    total_indexed = 0
    batch_lines = []
    errors = 0
    last_report = start_time

    with open(data_path, 'r') as f:
        for line in f:
            line = line.strip()
            if not line:
                continue

            doc_id = str(total_indexed + len(batch_lines) // 2)
            batch_lines.append(json.dumps({"index": {"_index": index, "_id": doc_id}}))
            batch_lines.append(line)

            if len(batch_lines) >= batch_size * 2:
                body = '\n'.join(batch_lines) + '\n'
                try:
                    resp = session.post(f'{url}/_bulk', data=body, timeout=120)
                    if resp.status_code != 200:
                        errors += 1
                        print(f"\n  HTTP {resp.status_code}: {resp.text[:200]}", file=sys.stderr)
                    else:
                        result = resp.json()
                        if result.get('errors'):
                            err_items = [i for i in result.get('items', [])
                                        if 'error' in i.get('index', {})]
                            if err_items:
                                errors += len(err_items)
                except requests.exceptions.Timeout:
                    errors += 1
                    print(f"\n  Timeout at {total_indexed} docs", file=sys.stderr)
                except Exception as e:
                    errors += 1
                    print(f"\n  Error: {e}", file=sys.stderr)

                total_indexed += len(batch_lines) // 2
                batch_lines = []

                now = time.time()
                if now - last_report >= 2.0:
                    elapsed = now - start_time
                    rate = total_indexed / elapsed if elapsed > 0 else 0
                    print(f"\r  {total_indexed:>12,} docs | {rate:>8,.0f} docs/sec | {elapsed:>6.1f}s",
                          end='', flush=True)
                    last_report = now

    # Flush remaining
    if batch_lines:
        body = '\n'.join(batch_lines) + '\n'
        try:
            resp = session.post(f'{url}/_bulk', data=body, timeout=120)
            if resp.status_code == 200:
                result = resp.json()
                if result.get('errors'):
                    errors += sum(1 for i in result.get('items', [])
                                 if 'error' in i.get('index', {}))
        except Exception as e:
            errors += 1
            print(f"\n  Final batch error: {e}", file=sys.stderr)
        total_indexed += len(batch_lines) // 2

    elapsed = time.time() - start_time
    rate = total_indexed / elapsed if elapsed > 0 else 0

    print(f"\n  Indexed {total_indexed:,} docs in {elapsed:.1f}s ({rate:,.0f} docs/sec), {errors} errors")
    return total_indexed, errors


def refresh_and_verify(session, url, index):
    """Refresh index and verify doc count."""
    print("  Refreshing index...")
    try:
        session.post(f'{url}/{index}/_refresh', timeout=60)
    except Exception:
        pass

    try:
        resp = session.get(f'{url}/{index}/_count', timeout=30)
        if resp.status_code == 200:
            count = resp.json().get('count', 0)
            print(f"  Verified doc count: {count:,}")
            return count
    except Exception:
        pass
    return 0


def main():
    parser = argparse.ArgumentParser(description='Index shared Big5 data')
    parser.add_argument('--url', default='http://localhost:9200')
    parser.add_argument('--index', default='big5-benchmark')
    parser.add_argument('--data', required=True, help='Path to NDJSON file')
    parser.add_argument('--batch-size', type=int, default=5000)
    parser.add_argument('--skip-delete', action='store_true',
                       help='Skip deleting existing index')
    args = parser.parse_args()

    if not os.path.exists(args.data):
        print(f"ERROR: Data file not found: {args.data}")
        sys.exit(1)

    session = requests.Session()
    session.headers.update({'Content-Type': 'application/json'})

    print(f"\n{'='*60}")
    print(f"  Indexing shared Big5 data into {args.url}/{args.index}")
    print(f"{'='*60}\n")

    if not args.skip_delete:
        delete_index(session, args.url, args.index)

    create_index(session, args.url, args.index)

    session.headers.update({'Content-Type': 'application/x-ndjson'})
    total, errors = bulk_index(session, args.url, args.index, args.data, args.batch_size)

    session.headers.update({'Content-Type': 'application/json'})
    count = refresh_and_verify(session, args.url, args.index)

    print(f"\n{'='*60}")
    print(f"  Done: {count:,} docs indexed, {errors} errors")
    print(f"{'='*60}\n")

    session.close()


if __name__ == '__main__':
    main()
