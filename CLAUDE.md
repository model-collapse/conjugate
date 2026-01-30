# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CONJUGATE is a cloud-native distributed search engine providing 100% OpenSearch API compatibility while leveraging the high-performance Diagon search engine core (C++). The name stands for: **C**loud-native **O**bservability + **N**atural-language **J**oint **U**nderstanding **G**ranular search **A**nalytics **T**unable **E**ngine.

**Primary Language**: Go (with CGO bindings to C++ Diagon core)
**Architecture**: Distributed system with specialized node types (Master, Coordination, Data)

## Build Commands

### Building Binaries

```bash
# Build all binaries
make all

# Build specific components
make master        # Build master node binary (bin/conjugate-master)
make coordination  # Build coordination node binary (bin/conjugate-coordination)
make qctl          # Build CLI tool (bin/qctl)

# Build with release mode (optimized, no race detector)
BUILD_MODE=release make all
```

### Building C++ Diagon Integration

```bash
# Build Diagon C++ engine and C API wrapper
make diagon

# This will:
# 1. Build Diagon C++ core in pkg/data/diagon/upstream/build
# 2. Build C API wrapper via build_c_api.sh
# 3. Generate libdiagon.so shared library
```

### Testing

```bash
# Run all Go tests
make test-go

# Run Go tests with coverage
make coverage

# Run specific package tests
go test ./pkg/ppl/parser -v
go test ./pkg/ppl/executor -run TestFillnullOperator

# Run integration tests
make test-integration

# Run end-to-end tests
make test-e2e

# Run benchmarks
make bench
```

### Code Quality

```bash
# Run linter (uses golangci-lint with config in .golangci.yml)
make lint-go

# Format code
make fmt-go

# Run go vet
make vet
```

### Local Development

```bash
# Start local test cluster with Docker Compose
make test-cluster-up

# Stop test cluster
make test-cluster-down

# View cluster logs
make test-cluster-logs
```

## High-Level Architecture

### Node Types and Responsibilities

**Master Nodes** (pkg/master/)
- Cluster state management via Raft consensus
- Index metadata (mappings, settings, aliases)
- Shard allocation and rebalancing
- Node discovery and health monitoring
- Dual-mode support: Traditional Raft or K8S-native (Operator)

**Coordination Nodes** (pkg/coordination/)
- HTTP REST API (OpenSearch-compatible)
- Query parsing (DSL + PPL)
- Query planning and optimization
- Result aggregation from Data Nodes
- Python pipeline execution framework
- UDF registry and execution (Expression Trees, WASM, Python)

**Data Nodes** (pkg/data/)
- Shard management and storage
- Document indexing via Diagon C++ engine
- Query execution using Diagon
- Local aggregations and computations
- gRPC service for inter-node communication

### Query Processing Pipeline

```
Client Request (HTTP)
    ↓
Coordination Node REST API
    ↓
Query Parser (DSL or PPL) → AST
    ↓
Semantic Analyzer (type checking, schema validation)
    ↓
Logical Planner (build operator tree)
    ↓
Optimizer (push-down filters, projections)
    ↓
Physical Planner (execution strategy)
    ↓
Query Executor (distribute to Data Nodes via gRPC)
    ↓
Data Nodes (Diagon C++ engine executes locally)
    ↓
Result Aggregation (merge hits, aggregations)
    ↓
Response to Client
```

### Diagon Integration (CGO Architecture)

The codebase integrates with Diagon C++ search engine via CGO:

**Layer Structure** (pkg/data/diagon/):
1. **Diagon C++ Core** (upstream/ - git submodule)
   - Pure C++ search engine implementation
   - Inverted index, columnar storage, SIMD acceleration

