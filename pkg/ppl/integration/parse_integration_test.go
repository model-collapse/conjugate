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

func TestParseCommand_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	t.Run("BasicLogParsing", func(t *testing.T) {
		// Query: search source=logs | parse message "user (?P<username>\w+) from (?P<ip>\d+\.\d+\.\d+\.\d+)"
		query := `search source=logs | parse message "user (?P<username>\w+) from (?P<ip>\d+\.\d+\.\d+\.\d+)"`

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
		parsePlan, ok := logicalPlan.(*planner.LogicalParse)
		require.True(t, ok, "Expected LogicalParse at root")
		assert.Equal(t, "message", parsePlan.SourceField)
		assert.Contains(t, parsePlan.Pattern, "username")
		assert.Contains(t, parsePlan.Pattern, "ip")
		assert.Len(t, parsePlan.ExtractedFields, 2)
		assert.Contains(t, parsePlan.ExtractedFields, "username")
		assert.Contains(t, parsePlan.ExtractedFields, "ip")

		// Check output schema includes new fields
		outputSchema := parsePlan.Schema()
		_, err = outputSchema.GetField("username")
		assert.NoError(t, err, "username field should be in output schema")
		_, err = outputSchema.GetField("ip")
		assert.NoError(t, err, "ip field should be in output schema")

		// Build physical plan
		physicalPlanner := physical.NewPhysicalPlanner()
		physicalPlan, err := physicalPlanner.Plan(logicalPlan)
		require.NoError(t, err)
		require.NotNil(t, physicalPlan)

		// Verify physical plan structure
		physParse, ok := physicalPlan.(*physical.PhysicalParse)
		require.True(t, ok, "Expected PhysicalParse")
		assert.Equal(t, "message", physParse.SourceField)
		assert.Equal(t, physical.ExecuteOnCoordinator, physParse.Location())

		// Create test data
		testRows := []*executor.Row{
			executor.NewRow(map[string]interface{}{
				"message":   "user john from 192.168.1.100",
				"timestamp": "2024-01-15T10:00:00Z",
			}),
			executor.NewRow(map[string]interface{}{
				"message":   "user jane from 10.0.0.5",
				"timestamp": "2024-01-15T10:01:00Z",
			}),
			executor.NewRow(map[string]interface{}{
				"message":   "error: connection failed", // Won't match
				"timestamp": "2024-01-15T10:02:00Z",
			}),
		}

		// Build operator (manually since we're injecting mock data)
		mockScan := executor.NewSliceIterator(testRows)

		// Execute
		parseOp, err := executor.NewParseOperator(
			mockScan,
			physParse.SourceField,
			physParse.Pattern,
			physParse.ExtractedFields,
			logger,
		)
		require.NoError(t, err)

		err = parseOp.Open(ctx)
		require.NoError(t, err)

		// First row - should have username and ip extracted
		row, err := parseOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		username, exists := row.Get("username")
		assert.True(t, exists)
		assert.Equal(t, "john", username)

		ip, exists := row.Get("ip")
		assert.True(t, exists)
		assert.Equal(t, "192.168.1.100", ip)

		// Original fields should still be present
		message, exists := row.Get("message")
		assert.True(t, exists)
		assert.Equal(t, "user john from 192.168.1.100", message)

		// Second row
		row, err = parseOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		username, _ = row.Get("username")
		assert.Equal(t, "jane", username)

		ip, _ = row.Get("ip")
		assert.Equal(t, "10.0.0.5", ip)

		// Third row - no match, should not have extracted fields
		row, err = parseOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		_, exists = row.Get("username")
		assert.False(t, exists)

		_, exists = row.Get("ip")
		assert.False(t, exists)

		err = parseOp.Close()
		require.NoError(t, err)
	})

	t.Run("ComplexPatternWithMultipleFields", func(t *testing.T) {
		// Query with complex regex pattern
		query := `search source=logs | parse log "\[(?P<timestamp>[^\]]+)\] (?P<level>\w+) (?P<message>.+)"`

		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("logs")
		schema.AddField("log", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)

		parsePlan, ok := logicalPlan.(*planner.LogicalParse)
		require.True(t, ok)
		assert.Equal(t, "log", parsePlan.SourceField)
		assert.Len(t, parsePlan.ExtractedFields, 3)

		// Test execution
		testRows := []*executor.Row{
			executor.NewRow(map[string]interface{}{
				"log": "[2024-01-15 10:30:45] ERROR Connection timeout",
			}),
		}

		mockScan := executor.NewSliceIterator(testRows)
		parseOp, err := executor.NewParseOperator(
			mockScan,
			parsePlan.SourceField,
			parsePlan.Pattern,
			parsePlan.ExtractedFields,
			logger,
		)
		require.NoError(t, err)

		err = parseOp.Open(ctx)
		require.NoError(t, err)

		row, err := parseOp.Next(ctx)
		require.NoError(t, err)

		timestamp, _ := row.Get("timestamp")
		assert.Equal(t, "2024-01-15 10:30:45", timestamp)

		level, _ := row.Get("level")
		assert.Equal(t, "ERROR", level)

		message, _ := row.Get("message")
		assert.Equal(t, "Connection timeout", message)

		err = parseOp.Close()
		require.NoError(t, err)
	})

	t.Run("ParseWithFieldParameter", func(t *testing.T) {
		// Query with field parameter syntax
		query := `search source=logs | parse field=raw_log "status=(?P<status>\d+)"`

		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("logs")
		schema.AddField("raw_log", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)

		parsePlan, ok := logicalPlan.(*planner.LogicalParse)
		require.True(t, ok)
		assert.Equal(t, "raw_log", parsePlan.SourceField)
		assert.Contains(t, parsePlan.ExtractedFields, "status")
	})

	t.Run("ParseInPipeline", func(t *testing.T) {
		// Query with parse followed by projection
		// Note: fields command referencing extracted fields requires two-phase analysis
		query := `search source=logs | parse message "user (?P<user>\w+)" | fields message`

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

		// Verify plan structure: Project -> Parse -> Scan
		projectPlan, ok := logicalPlan.(*planner.LogicalProject)
		require.True(t, ok, "Expected LogicalProject at root")

		parsePlan, ok := projectPlan.Input.(*planner.LogicalParse)
		require.True(t, ok, "Expected LogicalParse under Project")
		assert.Equal(t, "message", parsePlan.SourceField)

		scanPlan, ok := parsePlan.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan under Parse")
		assert.Equal(t, "logs", scanPlan.Source)
	})
}
