// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package planner

import (
	"fmt"
	"strings"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
)

// LogicalPlan represents a logical query execution plan
type LogicalPlan interface {
	// Schema returns the output schema of this operator
	Schema() *analyzer.Schema

	// Children returns the child operators
	Children() []LogicalPlan

	// String returns a string representation
	String() string
}

// LogicalScan reads data from a source (index)
type LogicalScan struct {
	Source      string
	OutputSchema *analyzer.Schema
}

func (l *LogicalScan) Schema() *analyzer.Schema  { return l.OutputSchema }
func (l *LogicalScan) Children() []LogicalPlan   { return nil }
func (l *LogicalScan) String() string {
	return fmt.Sprintf("Scan(%s)", l.Source)
}

// LogicalFilter filters rows based on a predicate
type LogicalFilter struct {
	Condition ast.Expression
	Input     LogicalPlan
}

func (l *LogicalFilter) Schema() *analyzer.Schema { return l.Input.Schema() }
func (l *LogicalFilter) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalFilter) String() string {
	return fmt.Sprintf("Filter(%s)", l.Condition.String())
}

// LogicalProject selects specific fields
type LogicalProject struct {
	Fields       []ast.Expression
	OutputSchema *analyzer.Schema
	Input        LogicalPlan
	Exclude      bool // True if fields should be excluded
}

func (l *LogicalProject) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalProject) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalProject) String() string {
	fieldNames := make([]string, len(l.Fields))
	for i, f := range l.Fields {
		fieldNames[i] = f.String()
	}
	prefix := ""
	if l.Exclude {
		prefix = "-"
	}
	return fmt.Sprintf("Project(%s%s)", prefix, strings.Join(fieldNames, ", "))
}

// LogicalSort sorts rows by one or more keys
type LogicalSort struct {
	SortKeys []*ast.SortKey
	Input    LogicalPlan
}

func (l *LogicalSort) Schema() *analyzer.Schema { return l.Input.Schema() }
func (l *LogicalSort) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalSort) String() string {
	keys := make([]string, len(l.SortKeys))
	for i, k := range l.SortKeys {
		order := "ASC"
		if k.Descending {
			order = "DESC"
		}
		keys[i] = fmt.Sprintf("%s %s", k.Field.String(), order)
	}
	return fmt.Sprintf("Sort(%s)", strings.Join(keys, ", "))
}

// LogicalLimit limits the number of output rows
type LogicalLimit struct {
	Count int
	Input LogicalPlan
}

func (l *LogicalLimit) Schema() *analyzer.Schema { return l.Input.Schema() }
func (l *LogicalLimit) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalLimit) String() string {
	return fmt.Sprintf("Limit(%d)", l.Count)
}

// LogicalAggregate performs grouping and aggregation
type LogicalAggregate struct {
	GroupBy      []ast.Expression
	Aggregations []*ast.Aggregation
	OutputSchema *analyzer.Schema
	Input        LogicalPlan
}

func (l *LogicalAggregate) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalAggregate) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalAggregate) String() string {
	aggs := make([]string, len(l.Aggregations))
	for i, agg := range l.Aggregations {
		aggs[i] = agg.String()
	}

	if len(l.GroupBy) == 0 {
		return fmt.Sprintf("Aggregate(%s)", strings.Join(aggs, ", "))
	}

	groups := make([]string, len(l.GroupBy))
	for i, g := range l.GroupBy {
		groups[i] = g.String()
	}
	return fmt.Sprintf("Aggregate(%s, group_by=[%s])", strings.Join(aggs, ", "), strings.Join(groups, ", "))
}

// LogicalExplain wraps a plan to explain it
type LogicalExplain struct {
	Input LogicalPlan
}

func (l *LogicalExplain) Schema() *analyzer.Schema { return l.Input.Schema() }
func (l *LogicalExplain) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalExplain) String() string {
	return "Explain"
}

// Helper function to print plan tree
func PrintPlan(plan LogicalPlan, indent int) string {
	prefix := strings.Repeat("  ", indent)
	result := prefix + plan.String() + "\n"

	for _, child := range plan.Children() {
		result += PrintPlan(child, indent+1)
	}

	return result
}

// GetLeafScans returns all LogicalScan nodes in the plan
func GetLeafScans(plan LogicalPlan) []*LogicalScan {
	scans := make([]*LogicalScan, 0)

	if scan, ok := plan.(*LogicalScan); ok {
		scans = append(scans, scan)
		return scans
	}

	for _, child := range plan.Children() {
		scans = append(scans, GetLeafScans(child)...)
	}

	return scans
}

