// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package physical

import (
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/planner"
)

// PhysicalPlanner converts logical plans to physical plans
type PhysicalPlanner struct {
	enablePushDown bool
	maxPushDownOps int // Maximum number of operations to push down
}

// NewPhysicalPlanner creates a new physical planner
func NewPhysicalPlanner() *PhysicalPlanner {
	return &PhysicalPlanner{
		enablePushDown: true,
		maxPushDownOps: 4, // Default: filter, project, sort, limit
	}
}

// WithPushDown enables or disables push-down optimization
func (pp *PhysicalPlanner) WithPushDown(enable bool) *PhysicalPlanner {
	pp.enablePushDown = enable
	return pp
}

// Plan converts a logical plan to a physical plan
func (pp *PhysicalPlanner) Plan(logicalPlan planner.LogicalPlan) (PhysicalPlan, error) {
	if logicalPlan == nil {
		return nil, fmt.Errorf("logical plan is nil")
	}

	// Check if we can push down operations
	if pp.enablePushDown {
		return pp.planWithPushDown(logicalPlan)
	}

	// No push-down: create coordinator-side operators for everything
	return pp.planCoordinatorOnly(logicalPlan)
}

// planWithPushDown creates a physical plan with push-down optimization
func (pp *PhysicalPlanner) planWithPushDown(logicalPlan planner.LogicalPlan) (PhysicalPlan, error) {
	// Try to push down operations to the scan
	scan, pushedDownOps := pp.extractPushableOps(logicalPlan)
	if scan == nil {
		return pp.planCoordinatorOnly(logicalPlan)
	}

	// Create physical scan with pushed-down operations
	physicalScan := &PhysicalScan{
		Source:       scan.Source,
		OutputSchema: scan.OutputSchema,
	}

	// Apply pushed-down operations to the scan
	if pushedDownOps.filter != nil {
		physicalScan.Filter = pushedDownOps.filter
	}
	if len(pushedDownOps.fields) > 0 {
		physicalScan.Fields = pushedDownOps.fields
	}
	if len(pushedDownOps.sortKeys) > 0 {
		physicalScan.SortKeys = pushedDownOps.sortKeys
	}
	if pushedDownOps.limit > 0 {
		physicalScan.Limit = pushedDownOps.limit
	}
	if len(pushedDownOps.computedFields) > 0 {
		physicalScan.ComputedFields = pushedDownOps.computedFields
	}

	// Build coordinator-side operations for non-pushable ops
	currentPlan := PhysicalPlan(physicalScan)

	for _, op := range pushedDownOps.coordinatorOps {
		var err error
		currentPlan, err = pp.planCoordinatorOp(op, currentPlan)
		if err != nil {
			return nil, err
		}
	}

	return currentPlan, nil
}

// pushableOps holds operations that can be pushed down
type pushableOps struct {
	filter         ast.Expression
	fields         []string
	sortKeys       []*ast.SortKey
	limit          int
	computedFields []*ast.EvalAssignment // Eval assignments that can be pushed down
	coordinatorOps []planner.LogicalPlan // Operations that must run on coordinator
}

