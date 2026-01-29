// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

import (
	"fmt"
	"strings"
)

// Expression represents an expression in PPL
type Expression interface {
	Node
	expressionNode()
}

// BinaryExpression represents a binary operation (e.g., a + b, x = 5)
type BinaryExpression struct {
	BaseNode
	Left     Expression
	Operator string // =, !=, <, >, <=, >=, +, -, *, /, AND, OR, etc.
	Right    Expression
}

func (b *BinaryExpression) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinaryExpression(b)
}

func (b *BinaryExpression) Type() NodeType { return NodeTypeBinaryExpression }
func (b *BinaryExpression) expressionNode() {}
func (b *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Operator, b.Right.String())
}

// UnaryExpression represents a unary operation (e.g., NOT x, -5)
type UnaryExpression struct {
	BaseNode
	Operator string // NOT, -, +
	Operand  Expression
}

func (u *UnaryExpression) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitUnaryExpression(u)
}

func (u *UnaryExpression) Type() NodeType { return NodeTypeUnaryExpression }
func (u *UnaryExpression) expressionNode() {}
func (u *UnaryExpression) String() string {
	return fmt.Sprintf("(%s %s)", u.Operator, u.Operand.String())
}

// FunctionCall represents a function call (e.g., count(), sum(price))
type FunctionCall struct {
	BaseNode
	Name      string
	Arguments []Expression
	Distinct  bool // For COUNT(DISTINCT field)
}

func (f *FunctionCall) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFunctionCall(f)
}

func (f *FunctionCall) Type() NodeType { return NodeTypeFunctionCall }
func (f *FunctionCall) expressionNode() {}
func (f *FunctionCall) String() string {
	args := make([]string, len(f.Arguments))
	for i, arg := range f.Arguments {
		args[i] = arg.String()
	}
	argsStr := strings.Join(args, ", ")
	if f.Distinct {
		return fmt.Sprintf("%s(DISTINCT %s)", f.Name, argsStr)
	}
	return fmt.Sprintf("%s(%s)", f.Name, argsStr)
}

// FieldReference represents a field reference (e.g., user.name, price)
type FieldReference struct {
	BaseNode
	Name string // Can include dots for nested fields (e.g., "user.address.city")
}

func (f *FieldReference) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFieldReference(f)
}

func (f *FieldReference) Type() NodeType { return NodeTypeFieldReference }
func (f *FieldReference) expressionNode() {}
func (f *FieldReference) String() string  { return f.Name }

// LiteralType represents the type of a literal value
type LiteralType int

const (
	LiteralTypeNull LiteralType = iota
	LiteralTypeBool
	LiteralTypeInt
	LiteralTypeFloat
	LiteralTypeString
)

// Literal represents a literal value (e.g., 42, "hello", true)
type Literal struct {
	BaseNode
	Value      interface{} // The actual value
	LiteralTyp LiteralType
}

func (l *Literal) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLiteral(l)
}

func (l *Literal) Type() NodeType { return NodeTypeLiteral }
func (l *Literal) expressionNode() {}
func (l *Literal) String() string {
	if l.Value == nil {
		return "null"
	}
	switch l.LiteralTyp {
	case LiteralTypeString:
		return fmt.Sprintf("\"%v\"", l.Value)
	default:
		return fmt.Sprintf("%v", l.Value)
	}
}

// ListLiteral represents a list of values (e.g., IN (1, 2, 3))
type ListLiteral struct {
	BaseNode
	Values []Expression
}

func (l *ListLiteral) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitListLiteral(l)
}

func (l *ListLiteral) Type() NodeType { return NodeTypeListLiteral }
func (l *ListLiteral) expressionNode() {}
func (l *ListLiteral) String() string {
	vals := make([]string, len(l.Values))
	for i, v := range l.Values {
		vals[i] = v.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(vals, ", "))
}

// WhenClause represents a WHEN clause in a CASE expression
type WhenClause struct {
	BaseNode
	Condition Expression
	Result    Expression
}

func (w *WhenClause) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitWhenClause(w)
}

func (w *WhenClause) Type() NodeType { return NodeTypeWhenClause }
func (w *WhenClause) String() string {
	return fmt.Sprintf("WHEN %s THEN %s", w.Condition.String(), w.Result.String())
}

// CaseExpression represents a CASE expression
type CaseExpression struct {
	BaseNode
	WhenClauses []*WhenClause
	ElseResult  Expression // Optional
}

func (c *CaseExpression) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitCaseExpression(c)
}

func (c *CaseExpression) Type() NodeType { return NodeTypeCaseExpression }
func (c *CaseExpression) expressionNode() {}
func (c *CaseExpression) String() string {
	parts := []string{"CASE"}
	for _, when := range c.WhenClauses {
		parts = append(parts, when.String())
	}
	if c.ElseResult != nil {
		parts = append(parts, fmt.Sprintf("ELSE %s", c.ElseResult.String()))
	}
	parts = append(parts, "END")
	return strings.Join(parts, " ")
}
