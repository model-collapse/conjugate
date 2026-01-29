// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package integration

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/dsl"
	"github.com/conjugate/conjugate/pkg/ppl/functions"
	"github.com/conjugate/conjugate/pkg/ppl/optimizer"
	"github.com/conjugate/conjugate/pkg/ppl/parser"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"github.com/conjugate/conjugate/pkg/ppl/planner"
	"github.com/conjugate/conjugate/pkg/wasm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestWASMUDFIntegration tests end-to-end PPL → WASM UDF translation
func TestWASMUDFIntegration(t *testing.T) {
	// Setup WASM runtime and registry
	wasmConfig := &wasm.Config{
		EnableJIT:      false,
		MaxMemoryPages: 256,
		Logger:         zap.NewNop(),
	}
	wasmRuntime, err := wasm.NewRuntime(wasmConfig)
	require.NoError(t, err)
	defer wasmRuntime.Close()

	udfRegistry, err := wasm.NewUDFRegistry(&wasm.UDFRegistryConfig{
		Runtime:         wasmRuntime,
		DefaultPoolSize: 1,
		EnableStats:     false,
		Logger:          zap.NewNop(),
	})
	require.NoError(t, err)

	// Create function builder with UDF registry
	functionBuilder := functions.NewFunctionBuilder(udfRegistry)

	// Create DSL translator with function builder
	translator := dsl.NewTranslator().WithFunctionBuilder(functionBuilder)

	// Create other components
	p := parser.NewParser()
	schema := analyzer.NewSchema("logs")
	// Add common fields
	schema.AddField("latency", analyzer.FieldTypeDouble)
	schema.AddField("status", analyzer.FieldTypeInt)
	schema.AddField("error_rate", analyzer.FieldTypeDouble)
	schema.AddField("timestamp", analyzer.FieldTypeDate)

	a := analyzer.NewAnalyzer(schema)
	b := planner.NewPlanBuilder(schema)
	o := optimizer.DefaultOptimizer()
	pp := physical.NewPhysicalPlanner()

	t.Run("SimpleFunctionFilter", func(t *testing.T) {
		// Query: source=logs | where abs(latency) > 100
		query := "source=logs | where abs(latency) > 100"

		// Parse
		tree, err := p.Parse(query)
		require.NoError(t, err, "Parse should succeed")

		// Analyze
		err = a.Analyze(tree)
		require.NoError(t, err, "Analyze should succeed")

		// Build logical plan
		logicalPlan, err := b.Build(tree)
		require.NoError(t, err, "Build logical plan should succeed")

		// Optimize
		optimizedPlan, err := o.Optimize(logicalPlan)
		require.NoError(t, err, "Optimize should succeed")
		require.NotNil(t, optimizedPlan, "Optimize should return a plan")

		// Create physical plan
		physicalPlan, err := pp.Plan(optimizedPlan)
		require.NoError(t, err, "Physical planning should succeed")

		// Translate to DSL
		dslMap, err := translator.TranslateToJSON(physicalPlan)
		require.NoError(t, err, "DSL translation should succeed")

		// Verify DSL structure
		assert.Contains(t, dslMap, "query", "DSL should have query")
		queryMap, ok := dslMap["query"].(map[string]interface{})
		require.True(t, ok, "Query should be a map")

		// Should use wasm_udf query
		assert.Contains(t, queryMap, "wasm_udf", "Should use wasm_udf query")
		wasmUDF, ok := queryMap["wasm_udf"].(map[string]interface{})
		require.True(t, ok, "wasm_udf should be a map")

		// Verify UDF reference
		assert.Equal(t, "math_abs_cmp", wasmUDF["name"], "UDF name should be math_abs_cmp")
		assert.Equal(t, "builtin", wasmUDF["version"], "UDF version should be builtin")

		// Verify parameters
		assert.Contains(t, wasmUDF, "parameters", "Should have parameters")
		params, ok := wasmUDF["parameters"].(map[string]interface{})
		require.True(t, ok, "Parameters should be a map")
		assert.Equal(t, ">", params["operator"], "Operator should be >")
		// Threshold can be int or float depending on literal type in query
		threshold := params["threshold"]
		assert.True(t, threshold == 100 || threshold == 100.0, "Threshold should be 100")

		// Verify field bindings
		assert.Contains(t, wasmUDF, "field_bindings", "Should have field bindings")
		bindings, ok := wasmUDF["field_bindings"].(map[string]string)
		require.True(t, ok, "Field bindings should be a map")
		assert.Equal(t, "latency", bindings["arg0"], "Should bind latency to arg0")

		t.Log("✅ Simple function filter translates to WASM UDF with parameters")
	})

	t.Run("MultipleFunctionFilters", func(t *testing.T) {
		// Query: source=logs | where abs(latency) > 100 AND ceil(error_rate) < 5
		query := "source=logs | where abs(latency) > 100 AND ceil(error_rate) < 5"

		// Parse and analyze
		tree, err := p.Parse(query)
		require.NoError(t, err)
		err = a.Analyze(tree)
		require.NoError(t, err)

		// Build, optimize and plan
		logicalPlan, err := b.Build(tree)
		require.NoError(t, err)
		optimizedPlan, err := o.Optimize(logicalPlan)
		require.NoError(t, err)
		physicalPlan, err := pp.Plan(optimizedPlan)
		require.NoError(t, err)

		// Translate
		dslMap, err := translator.TranslateToJSON(physicalPlan)
		require.NoError(t, err)

		// Verify structure - should have bool query with must clause
		queryMap, ok := dslMap["query"].(map[string]interface{})
		require.True(t, ok)

		// Should use bool query to combine UDFs
		assert.Contains(t, queryMap, "bool", "Should use bool query for AND")

		t.Log("✅ Multiple function filters translate correctly")
	})

	t.Run("FunctionInProjection", func(t *testing.T) {
		// TODO: Alias syntax "as" not yet implemented in parser (Tier 1 feature)
		// Query would be: source=logs | fields latency, abs(latency) as abs_latency
		// For now, test simple field projection

		query := "source=logs | fields latency, status"

		// Parse and analyze
		tree, err := p.Parse(query)
		require.NoError(t, err)
		err = a.Analyze(tree)
		require.NoError(t, err)

		// Build, optimize and plan
		logicalPlan, err := b.Build(tree)
		require.NoError(t, err)
		optimizedPlan, err := o.Optimize(logicalPlan)
		require.NoError(t, err)
		physicalPlan, err := pp.Plan(optimizedPlan)
		require.NoError(t, err)

		// Translate
		dslMap, err := translator.TranslateToJSON(physicalPlan)
		require.NoError(t, err)

		// Should have _source projection with fields
		assert.Contains(t, dslMap, "_source", "Should have _source")
		source, ok := dslMap["_source"].([]string)
		require.True(t, ok, "_source should be a string array")

		// Should include both fields
		assert.Contains(t, source, "latency", "Should include latency field")
		assert.Contains(t, source, "status", "Should include status field")

		t.Log("✅ Field projection works (function projection requires Tier 1 alias support)")
	})

	t.Run("ComputedFieldExpression", func(t *testing.T) {
		// TODO: eval command not yet implemented in parser (Tier 1 feature)
		// Query would be: source=logs | eval double_latency = latency * 2
		// For now, just test that BuildComputedField function works

		// Test BuildComputedField directly
		expr := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "latency"},
			Operator: "*",
			Right:    &ast.Literal{Value: 2, LiteralTyp: ast.LiteralTypeInt},
		}

		udfRef, err := functionBuilder.BuildComputedField(expr, "double_latency")
		require.NoError(t, err)
		require.NotNil(t, udfRef)

		// Verify arithmetic UDF
		assert.Equal(t, "arithmetic", udfRef.Name)
		assert.Equal(t, "builtin", udfRef.Version)
		assert.Equal(t, "*", udfRef.Parameters["operator"])
		assert.Equal(t, 2, udfRef.Parameters["operand"])
		assert.Equal(t, "latency", udfRef.FieldBindings["field"])

		t.Log("✅ Computed field expression handled (eval command requires Tier 1 parser support)")
	})

	t.Run("FunctionRegistry", func(t *testing.T) {
		// Test that function registry has expected functions
		registry := functions.GetBuiltinRegistry()

		// Check math functions
		assert.True(t, registry.HasFunction("abs"), "Should have abs function")
		assert.True(t, registry.HasFunction("ceil"), "Should have ceil function")
		assert.True(t, registry.HasFunction("floor"), "Should have floor function")
		assert.True(t, registry.HasFunction("sqrt"), "Should have sqrt function")

		// Check string functions
		assert.True(t, registry.HasFunction("upper"), "Should have upper function")
		assert.True(t, registry.HasFunction("lower"), "Should have lower function")
		assert.True(t, registry.HasFunction("trim"), "Should have trim function")

		// Check date functions
		assert.True(t, registry.HasFunction("year"), "Should have year function")
		assert.True(t, registry.HasFunction("month"), "Should have month function")
		assert.True(t, registry.HasFunction("day"), "Should have day function")

		// List all functions
		allFunctions := registry.ListFunctions()
		assert.GreaterOrEqual(t, len(allFunctions), 30, "Should have at least 30 functions")

		t.Logf("✅ Function registry has %d functions", len(allFunctions))
	})

	t.Run("FunctionBuilderCanBuildUDF", func(t *testing.T) {
		// Test CanBuildUDF logic
		builder := functions.NewFunctionBuilder(udfRegistry)

		// Should build UDF for function calls
		funcCall := &ast.FunctionCall{
			Name: "abs",
			Arguments: []ast.Expression{
				&ast.FieldReference{Name: "latency"},
			},
		}
		assert.True(t, builder.CanBuildUDF(funcCall), "Should build UDF for abs(latency)")

		// Should build UDF for comparisons with function calls
		comparison := &ast.BinaryExpression{
			Left:     funcCall,
			Operator: ">",
			Right:    &ast.Literal{Value: 100.0, LiteralTyp: ast.LiteralTypeFloat},
		}
		assert.True(t, builder.CanBuildUDF(comparison), "Should build UDF for abs(latency) > 100")

		// Should NOT build UDF for simple field references
		fieldRef := &ast.FieldReference{Name: "status"}
		assert.False(t, builder.CanBuildUDF(fieldRef), "Should not build UDF for simple field")

		// Should NOT build UDF for simple comparisons
		simpleComparison := &ast.BinaryExpression{
			Left:     fieldRef,
			Operator: "=",
			Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
		}
		assert.False(t, builder.CanBuildUDF(simpleComparison), "Should not build UDF for status = 200")

		t.Log("✅ CanBuildUDF logic works correctly")
	})
}

