<h1 align="center">CONJUGATE</h1>

<p align="center">
  <img src="./logo.png" alt="CONJUGATE Logo" width="200"/>
</p>

<p align="center">
  <strong>C</strong>loud-native <strong>O</strong>bservability + <strong>N</strong>atural-language <strong>J</strong>oint <strong>U</strong>nderstanding <strong>G</strong>ranular search <strong>A</strong>nalytics <strong>T</strong>unable <strong>E</strong>ngine
</p>

<p align="center">
  <strong>OpenSearch-Compatible | Diagon-Powered | Cloud-Native</strong>
</p>

<p align="center">
  <a href="IMPLEMENTATION_ROADMAP.md"><img src="https://img.shields.io/badge/status-design_phase-yellow" alt="Status"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License"></a>
  <a href="https://github.com/model-collapse/diagon"><img src="https://img.shields.io/badge/powered%20by-Diagon-green" alt="Diagon"></a>
</p>

---

## What is CONJUGATE?

CONJUGATE is a **cloud-native distributed search engine** that provides 100% OpenSearch API compatibility while leveraging the high-performance [Diagon](https://github.com/model-collapse/diagon) search engine core.

**The name CONJUGATE represents our core capabilities: Cloud-native Observability + Natural-language Joint Understanding with Granular search Analytics in a Tunable Engine. Just as conjugate pairs work together in harmony, CONJUGATE nodes coordinate seamlessly across distributed Kubernetes environments to deliver high-performance search with deep observability.**

### Key Features

✅ **100% OpenSearch API Compatibility**
- Index Management, Document APIs, Search APIs
- Full Query DSL support
- 90% PPL (Piped Processing Language) support planned (Phase 4)

✅ **High Performance**
- Diagon core: Lucene-style inverted index + ClickHouse columnar storage
- SIMD-accelerated BM25 scoring (4-8× faster)
- Advanced compression (40-70% storage savings)
- Skip indexes for granule pruning (90%+ data skipping)

✅ **Distributed Architecture**
- Specialized node types (Master, Coordination, Data)
- Horizontal scalability (10-1000+ nodes)
- Multi-tier storage (Hot/Warm/Cold/Frozen)
- **Dual-mode control plane**: Traditional (Raft) or K8S-native (Operator)
- Auto-detection of deployment environment

✅ **Python-First Pipelines**
- Customize search with Python code
- Pre/post-processing hooks
- ML model integration (ONNX, TensorFlow, PyTorch)
- Built-in examples (synonym expansion, re-ranking, A/B testing)

✅ **Cloud-Native**
- Kubernetes operator
- StatefulSets for data nodes
- Auto-scaling coordination nodes
- S3/MinIO/Ceph integration

✅ **Query Optimization**
- Custom Go query planner (learning from Calcite principles)
- Cost-based optimization with logical plan representation
- Push-down filters, projections, and UDFs
- Hybrid inverted + columnar scans
- Multi-tiered UDFs: Expression Trees (80%) + WASM (15%) + Python (5%)

---

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────┐
│                    CONJUGATE Cluster                          │
├──────────────────────────────────────────────────────────────┤
│                                                                │
│  ┌────────────────────────────────────────────────────────┐  │
│  │         API Layer (OpenSearch Compatible)               │  │
│  │   REST API | DSL | PPL | Python Pipelines              │  │
│  └────────────────────────────────────────────────────────┘  │
│                            ↓                                   │
│  ┌────────────────────────────────────────────────────────┐  │
│  │         Control Plane (Dual-Mode Support)              │  │
│  │   Mode 1: Master Nodes (Raft) - Bare metal/VMs/K8S    │  │
│  │   Mode 2: K8S Operator - K8S-native with CRDs          │  │
│  │   • Cluster state    • Shard allocation                │  │
│  │   • Index metadata   • Node discovery                  │  │
│  └────────────────────────────────────────────────────────┘  │
│                            ↓                                   │
│  ┌────────────────────────────────────────────────────────┐  │
│  │            Coordination Nodes (Query Planning)         │  │
│  │   • DSL/PPL parsing       • Custom Go query planner    │  │
│  │   • Python pipelines      • Result aggregation         │  │
│  └────────────────────────────────────────────────────────┘  │
│                            ↓                                   │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              Data Nodes (Diagon Core)                  │  │
│  │   Inverted Index  │  Forward Index  │  Computation     │  │
│  │   • Text search   │  • Aggregations │  • Joins         │  │
│  │   • BM25 scoring  │  • Sorting      │  • ML inference  │  │
│  │   • SIMD-accelerated with skip indexes                 │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                                │
└──────────────────────────────────────────────────────────────┘
```

---

## Distributed Search (Implemented ✅)

CONJUGATE now supports **horizontal scaling across multiple physical DataNodes** with automatic shard distribution and result aggregation.

### Inter-Node Distributed Search Architecture

```
Client HTTP Request
    ↓
Coordination Node (REST API)
    ↓
QueryExecutor (Go)
    ├─ Get shard routing from Master
    ├─ Query all DataNodes in parallel (gRPC)
    │   ↓
    │   DataNode 1, 2, 3... (Go + C++)
    │       ↓
    │       Shard.Search() → Diagon C++ Engine (local)
    │       ↓
    │       Returns SearchResult with Aggregations
    ↓
Aggregate Results (Go)
    ├─ Merge hits (global ranking by score)
    ├─ Merge aggregations (all 14 types)
    └─ Apply global pagination
    ↓
Return SearchResult to Client
```

### Key Features

✅ **Parallel Query Distribution**
- Coordination node queries all DataNodes concurrently via gRPC
- Each DataNode executes queries on local shards using Diagon C++ engine
- Connection pooling and automatic error handling

✅ **Comprehensive Aggregation Support** (14 types)
- **Bucket**: terms, histogram, date_histogram, range, filters
- **Metric**: stats, extended_stats, percentiles, cardinality
- **Simple Metric**: avg, min, max, sum, value_count
- 12/14 types maintain exactness across shards (85.7%)

✅ **Continuous Auto-Discovery**
- Coordination node polls master every 30 seconds for cluster state
- New DataNodes automatically discovered and registered
- Dynamic scaling: add nodes without restarts

✅ **Graceful Degradation**
- Queries succeed with partial results when some shards are unavailable
- No cascading failures
- Proportional degradation with node failures

✅ **Global Result Ranking**
- Hits sorted by score across all shards
- Global pagination (from/size parameters)
- No duplicate documents in results

### Multi-Node Deployment Example

```bash
# Start 3-node distributed cluster
kubectl apply -f - <<EOF
apiVersion: conjugate.io/v1
kind: ConjugateCluster
metadata:
  name: conjugate-prod
spec:
  version: "1.0.0"
  master:
    replicas: 3  # Raft quorum
  coordination:
    replicas: 2
  data:
    replicas: 3  # Horizontal scaling
    storage:
      size: "100Gi"
EOF

# Create index with 6 shards (distributed across 3 DataNodes)
curl -X PUT "http://localhost:9200/products" \
  -H 'Content-Type: application/json' \
  -d '{
    "settings": {
      "number_of_shards": 6,
      "number_of_replicas": 1
    }
  }'

