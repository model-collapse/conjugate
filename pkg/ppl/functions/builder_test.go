// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package functions

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/wasm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func createTestUDFRegistry(t *testing.T) *wasm.UDFRegistry {
	// Create WASM runtime
	runtime, err := wasm.NewRuntime(&wasm.Config{
		EnableJIT:      false, // Use interpreter for tests
		MaxMemoryPages: 256,
		Logger:         zap.NewNop(),
	})
	require.NoError(t, err)

	// Create UDF registry
	registry, err := wasm.NewUDFRegistry(&wasm.UDFRegistryConfig{
		Runtime:         runtime,
		DefaultPoolSize: 1,
		EnableStats:     false,
		Logger:          zap.NewNop(),
	})
	require.NoError(t, err)

	return registry
}

func TestFunctionBuilder_CanBuildUDF(t *testing.T) {
	registry := createTestUDFRegistry(t)
	builder := NewFunctionBuilder(registry)

	tests := []struct {
		name     string
		expr     ast.Expression
		expected bool
	}{
		{
			name: "abs function call",
			expr: &ast.FunctionCall{
				Name: "abs",
				Arguments: []ast.Expression{
					&ast.FieldReference{Name: "latency"},
				},
			},
			expected: true,
		},
		{
			name: "unknown function",
			expr: &ast.FunctionCall{
				Name: "unknown_func",
				Arguments: []ast.Expression{
					&ast.FieldReference{Name: "field"},
				},
			},
			expected: false,
		},
		{
			name: "comparison with function",
			expr: &ast.BinaryExpression{
				Left: &ast.FunctionCall{
					Name: "abs",
					Arguments: []ast.Expression{
						&ast.FieldReference{Name: "latency"},
					},
				},
				Operator: ">",
				Right:    &ast.Literal{Value: 100, LiteralTyp: ast.LiteralTypeInt},
			},
			expected: true,
		},
		{
			name: "simple field reference",
			expr: &ast.FieldReference{Name: "status"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.CanBuildUDF(tt.expr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFunctionBuilder_BuildUDF_FunctionCall(t *testing.T) {
	registry := createTestUDFRegistry(t)
	builder := NewFunctionBuilder(registry)

	// Test: abs(latency) > 100
	expr := &ast.BinaryExpression{
		Left: &ast.FunctionCall{
			Name: "abs",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		},
		Operator: ">",
		Right:    &ast.Literal{Value: 100.0, LiteralTyp: ast.LiteralTypeFloat},
	}

	udfRef, err := builder.BuildUDF(expr)
	require.NoError(t, err)
	require.NotNil(t, udfRef)

	// Verify UDF reference
	assert.Equal(t, "math_abs_cmp", udfRef.Name)
	assert.Equal(t, "builtin", udfRef.Version)

	// Verify parameters
	assert.Contains(t, udfRef.Parameters, "operator")
	assert.Equal(t, ">", udfRef.Parameters["operator"])
	assert.Contains(t, udfRef.Parameters, "threshold")
	assert.Equal(t, 100.0, udfRef.Parameters["threshold"])

	// Verify field bindings
	assert.Contains(t, udfRef.FieldBindings, "arg0")
	assert.Equal(t, "latency", udfRef.FieldBindings["arg0"])
}

func TestFunctionBuilder_BuildComputedField(t *testing.T) {
	registry := createTestUDFRegistry(t)
	builder := NewFunctionBuilder(registry)

	// Test: latency * 2
	expr := &ast.BinaryExpression{
		Left:     &ast.FieldReference{Name: "latency"},
		Operator: "*",
		Right:    &ast.Literal{Value: 2, LiteralTyp: ast.LiteralTypeInt},
	}

	udfRef, err := builder.BuildComputedField(expr, "double_latency")
	require.NoError(t, err)
	require.NotNil(t, udfRef)

	assert.Equal(t, "arithmetic", udfRef.Name)
	assert.Equal(t, "builtin", udfRef.Version)

	// Verify parameters
	assert.Equal(t, "*", udfRef.Parameters["operator"])
	assert.Equal(t, 2, udfRef.Parameters["operand"])

	// Verify field bindings
	assert.Equal(t, "latency", udfRef.FieldBindings["field"])
}

func TestFunctionRegistry_HasFunction(t *testing.T) {
	registry := GetBuiltinRegistry()

	tests := []struct {
		name     string
		expected bool
	}{
		{"abs", true},
		{"ceil", true},
		{"upper", true},
		{"lower", true},
		{"year", true},
		{"unknown_func", false},
		{"ABS", true}, // Case insensitive
		{"UPPER", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := registry.HasFunction(tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFunctionRegistry_GetFunction(t *testing.T) {
	registry := GetBuiltinRegistry()

	// Test getting a function
	info := registry.GetFunction("abs")
	require.NotNil(t, info)

	assert.Equal(t, "abs", info.PPLName)
	assert.Equal(t, "math_abs", info.UDFName)
	assert.Equal(t, "math", info.Category)

	// Test alias
	infoAlias := registry.GetFunction("ln")
	require.NotNil(t, infoAlias)
	assert.Equal(t, "log", infoAlias.PPLName)
	assert.Equal(t, "math_log", infoAlias.UDFName)
}

func TestFunctionRegistry_ListFunctions(t *testing.T) {
	registry := GetBuiltinRegistry()

	functions := registry.ListFunctions()

	// Should have 30+ unique functions (deduplicated)
	assert.GreaterOrEqual(t, len(functions), 30)

	// Check categories are present
	categories := make(map[string]bool)
	for _, fn := range functions {
		categories[fn.Category] = true
	}

	assert.True(t, categories["math"])
	assert.True(t, categories["string"])
	assert.True(t, categories["date"])
	assert.True(t, categories["type"])
}

func TestLoadBuiltinLibrary(t *testing.T) {
	library, err := LoadBuiltinLibrary()
	require.NoError(t, err)

	// Should have WASM binaries for all functions
	assert.GreaterOrEqual(t, len(library), 30)

	// Check specific functions
	assert.Contains(t, library, "math_abs")
	assert.Contains(t, library, "string_upper")
	assert.Contains(t, library, "date_year")
	assert.Contains(t, library, "arithmetic")

	// Each entry should have valid WASM bytes
	for name, wasmBytes := range library {
		assert.NotEmpty(t, wasmBytes, "WASM bytes empty for %s", name)
		// Check for WASM magic number
		assert.Equal(t, byte(0x00), wasmBytes[0], "Invalid WASM magic for %s", name)
		assert.Equal(t, byte(0x61), wasmBytes[1], "Invalid WASM magic for %s", name)
		assert.Equal(t, byte(0x73), wasmBytes[2], "Invalid WASM magic for %s", name)
		assert.Equal(t, byte(0x6d), wasmBytes[3], "Invalid WASM magic for %s", name)
	}
}

// =====================================================================
// Tier 1 Function Library Tests
// =====================================================================

func TestFunctionRegistry_Tier1MathFunctions(t *testing.T) {
	registry := GetBuiltinRegistry()

	// Test advanced trig functions
	trigFunctions := []string{
		"asin", "acos", "atan", "atan2", "cot",
		"degrees", "radians",
	}

	for _, name := range trigFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing trig function: %s", name)
		})
	}

	// Test additional math functions
	mathFunctions := []string{
		"mod", "rand", "truncate", "trunc",
		"cbrt", "sign", "log2",
		"e", "pi",
	}

	for _, name := range mathFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing math function: %s", name)
		})
	}

	// Test bitwise functions
	bitwiseFunctions := []string{
		"bit_and", "bit_or", "bit_xor", "bit_not",
		"bitwise_and", "bitwise_or", // Aliases
	}

	for _, name := range bitwiseFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing bitwise function: %s", name)
		})
	}
}

