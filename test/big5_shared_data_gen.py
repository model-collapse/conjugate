#!/usr/bin/env python3
"""Deterministic NDJSON generator for Big5 benchmarking.

Generates identical data for indexing into both OpenSearch and Conjugate,
ensuring apple-to-apple benchmark comparisons.

Usage:
    python3 big5_shared_data_gen.py --docs 10000000 --output /tmp/big5_10m.ndjson
    python3 big5_shared_data_gen.py --docs 1000000 > /tmp/big5_1m.ndjson
"""
import argparse
import random
import sys
import time

# --- Field value pools ---

PROCESS_NAMES: list[str] = [
    "udev", "systemd", "cron", "sshd", "nginx", "dockerd", "kubelet",
]

CLOUD_REGIONS: list[str] = [
    "us-east-1", "us-east-2", "us-west-1", "us-west-2", "ca-central-1",
    "eu-west-1", "eu-west-2", "eu-west-3", "eu-central-1", "eu-central-2",
    "eu-north-1", "eu-south-1", "eu-south-2",
    "ap-east-1", "ap-southeast-1", "ap-southeast-2", "ap-southeast-3",
    "ap-northeast-1", "ap-northeast-2", "ap-northeast-3",
    "ap-south-1", "ap-south-2",
    "sa-east-1", "me-south-1", "me-central-1", "af-south-1",
]

HOST_NAMES: list[str] = [f"host-{i:03d}" for i in range(1, 51)]

LOG_STREAMS: list[str] = [
    "foulseer", "windcatcher", "stormwarden", "ashwalker", "tidecaller",
    "duskrunner", "frostweaver", "ironbinder", "flameguard", "shadowmend",
    "leafdancer", "stonecrier", "rainsinger", "cloudspinner", "thornkeeper",
    "nightbloom", "sunforger", "mistwalker", "deepwatcher", "skychaser",
    "rootwarden", "wavecutter", "dustshaper", "gleamfinder", "hollowsong",
    "peakstrider", "shellmender", "dawnpiercer", "voidtender", "sparkwelder",
    "mosscrafter", "sandscribe", "galehowler", "emberdrift", "poolseeker",
    "ridgepacer", "coralwinder", "plumecatcher", "quartzpolish", "breezeknot",
    "cliffhanger", "marshlight", "pebbleworn", "silkthreader", "cresthunter",
    "valleyecho", "snowtracker", "ferncurler", "oakheart", "tidemark",
]

LOG_GROUPS: list[str] = [
    "/var/log/messages", "/var/log/syslog", "/var/log/auth.log",
    "/var/log/kern.log", "/var/log/daemon.log",
]

AGENT_VERSIONS: list[str] = ["8.8.0", "8.9.1", "8.10.0", "8.7.2"]

AGENT_TYPES: list[str] = ["filebeat", "metricbeat", "packetbeat", "heartbeat"]

MESSAGE_WORDS: list[str] = [
    "chanter", "scale", "dog", "hero", "server", "request", "response",
    "connection", "timeout", "database", "query", "index", "search",
    "document", "field", "mapping", "cluster", "node", "shard", "replica",
    "primary", "allocation", "routing", "filter", "aggregate", "metric",
    "pipeline", "processor", "buffer", "overflow", "retry", "backoff",
    "healthy", "degraded", "critical", "warning", "notice", "alert",
    "recovery", "snapshot", "restore", "migrate", "compact", "flush",
    "refresh", "commit", "segment", "merge", "optimize", "rebalance",
]

# Pre-compute IP octets for realistic log messages
IP_TEMPLATES: list[str] = [
    f"ip-10-{a}-{b}-{c}"
    for a in range(40, 50)
    for b in range(80, 100)
    for c in range(1, 21)
]

# --- Timestamp helpers ---

# 2024-01-01T00:00:00Z epoch ms
TS_START_MS: int = 1704067200000
# 14 days in ms (2024-01-01 to 2024-01-14T23:59:59.999Z)
TS_RANGE_MS: int = 14 * 24 * 3600 * 1000 - 1

# Days in each month (non-leap year, only need Jan)
# Pre-computed hour/min/sec divisors
_MS_PER_DAY = 86400000
_MS_PER_HOUR = 3600000
_MS_PER_MIN = 60000
_MS_PER_SEC = 1000


def _format_timestamp(epoch_ms: int) -> str:
    """Format epoch_ms to ISO 8601 string. Avoids datetime/time module overhead."""
    offset = epoch_ms - TS_START_MS
    day = offset // _MS_PER_DAY
    remainder = offset - day * _MS_PER_DAY
    hour = remainder // _MS_PER_HOUR
    remainder -= hour * _MS_PER_HOUR
    minute = remainder // _MS_PER_MIN
    remainder -= minute * _MS_PER_MIN
    sec = remainder // _MS_PER_SEC
    ms = remainder - sec * _MS_PER_SEC
    # day 0 = Jan 1, day 13 = Jan 14
    return f"2024-01-{day + 1:02d}T{hour:02d}:{minute:02d}:{sec:02d}.{ms:03d}Z"


def _format_syslog_ts(day: int, hour: int, minute: int, sec: int) -> str:
    """Format syslog-style timestamp: 'Jan DD HH:MM:SS'."""
    return f"Jan {day:02d} {hour:02d}:{minute:02d}:{sec:02d}"


