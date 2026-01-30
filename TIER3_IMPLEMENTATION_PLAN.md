# Tier 3 Implementation Plan: Enterprise Features

**Date**: January 29, 2026
**Target**: 99% Query Coverage (188 functions, 36 commands)
**Timeline**: 8 weeks (estimated 6-7 weeks with current velocity)
**Current Status**: 3/12 commands complete (25%)

---

## Executive Summary

Implementing Tier 3 will bring CONJUGATE from **95% to 99% query coverage**, unlocking enterprise use cases including SIEM, grok pattern parsing, subqueries, and window functions.

**Completed Commands** (3/12):
- ✅ `flatten` - Flatten nested objects
- ✅ `table` - Format output table
- ✅ `reverse` - Reverse row order

**Remaining Commands** (9/12):
- `grok` - Pattern library parsing (HIGHEST PRIORITY)
- `spath` - JSON path navigation
- `subquery` - IN, EXISTS operations
- `eventstats` - Window functions
- `streamstats` - Running totals
- `addtotals` - Add summary rows
- `addcoltotals` - Add column totals
- `appendcol` - Add columns from query
- `appendpipe` - Process results further

---

## Implementation Strategy

### Phase 1: Quick Wins (Week 1) - 2.5 days
**Goal**: Knock out simple commands first

#### Task 1.1: `addtotals` (0.5 weeks = 2.5 days)
**Priority**: HIGH (easy, high value)
**Complexity**: LOW
**Purpose**: Add summary rows to results

```ppl
source=sales | stats sum(revenue) by category | addtotals
```

**Implementation**:
- Aggregate all numeric columns
- Add "Total" row at the end
- Simple operator, no complex logic

**Files to Create**:
- `pkg/ppl/executor/addtotals_operator.go` (~150 lines)
- `pkg/ppl/executor/addtotals_operator_test.go` (~200 lines)
- `pkg/ppl/physical/addtotals.go` (~80 lines)

#### Task 1.2: `addcoltotals` (0.5 weeks = 2.5 days)
**Priority**: MEDIUM (easy, moderate value)
**Complexity**: LOW
**Purpose**: Add column totals

```ppl
source=sales | chart sum(revenue) over category by month | addcoltotals
```

**Implementation**:
- Similar to addtotals but for columns
- Add "Total" column
- Matrix transposition logic

**Files to Create**:
- `pkg/ppl/executor/addcoltotals_operator.go` (~150 lines)
- `pkg/ppl/executor/addcoltotals_operator_test.go` (~200 lines)
- `pkg/ppl/physical/addcoltotals.go` (~80 lines)

**Week 1 Total**: 2 commands, ~860 lines

---

### Phase 2: JSON & Pattern Processing (Weeks 2-3) - 2 weeks

#### Task 2.1: `spath` (1 week)
**Priority**: HIGH (JSON navigation is common)
**Complexity**: MEDIUM
**Purpose**: Navigate and extract from JSON structures

```ppl
source=api_logs | spath path="response.data.user.id" output=user_id
source=json_docs | spath input=data path="$.items[*].name" output=item_names
```

**Implementation**:
- JSONPath library integration (gjson or similar)
- Auto-extraction mode (no path specified)
- Nested field creation
- Array handling

**Files to Create**:
- `pkg/ppl/executor/spath_operator.go` (~400 lines)
- `pkg/ppl/executor/spath_operator_test.go` (~500 lines)
- `pkg/ppl/physical/spath.go` (~120 lines)
- JSON path utility helpers

#### Task 2.2: `grok` (1 week)
**Priority**: HIGHEST (critical for enterprise SIEM)
**Complexity**: HIGH
**Purpose**: Pattern-based log parsing with built-in patterns

```ppl
source=access_logs | grok "%{COMMONAPACHELOG}"
source=app_logs | grok "user=%{USER:username} ip=%{IP:client_ip} latency=%{NUMBER:latency:float}ms"
```

