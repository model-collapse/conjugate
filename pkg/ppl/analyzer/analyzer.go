// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package analyzer

import (
	"fmt"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
)

// Analyzer performs semantic analysis on PPL AST
type Analyzer struct {
	schema      *Schema
	scope       *Scope
	typeChecker *TypeChecker
}

// NewAnalyzer creates a new analyzer with the given schema
func NewAnalyzer(schema *Schema) *Analyzer {
	scope := NewScope(nil)
	return &Analyzer{
		schema:      schema,
		scope:       scope,
		typeChecker: NewTypeChecker(schema, scope),
	}
}

// Analyze performs semantic analysis on a query AST
// Returns an error if the query is semantically invalid
func (a *Analyzer) Analyze(query *ast.Query) error {
	for _, cmd := range query.Commands {
		if err := a.analyzeCommand(cmd); err != nil {
			return fmt.Errorf("semantic analysis failed: %w", err)
		}
	}
	return nil
}

// analyzeCommand analyzes a single command
func (a *Analyzer) analyzeCommand(cmd ast.Command) error {
	switch c := cmd.(type) {
	case *ast.SearchCommand:
		return a.analyzeSearchCommand(c)
	case *ast.WhereCommand:
		return a.analyzeWhereCommand(c)
	case *ast.FieldsCommand:
		return a.analyzeFieldsCommand(c)
	case *ast.StatsCommand:
		return a.analyzeStatsCommand(c)
	case *ast.SortCommand:
		return a.analyzeSortCommand(c)
	case *ast.HeadCommand:
		return a.analyzeHeadCommand(c)
	case *ast.DescribeCommand:
		return a.analyzeDescribeCommand(c)
	case *ast.ShowDatasourcesCommand:
		return nil // No validation needed
	case *ast.ExplainCommand:
		return nil // No validation needed
	// Tier 1 Commands
	case *ast.TopCommand:
		return a.analyzeTopCommand(c)
	case *ast.RareCommand:
		return a.analyzeRareCommand(c)
	case *ast.DedupCommand:
		return a.analyzeDedupCommand(c)
	case *ast.EvalCommand:
		return a.analyzeEvalCommand(c)
	case *ast.RenameCommand:
		return a.analyzeRenameCommand(c)
	case *ast.ReplaceCommand:
		return a.analyzeReplaceCommand(c)
	case *ast.BinCommand:
		return a.analyzeBinCommand(c)
	case *ast.TimechartCommand:
		return a.analyzeTimechartCommand(c)
	case *ast.ParseCommand:
		return a.analyzeParseCommand(c)
	case *ast.RexCommand:
		return a.analyzeRexCommand(c)
	case *ast.LookupCommand:
		return a.analyzeLookupCommand(c)
	case *ast.AppendCommand:
		return a.analyzeAppendCommand(c)
	case *ast.JoinCommand:
		return a.analyzeJoinCommand(c)
	case *ast.ReverseCommand:
		return a.analyzeReverseCommand(c)
	case *ast.FlattenCommand:
		return a.analyzeFlattenCommand(c)
	case *ast.FillnullCommand:
		return a.analyzeFillnullCommand(c)
	default:
		return fmt.Errorf("unknown command type: %T", cmd)
	}
}

// analyzeSearchCommand validates the search command
func (a *Analyzer) analyzeSearchCommand(cmd *ast.SearchCommand) error {
	if cmd.Source == "" {
		return fmt.Errorf("search command requires a source")
	}

	// Update scope with schema fields from the source
	if a.schema != nil && a.schema.Source == cmd.Source {
		for name, field := range a.schema.Fields {
			if err := a.scope.Define(name, field.Type); err != nil {
				// Field already defined - this is OK, skip it
				continue
			}
		}
	}

	return nil
}

// analyzeWhereCommand validates the where clause
func (a *Analyzer) analyzeWhereCommand(cmd *ast.WhereCommand) error {
	if cmd.Condition == nil {
		return fmt.Errorf("where command requires a condition expression")
	}

	// Validate condition expression
	if err := a.analyzeExpression(cmd.Condition); err != nil {
		return fmt.Errorf("invalid WHERE condition: %w", err)
	}

	// Check that condition evaluates to boolean
	condType, err := a.typeChecker.InferType(cmd.Condition)
	if err != nil {
		return fmt.Errorf("WHERE condition type inference failed: %w", err)
	}

	if condType != FieldTypeBool {
		return fmt.Errorf("WHERE condition must be boolean, got %s", condType)
	}

	return nil
}

