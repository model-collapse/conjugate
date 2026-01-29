// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package analyzer

import (
	"fmt"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
)

// TypeChecker performs type inference and checking on expressions
type TypeChecker struct {
	schema *Schema
	scope  *Scope
}

// NewTypeChecker creates a new type checker
func NewTypeChecker(schema *Schema, scope *Scope) *TypeChecker {
	return &TypeChecker{
		schema: schema,
		scope:  scope,
	}
}

// InferType infers the type of an expression
func (tc *TypeChecker) InferType(expr ast.Expression) (FieldType, error) {
	switch e := expr.(type) {
	case *ast.Literal:
		return tc.inferLiteralType(e)

	case *ast.FieldReference:
		return tc.inferFieldRefType(e)

	case *ast.BinaryExpression:
		return tc.inferBinaryExprType(e)

	case *ast.UnaryExpression:
		return tc.inferUnaryExprType(e)

	case *ast.FunctionCall:
		return tc.inferFunctionType(e)

	case *ast.CaseExpression:
		return tc.inferCaseExprType(e)

	default:
		return FieldTypeUnknown, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// inferLiteralType infers the type of a literal
func (tc *TypeChecker) inferLiteralType(lit *ast.Literal) (FieldType, error) {
	switch lit.LiteralTyp {
	case ast.LiteralTypeInt:
		return FieldTypeLong, nil
	case ast.LiteralTypeFloat:
		return FieldTypeDouble, nil
	case ast.LiteralTypeString:
		return FieldTypeString, nil
	case ast.LiteralTypeBool:
		return FieldTypeBool, nil
	case ast.LiteralTypeNull:
		return FieldTypeUnknown, nil // NULL has no specific type
	default:
		return FieldTypeUnknown, fmt.Errorf("unknown literal type: %v", lit.LiteralTyp)
	}
}

// inferFieldRefType infers the type of a field reference
func (tc *TypeChecker) inferFieldRefType(ref *ast.FieldReference) (FieldType, error) {
	// Try scope first (for aliases)
	if tc.scope != nil && tc.scope.Has(ref.Name) {
		return tc.scope.GetType(ref.Name)
	}

	// Fall back to schema
	if tc.schema != nil {
		return tc.schema.FieldType(ref.Name)
	}

	return FieldTypeUnknown, fmt.Errorf("field %s not found", ref.Name)
}

// inferBinaryExprType infers the type of a binary expression
func (tc *TypeChecker) inferBinaryExprType(expr *ast.BinaryExpression) (FieldType, error) {
	leftType, err := tc.InferType(expr.Left)
	if err != nil {
		return FieldTypeUnknown, err
	}

	rightType, err := tc.InferType(expr.Right)
	if err != nil {
		return FieldTypeUnknown, err
	}

	// Comparison operators return boolean
	if isComparisonOp(expr.Operator) {
		// Validate types are compatible for comparison
		if err := tc.validateComparison(leftType, rightType, expr.Operator); err != nil {
			return FieldTypeUnknown, err
		}
		return FieldTypeBool, nil
	}

	// Logical operators (AND, OR) expect boolean operands and return boolean
	if isLogicalOp(expr.Operator) {
		if leftType != FieldTypeBool {
			return FieldTypeUnknown, fmt.Errorf("left operand of %s must be boolean, got %s", expr.Operator, leftType)
		}
		if rightType != FieldTypeBool {
			return FieldTypeUnknown, fmt.Errorf("right operand of %s must be boolean, got %s", expr.Operator, rightType)
		}
		return FieldTypeBool, nil
	}

	// Arithmetic operators
	if isArithmeticOp(expr.Operator) {
		return tc.inferArithmeticType(leftType, rightType, expr.Operator)
	}

	return FieldTypeUnknown, fmt.Errorf("unsupported operator: %s", expr.Operator)
}

// inferUnaryExprType infers the type of a unary expression
func (tc *TypeChecker) inferUnaryExprType(expr *ast.UnaryExpression) (FieldType, error) {
	operandType, err := tc.InferType(expr.Operand)
	if err != nil {
		return FieldTypeUnknown, err
	}

	switch expr.Operator {
	case "NOT", "!":
		if operandType != FieldTypeBool {
			return FieldTypeUnknown, fmt.Errorf("NOT operator requires boolean operand, got %s", operandType)
		}
		return FieldTypeBool, nil

	case "+", "-":
		if !operandType.IsNumeric() {
			return FieldTypeUnknown, fmt.Errorf("unary %s requires numeric operand, got %s", expr.Operator, operandType)
		}
		return operandType, nil

	default:
		return FieldTypeUnknown, fmt.Errorf("unsupported unary operator: %s", expr.Operator)
	}
}

// inferFunctionType infers the return type of a function call
func (tc *TypeChecker) inferFunctionType(call *ast.FunctionCall) (FieldType, error) {
	// For now, use a simple mapping based on function name
	// This will be expanded with a function registry in the future
	switch call.Name {
	// Aggregation functions
	case "count":
		return FieldTypeLong, nil
	case "sum":
		// Sum returns same type as argument (or double if mixed)
		if len(call.Arguments) > 0 {
			return tc.InferType(call.Arguments[0])
		}
		return FieldTypeDouble, nil
	case "avg", "mean":
		return FieldTypeDouble, nil
	case "min", "max":
		// Min/max return same type as argument
		if len(call.Arguments) > 0 {
			return tc.InferType(call.Arguments[0])
		}
		return FieldTypeUnknown, nil
	case "stddev", "variance", "percentile":
		return FieldTypeDouble, nil

	// String functions
	case "concat", "substring", "upper", "lower", "trim":
		return FieldTypeString, nil
	case "length":
		return FieldTypeInt, nil

	// Math functions
	case "abs":
		if len(call.Arguments) > 0 {
			return tc.InferType(call.Arguments[0])
		}
		return FieldTypeDouble, nil
	case "round", "ceil", "floor":
		return FieldTypeLong, nil
	case "sqrt", "pow", "exp", "log", "ln":
		return FieldTypeDouble, nil

	// Date functions
	case "now", "curdate", "curtime":
		return FieldTypeDate, nil
	case "year", "month", "day", "hour", "minute", "second":
		return FieldTypeInt, nil

	// Type conversion
	case "cast", "convert":
		// TODO: Infer based on target type in arguments
		return FieldTypeUnknown, nil

	default:
		return FieldTypeUnknown, fmt.Errorf("unknown function: %s", call.Name)
	}
}

// inferCaseExprType infers the type of a CASE expression
func (tc *TypeChecker) inferCaseExprType(caseExpr *ast.CaseExpression) (FieldType, error) {
	// All THEN expressions must have compatible types
	// WHEN conditions must be boolean
	var resultType FieldType

	for i, whenClause := range caseExpr.WhenClauses {
		// Check condition is boolean
		condType, err := tc.InferType(whenClause.Condition)
		if err != nil {
			return FieldTypeUnknown, err
		}
		if condType != FieldTypeBool {
			return FieldTypeUnknown, fmt.Errorf("WHEN condition must be boolean, got %s", condType)
		}

		// Infer result type from first THEN expression
		thenType, err := tc.InferType(whenClause.Result)
		if err != nil {
			return FieldTypeUnknown, err
		}

		if i == 0 {
			resultType = thenType
		} else if resultType != thenType {
			// Type mismatch - try to find compatible type
			resultType = tc.findCompatibleType(resultType, thenType)
			if resultType == FieldTypeUnknown {
				return FieldTypeUnknown, fmt.Errorf("incompatible types in CASE expression: %s and %s", resultType, thenType)
			}
		}
	}

	// Check ELSE clause if present
	if caseExpr.ElseResult != nil {
		elseType, err := tc.InferType(caseExpr.ElseResult)
		if err != nil {
			return FieldTypeUnknown, err
		}

		if resultType != elseType {
			resultType = tc.findCompatibleType(resultType, elseType)
			if resultType == FieldTypeUnknown {
				return FieldTypeUnknown, fmt.Errorf("incompatible types in CASE expression ELSE: %s and %s", resultType, elseType)
			}
		}
	}

	return resultType, nil
}

// inferArithmeticType infers the result type of arithmetic operations
func (tc *TypeChecker) inferArithmeticType(leftType, rightType FieldType, op string) (FieldType, error) {
	// Both operands must be numeric
	if !leftType.IsNumeric() {
		return FieldTypeUnknown, fmt.Errorf("left operand of %s must be numeric, got %s", op, leftType)
	}
	if !rightType.IsNumeric() {
		return FieldTypeUnknown, fmt.Errorf("right operand of %s must be numeric, got %s", op, rightType)
	}

	// Result type is the "wider" of the two types
	// Promotion order: int < long < float < double
	if leftType == FieldTypeDouble || rightType == FieldTypeDouble {
		return FieldTypeDouble, nil
	}
	if leftType == FieldTypeFloat || rightType == FieldTypeFloat {
		return FieldTypeFloat, nil
	}
	if leftType == FieldTypeLong || rightType == FieldTypeLong {
		return FieldTypeLong, nil
	}
	return FieldTypeInt, nil
}

// validateComparison checks if two types can be compared with the given operator
func (tc *TypeChecker) validateComparison(leftType, rightType FieldType, op string) error {
	// Allow NULL comparisons
	if leftType == FieldTypeUnknown || rightType == FieldTypeUnknown {
		return nil // NULL can be compared with anything
	}

	// Both types must be comparable
	if !leftType.IsComparable() {
		return fmt.Errorf("left operand of %s is not comparable: %s", op, leftType)
	}
	if !rightType.IsComparable() {
		return fmt.Errorf("right operand of %s is not comparable: %s", op, rightType)
	}

	// For equality/inequality, any comparable types can be compared
	if op == "=" || op == "!=" || op == "<>" {
		return nil
	}

	// For ordering comparisons (<, >, <=, >=), types must be compatible
	// Numeric types are compatible with each other
	if leftType.IsNumeric() && rightType.IsNumeric() {
		return nil
	}

	// String types are compatible with each other
	if leftType.IsString() && rightType.IsString() {
		return nil
	}

	// Same types are always compatible
	if leftType == rightType {
		return nil
	}

	return fmt.Errorf("incompatible types for %s: %s and %s", op, leftType, rightType)
}

// findCompatibleType finds a common type that both types can be coerced to
func (tc *TypeChecker) findCompatibleType(type1, type2 FieldType) FieldType {
	// Same type
	if type1 == type2 {
		return type1
	}

	// Numeric types can be widened
	if type1.IsNumeric() && type2.IsNumeric() {
		resultType, _ := tc.inferArithmeticType(type1, type2, "+") // Reuse arithmetic type promotion
		return resultType
	}

	// String types are compatible
	if type1.IsString() && type2.IsString() {
		return FieldTypeString
	}

	// No compatible type found
	return FieldTypeUnknown
}

// Helper functions

func isComparisonOp(op string) bool {
	switch op {
	case "=", "!=", "<>", "<", ">", "<=", ">=", "LIKE", "IN":
		return true
	default:
		return false
	}
}

func isLogicalOp(op string) bool {
	switch op {
	case "AND", "OR":
		return true
	default:
		return false
	}
}

func isArithmeticOp(op string) bool {
	switch op {
	case "+", "-", "*", "/", "%":
		return true
	default:
		return false
	}
}
