// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// evalOperator evaluates expressions and adds computed fields
type evalOperator struct {
	input       Operator
	assignments []*ast.EvalAssignment
	logger      *zap.Logger

	ctx    context.Context
	stats  *IteratorStats
	opened bool
	closed bool
}

// NewEvalOperator creates a new eval operator
func NewEvalOperator(
	input Operator,
	assignments []*ast.EvalAssignment,
	logger *zap.Logger,
) *evalOperator {
	return &evalOperator{
		input:       input,
		assignments: assignments,
		logger:      logger,
		stats:       &IteratorStats{},
	}
}

// Open initializes the operator
func (e *evalOperator) Open(ctx context.Context) error {
	if e.opened {
		return nil
	}

	e.ctx = ctx
	e.logger.Debug("Opening eval operator",
		zap.Int("num_assignments", len(e.assignments)))

	if err := e.input.Open(ctx); err != nil {
		return err
	}

	e.opened = true
	return nil
}

// Next returns the next row with computed fields added
func (e *evalOperator) Next(ctx context.Context) (*Row, error) {
	if e.closed {
		return nil, ErrClosed
	}

	if !e.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	row, err := e.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	e.stats.RowsRead++

	// Evaluate each assignment and add to row
	for _, assignment := range e.assignments {
		value, err := e.evaluate(row, assignment.Expression)
		if err != nil {
			e.logger.Warn("Failed to evaluate expression",
				zap.String("field", assignment.Field),
				zap.Error(err))
			continue
		}
		row.Set(assignment.Field, value)
	}

	e.stats.RowsReturned++
	return row, nil
}

// evaluate evaluates an expression against a row
func (e *evalOperator) evaluate(row *Row, expr ast.Expression) (interface{}, error) {
	switch ex := expr.(type) {
	case *ast.Literal:
		return ex.Value, nil

	case *ast.FieldReference:
		val, _ := row.Get(ex.Name)
		return val, nil

	case *ast.BinaryExpression:
		return e.evaluateBinaryExpr(row, ex)

	case *ast.UnaryExpression:
		return e.evaluateUnaryExpr(row, ex)

	case *ast.FunctionCall:
		return e.evaluateFunctionCall(row, ex)

	default:
		return nil, nil
	}
}

// evaluateBinaryExpr evaluates a binary expression
func (e *evalOperator) evaluateBinaryExpr(row *Row, expr *ast.BinaryExpression) (interface{}, error) {
	left, err := e.evaluate(row, expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := e.evaluate(row, expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case "+":
		leftNum, _ := toNumber(left)
		rightNum, _ := toNumber(right)
		return leftNum + rightNum, nil
	case "-":
		leftNum, _ := toNumber(left)
		rightNum, _ := toNumber(right)
		return leftNum - rightNum, nil
	case "*":
		leftNum, _ := toNumber(left)
		rightNum, _ := toNumber(right)
		return leftNum * rightNum, nil
	case "/":
		leftNum, _ := toNumber(left)
		rightNum, _ := toNumber(right)
		if rightNum == 0 {
			return nil, nil // Division by zero
		}
		return leftNum / rightNum, nil
	case "%":
		leftNum, _ := toNumber(left)
		rightNum, _ := toNumber(right)
		if rightNum == 0 {
			return nil, nil // Division by zero
		}
		return float64(int64(leftNum) % int64(rightNum)), nil
	case "=", "==":
		return left == right, nil
	case "!=", "<>":
		return left != right, nil
	case "<":
		return compare(left, right) < 0, nil
	case "<=":
		return compare(left, right) <= 0, nil
	case ">":
		return compare(left, right) > 0, nil
	case ">=":
		return compare(left, right) >= 0, nil
	case "AND", "and":
		leftBool := toBool(left)
		rightBool := toBool(right)
		return leftBool && rightBool, nil
	case "OR", "or":
		leftBool := toBool(left)
		rightBool := toBool(right)
		return leftBool || rightBool, nil
	case ".":
		// String concatenation
		return toString(left) + toString(right), nil
	default:
		return nil, nil
	}
}

// evaluateUnaryExpr evaluates a unary expression
func (e *evalOperator) evaluateUnaryExpr(row *Row, expr *ast.UnaryExpression) (interface{}, error) {
	operand, err := e.evaluate(row, expr.Operand)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case "NOT", "not", "!":
		return !toBool(operand), nil
	case "-":
		num, _ := toNumber(operand)
		return -num, nil
	default:
		return operand, nil
	}
}

// evaluateFunctionCall evaluates a function call
func (e *evalOperator) evaluateFunctionCall(row *Row, fc *ast.FunctionCall) (interface{}, error) {
	// Evaluate arguments
	args := make([]interface{}, len(fc.Arguments))
	for i, arg := range fc.Arguments {
		val, err := e.evaluate(row, arg)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	// Call function
	return evalFunction(fc.Name, args)
}

// Close releases resources
func (e *evalOperator) Close() error {
	e.closed = true
	return e.input.Close()
}

// Stats returns execution statistics
func (e *evalOperator) Stats() *IteratorStats {
	return e.stats
}