// ReplaceChild replaces a child operator with a new one
func ReplaceChild(plan LogicalPlan, oldChild, newChild LogicalPlan) LogicalPlan {
	switch p := plan.(type) {
	case *LogicalFilter:
		if p.Input == oldChild {
			return &LogicalFilter{Condition: p.Condition, Input: newChild}
		}

	case *LogicalProject:
		if p.Input == oldChild {
			return &LogicalProject{
				Fields:       p.Fields,
				OutputSchema: p.OutputSchema,
				Input:        newChild,
				Exclude:      p.Exclude,
			}
		}

	case *LogicalSort:
		if p.Input == oldChild {
			return &LogicalSort{SortKeys: p.SortKeys, Input: newChild}
		}

	case *LogicalLimit:
		if p.Input == oldChild {
			return &LogicalLimit{Count: p.Count, Input: newChild}
		}

	case *LogicalAggregate:
		if p.Input == oldChild {
			return &LogicalAggregate{
				GroupBy:      p.GroupBy,
				Aggregations: p.Aggregations,
				OutputSchema: p.OutputSchema,
				Input:        newChild,
			}
		}

	case *LogicalExplain:
		if p.Input == oldChild {
			return &LogicalExplain{Input: newChild}
		}

	case *LogicalDedup:
		if p.Input == oldChild {
			return &LogicalDedup{
				Fields:      p.Fields,
				Count:       p.Count,
				Consecutive: p.Consecutive,
				Input:       newChild,
			}
		}

	case *LogicalBin:
		if p.Input == oldChild {
			return &LogicalBin{
				Field: p.Field,
				Span:  p.Span,
				Bins:  p.Bins,
				Input: newChild,
			}
		}

	case *LogicalTop:
		if p.Input == oldChild {
			return &LogicalTop{
				Fields:       p.Fields,
				Limit:        p.Limit,
				GroupBy:      p.GroupBy,
				ShowCount:    p.ShowCount,
				ShowPercent:  p.ShowPercent,
				OutputSchema: p.OutputSchema,
				Input:        newChild,
			}
		}

	case *LogicalRare:
		if p.Input == oldChild {
			return &LogicalRare{
				Fields:       p.Fields,
				Limit:        p.Limit,
				GroupBy:      p.GroupBy,
				ShowCount:    p.ShowCount,
				ShowPercent:  p.ShowPercent,
				OutputSchema: p.OutputSchema,
				Input:        newChild,
			}
		}

	case *LogicalEval:
		if p.Input == oldChild {
			return &LogicalEval{
				Assignments:  p.Assignments,
				OutputSchema: p.OutputSchema,
				Input:        newChild,
			}
		}

	case *LogicalRename:
		if p.Input == oldChild {
			return &LogicalRename{
				Assignments:  p.Assignments,
				OutputSchema: p.OutputSchema,
				Input:        newChild,
			}
		}

	case *LogicalReplace:
		if p.Input == oldChild {
			return &LogicalReplace{
				Mappings: p.Mappings,
				Field:    p.Field,
				Input:    newChild,
			}
		}

	case *LogicalFillnull:
		if p.Input == oldChild {
			return &LogicalFillnull{
				Assignments:  p.Assignments,
				DefaultValue: p.DefaultValue,
				Fields:       p.Fields,
				Input:        newChild,
			}
		}

	case *LogicalTable:
		if p.Input == oldChild {
			return &LogicalTable{
				Fields:       p.Fields,
				OutputSchema: p.OutputSchema,
				Input:        newChild,
			}
		}

	case *LogicalEventstats:
		if p.Input == oldChild {
			return &LogicalEventstats{
				GroupBy:      p.GroupBy,
				Aggregations: p.Aggregations,
				OutputSchema: p.OutputSchema,
				Input:        newChild,
			}
		}

	case *LogicalStreamstats:
		if p.Input == oldChild {
			return &LogicalStreamstats{
				GroupBy:      p.GroupBy,
				Aggregations: p.Aggregations,
				Window:       p.Window,
				OutputSchema: p.OutputSchema,
				Input:        newChild,
			}
		}

	case *LogicalParse:
		if p.Input == oldChild {
			return &LogicalParse{
				SourceField:     p.SourceField,
				Pattern:         p.Pattern,
				ExtractedFields: p.ExtractedFields,
				OutputSchema:    p.OutputSchema,
				Input:           newChild,
			}
		}
	}

	return plan
}

