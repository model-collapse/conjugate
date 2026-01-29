// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/lookup"
	"go.uber.org/zap"
)

// lookupOperator enriches data with external lookup tables
type lookupOperator struct {
	input          Operator
	registry       *lookup.Registry
	tableName      string
	joinField      string
	joinFieldAlias string
	outputFields   []string
	outputAliases  []string
	logger         *zap.Logger

	lookupTable *lookup.LookupTable
	ctx         context.Context
	stats       *IteratorStats
	opened      bool
	closed      bool
}

// NewLookupOperator creates a new lookup operator
func NewLookupOperator(
	input Operator,
	registry *lookup.Registry,
	tableName string,
	joinField string,
	joinFieldAlias string,
	outputFields []string,
	outputAliases []string,
	logger *zap.Logger,
) (*lookupOperator, error) {
	// Get lookup table from registry
	lookupTable, err := registry.Get(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get lookup table: %w", err)
	}

	return &lookupOperator{
		input:          input,
		registry:       registry,
		tableName:      tableName,
		joinField:      joinField,
		joinFieldAlias: joinFieldAlias,
		outputFields:   outputFields,
		outputAliases:  outputAliases,
		lookupTable:    lookupTable,
		logger:         logger,
		stats:          &IteratorStats{},
	}, nil
}

// Open initializes the operator
func (l *lookupOperator) Open(ctx context.Context) error {
	if l.opened {
		return nil
	}

	l.ctx = ctx
	l.logger.Debug("Opening lookup operator",
		zap.String("table", l.tableName),
		zap.String("join_field", l.joinField),
		zap.Strings("output_fields", l.outputFields),
		zap.Int("table_size", l.lookupTable.Size()))

	if err := l.input.Open(ctx); err != nil {
		return err
	}

	l.opened = true
	return nil
}

// Next returns the next row with lookup fields enriched
func (l *lookupOperator) Next(ctx context.Context) (*Row, error) {
	if l.closed {
		return nil, ErrClosed
	}

	if !l.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	row, err := l.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	l.stats.RowsRead++

	// Get the join field value
	joinValue, exists := row.Get(l.joinField)
	if !exists {
		// Join field doesn't exist - log warning but don't fail
		l.logger.Debug("Join field not found in row",
			zap.String("field", l.joinField))
		l.stats.RowsReturned++
		return row, nil
	}

	// Convert to string for lookup key
	joinKey, ok := joinValue.(string)
	if !ok {
		// Try to convert to string
		joinKey = fmt.Sprintf("%v", joinValue)
	}

	// Perform lookup
	lookupRow, found := l.lookupTable.Lookup(joinKey)
	if !found {
		// No match in lookup table - log debug but don't fail
		l.logger.Debug("No lookup match found",
			zap.String("key", joinKey),
			zap.String("table", l.tableName))
		l.stats.RowsReturned++
		return row, nil
	}

	// Enrich row with lookup fields
	for i, outputField := range l.outputFields {
		// Get value from lookup table
		value, exists := lookupRow[outputField]
		if !exists {
			l.logger.Debug("Output field not found in lookup table",
				zap.String("field", outputField),
				zap.String("table", l.tableName))
			continue
		}

		// Determine the output field name (use alias if provided)
		outputName := outputField
		if i < len(l.outputAliases) && l.outputAliases[i] != "" {
			outputName = l.outputAliases[i]
		}

		// Add to row
		row.Set(outputName, value)
	}

	l.stats.RowsReturned++
	return row, nil
}

// Close releases resources
func (l *lookupOperator) Close() error {
	l.closed = true
	return l.input.Close()
}

// Stats returns execution statistics
func (l *lookupOperator) Stats() *IteratorStats {
	return l.stats
}
