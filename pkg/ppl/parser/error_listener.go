// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package parser

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
)

// SyntaxError represents a parsing error with location
type SyntaxError struct {
	Line    int
	Column  int
	Message string
	Token   string
}

func (e *SyntaxError) Error() string {
	if e.Token != "" {
		return fmt.Sprintf("syntax error at line %d, column %d: %s (near '%s')",
			e.Line, e.Column, e.Message, e.Token)
	}
	return fmt.Sprintf("syntax error at line %d, column %d: %s",
		e.Line, e.Column, e.Message)
}

// ErrorListener collects syntax errors during parsing
type ErrorListener struct {
	*antlr.DefaultErrorListener
	errors []*SyntaxError
}

// NewErrorListener creates a new error listener
func NewErrorListener() *ErrorListener {
	return &ErrorListener{
		DefaultErrorListener: antlr.NewDefaultErrorListener(),
		errors:               make([]*SyntaxError, 0),
	}
}

// SyntaxError is called when a syntax error is encountered
func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{},
	line, column int, msg string, e antlr.RecognitionException) {

	// Extract token text if available
	tokenText := ""
	if token, ok := offendingSymbol.(antlr.Token); ok {
		tokenText = token.GetText()
	}

	// Create syntax error
	syntaxErr := &SyntaxError{
		Line:    line,
		Column:  column,
		Message: l.enhanceErrorMessage(msg),
		Token:   tokenText,
	}

	l.errors = append(l.errors, syntaxErr)
}

// enhanceErrorMessage makes error messages more user-friendly
func (l *ErrorListener) enhanceErrorMessage(msg string) string {
	// Common error patterns and their improvements
	replacements := map[string]string{
		"mismatched input":    "unexpected token",
		"missing":             "expected",
		"extraneous input":    "unexpected",
		"no viable alternative": "invalid syntax",
	}

	enhanced := msg
	for old, new := range replacements {
		if strings.Contains(strings.ToLower(msg), strings.ToLower(old)) {
			enhanced = strings.Replace(enhanced, old, new, 1)
			break
		}
	}

	return enhanced
}

// HasErrors returns true if any errors were collected
func (l *ErrorListener) HasErrors() bool {
	return len(l.errors) > 0
}

// GetErrors returns all collected errors
func (l *ErrorListener) GetErrors() []*SyntaxError {
	return l.errors
}

// GetError returns a single error combining all errors
func (l *ErrorListener) GetError() error {
	if !l.HasErrors() {
		return nil
	}

	if len(l.errors) == 1 {
		return l.errors[0]
	}

	// Multiple errors - combine them
	messages := make([]string, len(l.errors))
	for i, err := range l.errors {
		messages[i] = err.Error()
	}

	return fmt.Errorf("multiple syntax errors:\n  - %s",
		strings.Join(messages, "\n  - "))
}

// Reset clears all collected errors
func (l *ErrorListener) Reset() {
	l.errors = make([]*SyntaxError, 0)
}

// GetPosition returns the position of the first error as an ast.Position
func (l *ErrorListener) GetPosition() ast.Position {
	if !l.HasErrors() {
		return ast.NoPos
	}

	firstErr := l.errors[0]
	return ast.NewPosition(firstErr.Line, firstErr.Column, -1)
}
