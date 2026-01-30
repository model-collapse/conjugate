// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// addcoltotalsOperator adds a column with the row-wise total of numeric fields
// For each row, it adds a new field (default: "Total") containing the sum of all numeric values in that row
type addcoltotalsOperator struct {
	input      Operator
	fields     []ast.Expression // Specific fields to include in total (if empty, all numeric fields)
	labelField string           // Name of the total column (default: "Total")
	label      string           // Not used for addcoltotals (kept for API compatibility)
	logger     *zap.Logger
	stats      *IteratorStats

	// State
	ctx    context.Context
	opened bool
	closed bool
}

// NewAddcoltotalsOperator creates a new addcoltotals operator
func NewAddcoltotalsOperator(
	input Operator,
	fields []ast.Expression,
	labelField string,
	label string,
	logger *zap.Logger,
) *addcoltotalsOperator {
	if labelField == "" {
		labelField = "Total"
	}
	return &addcoltotalsOperator{
		input:      input,
		fields:     fields,
		labelField: labelField,
		label:      label,
		logger:     logger,
		stats:      &IteratorStats{},
	}
}

// Open initializes the operator
func (a *addcoltotalsOperator) Open(ctx context.Context) error {
	if a.opened {
		return nil
	}

	a.ctx = ctx
	a.logger.Debug("Opening addcoltotals operator")

	// Open input operator
	if err := a.input.Open(ctx); err != nil {
		return err
	}

	a.opened = true
	return nil
}

// Next returns the next row with an additional total column
func (a *addcoltotalsOperator) Next(ctx context.Context) (*Row, error) {
	if a.closed {
		return nil, ErrClosed
	}

	if !a.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Get next row from input
	row, err := a.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	a.stats.RowsRead++

	// Add total column to this row
	rowWithTotal := a.addTotalColumn(row)

	a.stats.RowsReturned++
	return rowWithTotal, nil
}

// addTotalColumn adds a column with the row-wise sum of numeric fields
func (a *addcoltotalsOperator) addTotalColumn(row *Row) *Row {
	// Clone the row to avoid modifying the original
	newRow := row.Clone()

	// Determine which fields to sum
	fieldsToSum := make(map[string]bool)
	if len(a.fields) > 0 {
		// Specific fields requested
		for _, fieldExpr := range a.fields {
			if fieldRef, ok := fieldExpr.(*ast.FieldReference); ok {
				fieldsToSum[fieldRef.Name] = true
			}
		}
	} else {
		// Sum all numeric fields
		for _, field := range row.Fields() {
			fieldsToSum[field] = true
		}
	}

	// Calculate the sum of numeric fields
	var total float64
	for field := range fieldsToSum {
		value, exists := row.Get(field)
		if !exists {
			continue
		}

		// Try to convert to float64 and add to total
		floatVal, ok := toFloat64(value)
		if ok {
			total += floatVal
		}
	}

	// Add the total column
	newRow.Set(a.labelField, total)

	return newRow
}

// Close releases resources
func (a *addcoltotalsOperator) Close() error {
	a.closed = true

	// Close input
	if a.input != nil {
		return a.input.Close()
	}

	return nil
}

// Stats returns execution statistics
func (a *addcoltotalsOperator) Stats() *IteratorStats {
	return a.stats
}
