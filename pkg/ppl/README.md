# PPL (Piped Processing Language) Implementation

## Overview

This package implements PPL (Piped Processing Language) query processing for Quidditch, providing OpenSearch-compatible log analytics query capabilities.

## Architecture

```
Query String → Parser → AST → Analyzer → Logical Plan → Optimizer → Physical Plan → Translator → Execution
```

### Pipeline Stages

1. **Parser** (`parser/`)
   - ANTLR4-based lexer and parser
   - Converts PPL query string to Abstract Syntax Tree (AST)
   - Handles syntax errors and recovery

2. **AST** (`ast/`)
   - Node definitions for all PPL commands and expressions
   - Visitor pattern support
   - Position tracking for error reporting

3. **Analyzer** (`analyzer/`)
   - Semantic validation
   - Type checking and inference
   - Name resolution (fields, aliases)
   - Schema validation

4. **Planner** (`planner/`)
   - Converts AST to logical plan
   - Database-style relational operators
   - Logical equivalences

5. **Optimizer** (`optimizer/`)
   - Rule-based optimization
   - Filter push-down, projection pruning
   - Constant folding, expression simplification
   - Based on Apache Calcite patterns

6. **Physical Plan** (`physical/`)
   - Converts logical plan to physical execution plan
   - Algorithm selection (hash vs stream aggregation)
   - Push-down decision logic

7. **Translator** (`translator/`)
   - Translates physical plan to OpenSearch DSL
   - Handles push-down to data nodes
   - Fallback to coordinator execution

8. **Executor** (`executor/`)
   - Executes physical plan
   - Iterator-based streaming execution
   - Coordinator-side operations (eval, join, etc.)

## Directory Structure

```
pkg/ppl/
├── parser/
│   ├── PPLLexer.g4           # ANTLR4 lexer grammar
│   ├── PPLParser.g4          # ANTLR4 parser grammar
│   ├── parser.go             # Parser wrapper
│   └── error_listener.go     # Error handling
├── ast/
│   ├── node.go               # Base node interface
│   ├── command.go            # Command nodes (search, where, stats, etc.)
│   ├── expression.go         # Expression nodes
│   ├── visitor.go            # Visitor pattern
│   └── position.go           # Source position tracking
├── analyzer/
│   ├── analyzer.go           # Semantic analyzer
│   ├── type_checker.go       # Type inference and checking
│   └── scope.go              # Symbol table/scope
├── planner/
│   ├── logical_plan.go       # Logical operator definitions
│   ├── builder.go            # AST → Logical plan builder
│   └── schema.go             # Schema representation
├── optimizer/
│   ├── optimizer.go          # Optimizer interface
│   ├── rules.go              # Optimization rule definitions
│   ├── hep_planner.go        # Heuristic planner
│   └── cost.go               # Cost model
├── physical/
│   ├── physical_plan.go      # Physical operator definitions
│   ├── planner.go            # Logical → Physical conversion
│   └── pushdown.go           # Push-down decision logic
├── translator/
│   ├── translator.go         # Physical → DSL translator
│   ├── query_builder.go      # DSL query construction
│   └── agg_builder.go        # Aggregation DSL construction
├── executor/
│   ├── executor.go           # Execution engine
│   ├── iterator.go           # Iterator interface
│   └── operators.go          # Coordinator-side operators
└── README.md                 # This file
```

## Implementation Phases

### Phase 1: Tier 0 Foundation (Weeks 1-6)
**Commands:** search, where, fields, sort, head, describe, showdatasources, explain
**Functions:** 70 core functions (math, string, date, conditional, basic aggregation)

**Week 1-2: Parser**
- [x] ANTLR4 grammar for Tier 0 commands
- [x] AST node definitions
- [x] Parser wrapper and error handling
- [x] **AST unit tests (229 test cases, 39 edge cases)** ✅
- [ ] 50+ grammar tests (tests ready, requires ANTLR4 code generation)

**Week 3-4: Planner & Optimizer**
- [ ] Logical plan builder
- [ ] Basic optimization rules (FilterMerge, ProjectMerge, ReduceExpressions)
- [ ] Expression evaluator
- [ ] 100+ expression tests

**Week 5-6: Execution**
- [ ] Physical planner with push-down logic
- [ ] DSL translator for basic queries
- [ ] Coordinator execution for non-pushable operations
- [ ] Integration with Quidditch query executor
- [ ] End-to-end tests

