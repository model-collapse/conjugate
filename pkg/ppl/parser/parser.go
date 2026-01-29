// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package parser

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/parser/generated"
)

// Parser wraps the ANTLR4-generated parser and builds AST
type Parser struct {
	errorListener *ErrorListener
}

// NewParser creates a new PPL parser
func NewParser() *Parser {
	return &Parser{
		errorListener: NewErrorListener(),
	}
}

// Parse parses a PPL query string and returns an AST
func (p *Parser) Parse(query string) (*ast.Query, error) {
	// Reset error listener
	p.errorListener.Reset()

	// Create input stream
	input := antlr.NewInputStream(query)

	// Create lexer
	lexer := generated.NewPPLLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(p.errorListener)

	// Create token stream
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create parser
	parser := generated.NewPPLParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(p.errorListener)
	parser.BuildParseTrees = true

	// Parse query (with EOF requirement in grammar)
	tree := parser.Query()

	// Check for syntax errors
	if p.errorListener.HasErrors() {
		return nil, p.errorListener.GetError()
	}

	// Build AST from parse tree
	builder := NewASTBuilder()
	result := builder.VisitQuery(tree.(*generated.QueryContext))

	// Check if result is an error
	if errResult, ok := result.(error); ok {
		return nil, fmt.Errorf("failed to build AST: %w", errResult)
	}

	return result.(*ast.Query), nil
}

// ParseWithRecovery parses with error recovery enabled
func (p *Parser) ParseWithRecovery(query string) (*ast.Query, []error) {
	astQuery, err := p.Parse(query)
	if err != nil {
		return nil, []error{err}
	}
	return astQuery, nil
}

// ValidateSyntax checks if a query has valid syntax without building AST
func (p *Parser) ValidateSyntax(query string) error {
	p.errorListener.Reset()

	input := antlr.NewInputStream(query)
	lexer := generated.NewPPLLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(p.errorListener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := generated.NewPPLParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(p.errorListener)

	// Parse (with EOF requirement in grammar)
	_ = parser.Query()

	if p.errorListener.HasErrors() {
		return p.errorListener.GetError()
	}

	return nil
}

// GetErrorListener returns the error listener for inspection
func (p *Parser) GetErrorListener() *ErrorListener {
	return p.errorListener
}