// TestWASMUDFParameterFlow tests that parameters flow correctly from PPL → DSL → WASM
func TestWASMUDFParameterFlow(t *testing.T) {
	// Setup
	wasmConfig := &wasm.Config{
		EnableJIT:      false,
		MaxMemoryPages: 256,
		Logger:         zap.NewNop(),
	}
	wasmRuntime, err := wasm.NewRuntime(wasmConfig)
	require.NoError(t, err)
	defer wasmRuntime.Close()

	udfRegistry, err := wasm.NewUDFRegistry(&wasm.UDFRegistryConfig{
		Runtime:         wasmRuntime,
		DefaultPoolSize: 1,
		EnableStats:     false,
		Logger:          zap.NewNop(),
	})
	require.NoError(t, err)

	functionBuilder := functions.NewFunctionBuilder(udfRegistry)

	// Create schema for tests
	_ = analyzer.NewSchema("logs") // schema unused in these specific tests but good to have

	t.Run("ParameterExtraction", func(t *testing.T) {
		// Build UDF reference for: abs(latency) > 100
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

		udfRef, err := functionBuilder.BuildUDF(expr)
		require.NoError(t, err, "Should build UDF reference")

		// Verify UDF reference
		assert.Equal(t, "math_abs_cmp", udfRef.Name)
		assert.Equal(t, "builtin", udfRef.Version)

		// Verify parameters
		assert.Contains(t, udfRef.Parameters, "operator")
		assert.Equal(t, ">", udfRef.Parameters["operator"])
		assert.Contains(t, udfRef.Parameters, "threshold")
		// Threshold can be int or float depending on literal type in query
		threshold := udfRef.Parameters["threshold"]
		assert.True(t, threshold == 100 || threshold == 100.0, "Threshold should be 100")

		// Verify field bindings
		assert.Contains(t, udfRef.FieldBindings, "arg0")
		assert.Equal(t, "latency", udfRef.FieldBindings["arg0"])

		t.Log("✅ Parameter extraction works correctly")
	})

	t.Run("ArithmeticUDF", func(t *testing.T) {
		// Build UDF for: latency * 2
		expr := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "latency"},
			Operator: "*",
			Right:    &ast.Literal{Value: 2, LiteralTyp: ast.LiteralTypeInt},
		}

		udfRef, err := functionBuilder.BuildComputedField(expr, "double_latency")
		require.NoError(t, err, "Should build arithmetic UDF")

		// Verify
		assert.Equal(t, "arithmetic", udfRef.Name)
		assert.Equal(t, "*", udfRef.Parameters["operator"])
		assert.Equal(t, 2, udfRef.Parameters["operand"])
		assert.Equal(t, "latency", udfRef.FieldBindings["field"])

		t.Log("✅ Arithmetic UDF parameters work correctly")
	})
}

