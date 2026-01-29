// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package optimizer

import (
	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/planner"
)

// FilterMergeRule combines consecutive filters into a single filter with AND
type FilterMergeRule struct{}

func NewFilterMergeRule() *FilterMergeRule {
	return &FilterMergeRule{}
}

func (r *FilterMergeRule) Name() string {
	return "FilterMerge"
}

func (r *FilterMergeRule) Description() string {
	return "Merges consecutive filters into a single filter with AND"
}

func (r *FilterMergeRule) Apply(plan planner.LogicalPlan) planner.LogicalPlan {
	filter, ok := plan.(*planner.LogicalFilter)
	if !ok {
		return nil
	}

	// Check if the input is also a filter
	childFilter, ok := filter.Input.(*planner.LogicalFilter)
	if !ok {
		return nil
	}

	// Merge: Filter(A) -> Filter(B) => Filter(A AND B) -> B.input
	mergedCondition := &ast.BinaryExpression{
		Left:     filter.Condition,
		Operator: "AND",
		Right:    childFilter.Condition,
	}

	return &planner.LogicalFilter{
		Condition: mergedCondition,
		Input:     childFilter.Input,
	}
}

// FilterPushDownRule pushes filters down past projections and sorts when possible
type FilterPushDownRule struct{}

func NewFilterPushDownRule() *FilterPushDownRule {
	return &FilterPushDownRule{}
}

func (r *FilterPushDownRule) Name() string {
	return "FilterPushDown"
}

func (r *FilterPushDownRule) Description() string {
	return "Pushes filters closer to the data source for early filtering"
}

func (r *FilterPushDownRule) Apply(plan planner.LogicalPlan) planner.LogicalPlan {
	filter, ok := plan.(*planner.LogicalFilter)
	if !ok {
		return nil
	}

	// Can push down past Project if filter only references fields in projection
	if project, ok := filter.Input.(*planner.LogicalProject); ok {
		// For simplicity, always push down (TODO: check field dependencies)
		// Filter -> Project => Project -> Filter
		return &planner.LogicalProject{
			Fields:       project.Fields,
			OutputSchema: project.OutputSchema,
			Exclude:      project.Exclude,
			Input: &planner.LogicalFilter{
				Condition: filter.Condition,
				Input:     project.Input,
			},
		}
	}

	// Can push down past Sort (sorting doesn't affect filtering)
	if sort, ok := filter.Input.(*planner.LogicalSort); ok {
		// Filter -> Sort => Sort -> Filter
		return &planner.LogicalSort{
			SortKeys: sort.SortKeys,
			Input: &planner.LogicalFilter{
				Condition: filter.Condition,
				Input:     sort.Input,
			},
		}
	}

	return nil
}

// ProjectMergeRule combines consecutive projections
type ProjectMergeRule struct{}

func NewProjectMergeRule() *ProjectMergeRule {
	return &ProjectMergeRule{}
}

func (r *ProjectMergeRule) Name() string {
	return "ProjectMerge"
}

func (r *ProjectMergeRule) Description() string {
	return "Merges consecutive projections into a single projection"
}

func (r *ProjectMergeRule) Apply(plan planner.LogicalPlan) planner.LogicalPlan {
	project, ok := plan.(*planner.LogicalProject)
	if !ok {
		return nil
	}

	// Check if the input is also a projection
	childProject, ok := project.Input.(*planner.LogicalProject)
	if !ok {
		return nil
	}

	// For now, just keep the outer projection (simpler)
	// TODO: Properly merge field lists based on include/exclude semantics
	return &planner.LogicalProject{
		Fields:       project.Fields,
		OutputSchema: project.OutputSchema,
		Input:        childProject.Input,
		Exclude:      project.Exclude,
	}
}

// ProjectionPruningRule removes unnecessary projections
type ProjectionPruningRule struct{}

func NewProjectionPruningRule() *ProjectionPruningRule {
	return &ProjectionPruningRule{}
}

func (r *ProjectionPruningRule) Name() string {
	return "ProjectionPruning"
}

func (r *ProjectionPruningRule) Description() string {
	return "Removes projections that don't change the schema"
}

func (r *ProjectionPruningRule) Apply(plan planner.LogicalPlan) planner.LogicalPlan {
	project, ok := plan.(*planner.LogicalProject)
	if !ok {
		return nil
	}

	// If projection includes all fields from input, it's redundant
	// For now, we'll keep projections (TODO: implement proper field set comparison)
	_ = project

	return nil // Don't optimize for now
}

// ConstantFoldingRule evaluates constant expressions at compile time
type ConstantFoldingRule struct{}

func NewConstantFoldingRule() *ConstantFoldingRule {
	return &ConstantFoldingRule{}
}

func (r *ConstantFoldingRule) Name() string {
	return "ConstantFolding"
}

func (r *ConstantFoldingRule) Description() string {
	return "Evaluates constant expressions at compile time"
}

func (r *ConstantFoldingRule) Apply(plan planner.LogicalPlan) planner.LogicalPlan {
	// Look for filters with constant conditions
	filter, ok := plan.(*planner.LogicalFilter)
	if !ok {
		return nil
	}

	// Fold constant expressions in the condition
	folded := r.foldExpression(filter.Condition)
	if folded != filter.Condition {
		return &planner.LogicalFilter{
			Condition: folded,
			Input:     filter.Input,
		}
	}

	return nil
}

