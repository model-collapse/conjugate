// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// addtotalsOperator adds row and/or column totals based on OpenSearch PPL specification
// - row=true (default): Adds a field to each row with the row-wise total
// - col=true: Adds a summary row at the end with column totals
type addtotalsOperator struct {
	input      Operator
	fields     []ast.Expression // Specific fields to total (if empty, total all numeric fields)
	row        bool             // If true, add row totals as new field (default: true)
	col        bool             // If true, add column totals in summary row (default: false)
	labelField string           // Field to use for the "Total" label in summary row
	label      string           // Label for the summary row (default: "Total")
	fieldName  string           // Name of field for row totals (used when row=true)
	logger     *zap.Logger
	stats      *IteratorStats

	// State
	ctx         context.Context
	buffer      []*Row         // Buffer to hold all rows (only used when col=true)
	index       int            // Current index for emitting rows
	totalsAdded bool           // True when totals row has been calculated and added
	opened      bool
	closed      bool
}

// NewAddtotalsOperator creates a new addtotals operator
func NewAddtotalsOperator(
	input Operator,
	fields []ast.Expression,
	row bool,            // Add row totals (default: true)
	col bool,            // Add column totals (default: false)
	labelField string,   // Field for summary row label
	label string,        // Text for summary row label (default: "Total")
	fieldName string,    // Field name for row totals
	logger *zap.Logger,
) *addtotalsOperator {
	if label == "" {
		label = "Total"
	}
	if fieldName == "" {
		fieldName = "total"
	}
	return &addtotalsOperator{
		input:      input,
		fields:     fields,
		row:        row,
		col:        col,
		labelField: labelField,
		label:      label,
		fieldName:  fieldName,
		logger:     logger,
		stats:      &IteratorStats{},
	}
}

// Open initializes the operator
func (a *addtotalsOperator) Open(ctx context.Context) error {
	if a.opened {
		return nil
	}

	a.ctx = ctx
	a.logger.Debug("Opening addtotals operator")

	// Open input operator
	if err := a.input.Open(ctx); err != nil {
		return err
	}

	a.opened = true
	return nil
}

// Next returns the next row with modifications based on row/col settings
func (a *addtotalsOperator) Next(ctx context.Context) (*Row, error) {
	if a.closed {
		return nil, ErrClosed
	}

	if !a.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Case 1: col=true - Need to buffer all rows to calculate column totals
	if a.col {
		return a.nextWithColumnTotals(ctx)
	}

	// Case 2: col=false, row=true - Streaming mode, just add row totals
	if a.row {
		return a.nextWithRowTotals(ctx)
	}

	// Case 3: col=false, row=false - Pass through (no modifications)
	row, err := a.input.Next(ctx)
	if err != nil {
		return nil, err
	}
	a.stats.RowsRead++
	a.stats.RowsReturned++
	return row, nil
}

// nextWithRowTotals adds row totals in streaming mode (O(1) memory)
func (a *addtotalsOperator) nextWithRowTotals(ctx context.Context) (*Row, error) {
	row, err := a.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	a.stats.RowsRead++

	// Calculate row total
	total := 0.0
	fieldsToTotal := a.getFieldsToTotal(row)

	for field := range fieldsToTotal {
		value, exists := row.Get(field)
		if !exists {
			continue
		}

		floatVal, ok := toFloat64(value)
		if ok {
			total += floatVal
		}
	}

	// Add row total field
	row.Set(a.fieldName, total)

	a.stats.RowsReturned++
	return row, nil
}

