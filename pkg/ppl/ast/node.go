// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

// Node is the base interface for all AST nodes
type Node interface {
	// Accept implements the Visitor pattern
	Accept(visitor Visitor) (interface{}, error)

	// Type returns the node type
	Type() NodeType

	// Position returns the source position of this node
	Position() Position

	// String returns a string representation for debugging
	String() string
}

// NodeType represents the type of an AST node
type NodeType int

const (
	// Command nodes
	NodeTypeQuery NodeType = iota
	NodeTypeSearchCommand
	NodeTypeWhereCommand
	NodeTypeFieldsCommand
	NodeTypeStatsCommand
	NodeTypeSortCommand
	NodeTypeHeadCommand
	NodeTypeTailCommand
	NodeTypeTopCommand
	NodeTypeRareCommand
	NodeTypeEvalCommand
	NodeTypeRenameCommand
	NodeTypeReplaceCommand
	NodeTypeFillnullCommand
	NodeTypeParseCommand
	NodeTypeRexCommand
	NodeTypeGrokCommand
	NodeTypeJoinCommand
	NodeTypeLookupCommand
	NodeTypeAppendCommand
	NodeTypeDedupCommand
	NodeTypeBinCommand
	NodeTypeTimechartCommand
	NodeTypeChartCommand
	NodeTypeDescribeCommand
	NodeTypeShowDatasourcesCommand
	NodeTypeExplainCommand
	NodeTypeReverseCommand
	NodeTypeFlattenCommand
	NodeTypeAddtotalsCommand
	NodeTypeAddcoltotalsCommand

	// Expression nodes
	NodeTypeBinaryExpression
	NodeTypeUnaryExpression
	NodeTypeFunctionCall
	NodeTypeFieldReference
	NodeTypeLiteral
	NodeTypeListLiteral
	NodeTypeCaseExpression
	NodeTypeWhenClause

	// Other nodes
	NodeTypeAggregation
	NodeTypeSortKey
	NodeTypeFieldAlias
)

// String returns the string representation of the node type
func (t NodeType) String() string {
	switch t {
	case NodeTypeQuery:
		return "Query"
	case NodeTypeSearchCommand:
		return "SearchCommand"
	case NodeTypeWhereCommand:
		return "WhereCommand"
	case NodeTypeFieldsCommand:
		return "FieldsCommand"
	case NodeTypeStatsCommand:
		return "StatsCommand"
	case NodeTypeSortCommand:
		return "SortCommand"
	case NodeTypeHeadCommand:
		return "HeadCommand"
	case NodeTypeDescribeCommand:
		return "DescribeCommand"
	case NodeTypeShowDatasourcesCommand:
		return "ShowDatasourcesCommand"
	case NodeTypeExplainCommand:
		return "ExplainCommand"
	// Tier 1 commands
	case NodeTypeTailCommand:
		return "TailCommand"
	case NodeTypeTopCommand:
		return "TopCommand"
	case NodeTypeRareCommand:
		return "RareCommand"
	case NodeTypeEvalCommand:
		return "EvalCommand"
	case NodeTypeRenameCommand:
		return "RenameCommand"
	case NodeTypeReplaceCommand:
		return "ReplaceCommand"
	case NodeTypeFillnullCommand:
		return "FillnullCommand"
	case NodeTypeParseCommand:
		return "ParseCommand"
	case NodeTypeRexCommand:
		return "RexCommand"
	case NodeTypeGrokCommand:
		return "GrokCommand"
	case NodeTypeJoinCommand:
		return "JoinCommand"
	case NodeTypeLookupCommand:
		return "LookupCommand"
	case NodeTypeAppendCommand:
		return "AppendCommand"
	case NodeTypeDedupCommand:
		return "DedupCommand"
	case NodeTypeBinCommand:
		return "BinCommand"
	case NodeTypeTimechartCommand:
		return "TimechartCommand"
	case NodeTypeChartCommand:
		return "ChartCommand"
	case NodeTypeReverseCommand:
		return "ReverseCommand"
	case NodeTypeFlattenCommand:
		return "FlattenCommand"
	case NodeTypeAddtotalsCommand:
		return "AddtotalsCommand"
	case NodeTypeAddcoltotalsCommand:
		return "AddcoltotalsCommand"
	// Expression nodes
	case NodeTypeBinaryExpression:
		return "BinaryExpression"
	case NodeTypeUnaryExpression:
		return "UnaryExpression"
	case NodeTypeFunctionCall:
		return "FunctionCall"
	case NodeTypeFieldReference:
		return "FieldReference"
	case NodeTypeLiteral:
		return "Literal"
	case NodeTypeListLiteral:
		return "ListLiteral"
	case NodeTypeCaseExpression:
		return "CaseExpression"
	case NodeTypeWhenClause:
		return "WhenClause"
	case NodeTypeAggregation:
		return "Aggregation"
	case NodeTypeSortKey:
		return "SortKey"
	default:
		return "Unknown"
	}
}

// BaseNode provides common functionality for all nodes
type BaseNode struct {
	Pos Position
}

// Position returns the source position
func (b *BaseNode) Position() Position {
	return b.Pos
}
