# Script Pushdown Design - AGGRESSIVE PUSHDOWN

**Date**: January 28, 2026
**Status**: CORRECT DESIGN - Implement Now
**Decision**: MAXIMUM PUSHDOWN to OpenSearch using Painless scripts

## Core Principle

**PUSH EVERYTHING POSSIBLE TO OPENSEARCH**

Minimize coordinator work, minimize data transfer, maximize OpenSearch compute utilization.

## What Gets Pushed Down

### 1. Filters with Functions → Script Query ✅

```sql
source=logs | where abs(latency) > 100
```
→
```json
{
  "query": {
    "script": {
      "script": {
        "source": "Math.abs(doc['latency'].value) > params.threshold",
        "lang": "painless",
        "params": {"threshold": 100}
      }
    }
  }
}
```

**Benefit**: Only matching documents transferred (99%+ reduction in data transfer)

### 2. Computed Fields → Script Fields ✅

```sql
source=logs | fields status, latency, latency * 2 as double_latency
```
→
```json
{
  "_source": ["status", "latency"],
  "script_fields": {
    "double_latency": {
      "script": {
        "source": "doc['latency'].value * 2",
        "lang": "painless"
      }
    }
  }
}
```

**Benefit**: Computation happens on OpenSearch, coordinator just passes through

### 3. Aggregations with Computed Fields → Script Aggregations ✅

```sql
source=logs | stats avg(latency * 1000) as avg_ms by status
```
→
```json
{
  "size": 0,
  "aggs": {
    "group_by_status": {
      "terms": {"field": "status"},
      "aggs": {
        "avg_ms": {
          "avg": {
            "script": {
              "source": "doc['latency'].value * 1000",
              "lang": "painless"
            }
          }
        }
      }
    }
  }
}
```

**Benefit**: Only aggregation results transferred (MB → KB), no raw documents

### 4. Complex Expressions → Painless Translation ✅

```sql
source=logs | where (latency > 100 AND status >= 400) OR response_time < 50
```
→
```json
{
  "query": {
    "bool": {
      "should": [
        {
          "bool": {
            "must": [
              {"range": {"latency": {"gt": 100}}},
              {"range": {"status": {"gte": 400}}}
            ]
          }
        },
        {"range": {"response_time": {"lt": 50}}}
      ]
    }
  }
}
```

**Note**: No script needed here - use native DSL. Only use scripts when necessary.

## Function Mapping to Painless

### Math Functions

| PPL Function | Painless Equivalent | Example |
|--------------|-------------------|---------|
| `abs(x)` | `Math.abs(x)` | `Math.abs(doc['latency'].value)` |
| `ceil(x)` | `Math.ceil(x)` | `Math.ceil(doc['latency'].value)` |
| `floor(x)` | `Math.floor(x)` | `Math.floor(doc['latency'].value)` |
| `round(x)` | `Math.round(x)` | `Math.round(doc['latency'].value)` |
| `sqrt(x)` | `Math.sqrt(x)` | `Math.sqrt(doc['latency'].value)` |
| `pow(x, y)` | `Math.pow(x, y)` | `Math.pow(doc['latency'].value, 2)` |
| `log(x)` | `Math.log(x)` | `Math.log(doc['latency'].value)` |
| `log10(x)` | `Math.log10(x)` | `Math.log10(doc['latency'].value)` |
| `exp(x)` | `Math.exp(x)` | `Math.exp(doc['latency'].value)` |
| `sin(x)` | `Math.sin(x)` | `Math.sin(doc['angle'].value)` |
| `cos(x)` | `Math.cos(x)` | `Math.cos(doc['angle'].value)` |
| `tan(x)` | `Math.tan(x)` | `Math.tan(doc['angle'].value)` |

### String Functions

| PPL Function | Painless Equivalent | Example |
|--------------|-------------------|---------|
| `upper(s)` | `s.toUpperCase()` | `doc['host'].value.toUpperCase()` |
| `lower(s)` | `s.toLowerCase()` | `doc['host'].value.toLowerCase()` |
| `trim(s)` | `s.trim()` | `doc['message'].value.trim()` |
| `length(s)` | `s.length()` | `doc['message'].value.length()` |
| `substring(s, i, n)` | `s.substring(i, i+n)` | `doc['message'].value.substring(0, 10)` |
| `concat(s1, s2)` | `s1 + s2` | `doc['first'].value + doc['last'].value` |
| `replace(s, old, new)` | `s.replace(old, new)` | `doc['message'].value.replace('old', 'new')` |

