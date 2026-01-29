// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package planner

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestSchema() *analyzer.Schema {
	schema := analyzer.NewSchema("logs")
	schema.AddField("status", analyzer.FieldTypeInt)
	schema.AddField("host", analyzer.FieldTypeString)
	schema.AddField("timestamp", analyzer.FieldTypeDate)
	schema.AddField("latency", analyzer.FieldTypeDouble)
	schema.AddField("message", analyzer.FieldTypeText)
	return schema
}

func TestPlanBuilder_SearchCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be a LogicalScan
	scan, ok := plan.(*LogicalScan)
	assert.True(t, ok, "Expected LogicalScan, got %T", plan)
	assert.Equal(t, "logs", scan.Source)
	assert.Equal(t, schema, scan.OutputSchema)
}

func TestPlanBuilder_WhereCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

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

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Filter -> Scan
	filter, ok := plan.(*LogicalFilter)
	assert.True(t, ok, "Expected LogicalFilter, got %T", plan)
	assert.NotNil(t, filter.Condition)

	scan, ok := filter.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan, got %T", filter.Input)
	assert.Equal(t, "logs", scan.Source)
}

func TestPlanBuilder_FieldsCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.FieldsCommand{
				Fields: []ast.Expression{
					&ast.FieldReference{Name: "status"},
					&ast.FieldReference{Name: "host"},
				},
				Includes: true,
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Project -> Scan
	project, ok := plan.(*LogicalProject)
	assert.True(t, ok, "Expected LogicalProject, got %T", plan)
	assert.Equal(t, 2, len(project.Fields))
	assert.False(t, project.Exclude)

	_, ok = project.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan, got %T", project.Input)
}

func TestPlanBuilder_SortCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

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

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Sort -> Scan
	sort, ok := plan.(*LogicalSort)
	assert.True(t, ok, "Expected LogicalSort, got %T", plan)
	assert.Equal(t, 1, len(sort.SortKeys))
	assert.True(t, sort.SortKeys[0].Descending)

	_, ok = sort.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan, got %T", sort.Input)
}

func TestPlanBuilder_HeadCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.HeadCommand{Count: 10},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Limit -> Scan
	limit, ok := plan.(*LogicalLimit)
	assert.True(t, ok, "Expected LogicalLimit, got %T", plan)
	assert.Equal(t, 10, limit.Count)

	_, ok = limit.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan, got %T", limit.Input)
}

func TestPlanBuilder_StatsCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.StatsCommand{
				Aggregations: []*ast.Aggregation{
					{
						Func:  &ast.FunctionCall{Name: "count"},
						Alias: "total",
					},
				},
				GroupBy: []ast.Expression{
					&ast.FieldReference{Name: "host"},
				},
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Aggregate -> Scan
	agg, ok := plan.(*LogicalAggregate)
	assert.True(t, ok, "Expected LogicalAggregate, got %T", plan)
	assert.Equal(t, 1, len(agg.Aggregations))
	assert.Equal(t, 1, len(agg.GroupBy))

	_, ok = agg.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan, got %T", agg.Input)
}

func TestPlanBuilder_ComplexPipeline(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

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
						Func:  &ast.FunctionCall{Name: "count"},
						Alias: "total",
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

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Verify structure: Limit -> Sort -> Aggregate -> Filter -> Scan
	limit, ok := plan.(*LogicalLimit)
	require.True(t, ok, "Expected LogicalLimit at top, got %T", plan)
	assert.Equal(t, 10, limit.Count)

	sort, ok := limit.Input.(*LogicalSort)
	require.True(t, ok, "Expected LogicalSort, got %T", limit.Input)
	assert.Equal(t, 1, len(sort.SortKeys))

	agg, ok := sort.Input.(*LogicalAggregate)
	require.True(t, ok, "Expected LogicalAggregate, got %T", sort.Input)
	assert.Equal(t, 1, len(agg.Aggregations))

	filter, ok := agg.Input.(*LogicalFilter)
	require.True(t, ok, "Expected LogicalFilter, got %T", agg.Input)

	scan, ok := filter.Input.(*LogicalScan)
	require.True(t, ok, "Expected LogicalScan at bottom, got %T", filter.Input)
	assert.Equal(t, "logs", scan.Source)
}

