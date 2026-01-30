// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// streamstatsOperator computes running statistics in a streaming fashion
// Unlike eventstats which computes over all events, streamstats computes
// running aggregations as events are processed (running total, moving average, etc.)
type streamstatsOperator struct {
	input        Operator
	groupBy      []ast.Expression
	aggregations []*ast.Aggregation
	window       int         // Window size for rolling aggregations (0 = unbounded)
	current      bool        // Include current event in calculation (default: true)
	global       bool        // Compute stats globally, ignore grouping (default: false)
	resetBefore  ast.Expression // Reset statistics before this condition
	resetAfter   ast.Expression // Reset statistics after this condition
	logger       *zap.Logger

	ctx    context.Context
	stats  *IteratorStats
	opened bool
	closed bool

	// State for tracking running aggregations
	groupStates map[string]*streamState // State per group
}

// streamState tracks running aggregation state for a group
type streamState struct {
	windowValues map[string][]interface{} // Windowed values per aggregation
	runningCount int64                    // Running count
	runningSum   map[string]float64       // Running sum per aggregation
}

// NewStreamstatsOperator creates a new streamstats operator
func NewStreamstatsOperator(
	input Operator,
	groupBy []ast.Expression,
	aggregations []*ast.Aggregation,
	window int,
	global bool,
	resetBefore ast.Expression,
	resetAfter ast.Expression,
	logger *zap.Logger,
) *streamstatsOperator {
	return &streamstatsOperator{
		input:        input,
		groupBy:      groupBy,
		aggregations: aggregations,
		window:       window,
		current:      true, // Default: include current event
		global:       global,
		resetBefore:  resetBefore,
		resetAfter:   resetAfter,
		logger:       logger,
		stats:        &IteratorStats{},
		groupStates:  make(map[string]*streamState),
	}
}

// Open initializes the operator
func (s *streamstatsOperator) Open(ctx context.Context) error {
	if s.opened {
		return nil
	}

	s.ctx = ctx
	s.logger.Debug("Opening streamstats operator",
		zap.Int("num_aggregations", len(s.aggregations)),
		zap.Int("num_group_by", len(s.groupBy)),
		zap.Int("window", s.window))

	if err := s.input.Open(ctx); err != nil {
		return err
	}

	s.opened = true
	return nil
}

// Next returns the next row with running statistics added
func (s *streamstatsOperator) Next(ctx context.Context) (*Row, error) {
	if s.closed {
		return nil, ErrClosed
	}

	if !s.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	row, err := s.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	s.stats.RowsRead++

	// Determine group key for this row
	groupKey := s.getGroupKey(row)

	// Get or create state for this group
	state := s.groupStates[groupKey]
	if state == nil {
		state = &streamState{
			windowValues: make(map[string][]interface{}),
			runningSum:   make(map[string]float64),
			runningCount: 0,
		}
		s.groupStates[groupKey] = state
	}

	// Compute and add running statistics to the row
	for _, agg := range s.aggregations {
		aggName := agg.Alias
		if aggName == "" && agg.Func != nil {
			aggName = agg.Func.Name
		}

		// Compute running aggregation value
		value, err := s.computeRunningAggregation(agg, row, state)
		if err != nil {
			s.logger.Warn("Failed to compute running aggregation",
				zap.String("aggregation", aggName),
				zap.Error(err))
			continue
		}

		// Add result to row
		row.Set(aggName, value)
	}

	// Update state after computing aggregations
	state.runningCount++

	s.stats.RowsReturned++
	return row, nil
}

