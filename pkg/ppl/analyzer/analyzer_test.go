// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package analyzer

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestSchema() *Schema {
	schema := NewSchema("logs")
	schema.AddField("status", FieldTypeInt)
	schema.AddField("host", FieldTypeString)
	schema.AddField("timestamp", FieldTypeDate)
	schema.AddField("latency", FieldTypeDouble)
	schema.AddField("level", FieldTypeString)
	schema.AddField("message", FieldTypeText)
	return schema
}

func TestAnalyzer_SearchCommand(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
		},
	}

	err := analyzer.Analyze(query)
	assert.NoError(t, err)
}

func TestAnalyzer_WhereCommand_Valid(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.WhereCommand{
				Condition: &ast.BinaryExpression{
					Left:     &ast.FieldReference{Name: "status"},
					Operator: "=",
					Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
				},
			},
		},
	}

	err := analyzer.Analyze(query)
	assert.NoError(t, err)
}

func TestAnalyzer_WhereCommand_NonBooleanCondition(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.WhereCommand{
				Condition: &ast.FieldReference{Name: "status"}, // Not boolean
			},
		},
	}

	err := analyzer.Analyze(query)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be boolean")
}

func TestAnalyzer_WhereCommand_UnknownField(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.WhereCommand{
				Condition: &ast.BinaryExpression{
					Left:     &ast.FieldReference{Name: "unknown_field"},
					Operator: "=",
					Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
				},
			},
		},
	}

	err := analyzer.Analyze(query)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestAnalyzer_FieldsCommand_Valid(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

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
}

func TestAnalyzer_StatsCommand_Valid(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.StatsCommand{
				Aggregations: []*ast.Aggregation{
					{
						Func: &ast.FunctionCall{Name: "count"},
						Alias:      "total",
					},
				},
				GroupBy: []ast.Expression{
					&ast.FieldReference{Name: "host"},
				},
			},
		},
	}

	err := analyzer.Analyze(query)
	assert.NoError(t, err)
}

func TestAnalyzer_SortCommand_Valid(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.SortCommand{
				SortKeys: []*ast.SortKey{
					{
						Field:      &ast.FieldReference{Name: "timestamp"},
						Descending: true,
					},
				},
			},
		},
	}

	err := analyzer.Analyze(query)
	assert.NoError(t, err)
}

func TestAnalyzer_HeadCommand_Valid(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.HeadCommand{Count: 10},
		},
	}

	err := analyzer.Analyze(query)
	assert.NoError(t, err)
}

func TestAnalyzer_HeadCommand_InvalidCount(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.HeadCommand{Count: -1},
		},
	}

	err := analyzer.Analyze(query)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "positive count")
}

func TestAnalyzer_TypeMismatch_Comparison(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.WhereCommand{
				Condition: &ast.BinaryExpression{
					Left:     &ast.FieldReference{Name: "status"}, // int
					Operator: ">",
					Right:    &ast.Literal{Value: "abc", LiteralTyp: ast.LiteralTypeString}, // string
				},
			},
		},
	}

	err := analyzer.Analyze(query)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incompatible")
}

func TestAnalyzer_ArithmeticExpression(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.WhereCommand{
				Condition: &ast.BinaryExpression{
					Left: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "latency"},
						Operator: "*",
						Right:    &ast.Literal{Value: 2.0, LiteralTyp: ast.LiteralTypeFloat},
					},
					Operator: ">",
					Right:    &ast.Literal{Value: 100.0, LiteralTyp: ast.LiteralTypeFloat},
				},
			},
		},
	}

	err := analyzer.Analyze(query)
	assert.NoError(t, err)
}

func TestAnalyzer_FunctionCall_Count(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.StatsCommand{
				Aggregations: []*ast.Aggregation{
					{
						Func: &ast.FunctionCall{Name: "count"},
					},
				},
			},
		},
	}

	err := analyzer.Analyze(query)
	assert.NoError(t, err)
}

