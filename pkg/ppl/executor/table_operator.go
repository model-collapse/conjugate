// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// tableOperator selects specific fields for display
// Similar to project operator but always includes (never excludes)
type tableOperator struct {
	input  Operator
	fields []ast.Expression
	logger *zap.Logger

	ctx    context.Context
	stats  *IteratorStats
	opened bool
	closed bool
}

// NewTableOperator creates a new table operator
func NewTableOperator(input Operator, fields []ast.Expression, logger *zap.Logger) *tableOperator {
	return &tableOperator{
		input:  input,
		fields: fields,
		logger: logger,
		stats:  &IteratorStats{},
	}
}

// Open initializes the operator
func (t *tableOperator) Open(ctx context.Context) error {
	if t.opened {
		return nil
	}

	t.ctx = ctx
	t.logger.Debug("Opening table operator",
		zap.Int("num_fields", len(t.fields)))

	if err := t.input.Open(ctx); err != nil {
		return err
	}

	t.opened = true
	return nil
}

// Next returns the next row with only selected fields
func (t *tableOperator) Next(ctx context.Context) (*Row, error) {
	if t.closed {
		return nil, ErrClosed
	}

	if !t.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	row, err := t.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	t.stats.RowsRead++

	// Create new row with only selected fields
	newFields := make(map[string]interface{})

	for _, fieldExpr := range t.fields {
		// For now, only handle simple field references
		// More complex expressions would need evaluation
		switch expr := fieldExpr.(type) {
		case *ast.FieldReference:
			// Get field value from input row
			if val, exists := row.Get(expr.Name); exists {
				newFields[expr.Name] = val
			} else {
				// Field doesn't exist, set to nil
				newFields[expr.Name] = nil
			}

		case *ast.FunctionCall:
			// For function calls, use the function name or alias as field name
			// The function should already be evaluated in the input row
			fieldName := expr.Name
			if val, exists := row.Get(fieldName); exists {
				newFields[fieldName] = val
			}

		default:
			// For other expressions, use the string representation as field name
			fieldName := expr.String()
			if val, exists := row.Get(fieldName); exists {
				newFields[fieldName] = val
			}
		}
	}

	t.stats.RowsReturned++
	return NewRow(newFields), nil
}

// Close releases resources
func (t *tableOperator) Close() error {
	t.closed = true
	return t.input.Close()
}

// Stats returns execution statistics
func (t *tableOperator) Stats() *IteratorStats {
	return t.stats
}
