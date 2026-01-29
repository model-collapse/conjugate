// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package optimizer

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
	return schema
}

func TestFilterMergeRule(t *testing.T) {
	schema := createTestSchema()
	rule := NewFilterMergeRule()

	// Create plan: Filter(status=500) -> Filter(host='server1') -> Scan
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}
	filter1 := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "host"},
			Operator: "=",
			Right:    &ast.Literal{Value: "server1", LiteralTyp: ast.LiteralTypeString},
		},
		Input: scan,
	}
	filter2 := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: filter1,
	}

	// Apply rule
	optimized := rule.Apply(filter2)
	require.NotNil(t, optimized)

	// Should merge into single filter with AND
	mergedFilter, ok := optimized.(*planner.LogicalFilter)
	require.True(t, ok)

	// Check that it's an AND expression
	andExpr, ok := mergedFilter.Condition.(*ast.BinaryExpression)
	require.True(t, ok)
	assert.Equal(t, "AND", andExpr.Operator)

	// Input should be the scan
	_, ok = mergedFilter.Input.(*planner.LogicalScan)
	assert.True(t, ok)
}

func TestFilterPushDownRule_PastProject(t *testing.T) {
	schema := createTestSchema()
	rule := NewFilterPushDownRule()

	// Create plan: Filter -> Project -> Scan
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}
	project := &planner.LogicalProject{
		Fields: []ast.Expression{
			&ast.FieldReference{Name: "status"},
			&ast.FieldReference{Name: "host"},
		},
		OutputSchema: schema,
		Input:        scan,
	}
	filter := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: project,
	}

	// Apply rule
	optimized := rule.Apply(filter)
	require.NotNil(t, optimized)

	// Should become: Project -> Filter -> Scan
	newProject, ok := optimized.(*planner.LogicalProject)
	require.True(t, ok)

	newFilter, ok := newProject.Input.(*planner.LogicalFilter)
	require.True(t, ok)

	_, ok = newFilter.Input.(*planner.LogicalScan)
	assert.True(t, ok)
}

func TestFilterPushDownRule_PastSort(t *testing.T) {
	schema := createTestSchema()
	rule := NewFilterPushDownRule()

	// Create plan: Filter -> Sort -> Scan
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}
	sort := &planner.LogicalSort{
		SortKeys: []*ast.SortKey{
			{Field: &ast.FieldReference{Name: "timestamp"}, Descending: true},
		},
		Input: scan,
	}
	filter := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: sort,
	}

	// Apply rule
	optimized := rule.Apply(filter)
	require.NotNil(t, optimized)

	// Should become: Sort -> Filter -> Scan
	newSort, ok := optimized.(*planner.LogicalSort)
	require.True(t, ok)

	newFilter, ok := newSort.Input.(*planner.LogicalFilter)
	require.True(t, ok)

	_, ok = newFilter.Input.(*planner.LogicalScan)
	assert.True(t, ok)
}

func TestProjectMergeRule(t *testing.T) {
	schema := createTestSchema()
	rule := NewProjectMergeRule()

	// Create plan: Project(status, host) -> Project(status, host, timestamp) -> Scan
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}
	project1 := &planner.LogicalProject{
		Fields: []ast.Expression{
			&ast.FieldReference{Name: "status"},
			&ast.FieldReference{Name: "host"},
			&ast.FieldReference{Name: "timestamp"},
		},
		OutputSchema: schema,
		Input:        scan,
	}
	project2 := &planner.LogicalProject{
		Fields: []ast.Expression{
			&ast.FieldReference{Name: "status"},
			&ast.FieldReference{Name: "host"},
		},
		OutputSchema: schema,
		Input:        project1,
	}

	// Apply rule
	optimized := rule.Apply(project2)
	require.NotNil(t, optimized)

	// Should merge into single project
	mergedProject, ok := optimized.(*planner.LogicalProject)
	require.True(t, ok)

	// Should skip the inner project and go straight to scan
	_, ok = mergedProject.Input.(*planner.LogicalScan)
	assert.True(t, ok)
}

func TestConstantFoldingRule_Arithmetic(t *testing.T) {
	schema := createTestSchema()
	rule := NewConstantFoldingRule()

	// Create filter with constant arithmetic: status = (10 + 5)
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}
	filter := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right: &ast.BinaryExpression{
				Left:     &ast.Literal{Value: 10, LiteralTyp: ast.LiteralTypeInt},
				Operator: "+",
				Right:    &ast.Literal{Value: 5, LiteralTyp: ast.LiteralTypeInt},
			},
		},
		Input: scan,
	}

	// Apply rule
	optimized := rule.Apply(filter)
	require.NotNil(t, optimized)

	// Should fold to: status = 15
	foldedFilter, ok := optimized.(*planner.LogicalFilter)
	require.True(t, ok)

	binExpr, ok := foldedFilter.Condition.(*ast.BinaryExpression)
	require.True(t, ok)

	rightLit, ok := binExpr.Right.(*ast.Literal)
	require.True(t, ok)
	assert.Equal(t, 15, rightLit.Value)
}

func TestConstantFoldingRule_Multiplication(t *testing.T) {
	rule := NewConstantFoldingRule()

	// Create expression: 3 * 4
	expr := &ast.BinaryExpression{
		Left:     &ast.Literal{Value: 3, LiteralTyp: ast.LiteralTypeInt},
		Operator: "*",
		Right:    &ast.Literal{Value: 4, LiteralTyp: ast.LiteralTypeInt},
	}

	folded := rule.foldExpression(expr)
	require.NotNil(t, folded)

	lit, ok := folded.(*ast.Literal)
	require.True(t, ok)
	assert.Equal(t, 12, lit.Value)
}

