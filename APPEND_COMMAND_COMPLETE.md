# Append Command Implementation Complete

**Date**: 2026-01-29
**Status**: ✅ Complete
**Test Results**: 9/9 tests passing (5 unit + 4 integration)

## Overview

Implemented the **append** command for Tier 2 PPL, enabling result set concatenation by appending results from a subsearch to the main query results.

## Syntax

```ppl
append [subsearch]
```

### Examples

```ppl
# Basic append
search source=logs_2024 | append [search source=logs_2023]

# Append with processing commands
search source=logs | fields timestamp, message | append [search source=archive | fields timestamp, message]

# Append with eval
search source=main | eval status="main" | append [search source=archive | eval status="archive"]

# Multiple appends (chained)
search source=logs_2024 | append [search source=logs_2023] | append [search source=logs_2022]
```

## Implementation Details

### 1. Grammar Layer (`pkg/ppl/parser/PPLLexer.g4`, `PPLParser.g4`)
- Added `APPEND` keyword to lexer
- Created `appendCommand` parser rule:
  ```antlr
  appendCommand
      : APPEND LBRACKET searchQuery RBRACKET
      ;
  ```
- Integrated into `processingCommand` alternatives

### 2. AST Layer (`pkg/ppl/ast/command.go`)
- Created `AppendCommand` struct:
  ```go
  type AppendCommand struct {
      BaseNode
      Subsearch *Query // The subsearch query to append
  }
  ```
- Added `VisitAppendCommand` to visitor interface and `BaseVisitor`
- NodeTypeAppendCommand already existed in node types

### 3. Parser Integration (`pkg/ppl/parser/ast_builder.go`)
- Implemented `VisitAppendCommand`:
  - Extracts subsearch `SearchQuery` context
  - Visits subsearch to build its AST
  - Wraps command list in a `Query` object
  - Creates `AppendCommand` node
- Added dispatch in `VisitProcessingCommand`

### 4. Analyzer (`pkg/ppl/analyzer/analyzer.go`)
- Added `analyzeAppendCommand`:
  - Validates subsearch exists and is not empty
  - Verifies subsearch starts with a `SearchCommand`
  - Note: Full subsearch analysis happens at planning time (separate source/schema)

### 5. Logical Planning (`pkg/ppl/planner/logical_plan.go`, `builder.go`)

**LogicalAppend Operator**:
```go
type LogicalAppend struct {
    Subsearch    LogicalPlan      // Plan for the subsearch query
    OutputSchema *analyzer.Schema // Unified schema (union of fields)
    Input        LogicalPlan      // Input from main query
}
```

**buildAppendCommand**:
- Validates input plan exists
- Validates subsearch structure
- Creates new `PlanBuilder` with subsearch's schema
- Builds independent subsearch plan using `Build()`
- Merges schemas (union of all fields from both queries)
- Creates `LogicalAppend` operator

### 6. Physical Planning (`pkg/ppl/physical/physical_plan.go`, `planner.go`)

**PhysicalAppend Operator**:
```go
type PhysicalAppend struct {
    Subsearch    PhysicalPlan
    OutputSchema *analyzer.Schema
    Input        PhysicalPlan
}
```

- Location: `ExecuteOnCoordinator` (both queries execute on coordinator)
- Conversion: Plans both input and subsearch using `planCoordinatorOnly()`
- Added to both planner switch statements (2 locations)

### 7. Execution (`pkg/ppl/executor/executor.go`, `append_operator.go`)

**AppendOperator**:
- Implements streaming concatenation
- State machine with `inputDone` flag
- Execution flow:
  1. Opens main input operator
  2. Streams rows from main input until exhausted (ErrNoMoreRows)
  3. When main exhausted, opens subsearch operator
  4. Streams rows from subsearch until exhausted
  5. Returns final ErrNoMoreRows

**Key Features**:
- **Lazy subsearch initialization**: Subsearch not opened until main input exhausted
- **Memory efficient**: One pass, streaming execution
- **Schema flexibility**: Handles different schemas (union of fields)
- **Proper cleanup**: Closes both input and subsearch on Close()

**Integration in Executor**:
- Added `PhysicalAppend` case to `buildOperator()`
- Builds both input and subsearch operators recursively
- Creates `AppendOperator` linking them

## Test Coverage

### Unit Tests (`pkg/ppl/executor/append_operator_test.go`)

5 tests covering operator behavior:
1. ✅ **BasicAppend** - Concatenate 2+2 rows from main and subsearch
2. ✅ **EmptyMainInput** - Main empty, subsearch has data
3. ✅ **EmptySubsearch** - Main has data, subsearch empty
4. ✅ **BothEmpty** - Both empty (immediate EOF)
5. ✅ **DifferentSchemas** - Union of fields, missing fields handled gracefully

### Integration Tests (`pkg/ppl/integration/append_integration_test.go`)

4 tests covering end-to-end scenarios:
1. ✅ **BasicAppend** - Parse → Analyze → Logical Plan → Physical Plan
2. ✅ **AppendWithProcessingCommands** - Append with fields/projection commands
3. ✅ **AppendWithEval** - Append with eval commands in both queries
4. ✅ **MultipleAppends** - Chained appends (3 sources)

### Test Results