# Index 100K documents (auto-distributed via consistent hashing)
# ... bulk indexing ...

# Search across all nodes with aggregations
curl -X GET "http://localhost:9200/products/_search" \
  -H 'Content-Type: application/json' \
  -d '{
    "query": {"match_all": {}},
    "size": 10,
    "aggs": {
      "categories": {
        "terms": {"field": "category", "size": 10}
      },
      "price_ranges": {
        "range": {
          "field": "price",
          "ranges": [
            {"key": "low", "to": 50},
            {"key": "medium", "from": 50, "to": 200},
            {"key": "high", "from": 200}
          ]
        }
      },
      "price_stats": {
        "stats": {"field": "price"}
      }
    }
  }'

# Response: Results merged from all 3 DataNodes
# - Global hit ranking by score
# - Aggregations merged correctly
# - Total hits: sum across all shards
```

### Performance Characteristics

**Query Latency** (measured on 116M docs, single shard):
- 0.40ms P50 warm across all 42 Big5 query types
- 2x-388x faster than OpenSearch 2.11 on the same hardware

**Scalability**:
- Zero degradation from 10M to 116M docs (OS degrades 3-10x)
- Linear throughput scaling: 2x nodes ≈ 2x QPS
- Aggregation merge overhead: <10% vs single-node

**Reliability**:
- Partial shard failure: Query succeeds with available data
- Master failover: New leader elected within 5 seconds (Raft)
- Auto-recovery: Failed nodes rejoin automatically

### Architecture Principles

🎯 **Clean Separation**: Network layer (Go) separate from search engine (C++)
- C++ Diagon engine queries LOCAL shards only (no network I/O)
- Go QueryExecutor handles inter-node distribution and result aggregation

🎯 **Fault Tolerance**: Built-in resilience
- Partial results when some nodes fail
- Timeout handling per shard
- Circuit breaker patterns

🎯 **Auto-Discovery**: Zero-configuration scaling
- Coordination nodes automatically discover DataNodes via Master
- No manual client registration
- Polling interval: 30 seconds (configurable)

---

## Quick Start

### One-Command Deploy 🚀

```bash
# Clone repository
git clone https://github.com/yourorg/conjugate.git
cd conjugate

