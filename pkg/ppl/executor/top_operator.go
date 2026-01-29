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

// topOperator returns the most frequent values for specified fields
type topOperator struct {
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

// NewTopOperator creates a new top operator
func NewTopOperator(
	input Operator,
	fields []ast.Expression,
	limit int,
	groupBy []ast.Expression,
	showCount bool,
	showPercent bool,
	logger *zap.Logger,
) *topOperator {
	if limit <= 0 {
		limit = 10 // Default top 10
	}
	return &topOperator{
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
func (t *topOperator) Open(ctx context.Context) error {
	if t.opened {
		return nil
	}

	t.ctx = ctx
	t.logger.Debug("Opening top operator",
		zap.Int("limit", t.limit),
		zap.Int("num_fields", len(t.fields)))

	if err := t.input.Open(ctx); err != nil {
		return err
	}

	// Count all values
	if err := t.computeTop(); err != nil {
		return err
	}

	t.index = 0
	t.opened = true
	return nil
}

// computeTop counts values and computes top N
func (t *topOperator) computeTop() error {
	// Count values by key
	counts := make(map[string]int64)
	var total int64

	for {
		row, err := t.input.Next(t.ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return err
		}

		t.stats.RowsRead++
		total++

		// Compute key from fields
		key := t.computeKey(row)
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

	// Sort by count descending (top = highest counts)
	sort.Slice(sortedCounts, func(i, j int) bool {
		return sortedCounts[i].count > sortedCounts[j].count
	})

	// Take top N
	limit := t.limit
	if limit > len(sortedCounts) {
		limit = len(sortedCounts)
	}

	// Build result rows
	t.results = make([]*Row, 0, limit)
	for i := 0; i < limit; i++ {
		kc := sortedCounts[i]
		row := NewRow(nil)

		// Parse key back into field values
		keyParts := strings.Split(kc.key, "|")
		for j, field := range t.fields {
			if ref, ok := field.(*ast.FieldReference); ok && j < len(keyParts) {
				row.Set(ref.Name, keyParts[j])
			}
		}

		// Add count
		row.Set("count", kc.count)

		// Add percent if requested
		if t.showPercent && total > 0 {
			percent := float64(kc.count) / float64(total) * 100.0
			row.Set("percent", percent)
		}

		t.results = append(t.results, row)
	}

	t.stats.RowsReturned = int64(len(t.results))
	return nil
}

// computeKey computes a unique key for the row based on top fields
func (t *topOperator) computeKey(row *Row) string {
	parts := make([]string, len(t.fields))
	for i, field := range t.fields {
		if ref, ok := field.(*ast.FieldReference); ok {
			val, _ := row.Get(ref.Name)
			parts[i] = toString(val)
		}
	}
	return strings.Join(parts, "|")
}

// Next returns the next top result
func (t *topOperator) Next(ctx context.Context) (*Row, error) {
	if t.closed {
		return nil, ErrClosed
	}

	if !t.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if t.index >= len(t.results) {
		return nil, ErrNoMoreRows
	}

	row := t.results[t.index]
	t.index++
	return row, nil
}

// Close releases resources
func (t *topOperator) Close() error {
	t.closed = true
	t.results = nil
	return t.input.Close()
}

// Stats returns execution statistics
func (t *topOperator) Stats() *IteratorStats {
	return t.stats
}