// LogicalDedup removes duplicate rows based on specified fields
type LogicalDedup struct {
	Fields      []ast.Expression
	Count       int  // Number of duplicates to keep (default 1)
	Consecutive bool // Only remove consecutive duplicates
	Input       LogicalPlan
}

func (l *LogicalDedup) Schema() *analyzer.Schema { return l.Input.Schema() }
func (l *LogicalDedup) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalDedup) String() string {
	fieldNames := make([]string, len(l.Fields))
	for i, f := range l.Fields {
		fieldNames[i] = f.String()
	}
	result := fmt.Sprintf("Dedup(%s", strings.Join(fieldNames, ", "))
	if l.Count > 1 {
		result += fmt.Sprintf(", count=%d", l.Count)
	}
	if l.Consecutive {
		result += ", consecutive=true"
	}
	result += ")"
	return result
}

// LogicalBin bins a field into ranges or time buckets
type LogicalBin struct {
	Field ast.Expression // Field to bin
	Span  *ast.TimeSpan
	Bins  int
	Input LogicalPlan
}

func (l *LogicalBin) Schema() *analyzer.Schema { return l.Input.Schema() }
func (l *LogicalBin) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalBin) String() string {
	result := fmt.Sprintf("Bin(%s", l.Field.String())
	if l.Span != nil {
		result += fmt.Sprintf(", span=%s", l.Span.String())
	} else if l.Bins > 0 {
		result += fmt.Sprintf(", bins=%d", l.Bins)
	}
	result += ")"
	return result
}

// LogicalTop returns the most frequent values for specified fields
type LogicalTop struct {
	Fields      []ast.Expression
	Limit       int
	GroupBy     []ast.Expression
	ShowCount   bool
	ShowPercent bool
	OutputSchema *analyzer.Schema
	Input       LogicalPlan
}

func (l *LogicalTop) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalTop) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalTop) String() string {
	fieldNames := make([]string, len(l.Fields))
	for i, f := range l.Fields {
		fieldNames[i] = f.String()
	}
	result := fmt.Sprintf("Top(%d, %s", l.Limit, strings.Join(fieldNames, ", "))
	if len(l.GroupBy) > 0 {
		groups := make([]string, len(l.GroupBy))
		for i, g := range l.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(", by=[%s]", strings.Join(groups, ", "))
	}
	result += ")"
	return result
}

// LogicalRare returns the least frequent values for specified fields
type LogicalRare struct {
	Fields      []ast.Expression
	Limit       int
	GroupBy     []ast.Expression
	ShowCount   bool
	ShowPercent bool
	OutputSchema *analyzer.Schema
	Input       LogicalPlan
}

func (l *LogicalRare) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalRare) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalRare) String() string {
	fieldNames := make([]string, len(l.Fields))
	for i, f := range l.Fields {
		fieldNames[i] = f.String()
	}
	result := fmt.Sprintf("Rare(%d, %s", l.Limit, strings.Join(fieldNames, ", "))
	if len(l.GroupBy) > 0 {
		groups := make([]string, len(l.GroupBy))
		for i, g := range l.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(", by=[%s]", strings.Join(groups, ", "))
	}
	result += ")"
	return result
}

// LogicalEval evaluates expressions and adds computed fields
type LogicalEval struct {
	Assignments  []*ast.EvalAssignment
	OutputSchema *analyzer.Schema
	Input        LogicalPlan
}

func (l *LogicalEval) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalEval) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalEval) String() string {
	assignments := make([]string, len(l.Assignments))
	for i, a := range l.Assignments {
		assignments[i] = fmt.Sprintf("%s=%s", a.Field, a.Expression.String())
	}
	return fmt.Sprintf("Eval(%s)", strings.Join(assignments, ", "))
}

// LogicalRename renames fields
type LogicalRename struct {
	Assignments  []*ast.RenameAssignment
	OutputSchema *analyzer.Schema
	Input        LogicalPlan
}

func (l *LogicalRename) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalRename) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalRename) String() string {
	assignments := make([]string, len(l.Assignments))
	for i, a := range l.Assignments {
		assignments[i] = fmt.Sprintf("%s→%s", a.OldName, a.NewName)
	}
	return fmt.Sprintf("Rename(%s)", strings.Join(assignments, ", "))
}