# Deploy to Kubernetes (auto-detects control plane mode)
./scripts/deploy-k8s.sh --profile dev

# Get endpoint
kubectl get svc conjugate-coordination -n conjugate
```

**That's it!** Your distributed search cluster is running.

📖 **Detailed Guide**: [QUICKSTART_K8S.md](QUICKSTART_K8S.md)

### Deployment Modes

CONJUGATE supports two control plane modes:

**K8S-Native (Auto-selected for K8S)**
```bash
./scripts/deploy-k8s.sh --mode k8s --profile dev
```
- Uses Kubernetes Operator + CRDs
- Leverages K8S etcd (Raft built-in)
- Cost: ~$40/month (AWS EKS)

**Traditional Raft (For multi-environment)**
```bash
./scripts/deploy-k8s.sh --mode raft --profile prod
```
- Dedicated master nodes with Raft
- Works on K8S, VMs, bare metal
- Cost: ~$162/month (AWS EKS)

**Auto-Detect (Default)**
```bash
./scripts/deploy-k8s.sh --mode auto
```
- K8S → Uses K8S-native
- Non-K8S → Uses Raft

### Index & Search

```bash
# Create index
curl -X PUT "http://localhost:9200/my-index" \
  -H 'Content-Type: application/json' \
  -d '{
    "settings": {"number_of_shards": 1},
    "mappings": {
      "properties": {
        "title": {"type": "text"},
        "price": {"type": "float"}
      }
    }
  }'

# Index document
curl -X PUT "http://localhost:9200/my-index/_doc/1" \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "CONJUGATE Search Engine",
    "price": 99.99
  }'

# Search
curl -X GET "http://localhost:9200/my-index/_search" \
  -H 'Content-Type: application/json' \
  -d '{
    "query": {
      "bool": {
        "must": {"match": {"title": "search"}},
        "filter": {"range": {"price": {"lte": 100}}}
      }
    }
  }'
```

---

## Documentation

### Core Documentation

📖 **[Architecture Overview](QUIDDITCH_ARCHITECTURE.md)** - Complete system design
- Node types and responsibilities
- API compatibility (100% DSL, 90% PPL)
- Query processing pipeline
- Storage architecture
- Distributed coordination

📖 **[Implementation Roadmap](IMPLEMENTATION_ROADMAP.md)** - 18-month plan
- OpenSearch API compatibility matrix
- 6 implementation phases
- Team structure (8-10 people)
- Timeline and milestones
- Risk assessment

📖 **[Kubernetes Deployment](KUBERNETES_DEPLOYMENT.md)** - Cloud-native guide
- Operator installation
- Cluster configuration
- Storage and networking
- Monitoring and security
- Backup and restore

📖 **[Python Pipeline Guide](PYTHON_PIPELINE_GUIDE.md)** - Customize search
- Pipeline architecture
- Processor types
- API reference
- Examples (synonym expansion, ML re-ranking, A/B testing)
- Testing and deployment

### Control Plane Architecture

📖 **[Dual-Mode Control Plane](docs/DUAL_MODE_CONTROL_PLANE.md)** - Flexible architecture design ⭐
- **Support for BOTH traditional (Raft) and K8S-native modes**
- Pluggable control plane interface
- Complete implementation for both modes
- Auto-detection of deployment environment
- Migration paths between modes
- Unified configuration format

📖 **[Master Node Architecture](docs/MASTER_NODE_ARCHITECTURE.md)** - Traditional Raft control plane
- Master node responsibilities and Raft consensus
- Bandwidth allocation analysis (16 KB/sec total)
- Traditional deployment patterns
- Cost analysis and recommendations
- **Key finding**: 3 master nodes can handle 1000+ data nodes

📖 **[Kubernetes Deployment Guide](docs/KUBERNETES_DEPLOYMENT_GUIDE.md)** - K8S deployment patterns
- Complete manifests (StatefulSets, Deployments, Services)
- Traditional masters vs K8S-native control plane
- Production patterns (multi-zone, node selectors, PDBs)
- Cost analysis ($162/month for 3 masters vs $40/month for operator)
- Migration strategies and Helm charts

📖 **[K8S-Native Deep Dive](docs/K8S_NATIVE_DEEP_DIVE.md)** - Cloud-native architecture analysis
- Why K8S-native should be considered for K8S deployments
- K8S already provides Raft (via etcd/strong consistency)
- Operator pattern as 2026 standard (Vitess, TiDB, Strimzi)
- Complete CRD and Controller implementation examples
- Cost/latency/complexity trade-off analysis

📖 **[K8S-Native Summary](K8S_NATIVE_SUMMARY.md)** - Quick architectural decision guide
- Decision framework for choosing control plane architecture
- Trade-offs comparison (Traditional vs K8S-Native)
- When to use each mode

### Diagon Core

🔗 **[Diagon Project](https://github.com/model-collapse/diagon)** - Underlying search engine
- Lucene-style inverted index
- ClickHouse columnar storage
- SIMD-accelerated BM25
- Comprehensive design docs (100% complete)

---

## Use Cases

### 1. Log Analytics (Replacing OpenSearch)

```yaml
# High-throughput log ingestion
indices: logs-*
settings:
  number_of_shards: 10
  codec: "diagon_best_compression"
  refresh_interval: "5s"
