// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"go.uber.org/zap"
)

// aggregationOperator computes aggregations over rows
type aggregationOperator struct {
	input        Operator
	groupBy      []ast.Expression
	aggregations []*ast.Aggregation
	algorithm    physical.AggregationAlgorithm
	logger       *zap.Logger

	ctx          context.Context
	results      []*Row
	aggResults   map[string]*AggregationValue
	index        int
	stats        *IteratorStats
	opened       bool
	closed       bool
}

// NewAggregationOperator creates a new aggregation operator
func NewAggregationOperator(
	input Operator,
	groupBy []ast.Expression,
	aggregations []*ast.Aggregation,
	algorithm physical.AggregationAlgorithm,
	logger *zap.Logger,
) *aggregationOperator {
	return &aggregationOperator{
		input:        input,
		groupBy:      groupBy,
		aggregations: aggregations,
		algorithm:    algorithm,
		logger:       logger,
		stats:        &IteratorStats{},
		aggResults:   make(map[string]*AggregationValue),
	}
}

// Open initializes the operator by computing aggregations
func (a *aggregationOperator) Open(ctx context.Context) error {
	if a.opened {
		return nil
	}

	a.ctx = ctx
	a.logger.Debug("Opening aggregation operator",
		zap.Int("num_groups", len(a.groupBy)),
		zap.Int("num_aggs", len(a.aggregations)),
		zap.String("algorithm", a.algorithm.String()))

	// Open input
	if err := a.input.Open(ctx); err != nil {
		return err
	}

	// Compute aggregations
	if len(a.groupBy) == 0 {
		// No GROUP BY - compute global aggregations
		if err := a.computeGlobalAggregations(); err != nil {
			return err
		}
	} else {
		// With GROUP BY - use hash or stream aggregation
		if err := a.computeGroupedAggregations(); err != nil {
			return err
		}
	}

	a.index = 0
	a.opened = true

	return nil
}

// computeGlobalAggregations computes aggregations without grouping
func (a *aggregationOperator) computeGlobalAggregations() error {
	// Initialize accumulators
	accumulators := make(map[string]*aggregationAccumulator)
	for _, agg := range a.aggregations {
		name := agg.Alias
		if name == "" {
			name = agg.Func.Name
		}
		accumulators[name] = newAccumulator(agg.Func.Name)
	}

	// Process all rows
	for {
		row, err := a.input.Next(a.ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return err
		}

		a.stats.RowsRead++

		// Update accumulators
		for _, agg := range a.aggregations {
			name := agg.Alias
			if name == "" {
				name = agg.Func.Name
			}
			acc := accumulators[name]

			// Get value from row
			var value interface{}
			if len(agg.Func.Arguments) > 0 {
				if ref, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
					value, _ = row.Get(ref.Name)
				}
			}

			acc.update(value)
		}
	}

	// Create result row
	resultRow := NewRow(nil)
	for _, agg := range a.aggregations {
		name := agg.Alias
		if name == "" {
			name = agg.Func.Name
		}
		acc := accumulators[name]
		resultRow.Set(name, acc.result())

		// Store in aggregation results
		a.aggResults[name] = &AggregationValue{
			Type:  strings.ToLower(agg.Func.Name),
			Value: toFloat(acc.result()),
			Count: acc.count,
		}
	}

	a.results = []*Row{resultRow}
	a.stats.RowsReturned = 1

	return nil
}

