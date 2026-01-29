// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// fillnullOperator fills NULL/missing values with a default value
type fillnullOperator struct {
	input      Operator
	value      ast.Expression   // Default value to fill
	fieldExprs []ast.Expression // Optional list of field expressions (empty = all fields)
	fieldNames []string         // Extracted field names from fieldExprs
	logger     *zap.Logger
	stats      *IteratorStats
	fieldSet   map[string]bool // Set of fields to fill (for quick lookup)
	fillValue  interface{}     // Evaluated fill value
	ctx        context.Context
	opened     bool
	closed     bool
}

// NewFillnullOperator creates a new fillnull operator
func NewFillnullOperator(input Operator, value ast.Expression, fields []ast.Expression, logger *zap.Logger) *fillnullOperator {
	// Extract field names from field expressions
	fieldNames := make([]string, 0, len(fields))
	for _, fieldExpr := range fields {
		if fieldRef, ok := fieldExpr.(*ast.FieldReference); ok {
			fieldNames = append(fieldNames, fieldRef.Name)
		}
	}

	// Convert fields slice to set for O(1) lookup
	fieldSet := make(map[string]bool)
	for _, field := range fieldNames {
		fieldSet[field] = true
	}

	return &fillnullOperator{
		input:      input,
		value:      value,
		fieldExprs: fields,
		fieldNames: fieldNames,
		fieldSet:   fieldSet,
		logger:     logger,
		stats:      &IteratorStats{},
	}
}

// Open initializes the operator
func (f *fillnullOperator) Open(ctx context.Context) error {
	if f.opened {
		return nil
	}

	f.ctx = ctx
	f.logger.Debug("Opening fillnull operator",
		zap.String("value", f.value.String()),
		zap.Strings("fields", f.fieldNames))

	// Evaluate the fill value (it should be a literal)
	if lit, ok := f.value.(*ast.Literal); ok {
		f.fillValue = lit.Value
	} else {
		return fmt.Errorf("fillnull value must be a literal, got %T", f.value)
	}

	// Open input operator
	if err := f.input.Open(ctx); err != nil {
		return err
	}

	f.opened = true
	return nil
}

// Next returns the next row with NULL values filled
func (f *fillnullOperator) Next(ctx context.Context) (*Row, error) {
	if f.closed {
		return nil, ErrClosed
	}

	if !f.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Get next row from input
	row, err := f.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	f.stats.RowsRead++

	// Fill NULL values in the row
	f.fillNullValues(row)

	f.stats.RowsReturned++
	return row, nil
}

// fillNullValues fills NULL values in the row according to the configuration
func (f *fillnullOperator) fillNullValues(row *Row) {
	// If no fields specified, fill all NULL fields
	if len(f.fieldNames) == 0 {
		// Get all field names from the row
		rowMap := row.ToMap()
		for fieldName, fieldValue := range rowMap {
			if fieldValue == nil {
				row.Set(fieldName, f.fillValue)
			}
		}
	} else {
		// Fill only specified fields
		for _, fieldName := range f.fieldNames {
			value, exists := row.Get(fieldName)
			if exists {
				if value == nil {
					row.Set(fieldName, f.fillValue)
				}
			} else {
				// Field doesn't exist, create it with fill value
				row.Set(fieldName, f.fillValue)
			}
		}
	}
}

// Close releases resources
func (f *fillnullOperator) Close() error {
	f.closed = true

	// Close input
	if f.input != nil {
		return f.input.Close()
	}

	return nil
}

// Stats returns execution statistics
func (f *fillnullOperator) Stats() *IteratorStats {
	return f.stats
}
