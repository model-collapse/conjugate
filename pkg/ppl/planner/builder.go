// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package planner

import (
	"fmt"
	"regexp"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
)

// PlanBuilder builds a logical plan from an analyzed AST
type PlanBuilder struct {
	schema *analyzer.Schema
}

// NewPlanBuilder creates a new plan builder
func NewPlanBuilder(schema *analyzer.Schema) *PlanBuilder {
	return &PlanBuilder{schema: schema}
}

// Build builds a logical plan from a query AST
func (b *PlanBuilder) Build(query *ast.Query) (LogicalPlan, error) {
	if len(query.Commands) == 0 {
		return nil, fmt.Errorf("query has no commands")
	}

	// Start with the first command (should be SearchCommand or metadata command)
	var plan LogicalPlan
	var err error

	firstCmd := query.Commands[0]
	switch cmd := firstCmd.(type) {
	case *ast.SearchCommand:
		plan, err = b.buildSearchCommand(cmd)
		if err != nil {
			return nil, err
		}

	case *ast.DescribeCommand:
		// Describe doesn't produce a data plan - handle separately
		return nil, fmt.Errorf("DESCRIBE command not yet supported in planner")

	case *ast.ShowDatasourcesCommand:
		// ShowDatasources doesn't produce a data plan - handle separately
		return nil, fmt.Errorf("SHOWDATASOURCES command not yet supported in planner")

	case *ast.ExplainCommand:
		// EXPLAIN wraps the actual query - skip it for now
		// The explain will be handled at the top level after building the plan
		if len(query.Commands) < 2 {
			return nil, fmt.Errorf("EXPLAIN requires a query to explain")
		}
		// Process the rest of the commands
		// (Will handle EXPLAIN wrapper later)

	default:
		return nil, fmt.Errorf("query must start with search command, got %T", firstCmd)
	}

	// Process remaining commands in pipeline order
	startIdx := 1
	if _, isExplain := firstCmd.(*ast.ExplainCommand); isExplain {
		startIdx = 1 // Start from second command after EXPLAIN
	}

	for i := startIdx; i < len(query.Commands); i++ {
		cmd := query.Commands[i]
		plan, err = b.buildCommand(cmd, plan)
		if err != nil {
			return nil, fmt.Errorf("failed to build command %d: %w", i, err)
		}
	}

	// If the first command was EXPLAIN, wrap the plan
	if _, isExplain := firstCmd.(*ast.ExplainCommand); isExplain {
		plan = &LogicalExplain{Input: plan}
	}

	return plan, nil
}

// buildCommand builds a logical operator for a command, chaining onto the input plan
func (b *PlanBuilder) buildCommand(cmd ast.Command, input LogicalPlan) (LogicalPlan, error) {
	switch c := cmd.(type) {
	case *ast.SearchCommand:
		// SearchCommand should only be the first command
		return nil, fmt.Errorf("search command can only be the first command in a pipeline")

	case *ast.WhereCommand:
		return b.buildWhereCommand(c, input)

	case *ast.FieldsCommand:
		return b.buildFieldsCommand(c, input)

	case *ast.StatsCommand:
		return b.buildStatsCommand(c, input)

	case *ast.SortCommand:
		return b.buildSortCommand(c, input)

	case *ast.HeadCommand:
		return b.buildHeadCommand(c, input)

	case *ast.DedupCommand:
		return b.buildDedupCommand(c, input)

	case *ast.BinCommand:
		return b.buildBinCommand(c, input)

	case *ast.TopCommand:
		return b.buildTopCommand(c, input)

	case *ast.RareCommand:
		return b.buildRareCommand(c, input)

	case *ast.ChartCommand:
		return b.buildChartCommand(c, input)

	case *ast.TimechartCommand:
		return b.buildTimechartCommand(c, input)

	case *ast.EvalCommand:
		return b.buildEvalCommand(c, input)

	case *ast.RenameCommand:
		return b.buildRenameCommand(c, input)

	case *ast.ReplaceCommand:
		return b.buildReplaceCommand(c, input)

	case *ast.FillnullCommand:
		return b.buildFillnullCommand(c, input)

	case *ast.ParseCommand:
		return b.buildParseCommand(c, input)

	case *ast.RexCommand:
		return b.buildRexCommand(c, input)

	case *ast.LookupCommand:
		return b.buildLookupCommand(c, input)

	case *ast.AppendCommand:
		return b.buildAppendCommand(c, input)

	case *ast.JoinCommand:
		return b.buildJoinCommand(c, input)

	case *ast.TableCommand:
		return b.buildTableCommand(c, input)

	case *ast.EventstatsCommand:
		return b.buildEventstatsCommand(c, input)

	case *ast.StreamstatsCommand:
		return b.buildStreamstatsCommand(c, input)

	case *ast.ReverseCommand:
		return b.buildReverseCommand(c, input)

	case *ast.FlattenCommand:
		return b.buildFlattenCommand(c, input)

	default:
		return nil, fmt.Errorf("unsupported command type: %T", cmd)
	}
}