func TestPrintPlan(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

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
			&ast.HeadCommand{Count: 10},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)

	planStr := PrintPlan(plan, 0)
	t.Logf("Plan:\n%s", planStr)

	// Should contain all operators
	assert.Contains(t, planStr, "Limit(10)")
	assert.Contains(t, planStr, "Filter")
	assert.Contains(t, planStr, "Scan(logs)")
}

func TestGetLeafScans(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

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

	plan, err := builder.Build(query)
	require.NoError(t, err)

	scans := GetLeafScans(plan)
	assert.Equal(t, 1, len(scans))
	assert.Equal(t, "logs", scans[0].Source)
}

func TestReplaceChild(t *testing.T) {
	schema := createTestSchema()

	scan1 := &LogicalScan{Source: "logs", OutputSchema: schema}
	scan2 := &LogicalScan{Source: "metrics", OutputSchema: schema}

	filter := &LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: scan1,
	}

	// Replace scan1 with scan2
	newFilter := ReplaceChild(filter, scan1, scan2)
	require.NotNil(t, newFilter)

	newFilterTyped, ok := newFilter.(*LogicalFilter)
	require.True(t, ok)

	newScan, ok := newFilterTyped.Input.(*LogicalScan)
	require.True(t, ok)
	assert.Equal(t, "metrics", newScan.Source)
}

func TestLogicalPlan_Schema(t *testing.T) {
	schema := createTestSchema()

	scan := &LogicalScan{Source: "logs", OutputSchema: schema}
	assert.Equal(t, schema, scan.Schema())

	filter := &LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: scan,
	}
	// Filter preserves input schema
	assert.Equal(t, schema, filter.Schema())
}

// Tier 1 Command Tests

func TestPlanBuilder_DedupCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.DedupCommand{
				Fields:      []ast.Expression{&ast.FieldReference{Name: "host"}},
				Count:       2,
				Consecutive: false,
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Dedup -> Scan
	dedup, ok := plan.(*LogicalDedup)
	assert.True(t, ok, "Expected LogicalDedup, got %T", plan)
	assert.Len(t, dedup.Fields, 1)
	assert.Equal(t, 2, dedup.Count)
	assert.False(t, dedup.Consecutive)

	_, ok = dedup.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan")
}

func TestPlanBuilder_BinCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.BinCommand{
				Field: &ast.FieldReference{Name: "latency"},
				Bins:  10,
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Bin -> Scan
	bin, ok := plan.(*LogicalBin)
	assert.True(t, ok, "Expected LogicalBin, got %T", plan)
	assert.NotNil(t, bin.Field)
	assert.Equal(t, 10, bin.Bins)

	_, ok = bin.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan")
}

func TestPlanBuilder_TopCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.TopCommand{
				Fields:      []ast.Expression{&ast.FieldReference{Name: "host"}},
				Limit:       10,
				ShowCount:   true,
				ShowPercent: true,
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Top -> Scan
	top, ok := plan.(*LogicalTop)
	assert.True(t, ok, "Expected LogicalTop, got %T", plan)
	assert.Len(t, top.Fields, 1)
	assert.Equal(t, 10, top.Limit)
	assert.True(t, top.ShowCount)
	assert.True(t, top.ShowPercent)
	assert.NotNil(t, top.OutputSchema)

	_, ok = top.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan")
}

func TestPlanBuilder_RareCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.RareCommand{
				Fields:      []ast.Expression{&ast.FieldReference{Name: "status"}},
				Limit:       5,
				ShowCount:   true,
				ShowPercent: false,
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Rare -> Scan
	rare, ok := plan.(*LogicalRare)
	assert.True(t, ok, "Expected LogicalRare, got %T", plan)
	assert.Len(t, rare.Fields, 1)
	assert.Equal(t, 5, rare.Limit)
	assert.True(t, rare.ShowCount)
	assert.False(t, rare.ShowPercent)

	_, ok = rare.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan")
}

func TestPlanBuilder_ChartCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.ChartCommand{
				Aggregations: []*ast.Aggregation{
					{
						Func: &ast.FunctionCall{
							Name:      "count",
							Arguments: nil,
						},
						Alias: "total",
					},
				},
				GroupBy: []ast.Expression{&ast.FieldReference{Name: "host"}},
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Chart should be converted to Aggregate
	agg, ok := plan.(*LogicalAggregate)
	assert.True(t, ok, "Expected LogicalAggregate, got %T", plan)
	assert.Len(t, agg.GroupBy, 1)
	assert.Len(t, agg.Aggregations, 1)
	assert.NotNil(t, agg.OutputSchema)
}