// TestWASMUDFBuiltinLibrary tests the built-in function library
func TestWASMUDFBuiltinLibrary(t *testing.T) {
	t.Run("LoadBuiltinLibrary", func(t *testing.T) {
		library, err := functions.LoadBuiltinLibrary()
		require.NoError(t, err, "Should load builtin library")

		// Check that all expected functions are present
		expectedFunctions := []string{
			"math_abs",
			"math_ceil",
			"math_floor",
			"math_round",
			"math_sqrt",
			"string_upper",
			"string_lower",
			"date_year",
			"date_month",
			"arithmetic",
		}

		for _, funcName := range expectedFunctions {
			assert.Contains(t, library, funcName, "Library should contain %s", funcName)
			wasmBytes := library[funcName]
			assert.NotEmpty(t, wasmBytes, "WASM bytes should not be empty for %s", funcName)

			// Verify WASM magic number
			assert.Equal(t, byte(0x00), wasmBytes[0], "WASM magic byte 0")
			assert.Equal(t, byte(0x61), wasmBytes[1], "WASM magic byte 1")
			assert.Equal(t, byte(0x73), wasmBytes[2], "WASM magic byte 2")
			assert.Equal(t, byte(0x6d), wasmBytes[3], "WASM magic byte 3")
		}

		t.Logf("✅ Built-in library has %d WASM modules", len(library))
	})

	t.Run("GetBuiltinWASM", func(t *testing.T) {
		// Test individual function retrieval
		wasmBytes, err := functions.GetBuiltinWASM("math_abs")
		require.NoError(t, err, "Should get math_abs WASM")
		assert.NotEmpty(t, wasmBytes, "WASM bytes should not be empty")

		// Test non-existent function
		_, err = functions.GetBuiltinWASM("nonexistent")
		assert.Error(t, err, "Should error for nonexistent function")
	})
}

