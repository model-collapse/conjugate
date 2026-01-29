// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"
	"strings"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// filterOperator filters rows based on a condition
type filterOperator struct {
	input     Operator
	condition ast.Expression
	logger    *zap.Logger

	ctx    context.Context
	stats  *IteratorStats
	opened bool
	closed bool
}

// NewFilterOperator creates a new filter operator
func NewFilterOperator(input Operator, condition ast.Expression, logger *zap.Logger) *filterOperator {
	return &filterOperator{
		input:     input,
		condition: condition,
		logger:    logger,
		stats:     &IteratorStats{},
	}
}

// Open initializes the operator
func (f *filterOperator) Open(ctx context.Context) error {
	if f.opened {
		return nil
	}

	f.ctx = ctx
	f.logger.Debug("Opening filter operator",
		zap.String("condition", f.condition.String()))

	// Open input
	if err := f.input.Open(ctx); err != nil {
		return err
	}

	f.opened = true
	return nil
}

// Next returns the next row that matches the condition
func (f *filterOperator) Next(ctx context.Context) (*Row, error) {
	if f.closed {
		return nil, ErrClosed
	}

	if !f.opened {
		return nil, ErrClosed
	}

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// Get next row from input
		row, err := f.input.Next(ctx)
		if err != nil {
			return nil, err
		}

		f.stats.RowsRead++

		// Evaluate condition
		matches, err := f.evaluateCondition(row, f.condition)
		if err != nil {
			f.logger.Warn("Error evaluating filter condition",
				zap.Error(err))
			continue // Skip rows with evaluation errors
		}

		if matches {
			f.stats.RowsReturned++
			return row, nil
		}
	}
}

// evaluateCondition evaluates an expression against a row
func (f *filterOperator) evaluateCondition(row *Row, expr ast.Expression) (bool, error) {
	result, err := f.evaluate(row, expr)
	if err != nil {
		return false, err
	}

	// Convert result to bool
	switch v := result.(type) {
	case bool:
		return v, nil
	case int64:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return v != "", nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", result)
	}
}

// evaluate evaluates an expression and returns the result
func (f *filterOperator) evaluate(row *Row, expr ast.Expression) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.Literal:
		return e.Value, nil

	case *ast.FieldReference:
		value, _ := row.Get(e.Name)
		return value, nil

	case *ast.BinaryExpression:
		return f.evaluateBinaryExpr(row, e)

	case *ast.UnaryExpression:
		return f.evaluateUnaryExpr(row, e)

	case *ast.FunctionCall:
		return f.evaluateFunctionCall(row, e)

	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// evaluateBinaryExpr evaluates a binary expression
