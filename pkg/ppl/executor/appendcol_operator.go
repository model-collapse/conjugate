// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// appendcolOperator merges columns from a subsearch horizontally into the main result
// Unlike append (which concatenates rows vertically), appendcol adds columns side-by-side
type appendcolOperator struct {
	input          Operator
	subsearch      Operator
	logger         *zap.Logger
	ctx            context.Context
	stats          *IteratorStats
	opened         bool
	closed         bool
	subsearchRows  []*Row  // Buffer subsearch results
	subsearchIndex int     // Current position in subsearch results
	override       bool    // If true, subsearch columns override main columns
}

// NewAppendcolOperator creates a new appendcol operator
// If override is true, columns from subsearch will override columns from main input with same name
func NewAppendcolOperator(input Operator, subsearch Operator, override bool, logger *zap.Logger) *appendcolOperator {
	return &appendcolOperator{
		input:     input,
		subsearch: subsearch,
		override:  override,
		logger:    logger,
		stats:     &IteratorStats{},
	}
}

// Open initializes the operator and buffers subsearch results
func (a *appendcolOperator) Open(ctx context.Context) error {
	if a.opened {
		return nil
	}

	a.ctx = ctx
	a.logger.Debug("Opening appendcol operator")

	// Open input
	if err := a.input.Open(ctx); err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}

	// Open subsearch and buffer all results
	if err := a.subsearch.Open(ctx); err != nil {
		return fmt.Errorf("failed to open subsearch: %w", err)
	}

	// Buffer all subsearch results
	a.subsearchRows = make([]*Row, 0)
	for {
		row, err := a.subsearch.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read subsearch: %w", err)
		}
		a.subsearchRows = append(a.subsearchRows, row)
	}

	a.logger.Debug("Buffered subsearch results",
		zap.Int("count", len(a.subsearchRows)))

	a.opened = true
	a.subsearchIndex = 0
	return nil
}

// Next returns the next row with columns merged from both input and subsearch
func (a *appendcolOperator) Next(ctx context.Context) (*Row, error) {
	if a.closed {
		return nil, ErrClosed
	}

	if !a.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Read next row from main input
	mainRow, err := a.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	a.stats.RowsRead++

	// If no subsearch results, return main row as-is
	if len(a.subsearchRows) == 0 {
		a.stats.RowsReturned++
		return mainRow, nil
	}

	// Merge with corresponding subsearch row (row-by-row alignment)
	// If we've exhausted subsearch rows, the remaining main rows get no additional columns
	var mergedRow *Row
	if a.subsearchIndex < len(a.subsearchRows) {
		subsearchRow := a.subsearchRows[a.subsearchIndex]
		mergedRow = a.mergeRows(mainRow, subsearchRow)
		a.subsearchIndex++
	} else {
		// No more subsearch rows, return main row only
		mergedRow = mainRow
	}

	a.stats.RowsReturned++
	return mergedRow, nil
}

// mergeRows combines columns from mainRow and subsearchRow
func (a *appendcolOperator) mergeRows(mainRow, subsearchRow *Row) *Row {
	// Start with a clone of the main row
	merged := mainRow.Clone()

	// Add/override with subsearch fields
	for _, fieldName := range subsearchRow.Fields() {
		value, _ := subsearchRow.Get(fieldName)

		// Check for conflicts
		if _, exists := merged.Get(fieldName); exists {
			if a.override {
				// Override: subsearch wins
				merged.Set(fieldName, value)
				a.logger.Debug("Column conflict resolved by override",
					zap.String("field", fieldName))
			} else {
				// Default: main wins, log warning
				a.logger.Debug("Column conflict: keeping main value",
					zap.String("field", fieldName))
			}
		} else {
			// No conflict, add field
			merged.Set(fieldName, value)
		}
	}

	return merged
}

// Close releases resources
func (a *appendcolOperator) Close() error {
	a.closed = true

	// Close both input and subsearch
	var inputErr, subsearchErr error
	if a.input != nil {
		inputErr = a.input.Close()
	}
	if a.subsearch != nil {
		subsearchErr = a.subsearch.Close()
	}

	// Clear buffer
	a.subsearchRows = nil

	if inputErr != nil {
		return inputErr
	}
	return subsearchErr
}

// Stats returns execution statistics
func (a *appendcolOperator) Stats() *IteratorStats {
	return a.stats
}
