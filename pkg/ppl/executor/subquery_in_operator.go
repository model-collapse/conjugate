// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// subqueryInOperator filters rows based on IN subquery
// Example: where status IN [search source=valid_statuses | fields status]
type subqueryInOperator struct {
	input          Operator
	logger         *zap.Logger
	ctx            context.Context
	stats          *IteratorStats
	opened         bool
	closed         bool
	fieldName      string            // Field to check
	subqueryExec   *SubqueryExecutor
	valueSet       map[interface{}]bool // Hash set for fast lookup
	negate         bool              // If true, NOT IN
}

// SubqueryInConfig holds configuration for IN subquery operator
type SubqueryInConfig struct {
	FieldName string          // Field to check
	Subsearch Operator        // Subsearch pipeline
	Negate    bool            // If true, NOT IN
}

// NewSubqueryInOperator creates a new IN subquery operator
func NewSubqueryInOperator(input Operator, config SubqueryInConfig, logger *zap.Logger) *subqueryInOperator {
	subqueryExec := NewSubqueryExecutor(config.Subsearch, SubqueryTypeIN, logger)

	return &subqueryInOperator{
		input:        input,
		logger:       logger,
		stats:        &IteratorStats{},
		fieldName:    config.FieldName,
		subqueryExec: subqueryExec,
		negate:       config.Negate,
	}
}

// Open initializes the operator and executes the subquery
func (s *subqueryInOperator) Open(ctx context.Context) error {
	if s.opened {
		return nil
	}

	s.ctx = ctx
	s.logger.Debug("Opening IN subquery operator",
		zap.String("field", s.fieldName),
		zap.Bool("negate", s.negate))

	// Open input
	if err := s.input.Open(ctx); err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}

	// Execute subquery once
	if err := s.subqueryExec.Execute(ctx); err != nil {
		return fmt.Errorf("failed to execute subquery: %w", err)
	}

	// Build hash set for O(1) lookup
	values := s.subqueryExec.GetFieldValues(s.fieldName)
	s.valueSet = make(map[interface{}]bool, len(values))

	for _, value := range values {
		// Normalize value for comparison
		normalizedValue := s.normalizeValue(value)
		s.valueSet[normalizedValue] = true
	}

	s.logger.Debug("Built IN value set",
		zap.Int("values", len(s.valueSet)))

	s.opened = true
	return nil
}

// Next filters rows based on IN condition
func (s *subqueryInOperator) Next(ctx context.Context) (*Row, error) {
	if s.closed {
		return nil, ErrClosed
	}

	if !s.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
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

		// Normalize and check if in set
		normalizedValue := s.normalizeValue(fieldValue)
		inSet := s.valueSet[normalizedValue]

		// Apply negate logic
		if s.negate {
			inSet = !inSet
		}

		if inSet {
			s.stats.RowsReturned++
			return row, nil
		}

		// Not in set, skip to next row
	}
}

// normalizeValue converts values to comparable form
func (s *subqueryInOperator) normalizeValue(value interface{}) interface{} {
	// Handle numeric type conversions for comparison
	switch v := value.(type) {
	case int:
		return int64(v)
	case int32:
		return int64(v)
	case uint16:
		return int64(v)
	case float32:
		return float64(v)
	default:
		return v
	}
}

// Close releases resources
func (s *subqueryInOperator) Close() error {
	s.closed = true
	if s.input != nil {
		return s.input.Close()
	}
	return nil
}

// Stats returns execution statistics
func (s *subqueryInOperator) Stats() *IteratorStats {
	return s.stats
}
