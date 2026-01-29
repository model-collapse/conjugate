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
// Additional Function Tests for Coverage
// =====================================================================

func TestTypeChecker_InferFunctionType_Extended(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("StddevFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "stddev",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("VarianceFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "variance",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("UpperFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "upper",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "host"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, fieldType)
	})

	t.Run("LowerFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "lower",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "host"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, fieldType)
	})

	t.Run("LengthFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "length",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "host"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("ConcatFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "concat",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "host"},
				&ast.Literal{Value: "-", LiteralTyp: ast.LiteralTypeString},
				&ast.FieldReference{Name: "level"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, fieldType)
	})

	t.Run("AbsFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "abs",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "status"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("RoundFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "round",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("FloorFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "floor",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("CeilFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "ceil",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("PercentileFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "percentile",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
				&ast.Literal{Value: 95, LiteralTyp: ast.LiteralTypeInt},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("CountNoArgs", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "count",
			Arguments: []ast.Expression{},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("FunctionWithNoArguments", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "now",
			Arguments: []ast.Expression{},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		// now() returns date type
		assert.Equal(t, FieldTypeDate, fieldType)
	})
}

func TestAnalyzer_TimechartCommand_Extended(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("TimechartMultipleAggregations", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TimechartCommand{
					Span: &ast.TimeSpan{Value: 5, Unit: "m"},
					Aggregations: []*ast.Aggregation{
						{
							Func:  &ast.FunctionCall{Name: "count", Arguments: []ast.Expression{}},
							Alias: "total_count",
						},
						{
							Func:  &ast.FunctionCall{Name: "avg", Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}}},
							Alias: "avg_latency",
						},
						{
							Func:  &ast.FunctionCall{Name: "max", Arguments: []ast.Expression{&ast.FieldReference{Name: "status"}}},
							Alias: "max_status",
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
		assert.True(t, analyzer.scope.Has("_time"))
		assert.True(t, analyzer.scope.Has("total_count"))
		assert.True(t, analyzer.scope.Has("avg_latency"))
		assert.True(t, analyzer.scope.Has("max_status"))
	})

	t.Run("TimechartWithInvalidAggregation", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TimechartCommand{
					Span: &ast.TimeSpan{Value: 1, Unit: "h"},
					Aggregations: []*ast.Aggregation{
						{
							Func:  &ast.FunctionCall{Name: "avg", Arguments: []ast.Expression{&ast.FieldReference{Name: "nonexistent"}}},
							Alias: "avg_val",
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("TimechartWithComplexSpan", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TimechartCommand{
					Span: &ast.TimeSpan{Value: 30, Unit: "s"},
					Aggregations: []*ast.Aggregation{
						{
							Func:  &ast.FunctionCall{Name: "sum", Arguments: []ast.Expression{&ast.FieldReference{Name: "status"}}},
							Alias: "total_status",
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})
}

func TestAnalyzer_SortCommand_Extended(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("SortMultipleFields", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.SortCommand{
					SortKeys: []*ast.SortKey{
						{Field: &ast.FieldReference{Name: "status"}, Descending: false},
						{Field: &ast.FieldReference{Name: "latency"}, Descending: true},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("SortByExpression", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.SortCommand{
					SortKeys: []*ast.SortKey{
						{
							Field: &ast.BinaryExpression{
								Left:     &ast.FieldReference{Name: "status"},
								Operator: "+",
								Right:    &ast.Literal{Value: 100, LiteralTyp: ast.LiteralTypeInt},
							},
							Descending: false,
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("SortByInvalidField", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.SortCommand{
					SortKeys: []*ast.SortKey{
						{Field: &ast.FieldReference{Name: "nonexistent"}, Descending: false},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestAnalyzer_FieldsCommand_Extended(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("FieldsWithWildcard", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.FieldsCommand{
					Fields: []ast.Expression{
						&ast.FieldReference{Name: "status"},
						&ast.FieldReference{Name: "host"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("FieldsWithInvalidField", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.FieldsCommand{
					Fields: []ast.Expression{
						&ast.FieldReference{Name: "nonexistent"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})

	t.Run("FieldsWithExpressions", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.FieldsCommand{
					Fields: []ast.Expression{
						&ast.BinaryExpression{
							Left:     &ast.FieldReference{Name: "status"},
							Operator: "+",
							Right:    &ast.Literal{Value: 1, LiteralTyp: ast.LiteralTypeInt},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})
}

func TestAnalyzer_StatsCommand_Extended(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("StatsMultipleGroupBy", func(t *testing.T) {
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
						&ast.FieldReference{Name: "host"},
						&ast.FieldReference{Name: "level"},
						&ast.FieldReference{Name: "status"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("StatsWithComplexAggregation", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.StatsCommand{
					Aggregations: []*ast.Aggregation{
						{
							Func: &ast.FunctionCall{
								Name: "avg",
								Arguments: []ast.Expression{
									&ast.BinaryExpression{
										Left:     &ast.FieldReference{Name: "latency"},
										Operator: "*",
										Right:    &ast.Literal{Value: 1000.0, LiteralTyp: ast.LiteralTypeFloat},
									},
								},
							},
							Alias: "avg_ms",
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})
}

func TestSchema_GetField_ErrorCases(t *testing.T) {
	schema := createTestSchema()

	t.Run("GetNonexistentField", func(t *testing.T) {
		_, err := schema.GetField("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("GetExistingField", func(t *testing.T) {
		field, err := schema.GetField("status")
		require.NoError(t, err)
		assert.Equal(t, "status", field.Name)
		assert.Equal(t, FieldTypeInt, field.Type)
	})

	t.Run("GetMultipleFields", func(t *testing.T) {
		fields := []string{"status", "host", "latency", "level"}
		for _, fieldName := range fields {
			field, err := schema.GetField(fieldName)
			require.NoError(t, err, "Field %s should exist", fieldName)
			assert.Equal(t, fieldName, field.Name)
		}
	})
}

func TestTypeChecker_ArithmeticOperators_Extended(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	operators := []string{"+", "-", "*", "/", "%"}

	for _, op := range operators {
		t.Run("Operator_"+op, func(t *testing.T) {
			expr := &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: op,
				Right:    &ast.FieldReference{Name: "status"},
			}
			_, err := typeChecker.InferType(expr)
			assert.NoError(t, err)
		})
	}

	t.Run("PowerOperator", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: 2, LiteralTyp: ast.LiteralTypeInt},
			Operator: "^",
			Right:    &ast.Literal{Value: 3, LiteralTyp: ast.LiteralTypeInt},
		}
		_, err := typeChecker.InferType(expr)
		// Power operator may not be supported
		if err != nil {
			assert.Contains(t, err.Error(), "unsupported")
		}
	})
}

func TestTypeChecker_CaseExpression_Extended(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("CaseWithoutElse", func(t *testing.T) {
		caseExpr := &ast.CaseExpression{
			WhenClauses: []*ast.WhenClause{
				{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: ">=",
						Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
					},
					Result: &ast.Literal{Value: "success", LiteralTyp: ast.LiteralTypeString},
				},
			},
			ElseResult: nil,
		}

		fieldType, err := typeChecker.InferType(caseExpr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, fieldType)
	})

	t.Run("CaseWithMultipleConditions", func(t *testing.T) {
		caseExpr := &ast.CaseExpression{
			WhenClauses: []*ast.WhenClause{
				{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: "<",
						Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
					},
					Result: &ast.Literal{Value: 1, LiteralTyp: ast.LiteralTypeInt},
				},
				{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: "<",
						Right:    &ast.Literal{Value: 300, LiteralTyp: ast.LiteralTypeInt},
					},
					Result: &ast.Literal{Value: 2, LiteralTyp: ast.LiteralTypeInt},
				},
				{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: "<",
						Right:    &ast.Literal{Value: 400, LiteralTyp: ast.LiteralTypeInt},
					},
					Result: &ast.Literal{Value: 3, LiteralTyp: ast.LiteralTypeInt},
				},
			},
			ElseResult: &ast.Literal{Value: 4, LiteralTyp: ast.LiteralTypeInt},
		}

		fieldType, err := typeChecker.InferType(caseExpr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("CaseWithNonBooleanCondition", func(t *testing.T) {
		caseExpr := &ast.CaseExpression{
			WhenClauses: []*ast.WhenClause{
				{
					Condition: &ast.Literal{Value: 123, LiteralTyp: ast.LiteralTypeInt},
					Result:    &ast.Literal{Value: "result", LiteralTyp: ast.LiteralTypeString},
				},
			},
			ElseResult: &ast.Literal{Value: "else", LiteralTyp: ast.LiteralTypeString},
		}

		_, err := typeChecker.InferType(caseExpr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be boolean")
	})

	t.Run("CaseWithInvalidThenType", func(t *testing.T) {
		caseExpr := &ast.CaseExpression{
			WhenClauses: []*ast.WhenClause{
				{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: "=",
						Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
					},
					Result: &ast.FieldReference{Name: "nonexistent"},
				},
			},
			ElseResult: &ast.Literal{Value: "else", LiteralTyp: ast.LiteralTypeString},
		}

		_, err := typeChecker.InferType(caseExpr)
		assert.Error(t, err)
	})
}

func TestTypeChecker_LiteralType_EdgeCases(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("NullLiteral", func(t *testing.T) {
		lit := &ast.Literal{Value: nil, LiteralTyp: ast.LiteralTypeNull}
		fieldType, err := typeChecker.InferType(lit)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeUnknown, fieldType)
	})

	t.Run("IntegerLiteral", func(t *testing.T) {
		lit := &ast.Literal{Value: 42, LiteralTyp: ast.LiteralTypeInt}
		fieldType, err := typeChecker.InferType(lit)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("FloatLiteral", func(t *testing.T) {
		lit := &ast.Literal{Value: 3.14, LiteralTyp: ast.LiteralTypeFloat}
		fieldType, err := typeChecker.InferType(lit)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("BooleanLiteral", func(t *testing.T) {
		lit := &ast.Literal{Value: true, LiteralTyp: ast.LiteralTypeBool}
		fieldType, err := typeChecker.InferType(lit)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("StringLiteral", func(t *testing.T) {
		lit := &ast.Literal{Value: "test", LiteralTyp: ast.LiteralTypeString}
		fieldType, err := typeChecker.InferType(lit)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, fieldType)
	})
}

func TestTypeChecker_ValidateComparison_Extended(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	comparisonOps := []string{"=", "!=", "<", "<=", ">", ">="}

	for _, op := range comparisonOps {
		t.Run("Comparison_"+op, func(t *testing.T) {
			expr := &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: op,
				Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
			}
			fieldType, err := typeChecker.InferType(expr)
			require.NoError(t, err)
			assert.Equal(t, FieldTypeBool, fieldType)
		})
	}

	t.Run("CompareStringWithNumber", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "host"},
			Operator: "<",
			Right:    &ast.Literal{Value: 100, LiteralTyp: ast.LiteralTypeInt},
		}
		_, err := typeChecker.InferType(expr)
		assert.Error(t, err)
	})
}
