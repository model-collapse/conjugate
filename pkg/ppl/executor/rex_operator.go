// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"
	"regexp"

	"go.uber.org/zap"
)

// rexOperator extracts fields using regular expressions
type rexOperator struct {
	input           Operator
	sourceField     string // Can be empty (defaults to _raw)
	pattern         *regexp.Regexp
	extractedFields []string
	logger          *zap.Logger

	ctx    context.Context
	stats  *IteratorStats
	opened bool
	closed bool
}

// NewRexOperator creates a new rex operator
func NewRexOperator(
	input Operator,
	sourceField string, // If empty, defaults to _raw
	patternStr string,
	extractedFields []string,
	logger *zap.Logger,
) (*rexOperator, error) {
	// Compile the regex pattern
	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	// Default to _raw if sourceField is empty
	if sourceField == "" {
		sourceField = "_raw"
	}

	return &rexOperator{
		input:           input,
		sourceField:     sourceField,
		pattern:         pattern,
		extractedFields: extractedFields,
		logger:          logger,
		stats:           &IteratorStats{},
	}, nil
}

// Open initializes the operator
func (r *rexOperator) Open(ctx context.Context) error {
	if r.opened {
		return nil
	}

	r.ctx = ctx
	r.logger.Debug("Opening rex operator",
		zap.String("source_field", r.sourceField),
		zap.String("pattern", r.pattern.String()),
		zap.Strings("extracted_fields", r.extractedFields))

	if err := r.input.Open(ctx); err != nil {
		return err
	}

	r.opened = true
	return nil
}

// Next returns the next row with extracted fields added
func (r *rexOperator) Next(ctx context.Context) (*Row, error) {
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

	// Get the source field value
	sourceValue, exists := row.Get(r.sourceField)
	if !exists {
		// Source field doesn't exist - log debug but don't fail
		// Rex is lenient and just continues without extraction
		r.logger.Debug("Source field not found in row",
			zap.String("field", r.sourceField))
		r.stats.RowsReturned++
		return row, nil
	}

	// Convert to string
	sourceStr, ok := sourceValue.(string)
	if !ok {
		// Try to convert to string
		sourceStr = fmt.Sprintf("%v", sourceValue)
	}

	// Apply regex pattern
	matches := r.pattern.FindStringSubmatch(sourceStr)
	if matches == nil {
		// No match - log debug but don't fail
		r.logger.Debug("Pattern did not match",
			zap.String("field", r.sourceField),
			zap.String("value", sourceStr))
		r.stats.RowsReturned++
		return row, nil
	}

	// Extract named capture groups
	// matches[0] is the full match, matches[1:] are the capture groups
	subexpNames := r.pattern.SubexpNames()
	for i := 1; i < len(matches) && i < len(subexpNames); i++ {
		if subexpNames[i] != "" && i < len(matches) {
			// Add the extracted field to the row
			row.Set(subexpNames[i], matches[i])
		}
	}

	r.stats.RowsReturned++
	return row, nil
}

// Close releases resources
func (r *rexOperator) Close() error {
	r.closed = true
	return r.input.Close()
}

// Stats returns execution statistics
func (r *rexOperator) Stats() *IteratorStats {
	return r.stats
}
