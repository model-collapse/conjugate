// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package optimizer

import (
	"github.com/conjugate/conjugate/pkg/ppl/planner"
)

// Rule represents a query optimization rule
type Rule interface {
	// Name returns the rule name
	Name() string

	// Apply applies the rule to a plan and returns the optimized plan
	// Returns nil if the rule doesn't apply
	Apply(plan planner.LogicalPlan) planner.LogicalPlan

	// Description returns a human-readable description of the rule
	Description() string
}

// Optimizer optimizes logical query plans
type Optimizer interface {
	// Optimize applies optimization rules to the plan
	Optimize(plan planner.LogicalPlan) (planner.LogicalPlan, error)
}

// HepOptimizer implements the HEP (Heuristic Execution Planner) pattern
// It applies rules iteratively until no changes occur or max iterations reached
type HepOptimizer struct {
	rules         []Rule
	maxIterations int
}

// NewHepOptimizer creates a new HEP optimizer with the given rules
func NewHepOptimizer(rules []Rule) *HepOptimizer {
	return &HepOptimizer{
		rules:         rules,
		maxIterations: 10, // Default max iterations
	}
}

// WithMaxIterations sets the maximum number of optimization iterations
func (h *HepOptimizer) WithMaxIterations(max int) *HepOptimizer {
	h.maxIterations = max
	return h
}

// Optimize applies all rules iteratively until the plan stabilizes
func (h *HepOptimizer) Optimize(plan planner.LogicalPlan) (planner.LogicalPlan, error) {
	if plan == nil {
		return nil, nil
	}

	currentPlan := plan
	iteration := 0

	for iteration < h.maxIterations {
		changed := false
		iteration++

		// Apply each rule once per iteration
		for _, rule := range h.rules {
			optimized := h.applyRuleRecursively(rule, currentPlan)
			if optimized != nil && optimized != currentPlan {
				currentPlan = optimized
				changed = true
			}
		}

		// If no changes were made, optimization is complete
		if !changed {
			break
		}
	}

	return currentPlan, nil
}

// applyRuleRecursively applies a rule to the plan and all its children
func (h *HepOptimizer) applyRuleRecursively(rule Rule, plan planner.LogicalPlan) planner.LogicalPlan {
	if plan == nil {
		return nil
	}

	// Try to apply rule to current node
	optimized := rule.Apply(plan)
	if optimized != nil && optimized != plan {
		// Rule was applied, return optimized plan
		return optimized
	}

	// Rule didn't apply, try children
	children := plan.Children()
	if len(children) == 0 {
		return plan // No children, return as-is
	}

	// Recursively optimize children
	hasChanges := false
	optimizedChildren := make([]planner.LogicalPlan, len(children))
	for i, child := range children {
		optimizedChild := h.applyRuleRecursively(rule, child)
		optimizedChildren[i] = optimizedChild
		if optimizedChild != child {
			hasChanges = true
		}
	}

	// If any child changed, rebuild the parent with new children
	if hasChanges {
		return h.rebuildWithChildren(plan, optimizedChildren)
	}

	return plan
}

// rebuildWithChildren rebuilds a plan node with new children
func (h *HepOptimizer) rebuildWithChildren(plan planner.LogicalPlan, newChildren []planner.LogicalPlan) planner.LogicalPlan {
	if len(newChildren) == 0 {
		return plan
	}

	switch p := plan.(type) {
	case *planner.LogicalFilter:
		return &planner.LogicalFilter{
			Condition: p.Condition,
			Input:     newChildren[0],
		}

	case *planner.LogicalProject:
		return &planner.LogicalProject{
			Fields:       p.Fields,
			OutputSchema: p.OutputSchema,
			Input:        newChildren[0],
			Exclude:      p.Exclude,
		}

	case *planner.LogicalSort:
		return &planner.LogicalSort{
			SortKeys: p.SortKeys,
			Input:    newChildren[0],
		}

	case *planner.LogicalLimit:
		return &planner.LogicalLimit{
			Count: p.Count,
			Input: newChildren[0],
		}

	case *planner.LogicalAggregate:
		return &planner.LogicalAggregate{
			GroupBy:      p.GroupBy,
			Aggregations: p.Aggregations,
			OutputSchema: p.OutputSchema,
			Input:        newChildren[0],
		}

	case *planner.LogicalExplain:
		return &planner.LogicalExplain{
			Input: newChildren[0],
		}

	default:
		// Unknown type, return as-is
		return plan
	}
}

// DefaultOptimizer returns an optimizer with standard optimization rules
func DefaultOptimizer() *HepOptimizer {
	rules := []Rule{
		NewFilterMergeRule(),
		NewFilterPushDownRule(),
		NewProjectMergeRule(),
		NewProjectionPruningRule(),
		NewConstantFoldingRule(),
	}
	return NewHepOptimizer(rules)
}
