// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package functions

import (
	"fmt"
	"strings"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/dsl"
	"github.com/conjugate/conjugate/pkg/wasm"
)

// FunctionBuilder builds WASM UDF references from PPL expressions
type FunctionBuilder struct {
	registry    *FunctionRegistry
	udfRegistry *wasm.UDFRegistry
}

// NewFunctionBuilder creates a new function builder
func NewFunctionBuilder(udfRegistry *wasm.UDFRegistry) *FunctionBuilder {
	return &FunctionBuilder{
		registry:    GetBuiltinRegistry(),
		udfRegistry: udfRegistry,
	}
}

// CanBuildUDF checks if an expression can be converted to a WASM UDF
func (fb *FunctionBuilder) CanBuildUDF(expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.FunctionCall:
		// Check if function is in registry
		return fb.registry.HasFunction(e.Name)

	case *ast.BinaryExpression:
		// Check if both sides can be built
		// For now, only support function calls with comparisons
		switch e.Operator {
		case ">", ">=", "<", "<=", "=", "!=":
			// Check if left side is a function call
			if fn, ok := e.Left.(*ast.FunctionCall); ok {
				return fb.registry.HasFunction(fn.Name)
			}
		}
		return false

	case *ast.FieldReference, *ast.Literal:
		// Simple expressions don't need UDFs
		return false

	default:
		return false
	}
}

// BuildUDF creates a UDF reference from an expression
func (fb *FunctionBuilder) BuildUDF(expr ast.Expression) (*dsl.UDFReference, error) {
	switch e := expr.(type) {
	case *ast.FunctionCall:
		return fb.buildFunctionUDF(e, nil)

	case *ast.BinaryExpression:
		return fb.buildComparisonUDF(e)

	default:
		return nil, fmt.Errorf("unsupported expression type for UDF: %T", expr)
	}
}

// buildFunctionUDF builds a UDF for a function call
func (fb *FunctionBuilder) buildFunctionUDF(fn *ast.FunctionCall, comparison *ast.BinaryExpression) (*dsl.UDFReference, error) {
	funcName := strings.ToLower(fn.Name)

	// Get function info from registry
	funcInfo := fb.registry.GetFunction(funcName)
	if funcInfo == nil {
		return nil, fmt.Errorf("function not found in registry: %s", funcName)
	}

	// Ensure UDF is registered
	if err := fb.ensureUDFRegistered(funcInfo); err != nil {
		return nil, err
	}

	// Extract field bindings and parameters
	fieldBindings := make(map[string]string)
	params := make(map[string]interface{})

	// Process function arguments
	for i, arg := range fn.Arguments {
		paramName := fmt.Sprintf("arg%d", i)

		switch a := arg.(type) {
		case *ast.FieldReference:
			// Bind field to parameter
			fieldBindings[paramName] = a.Name

		case *ast.Literal:
			// Add literal as parameter
			params[paramName] = a.Value

		default:
			return nil, fmt.Errorf("unsupported argument type: %T", arg)
		}
	}

	// Add comparison operator and value if present
	if comparison != nil {
		params["operator"] = comparison.Operator

		if lit, ok := comparison.Right.(*ast.Literal); ok {
			params["threshold"] = lit.Value
		}
	}

	udfName := funcInfo.UDFName
	if comparison != nil {
		// Use comparison-specific UDF variant
		udfName = fmt.Sprintf("%s_cmp", funcInfo.UDFName)
	}

	return &dsl.UDFReference{
		Name:          udfName,
		Version:       "builtin",
		Parameters:    params,
		FieldBindings: fieldBindings,
	}, nil
}

