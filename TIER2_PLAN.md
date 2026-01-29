# Tier 2 PPL Implementation Plan

**Date**: January 29, 2026
**Status**: 🚀 Ready to Start

## Executive Summary

Tier 2 adds advanced analytics capabilities with data transformation and combination operations. From the 9 planned commands, **8 are now complete**, leaving **1 final command** to implement.

## Completion Status

### ✅ Already Complete (8/9 commands - 89%)

| Command | Status | Files | Purpose |
|---------|--------|-------|---------|
| `eval` | ✅ Complete | executor/eval_operator.go | Calculate new fields |
| `rename` | ✅ Complete | executor/rename_operator.go | Rename fields |
| `replace` | ✅ Complete | executor/replace_operator.go | Substitute values |
| `fillnull` | ✅ Complete | executor/fillnull_operator.go | Handle missing values |
| **`parse`** | ✅ **Complete** | executor/parse_operator.go | Extract structured data (key=value) |
| **`rex`** | ✅ **Complete** | executor/rex_operator.go | Regex extraction with named groups |
| **`lookup`** | ✅ **Complete** | executor/lookup_operator.go, lookup/ | Reference external data/tables |
| **`append`** | ✅ **Complete** | executor/append_operator.go | Concatenate result sets |

### 🔨 To Implement (1/9 commands - 11%)

| Command | Complexity | Weeks | Priority | Purpose |
|---------|------------|-------|----------|---------|
| `join` | Very High | 6 | 1 | Combine datasets (inner/left join) |

**Total Effort**: 6 weeks for remaining command

## Implementation Order (Recommended)

### ✅ Phase 1: Text Processing & Data Enrichment (COMPLETE)

**Commands**: parse, rex, lookup, append
**Status**: ✅ All commands implemented and tested

These commands add powerful data extraction, enrichment, and combination capabilities:
- **parse**: Extract structured data from unstructured text (5 unit + 3 integration tests)
- **rex**: Regex extraction with named groups (6 unit + 3 integration tests)
- **lookup**: Reference external lookup tables (6 unit + 4 integration tests)
- **append**: Concatenate result sets from subsearches (5 unit + 4 integration tests)

**Documentation**:
- PARSE_COMMAND_COMPLETE.md
- REX_COMMAND_COMPLETE.md
- LOOKUP_COMMAND_COMPLETE.md
- APPEND_COMMAND_COMPLETE.md

#### 1. parse Command (2 weeks)

**Purpose**: Extract structured data from unstructured text
```sql
-- Extract key=value pairs
source=logs | parse message "user=(?<user>\w+) action=(?<action>\w+)"

-- Extract with pattern
source=apache | parse line "%{IP:client} - - \[%{TIMESTAMP:timestamp}\]"
```

**Implementation**:
1. **Parser** (0.5 weeks)
   - Add parse command to grammar
   - Support pattern syntax with named captures
   - Support field extraction syntax

2. **Logical Plan** (0.5 weeks)
   - LogicalParse operator
   - Pattern compilation and validation
   - Field name extraction

3. **Physical Plan & Executor** (1 week)
   - PhysicalParse operator
   - Regex pattern matching engine
   - Field extraction and type conversion
   - Handle parse failures gracefully

**Files to Create**:
- `pkg/ppl/planner/logical_plan.go` - Add LogicalParse
- `pkg/ppl/physical/physical_plan.go` - Add PhysicalParse
- `pkg/ppl/executor/parse_operator.go` - Parse implementation
- `pkg/ppl/executor/parse_operator_test.go` - Unit tests
- `pkg/ppl/integration/parse_integration_test.go` - Integration tests

#### 2. rex Command (1.5 weeks)

**Purpose**: Extract fields using regular expressions
```sql
-- Named capture groups
source=logs | rex field=message "(?<error_code>\d{3}): (?<error_msg>.*)"

-- Multiple extractions
source=events | rex "user=(?<user>\w+)" | rex "time=(?<time>\d+)"
```

**Implementation**:
1. **Parser** (0.3 weeks)
   - Add rex command to grammar
   - Support field parameter
   - Support inline patterns

