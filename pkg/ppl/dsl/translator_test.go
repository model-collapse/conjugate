// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package dsl

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestSchema() *analyzer.Schema {
	schema := analyzer.NewSchema("logs")
	schema.AddField("status", analyzer.FieldTypeInt)
	schema.AddField("host", analyzer.FieldTypeString)
	schema.AddField("timestamp", analyzer.FieldTypeDate)
	schema.AddField("latency", analyzer.FieldTypeDouble)
	schema.AddField("message", analyzer.FieldTypeText)
	return schema
}

func TestTranslator_SimpleScan(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)
	require.NotNil(t, dsl)

	// Should have match_all query
	assert.NotNil(t, dsl.Query)
	assert.Equal(t, map[string]interface{}{
		"match_all": map[string]interface{}{},
	}, dsl.Query)

	// No _source, sort, size
	assert.Nil(t, dsl.Source)
	assert.Nil(t, dsl.Sort)
	assert.Nil(t, dsl.Size)
}

func TestTranslator_FilterTerm(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// status = 500
	filter := &ast.BinaryExpression{
		Left:     &ast.FieldReference{Name: "status"},
		Operator: "=",
		Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
	}

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		Filter:       filter,
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"term": map[string]interface{}{
			"status": 500,
		},
	}
	assert.Equal(t, expected, dsl.Query)
}

func TestTranslator_FilterRange(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// latency > 100.0
	filter := &ast.BinaryExpression{
		Left:     &ast.FieldReference{Name: "latency"},
		Operator: ">",
		Right:    &ast.Literal{Value: 100.0, LiteralTyp: ast.LiteralTypeFloat},
	}

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		Filter:       filter,
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"range": map[string]interface{}{
			"latency": map[string]interface{}{
				"gt": 100.0,
			},
		},
	}
	assert.Equal(t, expected, dsl.Query)
}

func TestTranslator_FilterAND(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// status = 500 AND host = "server1"
	filter := &ast.BinaryExpression{
		Left: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Operator: "AND",
		Right: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "host"},
			Operator: "=",
			Right:    &ast.Literal{Value: "server1", LiteralTyp: ast.LiteralTypeString},
		},
	}

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		Filter:       filter,
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should be a bool query with must
	boolQuery, ok := dsl.Query["bool"].(map[string]interface{})
	require.True(t, ok)

	must, ok := boolQuery["must"].([]interface{})
	require.True(t, ok)
	assert.Len(t, must, 2)
}

func TestTranslator_FilterOR(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// status = 500 OR status = 404
	filter := &ast.BinaryExpression{
		Left: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Operator: "OR",
		Right: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 404, LiteralTyp: ast.LiteralTypeInt},
		},
	}

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		Filter:       filter,
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should be a bool query with should
	boolQuery, ok := dsl.Query["bool"].(map[string]interface{})
	require.True(t, ok)

	should, ok := boolQuery["should"].([]interface{})
	require.True(t, ok)
	assert.Len(t, should, 2)
}

func TestTranslator_FilterNOT(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// NOT status = 200
	filter := &ast.UnaryExpression{
		Operator: "NOT",
		Operand: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
		},
	}

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		Filter:       filter,
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should be a bool query with must_not
	boolQuery, ok := dsl.Query["bool"].(map[string]interface{})
	require.True(t, ok)

	mustNot, ok := boolQuery["must_not"].([]interface{})
	require.True(t, ok)
	assert.Len(t, mustNot, 1)
}

func TestTranslator_Projection(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		Fields:       []string{"status", "host", "timestamp"},
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	assert.Equal(t, []string{"status", "host", "timestamp"}, dsl.Source)
}

func TestTranslator_Sort(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		SortKeys: []*ast.SortKey{
			{
				Field:      &ast.FieldReference{Name: "timestamp"},
				Descending: true,
			},
			{
				Field:      &ast.FieldReference{Name: "status"},
				Descending: false,
			},
		},
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	require.Len(t, dsl.Sort, 2)

	// First sort: timestamp DESC
	sort1, ok := dsl.Sort[0]["timestamp"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "desc", sort1["order"])

	// Second sort: status ASC
	sort2, ok := dsl.Sort[1]["status"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "asc", sort2["order"])
}