```bash
# Unit tests
$ go test ./pkg/ppl/executor -run TestAppendOperator -v
PASS: TestAppendOperator (5/5 tests)

# Integration tests
$ go test ./pkg/ppl/integration -run TestAppendCommand_Integration -v
PASS: TestAppendCommand_Integration (4/4 tests)

Total: 9/9 tests passing ✅
```

## Key Features

### 1. Independent Subsearch Execution
- Subsearch has its own source and schema
- Subsearch analyzed and planned independently
- No shared state between main query and subsearch

### 2. Schema Union
- Output schema is union of all fields from both queries
- Fields present in one query but not the other are handled gracefully
- Type conflicts resolved (first-wins strategy)

### 3. Streaming Execution
- No buffering of all rows
- One pass through data
- Memory efficient for large result sets

### 4. Lazy Initialization
- Subsearch not opened until main input exhausted
- Saves resources if query is limited before subsearch

### 5. Proper Error Handling
- Errors from subsearch propagated correctly
- Cleanup on failure
- Context cancellation supported

## Files Created

```
pkg/ppl/executor/
├── append_operator.go      - Append operator implementation
└── append_operator_test.go - Unit tests (5 tests)

pkg/ppl/integration/
└── append_integration_test.go - Integration tests (4 tests)
```

## Files Modified

```
pkg/ppl/parser/
├── PPLLexer.g4          - Added APPEND keyword
├── PPLParser.g4         - Added appendCommand rule
├── ast_builder.go       - Implemented VisitAppendCommand
└── generated/           - Regenerated parser

pkg/ppl/ast/
├── command.go           - Added AppendCommand struct
└── visitor.go           - Added VisitAppendCommand to interface

pkg/ppl/analyzer/
└── analyzer.go          - Added analyzeAppendCommand

pkg/ppl/planner/
├── logical_plan.go      - Added LogicalAppend operator
└── builder.go           - Implemented buildAppendCommand

pkg/ppl/physical/
├── physical_plan.go     - Added PhysicalAppend operator
└── planner.go           - Added LogicalAppend → PhysicalAppend conversion (2 locations)

pkg/ppl/executor/
└── executor.go          - Integrated append operator into buildOperator()
```

## Usage Example

```ppl
# Union logs from multiple years
search source=logs_2024
| fields timestamp, level, message
| append [search source=logs_2023 | fields timestamp, level, message]
| append [search source=logs_2022 | fields timestamp, level, message]
| sort timestamp desc
| head 100

# Combine errors and warnings
search source=errors
| eval type="error"
| append [search source=warnings | eval type="warning"]
| stats count() by type

# Merge data from different indices
search source=production_logs
| append [search source=staging_logs | eval env="staging"]
| where env="staging" OR level="ERROR"
```

## Performance Characteristics

- **Time Complexity**: O(n + m) where n = main rows, m = subsearch rows
- **Space Complexity**: O(1) for streaming (no buffering)
- **Execution**: Single pass, no materialization
- **Location**: Coordinator (no distribution across shards)

## Comparison with SQL

```sql
-- SQL UNION
SELECT * FROM logs_2024
UNION ALL
SELECT * FROM logs_2023;

-- PPL equivalent
search source=logs_2024 | append [search source=logs_2023]
```

**Differences**:
- PPL `append` = SQL `UNION ALL` (keeps duplicates)
- PPL subsearch can have its own pipeline
- PPL supports nested processing in subsearch

## Limitations & Future Enhancements

### Current Limitations
1. Executes on coordinator only (no distributed execution)
2. No deduplication (like UNION ALL, not UNION)
3. Subsearch must start with `search` command
4. Schema union uses simple first-wins type resolution

### Future Enhancements
1. **Distributed Append**: Execute subsearches on data nodes
2. **Dedup Option**: `append dedup=true` for UNION-like behavior
3. **Schema Validation**: Warn on schema mismatches
4. **Append Multiple**: `append [s1], [s2], [s3]` syntax
5. **Append Stats**: Track append performance metrics
6. **Parallel Subsearches**: Execute subsearches in parallel

## Tier 2 Progress

| Command | Status | Tests |
|---------|--------|-------|
| eval | ✅ Complete | - |
| rename | ✅ Complete | - |
| replace | ✅ Complete | - |
| fillnull | ✅ Complete | - |
| parse | ✅ Complete | 5 unit + 3 integration |
| rex | ✅ Complete | 6 unit + 3 integration |
| lookup | ✅ Complete | 6 unit + 4 integration |
| **append** | ✅ **Complete** | 5 unit + 4 integration |
| join | ⏳ Pending | - |

**Progress**: 8/9 Tier 2 commands complete (89%)

## Conclusion

The append command is fully implemented and tested with:
- ✅ Complete grammar and parser support
- ✅ Robust AST representation
- ✅ Semantic analysis with independent subsearch validation
- ✅ Logical and physical planning with schema unification
- ✅ Efficient streaming execution
- ✅ Comprehensive test coverage (9/9 tests passing)
- ✅ Production-ready error handling
- ✅ Schema flexibility
- ✅ Pipeline integration

The implementation follows established patterns from other commands (parse, rex, lookup), maintaining consistency across the codebase. The append command is ready for production use and completes 8 out of 9 Tier 2 commands.