// buildComparisonUDF builds a UDF for a comparison with a function call
func (fb *FunctionBuilder) buildComparisonUDF(expr *ast.BinaryExpression) (*dsl.UDFReference, error) {
	// Check if left side is a function call
	fn, ok := expr.Left.(*ast.FunctionCall)
	if !ok {
		return nil, fmt.Errorf("comparison left side must be a function call")
	}

	// Check if right side is a literal
	_, ok = expr.Right.(*ast.Literal)
	if !ok {
		return nil, fmt.Errorf("comparison right side must be a literal")
	}

	return fb.buildFunctionUDF(fn, expr)
}

// ensureUDFRegistered ensures a UDF is registered in the WASM registry
func (fb *FunctionBuilder) ensureUDFRegistered(funcInfo *FunctionInfo) error {
	// Check if already registered
	if _, err := fb.udfRegistry.Get(funcInfo.UDFName, "builtin"); err == nil {
		// Already registered
		return nil
	}

	// Get WASM bytes from embedded library
	wasmBytes, err := GetBuiltinWASM(funcInfo.UDFName)
	if err != nil {
		return fmt.Errorf("failed to get WASM for %s: %w", funcInfo.UDFName, err)
	}

	// Create metadata with WASM bytes
	metadata := &wasm.UDFMetadata{
		Name:         funcInfo.UDFName,
		Version:      "builtin",
		Description:  funcInfo.Description,
		Author:       "CONJUGATE",
		Category:     funcInfo.Category,
		Tags:         []string{"builtin", funcInfo.Category},
		FunctionName: "execute", // Standard entry point for all built-in UDFs
		WASMBytes:    wasmBytes,
		WASMSize:     len(wasmBytes),
	}

	// Register with WASM registry
	if err := fb.udfRegistry.Register(metadata); err != nil {
		return fmt.Errorf("failed to register UDF %s: %w", funcInfo.UDFName, err)
	}

	return nil
}

// BuildComputedField builds a UDF reference for a computed field (projection)
func (fb *FunctionBuilder) BuildComputedField(expr ast.Expression, alias string) (*dsl.UDFReference, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpression:
		return fb.buildArithmeticUDF(e, alias)

	case *ast.FunctionCall:
		return fb.buildFunctionUDF(e, nil)

	default:
		return nil, fmt.Errorf("unsupported computed field expression: %T", expr)
	}
}

// buildArithmeticUDF builds a UDF for arithmetic expressions
func (fb *FunctionBuilder) buildArithmeticUDF(expr *ast.BinaryExpression, alias string) (*dsl.UDFReference, error) {
	// For simple arithmetic like "latency * 2", use a generic arithmetic UDF
	switch expr.Operator {
	case "+", "-", "*", "/", "%":
		fieldBindings := make(map[string]string)
		params := make(map[string]interface{})

		// Extract left operand
		if field, ok := expr.Left.(*ast.FieldReference); ok {
			fieldBindings["field"] = field.Name
		} else {
			return nil, fmt.Errorf("arithmetic left operand must be a field")
		}

		// Extract right operand
		if lit, ok := expr.Right.(*ast.Literal); ok {
			params["operand"] = lit.Value
		} else {
			return nil, fmt.Errorf("arithmetic right operand must be a literal")
		}

		params["operator"] = expr.Operator

		return &dsl.UDFReference{
			Name:          "arithmetic",
			Version:       "builtin",
			Parameters:    params,
			FieldBindings: fieldBindings,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported operator for arithmetic UDF: %s", expr.Operator)
	}
}

// BuildAggregationUDF builds a UDF reference for aggregation expressions
func (fb *FunctionBuilder) BuildAggregationUDF(agg *ast.Aggregation) (*dsl.UDFReference, error) {
	// Check if aggregation function has a complex expression
	if len(agg.Func.Arguments) == 0 {
		return nil, nil // Simple count() - no UDF needed
	}

	arg := agg.Func.Arguments[0]

	// Check if argument is a simple field reference
	if _, ok := arg.(*ast.FieldReference); ok {
		return nil, nil // Simple field aggregation - no UDF needed
	}

	// Build UDF for complex expression
	return fb.BuildComputedField(arg, agg.Alias)
}
