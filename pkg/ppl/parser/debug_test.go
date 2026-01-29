// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package parser

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/conjugate/conjugate/pkg/ppl/parser/generated"
)

func TestLexerTokenization(t *testing.T) {
	query := "source=logs"

	input := antlr.NewInputStream(query)
	lexer := generated.NewPPLLexer(input)

	t.Logf("Tokenizing: %s", query)
	for {
		token := lexer.NextToken()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		tokenName := lexer.SymbolicNames[token.GetTokenType()]
		t.Logf("Token: %s = '%s' (type %d)", tokenName, token.GetText(), token.GetTokenType())
	}
}

func TestSimpleParse(t *testing.T) {
	query := "source=logs"

	parser := NewParser()
	_, err := parser.Parse(query)
	if err != nil {
		t.Logf("Parse error: %v", err)
		t.Fail()
	}
}
