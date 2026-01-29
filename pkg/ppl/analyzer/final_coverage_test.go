// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package analyzer

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================================
// Final Coverage Tests to Reach 85%
// =====================================================================

func TestTypeChecker_InferUnaryExprType_AllBranches(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("UnaryPlusOperator", func(t *testing.T) {
		expr := &ast.UnaryExpression{
			Operator: "+",
			Operand:  &ast.Literal{Value: 42, LiteralTyp: ast.LiteralTypeInt},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("UnaryMinusOnFloat", func(t *testing.T) {
		expr := &ast.UnaryExpression{
			Operator: "-",
			Operand:  &ast.FieldReference{Name: "latency"},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("UnaryNotOnExpression", func(t *testing.T) {
		expr := &ast.UnaryExpression{
			Operator: "NOT",
			Operand: &ast.BinaryExpression{
				Left:     &ast.Literal{Value: true, LiteralTyp: ast.LiteralTypeBool},
				Operator: "AND",
				Right:    &ast.Literal{Value: false, LiteralTyp: ast.LiteralTypeBool},
			},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("UnaryInvalidOperator", func(t *testing.T) {
		expr := &ast.UnaryExpression{
			Operator: "INVALID",
			Operand:  &ast.Literal{Value: 42, LiteralTyp: ast.LiteralTypeInt},
		}
		_, err := typeChecker.InferType(expr)
		assert.Error(t, err)
	})
}

func TestTypeChecker_InferFunctionType_MoreFunctions(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("MeanFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "mean",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("SubstringFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "substring",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "host"},
				&ast.Literal{Value: 0, LiteralTyp: ast.LiteralTypeInt},
				&ast.Literal{Value: 5, LiteralTyp: ast.LiteralTypeInt},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, fieldType)
	})

	t.Run("TrimFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "trim",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "host"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, fieldType)
	})

	t.Run("SqrtFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "sqrt",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("PowFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "pow",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
				&ast.Literal{Value: 2, LiteralTyp: ast.LiteralTypeInt},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("ExpFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "exp",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("LogFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "log",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("LnFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "ln",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("CurdateFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "curdate",
			Arguments: []ast.Expression{},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDate, fieldType)
	})

	t.Run("CurtimeFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "curtime",
			Arguments: []ast.Expression{},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDate, fieldType)
	})

	t.Run("YearFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "year",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "timestamp"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("MonthFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "month",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "timestamp"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("DayFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "day",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "timestamp"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("HourFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "hour",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "timestamp"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("MinuteFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "minute",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "timestamp"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("SecondFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "second",
			Arguments: []ast.Expression{&ast.FieldReference{Name: "timestamp"}},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("CastFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "cast",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "status"},
				&ast.Literal{Value: "string", LiteralTyp: ast.LiteralTypeString},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeUnknown, fieldType)
	})

	t.Run("ConvertFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "convert",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "status"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeUnknown, fieldType)
	})
}

func TestAnalyzer_SortCommand_MoreCases(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("SortEmptyFields", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.SortCommand{
					SortKeys: []*ast.SortKey{},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})

	t.Run("SortWithFunctionCall", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.SortCommand{
					SortKeys: []*ast.SortKey{
						{
							Field: &ast.FunctionCall{
								Name:      "abs",
								Arguments: []ast.Expression{&ast.FieldReference{Name: "status"}},
							},
							Descending: true,
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})
}

func TestAnalyzer_CaseExpression_MoreCases(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("CaseEmptyWhenClauses", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "result",
							Expression: &ast.CaseExpression{
								WhenClauses: []*ast.WhenClause{},
								ElseResult:  &ast.Literal{Value: "default", LiteralTyp: ast.LiteralTypeString},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})

	t.Run("CaseWithInvalidConditionType", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "result",
							Expression: &ast.CaseExpression{
								WhenClauses: []*ast.WhenClause{
									{
										Condition: &ast.FieldReference{Name: "nonexistent"},
										Result:    &ast.Literal{Value: "result", LiteralTyp: ast.LiteralTypeString},
									},
								},
								ElseResult: &ast.Literal{Value: "else", LiteralTyp: ast.LiteralTypeString},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestAnalyzer_UnaryExpression_InvalidCases(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("UnaryInvalidOperand", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "result",
							Expression: &ast.UnaryExpression{
								Operator: "-",
								Operand:  &ast.FieldReference{Name: "nonexistent"},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})

	t.Run("UnaryNegateString", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "result",
							Expression: &ast.UnaryExpression{
								Operator: "-",
								Operand:  &ast.FieldReference{Name: "host"},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestTypeChecker_ValidateComparison_MoreEdgeCases(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("CompareNullWithNull", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: nil, LiteralTyp: ast.LiteralTypeNull},
			Operator: "=",
			Right:    &ast.Literal{Value: nil, LiteralTyp: ast.LiteralTypeNull},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("CompareIntWithString", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.FieldReference{Name: "host"},
		}
		// Comparison still returns bool even with different types
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})
}