**Implementation**:
- Port grok pattern library (50+ patterns)
- Pattern definitions: HOSTNAME, IP, NUMBER, TIMESTAMP, etc.
- Named capture groups
- Type coercion (string, int, float)
- Custom pattern support
- Pattern composition

**Key Patterns to Support**:
- COMMONAPACHELOG, COMBINEDAPACHELOG
- IP, HOSTNAME, USER, EMAIL
- NUMBER, INT, BASE10NUM, BASE16NUM
- TIMESTAMP_ISO8601, SYSLOGTIMESTAMP
- PATH, URI, URL
- LOGLEVEL, UUID

**Files to Create**:
- `pkg/ppl/executor/grok_operator.go` (~600 lines)
- `pkg/ppl/executor/grok_operator_test.go` (~800 lines)
- `pkg/ppl/grok/patterns.go` (~1000 lines - pattern library)
- `pkg/ppl/grok/parser.go` (~400 lines)
- `pkg/ppl/physical/grok.go` (~150 lines)

**Week 2-3 Total**: 2 commands, ~3,970 lines

---

### Phase 3: Result Processing (Week 4) - 1 week

#### Task 3.1: `appendcol` (0.5 weeks = 2.5 days)
**Priority**: MEDIUM (moderate complexity)
**Complexity**: MEDIUM
**Purpose**: Add columns from subsearch

```ppl
source=orders | appendcol [search source=customers | fields customer_id, name, email]
```

**Implementation**:
- Similar to join but simpler (no join key)
- Column-wise concatenation
- Schema merging with conflict handling
- Subsearch execution
- Result alignment

**Files to Create**:
- `pkg/ppl/executor/appendcol_operator.go` (~350 lines)
- `pkg/ppl/executor/appendcol_operator_test.go` (~400 lines)
- `pkg/ppl/physical/appendcol.go` (~120 lines)

#### Task 3.2: `appendpipe` (0.5 weeks = 2.5 days)
**Priority**: MEDIUM (powerful for complex pipelines)
**Complexity**: MEDIUM
**Purpose**: Process results and append back

```ppl
source=errors
| stats count() by error_code
| appendpipe [stats sum(count) as total | eval percentage = count / total * 100]
```

**Implementation**:
- Execute subsearch on current results
- Append subsearch results to main results
- Context passing
- Schema union

**Files to Create**:
- `pkg/ppl/executor/appendpipe_operator.go` (~350 lines)
- `pkg/ppl/executor/appendpipe_operator_test.go` (~400 lines)
- `pkg/ppl/physical/appendpipe.go` (~120 lines)

**Week 4 Total**: 2 commands, ~1,740 lines

---

### Phase 4: Window Functions (Weeks 5-6) - 2 weeks

#### Task 4.1: `eventstats` (1 week)
**Priority**: HIGH (commonly requested)
**Complexity**: HIGH
**Purpose**: Add aggregated statistics to each event

```ppl
source=response_times
| eventstats avg(latency) as avg_latency, stddev(latency) as stddev_latency by endpoint
| eval zscore = (latency - avg_latency) / stddev_latency
| where abs(zscore) > 2
```

**Implementation**:
- Window aggregation framework
- Partition by fields (optional)
- Calculate aggregations per partition
- Add results as new fields to ALL events
- No row reduction (unlike stats)

**Window Functions to Support**:
- All aggregation functions (count, sum, avg, min, max)
- Statistical: stddev, variance, percentile
- Advanced: first, last, values, distinct_count

**Files to Create**:
- `pkg/ppl/executor/eventstats_operator.go` (~600 lines)
- `pkg/ppl/executor/eventstats_operator_test.go` (~700 lines)
- `pkg/ppl/executor/window_framework.go` (~500 lines - shared)
- `pkg/ppl/physical/eventstats.go` (~150 lines)

#### Task 4.2: `streamstats` (1 week)
**Priority**: HIGH (financial analysis critical)
**Complexity**: HIGH
**Purpose**: Running/cumulative statistics