// LogicalReplace replaces values in a field
type LogicalReplace struct {
	Mappings []*ast.ReplaceMapping
	Field    string
	Input    LogicalPlan
}

func (l *LogicalReplace) Schema() *analyzer.Schema { return l.Input.Schema() }
func (l *LogicalReplace) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalReplace) String() string {
	mappings := make([]string, len(l.Mappings))
	for i, m := range l.Mappings {
		mappings[i] = fmt.Sprintf("%s→%s", m.OldValue.String(), m.NewValue.String())
	}
	return fmt.Sprintf("Replace(%s in %s)", strings.Join(mappings, ", "), l.Field)
}

// LogicalFillnull fills null/missing values in fields
type LogicalFillnull struct {
	Assignments  []*ast.FillnullAssignment
	DefaultValue ast.Expression
	Fields       []ast.Expression
	Input        LogicalPlan
}

func (l *LogicalFillnull) Schema() *analyzer.Schema { return l.Input.Schema() }
func (l *LogicalFillnull) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalFillnull) String() string {
	if len(l.Assignments) > 0 {
		assignments := make([]string, len(l.Assignments))
		for i, a := range l.Assignments {
			assignments[i] = fmt.Sprintf("%s=%s", a.Field, a.Value.String())
		}
		return fmt.Sprintf("Fillnull(%s)", strings.Join(assignments, ", "))
	}
	if l.DefaultValue != nil {
		result := fmt.Sprintf("Fillnull(value=%s", l.DefaultValue.String())
		if len(l.Fields) > 0 {
			fieldStrs := make([]string, len(l.Fields))
			for i, f := range l.Fields {
				fieldStrs[i] = f.String()
			}
			result += fmt.Sprintf(", fields=[%s]", strings.Join(fieldStrs, ", "))
		}
		result += ")"
		return result
	}
	return "Fillnull()"
}

// LogicalParse extracts fields from text using regex patterns
type LogicalParse struct {
	SourceField  string
	Pattern      string
	ExtractedFields []string // Field names extracted from named capture groups
	OutputSchema *analyzer.Schema
	Input        LogicalPlan
}

func (l *LogicalParse) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalParse) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalParse) String() string {
	if len(l.ExtractedFields) > 0 {
		return fmt.Sprintf("Parse(%s: %q -> [%s])",
			l.SourceField, l.Pattern, strings.Join(l.ExtractedFields, ", "))
	}
	return fmt.Sprintf("Parse(%s: %q)", l.SourceField, l.Pattern)
}

// LogicalRex extracts fields using regular expressions
type LogicalRex struct {
	SourceField     string   // Optional field to apply regex to (empty = _raw)
	Pattern         string   // Regex pattern with named captures
	ExtractedFields []string // Field names extracted from named capture groups
	OutputSchema    *analyzer.Schema
	Input           LogicalPlan
}

func (l *LogicalRex) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalRex) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalRex) String() string {
	sourceField := l.SourceField
	if sourceField == "" {
		sourceField = "_raw"
	}
	if len(l.ExtractedFields) > 0 {
		return fmt.Sprintf("Rex(%s: %q -> [%s])",
			sourceField, l.Pattern, strings.Join(l.ExtractedFields, ", "))
	}
	return fmt.Sprintf("Rex(%s: %q)", sourceField, l.Pattern)
}

// LogicalLookup enriches data with external lookup tables
type LogicalLookup struct {
	TableName      string   // Name of the lookup table
	JoinField      string   // Field from input data to join on
	JoinFieldAlias string   // Optional alias for join field
	OutputFields   []string // Fields to extract from lookup table (field names)
	OutputAliases  []string // Optional aliases for output fields
	OutputSchema   *analyzer.Schema
	Input          LogicalPlan
}

func (l *LogicalLookup) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalLookup) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalLookup) String() string {
	result := fmt.Sprintf("Lookup(table=%s, key=%s", l.TableName, l.JoinField)
	if l.JoinFieldAlias != "" {
		result += fmt.Sprintf(" AS %s", l.JoinFieldAlias)
	}
	result += " -> ["
	for i, field := range l.OutputFields {
		if i > 0 {
			result += ", "
		}
		result += field
		if i < len(l.OutputAliases) && l.OutputAliases[i] != "" {
			result += " AS " + l.OutputAliases[i]
		}
	}
	result += "])"
	return result
}

