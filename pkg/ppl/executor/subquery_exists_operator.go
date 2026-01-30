// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// subqueryExistsOperator filters rows based on EXISTS subquery
// Example: where EXISTS [search source=related_data | where related.id = main.id]
//
// Note: This is a simplified implementation that only supports uncorrelated EXISTS.
// Correlated EXISTS (with references to outer query) requires more complex implementation.
type subqueryExistsOperator struct {
	input          Operator
	logger         *zap.Logger
	ctx            context.Context
	stats          *IteratorStats
	opened         bool
	closed         bool
	subqueryExec   *SubqueryExecutor
	negate         bool              // If true, NOT EXISTS
	existsResult   bool              // Cached result of EXISTS check
}

// SubqueryExistsConfig holds configuration for EXISTS subquery operator
type SubqueryExistsConfig struct {
	Subsearch Operator        // Subsearch pipeline
	Negate    bool            // If true, NOT EXISTS
}

// NewSubqueryExistsOperator creates a new EXISTS subquery operator
func NewSubqueryExistsOperator(input Operator, config SubqueryExistsConfig, logger *zap.Logger) *subqueryExistsOperator {
	subqueryExec := NewSubqueryExecutor(config.Subsearch, SubqueryTypeEXISTS, logger)

	return &subqueryExistsOperator{
		input:        input,
		logger:       logger,
		stats:        &IteratorStats{},
		subqueryExec: subqueryExec,
		negate:       config.Negate,
	}
}

// Open initializes the operator and executes the subquery
func (s *subqueryExistsOperator) Open(ctx context.Context) error {
	if s.opened {
		return nil
	}

	s.ctx = ctx
	s.logger.Debug("Opening EXISTS subquery operator",
		zap.Bool("negate", s.negate))

	// Open input
	if err := s.input.Open(ctx); err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}

	// Execute subquery once
	if err := s.subqueryExec.Execute(ctx); err != nil {
		return fmt.Errorf("failed to execute subquery: %w", err)
	}

	// Check if subquery returned any results
	s.existsResult = s.subqueryExec.HasResults()

	// Apply negate logic
	if s.negate {
		s.existsResult = !s.existsResult
	}

	s.logger.Debug("EXISTS check result",
		zap.Bool("exists", s.existsResult),
		zap.Int("subquery_rows", len(s.subqueryExec.GetResults())))

	s.opened = true
	return nil
}

// Next filters rows based on EXISTS condition
func (s *subqueryExistsOperator) Next(ctx context.Context) (*Row, error) {
	if s.closed {
		return nil, ErrClosed
	}

	if !s.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// If EXISTS is false, return no rows
	if !s.existsResult {
		return nil, ErrNoMoreRows
	}

	// If EXISTS is true, pass through all rows
	row, err := s.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	s.stats.RowsRead++
	s.stats.RowsReturned++
	return row, nil
}

// Close releases resources
func (s *subqueryExistsOperator) Close() error {
	s.closed = true
	if s.input != nil {
		return s.input.Close()
	}
	return nil
}

// Stats returns execution statistics
func (s *subqueryExistsOperator) Stats() *IteratorStats {
	return s.stats
}
