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
// Type Checker Extended Tests (for missing coverage)
// =====================================================================

func TestTypeChecker_InferFunctionType(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("CountFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "count",
			Arguments: []ast.Expression{},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("SumFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "sum",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "status"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		// Should infer from argument type
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("AvgFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "avg",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("MinFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "min",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "status"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("MaxFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "max",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		}
		fieldType, err := typeChecker.InferType(fn)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("UnknownFunction", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name:      "unknown_func",
			Arguments: []ast.Expression{},
		}
		_, err := typeChecker.InferType(fn)
		// Unknown functions cause an error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown function")
	})

	t.Run("FunctionWithInvalidArgument", func(t *testing.T) {
		fn := &ast.FunctionCall{
			Name: "sum",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "nonexistent_field"},
			},
		}
		_, err := typeChecker.InferType(fn)
		assert.Error(t, err)
	})
}

func TestTypeChecker_InferArithmeticType(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("IntPlusInt", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
			Operator: "+",
			Right:    &ast.Literal{Value: 20, LiteralTyp: ast.LiteralTypeInt},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("IntPlusFloat", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
			Operator: "+",
			Right:    &ast.Literal{Value: 20.5, LiteralTyp: ast.LiteralTypeFloat},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		// Should promote to double
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("FloatTimesFloat", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: 3.14, LiteralTyp: ast.LiteralTypeFloat},
			Operator: "*",
			Right:    &ast.Literal{Value: 2.5, LiteralTyp: ast.LiteralTypeFloat},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("StringArithmeticInvalid", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: "hello", LiteralTyp: ast.LiteralTypeString},
			Operator: "+",
			Right:    &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
		}
		_, err := typeChecker.InferType(expr)
		assert.Error(t, err)
	})

	t.Run("DivisionProducesFloat", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
			Operator: "/",
			Right:    &ast.Literal{Value: 3, LiteralTyp: ast.LiteralTypeInt},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		// Division should always produce float type
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("ModuloOperation", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
			Operator: "%",
			Right:    &ast.Literal{Value: 3, LiteralTyp: ast.LiteralTypeInt},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("SubtractionOperation", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "-",
			Right:    &ast.Literal{Value: 100, LiteralTyp: ast.LiteralTypeInt},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		// Arithmetic operations promote to long
		assert.Equal(t, FieldTypeLong, fieldType)
	})
}

func TestTypeChecker_ValidateComparison(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("CompareIntWithInt", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
			Operator: ">",
			Right:    &ast.Literal{Value: 5, LiteralTyp: ast.LiteralTypeInt},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("CompareStringWithString", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: "hello", LiteralTyp: ast.LiteralTypeString},
			Operator: "=",
			Right:    &ast.Literal{Value: "world", LiteralTyp: ast.LiteralTypeString},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("CompareIntWithFloat", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
			Operator: "<",
			Right:    &ast.Literal{Value: 10.5, LiteralTyp: ast.LiteralTypeFloat},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("CompareIncompatibleTypes", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left:     &ast.Literal{Value: "hello", LiteralTyp: ast.LiteralTypeString},
			Operator: ">",
			Right:    &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
		}
		_, err := typeChecker.InferType(expr)
		assert.Error(t, err)
	})
}

