// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package parser_test

import (
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/parser"
)

// Example demonstrates basic parser usage
func Example_basicParsing() {
	// Create parser
	p := parser.NewParser()

	// Parse a simple query
	query := "source=logs | where status = 200 | head 10"
	ast, err := p.Parse(query)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	fmt.Printf("Successfully parsed query with %d commands\n", len(ast.Commands))
	// Output: Successfully parsed query with 3 commands
}

// Example_syntaxValidation demonstrates syntax checking without AST building
func Example_syntaxValidation() {
	p := parser.NewParser()

	// Valid query
	if err := p.ValidateSyntax("source=logs | where status = 200"); err == nil {
		fmt.Println("Valid syntax")
	}

	// Invalid query
	if err := p.ValidateSyntax("source=logs where"); err != nil {
		fmt.Println("Invalid syntax detected")
	}

	// Output:
	// Valid syntax
	// Invalid syntax detected
}

// Example_errorHandling demonstrates detailed error reporting
func Example_errorHandling() {
	p := parser.NewParser()

	// Parse invalid query
	_, err := p.Parse("source=logs | where")
	if err != nil {
		// Error includes line and column information
		fmt.Println("Error:", err)
	}
}

// Example_complexQuery demonstrates parsing a multi-command analytics query
func Example_complexQuery() {
	p := parser.NewParser()

	query := `
		source=logs
		| where timestamp > '2024-01-01' AND method = 'GET'
		| stats count() as requests, avg(response_time) as avg_time by endpoint
		| sort requests desc
		| head 10
	`

	ast, err := p.Parse(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Inspect parsed commands
	fmt.Printf("Parsed %d commands:\n", len(ast.Commands))
	for i, cmd := range ast.Commands {
		fmt.Printf("  %d. %s\n", i+1, cmd.Type())
	}

	// Output:
	// Parsed 5 commands:
	//   1. SearchCommand
	//   2. WhereCommand
	//   3. StatsCommand
	//   4. SortCommand
	//   5. HeadCommand
}

// Example_astTraversal demonstrates walking the AST
func Example_astTraversal() {
	// Note: This example shows the pattern, but requires generated code to run
	p := parser.NewParser()

	query := "source=logs | where status = 200"
	astQuery, err := p.Parse(query)
	if err != nil {
		return
	}

	// Access specific command types
	for _, cmd := range astQuery.Commands {
		switch c := cmd.(type) {
		case *ast.SearchCommand:
			fmt.Printf("Search source: %s\n", c.Source)
		case *ast.WhereCommand:
			fmt.Printf("Where condition: %s\n", c.Condition.String())
		}
	}
}
