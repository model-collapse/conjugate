// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// renameOperator renames fields in rows
type renameOperator struct {
	input       Operator
	assignments []*ast.RenameAssignment
	logger      *zap.Logger

	ctx    context.Context
	stats  *IteratorStats
	opened bool
	closed bool
}

// NewRenameOperator creates a new rename operator
func NewRenameOperator(
	input Operator,
	assignments []*ast.RenameAssignment,
	logger *zap.Logger,
) *renameOperator {
	return &renameOperator{
		input:       input,
		assignments: assignments,
		logger:      logger,
		stats:       &IteratorStats{},
	}
}

// Open initializes the operator
func (r *renameOperator) Open(ctx context.Context) error {
	if r.opened {
		return nil
	}

	r.ctx = ctx
	r.logger.Debug("Opening rename operator",
		zap.Int("num_assignments", len(r.assignments)))

	if err := r.input.Open(ctx); err != nil {
		return err
	}

	r.opened = true
	return nil
}

// Next returns the next row with renamed fields
func (r *renameOperator) Next(ctx context.Context) (*Row, error) {
	if r.closed {
		return nil, ErrClosed
	}

	if !r.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	row, err := r.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	r.stats.RowsRead++

	// Apply renames
	for _, assignment := range r.assignments {
		// Get value from old field
		val, ok := row.Get(assignment.OldName)
		if ok {
			// Set under new name
			row.Set(assignment.NewName, val)
			// Remove old name
			row.Delete(assignment.OldName)
		}
	}

	r.stats.RowsReturned++
	return row, nil
}

// Close releases resources
func (r *renameOperator) Close() error {
	r.closed = true
	return r.input.Close()
}

// Stats returns execution statistics
func (r *renameOperator) Stats() *IteratorStats {
	return r.stats
}
