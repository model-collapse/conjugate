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
// Tier 1 Command Analyzer Tests
// =====================================================================

func TestAnalyzer_TopCommand(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("ValidTop", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TopCommand{
					Fields: []ast.Expression{
						&ast.FieldReference{Name: "host"},
					},
					Limit:       10,
					ShowCount:   true,
					ShowPercent: true,
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)

		// Scope should have count and percent fields
		assert.True(t, analyzer.scope.Has("count"))
		assert.True(t, analyzer.scope.Has("percent"))
	})

	t.Run("TopWithGroupBy", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TopCommand{
					Fields:  []ast.Expression{&ast.FieldReference{Name: "status"}},
					Limit:   5,
					GroupBy: []ast.Expression{&ast.FieldReference{Name: "level"}},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("TopWithNoFields", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TopCommand{
					Fields: []ast.Expression{},
					Limit:  10,
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one field")
	})
}

func TestAnalyzer_RareCommand(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("ValidRare", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.RareCommand{
					Fields:      []ast.Expression{&ast.FieldReference{Name: "status"}},
					Limit:       10,
					ShowCount:   true,
					ShowPercent: false,
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)

		// Should have count field
		assert.True(t, analyzer.scope.Has("count"))
	})

	t.Run("RareWithNoFields", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.RareCommand{
					Fields: []ast.Expression{},
					Limit:  5,
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestAnalyzer_DedupCommand(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("ValidDedup", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.DedupCommand{
					Fields: []ast.Expression{&ast.FieldReference{Name: "host"}},
					Count:  2,
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("DedupWithNegativeCount", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.DedupCommand{
					Fields: []ast.Expression{&ast.FieldReference{Name: "host"}},
					Count:  -1,
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")
	})

	t.Run("DedupWithNoFields", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.DedupCommand{
					Fields: []ast.Expression{},
					Count:  1,
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestAnalyzer_EvalCommand(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("ValidEval", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "response_time",
							Expression: &ast.BinaryExpression{
								Left:     &ast.FieldReference{Name: "latency"},
								Operator: "*",
								Right:    &ast.Literal{Value: 1000.0, LiteralTyp: ast.LiteralTypeFloat},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)

		// New field should be in scope
		assert.True(t, analyzer.scope.Has("response_time"))
	})

	t.Run("EvalMultipleAssignments", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field:      "is_error",
							Expression: &ast.BinaryExpression{Left: &ast.FieldReference{Name: "status"}, Operator: ">=", Right: &ast.Literal{Value: 400, LiteralTyp: ast.LiteralTypeInt}},
						},
						{
							Field:      "is_critical",
							Expression: &ast.BinaryExpression{Left: &ast.FieldReference{Name: "status"}, Operator: ">=", Right: &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt}},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
		assert.True(t, analyzer.scope.Has("is_error"))
		assert.True(t, analyzer.scope.Has("is_critical"))
	})

	t.Run("EvalWithNoAssignments", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestAnalyzer_RenameCommand(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("ValidRename", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.RenameCommand{
					Assignments: []*ast.RenameAssignment{
						{OldName: "latency", NewName: "response_time"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)

		// New name should be in scope
		assert.True(t, analyzer.scope.Has("response_time"))
	})

	t.Run("RenameNonexistentField", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.RenameCommand{
					Assignments: []*ast.RenameAssignment{
						{OldName: "nonexistent", NewName: "new_name"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("RenameWithNoAssignments", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.RenameCommand{
					Assignments: []*ast.RenameAssignment{},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestAnalyzer_BinCommand(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("BinWithSpan", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.BinCommand{
					Field: &ast.FieldReference{Name: "timestamp"},
					Span:  &ast.TimeSpan{Value: 1, Unit: "h"},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("BinWithBins", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.BinCommand{
					Field: &ast.FieldReference{Name: "latency"},
					Bins:  10,
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("BinWithInvalidFieldType", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.BinCommand{
					Field: &ast.FieldReference{Name: "level"}, // String field, not numeric/date
					Span:  &ast.TimeSpan{Value: 1, Unit: "h"},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be numeric or date")
	})

	t.Run("BinWithNoSpanOrBins", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.BinCommand{
					Field: &ast.FieldReference{Name: "timestamp"},
					Span:  nil,
					Bins:  0,
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "requires either span or bins")
	})
}

func TestAnalyzer_TimechartCommand(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("ValidTimechart", func(t *testing.T) {
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
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)

		// Should have _time and total fields
		assert.True(t, analyzer.scope.Has("_time"))
		assert.True(t, analyzer.scope.Has("total"))
	})

	t.Run("TimechartWithNoSpan", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TimechartCommand{
					Span: nil,
					Aggregations: []*ast.Aggregation{
						{Func: &ast.FunctionCall{Name: "count", Arguments: []ast.Expression{}}, Alias: "total"},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "requires a span")
	})

	t.Run("TimechartWithNoAggregations", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.TimechartCommand{
					Span:         &ast.TimeSpan{Value: 1, Unit: "h"},
					Aggregations: []*ast.Aggregation{},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
	})
}

func TestAnalyzer_DescribeCommand(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("ValidDescribe", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.DescribeCommand{Source: "logs"},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("DescribeWithNoSource", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.DescribeCommand{Source: ""},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "requires a source")
	})
}

// =====================================================================
// Expression Analyzer Tests
// =====================================================================

func TestAnalyzer_UnaryExpression(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("UnaryNOT", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.WhereCommand{
					Condition: &ast.UnaryExpression{
						Operator: "NOT",
						Operand: &ast.BinaryExpression{
							Left:     &ast.FieldReference{Name: "status"},
							Operator: "=",
							Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.NoError(t, err)
	})

	t.Run("UnaryNegate", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "neg_latency",
							Expression: &ast.UnaryExpression{
								Operator: "-",
								Operand:  &ast.FieldReference{Name: "latency"},
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

func TestAnalyzer_CaseExpression(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("ValidCaseExpression", func(t *testing.T) {
		query := &ast.Query{
			Commands: []ast.Command{
				&ast.SearchCommand{Source: "logs"},
				&ast.EvalCommand{
					Assignments: []*ast.EvalAssignment{
						{
							Field: "severity",
							Expression: &ast.CaseExpression{
								WhenClauses: []*ast.WhenClause{
									{
										Condition: &ast.BinaryExpression{Left: &ast.FieldReference{Name: "status"}, Operator: ">=", Right: &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt}},
										Result:    &ast.Literal{Value: "critical", LiteralTyp: ast.LiteralTypeString},
									},
									{
										Condition: &ast.BinaryExpression{Left: &ast.FieldReference{Name: "status"}, Operator: ">=", Right: &ast.Literal{Value: 400, LiteralTyp: ast.LiteralTypeInt}},
										Result:    &ast.Literal{Value: "error", LiteralTyp: ast.LiteralTypeString},
									},
								},
								ElseResult: &ast.Literal{Value: "normal", LiteralTyp: ast.LiteralTypeString},
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

func TestAnalyzer_ListLiteral(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("ListLiteralNotSupportedInTypeChecker", func(t *testing.T) {
		// ListLiteral is validated by analyzer but not supported by type checker
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
								&ast.Literal{Value: 201, LiteralTyp: ast.LiteralTypeInt},
								&ast.Literal{Value: 204, LiteralTyp: ast.LiteralTypeInt},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported expression type")
	})

	t.Run("ListWithIncompatibleTypes", func(t *testing.T) {
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
								&ast.Literal{Value: "error", LiteralTyp: ast.LiteralTypeString},
							},
						},
					},
				},
			},
		}

		err := analyzer.Analyze(query)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "incompatible type")
	})
}

// =====================================================================
// Analyzer Getter Tests
// =====================================================================

func TestAnalyzer_Getters(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	t.Run("GetScope", func(t *testing.T) {
		scope := analyzer.GetScope()
		require.NotNil(t, scope)
	})

	t.Run("GetSchema", func(t *testing.T) {
		retrievedSchema := analyzer.GetSchema()
		require.NotNil(t, retrievedSchema)
		assert.Equal(t, "logs", retrievedSchema.Source)
	})

	t.Run("GetTypeChecker", func(t *testing.T) {
		typeChecker := analyzer.GetTypeChecker()
		require.NotNil(t, typeChecker)
	})
}
