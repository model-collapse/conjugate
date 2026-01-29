// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"sort"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// sortOperator sorts rows by specified keys
type sortOperator struct {
	input    Operator
	sortKeys []*ast.SortKey
	logger   *zap.Logger

	ctx    context.Context
	rows   []*Row
	index  int
	stats  *IteratorStats
	opened bool
	closed bool
}

// NewSortOperator creates a new sort operator
func NewSortOperator(input Operator, sortKeys []*ast.SortKey, logger *zap.Logger) *sortOperator {
	return &sortOperator{
		input:    input,
		sortKeys: sortKeys,
		logger:   logger,
		stats:    &IteratorStats{},
	}
}

// Open initializes the operator by loading and sorting all rows
func (s *sortOperator) Open(ctx context.Context) error {
	if s.opened {
		return nil
	}

	s.ctx = ctx
	s.logger.Debug("Opening sort operator",
		zap.Int("num_keys", len(s.sortKeys)))

	// Open input
	if err := s.input.Open(ctx); err != nil {
		return err
	}

	// Load all rows (sort requires full materialization)
	s.rows = make([]*Row, 0)
	for {
		row, err := s.input.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return err
		}
		s.rows = append(s.rows, row)
		s.stats.RowsRead++
	}

	// Sort rows
	sort.SliceStable(s.rows, func(i, j int) bool {
		return s.compareRows(s.rows[i], s.rows[j]) < 0
	})

	s.index = 0
	s.opened = true

	s.logger.Debug("Sort operator opened",
		zap.Int("rows_sorted", len(s.rows)))

	return nil
}

// compareRows compares two rows based on sort keys
func (s *sortOperator) compareRows(a, b *Row) int {
	for _, key := range s.sortKeys {
		fieldName := extractFieldName(key.Field)
		aVal, _ := a.Get(fieldName)
		bVal, _ := b.Get(fieldName)

		cmp := compare(aVal, bVal)
		if cmp != 0 {
			if key.Descending {
				return -cmp
			}
			return cmp
		}
	}
	return 0
}

// Next returns the next sorted row
func (s *sortOperator) Next(ctx context.Context) (*Row, error) {
	if s.closed {
		return nil, ErrClosed
	}

	if !s.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if s.index >= len(s.rows) {
		return nil, ErrNoMoreRows
	}

	row := s.rows[s.index]
	s.index++
	s.stats.RowsReturned++

	return row, nil
}

// Close releases resources
func (s *sortOperator) Close() error {
	s.closed = true
	s.rows = nil
	return s.input.Close()
}

// Stats returns execution statistics
func (s *sortOperator) Stats() *IteratorStats {
	return s.stats
}

// extractFieldName extracts the field name from an expression
func extractFieldName(expr ast.Expression) string {
	if ref, ok := expr.(*ast.FieldReference); ok {
		return ref.Name
	}
	return ""
}
