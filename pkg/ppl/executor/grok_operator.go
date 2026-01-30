// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/grok"
	"go.uber.org/zap"
)

// grokOperator applies grok patterns to parse unstructured log data
// Grok uses named regular expression patterns to extract structured fields
//
// Examples:
//   grok "%{COMMONAPACHELOG}"
//   grok "%{IP:client_ip} - - \[%{HTTPDATE:timestamp}\]"
//   grok pattern="%{LOGLEVEL:level} %{GREEDYDATA:message}"
type grokOperator struct {
	input          Operator
	logger         *zap.Logger
	ctx            context.Context
	stats          *IteratorStats
	opened         bool
	closed         bool
	pattern        string            // Grok pattern to match
	inputField     string            // Field to parse (default: "_raw")
	grokParser     *grok.Grok        // Compiled grok pattern
	keepOriginal   bool              // Keep original field after parsing
	customPatterns map[string]string // Custom grok patterns
}

// GrokConfig holds configuration for grok operator
type GrokConfig struct {
	Pattern        string            // Grok pattern (required)
	InputField     string            // Field to parse (default: "_raw")
	KeepOriginal   bool              // Keep original field (default: false)
	CustomPatterns map[string]string // Custom pattern definitions
}

// NewGrokOperator creates a new grok operator
func NewGrokOperator(input Operator, config GrokConfig, logger *zap.Logger) (*grokOperator, error) {
	// Validate pattern
	if config.Pattern == "" {
		return nil, fmt.Errorf("grok pattern is required")
	}

	// Set defaults
	if config.InputField == "" {
		config.InputField = "_raw"
	}

	// Compile the grok pattern
	var grokParser *grok.Grok
	var err error

	if config.CustomPatterns != nil && len(config.CustomPatterns) > 0 {
		grokParser, err = grok.NewGrokWithPatterns(config.Pattern, config.CustomPatterns)
	} else {
		grokParser, err = grok.NewGrok(config.Pattern)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to compile grok pattern: %w", err)
	}

	return &grokOperator{
		input:          input,
		logger:         logger,
		stats:          &IteratorStats{},
		pattern:        config.Pattern,
		inputField:     config.InputField,
		grokParser:     grokParser,
		keepOriginal:   config.KeepOriginal,
		customPatterns: config.CustomPatterns,
	}, nil
}

// Open initializes the operator
func (g *grokOperator) Open(ctx context.Context) error {
	if g.opened {
		return nil
	}

	g.ctx = ctx
	g.logger.Debug("Opening grok operator",
		zap.String("pattern", g.pattern),
		zap.String("input_field", g.inputField),
		zap.Bool("keep_original", g.keepOriginal))

	if err := g.input.Open(ctx); err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}

	g.opened = true
	return nil
}

// Next processes the next row, applying grok pattern
func (g *grokOperator) Next(ctx context.Context) (*Row, error) {
	if g.closed {
		return nil, ErrClosed
	}

	if !g.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Read next row
	row, err := g.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	g.stats.RowsRead++

	// Get the input field value
	inputValue, exists := row.Get(g.inputField)
	if !exists {
		// Input field doesn't exist, return row as-is
		g.stats.RowsReturned++
		return row, nil
	}

	// Convert to string
	inputStr, ok := inputValue.(string)
	if !ok {
		// Not a string, try to convert
		inputStr = fmt.Sprintf("%v", inputValue)
	}

	// Apply grok pattern
	match, matched := g.grokParser.Match(inputStr)
	if !matched {
		// No match, return row as-is
		g.logger.Debug("Grok pattern did not match",
			zap.String("input", inputStr))
		g.stats.RowsReturned++
		return row, nil
	}

	// Add extracted fields to the row
	for fieldName, fieldValue := range match.Fields {
		row.Set(fieldName, fieldValue)
	}

	// Remove original field if requested
	if !g.keepOriginal {
		row.Delete(g.inputField)
	}

	g.stats.RowsReturned++
	return row, nil
}

// Close releases resources
func (g *grokOperator) Close() error {
	g.closed = true
	if g.input != nil {
		return g.input.Close()
	}
	return nil
}

// Stats returns execution statistics
func (g *grokOperator) Stats() *IteratorStats {
	return g.stats
}