```ppl
source=sales
| sort timestamp
| streamstats sum(revenue) as cumulative_revenue by region
| streamstats current=f avg(revenue) as moving_avg window=5
```

**Implementation**:
- Ordered streaming computation
- Running totals (cumulative)
- Moving windows (window=N)
- Current vs non-current semantics
- Partition support (by fields)
- Reset on partition boundaries

**Window Types**:
- Cumulative: sum, count, avg
- Moving: last N values
- Global vs partitioned

**Files to Create**:
- `pkg/ppl/executor/streamstats_operator.go` (~700 lines)
- `pkg/ppl/executor/streamstats_operator_test.go` (~800 lines)
- `pkg/ppl/physical/streamstats.go` (~150 lines)

**Week 5-6 Total**: 2 commands, ~3,600 lines

---

### Phase 5: Subqueries (Weeks 7-8) - 2 weeks

#### Task 5.1: `subquery` Support (2 weeks)
**Priority**: HIGHEST (enterprise critical)
**Complexity**: VERY HIGH
**Purpose**: Nested queries with IN, EXISTS, scalar

```ppl
-- IN subquery
source=alerts | where severity IN [search source=critical_errors | fields error_code]

-- EXISTS subquery (correlated)
source=users | where EXISTS [search source=orders | where orders.user_id = users.id]

-- Scalar subquery
source=orders | where amount > [search source=orders | stats avg(amount) as threshold | fields threshold]
```

**Implementation**:

**5.1.1: Subquery Parser Extensions (2 days)**
- Extend ANTLR grammar for `[search ...]` syntax
- AST node for SubqueryExpression
- Parse IN, EXISTS, scalar contexts
- Validation and type checking

**5.1.2: Subquery Executor (3 days)**
- Subsearch execution engine
- Result materialization (10K row limit)
- Context isolation
- Error handling

**5.1.3: IN Subquery (1 day)**
- Execute subsearch → list of values
- Transform to `field IN (v1, v2, v3, ...)`
- Optimization: convert to hash lookup

**5.1.4: EXISTS Subquery (3 days)**
- Correlated vs uncorrelated
- Correlation detection
- Semi-join implementation
- Context variable passing

**5.1.5: Scalar Subquery (1 day)**
- Execute subsearch → single value
- Validation (must return 1 row, 1 column)
- Type coercion
- Use as literal in expressions

**Files to Create**:
- `pkg/ppl/parser/ppl_parser.g4` - Updates for subquery syntax
- `pkg/ppl/ast/subquery.go` (~200 lines)
- `pkg/ppl/analyzer/subquery_analyzer.go` (~300 lines)
- `pkg/ppl/executor/subquery_executor.go` (~800 lines)
- `pkg/ppl/executor/subquery_in_operator.go` (~250 lines)
- `pkg/ppl/executor/subquery_exists_operator.go` (~350 lines)
- `pkg/ppl/executor/subquery_scalar_operator.go` (~200 lines)
- `pkg/ppl/executor/subquery_test.go` (~1000 lines)
- `pkg/ppl/physical/subquery.go` (~200 lines)
- `pkg/ppl/optimizer/subquery_optimizer.go` (~400 lines)

**Week 7-8 Total**: 1 command (4 subtypes), ~3,700 lines

---

## Complete Timeline

| Week | Phase | Commands | Effort | Cumulative |
|------|-------|----------|--------|------------|
| **1** | Quick Wins | addtotals, addcoltotals | 2.5d+2.5d | 2 (17%) |
| **2-3** | JSON & Patterns | spath, grok | 1w+1w | 4 (33%) |
| **4** | Result Processing | appendcol, appendpipe | 2.5d+2.5d | 6 (50%) |
| **5-6** | Window Functions | eventstats, streamstats | 1w+1w | 8 (67%) |
| **7-8** | Subqueries | subquery (IN/EXISTS/scalar) | 2w | 9 (75%) |

**Total**: 8 weeks, 9 commands
**Already Complete**: 3 commands (flatten, table, reverse)
**Grand Total**: 12/12 commands (100%)

