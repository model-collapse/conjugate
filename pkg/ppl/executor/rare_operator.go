// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"sort"
	"strings"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// rareOperator returns the least frequent values for specified fields
type rareOperator struct {
	input       Operator
	fields      []ast.Expression
	limit       int
	groupBy     []ast.Expression
	showCount   bool
	showPercent bool
	logger      *zap.Logger

	ctx     context.Context
	results []*Row
	index   int
	stats   *IteratorStats
	opened  bool
	closed  bool
}

// NewRareOperator creates a new rare operator
func NewRareOperator(
	input Operator,
	fields []ast.Expression,
	limit int,
	groupBy []ast.Expression,
	showCount bool,
	showPercent bool,
	logger *zap.Logger,
) *rareOperator {
	if limit <= 0 {
		limit = 10 // Default rare 10
	}
	return &rareOperator{
		input:       input,
		fields:      fields,
		limit:       limit,
		groupBy:     groupBy,
		showCount:   showCount,
		showPercent: showPercent,
		logger:      logger,
		stats:       &IteratorStats{},
	}
}

// Open initializes the operator by counting all values
func (r *rareOperator) Open(ctx context.Context) error {
	if r.opened {
		return nil
	}

	r.ctx = ctx
	r.logger.Debug("Opening rare operator",
		zap.Int("limit", r.limit),
		zap.Int("num_fields", len(r.fields)))

	if err := r.input.Open(ctx); err != nil {
		return err
	}

	// Count all values
	if err := r.computeRare(); err != nil {
		return err
	}

	r.index = 0
	r.opened = true
	return nil
}

// computeRare counts values and computes rare N (lowest counts)
func (r *rareOperator) computeRare() error {
	// Count values by key
	counts := make(map[string]int64)
	var total int64

	for {
		row, err := r.input.Next(r.ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return err
		}

		r.stats.RowsRead++
		total++

		// Compute key from fields
		key := r.computeKey(row)
		counts[key]++
	}

	// Convert to sorted list
	type keyCount struct {
		key   string
		count int64
	}
	var sortedCounts []keyCount
	for k, c := range counts {
		sortedCounts = append(sortedCounts, keyCount{k, c})
	}

	// Sort by count ascending (rare = lowest counts)
	sort.Slice(sortedCounts, func(i, j int) bool {
		return sortedCounts[i].count < sortedCounts[j].count
	})

	// Take rare N
	limit := r.limit
	if limit > len(sortedCounts) {
		limit = len(sortedCounts)
	}

	// Build result rows
	r.results = make([]*Row, 0, limit)
	for i := 0; i < limit; i++ {
		kc := sortedCounts[i]
		row := NewRow(nil)

		// Parse key back into field values
		keyParts := strings.Split(kc.key, "|")
		for j, field := range r.fields {
			if ref, ok := field.(*ast.FieldReference); ok && j < len(keyParts) {
				row.Set(ref.Name, keyParts[j])
			}
		}

		// Add count
		row.Set("count", kc.count)

		// Add percent if requested
		if r.showPercent && total > 0 {
			percent := float64(kc.count) / float64(total) * 100.0
			row.Set("percent", percent)
		}

		r.results = append(r.results, row)
	}

	r.stats.RowsReturned = int64(len(r.results))
	return nil
}

// computeKey computes a unique key for the row based on rare fields
func (r *rareOperator) computeKey(row *Row) string {
	parts := make([]string, len(r.fields))
	for i, field := range r.fields {
		if ref, ok := field.(*ast.FieldReference); ok {
			val, _ := row.Get(ref.Name)
			parts[i] = toString(val)
		}
	}
	return strings.Join(parts, "|")
}

// Next returns the next rare result
func (r *rareOperator) Next(ctx context.Context) (*Row, error) {
	if r.closed {
		return nil, ErrClosed
	}

	if !r.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if r.index >= len(r.results) {
		return nil, ErrNoMoreRows
	}

	row := r.results[r.index]
	r.index++
	return row, nil
}

// Close releases resources
func (r *rareOperator) Close() error {
	r.closed = true
	r.results = nil
	return r.input.Close()
}

// Stats returns execution statistics
func (r *rareOperator) Stats() *IteratorStats {
	return r.stats
}