2. **Logical Plan** (0.4 weeks)
   - LogicalRex operator
   - Pattern validation
   - Named group extraction

3. **Physical Plan & Executor** (0.8 weeks)
   - PhysicalRex operator
   - Regex compilation with named groups
   - Field extraction
   - Support multiple rex commands in pipeline

**Files to Create**:
- `pkg/ppl/planner/logical_plan.go` - Add LogicalRex
- `pkg/ppl/physical/physical_plan.go` - Add PhysicalRex
- `pkg/ppl/executor/rex_operator.go` - Rex implementation
- `pkg/ppl/executor/rex_operator_test.go` - Unit tests
- `pkg/ppl/integration/rex_integration_test.go` - Integration tests

### Phase 2: Data Combination (7 weeks)

**Commands**: append, join

These enable combining data from multiple sources.

#### 1. append Command (1 week)

**Purpose**: Concatenate result sets from multiple queries
```sql
-- Union two queries
source=logs_2024 | append [source=logs_2023] | stats count()

-- Append with field mapping
source=errors | append [source=warnings | rename level as severity]
```

**Implementation**:
1. **Parser** (0.3 weeks)
   - Add append command with subsearch syntax
   - Support square bracket notation

2. **Logical Plan** (0.3 weeks)
   - LogicalAppend operator
   - Subsearch plan storage
   - Schema unification

3. **Physical Plan & Executor** (0.4 weeks)
   - PhysicalAppend operator
   - Execute subsearch independently
   - Concatenate results
   - Handle schema mismatches

**Files to Create**:
- `pkg/ppl/planner/logical_plan.go` - Add LogicalAppend
- `pkg/ppl/physical/physical_plan.go` - Add PhysicalAppend
- `pkg/ppl/executor/append_operator.go` - Append implementation
- `pkg/ppl/executor/append_operator_test.go` - Unit tests

### Phase 3: Join Operations (6 weeks)

**Command**: join (most complex Tier 2 feature)

#### 2. join Command (6 weeks)

**Purpose**: Combine datasets with SQL-like joins
```sql
-- Inner join
source=orders | join type=inner user_id [source=users]

-- Left join with field selection
source=logs | join type=left left=user_id right=id [source=users | fields id, name]

-- Multiple join conditions
source=sales | join region, product_id [source=inventory]
```

**Implementation**:
1. **Parser** (1 week)
   - Add join command with complex syntax
   - Support join types (inner, left)
   - Support join conditions (single/multiple fields)
   - Support subsearch syntax
   - Support field renaming

2. **Logical Plan** (1 week)
   - LogicalJoin operator
   - Join condition storage
   - Join type specification
   - Subsearch plan integration
   - Schema inference for joined output

3. **Physical Plan** (1 week)
   - PhysicalJoin operator
   - Choose join algorithm (hash vs nested loop)
   - Memory estimation
   - Cardinality estimation

4. **Executor - Hash Join** (2 weeks)
   - Build phase: Create hash table from right side
   - Probe phase: Lookup left side in hash table
   - Support inner join
   - Support left join (with nulls)
   - Handle large right side (spill to disk if needed)
   - Memory management

5. **Testing & Optimization** (1 week)
   - Comprehensive unit tests
   - Integration tests with various join scenarios
   - Performance testing
   - Memory usage optimization

**Files to Create**:
- `pkg/ppl/planner/logical_plan.go` - Add LogicalJoin
- `pkg/ppl/physical/physical_plan.go` - Add PhysicalJoin
- `pkg/ppl/executor/join_operator.go` - Join implementation (~400 lines)
- `pkg/ppl/executor/join_hash_builder.go` - Hash table builder
- `pkg/ppl/executor/join_operator_test.go` - Unit tests
- `pkg/ppl/integration/join_integration_test.go` - Integration tests

**Complexity Factors**:
- Multiple join types (inner, left, eventually right, full)
- Efficient hash table implementation
- Memory management for large joins
- Handling duplicate join keys
- NULL handling in left joins
- Schema merging and collision handling

## Functions (30 new functions)