func TestTranslator_Limit(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		Limit:        10,
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	require.NotNil(t, dsl.Size)
	assert.Equal(t, 10, *dsl.Size)
}

func TestTranslator_CombinedQuery(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// Filter, projection, sort, and limit all together
	filter := &ast.BinaryExpression{
		Left:     &ast.FieldReference{Name: "status"},
		Operator: "=",
		Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
	}

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		Filter:       filter,
		Fields:       []string{"status", "host"},
		SortKeys: []*ast.SortKey{
			{Field: &ast.FieldReference{Name: "timestamp"}, Descending: true},
		},
		Limit: 10,
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should have all components
	assert.NotNil(t, dsl.Query)
	assert.Equal(t, []string{"status", "host"}, dsl.Source)
	assert.Len(t, dsl.Sort, 1)
	assert.Equal(t, 10, *dsl.Size)
}

func TestTranslator_AggregationSimple(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// Simple aggregation: count()
	agg := &physical.PhysicalAggregate{
		Aggregations: []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name:      "count",
					Arguments: []ast.Expression{},
				},
				Alias: "total",
			},
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(agg)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)

	// Size should be 0 for aggregation queries
	require.NotNil(t, dsl.Size)
	assert.Equal(t, 0, *dsl.Size)

	// Check aggregation structure
	totalAgg, ok := dsl.Aggregations["total"].(map[string]interface{})
	require.True(t, ok)

	valueCount, ok := totalAgg["value_count"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "_id", valueCount["field"])
}

func TestTranslator_AggregationGroupBy(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// GROUP BY aggregation: count() by host
	agg := &physical.PhysicalAggregate{
		GroupBy: []ast.Expression{
			&ast.FieldReference{Name: "host"},
		},
		Aggregations: []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name:      "count",
					Arguments: []ast.Expression{},
				},
				Alias: "total",
			},
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(agg)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)

	// Should have group_by_host aggregation
	groupAgg, ok := dsl.Aggregations["group_by_host"].(map[string]interface{})
	require.True(t, ok)

	// Should be a terms aggregation
	termsAgg, ok := groupAgg["terms"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "host", termsAgg["field"])

	// Should have sub-aggregations
	subAggs, ok := groupAgg["aggs"].(map[string]interface{})
	require.True(t, ok)

	totalAgg, ok := subAggs["total"].(map[string]interface{})
	require.True(t, ok)
	assert.NotNil(t, totalAgg["value_count"])
}

func TestTranslator_AggregationMultipleMetrics(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// Multiple aggregations: count(), avg(latency), max(latency)
	agg := &physical.PhysicalAggregate{
		GroupBy: []ast.Expression{
			&ast.FieldReference{Name: "host"},
		},
		Aggregations: []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name:      "count",
					Arguments: []ast.Expression{},
				},
				Alias: "total",
			},
			{
				Func: &ast.FunctionCall{
					Name: "avg",
					Arguments: []ast.Expression{
						&ast.FieldReference{Name: "latency"},
					},
				},
				Alias: "avg_latency",
			},
			{
				Func: &ast.FunctionCall{
					Name: "max",
					Arguments: []ast.Expression{
						&ast.FieldReference{Name: "latency"},
					},
				},
				Alias: "max_latency",
			},
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(agg)
	require.NoError(t, err)

	// Get the group aggregation
	groupAgg, ok := dsl.Aggregations["group_by_host"].(map[string]interface{})
	require.True(t, ok)

	// Get sub-aggregations
	subAggs, ok := groupAgg["aggs"].(map[string]interface{})
	require.True(t, ok)

	// Should have all three metrics
	assert.Contains(t, subAggs, "total")
	assert.Contains(t, subAggs, "avg_latency")
	assert.Contains(t, subAggs, "max_latency")
}

func TestTranslator_TranslateToJSON(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	filter := &ast.BinaryExpression{
		Left:     &ast.FieldReference{Name: "status"},
		Operator: "=",
		Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
	}

	plan := &physical.PhysicalScan{
		Source:       "logs",
		OutputSchema: schema,
		Filter:       filter,
		Fields:       []string{"status", "host"},
		Limit:        10,
	}

	jsonDSL, err := translator.TranslateToJSON(plan)
	require.NoError(t, err)

	// Should have query, _source, and size
	assert.Contains(t, jsonDSL, "query")
	assert.Contains(t, jsonDSL, "_source")
	assert.Contains(t, jsonDSL, "size")

	// Values should match
	assert.Equal(t, []string{"status", "host"}, jsonDSL["_source"])
	assert.Equal(t, 10, jsonDSL["size"])
}

