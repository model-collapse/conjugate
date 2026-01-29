// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"go.uber.org/zap"
)

// appendOperator concatenates results from a subsearch to the main search
type appendOperator struct {
	input       Operator
	subsearch   Operator
	logger      *zap.Logger
	ctx         context.Context
	stats       *IteratorStats
	opened      bool
	closed      bool
	inputDone   bool // Track when input is exhausted
	currentIter Operator
}

// NewAppendOperator creates a new append operator
func NewAppendOperator(input Operator, subsearch Operator, logger *zap.Logger) *appendOperator {
	return &appendOperator{
		input:       input,
		subsearch:   subsearch,
		logger:      logger,
		stats:       &IteratorStats{},
		currentIter: input, // Start with main input
	}
}

// Open initializes the operator
func (a *appendOperator) Open(ctx context.Context) error {
	if a.opened {
		return nil
	}

	a.ctx = ctx
	a.logger.Debug("Opening append operator")

	// Open input
	if err := a.input.Open(ctx); err != nil {
		return err
	}

	// Don't open subsearch yet - we'll open it when input is exhausted
	a.opened = true
	return nil
}

// Next returns the next row, first from input, then from subsearch
func (a *appendOperator) Next(ctx context.Context) (*Row, error) {
	if a.closed {
		return nil, ErrClosed
	}

	if !a.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// If we're still reading from input
	if !a.inputDone {
		row, err := a.input.Next(ctx)
		if err == ErrNoMoreRows {
			// Input exhausted, switch to subsearch
			a.inputDone = true
			a.logger.Debug("Main input exhausted, switching to subsearch")

			// Open subsearch now
			if err := a.subsearch.Open(ctx); err != nil {
				return nil, err
			}
			a.currentIter = a.subsearch

			// Read first row from subsearch
			return a.subsearch.Next(ctx)
		}
		if err != nil {
			return nil, err
		}

		a.stats.RowsRead++
		a.stats.RowsReturned++
		return row, nil
	}

	// Reading from subsearch
	row, err := a.subsearch.Next(ctx)
	if err != nil {
		return nil, err
	}

	a.stats.RowsRead++
	a.stats.RowsReturned++
	return row, nil
}

// Close releases resources
func (a *appendOperator) Close() error {
	a.closed = true

	// Close both input and subsearch
	var inputErr, subsearchErr error
	if a.input != nil {
		inputErr = a.input.Close()
	}
	if a.subsearch != nil && a.inputDone {
		// Only close subsearch if it was opened
		subsearchErr = a.subsearch.Close()
	}

	if inputErr != nil {
		return inputErr
	}
	return subsearchErr
}

// Stats returns execution statistics
func (a *appendOperator) Stats() *IteratorStats {
	return a.stats
}
