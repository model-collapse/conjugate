// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package dsl

import (
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
)

// QueryBuilder builds OpenSearch query DSL from expressions
type QueryBuilder struct {
	functionBuilder FunctionBuilder
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// BuildFilter builds a query DSL from an expression
func (qb *QueryBuilder) BuildFilter(expr ast.Expression) (map[string]interface{}, error) {
	// Try native OpenSearch DSL first (faster than WASM)
	if nativeQuery, err := qb.tryNativeDSL(expr); err == nil {
		return nativeQuery, nil
	}

	// Try WASM UDF pushdown
	if qb.functionBuilder != nil && qb.functionBuilder.CanBuildUDF(expr) {
		return qb.buildWASMUDFQuery(expr)
	}

	// If we can't handle it, return error
	return nil, fmt.Errorf("cannot build filter for expression type: %T", expr)
}

// tryNativeDSL attempts to build native OpenSearch DSL
func (qb *QueryBuilder) tryNativeDSL(expr ast.Expression) (map[string]interface{}, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpression:
		return qb.buildBinaryExpression(e)

	case *ast.UnaryExpression:
		return qb.buildUnaryExpression(e)

	case *ast.FieldReference:
		// Field reference alone means "field exists and is true"
		return map[string]interface{}{
			"term": map[string]interface{}{
				e.Name: true,
			},
		}, nil

	case *ast.FunctionCall:
		// Function calls require WASM UDFs
		return nil, fmt.Errorf("function calls require WASM UDF")

	default:
		return nil, fmt.Errorf("unsupported filter expression type: %T", expr)
	}
}

// buildWASMUDFQuery builds a wasm_udf query
func (qb *QueryBuilder) buildWASMUDFQuery(expr ast.Expression) (map[string]interface{}, error) {
	udfRef, err := qb.functionBuilder.BuildUDF(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to build WASM UDF: %w", err)
	}

	query := map[string]interface{}{
		"wasm_udf": map[string]interface{}{
			"name":    udfRef.Name,
			"version": udfRef.Version,
		},
	}

	// Add parameters if present
	if len(udfRef.Parameters) > 0 {
		query["wasm_udf"].(map[string]interface{})["parameters"] = udfRef.Parameters
	}

	// Add field bindings if present
	if len(udfRef.FieldBindings) > 0 {
		query["wasm_udf"].(map[string]interface{})["field_bindings"] = udfRef.FieldBindings
	}

	return query, nil
}

// buildBinaryExpression builds DSL for binary expressions
func (qb *QueryBuilder) buildBinaryExpression(expr *ast.BinaryExpression) (map[string]interface{}, error) {
	switch expr.Operator {
	case "AND", "OR":
		return qb.buildLogicalExpression(expr)

	case "=", "!=", ">", ">=", "<", "<=":
		return qb.buildComparisonExpression(expr)

	case "LIKE":
		return qb.buildLikeExpression(expr)

	case "IN":
		return qb.buildInExpression(expr)

	default:
		return nil, fmt.Errorf("unsupported operator: %s", expr.Operator)
	}
}

// buildLogicalExpression builds bool query for AND/OR
func (qb *QueryBuilder) buildLogicalExpression(expr *ast.BinaryExpression) (map[string]interface{}, error) {
	leftQuery, err := qb.BuildFilter(expr.Left)
	if err != nil {
		return nil, err
	}

	rightQuery, err := qb.BuildFilter(expr.Right)
	if err != nil {
		return nil, err
	}

	if expr.Operator == "AND" {
		// Combine with must
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{leftQuery, rightQuery},
			},
		}, nil
	} else {
		// Combine with should
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{leftQuery, rightQuery},
			},
		}, nil
	}
}

// buildComparisonExpression builds term/range queries for comparisons
func (qb *QueryBuilder) buildComparisonExpression(expr *ast.BinaryExpression) (map[string]interface{}, error) {
	// Left side should be a field reference
	fieldRef, ok := expr.Left.(*ast.FieldReference)
	if !ok {
		return nil, fmt.Errorf("comparison left side must be a field reference")
	}

	// Right side should be a literal
	literal, ok := expr.Right.(*ast.Literal)
	if !ok {
		return nil, fmt.Errorf("comparison right side must be a literal")
	}

	fieldName := fieldRef.Name
	value := literal.Value

	switch expr.Operator {
	case "=":
		// Use term query for exact match
		return map[string]interface{}{
			"term": map[string]interface{}{
				fieldName: value,
			},
		}, nil

	case "!=":
		// Use must_not with term query
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							fieldName: value,
						},
					},
				},
			},
		}, nil

	case ">":
		return map[string]interface{}{
			"range": map[string]interface{}{
				fieldName: map[string]interface{}{
					"gt": value,
				},
			},
		}, nil

	case ">=":
		return map[string]interface{}{
			"range": map[string]interface{}{
				fieldName: map[string]interface{}{
					"gte": value,
				},
			},
		}, nil

	case "<":
		return map[string]interface{}{
			"range": map[string]interface{}{
				fieldName: map[string]interface{}{
					"lt": value,
				},
			},
		}, nil

	case "<=":
		return map[string]interface{}{
			"range": map[string]interface{}{
				fieldName: map[string]interface{}{
					"lte": value,
				},
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported comparison operator: %s", expr.Operator)
	}
}

// buildLikeExpression builds wildcard query for LIKE
func (qb *QueryBuilder) buildLikeExpression(expr *ast.BinaryExpression) (map[string]interface{}, error) {
	fieldRef, ok := expr.Left.(*ast.FieldReference)
	if !ok {
		return nil, fmt.Errorf("LIKE left side must be a field reference")
	}

	literal, ok := expr.Right.(*ast.Literal)
	if !ok {
		return nil, fmt.Errorf("LIKE right side must be a literal")
	}

	pattern, ok := literal.Value.(string)
	if !ok {
		return nil, fmt.Errorf("LIKE pattern must be a string")
	}

	// Convert SQL wildcards (%) to OpenSearch wildcards (*)
	// This is a simplified conversion
	osPattern := pattern

	return map[string]interface{}{
		"wildcard": map[string]interface{}{
			fieldRef.Name: map[string]interface{}{
				"value": osPattern,
			},
		},
	}, nil
}

// buildInExpression builds terms query for IN
func (qb *QueryBuilder) buildInExpression(expr *ast.BinaryExpression) (map[string]interface{}, error) {
	fieldRef, ok := expr.Left.(*ast.FieldReference)
	if !ok {
		return nil, fmt.Errorf("IN left side must be a field reference")
	}

	listLiteral, ok := expr.Right.(*ast.ListLiteral)
	if !ok {
		return nil, fmt.Errorf("IN right side must be a list")
	}

	// Extract values from list
	values := make([]interface{}, 0, len(listLiteral.Values))
	for _, elem := range listLiteral.Values {
		if lit, ok := elem.(*ast.Literal); ok {
			values = append(values, lit.Value)
		} else {
			return nil, fmt.Errorf("IN list must contain only literals")
		}
	}

	return map[string]interface{}{
		"terms": map[string]interface{}{
			fieldRef.Name: values,
		},
	}, nil
}

// buildUnaryExpression builds DSL for unary expressions
func (qb *QueryBuilder) buildUnaryExpression(expr *ast.UnaryExpression) (map[string]interface{}, error) {
	switch expr.Operator {
	case "NOT":
		operandQuery, err := qb.BuildFilter(expr.Operand)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": []interface{}{operandQuery},
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported unary operator: %s", expr.Operator)
	}
}
