// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package physical

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/planner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestSchema() *analyzer.Schema {
	schema := analyzer.NewSchema("logs")
	schema.AddField("status", analyzer.FieldTypeInt)
	schema.AddField("host", analyzer.FieldTypeString)
	schema.AddField("timestamp", analyzer.FieldTypeDate)
	schema.AddField("latency", analyzer.FieldTypeDouble)
	return schema
}

func TestPhysicalPlanner_SimpleScan(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	logicalPlan := &planner.LogicalScan{
		Source:       "logs",
		OutputSchema: schema,
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	scan, ok := physicalPlan.(*PhysicalScan)
	assert.True(t, ok)
	assert.Equal(t, "logs", scan.Source)
	assert.Equal(t, ExecuteOnDataNode, scan.Location())
}

func TestPhysicalPlanner_FilterPushDown(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	// Logical: Filter -> Scan
	condition := &ast.BinaryExpression{
		Left:     &ast.FieldReference{Name: "status"},
		Operator: "=",
		Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
	}

	logicalPlan := &planner.LogicalFilter{
		Condition: condition,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should push down to scan
	scan, ok := physicalPlan.(*PhysicalScan)
	assert.True(t, ok, "Expected PhysicalScan with pushed filter")
	assert.NotNil(t, scan.Filter)
	assert.Equal(t, condition, scan.Filter)
}

func TestPhysicalPlanner_ProjectPushDown(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	// Logical: Project -> Scan
	logicalPlan := &planner.LogicalProject{
		Fields: []ast.Expression{
			&ast.FieldReference{Name: "status"},
			&ast.FieldReference{Name: "host"},
		},
		OutputSchema: schema,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should push down to scan
	scan, ok := physicalPlan.(*PhysicalScan)
	assert.True(t, ok, "Expected PhysicalScan with pushed projection")
	assert.Equal(t, 2, len(scan.Fields))
	assert.Contains(t, scan.Fields, "status")
	assert.Contains(t, scan.Fields, "host")
}

func TestPhysicalPlanner_SortPushDown(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	// Logical: Sort -> Scan
	logicalPlan := &planner.LogicalSort{
		SortKeys: []*ast.SortKey{
			{Field: &ast.FieldReference{Name: "timestamp"}, Descending: true},
		},
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should push down to scan
	scan, ok := physicalPlan.(*PhysicalScan)
	assert.True(t, ok, "Expected PhysicalScan with pushed sort")
	assert.Equal(t, 1, len(scan.SortKeys))
}

func TestPhysicalPlanner_LimitPushDown(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	// Logical: Limit -> Scan
	logicalPlan := &planner.LogicalLimit{
		Count: 10,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should push down to scan
	scan, ok := physicalPlan.(*PhysicalScan)
	assert.True(t, ok, "Expected PhysicalScan with pushed limit")
	assert.Equal(t, 10, scan.Limit)
}

func TestPhysicalPlanner_MultiplePushDown(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	// Logical: Limit -> Sort -> Project -> Filter -> Scan
	logicalPlan := &planner.LogicalLimit{
		Count: 10,
		Input: &planner.LogicalSort{
			SortKeys: []*ast.SortKey{
				{Field: &ast.FieldReference{Name: "timestamp"}, Descending: true},
			},
			Input: &planner.LogicalProject{
				Fields: []ast.Expression{
					&ast.FieldReference{Name: "status"},
					&ast.FieldReference{Name: "host"},
				},
				OutputSchema: schema,
				Input: &planner.LogicalFilter{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: "=",
						Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
					},
					Input: &planner.LogicalScan{
						Source:       "logs",
						OutputSchema: schema,
					},
				},
			},
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	t.Logf("Physical plan:\n%s", PrintPlan(physicalPlan, 0))

	// Should push down all operations to scan
	scan, ok := physicalPlan.(*PhysicalScan)
	assert.True(t, ok, "Expected PhysicalScan with all ops pushed down")
	assert.NotNil(t, scan.Filter, "Filter should be pushed down")
	assert.Equal(t, 2, len(scan.Fields), "Project should be pushed down")
	assert.Equal(t, 1, len(scan.SortKeys), "Sort should be pushed down")
	assert.Equal(t, 10, scan.Limit, "Limit should be pushed down")
}

func TestPhysicalPlanner_AggregationNotPushedDown(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	// Logical: Aggregate -> Scan
	logicalPlan := &planner.LogicalAggregate{
		Aggregations: []*ast.Aggregation{
			{Func: &ast.FunctionCall{Name: "count"}, Alias: "total"},
		},
		GroupBy: []ast.Expression{
			&ast.FieldReference{Name: "host"},
		},
		OutputSchema: schema,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	t.Logf("Physical plan:\n%s", PrintPlan(physicalPlan, 0))

	// Should be PhysicalAggregate -> PhysicalScan (not pushed down)
	agg, ok := physicalPlan.(*PhysicalAggregate)
	assert.True(t, ok, "Expected PhysicalAggregate")
	assert.Equal(t, ExecuteOnCoordinator, agg.Location())

	scan, ok := agg.Input.(*PhysicalScan)
	assert.True(t, ok, "Expected PhysicalScan as input")
	assert.Nil(t, scan.Filter, "No filter to push down")
}

func TestPhysicalPlanner_NoPushDown(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner().WithPushDown(false)

	// Logical: Filter -> Scan
	logicalPlan := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should be PhysicalFilter -> PhysicalScan (no push-down)
	filter, ok := physicalPlan.(*PhysicalFilter)
	assert.True(t, ok, "Expected PhysicalFilter")
	assert.Equal(t, ExecuteOnCoordinator, filter.Location())

	scan, ok := filter.Input.(*PhysicalScan)
	assert.True(t, ok)
	assert.Nil(t, scan.Filter, "Filter should not be pushed down")
}

func TestPhysicalPlanner_ComplexFilterWithFunctionPushedDown(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	// Logical: Filter with function call -> Scan
	// Function calls CAN be pushed down now (for WASM UDF execution)
	logicalPlan := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left: &ast.FunctionCall{
				Name: "abs",
				Arguments: []ast.Expression{
					&ast.FieldReference{Name: "latency"},
				},
			},
			Operator: ">",
			Right:    &ast.Literal{Value: 100.0, LiteralTyp: ast.LiteralTypeFloat},
		},
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// With WASM UDF support, function calls CAN be pushed down to the scan
	scan, ok := physicalPlan.(*PhysicalScan)
	assert.True(t, ok, "Expected PhysicalScan (function pushed down as WASM UDF)")
	assert.NotNil(t, scan.Filter, "Filter with function should be pushed down")
}

func TestIsPushedDown(t *testing.T) {
	schema := createTestSchema()

	tests := []struct {
		name     string
		plan     PhysicalPlan
		expected bool
	}{
		{
			name: "simple scan",
			plan: &PhysicalScan{
				Source:       "logs",
				OutputSchema: schema,
			},
			expected: false,
		},
		{
			name: "scan with filter",
			plan: &PhysicalScan{
				Source:       "logs",
				OutputSchema: schema,
				Filter: &ast.BinaryExpression{
					Left:     &ast.FieldReference{Name: "status"},
					Operator: "=",
					Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
				},
			},
			expected: true,
		},
		{
			name: "scan with limit",
			plan: &PhysicalScan{
				Source:       "logs",
				OutputSchema: schema,
				Limit:        10,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPushedDown(tt.plan)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCountCoordinatorOps(t *testing.T) {
	schema := createTestSchema()

	tests := []struct {
		name     string
		plan     PhysicalPlan
		expected int
	}{
		{
			name: "only scan",
			plan: &PhysicalScan{
				Source:       "logs",
				OutputSchema: schema,
			},
			expected: 0, // Scan runs on data node
		},
		{
			name: "filter on coordinator",
			plan: &PhysicalFilter{
				Condition: &ast.BinaryExpression{
					Left:     &ast.FieldReference{Name: "status"},
					Operator: "=",
					Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
				},
				Input: &PhysicalScan{
					Source:       "logs",
					OutputSchema: schema,
				},
			},
			expected: 1,
		},
		{
			name: "multiple coordinator ops",
			plan: &PhysicalLimit{
				Count: 10,
				Input: &PhysicalSort{
					SortKeys: []*ast.SortKey{
						{Field: &ast.FieldReference{Name: "timestamp"}, Descending: true},
					},
					Input: &PhysicalFilter{
						Condition: &ast.BinaryExpression{
							Left:     &ast.FieldReference{Name: "status"},
							Operator: "=",
							Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
						},
						Input: &PhysicalScan{
							Source:       "logs",
							OutputSchema: schema,
						},
					},
				},
			},
			expected: 3, // Limit, Sort, Filter all on coordinator
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := CountCoordinatorOps(tt.plan)
			assert.Equal(t, tt.expected, count)
		})
	}
}

func TestPhysicalPlanner_PrintPlan(t *testing.T) {
	schema := createTestSchema()

	plan := &PhysicalLimit{
		Count: 10,
		Input: &PhysicalFilter{
			Condition: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: "=",
				Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
			},
			Input: &PhysicalScan{
				Source:       "logs",
				OutputSchema: schema,
			},
		},
	}

	planStr := PrintPlan(plan, 0)
	t.Logf("Plan:\n%s", planStr)

	assert.Contains(t, planStr, "PhysicalLimit(10)")
	assert.Contains(t, planStr, "PhysicalFilter")
	assert.Contains(t, planStr, "PhysicalScan(logs)")
	assert.Contains(t, planStr, "[Coordinator]")
	assert.Contains(t, planStr, "[DataNode]")
}

func TestPhysicalPlanner_ComplexQuery(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	// Logical: Limit -> Sort -> Aggregate -> Filter -> Scan
	logicalPlan := &planner.LogicalLimit{
		Count: 10,
		Input: &planner.LogicalSort{
			SortKeys: []*ast.SortKey{
				{Field: &ast.FieldReference{Name: "total"}, Descending: true},
			},
			Input: &planner.LogicalAggregate{
				Aggregations: []*ast.Aggregation{
					{Func: &ast.FunctionCall{Name: "count"}, Alias: "total"},
				},
				GroupBy: []ast.Expression{
					&ast.FieldReference{Name: "host"},
				},
				OutputSchema: schema,
				Input: &planner.LogicalFilter{
					Condition: &ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "status"},
						Operator: "=",
						Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
					},
					Input: &planner.LogicalScan{
						Source:       "logs",
						OutputSchema: schema,
					},
				},
			},
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	t.Logf("Physical plan:\n%s", PrintPlan(physicalPlan, 0))

	// Verify structure
	// Filter should be pushed down, aggregation runs on coordinator
	assert.NotNil(t, physicalPlan)

	// Check that we have pushed down operations
	isPushed := IsPushedDown(physicalPlan)
	assert.True(t, isPushed, "Should have pushed down operations")

	// Count coordinator operations (should have aggregate, sort, limit)
	coordOps := CountCoordinatorOps(physicalPlan)
	assert.GreaterOrEqual(t, coordOps, 2, "Should have coordinator operations")
}

func TestSelectAggregationAlgorithm(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	agg := &planner.LogicalAggregate{
		Aggregations: []*ast.Aggregation{
			{Func: &ast.FunctionCall{Name: "count"}, Alias: "total"},
		},
		GroupBy: []ast.Expression{
			&ast.FieldReference{Name: "host"},
		},
		OutputSchema: schema,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	algorithm := pp.selectAggregationAlgorithm(agg)
	// Should default to hash aggregation
	assert.Equal(t, HashAggregation, algorithm)
}

// Tier 1 Physical Operator Tests

func TestPhysicalPlanner_Dedup(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	logicalPlan := &planner.LogicalDedup{
		Fields:      []ast.Expression{&ast.FieldReference{Name: "host"}},
		Count:       2,
		Consecutive: false,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should be PhysicalDedup -> PhysicalScan
	dedup, ok := physicalPlan.(*PhysicalDedup)
	assert.True(t, ok, "Expected PhysicalDedup, got %T", physicalPlan)
	assert.Len(t, dedup.Fields, 1)
	assert.Equal(t, 2, dedup.Count)
	assert.False(t, dedup.Consecutive)
	assert.Equal(t, ExecuteOnCoordinator, dedup.Location())

	_, ok = dedup.Input.(*PhysicalScan)
	assert.True(t, ok, "Expected PhysicalScan")
}

func TestPhysicalPlanner_Bin(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	logicalPlan := &planner.LogicalBin{
		Field: &ast.FieldReference{Name: "latency"},
		Bins:  10,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should be PhysicalBin -> PhysicalScan
	bin, ok := physicalPlan.(*PhysicalBin)
	assert.True(t, ok, "Expected PhysicalBin, got %T", physicalPlan)
	assert.NotNil(t, bin.Field)
	assert.Equal(t, 10, bin.Bins)
	assert.Equal(t, ExecuteOnCoordinator, bin.Location())
}

func TestPhysicalPlanner_Top(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	outputSchema := analyzer.NewSchema("logs")
	outputSchema.AddField("host", analyzer.FieldTypeString)
	outputSchema.AddField("count", analyzer.FieldTypeLong)

	logicalPlan := &planner.LogicalTop{
		Fields:       []ast.Expression{&ast.FieldReference{Name: "host"}},
		Limit:        10,
		ShowCount:    true,
		ShowPercent:  false,
		OutputSchema: outputSchema,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should be PhysicalTop -> PhysicalScan
	top, ok := physicalPlan.(*PhysicalTop)
	assert.True(t, ok, "Expected PhysicalTop, got %T", physicalPlan)
	assert.Len(t, top.Fields, 1)
	assert.Equal(t, 10, top.Limit)
	assert.True(t, top.ShowCount)
	assert.False(t, top.ShowPercent)
	assert.Equal(t, TopRareHash, top.Algorithm)
	assert.Equal(t, ExecuteOnCoordinator, top.Location())
}

func TestPhysicalPlanner_Rare(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	outputSchema := analyzer.NewSchema("logs")
	outputSchema.AddField("status", analyzer.FieldTypeInt)
	outputSchema.AddField("count", analyzer.FieldTypeLong)

	logicalPlan := &planner.LogicalRare{
		Fields:       []ast.Expression{&ast.FieldReference{Name: "status"}},
		Limit:        5,
		ShowCount:    true,
		ShowPercent:  true,
		OutputSchema: outputSchema,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should be PhysicalRare -> PhysicalScan
	rare, ok := physicalPlan.(*PhysicalRare)
	assert.True(t, ok, "Expected PhysicalRare, got %T", physicalPlan)
	assert.Len(t, rare.Fields, 1)
	assert.Equal(t, 5, rare.Limit)
	assert.True(t, rare.ShowCount)
	assert.True(t, rare.ShowPercent)
	assert.Equal(t, TopRareHash, rare.Algorithm)
	assert.Equal(t, ExecuteOnCoordinator, rare.Location())
}

func TestPhysicalPlanner_Eval(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	outputSchema := schema.Clone()
	outputSchema.AddField("speed", analyzer.FieldTypeDouble)

	logicalPlan := &planner.LogicalEval{
		Assignments: []*ast.EvalAssignment{
			{
				Field: "speed",
				Expression: &ast.BinaryExpression{
					Left:     &ast.Literal{Value: 100.0, LiteralTyp: ast.LiteralTypeFloat},
					Operator: "/",
					Right:    &ast.FieldReference{Name: "latency"},
				},
			},
		},
		OutputSchema: outputSchema,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should be PhysicalEval -> PhysicalScan
	eval, ok := physicalPlan.(*PhysicalEval)
	assert.True(t, ok, "Expected PhysicalEval, got %T", physicalPlan)
	assert.Len(t, eval.Assignments, 1)
	assert.Equal(t, "speed", eval.Assignments[0].Field)
	assert.Equal(t, ExecuteOnCoordinator, eval.Location())
}

func TestPhysicalPlanner_Rename(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	outputSchema := schema.Clone()
	outputSchema.AddField("server", analyzer.FieldTypeString)

	logicalPlan := &planner.LogicalRename{
		Assignments: []*ast.RenameAssignment{
			{OldName: "host", NewName: "server"},
		},
		OutputSchema: outputSchema,
		Input: &planner.LogicalScan{
			Source:       "logs",
			OutputSchema: schema,
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Should be PhysicalRename -> PhysicalScan
	rename, ok := physicalPlan.(*PhysicalRename)
	assert.True(t, ok, "Expected PhysicalRename, got %T", physicalPlan)
	assert.Len(t, rename.Assignments, 1)
	assert.Equal(t, "host", rename.Assignments[0].OldName)
	assert.Equal(t, "server", rename.Assignments[0].NewName)
	assert.Equal(t, ExecuteOnCoordinator, rename.Location())
}

func TestPhysicalPlanner_Tier1Pipeline(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	outputSchema := schema.Clone()
	outputSchema.AddField("is_slow", analyzer.FieldTypeBool)

	// Complex Tier 1 pipeline: Scan -> Eval -> Filter -> Dedup
	logicalPlan := &planner.LogicalDedup{
		Fields: []ast.Expression{&ast.FieldReference{Name: "host"}},
		Count:  1,
		Input: &planner.LogicalFilter{
			Condition: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "is_slow"},
				Operator: "=",
				Right:    &ast.Literal{Value: true, LiteralTyp: ast.LiteralTypeBool},
			},
			Input: &planner.LogicalEval{
				Assignments: []*ast.EvalAssignment{
					{
						Field:      "is_slow",
						Expression: &ast.Literal{Value: true, LiteralTyp: ast.LiteralTypeBool},
					},
				},
				OutputSchema: outputSchema,
				Input: &planner.LogicalScan{
					Source:       "logs",
					OutputSchema: schema,
				},
			},
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	t.Logf("Physical plan:\n%s", PrintPlan(physicalPlan, 0))

	// Verify we have coordinator operations (eval and dedup are barriers)
	coordOps := CountCoordinatorOps(physicalPlan)
	assert.GreaterOrEqual(t, coordOps, 3, "Should have multiple coordinator operations")
}

func TestPhysicalPlanner_TopWithPushDown(t *testing.T) {
	schema := createTestSchema()
	pp := NewPhysicalPlanner()

	outputSchema := analyzer.NewSchema("logs")
	outputSchema.AddField("host", analyzer.FieldTypeString)
	outputSchema.AddField("count", analyzer.FieldTypeLong)

	// Pipeline: Scan -> Filter -> Top
	// Top creates a barrier, so filter stays on coordinator
	logicalPlan := &planner.LogicalTop{
		Fields:       []ast.Expression{&ast.FieldReference{Name: "host"}},
		Limit:        10,
		ShowCount:    true,
		OutputSchema: outputSchema,
		Input: &planner.LogicalFilter{
			Condition: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: "=",
				Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
			},
			Input: &planner.LogicalScan{
				Source:       "logs",
				OutputSchema: schema,
			},
		},
	}

	physicalPlan, err := pp.Plan(logicalPlan)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	t.Logf("Physical plan:\n%s", PrintPlan(physicalPlan, 0))

	// Verify the plan was created successfully
	// Note: Top creates a barrier, so filter stays on coordinator above top
	coordOps := CountCoordinatorOps(physicalPlan)
	assert.GreaterOrEqual(t, coordOps, 2, "Should have filter and top on coordinator")

	// Verify we have a scan at the leaf
	scans := GetLeafScans(physicalPlan)
	assert.Len(t, scans, 1, "Should have one leaf scan")
}
