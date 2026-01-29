// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeType_String(t *testing.T) {
	tests := []struct {
		name     string
		nodeType NodeType
		expected string
	}{
		{
			name:     "Query",
			nodeType: NodeTypeQuery,
			expected: "Query",
		},
		{
			name:     "SearchCommand",
			nodeType: NodeTypeSearchCommand,
			expected: "SearchCommand",
		},
		{
			name:     "WhereCommand",
			nodeType: NodeTypeWhereCommand,
			expected: "WhereCommand",
		},
		{
			name:     "FieldsCommand",
			nodeType: NodeTypeFieldsCommand,
			expected: "FieldsCommand",
		},
		{
			name:     "StatsCommand",
			nodeType: NodeTypeStatsCommand,
			expected: "StatsCommand",
		},
		{
			name:     "SortCommand",
			nodeType: NodeTypeSortCommand,
			expected: "SortCommand",
		},
		{
			name:     "HeadCommand",
			nodeType: NodeTypeHeadCommand,
			expected: "HeadCommand",
		},
		{
			name:     "DescribeCommand",
			nodeType: NodeTypeDescribeCommand,
			expected: "DescribeCommand",
		},
		{
			name:     "ShowDatasourcesCommand",
			nodeType: NodeTypeShowDatasourcesCommand,
			expected: "ShowDatasourcesCommand",
		},
		{
			name:     "ExplainCommand",
			nodeType: NodeTypeExplainCommand,
			expected: "ExplainCommand",
		},
		{
			name:     "BinaryExpression",
			nodeType: NodeTypeBinaryExpression,
			expected: "BinaryExpression",
		},
		{
			name:     "UnaryExpression",
			nodeType: NodeTypeUnaryExpression,
			expected: "UnaryExpression",
		},
		{
			name:     "FunctionCall",
			nodeType: NodeTypeFunctionCall,
			expected: "FunctionCall",
		},
		{
			name:     "FieldReference",
			nodeType: NodeTypeFieldReference,
			expected: "FieldReference",
		},
		{
			name:     "Literal",
			nodeType: NodeTypeLiteral,
			expected: "Literal",
		},
		{
			name:     "ListLiteral",
			nodeType: NodeTypeListLiteral,
			expected: "ListLiteral",
		},
		{
			name:     "CaseExpression",
			nodeType: NodeTypeCaseExpression,
			expected: "CaseExpression",
		},
		{
			name:     "WhenClause",
			nodeType: NodeTypeWhenClause,
			expected: "WhenClause",
		},
		{
			name:     "Aggregation",
			nodeType: NodeTypeAggregation,
			expected: "Aggregation",
		},
		{
			name:     "SortKey",
			nodeType: NodeTypeSortKey,
			expected: "SortKey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.nodeType.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNodeType_AllTypesCovered(t *testing.T) {
	// Ensure all node types have unique values and string representations
	types := []NodeType{
		NodeTypeQuery,
		NodeTypeSearchCommand,
		NodeTypeWhereCommand,
		NodeTypeFieldsCommand,
		NodeTypeStatsCommand,
		NodeTypeSortCommand,
		NodeTypeHeadCommand,
		NodeTypeDescribeCommand,
		NodeTypeShowDatasourcesCommand,
		NodeTypeExplainCommand,
		NodeTypeBinaryExpression,
		NodeTypeUnaryExpression,
		NodeTypeFunctionCall,
		NodeTypeFieldReference,
		NodeTypeLiteral,
		NodeTypeListLiteral,
		NodeTypeCaseExpression,
		NodeTypeWhenClause,
		NodeTypeAggregation,
		NodeTypeSortKey,
	}

	seen := make(map[NodeType]bool)
	seenStrings := make(map[string]bool)

	for _, typ := range types {
		// Check uniqueness
		assert.False(t, seen[typ], "Duplicate NodeType value: %v", typ)
		seen[typ] = true

		// Check string representation
		str := typ.String()
		assert.NotEmpty(t, str, "NodeType.String() returned empty for type %v", typ)
		assert.False(t, seenStrings[str], "Duplicate String() value: %s", str)
		seenStrings[str] = true
	}

	// Verify we have all 20 types
	assert.Len(t, types, 20, "Expected 20 node types")
}

func TestBaseNode_Position(t *testing.T) {
	tests := []struct {
		name     string
		baseNode BaseNode
		expected Position
	}{
		{
			name: "valid position",
			baseNode: BaseNode{
				Pos: Position{Line: 1, Column: 5, Offset: 4},
			},
			expected: Position{Line: 1, Column: 5, Offset: 4},
		},
		{
			name:     "zero position",
			baseNode: BaseNode{Pos: Position{Line: 0, Column: 0, Offset: 0}},
			expected: Position{Line: 0, Column: 0, Offset: 0},
		},
		{
			name: "large position",
			baseNode: BaseNode{
				Pos: Position{Line: 10000, Column: 500, Offset: 500000},
			},
			expected: Position{Line: 10000, Column: 500, Offset: 500000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.baseNode.Position()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBaseNode_SetPosition(t *testing.T) {
	node := &BaseNode{}

	// Initially zero
	assert.Equal(t, Position{}, node.Position())

	// Set position
	pos := Position{Line: 5, Column: 10, Offset: 50}
	node.Pos = pos

	// Verify
	assert.Equal(t, pos, node.Position())
}

func TestBaseNode_EdgeCases(t *testing.T) {
	t.Run("uninitialized basenode", func(t *testing.T) {
		node := &BaseNode{}
		pos := node.Position()
		assert.Equal(t, 0, pos.Line)
		assert.Equal(t, 0, pos.Column)
		assert.Equal(t, 0, pos.Offset)
	})

	t.Run("multiple position updates", func(t *testing.T) {
		node := &BaseNode{}

		pos1 := Position{Line: 1, Column: 1, Offset: 0}
		node.Pos = pos1
		assert.Equal(t, pos1, node.Position())

		pos2 := Position{Line: 2, Column: 5, Offset: 10}
		node.Pos = pos2
		assert.Equal(t, pos2, node.Position())

		// Original position should not affect current
		assert.NotEqual(t, pos1, node.Position())
	})
}
