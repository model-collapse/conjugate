// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// appendpipeOperator processes current results through a subsearch pipeline
// and appends those results as additional rows.
// Unlike append (which runs subsearch independently), appendpipe passes
// current results as input to the subsearch.
//
// Example:
//   source=errors | stats count() by error_code
//   | appendpipe [stats sum(count) as total]
//
// This would add a summary row with the total count.
type appendpipeOperator struct {
	input           Operator
	subsearchPlan   Operator // The subsearch pipeline (not yet opened)
	logger          *zap.Logger
	ctx             context.Context
	stats           *IteratorStats
	opened          bool
	closed          bool
	bufferedRows    []*Row     // Buffer all input rows
	currentIndex    int        // Current position in buffer
	subsearchOpened bool       // Whether subsearch has been opened
	subsearchDone   bool       // Whether we've finished reading subsearch
	readingSubsearch bool      // Whether we're currently reading subsearch
}

// NewAppendpipeOperator creates a new appendpipe operator
// The subsearchPlan should be a pipeline that will process the buffered rows
func NewAppendpipeOperator(input Operator, subsearchPlan Operator, logger *zap.Logger) *appendpipeOperator {
	return &appendpipeOperator{
		input:         input,
		subsearchPlan: subsearchPlan,
		logger:        logger,
		stats:         &IteratorStats{},
	}
}

// Open initializes the operator and buffers all input rows
func (a *appendpipeOperator) Open(ctx context.Context) error {
	if a.opened {
		return nil
	}

	a.ctx = ctx
	a.logger.Debug("Opening appendpipe operator")

	// Open input
	if err := a.input.Open(ctx); err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}

	// Buffer all input rows
	a.bufferedRows = make([]*Row, 0)
	for {
		row, err := a.input.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		a.bufferedRows = append(a.bufferedRows, row)
	}

	a.logger.Debug("Buffered input rows for appendpipe",
		zap.Int("count", len(a.bufferedRows)))

	a.opened = true
	a.currentIndex = 0
	return nil
}

// Next returns rows: first all buffered input rows, then all subsearch results
func (a *appendpipeOperator) Next(ctx context.Context) (*Row, error) {
	if a.closed {
		return nil, ErrClosed
	}

	if !a.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Phase 1: Return buffered rows from original input
	if !a.readingSubsearch {
		if a.currentIndex < len(a.bufferedRows) {
			row := a.bufferedRows[a.currentIndex]
			a.currentIndex++
			a.stats.RowsReturned++
			return row, nil
		}

		// Original rows exhausted, switch to subsearch
		a.readingSubsearch = true

		// Open and execute subsearch on the buffered rows
		if !a.subsearchOpened {
			// Create a slice iterator from buffered rows to feed to subsearch
			subsearchInput := NewSliceIterator(a.bufferedRows)

			// We need to wire the subsearch to use our buffered rows as input
			// This is tricky - the subsearch plan needs to be connected to our buffer
			// For now, we'll open the subsearch normally (it should handle its own input)
			if err := a.subsearchPlan.Open(ctx); err != nil {
				return nil, fmt.Errorf("failed to open subsearch: %w", err)
			}
			a.subsearchOpened = true

			// Note: In a real implementation, we'd need to inject the buffered rows
			// as the input to the subsearch pipeline. This might require special
			// handling in the query builder/planner.
			_ = subsearchInput // For now, this is a placeholder
		}
	}

	// Phase 2: Return rows from subsearch
	if a.readingSubsearch && !a.subsearchDone {
		row, err := a.subsearchPlan.Next(ctx)
		if err == ErrNoMoreRows {
			a.subsearchDone = true
			return nil, ErrNoMoreRows
		}
		if err != nil {
			return nil, fmt.Errorf("subsearch failed: %w", err)
		}

		a.stats.RowsReturned++
		return row, nil
	}

	return nil, ErrNoMoreRows
}

// Close releases resources
func (a *appendpipeOperator) Close() error {
	a.closed = true

	var inputErr, subsearchErr error
	if a.input != nil {
		inputErr = a.input.Close()
	}
	if a.subsearchPlan != nil && a.subsearchOpened {
		subsearchErr = a.subsearchPlan.Close()
	}

	// Clear buffer
	a.bufferedRows = nil

	if inputErr != nil {
		return inputErr
	}
	return subsearchErr
}

// Stats returns execution statistics
func (a *appendpipeOperator) Stats() *IteratorStats {
	return a.stats
}