// getGroupKey generates a group key from group by expressions
func (s *streamstatsOperator) getGroupKey(row *Row) string {
	if len(s.groupBy) == 0 {
		return "" // Single group for all rows
	}

	key := ""
	for i, expr := range s.groupBy {
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

// computeRunningAggregation computes a running aggregation value
func (s *streamstatsOperator) computeRunningAggregation(
	agg *ast.Aggregation,
	row *Row,
	state *streamState,
) (interface{}, error) {
	if agg.Func == nil {
		return nil, fmt.Errorf("aggregation has no function")
	}

	funcName := agg.Func.Name
	aggKey := agg.Alias
	if aggKey == "" {
		aggKey = funcName
	}

	// Extract current value from row if there's an argument
	var currentValue interface{}
	if len(agg.Func.Arguments) > 0 {
		if fieldRef, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
			if val, exists := row.Get(fieldRef.Name); exists {
				currentValue = val
			}
		}
	}

	// Compute aggregation based on function name
	switch funcName {
	case "count":
		// Running count
		if s.window > 0 {
			// Window-based count
			values := state.windowValues[aggKey]
			values = append(values, 1)
			if len(values) > s.window {
				values = values[1:]
			}
			state.windowValues[aggKey] = values
			return int64(len(values)), nil
		}
		// Unbounded running count
		return state.runningCount + 1, nil

	case "sum":
		// Running sum
		currentNum, ok := toFloat64(currentValue)
		if !ok {
			currentNum = 0
		}

		if s.window > 0 {
			// Window-based sum
			values := state.windowValues[aggKey]
			values = append(values, currentValue)
			if len(values) > s.window {
				values = values[1:]
			}
			state.windowValues[aggKey] = values

			// Sum windowed values
			sum := float64(0)
			for _, val := range values {
				if num, ok := toFloat64(val); ok {
					sum += num
				}
			}
			return sum, nil
		}

		// Unbounded running sum
		state.runningSum[aggKey] += currentNum
		return state.runningSum[aggKey], nil

	case "avg":
		// Running average
		currentNum, ok := toFloat64(currentValue)
		if !ok {
			currentNum = 0
		}

		if s.window > 0 {
			// Window-based average
			values := state.windowValues[aggKey]
			values = append(values, currentValue)
			if len(values) > s.window {
				values = values[1:]
			}
			state.windowValues[aggKey] = values

			// Average windowed values
			sum := float64(0)
			count := 0
			for _, val := range values {
				if num, ok := toFloat64(val); ok {
					sum += num
					count++
				}
			}
			if count == 0 {
				return float64(0), nil
			}
			return sum / float64(count), nil
		}

		// Unbounded running average
		state.runningSum[aggKey] += currentNum
		return state.runningSum[aggKey] / float64(state.runningCount+1), nil

	case "min":
		// Running minimum
		if s.window > 0 {
			// Window-based min
			values := state.windowValues[aggKey]
			values = append(values, currentValue)
			if len(values) > s.window {
				values = values[1:]
			}
			state.windowValues[aggKey] = values

			if len(values) == 0 {
				return nil, nil
			}

			min := values[0]
			for _, val := range values[1:] {
				if val != nil && (min == nil || compareValues(val, min) < 0) {
					min = val
				}
			}
			return min, nil
		}

		// Unbounded running min
		values := state.windowValues[aggKey]
		values = append(values, currentValue)
		state.windowValues[aggKey] = values

		if len(values) == 0 {
			return nil, nil
		}

		min := values[0]
		for _, val := range values[1:] {
			if val != nil && (min == nil || compareValues(val, min) < 0) {
				min = val
			}
		}
		return min, nil

	case "max":
		// Running maximum
		if s.window > 0 {
			// Window-based max
			values := state.windowValues[aggKey]
			values = append(values, currentValue)
			if len(values) > s.window {
				values = values[1:]
			}
			state.windowValues[aggKey] = values

			if len(values) == 0 {
				return nil, nil
			}

			max := values[0]
			for _, val := range values[1:] {
				if val != nil && (max == nil || compareValues(val, max) > 0) {
					max = val
				}
			}
			return max, nil
		}

		// Unbounded running max
		values := state.windowValues[aggKey]
		values = append(values, currentValue)
		state.windowValues[aggKey] = values

		if len(values) == 0 {
			return nil, nil
		}

		max := values[0]
		for _, val := range values[1:] {
			if val != nil && (max == nil || compareValues(val, max) > 0) {
				max = val
			}
		}
		return max, nil

	default:
		return nil, fmt.Errorf("unsupported aggregation function: %s", funcName)
	}
}

// Close releases resources
func (s *streamstatsOperator) Close() error {
	s.closed = true
	return s.input.Close()
}

// Stats returns execution statistics
func (s *streamstatsOperator) Stats() *IteratorStats {
	return s.stats
}