func TestTranslator_ComplexAggregationWithFilter(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// Filter + GROUP BY aggregation
	filter := &ast.BinaryExpression{
		Left:     &ast.FieldReference{Name: "status"},
		Operator: ">=",
		Right:    &ast.Literal{Value: 400, LiteralTyp: ast.LiteralTypeInt},
	}

	agg := &physical.PhysicalAggregate{
		GroupBy: []ast.Expression{
			&ast.FieldReference{Name: "host"},
		},
		Aggregations: []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name:      "count",
					Arguments: []ast.Expression{},
				},
				Alias: "error_count",
			},
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
			Filter:       filter,
		},
	}

	dsl, err := translator.Translate(agg)
	require.NoError(t, err)

	// Should have filter query
	assert.NotNil(t, dsl.Query)
	rangeQuery, ok := dsl.Query["range"].(map[string]interface{})
	require.True(t, ok)
	assert.NotNil(t, rangeQuery["status"])

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)
	assert.Contains(t, dsl.Aggregations, "group_by_host")
}

// =====================================================================
// Tier 1 DSL Translator Tests
// =====================================================================

func TestTranslator_MultiFieldGroupBy(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// GROUP BY host, status
	agg := &physical.PhysicalAggregate{
		GroupBy: []ast.Expression{
			&ast.FieldReference{Name: "host"},
			&ast.FieldReference{Name: "status"},
		},
		Aggregations: []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name:      "count",
					Arguments: []ast.Expression{},
				},
				Alias: "total",
			},
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(agg)
	require.NoError(t, err)

	// Should have nested aggregations
	require.NotNil(t, dsl.Aggregations)

	// Should have group_by_host at outer level
	assert.Contains(t, dsl.Aggregations, "group_by_host")

	hostAgg, ok := dsl.Aggregations["group_by_host"].(map[string]interface{})
	require.True(t, ok)

	// Verify it's a terms aggregation
	termsAgg, ok := hostAgg["terms"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "host", termsAgg["field"])

	// Verify nested aggregation exists
	nestedAggs, ok := hostAgg["aggs"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, nestedAggs, "group_by_status")

	// Verify the inner terms aggregation
	statusAgg, ok := nestedAggs["group_by_status"].(map[string]interface{})
	require.True(t, ok)
	statusTerms, ok := statusAgg["terms"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "status", statusTerms["field"])
}

func TestTranslator_TopAggregation(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// top 10 host
	plan := &physical.PhysicalTop{
		Fields: []ast.Expression{
			&ast.FieldReference{Name: "host"},
		},
		Limit:        10,
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)

	// Size should be 0
	require.NotNil(t, dsl.Size)
	assert.Equal(t, 0, *dsl.Size)

	// Should have top_host aggregation
	assert.Contains(t, dsl.Aggregations, "top_host")

	topAgg, ok := dsl.Aggregations["top_host"].(map[string]interface{})
	require.True(t, ok)

	// Verify it's a terms aggregation with correct ordering
	termsAgg, ok := topAgg["terms"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "host", termsAgg["field"])
	assert.Equal(t, 10, termsAgg["size"])

	// Should be ordered by count descending
	order, ok := termsAgg["order"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "desc", order["_count"])
}

func TestTranslator_RareAggregation(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// rare 10 host
	plan := &physical.PhysicalRare{
		Fields: []ast.Expression{
			&ast.FieldReference{Name: "host"},
		},
		Limit:        10,
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)

	// Size should be 0
	require.NotNil(t, dsl.Size)
	assert.Equal(t, 0, *dsl.Size)

	// Should have rare_host aggregation
	assert.Contains(t, dsl.Aggregations, "rare_host")

	rareAgg, ok := dsl.Aggregations["rare_host"].(map[string]interface{})
	require.True(t, ok)

	// Verify it's a terms aggregation with correct ordering
	termsAgg, ok := rareAgg["terms"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "host", termsAgg["field"])
	assert.Equal(t, 10, termsAgg["size"])

	// Should be ordered by count ascending (rare)
	order, ok := termsAgg["order"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "asc", order["_count"])
}

