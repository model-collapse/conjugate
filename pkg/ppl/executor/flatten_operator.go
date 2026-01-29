// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// flattenOperator flattens nested arrays/objects into separate rows
// It reads one row at a time and emits multiple rows if the field contains an array
type flattenOperator struct {
	input  Operator
	field  ast.Expression
	logger *zap.Logger
	stats  *IteratorStats

	// State for buffering expanded rows from current input row
	ctx           context.Context
	expandedRows  []*Row // Buffered rows from flattening current input row
	expandedIndex int    // Current index in expandedRows
	opened        bool
	closed        bool
}

// NewFlattenOperator creates a new flatten operator
func NewFlattenOperator(input Operator, field ast.Expression, logger *zap.Logger) *flattenOperator {
	return &flattenOperator{
		input:  input,
		field:  field,
		logger: logger,
		stats:  &IteratorStats{},
	}
}

// Open initializes the operator
func (f *flattenOperator) Open(ctx context.Context) error {
	if f.opened {
		return nil
	}

	f.ctx = ctx
	f.logger.Debug("Opening flatten operator",
		zap.String("field", f.field.String()))

	if err := f.input.Open(ctx); err != nil {
		return err
	}

	f.opened = true
	f.expandedIndex = 0
	return nil
}

// Next returns the next flattened row
func (f *flattenOperator) Next(ctx context.Context) (*Row, error) {
	if f.closed {
		return nil, ErrClosed
	}

	if !f.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// If we have buffered expanded rows, return the next one
	if f.expandedIndex < len(f.expandedRows) {
		row := f.expandedRows[f.expandedIndex]
		f.expandedIndex++
		f.stats.RowsReturned++
		return row, nil
	}

	// Need to read next input row and expand it
	inputRow, err := f.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	f.stats.RowsRead++

	// Extract the field value to flatten
	fieldName := f.getFieldName()
	fieldValue, exists := inputRow.Get(fieldName)

	// If field doesn't exist or is nil, return the row as-is
	if !exists || fieldValue == nil {
		f.stats.RowsReturned++
		return inputRow, nil
	}

	// Try to flatten the field value
	expanded := f.flattenValue(fieldValue, fieldName, inputRow)

	// If no expansion occurred (not an array), return original row
	if len(expanded) == 0 {
		f.stats.RowsReturned++
		return inputRow, nil
	}

	// Buffer the expanded rows
	f.expandedRows = expanded
	f.expandedIndex = 1 // Start from index 1 (return first row now)

	f.stats.RowsReturned++
	return f.expandedRows[0], nil
}

// getFieldName extracts the field name from the expression
func (f *flattenOperator) getFieldName() string {
	switch expr := f.field.(type) {
	case *ast.FieldReference:
		return expr.Name
	case *ast.FunctionCall:
		return expr.Name
	default:
		return f.field.String()
	}
}

// flattenValue attempts to flatten a field value
// Returns a slice of expanded rows, or empty slice if not flattenable
func (f *flattenOperator) flattenValue(value interface{}, fieldName string, originalRow *Row) []*Row {
	// Try to flatten as array
	if arrayVal, ok := value.([]interface{}); ok {
		return f.flattenArray(arrayVal, fieldName, originalRow)
	}

	// Try to flatten as slice of maps (common JSON structure)
	if arrayVal, ok := value.([]map[string]interface{}); ok {
		interfaceArray := make([]interface{}, len(arrayVal))
		for i, v := range arrayVal {
			interfaceArray[i] = v
		}
		return f.flattenArray(interfaceArray, fieldName, originalRow)
	}

	// Not flattenable - return empty slice
	return nil
}

// flattenArray flattens an array value into multiple rows
func (f *flattenOperator) flattenArray(array []interface{}, fieldName string, originalRow *Row) []*Row {
	if len(array) == 0 {
		// Empty array - return original row with field set to nil
		newRow := f.cloneRow(originalRow)
		newRow.Set(fieldName, nil)
		return []*Row{newRow}
	}

	// Create one row for each array element
	result := make([]*Row, len(array))
	for i, element := range array {
		newRow := f.cloneRow(originalRow)
		newRow.Set(fieldName, element)
		result[i] = newRow
	}

	f.logger.Debug("Flattened array",
		zap.String("field", fieldName),
		zap.Int("elements", len(array)))

	return result
}

// cloneRow creates a copy of a row
func (f *flattenOperator) cloneRow(row *Row) *Row {
	rowMap := row.ToMap()
	newMap := make(map[string]interface{}, len(rowMap))
	for k, v := range rowMap {
		newMap[k] = v
	}
	return NewRow(newMap)
}

// Close releases resources
func (f *flattenOperator) Close() error {
	f.closed = true

	// Clear buffered rows
	f.expandedRows = nil

	if f.input != nil {
		return f.input.Close()
	}

	return nil
}

// Stats returns execution statistics
func (f *flattenOperator) Stats() *IteratorStats {
	return f.stats
}
