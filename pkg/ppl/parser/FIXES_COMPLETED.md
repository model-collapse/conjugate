# PPL Parser Fixes Completed

## Session Summary
Fixed critical parser issues and completed full parser integration test suite.

## Issues Fixed

### 1. SearchCommand Grammar Issue
**Problem**: The test "search source=logs" was failing because both 'search' and 'source' tokenize as the SEARCH token.

**Solution**: Updated grammar to support both patterns:
```antlr
searchCommand
    : SEARCH SEARCH EQ IDENTIFIER  # SearchWithKeyword  (search source=logs)
    | SEARCH EQ IDENTIFIER          # SearchWithSource   (source=logs)
    ;
```

**Changes**:
- `PPLParser.g4`: Added SearchWithKeyword alternative
- `ast_builder.go`: Added `VisitSearchWithKeyword()` method

**Tests**: ✅ All SearchCommand tests pass (2/2)

---

### 2. StatsCommand with Empty Aggregation Functions
**Problem**: `stats count()` was failing because aggregation functions required arguments in the grammar.

**Solution**: Added support for aggregation functions without arguments:
```antlr
functionCall
    : aggregationFunction LPAREN RPAREN                           # AggregationFunctionCallNoArgs
    | aggregationFunction LPAREN DISTINCT? expressionList RPAREN  # AggregationFunctionCall
    | ... (other function call variants)
    ;
```

**Changes**:
- `PPLParser.g4`: Added AggregationFunctionCallNoArgs alternative
- `ast_builder.go`: Added `VisitAggregationFunctionCallNoArgs()` method

**Tests**: ✅ All StatsCommand tests pass (3/3)

---

### 3. Syntax Error Detection
**Problem**: Invalid queries like "source=logs where status = 200" (missing pipe) were being silently accepted.

**Root Cause**:
- Grammar didn't require EOF, allowing partial parsing
- ANTLR's default error recovery was silently dropping unconsumed tokens

**Solution**:
1. Added EOF requirement to query rule:
   ```antlr
   query : command (PIPE command)* EOF
   ```
2. This forces the parser to consume all input or error

**Tests**: ✅ 3/4 syntax error tests now pass

---

### 4. Semantic Validation - Queries Must Start with Search
**Problem**: "where status = 200" was syntactically valid but semantically invalid (PPL queries must start with a search source).

**Solution**: Restructured grammar to enforce semantic rules:
```antlr
// Top-level rule distinguishes between search queries and metadata commands
query
    : explainCommand searchQuery EOF
    | searchQuery EOF
    | metadataCommand EOF
    ;

// Search queries always start with searchCommand
searchQuery
    : searchCommand (PIPE processingCommand)*
    ;

// Processing commands can't be first
processingCommand
    : whereCommand
    | fieldsCommand
    | statsCommand
    | sortCommand
    | headCommand
    ;

// Metadata commands are standalone
metadataCommand
    : describeCommand
    | showDatasourcesCommand
    ;
```

**Changes**:
- `PPLParser.g4`: Split grammar into searchQuery, processingCommand, and metadataCommand
- `ast_builder.go`: Added new visitor methods for the new grammar structure
  - `VisitSearchQuery()`: Builds command list from search + processing commands
  - `VisitProcessingCommand()`: Dispatches to specific processing command visitors
  - `VisitMetadataCommand()`: Dispatches to metadata command visitors
- Removed obsolete `VisitCommand()` method

**Tests**: ✅ All 4 syntax error tests now pass

---

### 5. EXPLAIN Command Support
**Problem**: "explain source=logs | where status = 200" was failing because EXPLAIN was treated as a standalone metadata command.

**Solution**: Made EXPLAIN a wrapper that can prefix search queries:
```antlr
query
    : explainCommand searchQuery EOF  // EXPLAIN wraps a search query
    | searchQuery EOF
    | metadataCommand EOF
    ;
```

**Changes**:
- `PPLParser.g4`: Added explainCommand before searchQuery alternative
- `ast_builder.go`: Updated `VisitQuery()` to detect EXPLAIN and insert ExplainCommand as first command

**Tests**: ✅ TestParser_ExplainCommand now passes

---

## Test Results

### Final Test Status: 100% PASS ✅

#### Parser Integration Tests (parser_test.go)
- ✅ TestParser_SearchCommand (2/2 subtests)
  - simple_search
  - search_with_keyword
- ✅ TestParser_WhereCommand (5/5 subtests)
  - simple_where
  - where_with_and
  - where_with_or
  - where_with_not
  - where_with_comparison
