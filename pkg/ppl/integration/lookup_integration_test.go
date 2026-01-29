// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package integration

import (
	"context"
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/executor"
	"github.com/conjugate/conjugate/pkg/ppl/lookup"
	"github.com/conjugate/conjugate/pkg/ppl/parser"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"github.com/conjugate/conjugate/pkg/ppl/planner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLookupCommand_Integration(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	t.Run("BasicLookupEnrichment", func(t *testing.T) {
		// Query: search source=orders | lookup products product_id output name, price
		query := `search source=orders | lookup products product_id output name, price`

		// Parse query
		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)
		require.NotNil(t, ast)

		// Analyze
		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("product_id", analyzer.FieldTypeString)
		schema.AddField("quantity", analyzer.FieldTypeInt)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		// Build logical plan
		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)
		require.NotNil(t, logicalPlan)

		// Verify logical plan structure
		lookupPlan, ok := logicalPlan.(*planner.LogicalLookup)
		require.True(t, ok, "Expected LogicalLookup at root")
		assert.Equal(t, "products", lookupPlan.TableName)
		assert.Equal(t, "product_id", lookupPlan.JoinField)
		assert.Len(t, lookupPlan.OutputFields, 2)
		assert.Contains(t, lookupPlan.OutputFields, "name")
		assert.Contains(t, lookupPlan.OutputFields, "price")

		// Check output schema includes new fields
		outputSchema := lookupPlan.Schema()
		_, err = outputSchema.GetField("name")
		assert.NoError(t, err, "name field should be in output schema")
		_, err = outputSchema.GetField("price")
		assert.NoError(t, err, "price field should be in output schema")

		// Build physical plan
		physicalPlanner := physical.NewPhysicalPlanner()
		physicalPlan, err := physicalPlanner.Plan(logicalPlan)
		require.NoError(t, err)
		require.NotNil(t, physicalPlan)

		// Verify physical plan structure
		physLookup, ok := physicalPlan.(*physical.PhysicalLookup)
		require.True(t, ok, "Expected PhysicalLookup")
		assert.Equal(t, "products", physLookup.TableName)
		assert.Equal(t, physical.ExecuteOnCoordinator, physLookup.Location())

		// Create lookup registry and table
		registry := lookup.NewRegistry(logger)
		table := lookup.NewLookupTable("products", logger)

		table.AddRow("101", map[string]interface{}{
			"product_id": "101",
			"name":       "Laptop",
			"price":      "999.99",
		})
		table.AddRow("102", map[string]interface{}{
			"product_id": "102",
			"name":       "Mouse",
			"price":      "29.99",
		})

		err = registry.Register(table)
		require.NoError(t, err)

		// Create test data
		testRows := []*executor.Row{
			executor.NewRow(map[string]interface{}{
				"order_id":   "1",
				"product_id": "101",
				"quantity":   2,
			}),
			executor.NewRow(map[string]interface{}{
				"order_id":   "2",
				"product_id": "102",
				"quantity":   5,
			}),
			executor.NewRow(map[string]interface{}{
				"order_id":   "3",
				"product_id": "999", // Doesn't exist in lookup
				"quantity":   1,
			}),
		}

		// Build operator
		mockScan := executor.NewSliceIterator(testRows)

		// Execute
		lookupOp, err := executor.NewLookupOperator(
			mockScan,
			registry,
			physLookup.TableName,
			physLookup.JoinField,
			physLookup.JoinFieldAlias,
			physLookup.OutputFields,
			physLookup.OutputAliases,
			logger,
		)
		require.NoError(t, err)

		err = lookupOp.Open(ctx)
		require.NoError(t, err)

		// First row - should have name and price added
		row, err := lookupOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		name, exists := row.Get("name")
		assert.True(t, exists)
		assert.Equal(t, "Laptop", name)

		price, exists := row.Get("price")
		assert.True(t, exists)
		assert.Equal(t, "999.99", price)

		// Original fields should still be present
		orderID, exists := row.Get("order_id")
		assert.True(t, exists)
		assert.Equal(t, "1", orderID)

		// Second row
		row, err = lookupOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		name, _ = row.Get("name")
		assert.Equal(t, "Mouse", name)

		price, _ = row.Get("price")
		assert.Equal(t, "29.99", price)

		// Third row - no match, should not have lookup fields
		row, err = lookupOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		_, exists = row.Get("name")
		assert.False(t, exists)

		_, exists = row.Get("price")
		assert.False(t, exists)

		err = lookupOp.Close()
		require.NoError(t, err)
	})

	t.Run("LookupWithAliases", func(t *testing.T) {
		// Query with field aliases
		query := `search source=events | lookup users user_id as uid output username as user, email as contact`

		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("events")
		schema.AddField("event_id", analyzer.FieldTypeString)
		schema.AddField("user_id", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)

		lookupPlan, ok := logicalPlan.(*planner.LogicalLookup)
		require.True(t, ok)
		assert.Equal(t, "users", lookupPlan.TableName)
		assert.Equal(t, "user_id", lookupPlan.JoinField)
		assert.Equal(t, "uid", lookupPlan.JoinFieldAlias)
		assert.Contains(t, lookupPlan.OutputFields, "username")
		assert.Contains(t, lookupPlan.OutputAliases, "user")
		assert.Contains(t, lookupPlan.OutputAliases, "contact")
	})

	t.Run("LookupInPipeline", func(t *testing.T) {
		// Query with lookup followed by projection
		query := `search source=orders | lookup products product_id output name | fields order_id, name`

		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("product_id", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)

		// Verify plan structure: Project -> Lookup -> Scan
		projectPlan, ok := logicalPlan.(*planner.LogicalProject)
		require.True(t, ok, "Expected LogicalProject at root")

		lookupPlan, ok := projectPlan.Input.(*planner.LogicalLookup)
		require.True(t, ok, "Expected LogicalLookup under Project")
		assert.Equal(t, "products", lookupPlan.TableName)

		scanPlan, ok := lookupPlan.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan under Lookup")
		assert.Equal(t, "orders", scanPlan.Source)
	})

	t.Run("MultipleLookups", func(t *testing.T) {
		// Query with multiple lookup commands
		query := `search source=orders | lookup products product_id output name | lookup customers customer_id output customer_name`

		p := parser.NewParser()
		ast, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("product_id", analyzer.FieldTypeString)
		schema.AddField("customer_id", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(ast)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(ast)
		require.NoError(t, err)

		// Verify plan structure: Lookup -> Lookup -> Scan
		lookupPlan1, ok := logicalPlan.(*planner.LogicalLookup)
		require.True(t, ok, "Expected LogicalLookup at root")
		assert.Equal(t, "customers", lookupPlan1.TableName)

		lookupPlan2, ok := lookupPlan1.Input.(*planner.LogicalLookup)
		require.True(t, ok, "Expected LogicalLookup as input")
		assert.Equal(t, "products", lookupPlan2.TableName)

		scanPlan, ok := lookupPlan2.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan under second Lookup")
		assert.Equal(t, "orders", scanPlan.Source)
	})
}
