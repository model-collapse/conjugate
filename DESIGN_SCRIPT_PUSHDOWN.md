# Script Pushdown Design Analysis

**Date**: January 28, 2026
**Status**: Design Review

## Current Implementation

### Push-Down Decision Logic

**Location**: `pkg/ppl/physical/planner.go`

#### Filters
```go
case *ast.FunctionCall:
    // Some functions might not be supported in DSL
    // For now, allow simple functions
    return false // Be conservative
```

**Decision**: Function calls in filters are **NOT pushed down** to OpenSearch.

#### Projections
```go
func (pp *PhysicalPlanner) canPushDownProject(project *planner.LogicalProject) bool {
    // Can push down if all fields are simple field references
    for _, field := range project.Fields {
        if _, ok := field.(*ast.FieldReference); !ok {
            // Complex expression - can't push down
            return false
        }
    }
    return true
}
```

**Decision**: Computed fields (expressions) are **NOT pushed down** to OpenSearch.

### DSL Translator Capabilities

**Location**: `pkg/ppl/dsl/`

**Current Support**:
- ✅ Simple field references in filters
- ✅ Literal values
- ✅ Binary operators (=, !=, >, <, >=, <=, AND, OR, NOT)
- ✅ LIKE (wildcard queries)
- ✅ IN (terms queries)
- ❌ **NO Painless script generation**
- ❌ **NO function call translation**

## Query Examples

### Example 1: Filter with Function Call

**PPL Query**:
```sql
source=logs | where abs(latency) > 100
```

**Current Physical Plan**:
```
PhysicalFilter(abs(latency) > 100) [Coordinator]
  PhysicalScan(logs) [DataNode]
```

**Current DSL** (from scan only):
```json
{
  "query": {"match_all": {}}
}
```

**Result**:
- OpenSearch returns ALL documents
- Coordinator filters with `abs(latency) > 100`
- **High data transfer** if most documents don't match

### Example 2: Computed Field in Projection

**PPL Query**:
```sql
source=logs | fields status, latency, latency * 2 as double_latency
```

**Current Physical Plan**:
```
PhysicalProject(status, latency, latency * 2 as double_latency) [Coordinator]
  PhysicalScan(logs) [DataNode]
```

**Current DSL** (from scan only):
```json
{
  "_source": ["status", "latency"]
}
```

**Result**:
- OpenSearch returns status, latency
- Coordinator computes `latency * 2`
- **Minimal overhead** (computation is cheap)

### Example 3: Complex Aggregation

**PPL Query**:
```sql
source=logs | stats avg(latency * 1000) as avg_ms by status
```

**Current Physical Plan**:
```
PhysicalAggregate(avg(latency * 1000) as avg_ms, group_by=[status]) [Coordinator]
  PhysicalScan(logs) [DataNode]
```

**Current DSL** (from scan only):
```json
{
  "query": {"match_all": {}}
}
```

**Result**:
- OpenSearch returns ALL documents
- Coordinator performs aggregation
- **Very high data transfer** (aggregations typically filter dramatically)

## Design Alternatives

### Option 1: Current Design (No Script Pushdown) ✅ IMPLEMENTED

**Approach**: Run all functions and computed expressions on coordinator.

**Pros**:
- ✅ Simple to implement and maintain
- ✅ Guaranteed correctness (all functions work)
- ✅ No security concerns (no script injection)
- ✅ Easy to debug (all logic in Quidditch)
- ✅ Consistent behavior across queries

**Cons**:
- ❌ High network transfer for selective filters
- ❌ Can't leverage OpenSearch compute
- ❌ Coordinator becomes bottleneck for large result sets
- ❌ Inefficient for aggregations with computed fields

**Best For**:
- Tier 0 MVP (current phase)
- Small to medium datasets
- Queries where most data is needed anyway

### Option 2: Painless Script Pushdown (Full)

**Approach**: Translate all PPL expressions to Painless scripts and push to OpenSearch.

**Implementation**:
```json
{
  "query": {
    "script": {
      "script": {
        "source": "Math.abs(doc['latency'].value) > params.threshold",
        "params": {"threshold": 100}
      }
    }
  }
}
```

