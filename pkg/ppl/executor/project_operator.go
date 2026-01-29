// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// projectOperator selects or excludes fields from rows
type projectOperator struct {
	input   Operator
	fields  []ast.Expression
	exclude bool
	logger  *zap.Logger

	ctx    context.Context
	stats  *IteratorStats
	opened bool
	closed bool

	// Cached field names for simple projections
	fieldNames []string
}

// NewProjectOperator creates a new project operator
func NewProjectOperator(input Operator, fields []ast.Expression, exclude bool, logger *zap.Logger) *projectOperator {
	op := &projectOperator{
		input:   input,
		fields:  fields,
		exclude: exclude,
		logger:  logger,
		stats:   &IteratorStats{},
	}

	// Pre-compute simple field names
	op.fieldNames = make([]string, 0, len(fields))
	for _, f := range fields {
		if ref, ok := f.(*ast.FieldReference); ok {
			op.fieldNames = append(op.fieldNames, ref.Name)
		}
	}

	return op
}

// Open initializes the operator
func (p *projectOperator) Open(ctx context.Context) error {
	if p.opened {
		return nil
	}

	p.ctx = ctx
	p.logger.Debug("Opening project operator",
		zap.Int("num_fields", len(p.fields)),
		zap.Bool("exclude", p.exclude))

	// Open input
	if err := p.input.Open(ctx); err != nil {
		return err
	}

	p.opened = true
	return nil
}

// Next returns the next row with projected fields
func (p *projectOperator) Next(ctx context.Context) (*Row, error) {
	if p.closed {
		return nil, ErrClosed
	}

	if !p.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Get next row from input
	row, err := p.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	p.stats.RowsRead++
	p.stats.RowsReturned++

	// Apply projection
	return p.projectRow(row), nil
}

// projectRow applies the projection to a row
func (p *projectOperator) projectRow(row *Row) *Row {
	if p.exclude {
		// Exclude specified fields
		result := row.Clone()
		for _, name := range p.fieldNames {
			result.Delete(name)
		}
		return result
	}

	// Include only specified fields
	result := NewRow(nil)
	for _, name := range p.fieldNames {
		if val, ok := row.Get(name); ok {
			result.Set(name, val)
		}
	}

	// Also preserve _id and _score if not explicitly excluded
	if _, ok := row.Get("_id"); ok {
		result.Set("_id", func() interface{} { v, _ := row.Get("_id"); return v }())
	}
	if _, ok := row.Get("_score"); ok {
		result.Set("_score", func() interface{} { v, _ := row.Get("_score"); return v }())
	}

	return result
}

// Close releases resources
func (p *projectOperator) Close() error {
	p.closed = true
	return p.input.Close()
}

// Stats returns execution statistics
func (p *projectOperator) Stats() *IteratorStats {
	return p.stats
}
