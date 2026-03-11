package planner

import (
	"testing"
	"time"

	"github.com/conjugate/conjugate/pkg/coordination/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDateMathIntegration_RangeQueryConversion tests the full pipeline:
// JSON query -> Parser -> Converter (with date math) -> Logical Plan
func TestDateMathIntegration_RangeQueryConversion(t *testing.T) {
	// This simulates the big5 benchmark range query
	queryJSON := `{
		"query": {
			"bool": {
				"must": [
					{
						"range": {
							"metrics.size": {
								"gte": 1000
							}
						}
					},
					{
						"range": {
							"@timestamp": {
								"gte": "now-7d"
							}
						}
					}
				]
			}
		},
		"size": 10
	}`

	// Parse the query
	p := parser.NewQueryParser()
	req, err := p.ParseSearchRequest([]byte(queryJSON))
	require.NoError(t, err)

	// Convert to logical plan (date math should be parsed here)
	converter := NewConverter()
	plan, err := converter.ConvertSearchRequest(req, "logs", []int32{0})
	require.NoError(t, err)

	// Verify plan structure
	require.NotNil(t, plan)

	// Extract the scan node (should have filter pushed down)
	var scan *LogicalScan
	current := plan
	for current != nil {
		if s, ok := current.(*LogicalScan); ok {
			scan = s
			break
		}
		// Navigate down the plan tree
		switch node := current.(type) {
		case *LogicalLimit:
			current = node.Child
		case *LogicalSort:
			current = node.Child
		case *LogicalFilter:
			current = node.Child
		case *LogicalProject:
			current = node.Child
		default:
			break
		}
	}

	require.NotNil(t, scan, "Should have a scan node")
	require.NotNil(t, scan.Filter, "Scan should have filter pushed down")

	// The filter should be a boolean expression with two range queries
	require.Equal(t, ExprTypeBool, scan.Filter.Type)
	// Bool queries store the type ("must", "should", etc.) in Value field
	// But when must clauses are flattened, Value may be nil
	require.Len(t, scan.Filter.Children, 2)

	// First child: numeric range (no date math)
	numericRange := scan.Filter.Children[0]
	require.Equal(t, ExprTypeRange, numericRange.Type)
	require.Equal(t, "metrics.size", numericRange.Field)

	numericParams, ok := numericRange.Value.(map[string]interface{})
	require.True(t, ok)
	// JSON numbers are parsed as float64 in Go
	assert.Equal(t, float64(1000), numericParams["gte"]) // Should be unchanged

	// Second child: date range (should have date math parsed)
	dateRange := scan.Filter.Children[1]
	require.Equal(t, ExprTypeRange, dateRange.Type)
	require.Equal(t, "@timestamp", dateRange.Field)

	dateParams, ok := dateRange.Value.(map[string]interface{})
	require.True(t, ok)

	// The "now-7d" should have been converted to an absolute timestamp
	gteValue, exists := dateParams["gte"]
	require.True(t, exists, "Should have 'gte' parameter")

	// Verify it's now a timestamp string, not "now-7d"
	gteStr, ok := gteValue.(string)
	require.True(t, ok, "Date math should be converted to string timestamp")
	require.NotEqual(t, "now-7d", gteStr, "Should NOT be the original date math expression")

	// Verify it's a valid ISO8601 timestamp
	parsedTime, err := time.Parse("2006-01-02T15:04:05.000Z", gteStr)
	require.NoError(t, err, "Should be valid ISO8601 timestamp")

	// Verify it's approximately 7 days ago (within 1 minute tolerance)
	expectedTime := time.Now().UTC().AddDate(0, 0, -7)
	timeDiff := parsedTime.Sub(expectedTime).Abs()
	assert.Less(t, timeDiff.Minutes(), 1.0, "Timestamp should be approximately 7 days ago")

	t.Logf("Date math 'now-7d' converted to: %s", gteStr)
	t.Logf("Expected approximately: %s", expectedTime.Format("2006-01-02T15:04:05.000Z"))
}

// TestDateMathIntegration_SimpleRangeQuery tests a simple range query with date math
func TestDateMathIntegration_SimpleRangeQuery(t *testing.T) {
	queryJSON := `{
		"query": {
			"range": {
				"@timestamp": {
					"gte": "now-1d",
					"lte": "now"
				}
			}
		}
	}`

	// Parse
	p := parser.NewQueryParser()
	req, err := p.ParseSearchRequest([]byte(queryJSON))
	require.NoError(t, err)

	// Convert
	converter := NewConverter()
	expr, err := converter.ConvertQuery(req.ParsedQuery)
	require.NoError(t, err)

	// Verify
	require.Equal(t, ExprTypeRange, expr.Type)
	require.Equal(t, "@timestamp", expr.Field)

	params, ok := expr.Value.(map[string]interface{})
	require.True(t, ok)

	// Both gte and lte should be absolute timestamps
	gteStr, ok := params["gte"].(string)
	require.True(t, ok)
	require.NotEqual(t, "now-1d", gteStr)

	lteStr, ok := params["lte"].(string)
	require.True(t, ok)
	require.NotEqual(t, "now", lteStr)

	// Verify both are valid timestamps
	gteTime, err := time.Parse("2006-01-02T15:04:05.000Z", gteStr)
	require.NoError(t, err)

	lteTime, err := time.Parse("2006-01-02T15:04:05.000Z", lteStr)
	require.NoError(t, err)

	// lte should be approximately now
	nowTime := time.Now().UTC()
	lteDiff := lteTime.Sub(nowTime).Abs()
	assert.Less(t, lteDiff.Seconds(), 5.0, "lte should be approximately now")

	// gte should be approximately 1 day before lte
	gteExpected := lteTime.AddDate(0, 0, -1)
	gteDiff := gteTime.Sub(gteExpected).Abs()
	assert.Less(t, gteDiff.Seconds(), 5.0, "gte should be approximately 1 day ago")

	t.Logf("Date math conversion:")
	t.Logf("  'now-1d' -> %s", gteStr)
	t.Logf("  'now'    -> %s", lteStr)
}

// TestDateMathIntegration_NoDateMath verifies non-date-math values pass through
func TestDateMathIntegration_NoDateMath(t *testing.T) {
	queryJSON := `{
		"query": {
			"range": {
				"price": {
					"gte": 100,
					"lte": 500
				}
			}
		}
	}`

	// Parse
	p := parser.NewQueryParser()
	req, err := p.ParseSearchRequest([]byte(queryJSON))
	require.NoError(t, err)

	// Convert
	converter := NewConverter()
	expr, err := converter.ConvertQuery(req.ParsedQuery)
	require.NoError(t, err)

	// Verify
	require.Equal(t, ExprTypeRange, expr.Type)
	require.Equal(t, "price", expr.Field)

	params, ok := expr.Value.(map[string]interface{})
	require.True(t, ok)

	// Numeric values should pass through unchanged (JSON numbers are float64)
	assert.Equal(t, float64(100), params["gte"])
	assert.Equal(t, float64(500), params["lte"])
}

// BenchmarkDateMathIntegration measures the overhead of date math parsing
func BenchmarkDateMathIntegration(b *testing.B) {
	queryJSON := `{
		"query": {
			"bool": {
				"must": [
					{"range": {"metrics.size": {"gte": 1000}}},
					{"range": {"@timestamp": {"gte": "now-7d"}}}
				]
			}
		},
		"size": 10
	}`

	p := parser.NewQueryParser()
	req, _ := p.ParseSearchRequest([]byte(queryJSON))
	converter := NewConverter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = converter.ConvertSearchRequest(req, "logs", []int32{0})
	}
}
