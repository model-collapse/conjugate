// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinaryExpression_String(t *testing.T) {
	tests := []struct {
		name     string
		expr     *BinaryExpression
		expected string
	}{
		{
			name: "simple equality",
			expr: &BinaryExpression{
				Left:     &FieldReference{Name: "status"},
				Operator: "=",
				Right:    &Literal{Value: 200, LiteralTyp: LiteralTypeInt},
			},
			expected: "(status = 200)",
		},
		{
			name: "not equal",
			expr: &BinaryExpression{
				Left:     &FieldReference{Name: "method"},
				Operator: "!=",
				Right:    &Literal{Value: "POST", LiteralTyp: LiteralTypeString},
			},
			expected: "(method != \"POST\")",
		},
		{
			name: "greater than",
			expr: &BinaryExpression{
				Left:     &FieldReference{Name: "response_time"},
				Operator: ">",
				Right:    &Literal{Value: 1000, LiteralTyp: LiteralTypeInt},
			},
			expected: "(response_time > 1000)",
		},
		{
			name: "AND operator",
			expr: &BinaryExpression{
				Left: &BinaryExpression{
					Left:     &FieldReference{Name: "status"},
					Operator: "=",
					Right:    &Literal{Value: 200, LiteralTyp: LiteralTypeInt},
				},
				Operator: "AND",
				Right: &BinaryExpression{
					Left:     &FieldReference{Name: "method"},
					Operator: "=",
					Right:    &Literal{Value: "GET", LiteralTyp: LiteralTypeString},
				},
			},
			expected: "((status = 200) AND (method = \"GET\"))",
		},
		{
			name: "arithmetic addition",
			expr: &BinaryExpression{
				Left:     &FieldReference{Name: "value"},
				Operator: "+",
				Right:    &Literal{Value: 10, LiteralTyp: LiteralTypeInt},
			},
			expected: "(value + 10)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBinaryExpression_Type(t *testing.T) {
	expr := &BinaryExpression{}
	assert.Equal(t, NodeTypeBinaryExpression, expr.Type())
}

func TestUnaryExpression_String(t *testing.T) {
	tests := []struct {
		name     string
		expr     *UnaryExpression
		expected string
	}{
		{
			name: "NOT operator",
			expr: &UnaryExpression{
				Operator: "NOT",
				Operand: &BinaryExpression{
					Left:     &FieldReference{Name: "active"},
					Operator: "=",
					Right:    &Literal{Value: true, LiteralTyp: LiteralTypeBool},
				},
			},
			expected: "(NOT (active = true))",
		},
		{
			name: "negation",
			expr: &UnaryExpression{
				Operator: "-",
				Operand:  &Literal{Value: 5, LiteralTyp: LiteralTypeInt},
			},
			expected: "(- 5)",
		},
		{
			name: "positive",
			expr: &UnaryExpression{
				Operator: "+",
				Operand:  &FieldReference{Name: "count"},
			},
			expected: "(+ count)",
		},
		{
			name: "nested NOT",
			expr: &UnaryExpression{
				Operator: "NOT",
				Operand: &UnaryExpression{
					Operator: "NOT",
					Operand:  &FieldReference{Name: "flag"},
				},
			},
			expected: "(NOT (NOT flag))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnaryExpression_Type(t *testing.T) {
	expr := &UnaryExpression{}
	assert.Equal(t, NodeTypeUnaryExpression, expr.Type())
}

func TestFunctionCall_String(t *testing.T) {
	tests := []struct {
		name     string
		expr     *FunctionCall
		expected string
	}{
		{
			name: "no arguments",
			expr: &FunctionCall{
				Name:      "count",
				Arguments: []Expression{},
			},
			expected: "count()",
		},
		{
			name: "single argument",
			expr: &FunctionCall{
				Name: "max",
				Arguments: []Expression{
					&FieldReference{Name: "response_time"},
				},
			},
			expected: "max(response_time)",
		},
		{
			name: "multiple arguments",
			expr: &FunctionCall{
				Name: "concat",
				Arguments: []Expression{
					&FieldReference{Name: "first_name"},
					&Literal{Value: " ", LiteralTyp: LiteralTypeString},
					&FieldReference{Name: "last_name"},
				},
			},
			expected: "concat(first_name, \" \", last_name)",
		},
		{
			name: "distinct count",
			expr: &FunctionCall{
				Name: "count",
				Arguments: []Expression{
					&FieldReference{Name: "user_id"},
				},
				Distinct: true,
			},
			expected: "count(DISTINCT user_id)",
		},
		{
			name: "nested function",
			expr: &FunctionCall{
				Name: "round",
				Arguments: []Expression{
					&FunctionCall{
						Name: "avg",
						Arguments: []Expression{
							&FieldReference{Name: "price"},
						},
					},
				},
			},
			expected: "round(avg(price))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFunctionCall_Type(t *testing.T) {
	expr := &FunctionCall{}
	assert.Equal(t, NodeTypeFunctionCall, expr.Type())
}

func TestFieldReference_String(t *testing.T) {
	tests := []struct {
		name     string
		expr     *FieldReference
		expected string
	}{
		{
			name:     "simple field",
			expr:     &FieldReference{Name: "status"},
			expected: "status",
		},
		{
			name:     "nested field",
			expr:     &FieldReference{Name: "user.address.city"},
			expected: "user.address.city",
		},
		{
			name:     "empty name",
			expr:     &FieldReference{Name: ""},
			expected: "",
		},
		{
			name:     "special characters",
			expr:     &FieldReference{Name: "field_with_underscores"},
			expected: "field_with_underscores",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFieldReference_Type(t *testing.T) {
	expr := &FieldReference{}
	assert.Equal(t, NodeTypeFieldReference, expr.Type())
}

func TestLiteral_String(t *testing.T) {
	tests := []struct {
		name     string
		expr     *Literal
		expected string
	}{
		{
			name:     "null",
			expr:     &Literal{Value: nil, LiteralTyp: LiteralTypeNull},
			expected: "null",
		},
		{
			name:     "boolean true",
			expr:     &Literal{Value: true, LiteralTyp: LiteralTypeBool},
			expected: "true",
		},
		{
			name:     "boolean false",
			expr:     &Literal{Value: false, LiteralTyp: LiteralTypeBool},
			expected: "false",
		},
		{
			name:     "integer",
			expr:     &Literal{Value: 42, LiteralTyp: LiteralTypeInt},
			expected: "42",
		},
		{
			name:     "negative integer",
			expr:     &Literal{Value: -100, LiteralTyp: LiteralTypeInt},
			expected: "-100",
		},
		{
			name:     "zero",
			expr:     &Literal{Value: 0, LiteralTyp: LiteralTypeInt},
			expected: "0",
		},
		{
			name:     "float",
			expr:     &Literal{Value: 3.14, LiteralTyp: LiteralTypeFloat},
			expected: "3.14",
		},
		{
			name:     "string",
			expr:     &Literal{Value: "hello", LiteralTyp: LiteralTypeString},
			expected: "\"hello\"",
		},
		{
			name:     "empty string",
			expr:     &Literal{Value: "", LiteralTyp: LiteralTypeString},
			expected: "\"\"",
		},
		{
			name:     "string with quotes",
			expr:     &Literal{Value: "it's \"quoted\"", LiteralTyp: LiteralTypeString},
			expected: "\"it's \"quoted\"\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLiteral_Type(t *testing.T) {
	expr := &Literal{}
	assert.Equal(t, NodeTypeLiteral, expr.Type())
}

func TestListLiteral_String(t *testing.T) {
	tests := []struct {
		name     string
		expr     *ListLiteral
		expected string
	}{
		{
			name: "empty list",
			expr: &ListLiteral{
				Values: []Expression{},
			},
			expected: "()",
		},
		{
			name: "single value",
			expr: &ListLiteral{
				Values: []Expression{
					&Literal{Value: 1, LiteralTyp: LiteralTypeInt},
				},
			},
			expected: "(1)",
		},
		{
			name: "multiple integers",
			expr: &ListLiteral{
				Values: []Expression{
					&Literal{Value: 1, LiteralTyp: LiteralTypeInt},
					&Literal{Value: 2, LiteralTyp: LiteralTypeInt},
					&Literal{Value: 3, LiteralTyp: LiteralTypeInt},
				},
			},
			expected: "(1, 2, 3)",
		},
		{
			name: "mixed types",
			expr: &ListLiteral{
				Values: []Expression{
					&Literal{Value: 200, LiteralTyp: LiteralTypeInt},
					&Literal{Value: 404, LiteralTyp: LiteralTypeInt},
					&Literal{Value: 500, LiteralTyp: LiteralTypeInt},
				},
			},
			expected: "(200, 404, 500)",
		},
		{
			name: "string list",
			expr: &ListLiteral{
				Values: []Expression{
					&Literal{Value: "GET", LiteralTyp: LiteralTypeString},
					&Literal{Value: "POST", LiteralTyp: LiteralTypeString},
					&Literal{Value: "PUT", LiteralTyp: LiteralTypeString},
				},
			},
			expected: "(\"GET\", \"POST\", \"PUT\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestListLiteral_Type(t *testing.T) {
	expr := &ListLiteral{}
	assert.Equal(t, NodeTypeListLiteral, expr.Type())
}

func TestWhenClause_String(t *testing.T) {
	tests := []struct {
		name     string
		clause   *WhenClause
		expected string
	}{
		{
			name: "simple when",
			clause: &WhenClause{
				Condition: &BinaryExpression{
					Left:     &FieldReference{Name: "status"},
					Operator: "<",
					Right:    &Literal{Value: 300, LiteralTyp: LiteralTypeInt},
				},
				Result: &Literal{Value: "success", LiteralTyp: LiteralTypeString},
			},
			expected: "WHEN (status < 300) THEN \"success\"",
		},
		{
			name: "complex condition",
			clause: &WhenClause{
				Condition: &BinaryExpression{
					Left: &BinaryExpression{
						Left:     &FieldReference{Name: "status"},
						Operator: ">=",
						Right:    &Literal{Value: 400, LiteralTyp: LiteralTypeInt},
					},
					Operator: "AND",
					Right: &BinaryExpression{
						Left:     &FieldReference{Name: "status"},
						Operator: "<",
						Right:    &Literal{Value: 500, LiteralTyp: LiteralTypeInt},
					},
				},
				Result: &Literal{Value: "client_error", LiteralTyp: LiteralTypeString},
			},
			expected: "WHEN ((status >= 400) AND (status < 500)) THEN \"client_error\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.clause.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWhenClause_Type(t *testing.T) {
	clause := &WhenClause{}
	assert.Equal(t, NodeTypeWhenClause, clause.Type())
}

func TestCaseExpression_String(t *testing.T) {
	tests := []struct {
		name     string
		expr     *CaseExpression
		expected string
	}{
		{
			name: "simple case",
			expr: &CaseExpression{
				WhenClauses: []*WhenClause{
					{
						Condition: &BinaryExpression{
							Left:     &FieldReference{Name: "status"},
							Operator: "<",
							Right:    &Literal{Value: 300, LiteralTyp: LiteralTypeInt},
						},
						Result: &Literal{Value: "success", LiteralTyp: LiteralTypeString},
					},
				},
				ElseResult: nil,
			},
			expected: "CASE WHEN (status < 300) THEN \"success\" END",
		},
		{
			name: "case with else",
			expr: &CaseExpression{
				WhenClauses: []*WhenClause{
					{
						Condition: &BinaryExpression{
							Left:     &FieldReference{Name: "status"},
							Operator: "<",
							Right:    &Literal{Value: 300, LiteralTyp: LiteralTypeInt},
						},
						Result: &Literal{Value: "success", LiteralTyp: LiteralTypeString},
					},
				},
				ElseResult: &Literal{Value: "error", LiteralTyp: LiteralTypeString},
			},
			expected: "CASE WHEN (status < 300) THEN \"success\" ELSE \"error\" END",
		},
		{
			name: "multiple when clauses",
			expr: &CaseExpression{
				WhenClauses: []*WhenClause{
					{
						Condition: &BinaryExpression{
							Left:     &FieldReference{Name: "status"},
							Operator: "<",
							Right:    &Literal{Value: 300, LiteralTyp: LiteralTypeInt},
						},
						Result: &Literal{Value: "success", LiteralTyp: LiteralTypeString},
					},
					{
						Condition: &BinaryExpression{
							Left:     &FieldReference{Name: "status"},
							Operator: "<",
							Right:    &Literal{Value: 500, LiteralTyp: LiteralTypeInt},
						},
						Result: &Literal{Value: "client_error", LiteralTyp: LiteralTypeString},
					},
				},
				ElseResult: &Literal{Value: "server_error", LiteralTyp: LiteralTypeString},
			},
			expected: "CASE WHEN (status < 300) THEN \"success\" WHEN (status < 500) THEN \"client_error\" ELSE \"server_error\" END",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.expr.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCaseExpression_Type(t *testing.T) {
	expr := &CaseExpression{}
	assert.Equal(t, NodeTypeCaseExpression, expr.Type())
}

func TestExpression_EdgeCases(t *testing.T) {
	t.Run("nil field reference name", func(t *testing.T) {
		expr := &FieldReference{Name: ""}
		assert.Equal(t, "", expr.String())
	})

	t.Run("empty function arguments", func(t *testing.T) {
		expr := &FunctionCall{
			Name:      "count",
			Arguments: nil,
		}
		assert.Equal(t, "count()", expr.String())
	})

	t.Run("nil list values", func(t *testing.T) {
		expr := &ListLiteral{Values: nil}
		// Should not panic
		result := expr.String()
		assert.Equal(t, "()", result)
	})

	t.Run("case with no when clauses", func(t *testing.T) {
		expr := &CaseExpression{
			WhenClauses: nil,
			ElseResult:  &Literal{Value: "default", LiteralTyp: LiteralTypeString},
		}
		// Should not panic
		result := expr.String()
		assert.Contains(t, result, "CASE")
		assert.Contains(t, result, "END")
	})

	t.Run("deeply nested expressions", func(t *testing.T) {
		// Create a deeply nested binary expression: (((a = 1) AND (b = 2)) AND (c = 3))
		expr := &BinaryExpression{
			Left: &BinaryExpression{
				Left: &BinaryExpression{
					Left:     &FieldReference{Name: "a"},
					Operator: "=",
					Right:    &Literal{Value: 1, LiteralTyp: LiteralTypeInt},
				},
				Operator: "AND",
				Right: &BinaryExpression{
					Left:     &FieldReference{Name: "b"},
					Operator: "=",
					Right:    &Literal{Value: 2, LiteralTyp: LiteralTypeInt},
				},
			},
			Operator: "AND",
			Right: &BinaryExpression{
				Left:     &FieldReference{Name: "c"},
				Operator: "=",
				Right:    &Literal{Value: 3, LiteralTyp: LiteralTypeInt},
			},
		}

		result := expr.String()
		assert.NotEmpty(t, result)
		assert.Contains(t, result, "a")
		assert.Contains(t, result, "b")
		assert.Contains(t, result, "c")
	})
}

func TestExpression_Accept(t *testing.T) {
	// Create a simple visitor to test Accept
	visitor := &BaseVisitor{}

	tests := []struct {
		name string
		expr Expression
	}{
		{"BinaryExpression", &BinaryExpression{}},
		{"UnaryExpression", &UnaryExpression{}},
		{"FunctionCall", &FunctionCall{}},
		{"FieldReference", &FieldReference{}},
		{"Literal", &Literal{}},
		{"ListLiteral", &ListLiteral{}},
		{"CaseExpression", &CaseExpression{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result, err := tt.expr.Accept(visitor)
			require.NoError(t, err)
			assert.Nil(t, result)
		})
	}
}