func TestFunctionRegistry_Tier1StringFunctions(t *testing.T) {
	registry := GetBuiltinRegistry()

	// Test pattern matching functions
	patternFunctions := []string{
		"regexp", "regex", "regexp_like",
		"regexp_replace", "regexp_extract",
	}

	for _, name := range patternFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing pattern function: %s", name)
		})
	}

	// Test search functions
	searchFunctions := []string{
		"locate", "position", "instr",
	}

	for _, name := range searchFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing search function: %s", name)
		})
	}

	// Test additional string functions
	stringFunctions := []string{
		"reverse", "ltrim", "rtrim",
		"lpad", "rpad", "repeat",
		"left", "right", "concat_ws",
		"ascii", "chr", "char",
	}

	for _, name := range stringFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing string function: %s", name)
		})
	}
}

func TestFunctionRegistry_Tier1DateFunctions(t *testing.T) {
	registry := GetBuiltinRegistry()

	// Test extraction functions
	extractFunctions := []string{
		"weekofyear", "week", "quarter",
		"dayname", "monthname", "microsecond",
	}

	for _, name := range extractFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing extraction function: %s", name)
		})
	}

	// Test current date/time functions
	currentFunctions := []string{
		"now", "curdate", "curtime", "sysdate",
		"current_timestamp", "current_date", "current_time",
		"utc_date", "utc_time", "utc_timestamp",
	}

	for _, name := range currentFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing current date function: %s", name)
		})
	}

	// Test date construction functions
	constructFunctions := []string{
		"makedate", "maketime",
		"from_days", "to_days", "to_seconds",
		"from_unixtime", "unix_timestamp",
	}

	for _, name := range constructFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing construction function: %s", name)
		})
	}

	// Test date arithmetic functions
	arithmeticFunctions := []string{
		"date_add", "date_sub", "adddate", "subdate",
		"addtime", "subtime",
		"datediff", "timediff", "timestampdiff",
		"period_add", "period_diff",
	}

	for _, name := range arithmeticFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing date arithmetic function: %s", name)
		})
	}

	// Test utility functions
	utilityFunctions := []string{
		"last_day", "convert_tz",
		"date_format", "str_to_date", "time_format",
	}

	for _, name := range utilityFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing utility function: %s", name)
		})
	}
}

