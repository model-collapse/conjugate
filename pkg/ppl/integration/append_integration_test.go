// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package integration

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/parser"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"github.com/conjugate/conjugate/pkg/ppl/planner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendCommand_Integration(t *testing.T) {
	t.Run("BasicAppend", func(t *testing.T) {
		// Query: search source=logs_2024 | append [search source=logs_2023]
		query := `search source=logs_2024 | append [search source=logs_2023]`

		// Parse query
		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)
		require.NotNil(t, ast)

		// Analyze
		schema := analyzer.NewSchema("logs_2024")
		schema.AddField("timestamp", analyzer.FieldTypeString)
		schema.AddField("message", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		// Build logical plan
		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)
		require.NotNil(t, logicalPlan)

		// Verify logical plan structure: Append -> Scan (main)
		appendPlan, ok := logicalPlan.(*planner.LogicalAppend)
		require.True(t, ok, "Expected LogicalAppend at root")

		// Verify main input is a scan
		scanPlan, ok := appendPlan.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan as input")
		assert.Equal(t, "logs_2024", scanPlan.Source)

		// Verify subsearch is also a scan
		subsearchScan, ok := appendPlan.Subsearch.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan as subsearch")
		assert.Equal(t, "logs_2023", subsearchScan.Source)

		// Build physical plan
		physicalPlanner := physical.NewPhysicalPlanner()
		physicalPlan, err := physicalPlanner.Plan(logicalPlan)
		require.NoError(t, err)
		require.NotNil(t, physicalPlan)

		// Verify physical plan structure
		physAppend, ok := physicalPlan.(*physical.PhysicalAppend)
		require.True(t, ok, "Expected PhysicalAppend")
		assert.Equal(t, physical.ExecuteOnCoordinator, physAppend.Location())

		// Verify main input
		physScan, ok := physAppend.Input.(*physical.PhysicalScan)
		require.True(t, ok, "Expected PhysicalScan as input")
		assert.Equal(t, "logs_2024", physScan.Source)

		// Verify subsearch
		physSubsearch, ok := physAppend.Subsearch.(*physical.PhysicalScan)
		require.True(t, ok, "Expected PhysicalScan as subsearch")
		assert.Equal(t, "logs_2023", physSubsearch.Source)
	})

	t.Run("AppendWithProcessingCommands", func(t *testing.T) {
		// Query with processing commands in both main and subsearch
		query := `search source=logs_2024 | fields timestamp, message | append [search source=logs_2023 | fields timestamp, message]`

		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("logs_2024")
		schema.AddField("timestamp", analyzer.FieldTypeString)
		schema.AddField("message", analyzer.FieldTypeString)
		schema.AddField("level", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)

		// Verify plan structure: Append -> Project -> Scan
		appendPlan, ok := logicalPlan.(*planner.LogicalAppend)
		require.True(t, ok, "Expected LogicalAppend at root")

		// Main query has projection
		_, ok = appendPlan.Input.(*planner.LogicalProject)
		require.True(t, ok, "Expected LogicalProject in main query")

		// Subsearch also has projection
		_, ok = appendPlan.Subsearch.(*planner.LogicalProject)
		require.True(t, ok, "Expected LogicalProject in subsearch")
	})

	t.Run("AppendWithEval", func(t *testing.T) {
		// Query with eval command
		query := `search source=logs | eval status="main" | append [search source=archive | eval status="archive"]`

		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("logs")
		schema.AddField("message", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)

		// Verify it's an append command
		_, ok := logicalPlan.(*planner.LogicalAppend)
		require.True(t, ok, "Expected LogicalAppend")
	})

	t.Run("MultipleAppends", func(t *testing.T) {
		// Query with multiple append commands
		query := `search source=logs_2024 | append [search source=logs_2023] | append [search source=logs_2022]`

		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("logs_2024")
		schema.AddField("timestamp", analyzer.FieldTypeString)
		schema.AddField("message", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)

		// Verify plan structure: Append -> Append -> Scan
		outerAppend, ok := logicalPlan.(*planner.LogicalAppend)
		require.True(t, ok, "Expected LogicalAppend at root")

		// First append's subsearch is logs_2022
		subsearch1, ok := outerAppend.Subsearch.(*planner.LogicalScan)
		require.True(t, ok)
		assert.Equal(t, "logs_2022", subsearch1.Source)

		// Input is another append
		innerAppend, ok := outerAppend.Input.(*planner.LogicalAppend)
		require.True(t, ok, "Expected LogicalAppend as input")

		// Inner append's subsearch is logs_2023
		subsearch2, ok := innerAppend.Subsearch.(*planner.LogicalScan)
		require.True(t, ok)
		assert.Equal(t, "logs_2023", subsearch2.Source)

		// Inner append's input is logs_2024
		mainScan, ok := innerAppend.Input.(*planner.LogicalScan)
		require.True(t, ok)
		assert.Equal(t, "logs_2024", mainScan.Source)
	})
}