func TestTranslator_BinDateHistogram(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// bin timestamp span=1h
	plan := &physical.PhysicalBin{
		Field: &ast.FieldReference{Name: "timestamp"},
		Span: &ast.TimeSpan{
			Value: 1,
			Unit:  "h",
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)

	// Size should be 0
	require.NotNil(t, dsl.Size)
	assert.Equal(t, 0, *dsl.Size)

	// Should have bin_timestamp aggregation
	assert.Contains(t, dsl.Aggregations, "bin_timestamp")

	binAgg, ok := dsl.Aggregations["bin_timestamp"].(map[string]interface{})
	require.True(t, ok)

	// Verify it's a date_histogram aggregation
	dateHist, ok := binAgg["date_histogram"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "timestamp", dateHist["field"])
	assert.Equal(t, "1h", dateHist["calendar_interval"])
}

func TestTranslator_BinNumericHistogram(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// bin latency bins=10
	plan := &physical.PhysicalBin{
		Field:        &ast.FieldReference{Name: "latency"},
		Bins:         10,
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)

	// Should have bin_latency aggregation
	assert.Contains(t, dsl.Aggregations, "bin_latency")

	binAgg, ok := dsl.Aggregations["bin_latency"].(map[string]interface{})
	require.True(t, ok)

	// Verify it's an auto_date_histogram aggregation (when bins is specified)
	autoHist, ok := binAgg["auto_date_histogram"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "latency", autoHist["field"])
	assert.Equal(t, 10, autoHist["buckets"])
}

func TestTranslator_TopWithMultipleFields(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// top 10 host, status
	plan := &physical.PhysicalTop{
		Fields: []ast.Expression{
			&ast.FieldReference{Name: "host"},
			&ast.FieldReference{Name: "status"},
		},
		Limit:        10,
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)

	// Should have top_host aggregation
	assert.Contains(t, dsl.Aggregations, "top_host")

	topAgg, ok := dsl.Aggregations["top_host"].(map[string]interface{})
	require.True(t, ok)

	// Should have nested aggregation for status
	nestedAggs, ok := topAgg["aggs"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, nestedAggs, "by_status")
}

func TestAggregationBuilder_TimeSpanToInterval(t *testing.T) {
	ab := NewAggregationBuilder()

	tests := []struct {
		span     *ast.TimeSpan
		expected string
	}{
		{&ast.TimeSpan{Value: 1, Unit: "s"}, "1s"},
		{&ast.TimeSpan{Value: 30, Unit: "s"}, "30s"},
		{&ast.TimeSpan{Value: 1, Unit: "m"}, "1m"},
		{&ast.TimeSpan{Value: 5, Unit: "min"}, "5m"},
		{&ast.TimeSpan{Value: 1, Unit: "h"}, "1h"},
		{&ast.TimeSpan{Value: 24, Unit: "hour"}, "24h"},
		{&ast.TimeSpan{Value: 1, Unit: "d"}, "1d"},
		{&ast.TimeSpan{Value: 7, Unit: "day"}, "7d"},
		{&ast.TimeSpan{Value: 1, Unit: "w"}, "1w"},
		{&ast.TimeSpan{Value: 1, Unit: "mon"}, "1M"},
		{&ast.TimeSpan{Value: 1, Unit: "y"}, "1y"},
	}

	for _, tt := range tests {
		t.Run(tt.span.String(), func(t *testing.T) {
			result := ab.timeSpanToInterval(tt.span)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTranslator_TopWithFilter(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// where status >= 400 | top 10 host
	filter := &ast.BinaryExpression{
		Left:     &ast.FieldReference{Name: "status"},
		Operator: ">=",
		Right:    &ast.Literal{Value: 400, LiteralTyp: ast.LiteralTypeInt},
	}

	plan := &physical.PhysicalTop{
		Fields: []ast.Expression{
			&ast.FieldReference{Name: "host"},
		},
		Limit:        10,
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
			Filter:       filter,
		},
	}

	dsl, err := translator.Translate(plan)
	require.NoError(t, err)

	// Should have filter query
	assert.NotNil(t, dsl.Query)
	rangeQuery, ok := dsl.Query["range"].(map[string]interface{})
	require.True(t, ok)
	assert.NotNil(t, rangeQuery["status"])

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)
	assert.Contains(t, dsl.Aggregations, "top_host")
}

func TestTranslator_ThreeFieldGroupBy(t *testing.T) {
	schema := createTestSchema()
	schema.AddField("region", analyzer.FieldTypeString)
	translator := NewTranslator()

	// stats count() by host, status, region
	agg := &physical.PhysicalAggregate{
		GroupBy: []ast.Expression{
			&ast.FieldReference{Name: "host"},
			&ast.FieldReference{Name: "status"},
			&ast.FieldReference{Name: "region"},
		},
		Aggregations: []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name:      "count",
					Arguments: []ast.Expression{},
				},
				Alias: "total",
			},
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(agg)
	require.NoError(t, err)

	// Should have nested aggregations 3 levels deep
	require.NotNil(t, dsl.Aggregations)

	// Navigate the nested structure
	// Level 1: group_by_host
	hostAgg, ok := dsl.Aggregations["group_by_host"].(map[string]interface{})
	require.True(t, ok)
	require.NotNil(t, hostAgg["terms"])
	require.NotNil(t, hostAgg["aggs"])

	// Level 2: group_by_status
	hostSubAggs := hostAgg["aggs"].(map[string]interface{})
	statusAgg, ok := hostSubAggs["group_by_status"].(map[string]interface{})
	require.True(t, ok)
	require.NotNil(t, statusAgg["terms"])
	require.NotNil(t, statusAgg["aggs"])

	// Level 3: group_by_region (with metric)
	statusSubAggs := statusAgg["aggs"].(map[string]interface{})
	regionAgg, ok := statusSubAggs["group_by_region"].(map[string]interface{})
	require.True(t, ok)
	require.NotNil(t, regionAgg["terms"])

	// Metric should be at the innermost level
	regionSubAggs := regionAgg["aggs"].(map[string]interface{})
	_, hasTotal := regionSubAggs["total"]
	assert.True(t, hasTotal)
}

func TestAggregationBuilder_CardinalityAggregation(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// stats dc(host)
	agg := &physical.PhysicalAggregate{
		Aggregations: []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name: "dc",
					Arguments: []ast.Expression{
						&ast.FieldReference{Name: "host"},
					},
				},
				Alias: "unique_hosts",
			},
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(agg)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)
	assert.Contains(t, dsl.Aggregations, "unique_hosts")

	cardinalityAgg, ok := dsl.Aggregations["unique_hosts"].(map[string]interface{})
	require.True(t, ok)

	cardinality, ok := cardinalityAgg["cardinality"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "host", cardinality["field"])
}

