// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// SubqueryType represents the type of subquery
type SubqueryType int

const (
	SubqueryTypeIN SubqueryType = iota     // IN subquery: field IN [search ...]
	SubqueryTypeEXISTS                     // EXISTS subquery: EXISTS [search ...]
	SubqueryTypeScalar                     // Scalar subquery: field > [search ... | stats avg(x)]
)

// SubqueryExecutor executes a subquery and returns results
type SubqueryExecutor struct {
	subsearch      Operator          // The subsearch pipeline
	logger         *zap.Logger
	maxRows        int               // Maximum rows to materialize (default: 10000)
	resultCache    []*Row            // Cached results
	executed       bool              // Whether subquery has been executed
	subqueryType   SubqueryType
}

// NewSubqueryExecutor creates a new subquery executor
func NewSubqueryExecutor(subsearch Operator, subqueryType SubqueryType, logger *zap.Logger) *SubqueryExecutor {
	return &SubqueryExecutor{
		subsearch:    subsearch,
		subqueryType: subqueryType,
		logger:       logger,
		maxRows:      10000, // Default limit
	}
}

// Execute runs the subquery and materializes results
func (s *SubqueryExecutor) Execute(ctx context.Context) error {
	if s.executed {
		return nil // Already executed
	}

	s.logger.Debug("Executing subquery",
		zap.Int("type", int(s.subqueryType)),
		zap.Int("max_rows", s.maxRows))

	// Open subsearch
	if err := s.subsearch.Open(ctx); err != nil {
		return fmt.Errorf("failed to open subsearch: %w", err)
	}
	defer s.subsearch.Close()

	// Materialize results
	s.resultCache = make([]*Row, 0)
	rowCount := 0

	for {
		row, err := s.subsearch.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return fmt.Errorf("subsearch error: %w", err)
		}

		s.resultCache = append(s.resultCache, row)
		rowCount++

		// Check limit
		if rowCount >= s.maxRows {
			s.logger.Warn("Subquery hit row limit",
				zap.Int("limit", s.maxRows))
			break
		}
	}

	s.logger.Debug("Subquery executed",
		zap.Int("rows", len(s.resultCache)))

	s.executed = true
	return nil
}

// GetResults returns the cached results
func (s *SubqueryExecutor) GetResults() []*Row {
	return s.resultCache
}

// GetScalarValue returns a single scalar value from the subquery
// Used for scalar subqueries
func (s *SubqueryExecutor) GetScalarValue() (interface{}, error) {
	if !s.executed {
		return nil, fmt.Errorf("subquery not executed")
	}

	if len(s.resultCache) == 0 {
		return nil, fmt.Errorf("scalar subquery returned no rows")
	}

	if len(s.resultCache) > 1 {
		return nil, fmt.Errorf("scalar subquery returned multiple rows (%d)", len(s.resultCache))
	}

	row := s.resultCache[0]
	fields := row.Fields()

	if len(fields) == 0 {
		return nil, fmt.Errorf("scalar subquery returned no fields")
	}

	if len(fields) > 1 {
		return nil, fmt.Errorf("scalar subquery returned multiple fields (%d)", len(fields))
	}

	value, _ := row.Get(fields[0])
	return value, nil
}

// GetFieldValues extracts all values of a specific field from results
// Used for IN subqueries
func (s *SubqueryExecutor) GetFieldValues(fieldName string) []interface{} {
	if !s.executed {
		return nil
	}

	values := make([]interface{}, 0, len(s.resultCache))
	for _, row := range s.resultCache {
		if value, exists := row.Get(fieldName); exists {
			values = append(values, value)
		}
	}

	return values
}

// HasResults returns true if the subquery returned any rows
// Used for EXISTS subqueries
func (s *SubqueryExecutor) HasResults() bool {
	return s.executed && len(s.resultCache) > 0
}

// SetMaxRows sets the maximum number of rows to materialize
func (s *SubqueryExecutor) SetMaxRows(maxRows int) {
	s.maxRows = maxRows
}
