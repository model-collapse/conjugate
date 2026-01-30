// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// SubqueryScalarValue represents a scalar value from a subquery
// This can be used in comparison expressions
type SubqueryScalarValue struct {
	Value interface{}
	Valid bool
}

// scalarSubqueryOperator executes a scalar subquery and makes the result available
// Example: where revenue > [search source=benchmarks | stats avg(revenue) as threshold | fields threshold]
//
// Note: This operator doesn't filter rows itself; it executes the subquery and
// provides the scalar value for use in comparison expressions.
type scalarSubqueryOperator struct {
	input          Operator
	logger         *zap.Logger
	ctx            context.Context
	stats          *IteratorStats
	opened         bool
	closed         bool
	subqueryExec   *SubqueryExecutor
	scalarValue    SubqueryScalarValue
	comparisonOp   string            // Comparison operator: "=", "<", ">", "<=", ">=", "!="
	fieldName      string            // Field to compare
}

// SubqueryScalarConfig holds configuration for scalar subquery operator
type SubqueryScalarConfig struct {
	FieldName    string          // Field to compare
	Subsearch    Operator        // Subsearch pipeline
	ComparisonOp string          // Comparison operator
}

// NewScalarSubqueryOperator creates a new scalar subquery operator
func NewScalarSubqueryOperator(input Operator, config SubqueryScalarConfig, logger *zap.Logger) *scalarSubqueryOperator {
	subqueryExec := NewSubqueryExecutor(config.Subsearch, SubqueryTypeScalar, logger)

	return &scalarSubqueryOperator{
		input:        input,
		logger:       logger,
		stats:        &IteratorStats{},
		subqueryExec: subqueryExec,
		comparisonOp: config.ComparisonOp,
		fieldName:    config.FieldName,
	}
}

// Open initializes the operator and executes the scalar subquery
func (s *scalarSubqueryOperator) Open(ctx context.Context) error {
	if s.opened {
		return nil
	}

	s.ctx = ctx
	s.logger.Debug("Opening scalar subquery operator",
		zap.String("field", s.fieldName),
		zap.String("op", s.comparisonOp))

	// Open input
	if err := s.input.Open(ctx); err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}

	// Execute subquery
	if err := s.subqueryExec.Execute(ctx); err != nil {
		return fmt.Errorf("failed to execute subquery: %w", err)
	}

	// Get scalar value
	value, err := s.subqueryExec.GetScalarValue()
	if err != nil {
		s.logger.Warn("Scalar subquery error", zap.Error(err))
		s.scalarValue = SubqueryScalarValue{Valid: false}
	} else {
		s.scalarValue = SubqueryScalarValue{
			Value: value,
			Valid: true,
		}
		s.logger.Debug("Got scalar value",
			zap.Any("value", value))
	}

	s.opened = true
	return nil
}

// Next filters rows based on scalar comparison
func (s *scalarSubqueryOperator) Next(ctx context.Context) (*Row, error) {
	if s.closed {
		return nil, ErrClosed
	}

	if !s.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// If scalar value is invalid, return no rows
	if !s.scalarValue.Valid {
		return nil, ErrNoMoreRows
	}

	for {
		// Read next row
		row, err := s.input.Next(ctx)
		if err != nil {
			return nil, err
		}

		s.stats.RowsRead++

		// Get field value
		fieldValue, exists := row.Get(s.fieldName)
		if !exists {
			// Field doesn't exist - skip row
			continue
		}

		// Compare with scalar value
		if s.compareValues(fieldValue, s.scalarValue.Value) {
			s.stats.RowsReturned++
			return row, nil
		}

		// Comparison failed, skip to next row
	}
}

// compareValues performs the comparison based on the operator
func (s *scalarSubqueryOperator) compareValues(left, right interface{}) bool {
	// Convert to float64 for numeric comparisons
	leftNum, leftOk := toFloat64(left)
	rightNum, rightOk := toFloat64(right)

	if leftOk && rightOk {
		// Numeric comparison
		switch s.comparisonOp {
		case "=", "==":
			return leftNum == rightNum
		case "!=", "<>":
			return leftNum != rightNum
		case "<":
			return leftNum < rightNum
		case "<=":
			return leftNum <= rightNum
		case ">":
			return leftNum > rightNum
		case ">=":
			return leftNum >= rightNum
		default:
			s.logger.Warn("Unknown comparison operator", zap.String("op", s.comparisonOp))
			return false
		}
	}

	// String comparison (equality only)
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)

	switch s.comparisonOp {
	case "=", "==":
		return leftStr == rightStr
	case "!=", "<>":
		return leftStr != rightStr
	default:
		// Non-numeric comparison only supports equality
		s.logger.Warn("Non-numeric comparison only supports =, !=",
			zap.String("op", s.comparisonOp))
		return false
	}
}

// Close releases resources
func (s *scalarSubqueryOperator) Close() error {
	s.closed = true
	if s.input != nil {
		return s.input.Close()
	}
	return nil
}

// Stats returns execution statistics
func (s *scalarSubqueryOperator) Stats() *IteratorStats {
	return s.stats
}

// GetScalarValue returns the computed scalar value
func (s *scalarSubqueryOperator) GetScalarValue() SubqueryScalarValue {
	return s.scalarValue
}