// analyzeFieldsCommand validates the fields command
func (a *Analyzer) analyzeFieldsCommand(cmd *ast.FieldsCommand) error {
	if len(cmd.Fields) == 0 {
		return fmt.Errorf("fields command requires at least one field")
	}

	// Validate each field expression
	for i, field := range cmd.Fields {
		if err := a.analyzeExpression(field); err != nil {
			return fmt.Errorf("invalid field expression at position %d: %w", i, err)
		}
	}

	// Update scope with projected fields (remove excluded fields if Exclude=true)
	// For now, we'll keep all fields in scope and handle projection at execution time
	// TODO: Properly implement field projection in scope

	return nil
}

// analyzeStatsCommand validates the stats command
func (a *Analyzer) analyzeStatsCommand(cmd *ast.StatsCommand) error {
	if len(cmd.Aggregations) == 0 {
		return fmt.Errorf("stats command requires at least one aggregation")
	}

	// Create a new scope for the aggregation result
	aggScope := NewScope(a.scope)

	// Validate aggregations
	for i, agg := range cmd.Aggregations {
		// Validate the aggregation function call
		if agg.Func == nil {
			return fmt.Errorf("aggregation at position %d is missing function call", i)
		}

		if err := a.analyzeExpression(agg.Func); err != nil {
			return fmt.Errorf("invalid aggregation function at position %d: %w", i, err)
		}

		// Infer the result type
		aggType, err := a.typeChecker.InferType(agg.Func)
		if err != nil {
			return fmt.Errorf("aggregation type inference failed at position %d: %w", i, err)
		}

		// Add aggregation result to scope
		outputName := agg.Alias
		if outputName == "" {
			// Generate a default name
			outputName = fmt.Sprintf("agg_%d", i)
		}

		if err := aggScope.Define(outputName, aggType); err != nil {
			return fmt.Errorf("failed to define aggregation result %s: %w", outputName, err)
		}
	}

	// Validate GROUP BY fields
	for i, groupBy := range cmd.GroupBy {
		if err := a.analyzeExpression(groupBy); err != nil {
			return fmt.Errorf("invalid GROUP BY expression at position %d: %w", i, err)
		}

		// Add GROUP BY field to output scope
		if fieldRef, ok := groupBy.(*ast.FieldReference); ok {
			fieldType, err := a.typeChecker.InferType(fieldRef)
			if err != nil {
				return fmt.Errorf("GROUP BY field type inference failed: %w", err)
			}
			if err := aggScope.Define(fieldRef.Name, fieldType); err != nil {
				// Already defined, OK
			}
		}
	}

	// Update the analyzer's scope to the aggregation output scope
	a.scope = aggScope
	a.typeChecker = NewTypeChecker(a.schema, a.scope)

	return nil
}

// analyzeSortCommand validates the sort command
func (a *Analyzer) analyzeSortCommand(cmd *ast.SortCommand) error {
	if len(cmd.SortKeys) == 0 {
		return fmt.Errorf("sort command requires at least one sort key")
	}

	// Validate each sort key
	for i, sortKey := range cmd.SortKeys {
		// Validate the field expression
		if err := a.analyzeExpression(sortKey.Field); err != nil {
			return fmt.Errorf("invalid sort key at position %d: %w", i, err)
		}

		// Check that field type is comparable/sortable
		fieldType, err := a.typeChecker.InferType(sortKey.Field)
		if err != nil {
			return fmt.Errorf("sort key type inference failed at position %d: %w", i, err)
		}

		if !fieldType.IsComparable() {
			return fmt.Errorf("sort key at position %d has non-comparable type: %s", i, fieldType)
		}
	}

	return nil
}

// analyzeHeadCommand validates the head command
func (a *Analyzer) analyzeHeadCommand(cmd *ast.HeadCommand) error {
	if cmd.Count <= 0 {
		return fmt.Errorf("head command requires a positive count, got %d", cmd.Count)
	}
	return nil
}