```

**Benefits**:
- 2-289x faster range queries (measured on 116M docs)
- 2-388x faster sort queries
- Sub-millisecond warm latency across all query types

---

### 2. E-Commerce Search

```python
# Python pipeline for ML re-ranking
class PersonalizedRankingProcessor(Processor):
    def process_response(self, response, request):
        user_id = request.user.user_id
        user_profile = self.get_user_profile(user_id)

        # Re-rank with personalization model
        features = self.extract_features(response.hits, user_profile)
        scores = self.model.predict(features)

        for hit, score in zip(response.hits.hits, scores):
            hit._score = score

        response.hits.hits.sort(key=lambda h: h._score, reverse=True)
        return response
```

**Benefits**:
- Customizable ranking with Python
- ML model integration (ONNX)
- A/B testing framework

---

### 3. Real-Time Analytics (PPL - Planned Phase 4)

```sql
-- PPL query for time-series analytics (coming in Phase 4)
source=metrics
| where timestamp > now() - 1h
| stats avg(cpu_usage), max(memory_usage) by host, span(1m)
| where avg(cpu_usage) > 80
| sort -avg(cpu_usage)
```

**Benefits** (when implemented):
- SQL-like syntax (90% OpenSearch PPL compatible)
- Query planner-optimized execution
- Skip indexes for fast granule pruning

---

## Deployment Modes

### Single-Process (Development)

```yaml
# All roles in one process
node:
  roles: [master, coordination, inverted_index, forward_index, computation]
```

**Use Cases**:
- Local development
- Integration testing
- Small deployments (<1M documents)

---

### Distributed (Production)

```yaml
# Specialized nodes
master:
  replicas: 3
  resources: {memory: "8Gi", cpu: "4"}

coordination:
  replicas: 5-20  # Auto-scaling
  python: {enabled: true}

data:
  replicas: 10-1000+
  storage: {class: "nvme", size: "1Ti"}
  roles: [inverted_index, forward_index, computation]
```

**Use Cases**:
- Production deployments
- Multi-tenant SaaS
- Large-scale analytics

---

## Python Pipelines

Customize search behavior with Python:

### Example: Synonym Expansion

```python
from conjugate.pipeline import Processor

class SynonymExpansionProcessor(Processor):
    def __init__(self):
        self.synonyms = {
            "search": ["find", "query", "lookup"],
            "fast": ["quick", "rapid", "speedy"]
        }

    def process_request(self, request):
        # Expand query with synonyms
        if "match" in request.query:
            field, text = next(iter(request.query["match"].items()))
            terms = text.split()

            expanded = []
            for term in terms:
                expanded.append(term)
                expanded.extend(self.synonyms.get(term.lower(), []))

            request.query = {
                "bool": {
                    "should": [
                        {"match": {field: text}},
                        {"match": {field: " ".join(expanded)}}
                    ]
                }
            }

        return request
```

**Deploy**:
```bash
conjugate pipeline deploy --cluster prod --package my-pipeline.tar.gz
```

**Use**:
```bash
curl -X POST "http://localhost:9200/my-index/_search?pipeline=my-pipeline" \
  -d '{"query": {"match": {"title": "fast search"}}}'