// nextWithColumnTotals buffers all rows and adds column totals summary row
func (a *addtotalsOperator) nextWithColumnTotals(ctx context.Context) (*Row, error) {
	// First call: buffer all rows and calculate column totals
	if !a.totalsAdded {
		a.logger.Debug("Buffering all rows to calculate column totals")

		a.buffer = make([]*Row, 0, 1000)

		// Read all rows from input
		for {
			row, err := a.input.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			if err != nil {
				return nil, err
			}

			a.stats.RowsRead++

			// If row=true, also add row totals to each buffered row
			if a.row {
				total := 0.0
				fieldsToTotal := a.getFieldsToTotal(row)

				for field := range fieldsToTotal {
					value, exists := row.Get(field)
					if !exists {
						continue
					}

					floatVal, ok := toFloat64(value)
					if ok {
						total += floatVal
					}
				}

				row.Set(a.fieldName, total)
			}

			a.buffer = append(a.buffer, row)
		}

		a.logger.Debug("Buffered all rows",
			zap.Int("total_rows", len(a.buffer)))

		// Calculate and append column totals row
		if len(a.buffer) > 0 {
			totalsRow := a.calculateColumnTotals()
			a.buffer = append(a.buffer, totalsRow)
			a.logger.Debug("Added column totals row")
		}

		a.totalsAdded = true
		a.index = 0
	}

	// Emit rows including the totals row
	if a.index < len(a.buffer) {
		row := a.buffer[a.index]
		a.index++
		a.stats.RowsReturned++
		return row, nil
	}

	// No more rows
	return nil, ErrNoMoreRows
}

// getFieldsToTotal determines which fields to include in totals
func (a *addtotalsOperator) getFieldsToTotal(sampleRow *Row) map[string]bool {
	fieldsToTotal := make(map[string]bool)

	if len(a.fields) > 0 {
		// Specific fields requested
		for _, fieldExpr := range a.fields {
			if fieldRef, ok := fieldExpr.(*ast.FieldReference); ok {
				fieldsToTotal[fieldRef.Name] = true
			}
		}
	} else {
		// Total all numeric fields from sample row
		for _, field := range sampleRow.Fields() {
			// When calculating row totals (row=true, col=false), skip the row total field to avoid recursion
			// When calculating column totals (col=true), include all fields including row totals if present
			if field == a.fieldName && a.row && !a.col {
				continue
			}
			fieldsToTotal[field] = true
		}
	}

	return fieldsToTotal
}

// calculateColumnTotals calculates column totals for the summary row
func (a *addtotalsOperator) calculateColumnTotals() *Row {
	if len(a.buffer) == 0 {
		return NewRow(make(map[string]interface{}))
	}

	totalsData := make(map[string]interface{})

	// Determine which fields to total
	fieldsToTotal := a.getFieldsToTotal(a.buffer[0])

	// Initialize totals map
	totals := make(map[string]float64)
	isNumeric := make(map[string]bool)

	// Accumulate totals
	for _, row := range a.buffer {
		for field := range fieldsToTotal {
			value, exists := row.Get(field)
			if !exists {
				continue
			}

			// Try to convert to float64
			floatVal, ok := toFloat64(value)
			if ok {
				totals[field] += floatVal
				isNumeric[field] = true
			}
		}
	}

	// Build totals row
	// Add totals for numeric fields
	for field, total := range totals {
		if isNumeric[field] {
			totalsData[field] = total
		}
	}

	// Add label field
	if a.labelField != "" {
		// Use specified label field
		totalsData[a.labelField] = a.label
	} else if a.fieldName != "" {
		// Use specified field name
		totalsData[a.fieldName] = a.label
	} else {
		// Find first non-numeric field from first row (likely a group-by field)
		// and use that for the label
		firstRow := a.buffer[0]
		for _, field := range firstRow.Fields() {
			if !isNumeric[field] {
				totalsData[field] = a.label
				break
			}
		}
	}

	// If no label was added, add a default one
	if len(totalsData) == len(totals) {
		// No label was added
		totalsData["_total"] = a.label
	}

	return NewRow(totalsData)
}

// Close releases resources
func (a *addtotalsOperator) Close() error {
	a.closed = true

	// Clear buffer to free memory
	a.buffer = nil

	// Close input
	if a.input != nil {
		return a.input.Close()
	}

	return nil
}

// Stats returns execution statistics
func (a *addtotalsOperator) Stats() *IteratorStats {
	return a.stats
}
