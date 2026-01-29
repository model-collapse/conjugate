// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package dsl

import (
	"fmt"
	"strings"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
)

// AggregationBuilder builds OpenSearch aggregation DSL
type AggregationBuilder struct {
	functionBuilder FunctionBuilder
}

// NewAggregationBuilder creates a new aggregation builder
func NewAggregationBuilder() *AggregationBuilder {
	return &AggregationBuilder{}
}

// BuildTopAggregations builds aggregation DSL for top command
func (ab *AggregationBuilder) BuildTopAggregations(top *physical.PhysicalTop) (map[string]interface{}, error) {
	return ab.buildTopRareAggregations(top.Fields, top.Limit, top.GroupBy, true)
}

// BuildRareAggregations builds aggregation DSL for rare command
func (ab *AggregationBuilder) BuildRareAggregations(rare *physical.PhysicalRare) (map[string]interface{}, error) {
	return ab.buildTopRareAggregations(rare.Fields, rare.Limit, rare.GroupBy, false)
}

// buildTopRareAggregations builds terms aggregation with ordering for top/rare
func (ab *AggregationBuilder) buildTopRareAggregations(fields []ast.Expression, limit int, groupBy []ast.Expression, isTop bool) (map[string]interface{}, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("top/rare requires at least one field")
	}

	// Get the primary field to count
	primaryField, ok := fields[0].(*ast.FieldReference)
	if !ok {
		return nil, fmt.Errorf("top/rare field must be a simple field reference")
	}

	// Order by count (desc for top, asc for rare)
	order := "desc"
	if !isTop {
		order = "asc"
	}

	// Build the terms aggregation
	termsAgg := map[string]interface{}{
		"terms": map[string]interface{}{
			"field": primaryField.Name,
			"size":  limit,
			"order": map[string]interface{}{
				"_count": order,
			},
		},
	}

	// If there are multiple fields, create nested aggregations
	if len(fields) > 1 {
		innerAgg := termsAgg
		for i := 1; i < len(fields); i++ {
			fieldRef, ok := fields[i].(*ast.FieldReference)
			if !ok {
				continue
			}

			nestedTerms := map[string]interface{}{
				"terms": map[string]interface{}{
					"field": fieldRef.Name,
					"size":  limit,
					"order": map[string]interface{}{
						"_count": order,
					},
				},
			}

			if _, hasAggs := innerAgg["aggs"]; !hasAggs {
				innerAgg["aggs"] = make(map[string]interface{})
			}
			aggName := fmt.Sprintf("by_%s", fieldRef.Name)
			innerAgg["aggs"].(map[string]interface{})[aggName] = nestedTerms
			innerAgg = nestedTerms
		}
	}

	// If there's GROUP BY, wrap in the group terms
	if len(groupBy) > 0 {
		return ab.wrapWithGroupBy(termsAgg, groupBy, primaryField.Name)
	}

	aggName := fmt.Sprintf("top_%s", primaryField.Name)
	if !isTop {
		aggName = fmt.Sprintf("rare_%s", primaryField.Name)
	}

	return map[string]interface{}{
		aggName: termsAgg,
	}, nil
}

// BuildBinAggregations builds aggregation DSL for bin command
func (ab *AggregationBuilder) BuildBinAggregations(bin *physical.PhysicalBin) (map[string]interface{}, error) {
	fieldRef, ok := bin.Field.(*ast.FieldReference)
	if !ok {
		return nil, fmt.Errorf("bin field must be a simple field reference")
	}

	var agg map[string]interface{}

	if bin.Span != nil {
		// Time-based binning - use date_histogram
		interval := ab.timeSpanToInterval(bin.Span)
		agg = map[string]interface{}{
			"date_histogram": map[string]interface{}{
				"field":             fieldRef.Name,
				"calendar_interval": interval,
			},
		}
	} else {
		// Numeric binning - use histogram
		// Calculate interval based on bins count (will be refined at runtime)
		agg = map[string]interface{}{
			"histogram": map[string]interface{}{
				"field":         fieldRef.Name,
				"min_doc_count": 0,
			},
		}
		if bin.Bins > 0 {
			// Use auto_date_histogram if bins is specified for better interval selection
			agg = map[string]interface{}{
				"auto_date_histogram": map[string]interface{}{
					"field":   fieldRef.Name,
					"buckets": bin.Bins,
				},
			}
		}
	}

	aggName := fmt.Sprintf("bin_%s", fieldRef.Name)
	return map[string]interface{}{
		aggName: agg,
	}, nil
}