**Pros**:
- ✅ Minimal data transfer (filter at source)
- ✅ Leverages OpenSearch compute
- ✅ Efficient for selective filters
- ✅ Good for computed aggregations

**Cons**:
- ❌ Complex PPL → Painless translation
- ❌ Not all PPL functions map to Painless
- ❌ Painless scripts can be slow (interpreted)
- ❌ Security risk (script injection if not careful)
- ❌ Debugging is harder (logic split across systems)
- ❌ OpenSearch script compilation overhead

**Challenges**:
1. **Function Mapping**: PPL functions → Painless equivalents
   ```
   PPL: abs(x)        → Painless: Math.abs(x)
   PPL: upper(s)      → Painless: s.toUpperCase()
   PPL: substr(s,i,n) → Painless: s.substring(i, i+n)
   ```

2. **Field Access**: Different syntax
   ```
   PPL: latency       → Painless: doc['latency'].value
   PPL: nested.field  → Painless: doc['nested.field'].value
   ```

3. **Type Handling**: Painless is strongly typed
   ```
   PPL: auto-coercion → Painless: explicit casts needed
   ```

4. **Missing Fields**: Painless throws exceptions
   ```
   PPL: null handling → Painless: try/catch or containsKey checks
   ```

### Option 3: Selective Pushdown (Hybrid) 🎯 RECOMMENDED

**Approach**: Push down only safe, high-value operations as Painless scripts.

**Push Down** (to OpenSearch with scripts):
- ✅ Simple math functions in filters (abs, round, floor, ceil)
- ✅ String functions in filters (upper, lower, trim)
- ✅ Date functions (year, month, day)
- ✅ Computed fields in aggregations (high value)

**Don't Push Down** (run on coordinator):
- ❌ Complex string operations (regex, split)
- ❌ Nested function calls
- ❌ User-defined functions
- ❌ Functions with side effects

**Decision Criteria**:
```go
func (pp *PhysicalPlanner) canPushDownAsScript(expr ast.Expression) bool {
    switch e := expr.(type) {
    case *ast.FunctionCall:
        // Whitelist of safe, common functions
        safeFunctions := []string{
            "abs", "round", "floor", "ceil",    // Math
            "upper", "lower", "trim",           // String
            "year", "month", "day",             // Date
        }
        if contains(safeFunctions, strings.ToLower(e.Name)) {
            // Check all arguments are pushable
            for _, arg := range e.Arguments {
                if !pp.canPushDownAsScript(arg) {
                    return false
                }
            }
            return true
        }
        return false

    case *ast.FieldReference, *ast.Literal:
        return true

    case *ast.BinaryExpression:
        // Allow arithmetic and comparisons
        if e.Operator in ["+", "-", "*", "/", "%", ">", "<", ">=", "<=", "=", "!="] {
            return pp.canPushDownAsScript(e.Left) && pp.canPushDownAsScript(e.Right)
        }
        return false

    default:
        return false
    }
}
```

**Implementation Example**:
```go
// In dsl/query_builder.go
func (qb *QueryBuilder) buildScriptFilter(expr ast.Expression) (map[string]interface{}, error) {
    script := qb.generatePainlessScript(expr)

    return map[string]interface{}{
        "script": map[string]interface{}{
            "script": map[string]interface{}{
                "source": script.source,
                "lang":   "painless",
                "params": script.params,
            },
        },
    }, nil
}

func (qb *QueryBuilder) generatePainlessScript(expr ast.Expression) *PainlessScript {
    switch e := expr.(type) {
    case *ast.FunctionCall:
        switch strings.ToLower(e.Name) {
        case "abs":
            arg := qb.generatePainlessExpr(e.Arguments[0])
            return &PainlessScript{
                source: fmt.Sprintf("Math.abs(%s)", arg),
            }
        case "upper":
            arg := qb.generatePainlessExpr(e.Arguments[0])
            return &PainlessScript{
                source: fmt.Sprintf("%s.toUpperCase()", arg),
            }
        // ... more functions
        }

    case *ast.FieldReference:
        return &PainlessScript{
            source: fmt.Sprintf("doc['%s'].value", e.Name),
        }

    // ... more cases
    }
}
```

