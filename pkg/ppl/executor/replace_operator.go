// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// replaceOperator replaces values in a specified field
type replaceOperator struct {
	input    Operator
	mappings []*ast.ReplaceMapping
	field    string
	logger   *zap.Logger

	ctx    context.Context
	stats  *IteratorStats
	opened bool
	closed bool

	// Compiled regex patterns for efficient matching
	regexMappings []regexMapping
}

// regexMapping holds compiled regex and replacement value
type regexMapping struct {
	regex       *regexp.Regexp
	replacement string
	isRegex     bool
}

// NewReplaceOperator creates a new replace operator
func NewReplaceOperator(
	input Operator,
	mappings []*ast.ReplaceMapping,
	field string,
	logger *zap.Logger,
) *replaceOperator {
	return &replaceOperator{
		input:    input,
		mappings: mappings,
		field:    field,
		logger:   logger,
		stats:    &IteratorStats{},
	}
}

// Open initializes the operator
func (r *replaceOperator) Open(ctx context.Context) error {
	if r.opened {
		return nil
	}

	r.ctx = ctx
	r.logger.Debug("Opening replace operator",
		zap.Int("num_mappings", len(r.mappings)),
		zap.String("field", r.field))

	// Pre-compile any regex patterns
	r.regexMappings = make([]regexMapping, len(r.mappings))
	for i, mapping := range r.mappings {
		oldValue := r.evaluateExpression(nil, mapping.OldValue)
		newValue := r.evaluateExpression(nil, mapping.NewValue)

		oldStr := fmt.Sprintf("%v", oldValue)
		newStr := fmt.Sprintf("%v", newValue)

		// Check if old value looks like a regex pattern (starts with /)
		if strings.HasPrefix(oldStr, "/") && strings.HasSuffix(oldStr, "/") && len(oldStr) > 2 {
			// Extract regex pattern
			pattern := oldStr[1 : len(oldStr)-1]
			regex, err := regexp.Compile(pattern)
			if err != nil {
				r.logger.Warn("Invalid regex pattern, treating as literal",
					zap.String("pattern", pattern),
					zap.Error(err))
				r.regexMappings[i] = regexMapping{
					replacement: newStr,
					isRegex:     false,
				}
			} else {
				r.regexMappings[i] = regexMapping{
					regex:       regex,
					replacement: newStr,
					isRegex:     true,
				}
			}
		} else {
			// Literal string replacement
			r.regexMappings[i] = regexMapping{
				replacement: newStr,
				isRegex:     false,
			}
		}
	}

	if err := r.input.Open(ctx); err != nil {
		return err
	}

	r.opened = true
	return nil
}

// Next returns the next row with replacements applied
func (r *replaceOperator) Next(ctx context.Context) (*Row, error) {
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

	// Get the field value
	val, ok := row.Get(r.field)
	if !ok {
		// Field doesn't exist, return row unchanged
		r.stats.RowsReturned++
		return row, nil
	}

	// Convert value to string for replacement
	strVal := fmt.Sprintf("%v", val)

	// Apply all replacements in order
	for i, mapping := range r.mappings {
		regexMap := r.regexMappings[i]

		if regexMap.isRegex {
			// Regex replacement
			strVal = regexMap.regex.ReplaceAllString(strVal, regexMap.replacement)
		} else {
			// Literal string replacement
			oldValue := r.evaluateExpression(row, mapping.OldValue)
			oldStr := fmt.Sprintf("%v", oldValue)
			strVal = strings.ReplaceAll(strVal, oldStr, regexMap.replacement)
		}
	}

	// Set the modified value back
	row.Set(r.field, strVal)

	r.stats.RowsReturned++
	return row, nil
}

// evaluateExpression evaluates a simple expression (literal or field reference)
func (r *replaceOperator) evaluateExpression(row *Row, expr ast.Expression) interface{} {
	switch ex := expr.(type) {
	case *ast.Literal:
		return ex.Value

	case *ast.FieldReference:
		if row != nil {
			val, _ := row.Get(ex.Name)
			return val
		}
		return nil

	default:
		return nil
	}
}

// Close releases resources
func (r *replaceOperator) Close() error {
	r.closed = true
	return r.input.Close()
}

// Stats returns execution statistics
func (r *replaceOperator) Stats() *IteratorStats {
	return r.stats
}