// foldExpression attempts to evaluate constant expressions
func (r *ConstantFoldingRule) foldExpression(expr ast.Expression) ast.Expression {
	switch e := expr.(type) {
	case *ast.BinaryExpression:
		// Fold operands first
		leftFolded := r.foldExpression(e.Left)
		rightFolded := r.foldExpression(e.Right)

		// If both are literals, evaluate
		leftLit, leftIsLit := leftFolded.(*ast.Literal)
		rightLit, rightIsLit := rightFolded.(*ast.Literal)

		if leftIsLit && rightIsLit {
			// Try to evaluate constant expression
			result := r.evaluateConstant(leftLit, rightLit, e.Operator)
			if result != nil {
				return result
			}
		}

		// Return with folded operands
		if leftFolded != e.Left || rightFolded != e.Right {
			return &ast.BinaryExpression{
				Left:     leftFolded,
				Operator: e.Operator,
				Right:    rightFolded,
			}
		}

	case *ast.UnaryExpression:
		// Fold operand
		operandFolded := r.foldExpression(e.Operand)
		if operandLit, isLit := operandFolded.(*ast.Literal); isLit {
			// Evaluate unary constant
			if e.Operator == "NOT" && operandLit.LiteralTyp == ast.LiteralTypeBool {
				return &ast.Literal{
					Value:       !operandLit.Value.(bool),
					LiteralTyp: ast.LiteralTypeBool,
				}
			}
		}

		if operandFolded != e.Operand {
			return &ast.UnaryExpression{
				Operator: e.Operator,
				Operand:  operandFolded,
			}
		}
	}

	return expr
}

// evaluateConstant evaluates a binary operation on two literals
func (r *ConstantFoldingRule) evaluateConstant(left, right *ast.Literal, op string) *ast.Literal {
	// Only handle simple integer arithmetic for now
	if left.LiteralTyp != ast.LiteralTypeInt || right.LiteralTyp != ast.LiteralTypeInt {
		return nil
	}

	leftVal, ok1 := left.Value.(int)
	rightVal, ok2 := right.Value.(int)
	if !ok1 || !ok2 {
		return nil
	}

	var result int
	switch op {
	case "+":
		result = leftVal + rightVal
	case "-":
		result = leftVal - rightVal
	case "*":
		result = leftVal * rightVal
	case "/":
		if rightVal == 0 {
			return nil // Avoid division by zero
		}
		result = leftVal / rightVal
	case "%":
		if rightVal == 0 {
			return nil
		}
		result = leftVal % rightVal
	default:
		return nil // Unsupported operator
	}

	return &ast.Literal{
		Value:       result,
		LiteralTyp: ast.LiteralTypeInt,
	}
}

// LimitPushDownRule pushes limits down to reduce data early
type LimitPushDownRule struct{}

func NewLimitPushDownRule() *LimitPushDownRule {
	return &LimitPushDownRule{}
}

func (r *LimitPushDownRule) Name() string {
	return "LimitPushDown"
}

func (r *LimitPushDownRule) Description() string {
	return "Pushes limits down past non-expanding operators"
}

func (r *LimitPushDownRule) Apply(plan planner.LogicalPlan) planner.LogicalPlan {
	limit, ok := plan.(*planner.LogicalLimit)
	if !ok {
		return nil
	}

	// Can push down past Filter (filtering can only reduce rows)
	if filter, ok := limit.Input.(*planner.LogicalFilter); ok {
		// Limit -> Filter => Filter -> Limit
		return &planner.LogicalFilter{
			Condition: filter.Condition,
			Input: &planner.LogicalLimit{
				Count: limit.Count,
				Input: filter.Input,
			},
		}
	}

	// Can push down past Project (1:1 mapping)
	if project, ok := limit.Input.(*planner.LogicalProject); ok {
		// Limit -> Project => Project -> Limit
		return &planner.LogicalProject{
			Fields:       project.Fields,
			OutputSchema: project.OutputSchema,
			Exclude:      project.Exclude,
			Input: &planner.LogicalLimit{
				Count: limit.Count,
				Input: project.Input,
			},
		}
	}

	return nil
}

// EliminateRedundantSortRule removes sorts when output is already sorted
type EliminateRedundantSortRule struct{}

func NewEliminateRedundantSortRule() *EliminateRedundantSortRule {
	return &EliminateRedundantSortRule{}
}

func (r *EliminateRedundantSortRule) Name() string {
	return "EliminateRedundantSort"
}

func (r *EliminateRedundantSortRule) Description() string {
	return "Removes redundant sort operations"
}

func (r *EliminateRedundantSortRule) Apply(plan planner.LogicalPlan) planner.LogicalPlan {
	sort, ok := plan.(*planner.LogicalSort)
	if !ok {
		return nil
	}

	// Check if input is also a sort - outer sort supersedes inner sort
	if _, ok := sort.Input.(*planner.LogicalSort); ok {
		// Sort -> Sort => Sort (keep outer, discard inner)
		// For now, don't optimize (would need to check if inner sort is compatible)
		_ = sort
	}

	return nil
}