**Pros**:
- ✅ Balances performance and complexity
- ✅ High-value operations are optimized
- ✅ Reduces security risk (whitelist approach)
- ✅ Easier to debug than full pushdown
- ✅ Graceful degradation (fallback to coordinator)

**Cons**:
- ⚠️ More complex than Option 1
- ⚠️ Still requires Painless script generation
- ⚠️ Need to maintain function whitelist

### Option 4: Script Fields (Computed Projections Only)

**Approach**: Use OpenSearch script_fields for computed projections only, not filters.

**Implementation**:
```json
{
  "query": {"match_all": {}},
  "_source": ["status", "latency"],
  "script_fields": {
    "double_latency": {
      "script": {
        "source": "doc['latency'].value * 2"
      }
    }
  }
}
```

**Pros**:
- ✅ Good for computed fields in results
- ✅ Doesn't affect filtering performance
- ✅ Simple use case (just projections)

**Cons**:
- ❌ Doesn't help with filters (main bottleneck)
- ❌ Doesn't help with aggregations
- ❌ Still requires Painless generation

## Performance Analysis

### Data Transfer Comparison

**Scenario**: 1M documents, filter selects 1K (0.1%)

| Approach | Data Transferred | Computation | Notes |
|----------|------------------|-------------|-------|
| No Pushdown | 1M documents | Coordinator | Current implementation |
| Script Filter | 1K documents | OpenSearch | 1000x less transfer |
| Script Projection | 1M documents | OpenSearch | No benefit for filters |

**Winner**: Script Filter Pushdown (Option 2 or 3)

### Computation Comparison

**Scenario**: Aggregation with computed field `stats avg(latency * 1000) by status`

| Approach | Documents Moved | Aggregation Location | Cost |
|----------|-----------------|---------------------|------|
| No Pushdown | 1M | Coordinator | High transfer + compute |
| Script Agg | 0 (aggs only) | OpenSearch | Low transfer, OpenSearch compute |

**Winner**: Script Aggregation (Option 2 or 3)

## Security Considerations

### Script Injection Risk

**Problem**: User input in scripts could allow code execution.

**Mitigation**:
1. **Parameterization**: Never concatenate user input into scripts
   ```go
   // BAD
   script := fmt.Sprintf("doc['field'].value > %s", userInput)

   // GOOD
   script := "doc['field'].value > params.threshold"
   params := map[string]interface{}{"threshold": userInput}
   ```

2. **Whitelisting**: Only allow known-safe functions
3. **Validation**: Validate all expressions before script generation
4. **Sandboxing**: Rely on OpenSearch's Painless sandbox

### Resource Limits

**Problem**: Scripts can consume excessive CPU/memory.

**Mitigation**:
1. **Script Compilation Cache**: OpenSearch caches compiled scripts
2. **Circuit Breaker**: OpenSearch has built-in script execution limits
3. **Timeout**: Set script execution timeout
4. **Monitoring**: Track script execution time and errors

## Recommendation

### Phase 1 (Current - Tier 0): Option 1 - No Script Pushdown ✅

**Status**: Already implemented
**Rationale**: Fastest path to working MVP
**Trade-off**: Accept higher data transfer for simplicity

### Phase 2 (Tier 1): Option 3 - Selective Pushdown 🎯

**Implementation Plan**:

1. **Whitelist Common Functions** (Week 7-8)
   - Math: abs, round, floor, ceil, sqrt
   - String: upper, lower, trim, length
   - Date: year, month, day, hour
   - Total: ~20 functions

2. **Add Painless Generator** (2-3 days)
   ```
   pkg/ppl/dsl/script_generator.go  (300 lines)
   pkg/ppl/dsl/script_generator_test.go (200 lines)
   ```

3. **Update Physical Planner** (1 day)
   - Add `canPushDownAsScript()` logic
   - Push down whitelisted functions

