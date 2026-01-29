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
// Edge Case Tests for Maximum Coverage
// =====================================================================

func TestAnalyzer_UnaryExpression_AllOperators(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("NOTWithComplexCondition", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.WhereCommand{
					Condition: &ast.UnaryExpression{
						Operator: "NOT",
						Operand: &ast.BinaryExpression{
							Left: &ast.BinaryExpression{
								Left:     &ast.FieldReference{Name: "status"},
								Operator: ">=",
								Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
							},
							Operator: "AND",
							Right: &ast.BinaryExpression{
								Left:     &ast.FieldReference{Name: "status"},
								Operator: "<",
								Right:    &ast.Literal{Value: 300, LiteralTyp: ast.LiteralTypeInt},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("NegateFieldReference", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "neg_status",
							Expression: &ast.UnaryExpression{
								Operator: "-",
								Operand:  &ast.FieldReference{Name: "status"},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("NegateComplexExpression", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "result",
							Expression: &ast.UnaryExpression{
								Operator: "-",
								Operand: &ast.BinaryExpression{
									Left:     &ast.FieldReference{Name: "status"},
									Operator: "+",
									Right:    &ast.Literal{Value: 100, LiteralTyp: ast.LiteralTypeInt},
								},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})
}

func TestAnalyzer_CaseExpression_EdgeCases(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("CaseInWhere", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.WhereCommand{
					Condition: &ast.BinaryExpression{
						Left: &ast.CaseExpression{
							WhenClauses: []*ast.WhenClause{
								{
									Condition: &ast.BinaryExpression{
										Left:     &ast.FieldReference{Name: "status"},
										Operator: ">=",
										Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
									},
									Result: &ast.Literal{Value: "error", LiteralTyp: ast.LiteralTypeString},
								},
							},
							ElseResult: &ast.Literal{Value: "ok", LiteralTyp: ast.LiteralTypeString},
						},
						Operator: "=",
						Right:    &ast.Literal{Value: "error", LiteralTyp: ast.LiteralTypeString},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("NestedCaseExpression", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "category",
							Expression: &ast.CaseExpression{
								WhenClauses: []*ast.WhenClause{
									{
										Condition: &ast.BinaryExpression{
											Left:     &ast.FieldReference{Name: "status"},
											Operator: "<",
											Right:    &ast.Literal{Value: 300, LiteralTyp: ast.LiteralTypeInt},
										},
										Result: &ast.CaseExpression{
											WhenClauses: []*ast.WhenClause{
												{
													Condition: &ast.BinaryExpression{
														Left:     &ast.FieldReference{Name: "status"},
														Operator: "<",
														Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
													},
													Result: &ast.Literal{Value: "info", LiteralTyp: ast.LiteralTypeString},
												},
											},
											ElseResult: &ast.Literal{Value: "success", LiteralTyp: ast.LiteralTypeString},
										},
									},
								},
								ElseResult: &ast.Literal{Value: "other", LiteralTyp: ast.LiteralTypeString},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})
}

func TestAnalyzer_ListLiteral_EdgeCases(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("EmptyListLiteral", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.WhereCommand{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: "IN",
						Right: &ast.ListLiteral{
							Values: []ast.Expression{},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		// ListLiteral not supported in type checker
		assert.Error(t, err)
	})

	t.Run("ListWithSingleElement", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.WhereCommand{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: "IN",
						Right: &ast.ListLiteral{
							Values: []ast.Expression{
								&ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err) // IN with list not supported by type checker
	})

	t.Run("ListWithMixedNumericTypes", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.WhereCommand{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "latency"},
						Operator: "IN",
						Right: &ast.ListLiteral{
							Values: []ast.Expression{
								&ast.Literal{Value: 1, LiteralTyp: ast.LiteralTypeInt},
								&ast.Literal{Value: 1.5, LiteralTyp: ast.LiteralTypeFloat},
								&ast.Literal{Value: 2, LiteralTyp: ast.LiteralTypeInt},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err) // Mixed types not supported by type checker
	})
}

func TestAnalyzer_StatsCommand_AllEdgeCases(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("StatsWithNoAggregations", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.StatsCommand{
					Aggregations: []*ast.Aggregation{},
					GroupBy: []ast.Expression{
						&ast.FieldReference{Name: "host"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})

	t.Run("StatsGroupByWithInvalidField", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.StatsCommand{
					Aggregations: []*ast.Aggregation{
						{
							Func:  &ast.FunctionCall{Name: "count", Arguments: []ast.Expression{}},
							Alias: "total",
						},
					},
					GroupBy: []ast.Expression{
						&ast.FieldReference{Name: "nonexistent"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})

	t.Run("StatsWithDuplicateAlias", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.StatsCommand{
					Aggregations: []*ast.Aggregation{
						{
							Func:  &ast.FunctionCall{Name: "count", Arguments: []ast.Expression{}},
							Alias: "same_alias",
						},
						{
							Func:  &ast.FunctionCall{Name: "sum", Arguments: []ast.Expression{&ast.FieldReference{Name: "status"}}},
							Alias: "same_alias",
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestAnalyzer_EvalCommand_AllEdgeCases(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("EvalWithInvalidExpression", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "result",
							Expression: &ast.BinaryExpression{
								Left:     &ast.FieldReference{Name: "nonexistent"},
								Operator: "+",
								Right:    &ast.Literal{Value: 1, LiteralTyp: ast.LiteralTypeInt},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})

	t.Run("EvalRedefiningExistingField", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "status",
							Expression: &ast.BinaryExpression{
								Left:     &ast.FieldReference{Name: "status"},
								Operator: "+",
								Right:    &ast.Literal{Value: 1, LiteralTyp: ast.LiteralTypeInt},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err) // Redefining is allowed
	})

	t.Run("EvalDependentAssignments", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "first",
							Expression: &ast.BinaryExpression{
								Left:     &ast.FieldReference{Name: "status"},
								Operator: "+",
								Right:    &ast.Literal{Value: 100, LiteralTyp: ast.LiteralTypeInt},
							},
						},
						{
							Field: "second",
							Expression: &ast.BinaryExpression{
								Left:     &ast.FieldReference{Name: "first"},
								Operator: "*",
								Right:    &ast.Literal{Value: 2, LiteralTyp: ast.LiteralTypeInt},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
		assert.True(t, analyzer.scope.Has("first"))
		assert.True(t, analyzer.scope.Has("second"))
	})
}

func TestAnalyzer_TimechartCommand_AllEdgeCases(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("TimechartWithGroupBy", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TimechartCommand{
					Span: &ast.TimeSpan{Value: 1, Unit: "h"},
					Aggregations: []*ast.Aggregation{
						{
							Func:  &ast.FunctionCall{Name: "count", Arguments: []ast.Expression{}},
							Alias: "total",
						},
					},
					GroupBy: []ast.Expression{
						&ast.FieldReference{Name: "level"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("TimechartInvalidGroupBy", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TimechartCommand{
					Span: &ast.TimeSpan{Value: 1, Unit: "m"},
					Aggregations: []*ast.Aggregation{
						{
							Func:  &ast.FunctionCall{Name: "count", Arguments: []ast.Expression{}},
							Alias: "total",
						},
					},
					GroupBy: []ast.Expression{
						&ast.FieldReference{Name: "nonexistent"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})

	t.Run("TimechartNoAggregations", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TimechartCommand{
					Span:         &ast.TimeSpan{Value: 5, Unit: "m"},
					Aggregations: []*ast.Aggregation{},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})

	t.Run("TimechartDuplicateAlias", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TimechartCommand{
					Span: &ast.TimeSpan{Value: 1, Unit: "h"},
					Aggregations: []*ast.Aggregation{
						{
							Func:  &ast.FunctionCall{Name: "count", Arguments: []ast.Expression{}},
							Alias: "dup",
						},
						{
							Func:  &ast.FunctionCall{Name: "avg", Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}}},
							Alias: "dup",
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestTypeChecker_ValidateComparison_AllCases(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("CompareFloatWithInt", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "latency"},
			Operator: ">",
			Right:    &ast.Literal{Value: 100, LiteralTyp: ast.LiteralTypeInt},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("CompareWithNull", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: nil, LiteralTyp: ast.LiteralTypeNull},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("CompareStrings", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "host"},
			Operator: "=",
			Right:    &ast.FieldReference{Name: "level"},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("CompareWithExpression", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: "+",
				Right:    &ast.Literal{Value: 100, LiteralTyp: ast.LiteralTypeInt},
			},
			Operator: ">",
			Right:    &ast.Literal{Value: 300, LiteralTyp: ast.LiteralTypeInt},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})
}

func TestSchema_GetNestedField_AllCases(t *testing.T) {
	// Create nested schema
	field := &Field{
		Name: "metadata",
		Type: FieldTypeObject,
		Fields: map[string]*Field{
			"level1": {
				Name: "level1",
				Type: FieldTypeObject,
				Fields: map[string]*Field{
					"level2": {
						Name: "level2",
						Type: FieldTypeObject,
						Fields: map[string]*Field{
							"value": {
								Name: "value",
								Type: FieldTypeString,
							},
						},
					},
				},
			},
		},
	}

	t.Run("GetDeeplyNestedField", func(t *testing.T) {
		nested, err := field.GetNestedField("level1.level2.value")
		require.NoError(t, err)
		assert.Equal(t, "value", nested.Name)
		assert.Equal(t, FieldTypeString, nested.Type)
	})

	t.Run("GetPartialPath", func(t *testing.T) {
		nested, err := field.GetNestedField("level1.level2")
		require.NoError(t, err)
		assert.Equal(t, "level2", nested.Name)
		assert.Equal(t, FieldTypeObject, nested.Type)
	})

	t.Run("GetInvalidNestedPath", func(t *testing.T) {
		_, err := field.GetNestedField("level1.nonexistent")
		assert.Error(t, err)
	})

	t.Run("GetEmptyPath", func(t *testing.T) {
		_, err := field.GetNestedField("")
		assert.Error(t, err)
	})
}
