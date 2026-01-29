// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"
	"regexp"

	"go.uber.org/zap"
)

// parseOperator extracts fields from text using regex patterns
type parseOperator struct {
	input           Operator
	sourceField     string
	pattern         *regexp.Regexp
	extractedFields []string
	logger          *zap.Logger

	ctx    context.Context
	stats  *IteratorStats
	opened bool
	closed bool
}

// NewParseOperator creates a new parse operator
func NewParseOperator(
	input Operator,
	sourceField string,
	patternStr string,
	extractedFields []string,
	logger *zap.Logger,
) (*parseOperator, error) {
	// Compile the regex pattern
	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return &parseOperator{
		input:           input,
		sourceField:     sourceField,
		pattern:         pattern,
		extractedFields: extractedFields,
		logger:          logger,
		stats:           &IteratorStats{},
	}, nil
}

// Open initializes the operator
func (p *parseOperator) Open(ctx context.Context) error {
	if p.opened {
		return nil
	}

	p.ctx = ctx
	p.logger.Debug("Opening parse operator",
		zap.String("source_field", p.sourceField),
		zap.String("pattern", p.pattern.String()),
		zap.Strings("extracted_fields", p.extractedFields))

	if err := p.input.Open(ctx); err != nil {
		return err
	}

	p.opened = true
	return nil
}

// Next returns the next row with parsed fields added
func (p *parseOperator) Next(ctx context.Context) (*Row, error) {
	if p.closed {
		return nil, ErrClosed
	}

	if !p.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	row, err := p.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	p.stats.RowsRead++

	// Get the source field value
	sourceValue, exists := row.Get(p.sourceField)
	if !exists {
		// Source field doesn't exist - log warning but don't fail
		p.logger.Debug("Source field not found in row",
			zap.String("field", p.sourceField))
		p.stats.RowsReturned++
		return row, nil
	}

	// Convert to string
	sourceStr, ok := sourceValue.(string)
	if !ok {
		// Try to convert to string
		sourceStr = fmt.Sprintf("%v", sourceValue)
	}

	// Apply regex pattern
	matches := p.pattern.FindStringSubmatch(sourceStr)
	if matches == nil {
		// No match - log debug but don't fail
		p.logger.Debug("Pattern did not match",
			zap.String("field", p.sourceField),
			zap.String("value", sourceStr))
		p.stats.RowsReturned++
		return row, nil
	}

	// Extract named capture groups
	// matches[0] is the full match, matches[1:] are the capture groups
	subexpNames := p.pattern.SubexpNames()
	for i := 1; i < len(matches) && i < len(subexpNames); i++ {
		if subexpNames[i] != "" && i < len(matches) {
			// Add the extracted field to the row
			row.Set(subexpNames[i], matches[i])
		}
	}

	p.stats.RowsReturned++
	return row, nil
}

// Close releases resources
func (p *parseOperator) Close() error {
	p.closed = true
	return p.input.Close()
}

// Stats returns execution statistics
func (p *parseOperator) Stats() *IteratorStats {
	return p.stats
}
