// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// eventstatsOperator computes aggregations across all events and adds results to each row
// Unlike stats which groups rows, eventstats enriches each event with aggregate values
type eventstatsOperator struct {
	input        Operator
	groupBy      []ast.Expression
	aggregations []*ast.Aggregation
	logger       *zap.Logger

	ctx      context.Context
	stats    *IteratorStats
	opened   bool
	closed   bool
	rows     []*Row // Buffered rows
	rowIndex int    // Current position in buffered rows

	// Aggregation results keyed by group
	aggResults map[string]map[string]interface{}
}

// NewEventstatsOperator creates a new eventstats operator
func NewEventstatsOperator(
	input Operator,
	groupBy []ast.Expression,
	aggregations []*ast.Aggregation,
	logger *zap.Logger,
) *eventstatsOperator {
	return &eventstatsOperator{
		input:        input,
		groupBy:      groupBy,
		aggregations: aggregations,
		logger:       logger,
		stats:        &IteratorStats{},
		rows:         make([]*Row, 0),
		aggResults:   make(map[string]map[string]interface{}),
	}
}

// Open initializes the operator and computes aggregations
func (e *eventstatsOperator) Open(ctx context.Context) error {
	if e.opened {
		return nil
	}

	e.ctx = ctx
	e.logger.Debug("Opening eventstats operator",
		zap.Int("num_aggregations", len(e.aggregations)),
		zap.Int("num_group_by", len(e.groupBy)))

	if err := e.input.Open(ctx); err != nil {
		return err
	}

	// Read all rows and compute aggregations
	if err := e.computeAggregations(ctx); err != nil {
		return err
	}

	e.opened = true
	return nil
}

// computeAggregations reads all input rows, computes aggregations, and enriches rows
func (e *eventstatsOperator) computeAggregations(ctx context.Context) error {
	// Group rows by group keys
	groups := make(map[string][]*Row)
	rowGroups := make(map[*Row]string) // Track which group each row belongs to

	for {
		row, err := e.input.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return err
		}

		e.stats.RowsRead++

		// Determine group key for this row
		groupKey := e.getGroupKey(row)
		groups[groupKey] = append(groups[groupKey], row)
		rowGroups[row] = groupKey

		// Store original row
		e.rows = append(e.rows, row)
	}

	// Compute aggregations for each group
	for groupKey, groupRows := range groups {
		aggValues := make(map[string]interface{})

		for _, agg := range e.aggregations {
			aggName := agg.Alias
			if aggName == "" && agg.Func != nil {
				aggName = agg.Func.Name
			}

			// Compute aggregation value for this group
			value, err := e.computeAggregation(agg, groupRows)
			if err != nil {
				e.logger.Warn("Failed to compute aggregation",
					zap.String("aggregation", aggName),
					zap.Error(err))
				continue
			}

			aggValues[aggName] = value
		}

		e.aggResults[groupKey] = aggValues
	}

	// Enrich each row with aggregation results from its group
	for _, row := range e.rows {
		groupKey := rowGroups[row]
		aggValues := e.aggResults[groupKey]

		// Add aggregation results as new fields
		for fieldName, value := range aggValues {
			row.Set(fieldName, value)
		}
	}

	return nil
}

// getGroupKey generates a group key from group by expressions
func (e *eventstatsOperator) getGroupKey(row *Row) string {
	if len(e.groupBy) == 0 {
		return "" // Single group for all rows
	}

	key := ""
	for i, expr := range e.groupBy {
		if i > 0 {
			key += "|"
		}

		// Extract field value for group by
		switch ex := expr.(type) {
		case *ast.FieldReference:
			if val, exists := row.Get(ex.Name); exists {
				key += fmt.Sprintf("%v", val)
			} else {
				key += "NULL"
			}
		default:
			key += expr.String()
		}
	}

	return key
}

// computeAggregation computes an aggregation over a group of rows
func (e *eventstatsOperator) computeAggregation(agg *ast.Aggregation, rows []*Row) (interface{}, error) {
	if agg.Func == nil {
		return nil, fmt.Errorf("aggregation has no function")
	}

	funcName := agg.Func.Name
	var values []interface{}

	// Extract values from argument field if present
	if len(agg.Func.Arguments) > 0 {
		if fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
			for _, row := range rows {
				if val, exists := row.Get(fieldRef.Name); exists && val != nil {
					values = append(values, val)
				}
			}
		}
	}

	// Compute aggregation based on function name
	switch funcName {
	case "count":
		if agg.Func.Distinct {
			// Count distinct values
			seen := make(map[interface{}]bool)
			for _, val := range values {
				seen[val] = true
			}
			return int64(len(seen)), nil
		}
		// Count all rows in group
		return int64(len(rows)), nil

	case "sum":
		sum := float64(0)
		for _, val := range values {
			if num, ok := toFloat64(val); ok {
				sum += num
			}
		}
		return sum, nil

	case "avg":
		if len(values) == 0 {
			return float64(0), nil
		}
		sum := float64(0)
		for _, val := range values {
			if num, ok := toFloat64(val); ok {
				sum += num
			}
		}
		return sum / float64(len(values)), nil

	case "min":
		if len(values) == 0 {
			return nil, nil
		}
		min := values[0]
		for _, val := range values[1:] {
			if compareValues(val, min) < 0 {
				min = val
			}
		}
		return min, nil

	case "max":
		if len(values) == 0 {
			return nil, nil
		}
		max := values[0]
		for _, val := range values[1:] {
			if compareValues(val, max) > 0 {
				max = val
			}
		}
		return max, nil

	default:
		return nil, fmt.Errorf("unsupported aggregation function: %s", funcName)
	}
}

// toFloat64 is now in utils.go

// compareValues compares two values for min/max operations
func compareValues(a, b interface{}) int {
	// Try numeric comparison first
	aNum, aIsNum := toFloat64(a)
	bNum, bIsNum := toFloat64(b)
	if aIsNum && bIsNum {
		if aNum < bNum {
			return -1
		} else if aNum > bNum {
			return 1
		}
		return 0
	}

	// Fall back to string comparison
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)
	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

// Next returns the next enriched row
func (e *eventstatsOperator) Next(ctx context.Context) (*Row, error) {
	if e.closed {
		return nil, ErrClosed
	}

	if !e.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if e.rowIndex >= len(e.rows) {
		return nil, ErrNoMoreRows
	}

	row := e.rows[e.rowIndex]
	e.rowIndex++
	e.stats.RowsReturned++

	return row, nil
}

// Close releases resources
func (e *eventstatsOperator) Close() error {
	e.closed = true
	return e.input.Close()
}

// Stats returns execution statistics
func (e *eventstatsOperator) Stats() *IteratorStats {
	return e.stats
}
