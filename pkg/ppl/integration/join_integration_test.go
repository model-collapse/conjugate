// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package integration

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/parser"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"github.com/conjugate/conjugate/pkg/ppl/planner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoinCommand_Integration(t *testing.T) {
	t.Run("BasicInnerJoin", func(t *testing.T) {
		// Query: join orders and users on user_id
		query := `search source=orders | join user_id [search source=users]`

		// Parse query
		p := parser.NewParser()
		queryAST, err := p.Parse(query)
		require.NoError(t, err)
		require.NotNil(t, queryAST)

		// Setup schema for orders
		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("user_id", analyzer.FieldTypeString)
		schema.AddField("amount", analyzer.FieldTypeFloat)

		// Analyze
		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(queryAST)
		require.NoError(t, err)

		// Build logical plan
		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(queryAST)
		require.NoError(t, err)
		require.NotNil(t, logicalPlan)

		// Verify logical plan structure: Join -> Scan (orders)
		joinPlan, ok := logicalPlan.(*planner.LogicalJoin)
		require.True(t, ok, "Expected LogicalJoin at root")
		assert.Equal(t, ast.JoinTypeInner, joinPlan.JoinType)
		assert.Equal(t, "user_id", joinPlan.JoinField)
		assert.Equal(t, "user_id", joinPlan.RightField)

		// Verify main input is orders scan
		scanPlan, ok := joinPlan.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan as input")
		assert.Equal(t, "orders", scanPlan.Source)

		// Verify right side is users scan
		rightScan, ok := joinPlan.Right.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan as right side")
		assert.Equal(t, "users", rightScan.Source)

		// Build physical plan
		physicalPlanner := physical.NewPhysicalPlanner()
		physicalPlan, err := physicalPlanner.Plan(logicalPlan)
		require.NoError(t, err)
		require.NotNil(t, physicalPlan)

		// Verify physical plan structure
		physJoin, ok := physicalPlan.(*physical.PhysicalJoin)
		require.True(t, ok, "Expected PhysicalJoin")
		assert.Equal(t, ast.JoinTypeInner, physJoin.JoinType)
		assert.Equal(t, physical.ExecuteOnCoordinator, physJoin.Location())

		// Verify input
		physScan, ok := physJoin.Input.(*physical.PhysicalScan)
		require.True(t, ok, "Expected PhysicalScan as input")
		assert.Equal(t, "orders", physScan.Source)

		// Verify right side
		physRight, ok := physJoin.Right.(*physical.PhysicalScan)
		require.True(t, ok, "Expected PhysicalScan as right side")
		assert.Equal(t, "users", physRight.Source)
	})

	t.Run("LeftJoin", func(t *testing.T) {
		// Query with left join type
		query := `search source=orders | join type=left user_id [search source=users]`

		p := parser.NewParser()
		queryAST, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("user_id", analyzer.FieldTypeString)
		schema.AddField("amount", analyzer.FieldTypeFloat)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(queryAST)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(queryAST)
		require.NoError(t, err)

		// Verify join type is LEFT
		joinPlan, ok := logicalPlan.(*planner.LogicalJoin)
		require.True(t, ok, "Expected LogicalJoin")
		assert.Equal(t, ast.JoinTypeLeft, joinPlan.JoinType)
	})

	t.Run("JoinWithProcessingCommands", func(t *testing.T) {
		// Query with processing commands in both main and subsearch
		query := `search source=orders | where amount > 100 | join user_id [search source=users | fields user_id, name]`

		p := parser.NewParser()
		queryAST, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("user_id", analyzer.FieldTypeString)
		schema.AddField("amount", analyzer.FieldTypeFloat)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(queryAST)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(queryAST)
		require.NoError(t, err)

		// Verify plan structure: Join -> Filter -> Scan
		joinPlan, ok := logicalPlan.(*planner.LogicalJoin)
		require.True(t, ok, "Expected LogicalJoin at root")

		// Main query has filter
		filterPlan, ok := joinPlan.Input.(*planner.LogicalFilter)
		require.True(t, ok, "Expected LogicalFilter in main query")

		// Filter has scan as input
		_, ok = filterPlan.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan after filter")

		// Right side has projection
		projectPlan, ok := joinPlan.Right.(*planner.LogicalProject)
		require.True(t, ok, "Expected LogicalProject in right side")

		// Projection has scan as input
		_, ok = projectPlan.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan after projection")
	})

	t.Run("JoinFollowedByProjection", func(t *testing.T) {
		// Query with projection after join
		// Note: Only project fields that exist in the initial schema for semantic analysis
		query := `search source=orders | join user_id [search source=users] | fields order_id, user_id, amount`

		p := parser.NewParser()
		queryAST, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("user_id", analyzer.FieldTypeString)
		schema.AddField("amount", analyzer.FieldTypeFloat)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(queryAST)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(queryAST)
		require.NoError(t, err)

		// Verify plan structure: Project -> Join -> Scan
		projectPlan, ok := logicalPlan.(*planner.LogicalProject)
		require.True(t, ok, "Expected LogicalProject at root")

		// Join below projection
		joinPlan, ok := projectPlan.Input.(*planner.LogicalJoin)
		require.True(t, ok, "Expected LogicalJoin as input to projection")
		assert.Equal(t, ast.JoinTypeInner, joinPlan.JoinType)
	})

	t.Run("JoinWithEval", func(t *testing.T) {
		// Query with eval command before join
		query := `search source=orders | eval order_type="standard" | join user_id [search source=users]`

		p := parser.NewParser()
		queryAST, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("user_id", analyzer.FieldTypeString)
		schema.AddField("amount", analyzer.FieldTypeFloat)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(queryAST)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(queryAST)
		require.NoError(t, err)

		// Verify plan structure: Join -> Eval -> Scan
		joinPlan, ok := logicalPlan.(*planner.LogicalJoin)
		require.True(t, ok, "Expected LogicalJoin at root")

		// Eval before join
		evalPlan, ok := joinPlan.Input.(*planner.LogicalEval)
		require.True(t, ok, "Expected LogicalEval in main query")

		// Scan as input to eval
		_, ok = evalPlan.Input.(*planner.LogicalScan)
		require.True(t, ok, "Expected LogicalScan after eval")
	})

	t.Run("MultipleJoins", func(t *testing.T) {
		// Query with multiple joins
		query := `search source=orders | join user_id [search source=users] | join product_id [search source=products]`

		p := parser.NewParser()
		queryAST, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("user_id", analyzer.FieldTypeString)
		schema.AddField("product_id", analyzer.FieldTypeString)
		schema.AddField("amount", analyzer.FieldTypeFloat)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(queryAST)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(queryAST)
		require.NoError(t, err)

		// Verify plan structure: Join -> Join -> Scan
		outerJoin, ok := logicalPlan.(*planner.LogicalJoin)
		require.True(t, ok, "Expected LogicalJoin at root")
		assert.Equal(t, "product_id", outerJoin.JoinField)

		// Outer join's right side is products
		productsRight, ok := outerJoin.Right.(*planner.LogicalScan)
		require.True(t, ok)
		assert.Equal(t, "products", productsRight.Source)

		// Input is another join
		innerJoin, ok := outerJoin.Input.(*planner.LogicalJoin)
		require.True(t, ok, "Expected LogicalJoin as input")
		assert.Equal(t, "user_id", innerJoin.JoinField)

		// Inner join's right side is users
		usersRight, ok := innerJoin.Right.(*planner.LogicalScan)
		require.True(t, ok)
		assert.Equal(t, "users", usersRight.Source)

		// Inner join's input is orders
		ordersScan, ok := innerJoin.Input.(*planner.LogicalScan)
		require.True(t, ok)
		assert.Equal(t, "orders", ordersScan.Source)
	})

	t.Run("RightJoinType", func(t *testing.T) {
		// Query with right join type
		query := `search source=orders | join type=right user_id [search source=users]`

		p := parser.NewParser()
		queryAST, err := p.Parse(query)
		require.NoError(t, err)

		schema := analyzer.NewSchema("orders")
		schema.AddField("order_id", analyzer.FieldTypeString)
		schema.AddField("user_id", analyzer.FieldTypeString)

		semanticAnalyzer := analyzer.NewAnalyzer(schema)
		err = semanticAnalyzer.Analyze(queryAST)
		require.NoError(t, err)

		planBuilder := planner.NewPlanBuilder(schema)
		logicalPlan, err := planBuilder.Build(queryAST)
		require.NoError(t, err)

		joinPlan, ok := logicalPlan.(*planner.LogicalJoin)
		require.True(t, ok)
		assert.Equal(t, ast.JoinTypeRight, joinPlan.JoinType)
	})
}