// extractPushableOps extracts operations that can be pushed down to the scan
func (pp *PhysicalPlanner) extractPushableOps(plan planner.LogicalPlan) (*planner.LogicalScan, *pushableOps) {
	ops := &pushableOps{}
	current := plan
	barrierEncountered := false // Track if we've hit a non-pushable op

	// Walk down the plan tree collecting pushable operations
	for current != nil {
		switch p := current.(type) {
		case *planner.LogicalScan:
			// Found the scan - return it with collected ops
			return p, ops

		case *planner.LogicalFilter:
			// Filters can be pushed down if they're simple and no barrier encountered
			if !barrierEncountered && pp.canPushDownFilter(p.Condition) {
				if ops.filter == nil {
					ops.filter = p.Condition
				} else {
					// Merge with existing filter using AND
					ops.filter = &ast.BinaryExpression{
						Left:     ops.filter,
						Operator: "AND",
						Right:    p.Condition,
					}
				}
			} else {
				// Can't push down - must execute on coordinator
				ops.coordinatorOps = append(ops.coordinatorOps, p)
			}
			current = p.Input

		case *planner.LogicalProject:
			// Simple field projections can be pushed down if no barrier encountered
			if !barrierEncountered && pp.canPushDownProject(p) {
				ops.fields = pp.extractFieldNames(p.Fields)
			} else {
				// Complex projection - coordinator side
				ops.coordinatorOps = append(ops.coordinatorOps, p)
			}
			current = p.Input

		case *planner.LogicalSort:
			// Sorts can be pushed down if no barrier encountered
			if !barrierEncountered && pp.canPushDownSort(p.SortKeys) {
				ops.sortKeys = p.SortKeys
			} else {
				ops.coordinatorOps = append(ops.coordinatorOps, p)
			}
			current = p.Input

		case *planner.LogicalLimit:
			// Limits can be pushed down if no barrier encountered
			if !barrierEncountered {
				ops.limit = p.Count
			} else {
				ops.coordinatorOps = append(ops.coordinatorOps, p)
			}
			current = p.Input

		case *planner.LogicalAggregate:
			// Aggregations can sometimes be pushed down
			// For now, always run on coordinator
			ops.coordinatorOps = append(ops.coordinatorOps, p)
			barrierEncountered = true // Set barrier - ops above this must run on coordinator
			current = p.Input

		case *planner.LogicalExplain:
			// Explain wraps the plan - skip and continue
			current = p.Input

		// Tier 1 operators - all execute on coordinator
		case *planner.LogicalDedup:
			ops.coordinatorOps = append(ops.coordinatorOps, p)
			barrierEncountered = true // Dedup changes cardinality
			current = p.Input

		case *planner.LogicalBin:
			ops.coordinatorOps = append(ops.coordinatorOps, p)
			current = p.Input

		case *planner.LogicalTop:
			ops.coordinatorOps = append(ops.coordinatorOps, p)
			barrierEncountered = true // Top aggregates and limits
			current = p.Input

		case *planner.LogicalRare:
			ops.coordinatorOps = append(ops.coordinatorOps, p)
			barrierEncountered = true // Rare aggregates and limits
			current = p.Input

		case *planner.LogicalEval:
			// Eval can be pushed down if no barrier and expressions contain functions (WASM UDFs)
			if !barrierEncountered {
				pushable, nonPushable := pp.partitionEvalAssignments(p.Assignments)
				if len(pushable) > 0 {
					ops.computedFields = append(ops.computedFields, pushable...)
				}
				// If some assignments can't be pushed, keep them for coordinator
				if len(nonPushable) > 0 {
					// Create a new LogicalEval with only non-pushable assignments
					coordEval := &planner.LogicalEval{
						Assignments:  nonPushable,
						OutputSchema: p.OutputSchema,
						Input:        p.Input,
					}
					ops.coordinatorOps = append(ops.coordinatorOps, coordEval)
				}
			} else {
				// Barrier encountered - can't push down
				ops.coordinatorOps = append(ops.coordinatorOps, p)
			}
			current = p.Input

		case *planner.LogicalRename:
			ops.coordinatorOps = append(ops.coordinatorOps, p)
			current = p.Input

		default:
			// Unknown operator - can't push down
			ops.coordinatorOps = append(ops.coordinatorOps, p)
			barrierEncountered = true // Set barrier for unknown ops
			if len(p.Children()) > 0 {
				current = p.Children()[0]
			} else {
				current = nil
			}
		}
	}

	return nil, ops
}

// canPushDownFilter checks if a filter can be pushed down to data nodes
// This includes both native OpenSearch DSL and WASM UDF pushdown
func (pp *PhysicalPlanner) canPushDownFilter(condition ast.Expression) bool {
	// Check if it's a simple expression that can use native DSL
	if pp.isSimpleExpression(condition) {
		return true
	}

	// Check if it can be pushed down as WASM UDF
	// Note: FunctionBuilder is created in the DSL translator
	// For now, we conservatively allow function calls assuming they'll be checked later
	switch expr := condition.(type) {
	case *ast.FunctionCall:
		// Function calls can potentially be pushed as WASM UDFs
		return true

	case *ast.BinaryExpression:
		// Check if it's a comparison with a function call (e.g., abs(x) > 10)
		if _, ok := expr.Left.(*ast.FunctionCall); ok {
			return true
		}
		// Recursively check both sides
		return pp.canPushDownFilter(expr.Left) && pp.canPushDownFilter(expr.Right)

	case *ast.UnaryExpression:
		return pp.canPushDownFilter(expr.Operand)

	default:
		return false
	}
}

// isSimpleExpression checks if expression uses only simple field/literal comparisons
func (pp *PhysicalPlanner) isSimpleExpression(expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.FieldReference, *ast.Literal:
		return true

	case *ast.BinaryExpression:
		// Check if both sides are simple (no function calls)
		return pp.isSimpleExpression(e.Left) && pp.isSimpleExpression(e.Right)

	case *ast.UnaryExpression:
		return pp.isSimpleExpression(e.Operand)

	default:
		return false
	}
}