// timeSpanToInterval converts TimeSpan to OpenSearch interval
func (ab *AggregationBuilder) timeSpanToInterval(span *ast.TimeSpan) string {
	unit := strings.ToLower(span.Unit)
	value := span.Value

	switch unit {
	case "s", "sec", "second", "seconds":
		if value == 1 {
			return "1s"
		}
		return fmt.Sprintf("%ds", value)
	case "m", "min", "minute", "minutes":
		if value == 1 {
			return "1m"
		}
		return fmt.Sprintf("%dm", value)
	case "h", "hr", "hour", "hours":
		if value == 1 {
			return "1h"
		}
		return fmt.Sprintf("%dh", value)
	case "d", "day", "days":
		if value == 1 {
			return "1d"
		}
		return fmt.Sprintf("%dd", value)
	case "w", "week", "weeks":
		if value == 1 {
			return "1w"
		}
		return fmt.Sprintf("%dw", value)
	case "mon", "month", "months":
		if value == 1 {
			return "1M"
		}
		return fmt.Sprintf("%dM", value)
	case "q", "quarter", "quarters":
		if value == 1 {
			return "1q"
		}
		return fmt.Sprintf("%dq", value)
	case "y", "year", "years":
		if value == 1 {
			return "1y"
		}
		return fmt.Sprintf("%dy", value)
	default:
		// Default to minutes
		return fmt.Sprintf("%dm", value)
	}
}

// wrapWithGroupBy wraps an aggregation in group-by terms
func (ab *AggregationBuilder) wrapWithGroupBy(innerAgg map[string]interface{}, groupBy []ast.Expression, innerName string) (map[string]interface{}, error) {
	if len(groupBy) == 0 {
		return map[string]interface{}{
			innerName: innerAgg,
		}, nil
	}

	// Build nested group-by from outermost to innermost
	result := innerAgg
	for i := len(groupBy) - 1; i >= 0; i-- {
		groupField, ok := groupBy[i].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("GROUP BY field must be a simple field reference")
		}

		groupAgg := map[string]interface{}{
			"terms": map[string]interface{}{
				"field": groupField.Name,
				"size":  10000,
			},
			"aggs": map[string]interface{}{
				innerName: result,
			},
		}

		innerName = fmt.Sprintf("group_by_%s", groupField.Name)
		result = groupAgg
	}

	return map[string]interface{}{
		innerName: result,
	}, nil
}

// BuildAggregations builds aggregation DSL from a physical aggregate node
func (ab *AggregationBuilder) BuildAggregations(agg *physical.PhysicalAggregate) (map[string]interface{}, error) {
	aggs := make(map[string]interface{})

	// If there's a GROUP BY, create a terms aggregation
	if len(agg.GroupBy) > 0 {
		return ab.buildGroupByAggregations(agg)
	}

	// No GROUP BY - just compute metrics
	for _, aggFunc := range agg.Aggregations {
		metricAgg, err := ab.buildMetricAggregation(aggFunc)
		if err != nil {
			return nil, err
		}

		aggName := aggFunc.Alias
		if aggName == "" {
			aggName = ab.generateAggName(aggFunc)
		}

		aggs[aggName] = metricAgg
	}

	return aggs, nil
}

// buildGroupByAggregations builds nested aggregations for GROUP BY
func (ab *AggregationBuilder) buildGroupByAggregations(agg *physical.PhysicalAggregate) (map[string]interface{}, error) {
	if len(agg.GroupBy) == 0 {
		return nil, fmt.Errorf("GROUP BY requires at least one field")
	}

	// Build sub-aggregations for each metric
	subAggs := make(map[string]interface{})
	for _, aggFunc := range agg.Aggregations {
		metricAgg, err := ab.buildMetricAggregation(aggFunc)
		if err != nil {
			return nil, err
		}

		aggName := aggFunc.Alias
		if aggName == "" {
			aggName = ab.generateAggName(aggFunc)
		}

		subAggs[aggName] = metricAgg
	}

	// Build nested terms aggregations from innermost to outermost
	// For GROUP BY field1, field2, field3 we get:
	// terms(field1) -> terms(field2) -> terms(field3) -> metrics
	currentAggs := subAggs

	for i := len(agg.GroupBy) - 1; i >= 0; i-- {
		groupField, ok := agg.GroupBy[i].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("GROUP BY field must be a simple field reference")
		}

		termsAgg := map[string]interface{}{
			"terms": map[string]interface{}{
				"field": groupField.Name,
				"size":  10000, // Default size - could be configurable
			},
		}

		// If there are sub-aggregations (metrics or nested groups), add them
		if len(currentAggs) > 0 {
			termsAgg["aggs"] = currentAggs
		}

		groupName := fmt.Sprintf("group_by_%s", groupField.Name)

		// This becomes the current aggregation for the next (outer) level
		currentAggs = map[string]interface{}{
			groupName: termsAgg,
		}
	}

	return currentAggs, nil
}