func TestAggregationBuilder_ExtendedStatsAggregation(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// stats stats(latency)
	agg := &physical.PhysicalAggregate{
		Aggregations: []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name: "stats",
					Arguments: []ast.Expression{
						&ast.FieldReference{Name: "latency"},
					},
				},
				Alias: "latency_stats",
			},
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(agg)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)
	assert.Contains(t, dsl.Aggregations, "latency_stats")

	statsAgg, ok := dsl.Aggregations["latency_stats"].(map[string]interface{})
	require.True(t, ok)

	extStats, ok := statsAgg["extended_stats"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "latency", extStats["field"])
}

func TestAggregationBuilder_PercentilesAggregation(t *testing.T) {
	schema := createTestSchema()
	translator := NewTranslator()

	// stats percentiles(latency)
	agg := &physical.PhysicalAggregate{
		Aggregations: []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name: "percentiles",
					Arguments: []ast.Expression{
						&ast.FieldReference{Name: "latency"},
					},
				},
				Alias: "latency_percentiles",
			},
		},
		OutputSchema: schema,
		Input: &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	dsl, err := translator.Translate(agg)
	require.NoError(t, err)

	// Should have aggregations
	require.NotNil(t, dsl.Aggregations)
	assert.Contains(t, dsl.Aggregations, "latency_percentiles")

	percAgg, ok := dsl.Aggregations["latency_percentiles"].(map[string]interface{})
	require.True(t, ok)

	percentiles, ok := percAgg["percentiles"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "latency", percentiles["field"])
}
