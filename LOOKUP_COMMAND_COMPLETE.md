# Lookup Command Implementation Complete

**Date**: 2026-01-29
**Status**: ✅ Complete
**Test Results**: 10/10 tests passing (6 unit + 4 integration)

## Overview

Implemented the **lookup** command for Tier 2 PPL, enabling data enrichment by joining search results with external lookup tables.

## Syntax

```ppl
lookup <table_name> <join_field> [as <alias>] output <output_fields>
```

### Examples

```ppl
# Basic lookup
search source=orders | lookup products product_id output name, price

# With join field alias
search source=events | lookup users user_id as uid output username, email

# With output field aliases
search source=orders | lookup products product_id output name as product_name, price as cost

# In pipeline with other commands
search source=orders | lookup products product_id output name | fields order_id, name

# Multiple lookups
search source=orders | lookup products product_id output name | lookup customers customer_id output customer_name
```

## Implementation Details

### 1. Grammar Layer (`pkg/ppl/parser/PPLLexer.g4`, `PPLParser.g4`)
- Added `LOOKUP` and `OUTPUT` keywords to lexer
- Created `lookupCommand`, `lookupOutputList`, and `lookupOutputField` parser rules
- Integrated into `pipeCommand` alternatives

### 2. AST Layer (`pkg/ppl/ast/command.go`)
- Created `LookupCommand` struct with table name, join field, join field alias, and output fields
- Created `LookupOutputField` struct for representing output fields with optional aliases
- Added `VisitLookupCommand` to visitor interface

### 3. Lookup Storage (`pkg/ppl/lookup/`)

**New Package**: `pkg/ppl/lookup`

**Files**:
- `table.go`: LookupTable with hash-based indexing
- `registry.go`: Registry for managing multiple lookup tables

**Features**:
- Hash-based O(1) lookup using Go maps
- CSV file loading support
- Thread-safe operations with sync.RWMutex
- Field existence validation
- Statistics tracking

**Key Methods**:
```go
// LookupTable
func (lt *LookupTable) AddRow(key string, row map[string]interface{})
func (lt *LookupTable) Lookup(key string) (map[string]interface{}, bool)
func (lt *LookupTable) LoadFromCSV(filePath string, keyField string) error

// Registry
func (r *Registry) Register(table *LookupTable) error
func (r *Registry) Get(name string) (*LookupTable, error)
func (r *Registry) LoadFromCSV(name string, filePath string, keyField string) error
```

### 4. Logical Planning (`pkg/ppl/planner/logical_plan.go`, `builder.go`)
- Created `LogicalLookup` operator
- Implemented `buildLookupCommand` that:
  - Validates join field exists in input schema
  - Creates output schema by merging input fields with lookup output fields
  - Handles field aliasing

### 5. Analyzer (`pkg/ppl/analyzer/analyzer.go`)
- Added `analyzeLookupCommand` that:
  - Validates table name, join field, and output fields
  - Adds output fields to scope for subsequent command validation
  - Uses `FieldTypeUnknown` for lookup fields (actual types determined at runtime)

### 6. Physical Planning (`pkg/ppl/physical/physical_plan.go`, `planner.go`)
- Created `PhysicalLookup` operator
- Location: `ExecuteOnCoordinator` (lookup tables reside on coordinator)
- Conversion from LogicalLookup to PhysicalLookup

### 7. Execution (`pkg/ppl/executor/executor.go`, `lookup_operator.go`)

**Executor Updates**:
- Added `lookupRegistry` field to Executor
- Added `SetLookupRegistry()` and `GetLookupRegistry()` methods
- Integrated lookup operator into `buildOperator()` switch

**LookupOperator**:
- Implements streaming lookup enrichment
- For each input row:
  1. Extract join field value
  2. Look up row in lookup table by key
  3. Add output fields from lookup table to input row
  4. Use aliases if specified
- Graceful handling of:
  - Missing join field in input row (skip enrichment, continue)
  - No match in lookup table (skip enrichment, continue)
  - Missing output field in lookup table (skip that field, continue)

## Test Coverage

### Unit Tests (`pkg/ppl/executor/lookup_operator_test.go`)

6 tests covering:
1. ✅ BasicLookup - Standard lookup with multiple fields
2. ✅ LookupWithAliases - Output field aliasing
3. ✅ NoLookupMatch - Non-existent key handling
4. ✅ MissingJoinField - Input row without join field
5. ✅ InvalidLookupTable - Non-existent lookup table error
6. ✅ MultipleFields - Multiple output fields with mixed aliases

### Integration Tests (`pkg/ppl/integration/lookup_integration_test.go`)

4 tests covering end-to-end scenarios:
1. ✅ BasicLookupEnrichment - Parse → Analyze → Plan → Execute with lookup
2. ✅ LookupWithAliases - Join field and output field aliasing
3. ✅ LookupInPipeline - Lookup followed by fields command (validates scope updates)
4. ✅ MultipleLookups - Chained lookup commands

### Test Results

```bash
# Unit tests
$ go test ./pkg/ppl/executor -run TestLookupOperator -v
PASS: TestLookupOperator (6/6 tests)

# Integration tests
$ go test ./pkg/ppl/integration -run TestLookupCommand_Integration -v
PASS: TestLookupCommand_Integration (4/4 tests)

Total: 10/10 tests passing ✅
```