# --- JSON escaping ---
# Message words are pre-validated ASCII, no escaping needed for them.
# But we still need to handle the assembled message safely.
# Since all our pool values are pure ASCII alphanumeric + hyphens + dots + slashes,
# no JSON escaping is required. We use this knowledge to skip json.dumps entirely.


def generate(num_docs: int, output_file: str | None) -> None:
    """Generate deterministic NDJSON documents."""
    rng = random.Random(42)

    # Pre-compute random choices as arrays for batch efficiency
    # We'll use rng directly for per-doc randomness

    out = open(output_file, "w", buffering=1 << 20) if output_file else sys.stdout
    write = out.write

    # Pre-size message word selection
    num_words = len(MESSAGE_WORDS)
    num_regions = len(CLOUD_REGIONS)
    num_hosts = len(HOST_NAMES)
    num_processes = len(PROCESS_NAMES)
    num_streams = len(LOG_STREAMS)
    num_groups = len(LOG_GROUPS)
    num_agent_ver = len(AGENT_VERSIONS)
    num_agent_type = len(AGENT_TYPES)
    num_ips = len(IP_TEMPLATES)

    t0 = time.monotonic()
    last_report = t0
    _randrange = rng.randrange
    _randint = rng.randint

    for i in range(num_docs):
        # Timestamp: uniform over 14-day range
        ts_offset = _randrange(TS_RANGE_MS)

        # Decompose for both ISO and syslog format (avoid formatting twice)
        day = ts_offset // _MS_PER_DAY
        remainder = ts_offset - day * _MS_PER_DAY
        hour = remainder // _MS_PER_HOUR
        remainder -= hour * _MS_PER_HOUR
        minute = remainder // _MS_PER_MIN
        remainder -= minute * _MS_PER_MIN
        sec = remainder // _MS_PER_SEC
        ms = remainder - sec * _MS_PER_SEC
        day1 = day + 1  # 1-indexed

        ts_str = f"2024-01-{day1:02d}T{hour:02d}:{minute:02d}:{sec:02d}.{ms:03d}Z"

        # Field selections
        proc = PROCESS_NAMES[_randrange(num_processes)]
        region = CLOUD_REGIONS[_randrange(num_regions)]
        host_idx = _randrange(num_hosts)
        host = HOST_NAMES[host_idx]
        agent_name = host  # 1:1 mapping
        stream = LOG_STREAMS[_randrange(num_streams)]
        group = LOG_GROUPS[_randrange(num_groups)]
        agent_ver = AGENT_VERSIONS[_randrange(num_agent_ver)]
        agent_type = AGENT_TYPES[_randrange(num_agent_type)]
        ip = IP_TEMPLATES[_randrange(num_ips)]

        # Metrics
        size = _randint(1, 5000)
        tmin = _randrange(101)  # 0-100

        # Message: 5-8 random words
        nw = _randint(4, 8)
        words = " ".join(MESSAGE_WORDS[_randrange(num_words)] for _ in range(nw))
        msg = f"{ts_str} Jan {day1:02d} {hour:02d}:{minute:02d}:{sec:02d} {ip} {proc}: {words}"

        # Build JSON line directly via f-string (all values are safe ASCII)
        write(
            f'{{"@timestamp":"{ts_str}",'
            f'"process":{{"name":"{proc}"}},'
            f'"cloud":{{"region":"{region}"}},'
            f'"host":{{"name":"{host}"}},'
            f'"message":"{msg}",'
            f'"metrics":{{"size":{size},"tmin":{tmin}}},'
            f'"log":{{"file":{{"path":"{group}/{stream}"}}}},'
            f'"agent":{{"name":"{agent_name}","type":"{agent_type}","version":"{agent_ver}"}},'
            f'"aws":{{"cloudwatch":{{"log_stream":"{stream}","log_group":"{group}"}}}}}}\n'
        )

        # Progress every 1M docs
        if (i + 1) % 1_000_000 == 0:
            now = time.monotonic()
            elapsed = now - t0
            rate = (i + 1) / elapsed
            chunk_time = now - last_report
            chunk_rate = 1_000_000 / chunk_time if chunk_time > 0 else 0
            print(
                f"  {i + 1:>12,} docs | {elapsed:7.1f}s | "
                f"avg {rate:,.0f} docs/s | chunk {chunk_rate:,.0f} docs/s",
                file=sys.stderr,
            )
            last_report = now

    # Final summary
    elapsed = time.monotonic() - t0
    rate = num_docs / elapsed if elapsed > 0 else 0
    print(
        f"Done: {num_docs:,} docs in {elapsed:.1f}s ({rate:,.0f} docs/s)",
        file=sys.stderr,
    )

    if output_file:
        out.close()


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Generate deterministic Big5 NDJSON data for benchmarking",
    )
    parser.add_argument(
        "--docs", type=int, default=10_000_000,
        help="Number of documents to generate (default: 10,000,000)",
    )
    parser.add_argument(
        "--output", "-o", type=str, default=None,
        help="Output file path (default: stdout)",
    )
    args = parser.parse_args()

    if args.docs <= 0:
        print("ERROR: --docs must be positive", file=sys.stderr)
        sys.exit(1)

    print(
        f"Generating {args.docs:,} deterministic Big5 docs (seed=42)...",
        file=sys.stderr,
    )
    generate(args.docs, args.output)


if __name__ == "__main__":
    main()