---

## Priority Ranking

### Critical Path (Must Have)
1. **grok** - Enterprise SIEM requirement
2. **subquery** - Advanced analytics capability
3. **eventstats** - Common window function
4. **streamstats** - Financial/time-series analysis

### High Value (Should Have)
5. **spath** - JSON navigation is common
6. **addtotals** - Reporting/dashboards

### Nice to Have
7. **addcoltotals** - Matrix operations
8. **appendcol** - Advanced joins
9. **appendpipe** - Complex pipelines

---

## Technical Challenges

### Challenge 1: Grok Pattern Library
**Complexity**: Pattern translation from Ruby to Go
**Solution**:
- Port patterns from Logstash grok library
- Use Go regex (different from Ruby)
- Test thoroughly with real logs

**Estimated Effort**: 3 days of pattern testing

### Challenge 2: Correlated Subqueries
**Complexity**: Variable scope and context passing
**Solution**:
- Extend ExecutionContext with parent scope
- Implement variable resolution chain
- Use semi-join for EXISTS

**Estimated Effort**: 4 days

### Challenge 3: Window Function State Management
**Complexity**: Streaming state for running calculations
**Solution**:
- Circular buffer for moving windows
- Accumulator pattern for running totals
- Partition-aware state management

**Estimated Effort**: 3 days

### Challenge 4: Memory Management
**Complexity**: Subqueries and windows can be memory-intensive
**Solution**:
- 10K row limit on subqueries
- Streaming execution for windows
- Spill to disk for large partitions (future)

---

## Testing Strategy

### Unit Tests (per command)
- Parser: Valid/invalid syntax (10-15 tests)
- Executor: Core functionality (15-20 tests)
- Edge cases: Empty results, nulls, errors (10 tests)
- Performance: Memory usage, execution time (3 tests)

**Total**: ~40 tests per command × 9 commands = **360 unit tests**

### Integration Tests
- Multi-command pipelines (20 tests)
- Real-world log examples (15 tests)
- Large dataset handling (10 tests)
- Error handling (10 tests)

**Total**: **55 integration tests**

### E2E Tests
- SIEM use case with grok (5 tests)
- Financial analysis with streamstats (5 tests)
- Complex subquery scenarios (5 tests)

**Total**: **15 e2e tests**

**Grand Total**: **430 tests** for Tier 3

---

## Resource Requirements

### Engineers
- **Lead Engineer** (1): Complex commands (subquery, grok)
- **Mid-Level Engineer** (1): Window functions
- **Junior Engineer** (0.5): Simple commands, testing

**Total**: 2.5 engineers

### Effort Breakdown
- Implementation: ~11,000 lines of code
- Tests: ~8,000 lines of test code
- Documentation: ~2,000 lines
- **Total**: ~21,000 lines

### Timeline
- **Planned**: 8 weeks
- **Realistic**: 6-7 weeks (based on current velocity)
- **Aggressive**: 5 weeks (with 3 engineers)

---

## Success Criteria

### Functional
- ✅ All 12 Tier 3 commands implemented
- ✅ 188 total functions (98% coverage)
- ✅ All 430 tests passing
- ✅ Query coverage: 99%

### Performance
- ✅ Grok parsing: <10ms per event
- ✅ Window functions: <200ms for 10K rows
- ✅ Subqueries: <500ms for 10K row limit
- ✅ Memory usage: <1GB per query

### Quality
- ✅ Test coverage: >90%
- ✅ No memory leaks
- ✅ Error handling: Clear messages
- ✅ Documentation: Complete examples

---

## Risk Assessment

### High Risk
1. **Grok Pattern Compatibility** (Risk: Medium)
   - Mitigation: Test with real Logstash logs
   - Fallback: Support subset of critical patterns first

2. **Subquery Performance** (Risk: Medium)
   - Mitigation: Implement 10K row limit, caching
   - Fallback: Start with uncorrelated only