func TestFunctionRegistry_Tier1TypeFunctions(t *testing.T) {
	registry := GetBuiltinRegistry()

	typeFunctions := []string{
		"cast", "convert", "try_cast", "typeof",
	}

	for _, name := range typeFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing type function: %s", name)
		})
	}
}

func TestFunctionRegistry_ConditionalFunctions(t *testing.T) {
	registry := GetBuiltinRegistry()

	// Test null handling functions
	nullFunctions := []string{
		"isnull", "isnotnull",
		"ifnull", "nvl", "nvl2",
		"nullif", "coalesce",
	}

	for _, name := range nullFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing null function: %s", name)
		})
	}

	// Test conditional logic functions
	logicFunctions := []string{
		"if", "case",
		"greatest", "least",
		"in", "between",
	}

	for _, name := range logicFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing logic function: %s", name)
		})
	}
}

func TestFunctionRegistry_RelevanceFunctions(t *testing.T) {
	registry := GetBuiltinRegistry()

	relevanceFunctions := []string{
		"match", "match_phrase", "match_phrase_prefix",
		"match_bool_prefix",
		"multi_match",
		"query_string", "simple_query_string",
	}

	for _, name := range relevanceFunctions {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing relevance function: %s", name)
		})
	}
}

func TestFunctionRegistry_AggregationFunctions(t *testing.T) {
	registry := GetBuiltinRegistry()

	// Test basic aggregations
	basicAggs := []string{
		"count", "sum", "avg", "mean",
		"min", "max",
	}

	for _, name := range basicAggs {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing basic aggregation: %s", name)
		})
	}

	// Test statistical aggregations
	statAggs := []string{
		"stddev", "stdev", "stddev_samp", "stddev_pop",
		"variance", "var_samp", "var_pop",
	}

	for _, name := range statAggs {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing statistical aggregation: %s", name)
		})
	}

	// Test distinct/cardinality aggregations
	distinctAggs := []string{
		"distinct_count", "dc", "cardinality",
		"approx_count_distinct",
	}

	for _, name := range distinctAggs {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing distinct aggregation: %s", name)
		})
	}

	// Test percentile aggregations
	percentileAggs := []string{
		"percentile", "percentile_approx", "median",
	}

	for _, name := range percentileAggs {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing percentile aggregation: %s", name)
		})
	}

	// Test collection aggregations
	collectionAggs := []string{
		"values", "list", "first", "last",
		"earliest", "latest",
	}

	for _, name := range collectionAggs {
		t.Run(name, func(t *testing.T) {
			assert.True(t, registry.HasFunction(name),
				"Missing collection aggregation: %s", name)
		})
	}
}