// buildSearchCommand builds a LogicalScan operator
func (b *PlanBuilder) buildSearchCommand(cmd *ast.SearchCommand) (LogicalPlan, error) {
	// Create a scan operator with the schema
	return &LogicalScan{
		Source:       cmd.Source,
		OutputSchema: b.schema,
	}, nil
}

// buildWhereCommand builds a LogicalFilter operator
func (b *PlanBuilder) buildWhereCommand(cmd *ast.WhereCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("where command requires an input plan")
	}

	return &LogicalFilter{
		Condition: cmd.Condition,
		Input:     input,
	}, nil
}

// buildFieldsCommand builds a LogicalProject operator
func (b *PlanBuilder) buildFieldsCommand(cmd *ast.FieldsCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("fields command requires an input plan")
	}

	// Build output schema based on projected fields
	inputSchema := input.Schema()
	var outputSchema *analyzer.Schema

	if !cmd.Includes {
		// Exclude mode: remove specified fields
		outputSchema = inputSchema.Clone()
		for _, field := range cmd.Fields {
			if fieldRef, ok := field.(*ast.FieldReference); ok {
				delete(outputSchema.Fields, fieldRef.Name)
			}
		}
	} else {
		// Include mode: only keep specified fields
		fieldNames := make([]string, 0, len(cmd.Fields))
		for _, field := range cmd.Fields {
			if fieldRef, ok := field.(*ast.FieldReference); ok {
				fieldNames = append(fieldNames, fieldRef.Name)
			} else {
				// Complex expression - for now, just use the expression string as field name
				// TODO: Handle complex projections better
				fieldNames = append(fieldNames, field.String())
			}
		}

		var err error
		outputSchema, err = inputSchema.Project(fieldNames)
		if err != nil {
			// If projection fails, use input schema as fallback
			outputSchema = inputSchema
		}
	}

	return &LogicalProject{
		Fields:       cmd.Fields,
		OutputSchema: outputSchema,
		Input:        input,
		Exclude:      !cmd.Includes, // Exclude mode when Includes is false
	}, nil
}

// buildStatsCommand builds a LogicalAggregate operator
func (b *PlanBuilder) buildStatsCommand(cmd *ast.StatsCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("stats command requires an input plan")
	}

	// Build output schema: GROUP BY fields + aggregation results
	outputSchema := analyzer.NewSchema(input.Schema().Source + "_stats")

	// Add GROUP BY fields to output schema
	for _, groupBy := range cmd.GroupBy {
		if fieldRef, ok := groupBy.(*ast.FieldReference); ok {
			// Get field type from input schema
			if field, err := input.Schema().GetField(fieldRef.Name); err == nil {
				outputSchema.AddField(fieldRef.Name, field.Type)
			} else {
				// Unknown field type - use string as default
				outputSchema.AddField(fieldRef.Name, analyzer.FieldTypeString)
			}
		}
	}

	// Add aggregation results to output schema
	for _, agg := range cmd.Aggregations {
		aggName := agg.Alias
		if aggName == "" {
			// Generate default name from function
			if agg.Func != nil {
				aggName = agg.Func.Name
			} else {
				aggName = "agg"
			}
		}

		// Infer aggregation result type
		// TODO: Use type checker for proper type inference
		aggType := analyzer.FieldTypeLong // Default for count()
		if agg.Func != nil {
			switch agg.Func.Name {
			case "count":
				aggType = analyzer.FieldTypeLong
			case "sum", "avg", "stddev", "variance":
				aggType = analyzer.FieldTypeDouble
			case "min", "max":
				// Min/max preserve input type - for now use double
				aggType = analyzer.FieldTypeDouble
			}
		}

		outputSchema.AddField(aggName, aggType)
	}

	return &LogicalAggregate{
		GroupBy:      cmd.GroupBy,
		Aggregations: cmd.Aggregations,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildSortCommand builds a LogicalSort operator
func (b *PlanBuilder) buildSortCommand(cmd *ast.SortCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("sort command requires an input plan")
	}

	return &LogicalSort{
		SortKeys: cmd.SortKeys,
		Input:    input,
	}, nil
}