### Date Functions

| PPL Function | Painless Equivalent | Example |
|--------------|-------------------|---------|
| `year(d)` | `d.getYear()` | `doc['timestamp'].value.getYear()` |
| `month(d)` | `d.getMonthValue()` | `doc['timestamp'].value.getMonthValue()` |
| `day(d)` | `d.getDayOfMonth()` | `doc['timestamp'].value.getDayOfMonth()` |
| `hour(d)` | `d.getHour()` | `doc['timestamp'].value.getHour()` |
| `minute(d)` | `d.getMinute()` | `doc['timestamp'].value.getMinute()` |
| `second(d)` | `d.getSecond()` | `doc['timestamp'].value.getSecond()` |
| `dayofweek(d)` | `d.getDayOfWeek().getValue()` | `doc['timestamp'].value.getDayOfWeek().getValue()` |

### Type Conversion Functions

| PPL Function | Painless Equivalent | Example |
|--------------|-------------------|---------|
| `int(x)` | `(int)x` | `(int)doc['value'].value` |
| `long(x)` | `(long)x` | `(long)doc['value'].value` |
| `float(x)` | `(float)x` | `(float)doc['value'].value` |
| `double(x)` | `(double)x` | `(double)doc['value'].value` |
| `string(x)` | `String.valueOf(x)` | `String.valueOf(doc['value'].value)` |

### Conditional Functions

| PPL Function | Painless Equivalent | Example |
|--------------|-------------------|---------|
| `if(cond, a, b)` | `cond ? a : b` | `doc['status'].value >= 400 ? 'error' : 'ok'` |
| `coalesce(a, b)` | `a != null ? a : b` | Field handling with null checks |
| `isnull(x)` | `x == null` | `doc['field'].size() == 0` |
| `isnotnull(x)` | `x != null` | `doc['field'].size() > 0` |

## Implementation Architecture

### New Package: pkg/ppl/dsl/script/

```
pkg/ppl/dsl/script/
├── generator.go          # Main script generator
├── functions.go          # Function mapping table
├── translator.go         # Expression → Painless
├── security.go           # Parameterization & validation
└── generator_test.go     # Comprehensive tests
```

### Core Components

#### 1. PainlessGenerator

```go
type PainlessGenerator struct {
    functionMap map[string]PainlessFunction
    params      map[string]interface{}
    paramCount  int
}

type PainlessFunction struct {
    Name           string
    Template       string // e.g., "Math.abs(%s)"
    RequiredArgs   int
    Validator      func(args []ast.Expression) error
}

func (pg *PainlessGenerator) Generate(expr ast.Expression) (*PainlessScript, error)
```

#### 2. Expression Translator

```go
type PainlessScript struct {
    Source string                 // Painless script source
    Lang   string                 // "painless"
    Params map[string]interface{} // Parameters for security
}

func (pg *PainlessGenerator) translateExpression(expr ast.Expression) (string, error) {
    switch e := expr.(type) {
    case *ast.FunctionCall:
        return pg.translateFunction(e)
    case *ast.FieldReference:
        return pg.translateField(e)
    case *ast.Literal:
        return pg.translateLiteral(e)
    case *ast.BinaryExpression:
        return pg.translateBinary(e)
    // ...
    }
}
```

#### 3. Security Layer

```go
func (pg *PainlessGenerator) parameterize(value interface{}) string {
    paramName := fmt.Sprintf("param%d", pg.paramCount)
    pg.paramCount++
    pg.params[paramName] = value
    return "params." + paramName
}

// ALWAYS parameterize user input - NEVER concatenate
```

### Updated Physical Planner