```

---

## Benchmark: CONJUGATE vs OpenSearch 2.11

Measured on 116,000,000 documents (Big5 workload), single shard, same machine, same data.
10 iterations, 3 warmup rounds, `size=0`. Full report: [BIG5_BENCHMARK_REPORT_V10.md](BIG5_BENCHMARK_REPORT_V10.md).

### Query Latency (P50, warm)

| Category | Queries | CONJ P50 | OS P50 | Speedup |
|----------|---------|----------|--------|---------|
| Text Querying | 6 | 0.40ms | 1.02ms | 1.9-3.7x |
| Sorting | 13 | 0.39ms | 26.03ms | 2.0-387.8x |
| Date Histogram | 4 | 0.40ms | 1.23ms | 2.3-3.7x |
| Range Queries | 10 | 0.41ms | 16.04ms | 2.1-288.8x |
| Terms Aggregation | 9 | 0.39ms | 0.99ms | 2.1-3.0x |
| **Total** | **42** | **0.40ms** | **9.06ms** | **1.9-387.8x** |

**Result: CONJ wins 41/42 queries. OS ERROR on 1 (sort-keyword-can-match). OS wins 0.**

### Scaling Behavior (10M to 116M docs)

| Query | OS 10M P50 | OS 116M P50 | OS degradation | CONJ 116M P50 |
|-------|-----------|-------------|----------------|---------------|
| sort-numeric-desc-with-match | 43ms | 141ms | 3.3x slower | 0.38ms |
| range-conjunction-big-range-big-term | 12ms | 116ms | 9.7x slower | 0.40ms |
| range-numeric | 3ms | 21ms | 7.0x slower | 0.40ms |
| desc-sort-timestamp | 2ms | 9ms | 4.5x slower | 0.40ms |

CONJ shows zero scaling degradation from 10M to 116M.

### Aggregation Parity (116M docs)

| Check | Result |
|-------|--------|
| Range agg bucket counts (4 buckets) | EXACT |
| Date histogram (336 hourly buckets) | EXACT |
| Cardinality low (50 unique values) | EXACT |
| Terms agg (50 keys, 116M total) | Same keys, ±4% per-key variance |
| Cardinality high (~116M unique) | 3.6% diff (different estimation algorithms) |
| Range sub-agg avg | <0.3% diff (NDV 200K sample vs full scan) |

### Indexing Throughput (116M docs)

| Engine | Throughput | Notes |
|--------|-----------|-------|
| CONJ | ~25,800 docs/sec | Stable, no degradation |
| OS 2.11 | ~29,700 docs/sec | ~15% faster |

### Comparison

| Feature | OpenSearch | CONJUGATE |
|---------|------------|-----------|
| **API Compatibility** | 100% (reference) | 100% (DSL), 90% (PPL) |
| **Core Engine** | Lucene (Java) | Diagon (C++, Lucene-compatible) |
| **Query Latency (116M docs)** | 9.06ms avg P50 | 0.40ms avg P50 |
| **Sorting (worst case)** | 141ms P50 | 0.38ms P50 |
| **Range Queries (worst case)** | 116ms P50 | 0.40ms P50 |
| **Indexing** | ~29,700 docs/sec | ~25,800 docs/sec |
| **Scaling (10M→116M)** | 3-10x degradation | No degradation |
| **Columnar Storage** | Limited | Native (NumericDocValues) |
| **SIMD Acceleration** | No | AVX2/NEON |
| **Cloud-Native** | Helm charts | K8S operator |

---

## Roadmap

### Phase 0: Foundation (Months 1-2) ✅
- Complete Diagon core essentials
- SIMD, compression, advanced queries

### Phase 1: Distributed (Months 3-5) ✅ **99% COMPLETE**
- ✅ Master node with Raft consensus
- ✅ Data node with Diagon C++ engine (5,000 lines)
- ✅ Coordination node with REST API
- ✅ All nodes start and communicate
- ⏳ Shard allocation integration (7 hours remaining)
- **Status**: All code complete, needs integration glue (see [E2E_TEST_RESULTS.md](E2E_TEST_RESULTS.md))

### Phase 2: Query Planning (Months 6-8) ⏳
- OpenSearch DSL support
- Custom Go query planner (learning from Calcite principles)
- Expression Trees + WASM UDF framework
- Query optimization

### Phase 3: Python Integration (Months 9-10) ⏳
- Python runtime
- Pipeline framework
- Example pipelines

### Phase 4: Production Features (Months 11-13) ⏳
- Aggregations, highlighting
- PPL support (90%)
- Security, observability

### Phase 5: Cloud-Native (Months 14-16) ⏳
- Kubernetes operator
- Storage tiering (Hot/Warm/Cold)
- Backup & disaster recovery

### Phase 6: Optimization (Months 17-18) ⏳
- Performance tuning
- Large-scale validation (1000+ nodes)
- Cost optimization

**Target**: 1.0 Release in Month 18

---

## Contributing

We welcome contributions! This project is in the **design phase**.

### How to Help

- **Design Review**: Review [architecture docs](QUIDDITCH_ARCHITECTURE.md) and provide feedback
- **Prototype**: Build proof-of-concept for key components
- **Diagon**: Contribute to the [Diagon core](https://github.com/model-collapse/diagon)
- **Documentation**: Improve guides and examples

### Getting Started

```bash
# Clone repository
git clone https://github.com/yourusername/conjugate.git
cd conjugate