## Key Features

### 1. Hash-Based Indexing
- O(1) lookup performance using Go maps
- Key field indexed on table registration
- Efficient for large lookup tables

### 2. Field Aliasing
- Join field aliasing: `user_id as uid`
- Output field aliasing: `name as product_name`
- Supports both single and multiple field aliases

### 3. Graceful Degradation
- Missing join field → Skip enrichment, continue processing
- No lookup match → Skip enrichment, continue processing
- Missing output field → Skip that field, continue with others
- **Never fails the query** due to lookup issues

### 4. CSV Loading
```go
registry := lookup.NewRegistry(logger)
err := registry.LoadFromCSV("products", "products.csv", "product_id")
```

### 5. Thread-Safe Operations
- Concurrent lookup access supported
- RWMutex for optimal read performance
- Safe for multi-threaded query execution

### 6. Schema Integration
- Output fields added to analyzer scope
- Subsequent commands can reference lookup fields
- Type inference with `FieldTypeUnknown` (runtime resolution)

## Files Created

```
pkg/ppl/lookup/
├── table.go           - LookupTable implementation with hash indexing
└── registry.go        - Registry for managing lookup tables

pkg/ppl/executor/
├── lookup_operator.go      - Lookup operator implementation
└── lookup_operator_test.go - Unit tests (6 tests)

pkg/ppl/integration/
└── lookup_integration_test.go - Integration tests (4 tests)
```

## Files Modified

```
pkg/ppl/parser/
├── PPLLexer.g4        - Added LOOKUP, OUTPUT keywords
└── PPLParser.g4       - Added lookup grammar rules

pkg/ppl/ast/
├── command.go         - Added LookupCommand, LookupOutputField structs
└── visitor.go         - Added VisitLookupCommand to interface

pkg/ppl/parser/
└── ast_builder.go     - Implemented lookup AST building

pkg/ppl/analyzer/
└── analyzer.go        - Added analyzeLookupCommand with scope updates

pkg/ppl/planner/
├── logical_plan.go    - Added LogicalLookup operator
└── builder.go         - Implemented buildLookupCommand

pkg/ppl/physical/
├── physical_plan.go   - Added PhysicalLookup operator
└── planner.go         - Added LogicalLookup → PhysicalLookup conversion

pkg/ppl/executor/
└── executor.go        - Integrated lookup registry and operator
```

## Usage Example

```go
// Create lookup registry
registry := lookup.NewRegistry(logger)

// Load lookup table from CSV
err := registry.LoadFromCSV("products", "products.csv", "product_id")

// Or create table programmatically
table := lookup.NewLookupTable("products", logger)
table.AddRow("101", map[string]interface{}{
    "product_id": "101",
    "name":       "Laptop",
    "price":      "999.99",
})
registry.Register(table)

// Set registry on executor
executor.SetLookupRegistry(registry)

// Execute query with lookup
query := `search source=orders | lookup products product_id output name, price`
result, err := engine.Execute(query)
```

## Performance Characteristics

- **Lookup Performance**: O(1) hash table lookup
- **Memory**: O(n) where n = lookup table size
- **Streaming**: One pass through input data
- **No Buffering**: Results stream immediately
- **Thread-Safe**: Concurrent queries supported

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Missing join field in input | Skip enrichment, continue |
| No match in lookup table | Skip enrichment, continue |
| Missing output field in lookup | Skip field, continue |
| Invalid table name | Operator creation error |
| Missing table in registry | Operator creation error |

## Limitations & Future Enhancements

### Current Limitations
1. Lookup tables must fit in memory
2. Single key lookups only (no composite keys)
3. Exact match only (no fuzzy matching)
4. Static type inference (FieldTypeUnknown)

### Future Enhancements
1. **External Lookup Sources**: Support for database/API lookups
2. **Composite Keys**: Multi-field join keys
3. **Fuzzy Matching**: Approximate string matching
4. **Type Inference**: Infer actual types from lookup table schema
5. **Caching**: LRU cache for frequently accessed lookups
6. **Distributed Lookups**: Replicate lookup tables to data nodes
7. **Streaming Lookups**: Support for large lookup tables that don't fit in memory

## Tier 2 Progress

| Command | Status | Tests |
|---------|--------|-------|
| parse | ✅ Complete | 5 unit + 3 integration |
| rex | ✅ Complete | 6 unit + 3 integration |
| **lookup** | ✅ Complete | 6 unit + 4 integration |
| table | 🔄 Next | - |
| eventstats | ⏳ Pending | - |
| streamstats | ⏳ Pending | - |
| replace (command) | ⏳ Pending | - |

**Progress**: 3/7 Tier 2 commands complete (43%)

## Conclusion

The lookup command is fully implemented and tested with:
- ✅ Complete grammar and parser support
- ✅ Robust AST representation
- ✅ Semantic analysis with scope updates
- ✅ Logical and physical planning
- ✅ Efficient hash-based execution
- ✅ Comprehensive test coverage (10/10 tests passing)
- ✅ Production-ready error handling
- ✅ CSV loading support
- ✅ Field aliasing
- ✅ Pipeline integration

The implementation follows established patterns from parse and rex commands, maintaining consistency across the codebase. The lookup command is ready for production use.