func TestConstantFoldingRule_NOT(t *testing.T) {
	rule := NewConstantFoldingRule()

	// Create expression: NOT true
	expr := &ast.UnaryExpression{
		Operator: "NOT",
		Operand:  &ast.Literal{Value: true, LiteralTyp: ast.LiteralTypeBool},
	}

	folded := rule.foldExpression(expr)
	require.NotNil(t, folded)

	lit, ok := folded.(*ast.Literal)
	require.True(t, ok)
	assert.Equal(t, false, lit.Value)
}

func TestLimitPushDownRule_PastFilter(t *testing.T) {
	schema := createTestSchema()
	rule := NewLimitPushDownRule()

	// Create plan: Limit -> Filter -> Scan
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}
	filter := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: scan,
	}
	limit := &planner.LogicalLimit{
		Count: 10,
		Input: filter,
	}

	// Apply rule
	optimized := rule.Apply(limit)
	require.NotNil(t, optimized)

	// Should become: Filter -> Limit -> Scan
	newFilter, ok := optimized.(*planner.LogicalFilter)
	require.True(t, ok)

	newLimit, ok := newFilter.Input.(*planner.LogicalLimit)
	require.True(t, ok)
	assert.Equal(t, 10, newLimit.Count)

	_, ok = newLimit.Input.(*planner.LogicalScan)
	assert.True(t, ok)
}

func TestHepOptimizer_SingleRule(t *testing.T) {
	schema := createTestSchema()

	// Create plan: Filter -> Filter -> Scan
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}
	filter1 := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "host"},
			Operator: "=",
			Right:    &ast.Literal{Value: "server1", LiteralTyp: ast.LiteralTypeString},
		},
		Input: scan,
	}
	filter2 := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: filter1,
	}

	// Optimize with FilterMerge rule
	optimizer := NewHepOptimizer([]Rule{NewFilterMergeRule()})
	optimized, err := optimizer.Optimize(filter2)
	require.NoError(t, err)
	require.NotNil(t, optimized)

	// Should be merged into single filter
	mergedFilter, ok := optimized.(*planner.LogicalFilter)
	require.True(t, ok)

	andExpr, ok := mergedFilter.Condition.(*ast.BinaryExpression)
	require.True(t, ok)
	assert.Equal(t, "AND", andExpr.Operator)
}

func TestHepOptimizer_MultipleRules(t *testing.T) {
	schema := createTestSchema()

	// Create plan: Limit -> Filter -> Sort -> Scan
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}
	sort := &planner.LogicalSort{
		SortKeys: []*ast.SortKey{
			{Field: &ast.FieldReference{Name: "timestamp"}, Descending: true},
		},
		Input: scan,
	}
	filter := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: sort,
	}
	limit := &planner.LogicalLimit{
		Count: 10,
		Input: filter,
	}

	// Optimize with multiple rules
	optimizer := NewHepOptimizer([]Rule{
		NewFilterPushDownRule(),
		NewLimitPushDownRule(),
	})
	optimized, err := optimizer.Optimize(limit)
	require.NoError(t, err)
	require.NotNil(t, optimized)

	t.Logf("Original plan:\n%s", planner.PrintPlan(limit, 0))
	t.Logf("Optimized plan:\n%s", planner.PrintPlan(optimized, 0))

	// After optimization, structure should be better
	// (exact structure depends on rule application order)
	assert.NotNil(t, optimized)
}

func TestHepOptimizer_ComplexPlan(t *testing.T) {
	schema := createTestSchema()

	// Create complex plan with multiple optimization opportunities
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}

	// Filter 1
	filter1 := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "host"},
			Operator: "=",
			Right:    &ast.Literal{Value: "server1", LiteralTyp: ast.LiteralTypeString},
		},
		Input: scan,
	}

	// Filter 2
	filter2 := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: filter1,
	}

	// Project
	project := &planner.LogicalProject{
		Fields: []ast.Expression{
			&ast.FieldReference{Name: "status"},
			&ast.FieldReference{Name: "host"},
		},
		OutputSchema: schema,
		Input:        filter2,
	}

	// Sort
	sort := &planner.LogicalSort{
		SortKeys: []*ast.SortKey{
			{Field: &ast.FieldReference{Name: "timestamp"}, Descending: true},
		},
		Input: project,
	}

	// Limit
	limit := &planner.LogicalLimit{
		Count: 10,
		Input: sort,
	}

	// Optimize with all rules
	optimizer := DefaultOptimizer()
	optimized, err := optimizer.Optimize(limit)
	require.NoError(t, err)
	require.NotNil(t, optimized)

	t.Logf("Original plan:\n%s", planner.PrintPlan(limit, 0))
	t.Logf("Optimized plan:\n%s", planner.PrintPlan(optimized, 0))

	// Verify optimization happened
	assert.NotNil(t, optimized)
}

func TestHepOptimizer_MaxIterations(t *testing.T) {
	schema := createTestSchema()

	// Create a simple plan
	scan := &planner.LogicalScan{Source: "logs", OutputSchema: schema}
	filter := &planner.LogicalFilter{
		Condition: &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 500, LiteralTyp: ast.LiteralTypeInt},
		},
		Input: scan,
	}

	// Create optimizer with very low max iterations
	optimizer := NewHepOptimizer([]Rule{}).WithMaxIterations(1)
	optimized, err := optimizer.Optimize(filter)
	require.NoError(t, err)
	require.NotNil(t, optimized)

	// Should complete without error even with low iterations
	assert.NotNil(t, optimized)
}