```go
func (pp *PhysicalPlanner) canPushDownFilter(condition ast.Expression) bool {
    switch expr := condition.(type) {
    case *ast.BinaryExpression:
        // Simple comparisons - use native DSL
        if isSimpleComparison(expr) {
            return true
        }
        // Complex expressions - check if scriptable
        return pp.isScriptable(expr)

    case *ast.FunctionCall:
        // All mapped functions are scriptable
        return pp.isScriptableFunction(expr)

    default:
        return false
    }
}

func (pp *PhysicalPlanner) isScriptableFunction(fn *ast.FunctionCall) bool {
    // Check if function maps to Painless
    scriptableFunctions := map[string]bool{
        // Math
        "abs": true, "ceil": true, "floor": true, "round": true,
        "sqrt": true, "pow": true, "log": true, "log10": true,
        "exp": true, "sin": true, "cos": true, "tan": true,

        // String
        "upper": true, "lower": true, "trim": true, "length": true,
        "substring": true, "concat": true, "replace": true,

        // Date
        "year": true, "month": true, "day": true, "hour": true,
        "minute": true, "second": true, "dayofweek": true,

        // Type conversion
        "int": true, "long": true, "float": true, "double": true,
        "string": true,

        // Conditional
        "if": true, "coalesce": true, "isnull": true, "isnotnull": true,
    }

    return scriptableFunctions[strings.ToLower(fn.Name)]
}
```

### Updated Query Builder

```go
func (qb *QueryBuilder) BuildFilter(expr ast.Expression) (map[string]interface{}, error) {
    // Try native DSL first (better performance)
    if nativeQuery, ok := qb.tryNativeDSL(expr); ok {
        return nativeQuery, nil
    }

    // Fall back to script query
    return qb.buildScriptQuery(expr)
}

func (qb *QueryBuilder) buildScriptQuery(expr ast.Expression) (map[string]interface{}, error) {
    generator := script.NewPainlessGenerator()
    painlessScript, err := generator.Generate(expr)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "script": map[string]interface{}{
            "script": map[string]interface{}{
                "source": painlessScript.Source,
                "lang":   "painless",
                "params": painlessScript.Params,
            },
        },
    }, nil
}

func (qb *QueryBuilder) tryNativeDSL(expr ast.Expression) (map[string]interface{}, bool) {
    // Try to use native DSL (term, range, bool) which is faster than scripts
    switch e := expr.(type) {
    case *ast.BinaryExpression:
        if e.Operator == "=" && isFieldLiteralComparison(e) {
            // Use term query (faster than script)
            return qb.buildTermQuery(e), true
        }
        // ... other native DSL patterns
    }
    return nil, false
}
```

### Updated Aggregation Builder

```go
func (ab *AggregationBuilder) buildMetricAggregation(agg *ast.Aggregation) (map[string]interface{}, error) {
    // Check if argument is a simple field reference
    if len(agg.Func.Arguments) == 1 {
        if fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
            // Simple field - use native aggregation
            return ab.buildNativeMetric(agg.Func.Name, fieldRef.Name)
        }
    }

    // Complex expression - use script
    return ab.buildScriptMetric(agg)
}

func (ab *AggregationBuilder) buildScriptMetric(agg *ast.Aggregation) (map[string]interface{}, error) {
    generator := script.NewPainlessGenerator()

    // Generate script for the aggregation expression
    painlessScript, err := generator.Generate(agg.Func.Arguments[0])
    if err != nil {
        return nil, err
    }

    funcName := strings.ToLower(agg.Func.Name)

    return map[string]interface{}{
        funcName: map[string]interface{}{
            "script": map[string]interface{}{
                "source": painlessScript.Source,
                "lang":   "painless",
                "params": painlessScript.Params,
            },
        },
    }, nil
}
```

## Security Best Practices

### 1. Always Parameterize

```go
// WRONG - NEVER DO THIS
script := fmt.Sprintf("doc['latency'].value > %v", userInput)

// CORRECT - Always parameterize
script := "doc['latency'].value > params.threshold"
params := map[string]interface{}{"threshold": userInput}
```

### 2. Validate Before Generation

```go
func (pg *PainlessGenerator) Generate(expr ast.Expression) (*PainlessScript, error) {
    // Validate expression is safe
    if err := pg.validate(expr); err != nil {
        return nil, fmt.Errorf("unsafe expression: %w", err)
    }

    // Generate script
    source, err := pg.translateExpression(expr)
    if err != nil {
        return nil, err
    }

    return &PainlessScript{
        Source: source,
        Lang:   "painless",
        Params: pg.params,
    }, nil
}
```