func (f *filterOperator) evaluateBinaryExpr(row *Row, expr *ast.BinaryExpression) (interface{}, error) {
	left, err := f.evaluate(row, expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := f.evaluate(row, expr.Right)
	if err != nil {
		return nil, err
	}

	op := strings.ToUpper(expr.Operator)

	// Handle logical operators
	switch op {
	case "AND":
		return toBool(left) && toBool(right), nil
	case "OR":
		return toBool(left) || toBool(right), nil
	}

	// Handle comparison operators
	switch op {
	case "=", "==":
		return compare(left, right) == 0, nil
	case "!=", "<>":
		return compare(left, right) != 0, nil
	case "<":
		return compare(left, right) < 0, nil
	case "<=":
		return compare(left, right) <= 0, nil
	case ">":
		return compare(left, right) > 0, nil
	case ">=":
		return compare(left, right) >= 0, nil
	case "LIKE":
		return matchLike(toString(left), toString(right)), nil
	case "IN":
		return matchIn(left, right), nil
	}

	// Handle arithmetic operators
	switch op {
	case "+":
		return toFloat(left) + toFloat(right), nil
	case "-":
		return toFloat(left) - toFloat(right), nil
	case "*":
		return toFloat(left) * toFloat(right), nil
	case "/":
		r := toFloat(right)
		if r == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return toFloat(left) / r, nil
	case "%":
		return int64(toFloat(left)) % int64(toFloat(right)), nil
	}

	return nil, fmt.Errorf("unsupported operator: %s", expr.Operator)
}

// evaluateUnaryExpr evaluates a unary expression
func (f *filterOperator) evaluateUnaryExpr(row *Row, expr *ast.UnaryExpression) (interface{}, error) {
	operand, err := f.evaluate(row, expr.Operand)
	if err != nil {
		return nil, err
	}

	switch strings.ToUpper(expr.Operator) {
	case "NOT":
		return !toBool(operand), nil
	case "-":
		return -toFloat(operand), nil
	default:
		return nil, fmt.Errorf("unsupported unary operator: %s", expr.Operator)
	}
}

// evaluateFunctionCall evaluates a function call
func (f *filterOperator) evaluateFunctionCall(row *Row, fn *ast.FunctionCall) (interface{}, error) {
	// Evaluate arguments
	args := make([]interface{}, len(fn.Arguments))
	for i, arg := range fn.Arguments {
		val, err := f.evaluate(row, arg)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	// Execute built-in functions
	funcName := strings.ToLower(fn.Name)
	switch funcName {
	case "abs":
		if len(args) != 1 {
			return nil, fmt.Errorf("abs() requires 1 argument")
		}
		v := toFloat(args[0])
		if v < 0 {
			return -v, nil
		}
		return v, nil

	case "upper":
		if len(args) != 1 {
			return nil, fmt.Errorf("upper() requires 1 argument")
		}
		return strings.ToUpper(toString(args[0])), nil

	case "lower":
		if len(args) != 1 {
			return nil, fmt.Errorf("lower() requires 1 argument")
		}
		return strings.ToLower(toString(args[0])), nil

	case "length", "len":
		if len(args) != 1 {
			return nil, fmt.Errorf("length() requires 1 argument")
		}
		return int64(len(toString(args[0]))), nil

	case "isnull":
		if len(args) != 1 {
			return nil, fmt.Errorf("isnull() requires 1 argument")
		}
		return args[0] == nil, nil

	case "isnotnull":
		if len(args) != 1 {
			return nil, fmt.Errorf("isnotnull() requires 1 argument")
		}
		return args[0] != nil, nil

	default:
		return nil, fmt.Errorf("unsupported function: %s", fn.Name)
	}
}

// Close releases resources
func (f *filterOperator) Close() error {
	f.closed = true
	return f.input.Close()
}

// Stats returns execution statistics
func (f *filterOperator) Stats() *IteratorStats {
	return f.stats
}

// Helper functions

func toBool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case int64:
		return val != 0
	case float64:
		return val != 0
	case string:
		return val != ""
	default:
		return true
	}
}

func toFloat(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int64:
		return float64(val)
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

func compare(a, b interface{}) int {
	// Handle nil cases
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Try numeric comparison
	aNum, aIsNum := toNumber(a)
	bNum, bIsNum := toNumber(b)
	if aIsNum && bIsNum {
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
		return 0
	}

	// Fall back to string comparison
	aStr := toString(a)
	bStr := toString(b)
	if aStr < bStr {
		return -1
	}
	if aStr > bStr {
		return 1
	}
	return 0
}

func toNumber(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int64:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	default:
		return 0, false
	}
}

func matchLike(str, pattern string) bool {
	// Simple LIKE matching - convert % to .* for regex-like matching
	// This is a simplified implementation
	pattern = strings.ReplaceAll(pattern, "%", ".*")
	pattern = strings.ReplaceAll(pattern, "_", ".")

	// Simple prefix/suffix matching for common cases
	if strings.HasPrefix(pattern, ".*") && strings.HasSuffix(pattern, ".*") {
		return strings.Contains(str, pattern[2:len(pattern)-2])
	}
	if strings.HasPrefix(pattern, ".*") {
		return strings.HasSuffix(str, pattern[2:])
	}
	if strings.HasSuffix(pattern, ".*") {
		return strings.HasPrefix(str, pattern[:len(pattern)-2])
	}

	return str == pattern
}

func matchIn(value, list interface{}) bool {
	// Check if value is in a list
	switch l := list.(type) {
	case []interface{}:
		for _, item := range l {
			if compare(value, item) == 0 {
				return true
			}
		}
	}
	return false
}