// computeGroupedAggregations computes aggregations with grouping
func (a *aggregationOperator) computeGroupedAggregations() error {
	// Use hash-based grouping
	groups := make(map[string]*groupState)

	// Get group key field names
	groupFields := make([]string, len(a.groupBy))
	for i, expr := range a.groupBy {
		if ref, ok := expr.(*ast.FieldReference); ok {
			groupFields[i] = ref.Name
		}
	}

	// Process all rows
	for {
		row, err := a.input.Next(a.ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return err
		}

		a.stats.RowsRead++

		// Compute group key
		keyParts := make([]string, len(groupFields))
		for i, field := range groupFields {
			val, _ := row.Get(field)
			keyParts[i] = toString(val)
		}
		groupKey := strings.Join(keyParts, "|")

		// Get or create group state
		state, exists := groups[groupKey]
		if !exists {
			state = &groupState{
				keyValues:    make(map[string]interface{}),
				accumulators: make(map[string]*aggregationAccumulator),
			}
			for i, field := range groupFields {
				val, _ := row.Get(field)
				state.keyValues[field] = val
				state.keyParts = append(state.keyParts, keyParts[i])
			}
			for _, agg := range a.aggregations {
				name := agg.Alias
				if name == "" {
					name = agg.Func.Name
				}
				state.accumulators[name] = newAccumulator(agg.Func.Name)
			}
			groups[groupKey] = state
		}

		// Update accumulators
		for _, agg := range a.aggregations {
			name := agg.Alias
			if name == "" {
				name = agg.Func.Name
			}
			acc := state.accumulators[name]

			var value interface{}
			if len(agg.Func.Arguments) > 0 {
				if ref, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
					value, _ = row.Get(ref.Name)
				}
			}

			acc.update(value)
		}
	}

	// Create result rows
	a.results = make([]*Row, 0, len(groups))
	for _, state := range groups {
		row := NewRow(nil)

		// Add group key values
		for field, value := range state.keyValues {
			row.Set(field, value)
		}

		// Add aggregation results
		for _, agg := range a.aggregations {
			name := agg.Alias
			if name == "" {
				name = agg.Func.Name
			}
			acc := state.accumulators[name]
			row.Set(name, acc.result())
		}

		a.results = append(a.results, row)
	}

	a.stats.RowsReturned = int64(len(a.results))

	return nil
}

// Next returns the next aggregation result row
func (a *aggregationOperator) Next(ctx context.Context) (*Row, error) {
	if a.closed {
		return nil, ErrClosed
	}

	if !a.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if a.index >= len(a.results) {
		return nil, ErrNoMoreRows
	}

	row := a.results[a.index]
	a.index++

	return row, nil
}

// GetAggregations returns the aggregation results
func (a *aggregationOperator) GetAggregations() map[string]*AggregationValue {
	return a.aggResults
}

// Close releases resources
func (a *aggregationOperator) Close() error {
	a.closed = true
	a.results = nil
	return a.input.Close()
}

// Stats returns execution statistics
func (a *aggregationOperator) Stats() *IteratorStats {
	return a.stats
}

// groupState maintains state for a single group
type groupState struct {
	keyParts     []string
	keyValues    map[string]interface{}
	accumulators map[string]*aggregationAccumulator
}

// aggregationAccumulator maintains state for an aggregation function
type aggregationAccumulator struct {
	funcName string
	count    int64
	sum      float64
	min      float64
	max      float64
	hasMin   bool
	hasMax   bool
}

// newAccumulator creates a new accumulator for a function
func newAccumulator(funcName string) *aggregationAccumulator {
	return &aggregationAccumulator{
		funcName: strings.ToLower(funcName),
	}
}

// update updates the accumulator with a new value
func (acc *aggregationAccumulator) update(value interface{}) {
	acc.count++

	num, isNum := toNumber(value)
	if !isNum {
		return
	}

	acc.sum += num

	if !acc.hasMin || num < acc.min {
		acc.min = num
		acc.hasMin = true
	}

	if !acc.hasMax || num > acc.max {
		acc.max = num
		acc.hasMax = true
	}
}

// result returns the final aggregation result
func (acc *aggregationAccumulator) result() interface{} {
	switch acc.funcName {
	case "count":
		return acc.count
	case "sum":
		return acc.sum
	case "avg":
		if acc.count == 0 {
			return 0.0
		}
		return acc.sum / float64(acc.count)
	case "min":
		if !acc.hasMin {
			return nil
		}
		return acc.min
	case "max":
		if !acc.hasMax {
			return nil
		}
		return acc.max
	case "dc", "distinct_count", "cardinality":
		// Simplified - just return count for now
		return acc.count
	default:
		return fmt.Errorf("unsupported aggregation: %s", acc.funcName)
	}
}

// variance computes variance from sum and sum of squares
func variance(count int64, sum, sumSquares float64) float64 {
	if count <= 1 {
		return 0
	}
	mean := sum / float64(count)
	return (sumSquares / float64(count)) - (mean * mean)
}

// stdDev computes standard deviation
func stdDev(count int64, sum, sumSquares float64) float64 {
	return math.Sqrt(variance(count, sum, sumSquares))
}