### JSON Functions (11 functions)

**Implementation**: 1 week

- `JSON_EXTRACT` - Extract value from JSON path
- `JSON_EXTRACT_SCALAR` - Extract scalar value
- `GET_JSON_OBJECT` - Get JSON object by key
- `JSON_ARRAY` - Create JSON array
- `JSON_OBJECT` - Create JSON object
- `JSON_ARRAY_LENGTH` - Array length
- `JSON_KEYS` - Get object keys
- `JSON_VALID` - Validate JSON
- `JSON_TYPE` - Get JSON type
- Plus 2 more utility functions

**Files**:
- `pkg/ppl/functions/json.go` - JSON function implementations
- `pkg/ppl/functions/json_test.go` - Tests

### Collection Functions (15 functions)

**Implementation**: 1.5 weeks

**Array Functions**:
- `ARRAY` - Create array
- `ARRAY_CONTAINS` - Check element presence
- `ARRAY_LENGTH` - Array length
- `ARRAY_DISTINCT` - Remove duplicates
- `ARRAY_UNION` - Union of arrays
- `ARRAY_INTERSECT` - Intersection
- `ARRAY_EXCEPT` - Difference
- `ARRAY_JOIN` - Join to string
- `ARRAY_SORT` - Sort array
- `ARRAY_REVERSE` - Reverse array

**Map Functions**:
- `MAP` - Create map
- `MAP_KEYS` - Get keys
- `MAP_VALUES` - Get values
- `MAP_CONTAINS_KEY` - Check key
- `MAP_SIZE` - Map size

**Files**:
- `pkg/ppl/functions/array.go` - Array functions
- `pkg/ppl/functions/map.go` - Map functions
- `pkg/ppl/functions/collection_test.go` - Tests

### IP & Crypto Functions (4 functions)

**Implementation**: 0.5 weeks

- `INET_ATON` - IP to number
- `INET_NTOA` - Number to IP
- `MD5` - MD5 hash
- `SHA1` - SHA1 hash

**Files**:
- `pkg/ppl/functions/net.go` - IP functions
- `pkg/ppl/functions/crypto.go` - Crypto functions

## Testing Strategy

### Unit Tests
- Each operator: 100+ lines of tests
- Each function: 5-10 test cases
- Edge cases: empty data, nulls, errors

### Integration Tests
- End-to-end pipeline tests
- Complex query combinations
- Performance benchmarks

### Test Coverage Goals
- Operators: >80% coverage
- Functions: >90% coverage
- Integration: All Tier 2 queries working

## Timeline Summary

| Phase | Duration | Commands | Functions |
|-------|----------|----------|-----------|
| Phase 1: Text Processing | 3.5 weeks | parse, rex | - |
| Phase 2: Data Combination | 3 weeks | lookup, append | - |
| Phase 3: Join Operations | 6 weeks | join | - |
| Functions Implementation | 3 weeks | - | 30 functions |
| **Total** | **15.5 weeks** | **5 commands** | **30 functions** |

## Success Criteria

### Completion Metrics
- ✅ 5 new commands fully implemented
- ✅ 30 new functions added (165 total)
- ✅ 200+ new tests passing
- ✅ Integration tests for all Tier 2 scenarios
- ✅ Documentation with 50+ examples

### Performance Targets
- Parse/rex: <10ms overhead per record
- Lookup: <1ms per lookup (hash-based)
- Append: <50ms for result concatenation
- Join: <1s for 100K x 100K inner join (in-memory)

### Quality Targets
- Test coverage: >80%
- No memory leaks
- Graceful error handling
- Production-ready code quality

## Value Proposition

**What Tier 2 Enables**:
- Advanced log parsing and field extraction
- Data enrichment with lookup tables
- Combining data from multiple sources
- SQL-like join operations
- Complex data transformations

**Use Cases**:
- Parse unstructured logs into structured fields
- Enrich event data with user/product information
- Combine historical and real-time data
- Join security events with asset inventory
- Extract fields from JSON/XML logs

---

**Next Steps**: Start with Phase 1 (parse command)
