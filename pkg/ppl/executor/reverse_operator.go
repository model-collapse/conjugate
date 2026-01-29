// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"go.uber.org/zap"
)

// reverseOperator reverses the order of rows in the result set
// It buffers all rows from the input, then emits them in reverse order
type reverseOperator struct {
	input  Operator
	logger *zap.Logger
	stats  *IteratorStats

	// State
	ctx     context.Context
	buffer  []*Row // Buffer to hold all rows
	index   int    // Current index for emitting rows
	buffered bool   // True when all rows have been buffered
	opened  bool
	closed  bool
}

// NewReverseOperator creates a new reverse operator
func NewReverseOperator(input Operator, logger *zap.Logger) *reverseOperator {
	return &reverseOperator{
		input:  input,
		logger: logger,
		stats:  &IteratorStats{},
	}
}

// Open initializes the operator
func (r *reverseOperator) Open(ctx context.Context) error {
	if r.opened {
		return nil
	}

	r.ctx = ctx
	r.logger.Debug("Opening reverse operator")

	// Open input operator
	if err := r.input.Open(ctx); err != nil {
		return err
	}

	r.opened = true
	return nil
}

// Next returns the next row in reverse order
func (r *reverseOperator) Next(ctx context.Context) (*Row, error) {
	if r.closed {
		return nil, ErrClosed
	}

	if !r.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// First call: buffer all rows from input
	if !r.buffered {
		r.logger.Debug("Buffering all rows for reversal")

		r.buffer = make([]*Row, 0, 1000) // Pre-allocate with reasonable capacity

		// Read all rows from input
		for {
			row, err := r.input.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			if err != nil {
				return nil, err
			}

			r.buffer = append(r.buffer, row)
			r.stats.RowsRead++
		}

		r.logger.Debug("Buffered all rows",
			zap.Int("total_rows", len(r.buffer)))

		r.buffered = true
		r.index = len(r.buffer) - 1 // Start from the last row
	}

	// Emit rows in reverse order
	if r.index >= 0 {
		row := r.buffer[r.index]
		r.index--
		r.stats.RowsReturned++
		return row, nil
	}

	// No more rows
	return nil, ErrNoMoreRows
}

// Close releases resources
func (r *reverseOperator) Close() error {
	r.closed = true

	// Clear buffer to free memory
	r.buffer = nil

	// Close input
	if r.input != nil {
		return r.input.Close()
	}

	return nil
}

// Stats returns execution statistics
func (r *reverseOperator) Stats() *IteratorStats {
	return r.stats
}