// analyzeDescribeCommand validates the describe command
func (a *Analyzer) analyzeDescribeCommand(cmd *ast.DescribeCommand) error {
	if cmd.Source == "" {
		return fmt.Errorf("describe command requires a source")
	}
	return nil
}

// analyzeExpression performs recursive validation on expressions
func (a *Analyzer) analyzeExpression(expr ast.Expression) error {
	switch e := expr.(type) {
	case *ast.Literal:
		return nil // Literals are always valid

	case *ast.FieldReference:
		return a.analyzeFieldReference(e)

	case *ast.BinaryExpression:
		return a.analyzeBinaryExpression(e)

	case *ast.UnaryExpression:
		return a.analyzeUnaryExpression(e)

	case *ast.FunctionCall:
		return a.analyzeFunctionCall(e)

	case *ast.CaseExpression:
		return a.analyzeCaseExpression(e)

	case *ast.ListLiteral:
		return a.analyzeListLiteral(e)

	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}

// analyzeFieldReference validates a field reference
func (a *Analyzer) analyzeFieldReference(ref *ast.FieldReference) error {
	// Check if field exists in scope (for aliases)
	if a.scope.Has(ref.Name) {
		return nil
	}

	// Check if field exists in schema
	if a.schema != nil {
		if !a.schema.HasField(ref.Name) {
			return fmt.Errorf("field %s not found in schema", ref.Name)
		}
	}

	return nil
}

// analyzeBinaryExpression validates a binary expression
func (a *Analyzer) analyzeBinaryExpression(expr *ast.BinaryExpression) error {
	// Validate left operand
	if err := a.analyzeExpression(expr.Left); err != nil {
		return fmt.Errorf("invalid left operand: %w", err)
	}

	// Validate right operand
	if err := a.analyzeExpression(expr.Right); err != nil {
		return fmt.Errorf("invalid right operand: %w", err)
	}

	// Type check will be done by type checker
	_, err := a.typeChecker.InferType(expr)
	if err != nil {
		return fmt.Errorf("type checking failed for binary expression: %w", err)
	}

	return nil
}

// analyzeUnaryExpression validates a unary expression
func (a *Analyzer) analyzeUnaryExpression(expr *ast.UnaryExpression) error {
	// Validate operand
	if err := a.analyzeExpression(expr.Operand); err != nil {
		return fmt.Errorf("invalid operand: %w", err)
	}

	// Type check
	_, err := a.typeChecker.InferType(expr)
	if err != nil {
		return fmt.Errorf("type checking failed for unary expression: %w", err)
	}

	return nil
}

// analyzeFunctionCall validates a function call
func (a *Analyzer) analyzeFunctionCall(call *ast.FunctionCall) error {
	// Validate all arguments
	for i, arg := range call.Arguments {
		if err := a.analyzeExpression(arg); err != nil {
			return fmt.Errorf("invalid argument at position %d for function %s: %w", i, call.Name, err)
		}
	}

	// Type check (will validate function exists and argument types)
	_, err := a.typeChecker.InferType(call)
	if err != nil {
		return fmt.Errorf("type checking failed for function %s: %w", call.Name, err)
	}

	return nil
}

// analyzeCaseExpression validates a CASE expression
func (a *Analyzer) analyzeCaseExpression(caseExpr *ast.CaseExpression) error {
	// Validate all WHEN clauses
	for i, whenClause := range caseExpr.WhenClauses {
		// Validate condition
		if err := a.analyzeExpression(whenClause.Condition); err != nil {
			return fmt.Errorf("invalid WHEN condition at position %d: %w", i, err)
		}

		// Validate result expression
		if err := a.analyzeExpression(whenClause.Result); err != nil {
			return fmt.Errorf("invalid THEN expression at position %d: %w", i, err)
		}
	}

	// Validate ELSE clause if present
	if caseExpr.ElseResult != nil {
		if err := a.analyzeExpression(caseExpr.ElseResult); err != nil {
			return fmt.Errorf("invalid ELSE expression: %w", err)
		}
	}

	// Type check
	_, err := a.typeChecker.InferType(caseExpr)
	if err != nil {
		return fmt.Errorf("type checking failed for CASE expression: %w", err)
	}

	return nil
}

// analyzeListLiteral validates a list literal
func (a *Analyzer) analyzeListLiteral(list *ast.ListLiteral) error {
	// Validate all elements
	for i, elem := range list.Values {
		if err := a.analyzeExpression(elem); err != nil {
			return fmt.Errorf("invalid list element at position %d: %w", i, err)
		}
	}

	// All elements should have compatible types
	if len(list.Values) > 0 {
		firstType, err := a.typeChecker.InferType(list.Values[0])
		if err != nil {
			return fmt.Errorf("type inference failed for list element 0: %w", err)
		}

		for i := 1; i < len(list.Values); i++ {
			elemType, err := a.typeChecker.InferType(list.Values[i])
			if err != nil {
				return fmt.Errorf("type inference failed for list element %d: %w", i, err)
			}

			// Check type compatibility (allow some flexibility)
			if firstType.IsNumeric() && !elemType.IsNumeric() {
				return fmt.Errorf("list element %d has incompatible type: expected numeric, got %s", i, elemType)
			}
			if firstType.IsString() && !elemType.IsString() {
				return fmt.Errorf("list element %d has incompatible type: expected string, got %s", i, elemType)
			}
		}
	}

	return nil
}

// GetScope returns the current scope
func (a *Analyzer) GetScope() *Scope {
	return a.scope
}

// GetSchema returns the schema
func (a *Analyzer) GetSchema() *Schema {
	return a.schema
}

// GetTypeChecker returns the type checker
func (a *Analyzer) GetTypeChecker() *TypeChecker {
	return a.typeChecker
}

// =====================================================================
// Tier 1 Command Analysis Methods
// =====================================================================

// analyzeTopCommand validates the top command
func (a *Analyzer) analyzeTopCommand(cmd *ast.TopCommand) error {
	if len(cmd.Fields) == 0 {
		return fmt.Errorf("top command requires at least one field")
	}

	// Validate field expressions
	for i, field := range cmd.Fields {
		if err := a.analyzeExpression(field); err != nil {
			return fmt.Errorf("invalid top field at position %d: %w", i, err)
		}
	}

	// Validate group by fields
	for i, groupBy := range cmd.GroupBy {
		if err := a.analyzeExpression(groupBy); err != nil {
			return fmt.Errorf("invalid GROUP BY field at position %d: %w", i, err)
		}
	}

	// Update scope with output fields
	newScope := NewScope(a.scope)
	for _, field := range cmd.Fields {
		if fieldRef, ok := field.(*ast.FieldReference); ok {
			fieldType, err := a.typeChecker.InferType(fieldRef)
			if err == nil {
				newScope.Define(fieldRef.Name, fieldType)
			}
		}
	}
	// Add count field
	newScope.Define("count", FieldTypeLong)
	// Add percent field if requested
	if cmd.ShowPercent {
		newScope.Define("percent", FieldTypeDouble)
	}

	a.scope = newScope
	a.typeChecker = NewTypeChecker(a.schema, a.scope)

	return nil
}

// analyzeRareCommand validates the rare command
func (a *Analyzer) analyzeRareCommand(cmd *ast.RareCommand) error {
	if len(cmd.Fields) == 0 {
		return fmt.Errorf("rare command requires at least one field")
	}

	// Validate field expressions
	for i, field := range cmd.Fields {
		if err := a.analyzeExpression(field); err != nil {
			return fmt.Errorf("invalid rare field at position %d: %w", i, err)
		}
	}

	// Validate group by fields
	for i, groupBy := range cmd.GroupBy {
		if err := a.analyzeExpression(groupBy); err != nil {
			return fmt.Errorf("invalid GROUP BY field at position %d: %w", i, err)
		}
	}

	// Update scope with output fields
	newScope := NewScope(a.scope)
	for _, field := range cmd.Fields {
		if fieldRef, ok := field.(*ast.FieldReference); ok {
			fieldType, err := a.typeChecker.InferType(fieldRef)
			if err == nil {
				newScope.Define(fieldRef.Name, fieldType)
			}
		}
	}
	// Add count field
	newScope.Define("count", FieldTypeLong)
	// Add percent field if requested
	if cmd.ShowPercent {
		newScope.Define("percent", FieldTypeDouble)
	}

	a.scope = newScope
	a.typeChecker = NewTypeChecker(a.schema, a.scope)

	return nil
}

// analyzeDedupCommand validates the dedup command
func (a *Analyzer) analyzeDedupCommand(cmd *ast.DedupCommand) error {
	if len(cmd.Fields) == 0 {
		return fmt.Errorf("dedup command requires at least one field")
	}

	// Validate field expressions
	for i, field := range cmd.Fields {
		if err := a.analyzeExpression(field); err != nil {
			return fmt.Errorf("invalid dedup field at position %d: %w", i, err)
		}
	}

	// Validate count
	if cmd.Count < 0 {
		return fmt.Errorf("dedup count cannot be negative, got %d", cmd.Count)
	}

	// Dedup doesn't change the schema, so we don't update scope

	return nil
}

// analyzeEvalCommand validates the eval command
func (a *Analyzer) analyzeEvalCommand(cmd *ast.EvalCommand) error {
	if len(cmd.Assignments) == 0 {
		return fmt.Errorf("eval command requires at least one assignment")
	}

	// Validate each assignment
	for i, assignment := range cmd.Assignments {
		// Validate the expression
		if err := a.analyzeExpression(assignment.Expression); err != nil {
			return fmt.Errorf("invalid eval expression at position %d: %w", i, err)
		}

		// Infer the result type
		exprType, err := a.typeChecker.InferType(assignment.Expression)
		if err != nil {
			return fmt.Errorf("eval type inference failed at position %d: %w", i, err)
		}

		// Add the new field to scope
		if err := a.scope.Define(assignment.Field, exprType); err != nil {
			// Field already exists, update its type
			a.scope.Update(assignment.Field, exprType)
		}
	}

	// Recreate type checker with updated scope
	a.typeChecker = NewTypeChecker(a.schema, a.scope)

	return nil
}

// analyzeRenameCommand validates the rename command
func (a *Analyzer) analyzeRenameCommand(cmd *ast.RenameCommand) error {
	if len(cmd.Assignments) == 0 {
		return fmt.Errorf("rename command requires at least one assignment")
	}

	// Validate and update scope for each rename
	for i, assignment := range cmd.Assignments {
		// Check that old field exists
		oldFieldType := a.scope.Lookup(assignment.OldName)
		if oldFieldType == nil && a.schema != nil {
			field, err := a.schema.GetField(assignment.OldName)
			if err != nil {
				return fmt.Errorf("rename source field %s not found at position %d", assignment.OldName, i)
			}
			oldFieldType = &field.Type
		}

		if oldFieldType == nil {
			return fmt.Errorf("rename source field %s not found at position %d", assignment.OldName, i)
		}

		// Add new field name with same type
		if err := a.scope.Define(assignment.NewName, *oldFieldType); err != nil {
			// Field already exists, update its type
			a.scope.Update(assignment.NewName, *oldFieldType)
		}

		// Note: We don't remove the old name from scope since it may still be
		// referenced in subsequent pipeline stages (e.g., "rename a as b | where a > 10")
	}

	// Recreate type checker with updated scope
	a.typeChecker = NewTypeChecker(a.schema, a.scope)

	return nil
}

// analyzeReplaceCommand validates the replace command
func (a *Analyzer) analyzeReplaceCommand(cmd *ast.ReplaceCommand) error {
	if len(cmd.Mappings) == 0 {
		return fmt.Errorf("replace command requires at least one replacement mapping")
	}

	if cmd.Field == "" {
		return fmt.Errorf("replace command requires a target field")
	}

	// Check that target field exists
	targetFieldType := a.scope.Lookup(cmd.Field)
	if targetFieldType == nil && a.schema != nil {
		field, err := a.schema.GetField(cmd.Field)
		if err != nil {
			return fmt.Errorf("replace target field %s not found", cmd.Field)
		}
		targetFieldType = &field.Type
	}

	if targetFieldType == nil {
		return fmt.Errorf("replace target field %s not found", cmd.Field)
	}

	// Validate each mapping
	for i, mapping := range cmd.Mappings {
		// Validate old value expression
		if err := a.analyzeExpression(mapping.OldValue); err != nil {
			return fmt.Errorf("invalid old value expression in mapping %d: %w", i, err)
		}

		// Validate new value expression
		if err := a.analyzeExpression(mapping.NewValue); err != nil {
			return fmt.Errorf("invalid new value expression in mapping %d: %w", i, err)
		}
	}

	// Replace doesn't change the schema, field remains the same type
	return nil
}

// analyzeParseCommand validates the parse command
func (a *Analyzer) analyzeParseCommand(cmd *ast.ParseCommand) error {
	if cmd.SourceField == "" {
		return fmt.Errorf("parse command requires a source field")
	}

	if cmd.Pattern == "" {
		return fmt.Errorf("parse command requires a pattern")
	}

	// Check that source field exists
	sourceFieldType := a.scope.Lookup(cmd.SourceField)
	if sourceFieldType == nil && a.schema != nil {
		field, err := a.schema.GetField(cmd.SourceField)
		if err != nil {
			return fmt.Errorf("parse source field %s not found", cmd.SourceField)
		}
		sourceFieldType = &field.Type
	}

	if sourceFieldType == nil {
		return fmt.Errorf("parse source field %s not found", cmd.SourceField)
	}

	// Source field should be string-like for regex matching
	if sourceFieldType != nil && !sourceFieldType.IsString() {
		return fmt.Errorf("parse source field %s must be a string type, got %s",
			cmd.SourceField, sourceFieldType.String())
	}

	// Parse doesn't change existing fields, but will add extracted fields
	// The extracted fields will be determined during planning phase when the regex is compiled
	return nil
}

// analyzeRexCommand validates the rex command
func (a *Analyzer) analyzeRexCommand(cmd *ast.RexCommand) error {
	if cmd.Pattern == "" {
		return fmt.Errorf("rex command requires a pattern")
	}

	// Rex allows empty source field (defaults to _raw during planning)
	if cmd.SourceField != "" {
		// Check that source field exists
		sourceFieldType := a.scope.Lookup(cmd.SourceField)
		if sourceFieldType == nil && a.schema != nil {
			field, err := a.schema.GetField(cmd.SourceField)
			if err != nil {
				// Rex is more lenient than parse - field not existing is OK
				// It just won't extract anything at runtime
			} else {
				sourceFieldType = &field.Type
			}
		}

		// Source field should be string-like for regex matching
		if sourceFieldType != nil && !sourceFieldType.IsString() {
			return fmt.Errorf("rex source field %s must be a string type, got %s",
				cmd.SourceField, sourceFieldType.String())
		}
	}

	// Rex doesn't change existing fields, but will add extracted fields
	// The extracted fields will be determined during planning phase when the regex is compiled
	return nil
}

// analyzeLookupCommand validates the lookup command
func (a *Analyzer) analyzeLookupCommand(cmd *ast.LookupCommand) error {
	if cmd.TableName == "" {
		return fmt.Errorf("lookup command requires a table name")
	}

	if cmd.JoinField == "" {
		return fmt.Errorf("lookup command requires a join field")
	}

	if len(cmd.OutputFields) == 0 {
		return fmt.Errorf("lookup command requires at least one output field")
	}

	// Check that join field exists
	joinFieldType := a.scope.Lookup(cmd.JoinField)
	if joinFieldType == nil && a.schema != nil {
		field, err := a.schema.GetField(cmd.JoinField)
		if err != nil {
			return fmt.Errorf("lookup join field %s not found", cmd.JoinField)
		}
		joinFieldType = &field.Type
	}

	if joinFieldType == nil {
		return fmt.Errorf("lookup join field %s not found", cmd.JoinField)
	}

	// Add output fields to the scope so subsequent commands can reference them
	// The actual lookup table schema validation happens at runtime
	for _, outputField := range cmd.OutputFields {
		// Use the alias if provided, otherwise use the field name
		fieldName := outputField.Field
		if outputField.Alias != "" {
			fieldName = outputField.Alias
		}

		// Add to scope with unknown type (actual type determined at runtime)
		fieldType := FieldTypeUnknown
		if err := a.scope.Define(fieldName, fieldType); err != nil {
			// Field already exists, update its type
			a.scope.Update(fieldName, fieldType)
		}
	}

	// Recreate type checker with updated scope
	a.typeChecker = NewTypeChecker(a.schema, a.scope)

	return nil
}

// analyzeAppendCommand validates the append command
func (a *Analyzer) analyzeAppendCommand(cmd *ast.AppendCommand) error {
	if cmd.Subsearch == nil {
		return fmt.Errorf("append command requires a subsearch query")
	}

	// The subsearch is independent and starts from a different source
	// We need to analyze it separately with its own schema/scope
	// For now, we'll validate that it has at least a search command
	if len(cmd.Subsearch.Commands) == 0 {
		return fmt.Errorf("append subsearch requires at least a search command")
	}

	// Validate that the first command is a search command
	if _, ok := cmd.Subsearch.Commands[0].(*ast.SearchCommand); !ok {
		return fmt.Errorf("append subsearch must start with a search command")
	}

	// Note: Full subsearch analysis happens at planning/execution time
	// because it has its own source and schema
	return nil
}

// analyzeJoinCommand validates the join command
func (a *Analyzer) analyzeJoinCommand(cmd *ast.JoinCommand) error {
	if cmd.JoinField == "" {
		return fmt.Errorf("join command requires a join field")
	}

	if cmd.Subsearch == nil {
		return fmt.Errorf("join command requires a subsearch query")
	}

	// Validate that the join field exists in the current schema
	joinFieldType := a.scope.Lookup(cmd.JoinField)
	if joinFieldType == nil && a.schema != nil {
		field, err := a.schema.GetField(cmd.JoinField)
		if err != nil {
			return fmt.Errorf("join field '%s' not found in schema: %w", cmd.JoinField, err)
		}
		joinFieldType = &field.Type
	}

	if joinFieldType == nil {
		return fmt.Errorf("join field '%s' not found", cmd.JoinField)
	}

	// Validate subsearch structure
	if len(cmd.Subsearch.Commands) == 0 {
		return fmt.Errorf("join subsearch requires at least a search command")
	}

	// Validate that the first command is a search command
	if _, ok := cmd.Subsearch.Commands[0].(*ast.SearchCommand); !ok {
		return fmt.Errorf("join subsearch must start with a search command")
	}

	// Validate join type
	switch cmd.JoinType {
	case ast.JoinTypeInner, ast.JoinTypeLeft, ast.JoinTypeRight, ast.JoinTypeOuter, ast.JoinTypeFull:
		// Valid join types
	default:
		return fmt.Errorf("invalid join type: %s", cmd.JoinType)
	}

	// Note: Full subsearch analysis and schema merging happens at planning time
	// The subsearch has its own source and schema that will be analyzed independently
	return nil
}

// analyzeBinCommand validates the bin command
func (a *Analyzer) analyzeBinCommand(cmd *ast.BinCommand) error {
	// Validate the field expression
	if err := a.analyzeExpression(cmd.Field); err != nil {
		return fmt.Errorf("invalid bin field: %w", err)
	}

	// Check field type - should be numeric or date
	fieldType, err := a.typeChecker.InferType(cmd.Field)
	if err != nil {
		return fmt.Errorf("bin field type inference failed: %w", err)
	}

	if !fieldType.IsNumeric() && fieldType != FieldTypeDate {
		return fmt.Errorf("bin field must be numeric or date, got %s", fieldType)
	}

	// Validate span or bins parameter
	if cmd.Span == nil && cmd.Bins == 0 {
		return fmt.Errorf("bin command requires either span or bins parameter")
	}

	// Bin replaces the field value with the binned value, so type stays the same
	// No scope changes needed

	return nil
}

// analyzeTimechartCommand validates the timechart command
func (a *Analyzer) analyzeTimechartCommand(cmd *ast.TimechartCommand) error {
	if len(cmd.Aggregations) == 0 {
		return fmt.Errorf("timechart command requires at least one aggregation")
	}

	// Validate span
	if cmd.Span == nil {
		return fmt.Errorf("timechart command requires a span parameter")
	}

	// Create a new scope for the aggregation result
	aggScope := NewScope(a.scope)

	// Add _time field for the time bucket
	aggScope.Define("_time", FieldTypeDate)

	// Validate aggregations
	for i, agg := range cmd.Aggregations {
		// Validate the aggregation function call
		if agg.Func == nil {
			return fmt.Errorf("aggregation at position %d is missing function call", i)
		}

		if err := a.analyzeExpression(agg.Func); err != nil {
			return fmt.Errorf("invalid aggregation function at position %d: %w", i, err)
		}

		// Infer the result type
		aggType, err := a.typeChecker.InferType(agg.Func)
		if err != nil {
			return fmt.Errorf("aggregation type inference failed at position %d: %w", i, err)
		}

		// Add aggregation result to scope
		outputName := agg.Alias
		if outputName == "" {
			// Generate a default name
			outputName = fmt.Sprintf("agg_%d", i)
		}

		if err := aggScope.Define(outputName, aggType); err != nil {
			return fmt.Errorf("failed to define aggregation result %s: %w", outputName, err)
		}
	}

	// Validate GROUP BY fields (in addition to time)
	for i, groupBy := range cmd.GroupBy {
		if err := a.analyzeExpression(groupBy); err != nil {
			return fmt.Errorf("invalid GROUP BY expression at position %d: %w", i, err)
		}

		// Add GROUP BY field to output scope
		if fieldRef, ok := groupBy.(*ast.FieldReference); ok {
			fieldType, err := a.typeChecker.InferType(fieldRef)
			if err != nil {
				return fmt.Errorf("GROUP BY field type inference failed: %w", err)
			}
			if err := aggScope.Define(fieldRef.Name, fieldType); err != nil {
				// Already defined, OK
			}
		}
	}

	// Update the analyzer's scope to the aggregation output scope
	a.scope = aggScope
	a.typeChecker = NewTypeChecker(a.schema, a.scope)

	return nil
}

// analyzeReverseCommand validates the reverse command
func (a *Analyzer) analyzeReverseCommand(cmd *ast.ReverseCommand) error {
	// Reverse command has no parameters, so minimal validation needed
	// It simply reverses the order of rows in the result set
	return nil
}

// analyzeFlattenCommand validates the flatten command
func (a *Analyzer) analyzeFlattenCommand(cmd *ast.FlattenCommand) error {
	// Validate the field expression
	if cmd.Field == nil {
		return fmt.Errorf("flatten command requires a field expression")
	}

	// Analyze the field expression
	// The field should be a field reference (simple or nested)
	// We don't check if it exists in the schema - that's runtime behavior
	// But we do check that it's a valid expression
	switch fieldExpr := cmd.Field.(type) {
	case *ast.FieldReference:
		// Valid - simple or nested field reference
		if fieldExpr.Name == "" {
			return fmt.Errorf("flatten field reference has empty name")
		}
	case *ast.FunctionCall:
		// Could be valid if it returns an array/object
		// Allow it for now
	default:
		// Other expression types might be valid too
		// Be permissive - let runtime handle it
	}

	return nil
}

// analyzeFillnullCommand validates the fillnull command
func (a *Analyzer) analyzeFillnullCommand(cmd *ast.FillnullCommand) error {
	// Validate the default value expression
	if cmd.DefaultValue == nil {
		return fmt.Errorf("fillnull command requires a default value expression")
	}

	// The value should be a literal (string, number, bool, or null)
	// We validate this to ensure it's a constant value, not a complex expression
	if _, ok := cmd.DefaultValue.(*ast.Literal); !ok {
		return fmt.Errorf("fillnull value must be a literal (string, number, boolean, or null)")
	}

	// Validate field expressions if provided
	for i, fieldExpr := range cmd.Fields {
		if fieldExpr == nil {
			return fmt.Errorf("fillnull field at position %d is nil", i)
		}
		// Fields should be field references
		if _, ok := fieldExpr.(*ast.FieldReference); !ok {
			return fmt.Errorf("fillnull field at position %d must be a field reference, got %T", i, fieldExpr)
		}
	}

	// Note: We don't validate that fields exist in the schema
	// because schema might change during execution (e.g., after eval commands)
	// Runtime will handle missing fields gracefully

	return nil
}