// evalFunction evaluates a built-in function with given arguments
// This is a shared helper used by multiple operators (filter, eval)
func evalFunction(funcName string, args []interface{}) (interface{}, error) {
	name := strings.ToLower(funcName)
	switch name {
	case "abs":
		if len(args) != 1 {
			return nil, fmt.Errorf("abs() requires 1 argument")
		}
		v := toFloat(args[0])
		if v < 0 {
			return -v, nil
		}
		return v, nil

	case "upper":
		if len(args) != 1 {
			return nil, fmt.Errorf("upper() requires 1 argument")
		}
		return strings.ToUpper(toString(args[0])), nil

	case "lower":
		if len(args) != 1 {
			return nil, fmt.Errorf("lower() requires 1 argument")
		}
		return strings.ToLower(toString(args[0])), nil

	case "length", "len":
		if len(args) != 1 {
			return nil, fmt.Errorf("length() requires 1 argument")
		}
		return int64(len(toString(args[0]))), nil

	case "isnull":
		if len(args) != 1 {
			return nil, fmt.Errorf("isnull() requires 1 argument")
		}
		return args[0] == nil, nil

	case "isnotnull":
		if len(args) != 1 {
			return nil, fmt.Errorf("isnotnull() requires 1 argument")
		}
		return args[0] != nil, nil

	case "trim":
		if len(args) != 1 {
			return nil, fmt.Errorf("trim() requires 1 argument")
		}
		return strings.TrimSpace(toString(args[0])), nil

	case "ltrim":
		if len(args) != 1 {
			return nil, fmt.Errorf("ltrim() requires 1 argument")
		}
		return strings.TrimLeft(toString(args[0]), " \t\n\r"), nil

	case "rtrim":
		if len(args) != 1 {
			return nil, fmt.Errorf("rtrim() requires 1 argument")
		}
		return strings.TrimRight(toString(args[0]), " \t\n\r"), nil

	case "substr", "substring":
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("substr() requires 2-3 arguments")
		}
		s := toString(args[0])
		start := int(toFloat(args[1]))
		if start < 0 {
			start = 0
		}
		if start >= len(s) {
			return "", nil
		}
		if len(args) == 3 {
			length := int(toFloat(args[2]))
			if start+length > len(s) {
				length = len(s) - start
			}
			return s[start : start+length], nil
		}
		return s[start:], nil

	case "concat":
		var result strings.Builder
		for _, arg := range args {
			result.WriteString(toString(arg))
		}
		return result.String(), nil

	case "coalesce":
		for _, arg := range args {
			if arg != nil {
				return arg, nil
			}
		}
		return nil, nil

	case "if":
		if len(args) != 3 {
			return nil, fmt.Errorf("if() requires 3 arguments")
		}
		if toBool(args[0]) {
			return args[1], nil
		}
		return args[2], nil

	case "round":
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("round() requires 1-2 arguments")
		}
		v := toFloat(args[0])
		precision := 0
		if len(args) == 2 {
			precision = int(toFloat(args[1]))
		}
		multiplier := 1.0
		for i := 0; i < precision; i++ {
			multiplier *= 10
		}
		return float64(int64(v*multiplier+0.5)) / multiplier, nil

	case "floor":
		if len(args) != 1 {
			return nil, fmt.Errorf("floor() requires 1 argument")
		}
		return float64(int64(toFloat(args[0]))), nil

	case "ceil", "ceiling":
		if len(args) != 1 {
			return nil, fmt.Errorf("ceil() requires 1 argument")
		}
		v := toFloat(args[0])
		i := int64(v)
		if float64(i) < v {
			i++
		}
		return float64(i), nil

	case "sqrt":
		if len(args) != 1 {
			return nil, fmt.Errorf("sqrt() requires 1 argument")
		}
		v := toFloat(args[0])
		if v < 0 {
			return nil, fmt.Errorf("sqrt() requires non-negative argument")
		}
		return sqrt(v), nil

	case "pow", "power":
		if len(args) != 2 {
			return nil, fmt.Errorf("pow() requires 2 arguments")
		}
		return pow(toFloat(args[0]), toFloat(args[1])), nil

	default:
		return nil, fmt.Errorf("unsupported function: %s", funcName)
	}
}

// sqrt computes square root
func sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	// Newton's method
	z := x / 2
	for i := 0; i < 100; i++ {
		z2 := (z + x/z) / 2
		if z == z2 || (z2-z < 0.0001 && z-z2 < 0.0001) {
			break
		}
		z = z2
	}
	return z
}

// pow computes x^y
func pow(x, y float64) float64 {
	result := 1.0
	if y < 0 {
		x = 1 / x
		y = -y
	}
	for i := 0; i < int(y); i++ {
		result *= x
	}
	return result
}
