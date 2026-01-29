// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"strings"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// dedupOperator removes duplicate rows based on specified fields
type dedupOperator struct {
	input       Operator
	fields      []ast.Expression
	count       int  // Number of duplicates to keep
	consecutive bool // Only consider consecutive duplicates
	logger      *zap.Logger

	ctx     context.Context
	seen    map[string]int // Tracks count of each unique key
	prevKey string         // Previous key for consecutive mode
	stats   *IteratorStats
	opened  bool
	closed  bool
}

// NewDedupOperator creates a new dedup operator
func NewDedupOperator(
	input Operator,
	fields []ast.Expression,
	count int,
	consecutive bool,
	logger *zap.Logger,
) *dedupOperator {
	if count <= 0 {
		count = 1 // Default: keep only 1 of each
	}
	return &dedupOperator{
		input:       input,
		fields:      fields,
		count:       count,
		consecutive: consecutive,
		logger:      logger,
		seen:        make(map[string]int),
		stats:       &IteratorStats{},
	}
}

// Open initializes the operator
func (d *dedupOperator) Open(ctx context.Context) error {
	if d.opened {
		return nil
	}

	d.ctx = ctx
	d.logger.Debug("Opening dedup operator",
		zap.Int("num_fields", len(d.fields)),
		zap.Int("count", d.count),
		zap.Bool("consecutive", d.consecutive))

	if err := d.input.Open(ctx); err != nil {
		return err
	}

	d.opened = true
	return nil
}

// Next returns the next non-duplicate row
func (d *dedupOperator) Next(ctx context.Context) (*Row, error) {
	if d.closed {
		return nil, ErrClosed
	}

	if !d.opened {
		return nil, ErrClosed
	}

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		row, err := d.input.Next(ctx)
		if err != nil {
			return nil, err
		}

		d.stats.RowsRead++

		// Compute dedup key from fields
		key := d.computeKey(row)

		if d.consecutive {
			// Consecutive mode: only compare with previous row
			if key != d.prevKey {
				// New key, reset count
				d.seen = make(map[string]int)
				d.prevKey = key
			}
		}

		// Check if we've seen this key too many times
		count := d.seen[key]
		if count < d.count {
			d.seen[key] = count + 1
			d.stats.RowsReturned++
			return row, nil
		}

		// Skip this duplicate
	}
}

// computeKey computes a unique key for the row based on dedup fields
func (d *dedupOperator) computeKey(row *Row) string {
	parts := make([]string, len(d.fields))
	for i, field := range d.fields {
		if ref, ok := field.(*ast.FieldReference); ok {
			val, _ := row.Get(ref.Name)
			parts[i] = toString(val)
		}
	}
	return strings.Join(parts, "|")
}

// Close releases resources
func (d *dedupOperator) Close() error {
	d.closed = true
	d.seen = nil
	return d.input.Close()
}

// Stats returns execution statistics
func (d *dedupOperator) Stats() *IteratorStats {
	return d.stats
}