func TestPlanBuilder_TimechartCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.TimechartCommand{
				Aggregations: []*ast.Aggregation{
					{
						Func: &ast.FunctionCall{
							Name: "avg",
							Arguments: []ast.Expression{
								&ast.FieldReference{Name: "latency"},
							},
						},
						Alias: "avg_latency",
					},
				},
				Span: &ast.TimeSpan{Value: 1, Unit: "h"},
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Timechart should be converted to Aggregate with _time grouping
	agg, ok := plan.(*LogicalAggregate)
	assert.True(t, ok, "Expected LogicalAggregate, got %T", plan)
	assert.Len(t, agg.GroupBy, 1, "Should have _time grouping")
	assert.Len(t, agg.Aggregations, 1)

	// Check that _time field is in output schema
	_, err = agg.OutputSchema.GetField("_time")
	assert.NoError(t, err, "Output schema should have _time field")
}

func TestPlanBuilder_EvalCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.EvalCommand{
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
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Eval -> Scan
	eval, ok := plan.(*LogicalEval)
	assert.True(t, ok, "Expected LogicalEval, got %T", plan)
	assert.Len(t, eval.Assignments, 1)
	assert.Equal(t, "speed", eval.Assignments[0].Field)

	// Output schema should have the new field
	_, err = eval.OutputSchema.GetField("speed")
	assert.NoError(t, err, "Output schema should have 'speed' field")

	_, ok = eval.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan")
}

func TestPlanBuilder_RenameCommand(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.RenameCommand{
				Assignments: []*ast.RenameAssignment{
					{OldName: "host", NewName: "server"},
					{OldName: "status", NewName: "http_status"},
				},
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Should be Rename -> Scan
	rename, ok := plan.(*LogicalRename)
	assert.True(t, ok, "Expected LogicalRename, got %T", plan)
	assert.Len(t, rename.Assignments, 2)
	assert.Equal(t, "host", rename.Assignments[0].OldName)
	assert.Equal(t, "server", rename.Assignments[0].NewName)

	// Output schema should have the new field names
	_, err = rename.OutputSchema.GetField("server")
	assert.NoError(t, err, "Output schema should have 'server' field")
	_, err = rename.OutputSchema.GetField("http_status")
	assert.NoError(t, err, "Output schema should have 'http_status' field")

	_, ok = rename.Input.(*LogicalScan)
	assert.True(t, ok, "Expected input to be LogicalScan")
}

func TestPlanBuilder_Tier1ComplexPipeline(t *testing.T) {
	schema := createTestSchema()
	builder := NewPlanBuilder(schema)

	// Complex pipeline: search | eval | where | dedup | top
	query := &ast.Query{
		Commands: []ast.Command{
			&ast.SearchCommand{Source: "logs"},
			&ast.EvalCommand{
				Assignments: []*ast.EvalAssignment{
					{
						Field:      "is_slow",
						Expression: &ast.Literal{Value: true, LiteralTyp: ast.LiteralTypeBool},
					},
				},
			},
			&ast.WhereCommand{
				Condition: &ast.BinaryExpression{
					Left:     &ast.FieldReference{Name: "is_slow"},
					Operator: "=",
					Right:    &ast.Literal{Value: true, LiteralTyp: ast.LiteralTypeBool},
				},
			},
			&ast.DedupCommand{
				Fields: []ast.Expression{&ast.FieldReference{Name: "host"}},
				Count:  1,
			},
			&ast.TopCommand{
				Fields: []ast.Expression{&ast.FieldReference{Name: "status"}},
				Limit:  5,
			},
		},
	}

	plan, err := builder.Build(query)
	require.NoError(t, err)
	require.NotNil(t, plan)

	// Verify pipeline structure: Top -> Dedup -> Filter -> Eval -> Scan
	top, ok := plan.(*LogicalTop)
	assert.True(t, ok, "Expected LogicalTop at top")

	dedup, ok := top.Input.(*LogicalDedup)
	assert.True(t, ok, "Expected LogicalDedup")

	filter, ok := dedup.Input.(*LogicalFilter)
	assert.True(t, ok, "Expected LogicalFilter")

	eval, ok := filter.Input.(*LogicalEval)
	assert.True(t, ok, "Expected LogicalEval")

	scan, ok := eval.Input.(*LogicalScan)
	assert.True(t, ok, "Expected LogicalScan at bottom")
	assert.Equal(t, "logs", scan.Source)
}