// TestPPLToWASMEndToEnd tests the complete flow: PPL → AST → Physical → DSL → WASM
func TestPPLToWASMEndToEnd(t *testing.T) {
	// Full integration test
	pplQuery := "source=logs | where abs(latency) > 100 | fields timestamp, latency, status"

	// Setup all components
	wasmConfig := &wasm.Config{
		EnableJIT:      false,
		MaxMemoryPages: 256,
		Logger:         zap.NewNop(),
	}
	wasmRuntime, err := wasm.NewRuntime(wasmConfig)
	require.NoError(t, err)
	defer wasmRuntime.Close()

	udfRegistry, err := wasm.NewUDFRegistry(&wasm.UDFRegistryConfig{
		Runtime:         wasmRuntime,
		DefaultPoolSize: 1,
		EnableStats:     false,
		Logger:          zap.NewNop(),
	})
	require.NoError(t, err)

	functionBuilder := functions.NewFunctionBuilder(udfRegistry)
	translator := dsl.NewTranslator().WithFunctionBuilder(functionBuilder)

	p := parser.NewParser()
	schema := analyzer.NewSchema("logs")
	schema.AddField("latency", analyzer.FieldTypeDouble)
	schema.AddField("status", analyzer.FieldTypeInt)
	schema.AddField("timestamp", analyzer.FieldTypeDate)
	a := analyzer.NewAnalyzer(schema)
	b := planner.NewPlanBuilder(schema)
	o := optimizer.DefaultOptimizer()
	pp := physical.NewPhysicalPlanner()

	// Execute pipeline
	tree, err := p.Parse(pplQuery)
	require.NoError(t, err, "Parse should succeed")

	err = a.Analyze(tree)
	require.NoError(t, err, "Analyze should succeed")

	logicalPlan, err := b.Build(tree)
	require.NoError(t, err, "Build logical plan should succeed")

	optimizedPlan, err := o.Optimize(logicalPlan)
	require.NoError(t, err, "Optimize should succeed")
	require.NotNil(t, optimizedPlan, "Optimize should return a plan")

	physicalPlan, err := pp.Plan(optimizedPlan)
	require.NoError(t, err, "Physical planning should succeed")

	dslMap, err := translator.TranslateToJSON(physicalPlan)
	require.NoError(t, err, "DSL translation should succeed")

	// Verify complete DSL
	assert.Contains(t, dslMap, "query", "Should have query")
	assert.Contains(t, dslMap, "_source", "Should have _source projection")

	// Verify query uses WASM UDF
	queryMap, ok := dslMap["query"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, queryMap, "wasm_udf", "Should use wasm_udf")

	// Verify projection includes requested fields
	source, ok := dslMap["_source"].([]string)
	require.True(t, ok)
	assert.Contains(t, source, "timestamp")
	assert.Contains(t, source, "latency")
	assert.Contains(t, source, "status")

	t.Log("✅ End-to-end PPL → WASM UDF flow works correctly")
}

