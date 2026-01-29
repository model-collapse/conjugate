// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"go.uber.org/zap"
)

// limitOperator limits the number of rows returned
type limitOperator struct {
	input Operator
	count int
	logger *zap.Logger

	ctx     context.Context
	returned int
	stats   *IteratorStats
	opened  bool
	closed  bool
}

// NewLimitOperator creates a new limit operator
func NewLimitOperator(input Operator, count int, logger *zap.Logger) *limitOperator {
	return &limitOperator{
		input:  input,
		count:  count,
		logger: logger,
		stats:  &IteratorStats{},
	}
}

// Open initializes the operator
func (l *limitOperator) Open(ctx context.Context) error {
	if l.opened {
		return nil
	}

	l.ctx = ctx
	l.logger.Debug("Opening limit operator", zap.Int("count", l.count))

	// Open input
	if err := l.input.Open(ctx); err != nil {
		return err
	}

	l.returned = 0
	l.opened = true

	return nil
}

// Next returns the next row up to the limit
func (l *limitOperator) Next(ctx context.Context) (*Row, error) {
	if l.closed {
		return nil, ErrClosed
	}

	if !l.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Check if we've reached the limit
	if l.returned >= l.count {
		return nil, ErrNoMoreRows
	}

	// Get next row from input
	row, err := l.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	l.returned++
	l.stats.RowsRead++
	l.stats.RowsReturned++

	return row, nil
}

// Close releases resources
func (l *limitOperator) Close() error {
	l.closed = true
	return l.input.Close()
}

// Stats returns execution statistics
func (l *limitOperator) Stats() *IteratorStats {
	return l.stats
}