// canPushDownProject checks if a projection can be pushed down
func (pp *PhysicalPlanner) canPushDownProject(project *planner.LogicalProject) bool {
	// Can push down if all fields are simple field references
	for _, field := range project.Fields {
		if _, ok := field.(*ast.FieldReference); !ok {
			// Complex expression - can't push down
			return false
		}
	}
	return true
}

// canPushDownSort checks if a sort can be pushed down
func (pp *PhysicalPlanner) canPushDownSort(sortKeys []*ast.SortKey) bool {
	// Can push down if all sort keys are simple field references
	for _, key := range sortKeys {
		if _, ok := key.Field.(*ast.FieldReference); !ok {
			// Complex expression - can't push down
			return false
		}
	}
	return true
}

// extractFieldNames extracts field names from expressions
func (pp *PhysicalPlanner) extractFieldNames(fields []ast.Expression) []string {
	names := make([]string, 0, len(fields))
	for _, field := range fields {
		if fieldRef, ok := field.(*ast.FieldReference); ok {
			names = append(names, fieldRef.Name)
		}
	}
	return names
}

// partitionEvalAssignments separates eval assignments into pushable and non-pushable
// Assignments with function calls can be pushed down as WASM UDFs
func (pp *PhysicalPlanner) partitionEvalAssignments(assignments []*ast.EvalAssignment) (pushable, nonPushable []*ast.EvalAssignment) {
	pushable = make([]*ast.EvalAssignment, 0)
	nonPushable = make([]*ast.EvalAssignment, 0)

	for _, assignment := range assignments {
		if pp.canPushDownEvalExpression(assignment.Expression) {
			pushable = append(pushable, assignment)
		} else {
			nonPushable = append(nonPushable, assignment)
		}
	}

	return pushable, nonPushable
}

// canPushDownEvalExpression checks if an eval expression can be pushed down
// Expressions with function calls can be pushed down as WASM UDFs
func (pp *PhysicalPlanner) canPushDownEvalExpression(expr ast.Expression) bool {
	return pp.containsFunctionCall(expr)
}

// containsFunctionCall recursively checks if an expression contains a function call
func (pp *PhysicalPlanner) containsFunctionCall(expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.FunctionCall:
		return true

	case *ast.BinaryExpression:
		return pp.containsFunctionCall(e.Left) || pp.containsFunctionCall(e.Right)

	case *ast.UnaryExpression:
		return pp.containsFunctionCall(e.Operand)

	case *ast.FieldReference, *ast.Literal:
		return false

	default:
		return false
	}
}