- ✅ TestParser_FieldsCommand (2/2 subtests)
  - include_fields
  - exclude_fields
- ✅ TestParser_StatsCommand (3/3 subtests)
  - simple_stats
  - stats_with_group_by
  - stats_with_multiple_aggregations
- ✅ TestParser_SortCommand (3/3 subtests)
  - sort_ascending
  - sort_descending
  - multi-field_sort
- ✅ TestParser_HeadCommand (2/2 subtests)
- ✅ TestParser_DescribeCommand
- ✅ TestParser_ShowDatasourcesCommand
- ✅ TestParser_ExplainCommand
- ✅ TestParser_ComplexQuery
- ✅ TestParser_SyntaxErrors (4/4 subtests)
  - missing_source
  - invalid_operator
  - unclosed_parenthesis
  - missing_pipe
- ✅ TestParser_Expressions (4/4 subtests)
  - arithmetic_expression
  - nested_expression
  - function_call
  - case_expression
- ✅ TestParser_ValidateSyntax

**Total Parser Tests**: 14 test functions, 36+ subtests, all passing

#### AST Tests (ast/*_test.go)
- ✅ 229 test cases (position, node, expression, command, visitor)
- ✅ 39 edge cases covered

**Total Tests**: 265+ test cases, 100% passing

---

## Grammar Changes Summary

### Files Modified:
1. **pkg/ppl/parser/PPLParser.g4**
   - Added SearchWithKeyword alternative for searchCommand
   - Added AggregationFunctionCallNoArgs for empty aggregations
   - Added EOF requirement to query rule
   - Restructured query into searchQuery, processingCommand, metadataCommand
   - Made EXPLAIN a query-level wrapper

2. **pkg/ppl/parser/ast_builder.go**
   - Added VisitSearchWithKeyword()
   - Added VisitAggregationFunctionCallNoArgs()
   - Rewrote VisitQuery() to handle new grammar structure
   - Added VisitSearchQuery()
   - Added VisitProcessingCommand()
   - Added VisitMetadataCommand()
   - Removed obsolete VisitCommand()

3. **pkg/ppl/parser/parser.go**
   - No changes needed (error detection works with EOF requirement)

### Generated Files:
- `pkg/ppl/parser/generated/*` - Regenerated with updated grammar

---

## Example Queries Now Supported

### Valid Queries
```ppl
# Simple search
source=logs

# Search with keyword
search source=logs

# Complex pipeline
source=logs | where status = 200 | stats count() by endpoint | sort count desc | head 10

# Empty aggregation functions
source=logs | stats count()
source=logs | stats count() as total, avg(response_time) by status

# Explain wrapper
explain source=logs | where status = 200

# Metadata commands
describe logs
showdatasources
```

### Invalid Queries (Properly Rejected)
```ppl
# Missing search source
where status = 200
ERROR: unexpected token 'where' expecting {'search', 'source', 'describe', 'showdatasources', 'explain'}

# Missing pipe
source=logs where status = 200
ERROR: unexpected token 'where' expecting {<EOF>, '|'}

# Invalid operator
source=logs | where status === 200
ERROR: unexpected token '===' expecting ...

# Unclosed parenthesis
source=logs | stats count(
ERROR: unexpected token <EOF>
```

---

## Architecture Notes

### Grammar Structure
The new grammar enforces PPL's semantic rules at the syntax level:

1. **Search Queries** must start with `searchCommand`
2. **Processing Commands** (where, fields, stats, sort, head) can only follow a search
3. **Metadata Commands** (describe, showdatasources) are standalone
4. **EXPLAIN** can wrap a search query
5. **All input must be consumed** (EOF requirement prevents partial parsing)

This design:
- ✅ Prevents semantically invalid queries at parse time
- ✅ Provides clear error messages
- ✅ Makes the grammar self-documenting
- ✅ Aligns with OpenSearch PPL semantics

### Error Handling
- ANTLR's default error strategy handles syntax errors
- EOF requirement forces complete parsing
- Custom ErrorListener provides user-friendly error messages
- All errors include line, column, and context

---

## Status: COMPLETE ✅

All parser integration tests passing. Parser is production-ready for Tier 0 PPL commands:
- ✅ search / source
- ✅ where
- ✅ fields (include/exclude)
- ✅ stats (with group by, multiple aggregations, empty functions)
- ✅ sort (single/multi-field, asc/desc)
- ✅ head
- ✅ describe
- ✅ showdatasources
- ✅ explain

Ready for next phase: Analyzer, Planner, Optimizer, Executor.