func TestTypeChecker_LogicalOperators(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("ANDOperator", func(t *testing.T) {
		expr := &ast.BinaryExpression{
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
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("OROperator", func(t *testing.T) {
		expr := &ast.BinaryExpression{
			Left: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: "=",
				Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
			},
			Operator: "OR",
			Right: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: "=",
				Right:    &ast.Literal{Value: 404, LiteralTyp: ast.LiteralTypeInt},
			},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})
}

func TestTypeChecker_UnaryExpressions(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("NotOperator", func(t *testing.T) {
		expr := &ast.UnaryExpression{
			Operator: "NOT",
			Operand: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: "=",
				Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
			},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)
	})

	t.Run("NegateOperator", func(t *testing.T) {
		expr := &ast.UnaryExpression{
			Operator: "-",
			Operand:  &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeLong, fieldType)
	})

	t.Run("NegateFloat", func(t *testing.T) {
		expr := &ast.UnaryExpression{
			Operator: "-",
			Operand:  &ast.Literal{Value: 3.14, LiteralTyp: ast.LiteralTypeFloat},
		}
		fieldType, err := typeChecker.InferType(expr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeDouble, fieldType)
	})

	t.Run("NegateString", func(t *testing.T) {
		expr := &ast.UnaryExpression{
			Operator: "-",
			Operand:  &ast.Literal{Value: "hello", LiteralTyp: ast.LiteralTypeString},
		}
		_, err := typeChecker.InferType(expr)
		assert.Error(t, err)
	})
}

func TestTypeChecker_CaseExpression(t *testing.T) {
	schema := createTestSchema()
	scope := NewScope(nil)
	typeChecker := NewTypeChecker(schema, scope)

	t.Run("SimpleCaseExpression", func(t *testing.T) {
		caseExpr := &ast.CaseExpression{
			WhenClauses: []*ast.WhenClause{
				{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: ">=",
						Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
					},
					Result: &ast.Literal{Value: "critical", LiteralTyp: ast.LiteralTypeString},
				},
				{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: ">=",
						Right:    &ast.Literal{Value: 400, LiteralTyp: ast.LiteralTypeInt},
					},
					Result: &ast.Literal{Value: "error", LiteralTyp: ast.LiteralTypeString},
				},
			},
			ElseResult: &ast.Literal{Value: "normal", LiteralTyp: ast.LiteralTypeString},
		}

		fieldType, err := typeChecker.InferType(caseExpr)
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, fieldType)
	})

	t.Run("CaseExpressionWithInvalidCondition", func(t *testing.T) {
		caseExpr := &ast.CaseExpression{
			WhenClauses: []*ast.WhenClause{
				{
					Condition: &ast.Literal{Value: "not a boolean", LiteralTyp: ast.LiteralTypeString},
					Result:    &ast.Literal{Value: "result", LiteralTyp: ast.LiteralTypeString},
				},
			},
			ElseResult: &ast.Literal{Value: "else", LiteralTyp: ast.LiteralTypeString},
		}

		_, err := typeChecker.InferType(caseExpr)
		assert.Error(t, err)
	})
}

func TestFieldType_String(t *testing.T) {
	tests := []struct {
		fieldType FieldType
		expected  string
	}{
		{FieldTypeInt, "integer"},
		{FieldTypeLong, "long"},
		{FieldTypeFloat, "float"},
		{FieldTypeDouble, "double"},
		{FieldTypeString, "string"},
		{FieldTypeBool, "boolean"},
		{FieldTypeDate, "date"},
		{FieldTypeIP, "ip"},
		{FieldTypeObject, "object"},
		{FieldTypeArray, "array"},
		{FieldTypeGeoPoint, "geo_point"},
		{FieldTypeText, "text"},
		{FieldTypeKeyword, "keyword"},
		{FieldTypeUnknown, "unknown"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.fieldType.String()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestFieldType_IsComparable(t *testing.T) {
	comparableTypes := []FieldType{
		FieldTypeInt, FieldTypeLong, FieldTypeFloat, FieldTypeDouble,
		FieldTypeString, FieldTypeDate,
	}

	for _, ft := range comparableTypes {
		t.Run(ft.String(), func(t *testing.T) {
			assert.True(t, ft.IsComparable())
		})
	}

	nonComparableTypes := []FieldType{
		FieldTypeObject,
	}

	for _, ft := range nonComparableTypes {
		t.Run(ft.String(), func(t *testing.T) {
			assert.False(t, ft.IsComparable())
		})
	}
}

func TestField_GetNestedField(t *testing.T) {
	// Create a field with nested structure
	field := &Field{
		Name: "metadata",
		Type: FieldTypeObject,
		Fields: map[string]*Field{
			"category": {
				Name: "category",
				Type: FieldTypeString,
			},
			"tags": {
				Name:  "tags",
				Type:  FieldTypeString,
				Array: true,
			},
			"nested": {
				Name: "nested",
				Type: FieldTypeObject,
				Fields: map[string]*Field{
					"value": {
						Name: "value",
						Type: FieldTypeInt,
					},
				},
			},
		},
	}

	t.Run("GetDirectNestedField", func(t *testing.T) {
		nestedField, err := field.GetNestedField("category")
		require.NoError(t, err)
		assert.Equal(t, "category", nestedField.Name)
		assert.Equal(t, FieldTypeString, nestedField.Type)
	})

	t.Run("GetDeeplyNestedField", func(t *testing.T) {
		nestedField, err := field.GetNestedField("nested.value")
		require.NoError(t, err)
		assert.Equal(t, "value", nestedField.Name)
		assert.Equal(t, FieldTypeInt, nestedField.Type)
	})

	t.Run("GetNonexistentNestedField", func(t *testing.T) {
		_, err := field.GetNestedField("nonexistent")
		assert.Error(t, err)
	})

	t.Run("GetArrayField", func(t *testing.T) {
		nestedField, err := field.GetNestedField("tags")
		require.NoError(t, err)
		assert.True(t, nestedField.Array)
	})
}