// buildMetricAggregation builds a single metric aggregation
func (ab *AggregationBuilder) buildMetricAggregation(agg *ast.Aggregation) (map[string]interface{}, error) {
	funcName := strings.ToLower(agg.Func.Name)

	switch funcName {
	case "count":
		// Count aggregation
		if len(agg.Func.Arguments) == 0 || ab.isWildcard(agg.Func.Arguments[0]) {
			// count() or count(*) - use value_count on _id
			return map[string]interface{}{
				"value_count": map[string]interface{}{
					"field": "_id",
				},
			}, nil
		}

		// count(field) - use value_count on the field
		fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("count() argument must be a field reference")
		}

		return map[string]interface{}{
			"value_count": map[string]interface{}{
				"field": fieldRef.Name,
			},
		}, nil

	case "sum":
		if len(agg.Func.Arguments) != 1 {
			return nil, fmt.Errorf("sum() requires exactly one argument")
		}

		fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("sum() argument must be a field reference")
		}

		return map[string]interface{}{
			"sum": map[string]interface{}{
				"field": fieldRef.Name,
			},
		}, nil

	case "avg":
		if len(agg.Func.Arguments) != 1 {
			return nil, fmt.Errorf("avg() requires exactly one argument")
		}

		fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("avg() argument must be a field reference")
		}

		return map[string]interface{}{
			"avg": map[string]interface{}{
				"field": fieldRef.Name,
			},
		}, nil

	case "min":
		if len(agg.Func.Arguments) != 1 {
			return nil, fmt.Errorf("min() requires exactly one argument")
		}

		fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("min() argument must be a field reference")
		}

		return map[string]interface{}{
			"min": map[string]interface{}{
				"field": fieldRef.Name,
			},
		}, nil

	case "max":
		if len(agg.Func.Arguments) != 1 {
			return nil, fmt.Errorf("max() requires exactly one argument")
		}

		fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("max() argument must be a field reference")
		}

		return map[string]interface{}{
			"max": map[string]interface{}{
				"field": fieldRef.Name,
			},
		}, nil

	case "cardinality", "dc", "distinct_count":
		// Distinct count
		if len(agg.Func.Arguments) != 1 {
			return nil, fmt.Errorf("%s() requires exactly one argument", funcName)
		}

		fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("%s() argument must be a field reference", funcName)
		}

		return map[string]interface{}{
			"cardinality": map[string]interface{}{
				"field": fieldRef.Name,
			},
		}, nil

	case "stats":
		// Extended stats aggregation
		if len(agg.Func.Arguments) != 1 {
			return nil, fmt.Errorf("stats() requires exactly one argument")
		}

		fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("stats() argument must be a field reference")
		}

		return map[string]interface{}{
			"extended_stats": map[string]interface{}{
				"field": fieldRef.Name,
			},
		}, nil

	case "percentiles":
		if len(agg.Func.Arguments) != 1 {
			return nil, fmt.Errorf("percentiles() requires exactly one argument")
		}

		fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference)
		if !ok {
			return nil, fmt.Errorf("percentiles() argument must be a field reference")
		}

		return map[string]interface{}{
			"percentiles": map[string]interface{}{
				"field": fieldRef.Name,
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported aggregation function: %s", funcName)
	}
}

// isWildcard checks if an expression is a wildcard (*)
func (ab *AggregationBuilder) isWildcard(expr ast.Expression) bool {
	if fieldRef, ok := expr.(*ast.FieldReference); ok {
		return fieldRef.Name == "*"
	}
	return false
}

// generateAggName generates a default aggregation name
func (ab *AggregationBuilder) generateAggName(agg *ast.Aggregation) string {
	funcName := strings.ToLower(agg.Func.Name)

	if len(agg.Func.Arguments) == 0 {
		return funcName
	}

	if fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
		return fmt.Sprintf("%s_%s", funcName, fieldRef.Name)
	}

	return funcName
}
