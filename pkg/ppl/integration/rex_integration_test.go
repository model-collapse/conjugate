// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package integration

import (
	"context"
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/executor"
	"github.com/conjugate/conjugate/pkg/ppl/parser"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"github.com/conjugate/conjugate/pkg/ppl/planner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRexCommand_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	t.Run("BasicRexExtraction", func(t *testing.T) {
		// Query: search source=logs | rex field=message "(?P<code>\d{3}): (?P<msg>.*)"
		query := `search source=logs | rex field=message "(?P<code>\d{3}): (?P<msg>.*)"`

		// Parse query
		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)
		require.NotNil(t, ast)

		// Analyze
		schema := analyzer.NewSchema("logs")
		schema.AddField("message", analyzer.FieldTypeString)
		schema.AddField("timestamp", analyzer.FieldTypeDate)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		// Build logical plan
		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)
		require.NotNil(t, logicalPlan)

		// Verify logical plan structure
		rexPlan, ok := logicalPlan.(*planner.LogicalRex)
		require.True(t, ok, "Expected LogicalRex at root")
		assert.Equal(t, "message", rexPlan.SourceField)
		assert.Contains(t, rexPlan.Pattern, "code")
		assert.Contains(t, rexPlan.Pattern, "msg")
		assert.Len(t, rexPlan.ExtractedFields, 2)
		assert.Contains(t, rexPlan.ExtractedFields, "code")
		assert.Contains(t, rexPlan.ExtractedFields, "msg")

		// Check output schema includes new fields
		outputSchema := rexPlan.Schema()
		_, err = outputSchema.GetField("code")
		assert.NoError(t, err, "code field should be in output schema")
		_, err = outputSchema.GetField("msg")
		assert.NoError(t, err, "msg field should be in output schema")

		// Build physical plan
		physicalPlanner := physical.NewPhysicalPlanner()
		physicalPlan, err := physicalPlanner.Plan(logicalPlan)
		require.NoError(t, err)
		require.NotNil(t, physicalPlan)

		// Verify physical plan structure
		physRex, ok := physicalPlan.(*physical.PhysicalRex)
		require.True(t, ok, "Expected PhysicalRex")
		assert.Equal(t, "message", physRex.SourceField)
		assert.Equal(t, physical.ExecuteOnCoordinator, physRex.Location())

		// Create test data
		testRows := []*executor.Row{
			executor.NewRow(map[string]interface{}{
				"message":   "404: Not Found",
				"timestamp": "2024-01-15T10:00:00Z",
			}),
			executor.NewRow(map[string]interface{}{
				"message":   "500: Internal Server Error",
				"timestamp": "2024-01-15T10:01:00Z",
			}),
			executor.NewRow(map[string]interface{}{
				"message":   "invalid message", // Won't match
				"timestamp": "2024-01-15T10:02:00Z",
			}),
		}

		// Build operator (manually since we're injecting mock data)
		mockScan := executor.NewSliceIterator(testRows)

		// Execute
		rexOp, err := executor.NewRexOperator(
			mockScan,
			physRex.SourceField,
			physRex.Pattern,
			physRex.ExtractedFields,
			logger,
		)
		require.NoError(t, err)

		err = rexOp.Open(ctx)
		require.NoError(t, err)

		// First row - should have code and msg extracted
		row, err := rexOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		code, exists := row.Get("code")
		assert.True(t, exists)
		assert.Equal(t, "404", code)

		msg, exists := row.Get("msg")
		assert.True(t, exists)
		assert.Equal(t, "Not Found", msg)

		// Original fields should still be present
		message, exists := row.Get("message")
		assert.True(t, exists)
		assert.Equal(t, "404: Not Found", message)

		// Second row
		row, err = rexOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		code, _ = row.Get("code")
		assert.Equal(t, "500", code)

		msg, _ = row.Get("msg")
		assert.Equal(t, "Internal Server Error", msg)

		// Third row - no match, should not have extracted fields
		row, err = rexOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		_, exists = row.Get("code")
		assert.False(t, exists)

		_, exists = row.Get("msg")
		assert.False(t, exists)

		err = rexOp.Close()
		require.NoError(t, err)
	})

	t.Run("RexWithDefaultField", func(t *testing.T) {
		// Query without field parameter (should default to _raw)
		query := `search source=logs | rex "user=(?P<user>\w+)"`

		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("logs")
		schema.AddField("_raw", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)

		rexPlan, ok := logicalPlan.(*planner.LogicalRex)
		require.True(t, ok)
		assert.Equal(t, "_raw", rexPlan.SourceField)
		assert.Contains(t, rexPlan.ExtractedFields, "user")
	})

	t.Run("RexInPipeline", func(t *testing.T) {
		// Query with rex followed by projection
		// Note: fields command referencing extracted fields requires two-phase analysis
		query := `search source=logs | rex field=message "user=(?P<user>\w+)" | fields message`

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

		// Verify plan structure: Project -> Rex -> Scan
		projectPlan, ok := logicalPlan.(*planner.LogicalProject)
		require.True(t, ok, "Expected LogicalProject at root")

		rexPlan, ok := projectPlan.Input.(*planner.LogicalRex)
		require.True(t, ok, "Expected LogicalRex under Project")
		assert.Equal(t, "message", rexPlan.SourceField)

		scanPlan, ok := rexPlan.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan under Rex")
		assert.Equal(t, "logs", scanPlan.Source)
	})

	t.Run("MultipleRexCommands", func(t *testing.T) {
		// Test chaining multiple rex commands
		query := `search source=logs | rex field=message "user=(?P<user>\w+)" | rex field=message "action=(?P<action>\w+)"`

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

		// Verify plan structure: Rex -> Rex -> Scan
		rexPlan1, ok := logicalPlan.(*planner.LogicalRex)
		require.True(t, ok, "Expected LogicalRex at root")
		assert.Contains(t, rexPlan1.ExtractedFields, "action")

		rexPlan2, ok := rexPlan1.Input.(*planner.LogicalRex)
		require.True(t, ok, "Expected LogicalRex as input")
		assert.Contains(t, rexPlan2.ExtractedFields, "user")

		scanPlan, ok := rexPlan2.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan under second Rex")
		assert.Equal(t, "logs", scanPlan.Source)
	})
}