### 3. Sandbox Reliance

Painless is already sandboxed by OpenSearch:
- No file system access
- No network access
- No reflection
- No infinite loops (circuit breaker)
- CPU/memory limits enforced

## Performance Optimization

### 1. Prefer Native DSL Over Scripts

```go
// Order of preference:
// 1. Native DSL (term, range, bool) - FASTEST
// 2. Painless script - FAST
// 3. Coordinator execution - SLOW

// Example: Simple comparison
status = 500  →  {"term": {"status": 500}}  // Native DSL ✅

// Example: Function call
abs(latency) > 100  →  Script query  // Painless ✅

// Example: Complex function
custom_udf(x)  →  Coordinator  // No choice ❌
```

### 2. Script Compilation Caching

OpenSearch automatically caches compiled scripts:
- First execution: ~100ms (compilation)
- Subsequent executions: <1ms (cached)

### 3. Field Access Optimization

```go
// Prefer doc values (fast)
doc['field'].value

// Avoid source access (slow)
params._source.field
```

## Implementation Plan

### Phase 1: Core Script Generation (2-3 days)

1. **Create script package** (Day 1)
   - `pkg/ppl/dsl/script/generator.go` (400 lines)
   - `pkg/ppl/dsl/script/functions.go` (300 lines)
   - Function mapping table for 40+ functions

2. **Expression translation** (Day 2)
   - Translate all expression types to Painless
   - Handle field references, literals, operators
   - Implement parameterization

3. **Testing** (Day 2-3)
   - Unit tests for each function mapping
   - Security tests (no injection)
   - Edge cases (null handling, type coercion)

### Phase 2: Integration (2 days)

1. **Update Physical Planner** (Day 1)
   - Modify `canPushDownFilter()` to check scriptability
   - Update `canPushDownProject()` for computed fields
   - Add script metadata to physical plan

2. **Update Query Builder** (Day 1-2)
   - Integrate PainlessGenerator
   - Implement `buildScriptQuery()`
   - Implement `buildScriptFields()`

3. **Update Aggregation Builder** (Day 2)
   - Implement `buildScriptMetric()`
   - Handle computed aggregation expressions

### Phase 3: Testing & Validation (1-2 days)

1. **Integration tests**
   - End-to-end query tests with scripts
   - Performance benchmarks
   - Security validation

2. **OpenSearch compatibility**
   - Test against OpenSearch 2.x
   - Verify Painless syntax
   - Validate field access patterns

**Total Estimate**: 5-7 days

## Success Metrics

### Performance Targets

| Query Type | Before (No Script) | After (With Script) | Target |
|------------|-------------------|-------------------|--------|
| Filter with function | 500ms, 100MB | 50ms, 1MB | 10x faster, 100x less data |
| Computed aggregation | 2s, 500MB | 100ms, 1KB | 20x faster, 500,000x less data |
| Computed projection | 300ms, 50MB | 100ms, 50MB | 3x faster (compute offload) |

### Coverage Targets

- ✅ 40+ functions mapped to Painless
- ✅ 100% of math functions pushable
- ✅ 100% of string functions pushable
- ✅ 100% of date functions pushable
- ✅ All aggregations support computed fields

## Risk Mitigation

### Risk 1: Painless Syntax Errors

**Mitigation**: Comprehensive unit tests, syntax validation before sending to OpenSearch

### Risk 2: Performance Regression

**Mitigation**: Benchmarks, prefer native DSL when possible, monitoring

### Risk 3: Security Vulnerabilities

**Mitigation**: Always parameterize, never concatenate user input, validation layer

### Risk 4: OpenSearch Compatibility

**Mitigation**: Test against OpenSearch 2.x, document compatibility matrix

## Conclusion

**This is the CORRECT design**. Maximum pushdown to OpenSearch using Painless scripts will:

1. ✅ Minimize data transfer (100x+ reduction for selective queries)
2. ✅ Maximize OpenSearch compute utilization
3. ✅ Minimize coordinator bottleneck
4. ✅ Enable production-grade performance
5. ✅ Properly leverage OpenSearch capabilities

**No excuses. Implement this now.**

---

**Status**: APPROVED - Implement before Executor
**Estimated Effort**: 5-7 days
**Priority**: CRITICAL - Foundational architecture decision