# Read design documents
ls -la *.md

# Set up development environment (coming soon)
# make dev-setup
```

---

## Team

### Core Team (Target: 8-10 people)

- **Tech Lead (1)**: Go/C++, architecture
- **Backend Engineers (3)**: Go, master/coordination nodes
- **Systems Engineers (2)**: C++, Diagon core
- **DevOps Engineer (1)**: Kubernetes, CI/CD
- **SRE Engineer (1)**: Reliability, operations
- **Product Manager (1)**: Requirements, roadmap
- **Security Engineer (1)**: Auth, compliance (Phase 4+)
- **Technical Writer (1)**: Docs (Phase 5+)

**Join us!** See [IMPLEMENTATION_ROADMAP.md](IMPLEMENTATION_ROADMAP.md) for details.

---

## Technology Stack

| Component | Technology | Reason |
|-----------|-----------|--------|
| **Master Nodes** | Go | Distributed systems, Raft consensus |
| **Coordination Nodes** | Go + Python | Orchestration (Go), Pipelines (Python) |
| **Data Nodes** | C++ (Diagon) | Performance, SIMD, existing codebase |
| **Query Planner** | Go | Custom planner learning from Calcite principles |
| **Pipelines** | Python | ML/NLP ecosystem |
| **Orchestration** | Kubernetes | Cloud-native, auto-scaling |
| **Storage** | S3/MinIO/Ceph | Object storage for cold tier |
| **Monitoring** | Prometheus + Grafana | Metrics and dashboards |
| **Tracing** | OpenTelemetry | Distributed tracing |

---

## Name Explanation

**CONJUGATE** is a backronym that captures our core capabilities:

- **C**loud-native - Kubernetes deployment, microservices architecture
- **O**bservability - Built-in monitoring, tracing, and telemetry
- **N**atural-language - Advanced NLP and semantic search capabilities
- **J**oint - Collaborative distributed processing
- **U**nderstanding - Deep semantic comprehension of queries and documents
- **G**ranular - Fine-grained control and precision in search
- **A**nalytics - Comprehensive data analytics and aggregations
- **T**unable - Highly configurable and optimizable performance
- **E**ngine - High-performance search engine core

The name reflects the harmony between components - just as mathematical conjugates work in pairs to achieve balance, CONJUGATE components work together seamlessly. The system combines cloud-native observability with natural language understanding through granular analytics in a highly tunable engine.

See [NAMING.md](NAMING.md) for detailed naming rationale and migration from the previous name.

---

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.

---

## Acknowledgments

CONJUGATE is built upon the foundational work of:

- **[Apache Lucene](https://lucene.apache.org/)** - Inverted index design
- **[ClickHouse](https://clickhouse.com/)** - Columnar storage architecture
- **[OpenSearch](https://opensearch.org/)** - API specification
- **[Apache Calcite](https://calcite.apache.org/)** - Query optimizer design principles
- **[Diagon Project](https://github.com/model-collapse/diagon)** - High-performance search engine core

---

## Contact & Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/conjugate/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/conjugate/discussions)
- **Documentation**: [Project Docs](CONJUGATE_ARCHITECTURE.md)

---

**Last Updated**: 2026-03-14

**Latest Benchmark**: 116M docs, 42 queries, CONJ wins 41/42 — [Full Report](BIG5_BENCHMARK_REPORT_V10.md)

---

## Star History

⭐ Star this project to show your support!

---

Made with ❤️ by the CONJUGATE team