### Medium Risk
3. **Window Function Memory** (Risk: Medium)
   - Mitigation: Streaming execution, partition limits
   - Fallback: Add explicit limits in docs

### Low Risk
4. **Simple Commands** (addtotals, spath)
   - Low complexity, well-understood

---

## Dependencies

### Internal
- ✅ Parser infrastructure (complete)
- ✅ Executor framework (complete)
- ✅ Aggregation framework (complete)
- ⚠️ Window function framework (new - Week 5)
- ⚠️ Subquery execution engine (new - Week 7)

### External
- Grok pattern library (port from Logstash)
- JSONPath library (gjson or similar)
- Regex engine (Go stdlib)

---

## Deliverables

### Week 1
- ✅ addtotals command
- ✅ addcoltotals command
- 📄 Implementation guide

### Week 2-3
- ✅ spath command
- ✅ grok command
- ✅ Grok pattern library (50+ patterns)
- 📄 Grok pattern reference

### Week 4
- ✅ appendcol command
- ✅ appendpipe command
- 📄 Advanced pipeline guide

### Week 5-6
- ✅ eventstats command
- ✅ streamstats command
- ✅ Window function framework
- 📄 Window function guide

### Week 7-8
- ✅ subquery support (IN, EXISTS, scalar)
- ✅ Subquery optimizer
- 📄 Subquery best practices

### Week 8 (Final)
- 📄 Tier 3 Complete Summary
- 📄 99% Coverage Report
- 📊 Performance benchmarks
- 🎉 **TIER 3 COMPLETE**

---

## Next Steps

### Immediate (This Week)
1. **Review and Approve Plan**
   - Stakeholder review
   - Timeline confirmation
   - Resource allocation

2. **Prepare Development Environment**
   - Set up grok pattern library
   - Research JSONPath libraries
   - Create test data sets

3. **Start Week 1 Implementation**
   - Begin with addtotals (easiest)
   - Quick win to build momentum

### Week 2 Onward
- Follow phased implementation
- Weekly progress reviews
- Adjust timeline as needed

---

## Comparison: T2 vs T3

| Metric | Tier 2 | Tier 3 | Increase |
|--------|--------|--------|----------|
| **Commands** | 9 | 12 | +33% |
| **Functions** | 30 | 23 | -23% |
| **Effort (weeks)** | 10 | 8 | -20% |
| **Complexity** | Medium-High | High-Very High | +30% |
| **Code (LOC)** | ~5,500 | ~11,000 | +100% |
| **Use Cases** | Advanced Analytics | Enterprise + SIEM | Critical |

**Key Insight**: Tier 3 has fewer functions but more complex commands (grok, subquery, windows). Higher code volume due to complexity.

---

## Post-Tier 3 Options

After completing Tier 3, three paths forward:

### Option A: Production Hardening (Recommended)
- Parallel execution
- Memory optimization
- Large dataset handling (1B+ rows)
- Performance tuning
- **Timeline**: 4-6 weeks

### Option B: Tier 4 (ML Features)
- ML command framework
- Anomaly detection (RCF)
- Pattern discovery
- K-means clustering
- **Timeline**: 8 weeks

### Option C: API & Integration
- REST API for PPL
- Client libraries
- Query builder UI
- Documentation site
- **Timeline**: 4-6 weeks

**Recommendation**: Option A (Production Hardening) - ensures stability before adding more features

---

## Conclusion

Tier 3 implementation will complete the enterprise-grade PPL capabilities, bringing CONJUGATE to **99% query coverage**. With 3 commands already complete, we have **9 commands remaining** over **8 weeks**.

**Critical Path**: grok → subquery → eventstats → streamstats

**Quick Wins**: addtotals, addcoltotals (Week 1)

**Expected Completion**: Early April 2026

**Status**: Ready to begin implementation 🚀

---

**Document Version**: 1.0
**Last Updated**: January 29, 2026
**Author**: CONJUGATE Team
**Status**: Ready for Implementation