// buildHeadCommand builds a LogicalLimit operator
func (b *PlanBuilder) buildHeadCommand(cmd *ast.HeadCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("head command requires an input plan")
	}

	return &LogicalLimit{
		Count: cmd.Count,
		Input: input,
	}, nil
}

// buildDedupCommand builds a LogicalDedup operator
func (b *PlanBuilder) buildDedupCommand(cmd *ast.DedupCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("dedup command requires an input plan")
	}

	count := cmd.Count
	if count == 0 {
		count = 1 // Default to keeping 1 duplicate
	}

	return &LogicalDedup{
		Fields:      cmd.Fields,
		Count:       count,
		Consecutive: cmd.Consecutive,
		Input:       input,
	}, nil
}

// buildBinCommand builds a LogicalBin operator
func (b *PlanBuilder) buildBinCommand(cmd *ast.BinCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("bin command requires an input plan")
	}

	return &LogicalBin{
		Field: cmd.Field,
		Span:  cmd.Span,
		Bins:  cmd.Bins,
		Input: input,
	}, nil
}

// buildTopCommand builds a LogicalTop operator
func (b *PlanBuilder) buildTopCommand(cmd *ast.TopCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("top command requires an input plan")
	}

	// Build output schema for top command
	// Output includes: field values, count, and optionally percent
	outputSchema := analyzer.NewSchema(input.Schema().Source)

	// Add fields from the top command
	for _, field := range cmd.Fields {
		if ref, ok := field.(*ast.FieldReference); ok {
			if inField, err := input.Schema().GetField(ref.Name); err == nil {
				outputSchema.AddField(ref.Name, inField.Type)
			} else {
				outputSchema.AddField(ref.Name, analyzer.FieldTypeString)
			}
		}
	}

	// Add count field
	countField := cmd.CountField
	if countField == "" {
		countField = "count"
	}
	outputSchema.AddField(countField, analyzer.FieldTypeLong)

	// Add percent field if requested
	if cmd.ShowPercent {
		percentField := cmd.PercentField
		if percentField == "" {
			percentField = "percent"
		}
		outputSchema.AddField(percentField, analyzer.FieldTypeDouble)
	}

	limit := cmd.Limit
	if limit == 0 {
		limit = 10 // Default to top 10
	}

	return &LogicalTop{
		Fields:       cmd.Fields,
		Limit:        limit,
		GroupBy:      cmd.GroupBy,
		ShowCount:    cmd.ShowCount,
		ShowPercent:  cmd.ShowPercent,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildRareCommand builds a LogicalRare operator
func (b *PlanBuilder) buildRareCommand(cmd *ast.RareCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("rare command requires an input plan")
	}

	// Build output schema for rare command (same structure as top)
	outputSchema := analyzer.NewSchema(input.Schema().Source)

	// Add fields from the rare command
	for _, field := range cmd.Fields {
		if ref, ok := field.(*ast.FieldReference); ok {
			if inField, err := input.Schema().GetField(ref.Name); err == nil {
				outputSchema.AddField(ref.Name, inField.Type)
			} else {
				outputSchema.AddField(ref.Name, analyzer.FieldTypeString)
			}
		}
	}

	// Add count field
	countField := cmd.CountField
	if countField == "" {
		countField = "count"
	}
	outputSchema.AddField(countField, analyzer.FieldTypeLong)

	// Add percent field if requested
	if cmd.ShowPercent {
		percentField := cmd.PercentField
		if percentField == "" {
			percentField = "percent"
		}
		outputSchema.AddField(percentField, analyzer.FieldTypeDouble)
	}

	limit := cmd.Limit
	if limit == 0 {
		limit = 10 // Default to rare 10
	}

	return &LogicalRare{
		Fields:       cmd.Fields,
		Limit:        limit,
		GroupBy:      cmd.GroupBy,
		ShowCount:    cmd.ShowCount,
		ShowPercent:  cmd.ShowPercent,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildChartCommand builds a LogicalAggregate operator with grouping
func (b *PlanBuilder) buildChartCommand(cmd *ast.ChartCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("chart command requires an input plan")
	}

	// Chart is essentially a stats command with grouping
	// Build output schema
	outputSchema := analyzer.NewSchema(input.Schema().Source)

	// Add GROUP BY fields
	for _, groupExpr := range cmd.GroupBy {
		if ref, ok := groupExpr.(*ast.FieldReference); ok {
			if inField, err := input.Schema().GetField(ref.Name); err == nil {
				outputSchema.AddField(ref.Name, inField.Type)
			} else {
				outputSchema.AddField(ref.Name, analyzer.FieldTypeString)
			}
		}
	}

	// Add aggregation result fields
	for _, agg := range cmd.Aggregations {
		aggName := agg.Alias
		if aggName == "" {
			aggName = agg.Func.Name
			if len(agg.Func.Arguments) > 0 {
				if ref, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
					aggName += "_" + ref.Name
				}
			}
		}

		// Infer aggregation result type
		aggType := analyzer.FieldTypeLong
		if agg.Func != nil {
			switch agg.Func.Name {
			case "count":
				aggType = analyzer.FieldTypeLong
			case "sum", "avg", "stddev", "variance":
				aggType = analyzer.FieldTypeDouble
			case "min", "max":
				aggType = analyzer.FieldTypeDouble
			}
		}

		outputSchema.AddField(aggName, aggType)
	}

	return &LogicalAggregate{
		GroupBy:      cmd.GroupBy,
		Aggregations: cmd.Aggregations,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildTimechartCommand builds a LogicalAggregate operator with time-based grouping
func (b *PlanBuilder) buildTimechartCommand(cmd *ast.TimechartCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("timechart command requires an input plan")
	}

	// Timechart is an aggregation grouped by time buckets
	// We'll create a synthetic time bucket field for grouping
	outputSchema := analyzer.NewSchema(input.Schema().Source)

	// Add _time bucket field (the time binning result)
	outputSchema.AddField("_time", analyzer.FieldTypeDate)

	// Add any additional GROUP BY fields
	for _, groupExpr := range cmd.GroupBy {
		if ref, ok := groupExpr.(*ast.FieldReference); ok {
			if inField, err := input.Schema().GetField(ref.Name); err == nil {
				outputSchema.AddField(ref.Name, inField.Type)
			} else {
				outputSchema.AddField(ref.Name, analyzer.FieldTypeString)
			}
		}
	}

	// Add aggregation result fields
	for _, agg := range cmd.Aggregations {
		aggName := agg.Alias
		if aggName == "" {
			aggName = agg.Func.Name
			if len(agg.Func.Arguments) > 0 {
				if ref, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
					aggName += "_" + ref.Name
				}
			}
		}

		// Infer aggregation result type
		aggType := analyzer.FieldTypeLong
		if agg.Func != nil {
			switch agg.Func.Name {
			case "count":
				aggType = analyzer.FieldTypeLong
			case "sum", "avg", "stddev", "variance":
				aggType = analyzer.FieldTypeDouble
			case "min", "max":
				aggType = analyzer.FieldTypeDouble
			}
		}

		outputSchema.AddField(aggName, aggType)
	}

	// Build GROUP BY including the time bucket
	// Note: In physical planning, we'll add a time bucket expression
	groupBy := make([]ast.Expression, 0, len(cmd.GroupBy)+1)

	// Add time bucket as first group (we'll use a special marker for now)
	groupBy = append(groupBy, &ast.FieldReference{Name: "_time"})

	// Add user-specified GROUP BY fields
	groupBy = append(groupBy, cmd.GroupBy...)

	return &LogicalAggregate{
		GroupBy:      groupBy,
		Aggregations: cmd.Aggregations,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildEvalCommand builds a LogicalEval operator
func (b *PlanBuilder) buildEvalCommand(cmd *ast.EvalCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("eval command requires an input plan")
	}

	// Build output schema by adding computed fields
	outputSchema := input.Schema().Clone()

	for _, assignment := range cmd.Assignments {
		// For now, infer type as double (would need type checker for proper inference)
		// TODO: Use type checker to infer expression type
		fieldType := analyzer.FieldTypeDouble

		// Simple type inference based on expression type
		switch assignment.Expression.(type) {
		case *ast.Literal:
			lit := assignment.Expression.(*ast.Literal)
			switch lit.LiteralTyp {
			case ast.LiteralTypeInt:
				fieldType = analyzer.FieldTypeLong
			case ast.LiteralTypeFloat:
				fieldType = analyzer.FieldTypeDouble
			case ast.LiteralTypeString:
				fieldType = analyzer.FieldTypeString
			case ast.LiteralTypeBool:
				fieldType = analyzer.FieldTypeBool
			}
		case *ast.FieldReference:
			ref := assignment.Expression.(*ast.FieldReference)
			if field, err := input.Schema().GetField(ref.Name); err == nil {
				fieldType = field.Type
			}
		}

		outputSchema.AddField(assignment.Field, fieldType)
	}

	return &LogicalEval{
		Assignments:  cmd.Assignments,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildRenameCommand builds a LogicalRename operator
func (b *PlanBuilder) buildRenameCommand(cmd *ast.RenameCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("rename command requires an input plan")
	}

	// Build output schema with renamed fields
	outputSchema := input.Schema().Clone()

	for _, assignment := range cmd.Assignments {
		// Get the old field
		if oldField, err := input.Schema().GetField(assignment.OldName); err == nil {
			// Remove old field and add new field with same type
			outputSchema.AddField(assignment.NewName, oldField.Type)
			// Note: In a more complete implementation, we'd track the rename mapping
		}
	}

	return &LogicalRename{
		Assignments:  cmd.Assignments,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildReplaceCommand builds a LogicalReplace operator
func (b *PlanBuilder) buildReplaceCommand(cmd *ast.ReplaceCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("replace command requires an input plan")
	}

	// Replace doesn't change the schema - it only modifies field values
	// Verify that the target field exists in the input schema
	if _, err := input.Schema().GetField(cmd.Field); err != nil {
		return nil, fmt.Errorf("field '%s' not found in schema: %w", cmd.Field, err)
	}

	return &LogicalReplace{
		Mappings: cmd.Mappings,
		Field:    cmd.Field,
		Input:    input,
	}, nil
}

// buildFillnullCommand builds a LogicalFillnull operator
func (b *PlanBuilder) buildFillnullCommand(cmd *ast.FillnullCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("fillnull command requires an input plan")
	}

	// Fillnull doesn't change the schema - it only fills missing/null values
	// If using per-field assignments, verify that the fields exist
	if len(cmd.Assignments) > 0 {
		for _, assignment := range cmd.Assignments {
			if _, err := input.Schema().GetField(assignment.Field); err != nil {
				// Field doesn't exist - this is ok, fillnull can add new fields
				// No error needed here
			}
		}
	}

	// If using fields list, verify those fields exist
	if len(cmd.Fields) > 0 {
		for _, fieldExpr := range cmd.Fields {
			if ref, ok := fieldExpr.(*ast.FieldReference); ok {
				if _, err := input.Schema().GetField(ref.Name); err != nil {
					// Field doesn't exist - warning but not an error
					// fillnull can still be applied to non-existent fields
				}
			}
		}
	}

	return &LogicalFillnull{
		Assignments:  cmd.Assignments,
		DefaultValue: cmd.DefaultValue,
		Fields:       cmd.Fields,
		Input:        input,
	}, nil
}

// buildParseCommand builds a LogicalParse operator
func (b *PlanBuilder) buildParseCommand(cmd *ast.ParseCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("parse command requires an input plan")
	}

	// Verify the source field exists in the input schema
	inputSchema := input.Schema()
	if _, err := inputSchema.GetField(cmd.SourceField); err != nil {
		return nil, fmt.Errorf("source field '%s' not found in schema: %w", cmd.SourceField, err)
	}

	// Extract named capture groups from the regex pattern
	extractedFields, err := extractNamedCaptureGroups(cmd.Pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	if len(extractedFields) == 0 {
		return nil, fmt.Errorf("parse pattern must contain at least one named capture group (e.g., (?<fieldname>pattern))")
	}

	// Build output schema: input schema + new extracted fields
	outputSchema := analyzer.NewSchema(inputSchema.Source)

	// Copy all existing fields from input
	for _, field := range inputSchema.Fields {
		outputSchema.AddField(field.Name, field.Type)
	}

	// Add new extracted fields (default to string type from regex)
	for _, fieldName := range extractedFields {
		// Check if field already exists
		if _, err := outputSchema.GetField(fieldName); err == nil {
			// Field already exists - parse will overwrite it
			// This is allowed behavior
		} else {
			outputSchema.AddField(fieldName, analyzer.FieldTypeString)
		}
	}

	return &LogicalParse{
		SourceField:     cmd.SourceField,
		Pattern:         cmd.Pattern,
		ExtractedFields: extractedFields,
		OutputSchema:    outputSchema,
		Input:           input,
	}, nil
}

// extractNamedCaptureGroups extracts field names from regex named capture groups
// Example: "(?<user>\w+) from (?<ip>\d+\.\d+\.\d+\.\d+)" returns ["user", "ip"]
func extractNamedCaptureGroups(pattern string) ([]string, error) {
	// Compile the regex to validate it
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Get the names of all subexpressions (capture groups)
	// The first element is the full match (empty name), rest are named groups
	names := re.SubexpNames()

	var extractedFields []string
	for i, name := range names {
		if i > 0 && name != "" {
			extractedFields = append(extractedFields, name)
		}
	}

	return extractedFields, nil
}

// buildRexCommand builds a LogicalRex operator
func (b *PlanBuilder) buildRexCommand(cmd *ast.RexCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("rex command requires an input plan")
	}

	// Determine source field (default to "_raw" if not specified)
	sourceField := cmd.SourceField
	if sourceField == "" {
		sourceField = "_raw"
	}

	// Verify the source field exists in the input schema
	inputSchema := input.Schema()
	if _, err := inputSchema.GetField(sourceField); err != nil {
		// Field doesn't exist - this is OK for rex, we'll just not extract anything
		// Rex is more lenient than parse
		// But we'll still log a warning during execution
	}

	// Extract named capture groups from the regex pattern
	extractedFields, err := extractNamedCaptureGroups(cmd.Pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	if len(extractedFields) == 0 {
		return nil, fmt.Errorf("rex pattern must contain at least one named capture group (e.g., (?<fieldname>pattern))")
	}

	// Build output schema: input schema + new extracted fields
	outputSchema := analyzer.NewSchema(inputSchema.Source)

	// Copy all existing fields from input
	for _, field := range inputSchema.Fields {
		outputSchema.AddField(field.Name, field.Type)
	}

	// Add new extracted fields (default to string type from regex)
	for _, fieldName := range extractedFields {
		// Check if field already exists
		if _, err := outputSchema.GetField(fieldName); err == nil {
			// Field already exists - rex will overwrite it
			// This is allowed behavior
		} else {
			outputSchema.AddField(fieldName, analyzer.FieldTypeString)
		}
	}

	return &LogicalRex{
		SourceField:     sourceField,
		Pattern:         cmd.Pattern,
		ExtractedFields: extractedFields,
		OutputSchema:    outputSchema,
		Input:           input,
	}, nil
}

// buildLookupCommand builds a LogicalLookup operator
func (b *PlanBuilder) buildLookupCommand(cmd *ast.LookupCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("lookup command requires an input plan")
	}

	// Verify the join field exists in the input schema
	inputSchema := input.Schema()
	if _, err := inputSchema.GetField(cmd.JoinField); err != nil {
		return nil, fmt.Errorf("join field '%s' not found in schema: %w", cmd.JoinField, err)
	}

	// Extract output field names and aliases
	outputFields := make([]string, len(cmd.OutputFields))
	outputAliases := make([]string, len(cmd.OutputFields))
	for i, field := range cmd.OutputFields {
		outputFields[i] = field.Field
		outputAliases[i] = field.Alias
	}

	// Build output schema: input schema + lookup output fields
	outputSchema := analyzer.NewSchema(inputSchema.Source)

	// Copy all existing fields from input
	for _, field := range inputSchema.Fields {
		outputSchema.AddField(field.Name, field.Type)
	}

	// Add lookup output fields (default to string type since we don't know the lookup table schema yet)
	for i, fieldName := range outputFields {
		outputFieldName := fieldName
		if outputAliases[i] != "" {
			outputFieldName = outputAliases[i]
		}

		// Check if field already exists
		if _, err := outputSchema.GetField(outputFieldName); err == nil {
			// Field already exists - lookup will overwrite it
			// This is allowed behavior
		} else {
			outputSchema.AddField(outputFieldName, analyzer.FieldTypeString)
		}
	}

	return &LogicalLookup{
		TableName:      cmd.TableName,
		JoinField:      cmd.JoinField,
		JoinFieldAlias: cmd.JoinFieldAlias,
		OutputFields:   outputFields,
		OutputAliases:  outputAliases,
		OutputSchema:   outputSchema,
		Input:          input,
	}, nil
}

// buildAppendCommand builds a LogicalAppend operator
func (b *PlanBuilder) buildAppendCommand(cmd *ast.AppendCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("append command requires an input plan")
	}

	if cmd.Subsearch == nil {
		return nil, fmt.Errorf("append command requires a subsearch")
	}

	// Build the subsearch plan independently
	// The subsearch starts from a SearchCommand, so we need a new PlanBuilder
	// with the subsearch's source schema
	if len(cmd.Subsearch.Commands) == 0 {
		return nil, fmt.Errorf("append subsearch is empty")
	}

	searchCmd, ok := cmd.Subsearch.Commands[0].(*ast.SearchCommand)
	if !ok {
		return nil, fmt.Errorf("append subsearch must start with a search command")
	}

	// Create a new schema for the subsearch source
	subsearchSchema := analyzer.NewSchema(searchCmd.Source)

	// If the original schema has fields, we should try to match them
	// For now, we'll use the base schema with no fields (fields discovered at runtime)

	// Build subsearch plan
	subsearchBuilder := NewPlanBuilder(subsearchSchema)
	subsearchPlan, err := subsearchBuilder.Build(cmd.Subsearch)
	if err != nil {
		return nil, fmt.Errorf("failed to build subsearch plan: %w", err)
	}

	// Merge schemas: union of all fields from both main query and subsearch
	inputSchema := input.Schema()
	subsearchResultSchema := subsearchPlan.Schema()

	outputSchema := analyzer.NewSchema(fmt.Sprintf("%s_UNION_%s",
		inputSchema.Source, subsearchResultSchema.Source))

	// Add all fields from input schema
	for _, field := range inputSchema.Fields {
		outputSchema.AddField(field.Name, field.Type)
	}

	// Add fields from subsearch that aren't already present
	for _, field := range subsearchResultSchema.Fields {
		if !outputSchema.HasField(field.Name) {
			outputSchema.AddField(field.Name, field.Type)
		}
	}

	return &LogicalAppend{
		Subsearch:    subsearchPlan,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildJoinCommand builds a LogicalJoin operator
func (b *PlanBuilder) buildJoinCommand(cmd *ast.JoinCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("join command requires an input plan")
	}

	if cmd.Subsearch == nil {
		return nil, fmt.Errorf("join command requires a subsearch")
	}

	// Verify the join field exists in the input schema
	inputSchema := input.Schema()
	if _, err := inputSchema.GetField(cmd.JoinField); err != nil {
		return nil, fmt.Errorf("join field '%s' not found in input schema: %w", cmd.JoinField, err)
	}

	// Build the right side (subsearch) plan independently
	if len(cmd.Subsearch.Commands) == 0 {
		return nil, fmt.Errorf("join subsearch is empty")
	}

	searchCmd, ok := cmd.Subsearch.Commands[0].(*ast.SearchCommand)
	if !ok {
		return nil, fmt.Errorf("join subsearch must start with a search command")
	}

	// Create a new schema for the subsearch source
	rightSchema := analyzer.NewSchema(searchCmd.Source)

	// Build subsearch plan
	rightBuilder := NewPlanBuilder(rightSchema)
	rightPlan, err := rightBuilder.Build(cmd.Subsearch)
	if err != nil {
		return nil, fmt.Errorf("failed to build join right side plan: %w", err)
	}

	// Verify join field exists in right side schema
	// Only validate if the schema has fields (skip for empty schemas used in tests)
	rightResultSchema := rightPlan.Schema()
	if len(rightResultSchema.Fields) > 0 {
		if _, err := rightResultSchema.GetField(cmd.JoinField); err != nil {
			return nil, fmt.Errorf("join field '%s' not found in right side schema: %w", cmd.JoinField, err)
		}
	}

	// Merge schemas: combine fields from both sides
	// For conflicting field names, suffix right side fields with "_right"
	outputSchema := analyzer.NewSchema(fmt.Sprintf("%s_JOIN_%s",
		inputSchema.Source, rightResultSchema.Source))

	// Add all fields from left side
	for _, field := range inputSchema.Fields {
		outputSchema.AddField(field.Name, field.Type)
	}

	// Add fields from right side
	// If field name conflicts with left side (and not the join key), add "_right" suffix
	for _, field := range rightResultSchema.Fields {
		if field.Name == cmd.JoinField {
			// Skip the join key field from right side (already in left)
			continue
		}

		fieldName := field.Name
		if outputSchema.HasField(fieldName) {
			// Conflict: add suffix
			fieldName = fieldName + "_right"
		}
		outputSchema.AddField(fieldName, field.Type)
	}

	return &LogicalJoin{
		JoinType:     cmd.JoinType,
		JoinField:    cmd.JoinField,
		RightField:   cmd.JoinField, // Same field name on both sides for now
		Right:        rightPlan,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildTableCommand builds a LogicalTable operator
func (b *PlanBuilder) buildTableCommand(cmd *ast.TableCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("table command requires an input plan")
	}

	// Build output schema with only the selected fields
	inputSchema := input.Schema()
	outputSchema := analyzer.NewSchema(inputSchema.Source)

	for _, fieldExpr := range cmd.Fields {
		if ref, ok := fieldExpr.(*ast.FieldReference); ok {
			// Get field type from input schema
			if field, err := inputSchema.GetField(ref.Name); err == nil {
				outputSchema.AddField(ref.Name, field.Type)
			} else {
				// Field doesn't exist in input - add as unknown type
				outputSchema.AddField(ref.Name, analyzer.FieldTypeString)
			}
		} else {
			// Complex expression - use string representation as field name
			fieldName := fieldExpr.String()
			outputSchema.AddField(fieldName, analyzer.FieldTypeString)
		}
	}

	return &LogicalTable{
		Fields:       cmd.Fields,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildEventstatsCommand builds a LogicalEventstats operator
func (b *PlanBuilder) buildEventstatsCommand(cmd *ast.EventstatsCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("eventstats command requires an input plan")
	}

	// Build output schema: input fields + aggregation result fields
	outputSchema := input.Schema().Clone()

	// Add aggregation result fields
	for _, agg := range cmd.Aggregations {
		aggName := agg.Alias
		if aggName == "" && agg.Func != nil {
			aggName = agg.Func.Name
			if len(agg.Func.Arguments) > 0 {
				if ref, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
					aggName += "_" + ref.Name
				}
			}
		}

		// Infer aggregation result type
		aggType := analyzer.FieldTypeLong
		if agg.Func != nil {
			switch agg.Func.Name {
			case "count":
				aggType = analyzer.FieldTypeLong
			case "sum", "avg":
				aggType = analyzer.FieldTypeDouble
			case "min", "max":
				// Preserve input field type for min/max
				if len(agg.Func.Arguments) > 0 {
					if ref, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
						if field, err := input.Schema().GetField(ref.Name); err == nil {
							aggType = field.Type
						}
					}
				}
			}
		}

		outputSchema.AddField(aggName, aggType)
	}

	return &LogicalEventstats{
		GroupBy:      cmd.GroupBy,
		Aggregations: cmd.Aggregations,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildStreamstatsCommand builds a LogicalStreamstats operator
func (b *PlanBuilder) buildStreamstatsCommand(cmd *ast.StreamstatsCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("streamstats command requires an input plan")
	}

	// Build output schema: input fields + aggregation result fields
	outputSchema := input.Schema().Clone()

	// Add aggregation result fields
	for _, agg := range cmd.Aggregations {
		aggName := agg.Alias
		if aggName == "" && agg.Func != nil {
			aggName = agg.Func.Name
			if len(agg.Func.Arguments) > 0 {
				if ref, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
					aggName += "_" + ref.Name
				}
			}
		}

		// Infer aggregation result type
		aggType := analyzer.FieldTypeLong
		if agg.Func != nil {
			switch agg.Func.Name {
			case "count":
				aggType = analyzer.FieldTypeLong
			case "sum", "avg":
				aggType = analyzer.FieldTypeDouble
			case "min", "max":
				// Preserve input field type for min/max
				if len(agg.Func.Arguments) > 0 {
					if ref, ok := agg.Func.Arguments[0].(*ast.FieldReference); ok {
						if field, err := input.Schema().GetField(ref.Name); err == nil {
							aggType = field.Type
						}
					}
				}
			}
		}

		outputSchema.AddField(aggName, aggType)
	}

	return &LogicalStreamstats{
		GroupBy:      cmd.GroupBy,
		Aggregations: cmd.Aggregations,
		Window:       cmd.Window,
		OutputSchema: outputSchema,
		Input:        input,
	}, nil
}

// buildReverseCommand builds a LogicalReverse operator
func (b *PlanBuilder) buildReverseCommand(cmd *ast.ReverseCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("reverse command requires an input plan")
	}

	// Reverse command does not modify the schema, it just changes row order
	// Pass through the input schema
	inputSchema := input.Schema()
	outputSchema := inputSchema.Clone()

	return &LogicalReverse{
		Input:        input,
		OutputSchema: outputSchema,
	}, nil
}

// buildFlattenCommand builds a LogicalFlatten operator
func (b *PlanBuilder) buildFlattenCommand(cmd *ast.FlattenCommand, input LogicalPlan) (LogicalPlan, error) {
	if input == nil {
		return nil, fmt.Errorf("flatten command requires an input plan")
	}

	if cmd.Field == nil {
		return nil, fmt.Errorf("flatten command requires a field expression")
	}

	// Flatten command does not modify the schema structure, just creates more rows
	// The field being flattened becomes scalar instead of array, but schema fields stay same
	// Pass through the input schema
	inputSchema := input.Schema()
	outputSchema := inputSchema.Clone()

	return &LogicalFlatten{
		Input:        input,
		Field:        cmd.Field,
		OutputSchema: outputSchema,
	}, nil
}