func TestAnalyzer_FunctionCall_Avg(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.StatsCommand{
				Aggregations: []*ast.Aggregation{
					{
						Func: &ast.FunctionCall{
							Name: "avg",
							Arguments: []ast.Expression{
								&ast.FieldReference{Name: "latency"},
							},
						},
					},
				},
			},
		},
	}

	err := analyzer.Analyze(query)
	assert.NoError(t, err)
}

func TestAnalyzer_ComplexPipeline(t *testing.T) {
	schema := createTestSchema()
	analyzer := NewAnalyzer(schema)

	// source=logs | where status=500 | stats count() as total by host | sort total DESC | head 10
	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.WhereCommand{
				Condition: &ast.BinaryExpression{
					Left:     &ast.FieldReference{Name: "status"},
					Operator: "=",
					Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
				},
			},
			&ast.StatsCommand{
				Aggregations: []*ast.Aggregation{
					{
						Func: &ast.FunctionCall{Name: "count"},
						Alias:      "total",
					},
				},
				GroupBy: []ast.Expression{
					&ast.FieldReference{Name: "host"},
				},
			},
			&ast.SortCommand{
				SortKeys: []*ast.SortKey{
					{
						Field:      &ast.FieldReference{Name: "total"},
						Descending: true,
					},
				},
			},
			&ast.HeadCommand{Count: 10},
		},
	}

	err := analyzer.Analyze(query)
	assert.NoError(t, err)
}

func TestTypeChecker_InferType_Literal(t *testing.T) {
	schema := createTestSchema()
	tc := NewTypeChecker(schema, NewScope(nil))

	tests := []struct {
		name     string
		literal  *ast.Literal
		expected FieldType
	}{
		{"integer", &ast.Literal{Value: 42, LiteralTyp: ast.LiteralTypeInt}, FieldTypeLong},
		{"float", &ast.Literal{Value: 3.14, LiteralTyp: ast.LiteralTypeFloat}, FieldTypeDouble},
		{"string", &ast.Literal{Value: "hello", LiteralTyp: ast.LiteralTypeString}, FieldTypeString},
		{"boolean", &ast.Literal{Value: true, LiteralTyp: ast.LiteralTypeBool}, FieldTypeBool},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, err := tc.InferType(tt.literal)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, fieldType)
		})
	}
}

func TestTypeChecker_InferType_FieldReference(t *testing.T) {
	schema := createTestSchema()
	tc := NewTypeChecker(schema, NewScope(nil))

	tests := []struct {
		name     string
		field    string
		expected FieldType
	}{
		{"status", "status", FieldTypeInt},
		{"host", "host", FieldTypeString},
		{"timestamp", "timestamp", FieldTypeDate},
		{"latency", "latency", FieldTypeDouble},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, err := tc.InferType(&ast.FieldReference{Name: tt.field})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, fieldType)
		})
	}
}

func TestTypeChecker_InferType_BinaryExpression(t *testing.T) {
	schema := createTestSchema()
	tc := NewTypeChecker(schema, NewScope(nil))

	tests := []struct {
		name     string
		expr     *ast.BinaryExpression
		expected FieldType
	}{
		{
			name: "comparison",
			expr: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: "=",
				Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
			},
			expected: FieldTypeBool,
		},
		{
			name: "arithmetic",
			expr: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "latency"},
				Operator: "+",
				Right:    &ast.Literal{Value: 10.0, LiteralTyp: ast.LiteralTypeFloat},
			},
			expected: FieldTypeDouble,
		},
		{
			name: "logical",
			expr: &ast.BinaryExpression{
				Left: &ast.BinaryExpression{
					Left:     &ast.FieldReference{Name: "status"},
					Operator: "=",
					Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
				},
				Operator: "AND",
				Right: &ast.BinaryExpression{
					Left:     &ast.FieldReference{Name: "host"},
					Operator: "=",
					Right:    &ast.Literal{Value: "server1", LiteralTyp: ast.LiteralTypeString},
				},
			},
			expected: FieldTypeBool,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, err := tc.InferType(tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, fieldType)
		})
	}
}