2. **C API Layer** (upstream/src/core/include/diagon/*.h)
   - Opaque handle-based C API for language bindings
   - Exception-safe wrappers

3. **C++ Bridge** (c_api_src/)
   - Minimal bridge code for type conversion
   - Keep this layer thin

4. **Go Bindings** (*.go files)
   - CGO wrappers calling C API
   - analysis.go, bridge.go
   - Memory management with defer cleanup

**Important**: Never add Go code to Diagon upstream. Never add complex logic to c_api_src/. Prefer extending Diagon C API over adding bridge code.

### PPL (Piped Processing Language) Implementation

The PPL subsystem (pkg/ppl/) implements OpenSearch PPL with 90% coverage target:

**Components**:
- **parser/**: ANTLR4 grammar, AST generation (265+ tests passing)
- **analyzer/**: Semantic validation, type checking, schema propagation
- **planner/**: Logical plan builder (AST → operator tree)
- **optimizer/**: Query optimization rules (filter push-down, projection pruning)
- **physical/**: Physical execution plan generation
- **executor/**: Operator execution engine (70+ commands implemented)
- **functions/**: 135+ functions across 8 categories
- **translator/**: DSL ↔ PPL translation

**Tier 1 Status**: ✅ Complete (13 tasks, production-grade analytics)

### UDF (User-Defined Functions) Framework

Multi-tiered UDF execution (pkg/coordination/expressions/, pkg/wasm/):

1. **Expression Trees (80%)**: Go-native, compiled expressions
   - Fastest execution path
   - Type-safe compilation

2. **WASM (15%)**: Sandboxed compiled code
   - Wazero runtime
   - Good performance, secure

3. **Python (5%)**: Full flexibility
   - CPython embedding (future)
   - ML model integration
   - Complex transformations

### Distributed Search

**Architecture** (implemented in Phase 1):
- Coordination node queries all Data Nodes in parallel via gRPC
- Each Data Node executes locally using Diagon C++
- Results merged: global ranking, aggregation merging, pagination
- 14 aggregation types supported (12/14 exact across shards)
- Auto-discovery via Master node polling (30s interval)
- Graceful degradation with partial results

### Control Plane Modes

CONJUGATE supports two control plane architectures:

1. **Traditional Raft Mode**: Dedicated master nodes with Raft consensus
   - Works on K8S, VMs, bare metal
   - 3 master nodes can handle 1000+ data nodes

2. **K8S-Native Mode**: Kubernetes Operator + CRDs
   - Leverages K8S etcd (built-in Raft)
   - Lower cost (~$40/month vs $162/month on AWS)
   - Auto-detection in deployment scripts

## Important Conventions

### Code Organization

- `cmd/`: Entry points for binaries (main packages)
- `pkg/`: Reusable library code
- `pkg/common/`: Shared types, utilities, proto definitions
- `pkg/data/diagon/`: Diagon CGO bindings (keep bridge layer minimal)
- `test/integration/`: Integration tests requiring cluster setup
- `deployments/`: Docker and Kubernetes manifests
- `scripts/`: Build and deployment automation

### Testing Guidelines

- Tests are colocated: `foo.go` has `foo_test.go`
- Integration tests requiring Diagon: `*_integration_test.go`
- Use table-driven tests for multiple cases
- Mock external dependencies (Master, Data Nodes)
- Integration tests use real components via Docker Compose

### Error Handling

- Use custom error types in pkg/common/errors/
- Wrap errors with context: `fmt.Errorf("failed to parse query: %w", err)`
- Log errors with structured logging (zap)
- Don't panic in library code (only in main/init)

### CGO Best Practices

- Always free C memory: `defer C.free(unsafe.Pointer(cStr))`
- Convert Go strings to C: `cStr := C.CString(goStr)`
- Check for null handles before use
- Keep CGO calls in pkg/data/diagon/ only
- Never expose C types in public APIs

### Configuration

- Configuration via Viper (supports YAML, env vars, flags)
- Default config files in config/
- Environment variable prefix: `CONJUGATE_`
- Example: `CONJUGATE_MASTER_PORT=9300`

## Current Implementation Status

### Phase 1: Distributed Foundation ✅ 99% Complete
- Master node with Raft consensus
- Data node with Diagon integration (5,000 lines C++)
- Coordination node REST API
- gRPC inter-node communication
- Shard allocation (final integration pending)

### Phase 2: Query Planning ✅ Complete
- OpenSearch DSL parser
- Custom Go query planner
- Expression Trees + WASM UDF framework
- Query optimization (filter push-down, projection pruning)

### Phase 3: Infrastructure ✅ Configuration Complete
- CI/CD workflows (GitHub Actions)
- Docker images (GHCR automatic publishing)
- Kubernetes manifests
- Deployment scripts with auto-detection

### PPL Implementation
- **Tier 1**: ✅ Complete (70% function coverage, 135+ functions)
- **Tier 2**: ⏳ In Progress (advanced analytics, 20+ remaining commands)

### Aggregations: ✅ Complete
- 14 types: terms, histogram, date_histogram, range, filters, stats, extended_stats, percentiles, cardinality, avg, min, max, sum, value_count
- Distributed aggregation merging
- 12/14 types maintain exactness across shards (85.7%)

## Key Files to Understand

### Entry Points
- `cmd/master/main.go`: Master node startup
- `cmd/coordination/main.go`: Coordination node startup
- `cmd/data/main.go`: Data node startup

### Core Coordination Logic
- `pkg/coordination/coordination.go`: Main coordination logic (REST API, query routing)
- `pkg/coordination/query_service.go`: Query execution and distribution
- `pkg/coordination/executor/query_executor.go`: Distributed query execution

### PPL Implementation
- `pkg/ppl/parser/ppl_parser.g4`: ANTLR4 grammar
- `pkg/ppl/planner/builder.go`: Logical plan construction
- `pkg/ppl/executor/executor.go`: Operator execution engine
- `pkg/ppl/functions/registry.go`: Function registry (135+ functions)

### Diagon Integration
- `pkg/data/diagon/bridge.go`: Main CGO bridge (query execution)
- `pkg/data/diagon/analysis.go`: Text analysis integration
- `pkg/data/shard.go`: Shard management wrapper

### Data Node
- `pkg/data/data.go`: Data node main logic
- `pkg/data/grpc_service.go`: gRPC service for inter-node communication
- `pkg/data/udf_filter.go`: UDF execution in queries

### Master Node
- `pkg/master/master.go`: Master node logic
- `pkg/master/raft.go`: Raft consensus implementation
- `pkg/master/shard_allocator.go`: Shard allocation strategy

## Common Development Workflows

### Adding a New PPL Command

1. Update ANTLR grammar: `pkg/ppl/parser/ppl_parser.g4`
2. Regenerate parser: `cd pkg/ppl/parser && ./generate.sh`
3. Add AST node: `pkg/ppl/ast/ast.go`
4. Add analyzer validation: `pkg/ppl/analyzer/analyzer.go`
5. Add logical operator: `pkg/ppl/planner/logical_plan.go`
6. Add builder logic: `pkg/ppl/planner/builder.go`
7. Add physical operator: `pkg/ppl/physical/physical_plan.go`
8. Add executor: `pkg/ppl/executor/<command>_operator.go`
9. Wire up in executor: `pkg/ppl/executor/executor.go`
10. Add tests at each layer

### Adding a New Aggregation Type

1. Define proto message: `pkg/common/proto/aggregations.proto`
2. Implement in Coordination: `pkg/coordination/planner/aggregations.go`
3. Add Diagon integration: `pkg/data/diagon/bridge.go`
4. Implement merging logic: `pkg/coordination/executor/aggregation_merger.go`
5. Add tests: `test/integration/distributed_aggregations_test.go`

### Debugging Distributed Queries

1. Enable debug logging: `CONJUGATE_LOG_LEVEL=debug`
2. Check coordination logs for query distribution
3. Check data node logs for local execution
4. Use gRPC interceptors for request tracing
5. Verify shard routing: query Master for shard locations
6. Test with single node first, then multi-node

### Working with Diagon C++ Code

1. Diagon code lives in `pkg/data/diagon/upstream/` (git submodule)
2. Build Diagon: `cd pkg/data/diagon/upstream/build && cmake .. && make`
3. Build C API wrapper: `cd pkg/data/diagon && ./build_c_api.sh`
4. Test Go bindings: `cd pkg/data/diagon && go test -v`
5. **Never modify upstream/ directly** - contribute to Diagon repo
6. Update submodule: `git submodule update --remote`

## Performance Considerations

- **SIMD**: Diagon uses AVX2/NEON for 4-8× faster BM25 scoring
- **Compression**: 40-70% storage savings with LZ4/ZSTD
- **Skip Indexes**: 90%+ data skipping in columnar scans
- **Connection Pooling**: gRPC connections reused across queries
- **Query Caching**: 2-minute TTL cache in coordination node
- **Batch Indexing**: Use bulk API for high throughput

## Documentation References

- **Architecture**: REPOSITORY_ARCHITECTURE.md (clear boundaries between Diagon and CONJUGATE)
- **Rebranding**: NAMING.md, MIGRATION.md (renamed from Quidditch)
- **Roadmap**: README.md (18-month implementation plan)
- **PPL**: pkg/ppl/README.md (PPL implementation details)
- **Infrastructure**: PHASE3_INFRASTRUCTURE_GUIDE.md (deployment guide)
- **Control Plane**: docs/DUAL_MODE_CONTROL_PLANE.md (Raft vs K8S-native)