4. **Update Query Builder** (1-2 days)
   - Generate script filters
   - Handle script_fields for projections

5. **Testing & Validation** (1-2 days)
   - Unit tests for script generation
   - Integration tests with OpenSearch
   - Performance benchmarks

**Estimated Effort**: 5-8 days total

### Phase 3 (Tier 2+): Expand Whitelist

Add more functions as needed based on usage patterns and performance profiling.

## Implementation Complexity

### Code Changes Required for Option 3

1. **New File**: `pkg/ppl/dsl/script_generator.go`
   - `PainlessGenerator` struct
   - `GenerateScript(expr)` method
   - Function mapping table
   - Parameterization logic

2. **Modified**: `pkg/ppl/physical/planner.go`
   - Update `canPushDownFilter()` to check script eligibility
   - Add `canPushDownAsScript()` method
   - Mark expressions as "script-pushable"

3. **Modified**: `pkg/ppl/dsl/query_builder.go`
   - Add `buildScriptFilter()` method
   - Check if filter should use script
   - Integrate with PainlessGenerator

4. **Modified**: `pkg/ppl/dsl/agg_builder.go`
   - Support script-based aggregations
   - Handle computed fields in metrics

5. **Tests**: Comprehensive coverage
   - Script generation correctness
   - Security (no injection)
   - Performance benchmarks

**Estimated**: ~1,200 new lines of code, ~300 modified lines

## Benchmark Targets

### Filter with Function Call

**Query**: `source=logs | where abs(latency) > 100`

| Metric | No Pushdown | With Script | Target |
|--------|-------------|-------------|--------|
| Query Time | 500ms | 50ms | 10x improvement |
| Data Transfer | 100MB | 10MB | 10x reduction |
| Coordinator CPU | 80% | 10% | 8x reduction |

### Aggregation with Computed Field

**Query**: `source=logs | stats avg(latency * 1000) by status`

| Metric | No Pushdown | With Script | Target |
|--------|-------------|-------------|--------|
| Query Time | 2s | 200ms | 10x improvement |
| Data Transfer | 500MB | 1KB (agg only) | 500,000x reduction |
| Coordinator CPU | 95% | 5% | 19x reduction |

## Decision Matrix

| Criteria | Option 1 (Current) | Option 2 (Full) | Option 3 (Selective) | Option 4 (Projections) |
|----------|-------------------|-----------------|---------------------|----------------------|
| **Complexity** | ✅ Low | ❌ High | ⚠️ Medium | ⚠️ Low-Medium |
| **Performance** | ❌ Poor | ✅ Excellent | ✅ Good | ⚠️ Limited |
| **Security** | ✅ Safe | ❌ Risky | ✅ Controlled | ✅ Safe |
| **Maintainability** | ✅ Easy | ❌ Hard | ⚠️ Moderate | ✅ Easy |
| **MVP Ready** | ✅ Yes | ❌ No | ⚠️ Post-MVP | ⚠️ Post-MVP |
| **Tier 1 Ready** | ⚠️ Limited | ✅ Yes | ✅ Yes | ❌ No |

## Conclusion

### Current Status
✅ **Option 1 (No Script Pushdown)** is correctly implemented for Tier 0 MVP.

### Recommendation for Tier 1
🎯 **Option 3 (Selective Pushdown)** should be implemented in Tier 1 to achieve production-grade performance while maintaining reasonable complexity and security.

### Action Items
1. ✅ Complete Tier 0 with current design (no changes needed)
2. ⏳ In Tier 1, implement selective script pushdown for ~20 common functions
3. ⏳ Benchmark performance improvements
4. ⏳ Expand function whitelist based on real usage

### Success Criteria
- 10x performance improvement for filtered queries with functions
- 100x+ performance improvement for aggregations with computed fields
- No security vulnerabilities (parameterized scripts only)
- <5ms script generation overhead

---

**Author**: Claude (PPL Implementation)
**Reviewers**: [TBD]
**Status**: Design Review - Ready for Tier 0 Executor, Tier 1 Enhancement Planned