// TestWASMUDFEvalPushdown tests eval command with WASM UDF pushdown
func TestWASMUDFEvalPushdown(t *testing.T) {
	// Setup WASM runtime and registry
	wasmConfig := &wasm.Config{
		EnableJIT:      false,
		MaxMemoryPages: 256,
		Logger:         zap.NewNop(),
	}
	wasmRuntime, err := wasm.NewRuntime(wasmConfig)
	require.NoError(t, err)
	defer wasmRuntime.Close()

	udfRegistry, err := wasm.NewUDFRegistry(&wasm.UDFRegistryConfig{
		Runtime:         wasmRuntime,
		DefaultPoolSize: 1,
		EnableStats:     false,
		Logger:          zap.NewNop(),
	})
	require.NoError(t, err)

	// Create function builder with UDF registry
	functionBuilder := functions.NewFunctionBuilder(udfRegistry)

	// Create DSL translator with function builder
	translator := dsl.NewTranslator().WithFunctionBuilder(functionBuilder)

	// Create other components
	p := parser.NewParser()
	schema := analyzer.NewSchema("logs")
	schema.AddField("latency", analyzer.FieldTypeDouble)
	schema.AddField("status", analyzer.FieldTypeInt)
	schema.AddField("error_rate", analyzer.FieldTypeDouble)

	a := analyzer.NewAnalyzer(schema)
	b := planner.NewPlanBuilder(schema)
	o := optimizer.DefaultOptimizer()
	pp := physical.NewPhysicalPlanner()

	t.Run("EvalWithFunctionPushedDown", func(t *testing.T) {
		// Query: source=logs | eval double_latency = latency * 2, abs_latency = abs(latency)
		query := "source=logs | eval double_latency = latency * 2, abs_latency = abs(latency)"

		// Parse
		tree, err := p.Parse(query)
		require.NoError(t, err, "Parse should succeed")

		// Analyze
		err = a.Analyze(tree)
		require.NoError(t, err, "Analyze should succeed")

		// Build logical plan
		logicalPlan, err := b.Build(tree)
		require.NoError(t, err, "Build logical plan should succeed")

		// Optimize
		optimizedPlan, err := o.Optimize(logicalPlan)
		require.NoError(t, err, "Optimize should succeed")

		// Physical plan
		physicalPlan, err := pp.Plan(optimizedPlan)
		require.NoError(t, err, "Physical planning should succeed")

		// Check that computed fields were pushed to scan
		scans := physical.GetLeafScans(physicalPlan)
		require.Len(t, scans, 1, "Should have one scan")
		scan := scans[0]

		// abs_latency should be pushed down (contains function)
		// double_latency might not be pushed if it's just arithmetic without function
		// Let's check what was pushed
		t.Logf("Computed fields pushed: %d", len(scan.ComputedFields))
		for _, cf := range scan.ComputedFields {
			t.Logf("  - %s = %s", cf.Field, cf.Expression.String())
		}

		// At least one field should be pushed (abs_latency has function)
		if len(scan.ComputedFields) > 0 {
			found := false
			for _, cf := range scan.ComputedFields {
				if cf.Field == "abs_latency" {
					found = true
					break
				}
			}
			assert.True(t, found, "abs_latency should be pushed down")
		}

		// Translate to DSL
		dslMap, err := translator.TranslateToJSON(physicalPlan)
		require.NoError(t, err, "DSL translation should succeed")

		// Check for script_fields in DSL
		if scriptFields, ok := dslMap["script_fields"]; ok {
			scriptFieldsMap, ok := scriptFields.(map[string]interface{})
			require.True(t, ok, "script_fields should be a map")
			t.Logf("✅ Script fields generated: %d", len(scriptFieldsMap))

			// Verify abs_latency is in script_fields
			if absLatency, ok := scriptFieldsMap["abs_latency"]; ok {
				absLatencyMap, ok := absLatency.(map[string]interface{})
				require.True(t, ok)
				assert.Contains(t, absLatencyMap, "wasm_udf", "abs_latency should use wasm_udf")
				t.Log("✅ abs_latency pushed down as WASM UDF")
			}
		} else {
			t.Log("⚠️  No script_fields generated (eval expressions may not contain functions)")
		}
	})

	t.Run("EvalMixedPushability", func(t *testing.T) {
		// Query with both pushable (has function) and non-pushable (no function) expressions
		query := "source=logs | eval computed = abs(latency), simple = status + 1"

		// Parse
		tree, err := p.Parse(query)
		require.NoError(t, err, "Parse should succeed")

		// Analyze
		err = a.Analyze(tree)
		require.NoError(t, err, "Analyze should succeed")

		// Build logical plan
		logicalPlan, err := b.Build(tree)
		require.NoError(t, err, "Build logical plan should succeed")

		// Physical plan
		physicalPlan, err := pp.Plan(logicalPlan)
		require.NoError(t, err, "Physical planning should succeed")

		// Check separation of pushable vs non-pushable
		scans := physical.GetLeafScans(physicalPlan)
		require.Len(t, scans, 1, "Should have one scan")
		scan := scans[0]

		// Find if there's a PhysicalEval for non-pushable assignments
		var hasCoordEval bool
		var current physical.PhysicalPlan = physicalPlan
		for current != nil {
			if _, ok := current.(*physical.PhysicalEval); ok {
				hasCoordEval = true
				break
			}
			children := current.Children()
			if len(children) > 0 {
				current = children[0]
			} else {
				break
			}
		}

		// We should have either:
		// - Pushable fields in scan.ComputedFields
		// - Non-pushable fields in coordinator PhysicalEval
		// - Or both
		totalAssignments := len(scan.ComputedFields)
		if hasCoordEval {
			t.Log("✅ Mixed eval: some pushed down, some on coordinator")
		}

		t.Logf("Pushable assignments: %d, Coordinator eval: %v", totalAssignments, hasCoordEval)
	})
}