### Phase 2: Tier 1 Analytics (Weeks 7-14)
**Commands:** +stats, +chart, +timechart, +bin, +dedup, +top, +rare
**Functions:** +65 functions (complete math/string/date, statistical aggregations)

### Phase 3: Tier 2 Advanced (Weeks 15-24)
**Commands:** +eval, +rename, +parse, +rex, +join, +lookup, +append
**Functions:** +30 functions (JSON, collections, IP, crypto)

## Usage Example

```go
import "github.com/quidditch/quidditch/pkg/ppl"

// Parse query
query := "source=logs | where status=500 | stats count() by host"
ast, err := ppl.Parse(query)

// Analyze
analyzer := ppl.NewAnalyzer(schema)
err = analyzer.Analyze(ast)

// Plan
planner := ppl.NewLogicalPlanner()
logicalPlan := planner.Build(ast)

// Optimize
optimizer := ppl.NewHepOptimizer(rules)
optimizedPlan := optimizer.Optimize(logicalPlan)

// Physical planning
physicalPlanner := ppl.NewPhysicalPlanner()
physicalPlan := physicalPlanner.Plan(optimizedPlan)

// Translate to DSL
translator := ppl.NewDSLTranslator()
dslQuery, pushable := translator.Translate(physicalPlan)

// Execute
executor := ppl.NewExecutor(client)
results, err := executor.Execute(physicalPlan, dslQuery)
```

## Testing Strategy

### Unit Tests
- Parser: Grammar coverage for all commands
- AST: Node construction and visitor pattern
- Analyzer: Type checking and validation
- Optimizer: Rule application correctness
- Translator: DSL generation accuracy

### Integration Tests
- End-to-end query execution
- Push-down vs coordinator execution
- Complex query pipelines
- Error handling and recovery

### Benchmark Tests
- Parser performance
- Optimization time
- Query latency
- Memory usage

## Dependencies

- **ANTLR4 Go Runtime**: `github.com/antlr/antlr4/runtime/Go/antlr`
- **OpenSearch Go Client**: For DSL execution
- **Quidditch Query Executor**: Integration layer

## Configuration

```go
type Config struct {
    // Parser config
    MaxQueryLength    int           // Default: 1MB
    ParseTimeout      time.Duration // Default: 5s

    // Optimizer config
    OptimizationLevel int           // 0=none, 1=basic, 2=full
    MaxOptimizeRounds int           // Default: 10

    // Execution config
    QueryTimeout      time.Duration // Default: 30s
    MaxMemoryPerQuery int64         // Default: 500MB

    // Push-down config
    EnablePushDown    bool          // Default: true
    MaxPushDownSize   int           // Default: 10000 docs
}
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Parse Time | <1ms | Simple queries |
| Parse Time | <10ms | Complex queries |
| Optimization Time | <5ms | Tier 1 rules |
| Query Latency (push-down) | <100ms | Simple aggregations |
| Query Latency (coordinator) | <500ms | Joins, transforms |
| Memory Usage | <500MB | Per query |

## Error Handling

PPL errors are categorized:

```go
type ErrorType int

const (
    SyntaxError      ErrorType = iota // Parse errors
    SemanticError                     // Type errors, unknown fields
    ExecutionError                    // Runtime errors
    TimeoutError                      // Query timeout
    MemoryLimitError                  // Memory exceeded
)
```

## Observability

Metrics exposed:
- `ppl_queries_total{command}` - Total queries by command
- `ppl_parse_duration_seconds` - Parse time histogram
- `ppl_optimization_duration_seconds` - Optimization time
- `ppl_execution_duration_seconds{pushdown}` - Query latency
- `ppl_pushdown_rate` - Percentage of queries using push-down

## Development Guidelines

1. **Grammar Changes**: Update ANTLR4 files, regenerate parser
2. **New Commands**: Add AST node, planner case, physical operator
3. **New Functions**: Add to expression evaluator with tests
4. **Optimization Rules**: Document rule with before/after examples
5. **Testing**: Every component needs >80% test coverage

## References

- [PPL Tier Plan](../../design/PPL_TIER_PLAN.md)
- [PPL Architecture Design](../../design/PPL_ARCHITECTURE_DESIGN.md)
- [PPL Push-Down Research](../../design/PPL_PUSHDOWN_AND_OPTIMIZATION_RESEARCH.md)
- [OpenSearch PPL Documentation](https://opensearch.org/docs/latest/search-plugins/sql/ppl/)

---

**Status**: In Development (Tier 0)
**Last Updated**: January 28, 2026