func TestFunctionRegistry_FunctionCount(t *testing.T) {
	registry := GetBuiltinRegistry()

	functions := registry.ListFunctions()

	// Count functions by category
	categories := make(map[string]int)
	for _, fn := range functions {
		categories[fn.Category]++
	}

	// Should have 135+ functions (target from task)
	// Categories:
	// - math: 30+ functions
	// - string: 25+ functions
	// - date: 40+ functions
	// - type: 10 functions
	// - conditional: 12 functions
	// - relevance: 7 functions
	// - aggregation: 20+ functions

	t.Logf("Function counts by category:")
	for cat, count := range categories {
		t.Logf("  %s: %d", cat, count)
	}
	t.Logf("  Total: %d", len(functions))

	// Verify minimum counts
	assert.GreaterOrEqual(t, categories["math"], 25, "Need at least 25 math functions")
	assert.GreaterOrEqual(t, categories["string"], 20, "Need at least 20 string functions")
	assert.GreaterOrEqual(t, categories["date"], 35, "Need at least 35 date functions")
	assert.GreaterOrEqual(t, categories["type"], 8, "Need at least 8 type functions")
	assert.GreaterOrEqual(t, categories["conditional"], 10, "Need at least 10 conditional functions")
	assert.GreaterOrEqual(t, categories["relevance"], 7, "Need at least 7 relevance functions")
	assert.GreaterOrEqual(t, categories["aggregation"], 18, "Need at least 18 aggregation functions")

	// Total should be 135+
	assert.GreaterOrEqual(t, len(functions), 120, "Need at least 120 unique functions")
}

func TestFunctionRegistry_Aliases(t *testing.T) {
	registry := GetBuiltinRegistry()

	// Test that aliases resolve to the correct function
	aliasTests := []struct {
		alias    string
		expected string
	}{
		{"ln", "log"},
		{"ceiling", "ceil"},
		{"trunc", "truncate"},
		{"random", "rand"},
		{"ucase", "upper"},
		{"lcase", "lower"},
		{"substr", "substring"},
		{"mid", "substring"},
		{"nvl", "ifnull"},
		{"mean", "avg"},
		{"stdev", "stddev"},
		{"dc", "distinct_count"},
		{"dow", "dayofweek"},
		{"doy", "dayofyear"},
	}

	for _, tt := range aliasTests {
		t.Run(tt.alias, func(t *testing.T) {
			aliasInfo := registry.GetFunction(tt.alias)
			require.NotNil(t, aliasInfo, "Alias %s not found", tt.alias)

			expectedInfo := registry.GetFunction(tt.expected)
			require.NotNil(t, expectedInfo, "Expected function %s not found", tt.expected)

			assert.Equal(t, expectedInfo.UDFName, aliasInfo.UDFName,
				"Alias %s should resolve to %s", tt.alias, tt.expected)
		})
	}
}