// LogicalAppend concatenates results from a subsearch
type LogicalAppend struct {
	Subsearch    LogicalPlan      // Plan for the subsearch query
	OutputSchema *analyzer.Schema // Unified schema of main query and subsearch
	Input        LogicalPlan      // Input from main query
}

func (l *LogicalAppend) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalAppend) Children() []LogicalPlan  { return []LogicalPlan{l.Input, l.Subsearch} }
func (l *LogicalAppend) String() string {
	return fmt.Sprintf("Append(subsearch=%s)", l.Subsearch.String())
}

// LogicalJoin combines datasets with SQL-like joins
type LogicalJoin struct {
	JoinType     ast.JoinType     // Type of join (inner, left, right, outer, full)
	JoinField    string           // Field to join on from left side
	RightField   string           // Field to join on from right side (same as JoinField for now)
	Right        LogicalPlan      // Right side (subsearch) plan
	OutputSchema *analyzer.Schema // Merged schema from both sides
	Input        LogicalPlan      // Left side (input) plan
}

func (l *LogicalJoin) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalJoin) Children() []LogicalPlan  { return []LogicalPlan{l.Input, l.Right} }
func (l *LogicalJoin) String() string {
	return fmt.Sprintf("Join(type=%s, field=%s, right=%s)", l.JoinType, l.JoinField, l.Right.String())
}

// LogicalTable selects specific columns for display
type LogicalTable struct {
	Fields       []ast.Expression
	OutputSchema *analyzer.Schema
	Input        LogicalPlan
}

func (l *LogicalTable) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalTable) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalTable) String() string {
	fieldNames := make([]string, len(l.Fields))
	for i, f := range l.Fields {
		fieldNames[i] = f.String()
	}
	return fmt.Sprintf("Table(%s)", strings.Join(fieldNames, ", "))
}

// LogicalEventstats computes running statistics across all events
type LogicalEventstats struct {
	GroupBy      []ast.Expression
	Aggregations []*ast.Aggregation
	OutputSchema *analyzer.Schema
	Input        LogicalPlan
}

func (l *LogicalEventstats) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalEventstats) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalEventstats) String() string {
	aggs := make([]string, len(l.Aggregations))
	for i, agg := range l.Aggregations {
		aggs[i] = agg.String()
	}

	if len(l.GroupBy) == 0 {
		return fmt.Sprintf("Eventstats(%s)", strings.Join(aggs, ", "))
	}

	groups := make([]string, len(l.GroupBy))
	for i, g := range l.GroupBy {
		groups[i] = g.String()
	}
	return fmt.Sprintf("Eventstats(%s, by=[%s])", strings.Join(aggs, ", "), strings.Join(groups, ", "))
}

// LogicalStreamstats computes running statistics in streaming fashion
type LogicalStreamstats struct {
	GroupBy      []ast.Expression
	Aggregations []*ast.Aggregation
	Window       int
	OutputSchema *analyzer.Schema
	Input        LogicalPlan
}

func (l *LogicalStreamstats) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalStreamstats) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalStreamstats) String() string {
	result := "Streamstats("
	if l.Window > 0 {
		result += fmt.Sprintf("window=%d, ", l.Window)
	}

	aggs := make([]string, len(l.Aggregations))
	for i, agg := range l.Aggregations {
		aggs[i] = agg.String()
	}
	result += strings.Join(aggs, ", ")

	if len(l.GroupBy) > 0 {
		groups := make([]string, len(l.GroupBy))
		for i, g := range l.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(", by=[%s]", strings.Join(groups, ", "))
	}
	result += ")"
	return result
}
// LogicalReverse reverses the order of rows in the result set
type LogicalReverse struct {
	Input        LogicalPlan
	OutputSchema *analyzer.Schema
}

func (l *LogicalReverse) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalReverse) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalReverse) String() string {
	return "Reverse()"
}

// LogicalFlatten flattens nested arrays/objects into separate rows
type LogicalFlatten struct {
	Input        LogicalPlan
	Field        ast.Expression   // Field to flatten
	OutputSchema *analyzer.Schema
}

func (l *LogicalFlatten) Schema() *analyzer.Schema { return l.OutputSchema }
func (l *LogicalFlatten) Children() []LogicalPlan  { return []LogicalPlan{l.Input} }
func (l *LogicalFlatten) String() string {
	return fmt.Sprintf("Flatten(%s)", l.Field.String())
}