// planCoordinatorOnly creates a physical plan with all operations on coordinator
func (pp *PhysicalPlanner) planCoordinatorOnly(logicalPlan planner.LogicalPlan) (PhysicalPlan, error) {
	switch p := logicalPlan.(type) {
	case *planner.LogicalScan:
		return &PhysicalScan{
			Source:       p.Source,
			OutputSchema: p.OutputSchema,
		}, nil

	case *planner.LogicalFilter:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalFilter{
			Condition: p.Condition,
			Input:     input,
		}, nil

	case *planner.LogicalProject:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalProject{
			Fields:       p.Fields,
			OutputSchema: p.OutputSchema,
			Input:        input,
			Exclude:      p.Exclude,
		}, nil

	case *planner.LogicalSort:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalSort{
			SortKeys: p.SortKeys,
			Input:    input,
		}, nil

	case *planner.LogicalLimit:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalLimit{
			Count: p.Count,
			Input: input,
		}, nil

	case *planner.LogicalAggregate:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}

		// Choose aggregation algorithm based on characteristics
		algorithm := pp.selectAggregationAlgorithm(p)

		return &PhysicalAggregate{
			GroupBy:      p.GroupBy,
			Aggregations: p.Aggregations,
			OutputSchema: p.OutputSchema,
			Input:        input,
			Algorithm:    algorithm,
		}, nil

	case *planner.LogicalExplain:
		// Explain is a wrapper - plan the inner query
		return pp.planCoordinatorOnly(p.Input)

	// Tier 1 operators
	case *planner.LogicalDedup:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalDedup{
			Fields:      p.Fields,
			Count:       p.Count,
			Consecutive: p.Consecutive,
			Input:       input,
		}, nil

	case *planner.LogicalBin:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalBin{
			Field:        p.Field,
			Span:         p.Span,
			Bins:         p.Bins,
			OutputSchema: p.Input.Schema(),
			Input:        input,
		}, nil

	case *planner.LogicalTop:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalTop{
			Fields:       p.Fields,
			Limit:        p.Limit,
			GroupBy:      p.GroupBy,
			ShowCount:    p.ShowCount,
			ShowPercent:  p.ShowPercent,
			OutputSchema: p.OutputSchema,
			Input:        input,
			Algorithm:    TopRareHash, // Default to hash algorithm
		}, nil

	case *planner.LogicalRare:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalRare{
			Fields:       p.Fields,
			Limit:        p.Limit,
			GroupBy:      p.GroupBy,
			ShowCount:    p.ShowCount,
			ShowPercent:  p.ShowPercent,
			OutputSchema: p.OutputSchema,
			Input:        input,
			Algorithm:    TopRareHash, // Default to hash algorithm
		}, nil

	case *planner.LogicalEval:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalEval{
			Assignments:  p.Assignments,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalRename:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalRename{
			Assignments:  p.Assignments,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalReplace:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalReplace{
			Mappings: p.Mappings,
			Field:    p.Field,
			Input:    input,
		}, nil

	case *planner.LogicalFillnull:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalFillnull{
			Assignments:  p.Assignments,
			DefaultValue: p.DefaultValue,
			Fields:       p.Fields,
			Input:        input,
		}, nil

	case *planner.LogicalParse:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalParse{
			SourceField:     p.SourceField,
			Pattern:         p.Pattern,
			ExtractedFields: p.ExtractedFields,
			OutputSchema:    p.OutputSchema,
			Input:           input,
		}, nil

	case *planner.LogicalRex:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalRex{
			SourceField:     p.SourceField,
			Pattern:         p.Pattern,
			ExtractedFields: p.ExtractedFields,
			OutputSchema:    p.OutputSchema,
			Input:           input,
		}, nil

	case *planner.LogicalLookup:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalLookup{
			TableName:      p.TableName,
			JoinField:      p.JoinField,
			JoinFieldAlias: p.JoinFieldAlias,
			OutputFields:   p.OutputFields,
			OutputAliases:  p.OutputAliases,
			OutputSchema:   p.OutputSchema,
			Input:          input,
		}, nil

	case *planner.LogicalAppend:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		subsearch, err := pp.planCoordinatorOnly(p.Subsearch)
		if err != nil {
			return nil, err
		}
		return &PhysicalAppend{
			Subsearch:    subsearch,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalJoin:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		right, err := pp.planCoordinatorOnly(p.Right)
		if err != nil {
			return nil, err
		}
		return &PhysicalJoin{
			JoinType:     p.JoinType,
			JoinField:    p.JoinField,
			RightField:   p.RightField,
			Right:        right,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalReverse:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalReverse{
			Input:        input,
			OutputSchema: p.OutputSchema,
		}, nil

	case *planner.LogicalFlatten:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalFlatten{
			Input:        input,
			Field:        p.Field,
			OutputSchema: p.OutputSchema,
		}, nil

	case *planner.LogicalTable:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalTable{
			Fields:       p.Fields,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalEventstats:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalEventstats{
			GroupBy:      p.GroupBy,
			Aggregations: p.Aggregations,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalStreamstats:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalStreamstats{
			GroupBy:      p.GroupBy,
			Aggregations: p.Aggregations,
			Window:       p.Window,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported logical operator: %T", logicalPlan)
	}
}

// planCoordinatorOp creates a coordinator-side physical operator
func (pp *PhysicalPlanner) planCoordinatorOp(logicalOp planner.LogicalPlan, input PhysicalPlan) (PhysicalPlan, error) {
	switch p := logicalOp.(type) {
	case *planner.LogicalFilter:
		return &PhysicalFilter{
			Condition: p.Condition,
			Input:     input,
		}, nil

	case *planner.LogicalProject:
		return &PhysicalProject{
			Fields:       p.Fields,
			OutputSchema: p.OutputSchema,
			Input:        input,
			Exclude:      p.Exclude,
		}, nil

	case *planner.LogicalSort:
		return &PhysicalSort{
			SortKeys: p.SortKeys,
			Input:    input,
		}, nil

	case *planner.LogicalAggregate:
		algorithm := pp.selectAggregationAlgorithm(p)
		return &PhysicalAggregate{
			GroupBy:      p.GroupBy,
			Aggregations: p.Aggregations,
			OutputSchema: p.OutputSchema,
			Input:        input,
			Algorithm:    algorithm,
		}, nil

	case *planner.LogicalLimit:
		return &PhysicalLimit{
			Count: p.Count,
			Input: input,
		}, nil

	// Tier 1 operators
	case *planner.LogicalDedup:
		return &PhysicalDedup{
			Fields:      p.Fields,
			Count:       p.Count,
			Consecutive: p.Consecutive,
			Input:       input,
		}, nil

	case *planner.LogicalBin:
		return &PhysicalBin{
			Field:        p.Field,
			Span:         p.Span,
			Bins:         p.Bins,
			OutputSchema: p.Input.Schema(),
			Input:        input,
		}, nil

	case *planner.LogicalTop:
		return &PhysicalTop{
			Fields:       p.Fields,
			Limit:        p.Limit,
			GroupBy:      p.GroupBy,
			ShowCount:    p.ShowCount,
			ShowPercent:  p.ShowPercent,
			OutputSchema: p.OutputSchema,
			Input:        input,
			Algorithm:    TopRareHash,
		}, nil

	case *planner.LogicalRare:
		return &PhysicalRare{
			Fields:       p.Fields,
			Limit:        p.Limit,
			GroupBy:      p.GroupBy,
			ShowCount:    p.ShowCount,
			ShowPercent:  p.ShowPercent,
			OutputSchema: p.OutputSchema,
			Input:        input,
			Algorithm:    TopRareHash,
		}, nil

	case *planner.LogicalEval:
		return &PhysicalEval{
			Assignments:  p.Assignments,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalRename:
		return &PhysicalRename{
			Assignments:  p.Assignments,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalReplace:
		return &PhysicalReplace{
			Mappings: p.Mappings,
			Field:    p.Field,
			Input:    input,
		}, nil

	case *planner.LogicalFillnull:
		return &PhysicalFillnull{
			Assignments:  p.Assignments,
			DefaultValue: p.DefaultValue,
			Fields:       p.Fields,
			Input:        input,
		}, nil

	case *planner.LogicalParse:
		return &PhysicalParse{
			SourceField:     p.SourceField,
			Pattern:         p.Pattern,
			ExtractedFields: p.ExtractedFields,
			OutputSchema:    p.OutputSchema,
			Input:           input,
		}, nil

	case *planner.LogicalRex:
		return &PhysicalRex{
			SourceField:     p.SourceField,
			Pattern:         p.Pattern,
			ExtractedFields: p.ExtractedFields,
			OutputSchema:    p.OutputSchema,
			Input:           input,
		}, nil

	case *planner.LogicalLookup:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalLookup{
			TableName:      p.TableName,
			JoinField:      p.JoinField,
			JoinFieldAlias: p.JoinFieldAlias,
			OutputFields:   p.OutputFields,
			OutputAliases:  p.OutputAliases,
			OutputSchema:   p.OutputSchema,
			Input:          input,
		}, nil

	case *planner.LogicalAppend:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		subsearch, err := pp.planCoordinatorOnly(p.Subsearch)
		if err != nil {
			return nil, err
		}
		return &PhysicalAppend{
			Subsearch:    subsearch,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalJoin:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		right, err := pp.planCoordinatorOnly(p.Right)
		if err != nil {
			return nil, err
		}
		return &PhysicalJoin{
			JoinType:     p.JoinType,
			JoinField:    p.JoinField,
			RightField:   p.RightField,
			Right:        right,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalReverse:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalReverse{
			Input:        input,
			OutputSchema: p.OutputSchema,
		}, nil

	case *planner.LogicalFlatten:
		input, err := pp.planCoordinatorOnly(p.Input)
		if err != nil {
			return nil, err
		}
		return &PhysicalFlatten{
			Input:        input,
			Field:        p.Field,
			OutputSchema: p.OutputSchema,
		}, nil

	case *planner.LogicalTable:
		return &PhysicalTable{
			Fields:       p.Fields,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalEventstats:
		return &PhysicalEventstats{
			GroupBy:      p.GroupBy,
			Aggregations: p.Aggregations,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	case *planner.LogicalStreamstats:
		return &PhysicalStreamstats{
			GroupBy:      p.GroupBy,
			Aggregations: p.Aggregations,
			Window:       p.Window,
			OutputSchema: p.OutputSchema,
			Input:        input,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported coordinator operation: %T", logicalOp)
	}
}

// selectAggregationAlgorithm chooses between hash and stream aggregation
func (pp *PhysicalPlanner) selectAggregationAlgorithm(agg *planner.LogicalAggregate) AggregationAlgorithm {
	// If input is sorted by GROUP BY keys, use stream aggregation
	// For now, always use hash aggregation (more general)
	// TODO: Detect sorted input and use stream aggregation when possible
	return HashAggregation
}
