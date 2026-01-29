// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

import (
	"fmt"
	"strings"
)

// Command represents a PPL command in a pipe chain
type Command interface {
	Node
	commandNode()
}

// Query represents a complete PPL query
type Query struct {
	BaseNode
	Commands []Command
}

func (q *Query) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitQuery(q)
}

func (q *Query) Type() NodeType { return NodeTypeQuery }

func (q *Query) String() string {
	parts := make([]string, len(q.Commands))
	for i, cmd := range q.Commands {
		parts[i] = cmd.String()
	}
	return strings.Join(parts, " | ")
}

// SearchCommand: source=<index>
type SearchCommand struct {
	BaseNode
	Source string
}

func (s *SearchCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitSearchCommand(s)
}

func (s *SearchCommand) Type() NodeType     { return NodeTypeSearchCommand }
func (s *SearchCommand) commandNode()       {}
func (s *SearchCommand) String() string     { return fmt.Sprintf("source=%s", s.Source) }

// WhereCommand: where <condition>
type WhereCommand struct {
	BaseNode
	Condition Expression
}

func (w *WhereCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitWhereCommand(w)
}

func (w *WhereCommand) Type() NodeType     { return NodeTypeWhereCommand }
func (w *WhereCommand) commandNode()       {}
func (w *WhereCommand) String() string     { return fmt.Sprintf("where %s", w.Condition.String()) }

// FieldsCommand: fields <field1>, <field2>, ...
type FieldsCommand struct {
	BaseNode
	Fields   []Expression
	Includes bool // true for include (default), false for exclude (fields - x, y)
}

func (f *FieldsCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFieldsCommand(f)
}

func (f *FieldsCommand) Type() NodeType { return NodeTypeFieldsCommand }
func (f *FieldsCommand) commandNode()   {}
func (f *FieldsCommand) String() string {
	fields := make([]string, len(f.Fields))
	for i, field := range f.Fields {
		fields[i] = field.String()
	}
	if f.Includes {
		return fmt.Sprintf("fields %s", strings.Join(fields, ", "))
	}
	return fmt.Sprintf("fields - %s", strings.Join(fields, ", "))
}

// Aggregation represents an aggregation expression with optional alias
type Aggregation struct {
	BaseNode
	Func  *FunctionCall
	Alias string // Optional alias (as <name>)
}

func (a *Aggregation) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitAggregation(a)
}

func (a *Aggregation) Type() NodeType { return NodeTypeAggregation }
func (a *Aggregation) String() string {
	if a.Alias != "" {
		return fmt.Sprintf("%s as %s", a.Func.String(), a.Alias)
	}
	return a.Func.String()
}

// StatsCommand: stats <agg1>, <agg2> by <field1>, <field2>
type StatsCommand struct {
	BaseNode
	Aggregations []*Aggregation
	GroupBy      []Expression
}

func (s *StatsCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitStatsCommand(s)
}

func (s *StatsCommand) Type() NodeType { return NodeTypeStatsCommand }
func (s *StatsCommand) commandNode()   {}
func (s *StatsCommand) String() string {
	aggs := make([]string, len(s.Aggregations))
	for i, agg := range s.Aggregations {
		aggs[i] = agg.String()
	}
	result := fmt.Sprintf("stats %s", strings.Join(aggs, ", "))
	if len(s.GroupBy) > 0 {
		groups := make([]string, len(s.GroupBy))
		for i, g := range s.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(" by %s", strings.Join(groups, ", "))
	}
	return result
}

// SortKey represents a field to sort by with order
type SortKey struct {
	BaseNode
	Field      Expression
	Descending bool
}

func (s *SortKey) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitSortKey(s)
}

func (s *SortKey) Type() NodeType { return NodeTypeSortKey }
func (s *SortKey) String() string {
	if s.Descending {
		return fmt.Sprintf("%s DESC", s.Field.String())
	}
	return s.Field.String()
}

// SortCommand: sort <field1> [DESC], <field2> [ASC]
type SortCommand struct {
	BaseNode
	SortKeys []*SortKey
}

func (s *SortCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitSortCommand(s)
}

func (s *SortCommand) Type() NodeType { return NodeTypeSortCommand }
func (s *SortCommand) commandNode()   {}
func (s *SortCommand) String() string {
	keys := make([]string, len(s.SortKeys))
	for i, key := range s.SortKeys {
		keys[i] = key.String()
	}
	return fmt.Sprintf("sort %s", strings.Join(keys, ", "))
}

// HeadCommand: head <n>
type HeadCommand struct {
	BaseNode
	Count int
}

func (h *HeadCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitHeadCommand(h)
}

func (h *HeadCommand) Type() NodeType { return NodeTypeHeadCommand }
func (h *HeadCommand) commandNode()   {}
func (h *HeadCommand) String() string { return fmt.Sprintf("head %d", h.Count) }

// DescribeCommand: describe <source>
type DescribeCommand struct {
	BaseNode
	Source string
}

func (d *DescribeCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitDescribeCommand(d)
}

func (d *DescribeCommand) Type() NodeType     { return NodeTypeDescribeCommand }
func (d *DescribeCommand) commandNode()       {}
func (d *DescribeCommand) String() string     { return fmt.Sprintf("describe %s", d.Source) }

// ShowDatasourcesCommand: showdatasources (no args)
type ShowDatasourcesCommand struct {
	BaseNode
}

func (s *ShowDatasourcesCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitShowDatasourcesCommand(s)
}

func (s *ShowDatasourcesCommand) Type() NodeType     { return NodeTypeShowDatasourcesCommand }
func (s *ShowDatasourcesCommand) commandNode()       {}
func (s *ShowDatasourcesCommand) String() string     { return "showdatasources" }

// ExplainCommand: explain (applied to entire query)
type ExplainCommand struct {
	BaseNode
}

func (e *ExplainCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitExplainCommand(e)
}

func (e *ExplainCommand) Type() NodeType     { return NodeTypeExplainCommand }
func (e *ExplainCommand) commandNode()       {}
func (e *ExplainCommand) String() string     { return "explain" }

// ============================================================================
// Tier 1 Commands
// ============================================================================

// TimeSpan represents a time duration (e.g., "1h", "30m", "1d")
type TimeSpan struct {
	BaseNode
	Value int
	Unit  string // s, m, h, d, w, mon
}

func (ts *TimeSpan) String() string { return fmt.Sprintf("%d%s", ts.Value, ts.Unit) }

// ChartCommand: chart <aggregations> [by <fields>] [span=<timespan>]
type ChartCommand struct {
	BaseNode
	Aggregations []*Aggregation
	GroupBy      []Expression
	Span         *TimeSpan // Optional time span
	Limit        int       // Optional limit on groups
	UseOther     bool      // Include "other" bucket
	OtherStr     string    // Label for "other" bucket
	NullStr      string    // Label for null values
}

func (c *ChartCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitChartCommand(c)
}

func (c *ChartCommand) Type() NodeType { return NodeTypeChartCommand }
func (c *ChartCommand) commandNode()   {}
func (c *ChartCommand) String() string {
	aggs := make([]string, len(c.Aggregations))
	for i, agg := range c.Aggregations {
		aggs[i] = agg.String()
	}
	result := fmt.Sprintf("chart %s", strings.Join(aggs, ", "))
	if len(c.GroupBy) > 0 {
		groups := make([]string, len(c.GroupBy))
		for i, g := range c.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(" by %s", strings.Join(groups, ", "))
	}
	if c.Span != nil {
		result += fmt.Sprintf(" span=%s", c.Span.String())
	}
	return result
}

// TimechartCommand: timechart [span=<timespan>] <aggregations> [by <fields>]
type TimechartCommand struct {
	BaseNode
	Aggregations []*Aggregation
	GroupBy      []Expression
	Span         *TimeSpan // Time bucket size
	Bins         int       // Number of bins (alternative to span)
	Limit        int       // Limit on groups
	UseOther     bool      // Include "other" bucket
}

func (t *TimechartCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitTimechartCommand(t)
}

func (t *TimechartCommand) Type() NodeType { return NodeTypeTimechartCommand }
func (t *TimechartCommand) commandNode()   {}
func (t *TimechartCommand) String() string {
	result := "timechart"
	if t.Span != nil {
		result += fmt.Sprintf(" span=%s", t.Span.String())
	}
	aggs := make([]string, len(t.Aggregations))
	for i, agg := range t.Aggregations {
		aggs[i] = agg.String()
	}
	result += " " + strings.Join(aggs, ", ")
	if len(t.GroupBy) > 0 {
		groups := make([]string, len(t.GroupBy))
		for i, g := range t.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(" by %s", strings.Join(groups, ", "))
	}
	return result
}

// BinCommand: bin <field> [span=<timespan> | bins=<n>]
type BinCommand struct {
	BaseNode
	Field Expression // Field to bin
	Span  *TimeSpan  // Time span (for time fields)
	Bins  int        // Number of bins (for numeric fields)
	Auto  bool       // Auto-detect bin size
}

func (b *BinCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinCommand(b)
}

func (b *BinCommand) Type() NodeType { return NodeTypeBinCommand }
func (b *BinCommand) commandNode()   {}
func (b *BinCommand) String() string {
	result := fmt.Sprintf("bin %s", b.Field.String())
	if b.Span != nil {
		result += fmt.Sprintf(" span=%s", b.Span.String())
	} else if b.Bins > 0 {
		result += fmt.Sprintf(" bins=%d", b.Bins)
	} else if b.Auto {
		result += " span=auto"
	}
	return result
}

// DedupCommand: dedup [<n>] <fields> [keepevents=<bool>] [consecutive=<bool>]
type DedupCommand struct {
	BaseNode
	Fields      []Expression // Fields to deduplicate on
	Count       int          // Number of duplicates to keep (default 1)
	KeepEvents  bool         // Keep events even if duplicates
	Consecutive bool         // Only consider consecutive duplicates
	SortBy      []*SortKey   // Optional sort before dedup
}

func (d *DedupCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitDedupCommand(d)
}

func (d *DedupCommand) Type() NodeType { return NodeTypeDedupCommand }
func (d *DedupCommand) commandNode()   {}
func (d *DedupCommand) String() string {
	fields := make([]string, len(d.Fields))
	for i, f := range d.Fields {
		fields[i] = f.String()
	}
	result := "dedup"
	if d.Count > 1 {
		result += fmt.Sprintf(" %d", d.Count)
	}
	result += " " + strings.Join(fields, ", ")
	if d.KeepEvents {
		result += " keepevents=true"
	}
	if d.Consecutive {
		result += " consecutive=true"
	}
	return result
}

// TopCommand: top [<n>] <fields> [by <groupby>]
type TopCommand struct {
	BaseNode
	Fields       []Expression // Fields to get top values for
	Limit        int          // Number of top values (default 10)
	GroupBy      []Expression // Optional grouping
	CountField   string       // Name of count field
	PercentField string       // Name of percent field
	ShowCount    bool         // Show count column
	ShowPercent  bool         // Show percent column
	UseOther     bool         // Include "other" bucket
	OtherStr     string       // Label for "other" bucket
}

func (t *TopCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitTopCommand(t)
}

func (t *TopCommand) Type() NodeType { return NodeTypeTopCommand }
func (t *TopCommand) commandNode()   {}
func (t *TopCommand) String() string {
	fields := make([]string, len(t.Fields))
	for i, f := range t.Fields {
		fields[i] = f.String()
	}
	result := "top"
	if t.Limit > 0 {
		result += fmt.Sprintf(" %d", t.Limit)
	}
	result += " " + strings.Join(fields, ", ")
	if len(t.GroupBy) > 0 {
		groups := make([]string, len(t.GroupBy))
		for i, g := range t.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(" by %s", strings.Join(groups, ", "))
	}
	return result
}

// RareCommand: rare [<n>] <fields> [by <groupby>]
type RareCommand struct {
	BaseNode
	Fields       []Expression // Fields to get rare values for
	Limit        int          // Number of rare values (default 10)
	GroupBy      []Expression // Optional grouping
	CountField   string       // Name of count field
	PercentField string       // Name of percent field
	ShowCount    bool         // Show count column
	ShowPercent  bool         // Show percent column
}

func (r *RareCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitRareCommand(r)
}

func (r *RareCommand) Type() NodeType { return NodeTypeRareCommand }
func (r *RareCommand) commandNode()   {}
func (r *RareCommand) String() string {
	fields := make([]string, len(r.Fields))
	for i, f := range r.Fields {
		fields[i] = f.String()
	}
	result := "rare"
	if r.Limit > 0 {
		result += fmt.Sprintf(" %d", r.Limit)
	}
	result += " " + strings.Join(fields, ", ")
	if len(r.GroupBy) > 0 {
		groups := make([]string, len(r.GroupBy))
		for i, g := range r.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(" by %s", strings.Join(groups, ", "))
	}
	return result
}

// EvalAssignment represents a single eval assignment: field = expression
type EvalAssignment struct {
	BaseNode
	Field      string     // Target field name
	Expression Expression // Value expression
}

func (e *EvalAssignment) String() string {
	return fmt.Sprintf("%s = %s", e.Field, e.Expression.String())
}

// EvalCommand: eval <field> = <expression>, ...
type EvalCommand struct {
	BaseNode
	Assignments []*EvalAssignment
}

func (e *EvalCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitEvalCommand(e)
}

func (e *EvalCommand) Type() NodeType { return NodeTypeEvalCommand }
func (e *EvalCommand) commandNode()   {}
func (e *EvalCommand) String() string {
	assignments := make([]string, len(e.Assignments))
	for i, a := range e.Assignments {
		assignments[i] = a.String()
	}
	return fmt.Sprintf("eval %s", strings.Join(assignments, ", "))
}

// RenameAssignment represents a single rename: old as new
type RenameAssignment struct {
	BaseNode
	OldName string
	NewName string
}

func (r *RenameAssignment) String() string {
	return fmt.Sprintf("%s as %s", r.OldName, r.NewName)
}

// RenameCommand: rename <old> as <new>, ...
type RenameCommand struct {
	BaseNode
	Assignments []*RenameAssignment
}

func (r *RenameCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitRenameCommand(r)
}

func (r *RenameCommand) Type() NodeType { return NodeTypeRenameCommand }
func (r *RenameCommand) commandNode()   {}
func (r *RenameCommand) String() string {
	assignments := make([]string, len(r.Assignments))
	for i, a := range r.Assignments {
		assignments[i] = a.String()
	}
	return fmt.Sprintf("rename %s", strings.Join(assignments, ", "))
}

// ReplaceMapping represents a single value replacement: old with new
type ReplaceMapping struct {
	BaseNode
	OldValue Expression // Value to replace (can be string literal or pattern)
	NewValue Expression // Replacement value
}

func (r *ReplaceMapping) String() string {
	return fmt.Sprintf("%s with %s", r.OldValue.String(), r.NewValue.String())
}

// ReplaceCommand: replace <oldval1> with <newval1>, <oldval2> with <newval2> in <field>
type ReplaceCommand struct {
	BaseNode
	Mappings []* ReplaceMapping
	Field    string // Target field to apply replacements
}

func (r *ReplaceCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitReplaceCommand(r)
}

func (r *ReplaceCommand) Type() NodeType { return NodeTypeReplaceCommand }
func (r *ReplaceCommand) commandNode()   {}
func (r *ReplaceCommand) String() string {
	mappings := make([]string, len(r.Mappings))
	for i, m := range r.Mappings {
		mappings[i] = m.String()
	}
	return fmt.Sprintf("replace %s in %s", strings.Join(mappings, ", "), r.Field)
}

// FillnullAssignment represents filling null values in a field: field=value
type FillnullAssignment struct {
	BaseNode
	Field string
	Value Expression
}

func (f *FillnullAssignment) String() string {
	return fmt.Sprintf("%s=%s", f.Field, f.Value.String())
}

// FillnullCommand: fillnull <field1>=<value1>, <field2>=<value2>, ... or fillnull value=<default> fields <field1>, <field2>
type FillnullCommand struct {
	BaseNode
	Assignments  []*FillnullAssignment // Per-field assignments
	DefaultValue Expression            // Default value for all fields
	Fields       []Expression          // Specific fields to fill (when using default value)
}

func (f *FillnullCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFillnullCommand(f)
}

func (f *FillnullCommand) Type() NodeType { return NodeTypeFillnullCommand }
func (f *FillnullCommand) commandNode()   {}
func (f *FillnullCommand) String() string {
	if len(f.Assignments) > 0 {
		assignments := make([]string, len(f.Assignments))
		for i, a := range f.Assignments {
			assignments[i] = a.String()
		}
		return fmt.Sprintf("fillnull %s", strings.Join(assignments, ", "))
	}
	if f.DefaultValue != nil {
		result := fmt.Sprintf("fillnull value=%s", f.DefaultValue.String())
		if len(f.Fields) > 0 {
			fields := make([]string, len(f.Fields))
			for i, field := range f.Fields {
				fields[i] = field.String()
			}
			result += fmt.Sprintf(" fields %s", strings.Join(fields, ", "))
		}
		return result
	}
	return "fillnull"
}

// ParseCommand: parse [field=]<source_field> "<pattern>"
// Extracts fields from text using regex patterns with named captures
// Example: parse message "(?<user>\w+) logged in from (?<ip>\d+\.\d+\.\d+\.\d+)"
type ParseCommand struct {
	BaseNode
	SourceField string // Field to parse (e.g., "message", "_raw")
	Pattern     string // Regex pattern with named captures
	FieldParam  string // Optional field parameter name (usually "field")
}

func (p *ParseCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitParseCommand(p)
}

func (p *ParseCommand) Type() NodeType { return NodeTypeParseCommand }
func (p *ParseCommand) commandNode()   {}
func (p *ParseCommand) String() string {
	if p.FieldParam != "" {
		return fmt.Sprintf("parse %s=%s \"%s\"", p.FieldParam, p.SourceField, p.Pattern)
	}
	return fmt.Sprintf("parse %s \"%s\"", p.SourceField, p.Pattern)
}

// RexCommand: rex [field=<source_field>] "<pattern>"
// Extracts fields using regular expressions
// Example: rex "(?<error_code>\d{3}): (?<error_msg>.*)"
// Example: rex field=message "user=(?<user>\w+)"
type RexCommand struct {
	BaseNode
	SourceField string // Field to apply regex to (optional, defaults to _raw)
	Pattern     string // Regex pattern with named captures
	FieldParam  string // Optional field parameter name (usually "field")
}

func (r *RexCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitRexCommand(r)
}

func (r *RexCommand) Type() NodeType { return NodeTypeRexCommand }
func (r *RexCommand) commandNode()   {}
func (r *RexCommand) String() string {
	if r.SourceField != "" {
		if r.FieldParam != "" {
			return fmt.Sprintf("rex %s=%s \"%s\"", r.FieldParam, r.SourceField, r.Pattern)
		}
		return fmt.Sprintf("rex field=%s \"%s\"", r.SourceField, r.Pattern)
	}
	return fmt.Sprintf("rex \"%s\"", r.Pattern)
}

// LookupOutputField represents an output field from a lookup table
type LookupOutputField struct {
	Field string // Field name from lookup table
	Alias string // Optional alias for the field in output
}

func (l *LookupOutputField) String() string {
	if l.Alias != "" {
		return fmt.Sprintf("%s AS %s", l.Field, l.Alias)
	}
	return l.Field
}

// LookupCommand: lookup <table> <join_field> [AS <alias>] OUTPUT <fields>
// Enriches data with external lookup tables
// Example: lookup products product_id OUTPUT name, price
// Example: lookup users user_id AS uid OUTPUT username AS user
type LookupCommand struct {
	BaseNode
	TableName     string               // Name of the lookup table
	JoinField     string               // Field from input data to join on
	JoinFieldAlias string              // Optional alias for join field
	OutputFields  []*LookupOutputField // Fields to extract from lookup table
}

func (l *LookupCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLookupCommand(l)
}

func (l *LookupCommand) Type() NodeType { return NodeTypeLookupCommand }
func (l *LookupCommand) commandNode()   {}
func (l *LookupCommand) String() string {
	result := fmt.Sprintf("lookup %s %s", l.TableName, l.JoinField)
	if l.JoinFieldAlias != "" {
		result += fmt.Sprintf(" AS %s", l.JoinFieldAlias)
	}
	result += " OUTPUT "
	outputStrs := make([]string, len(l.OutputFields))
	for i, field := range l.OutputFields {
		outputStrs[i] = field.String()
	}
	result += strings.Join(outputStrs, ", ")
	return result
}

// AppendCommand: append [subsearch]
// Concatenates results from a subsearch to the main search results
type AppendCommand struct {
	BaseNode
	Subsearch *Query // The subsearch query to append
}

func (a *AppendCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitAppendCommand(a)
}

func (a *AppendCommand) Type() NodeType { return NodeTypeAppendCommand }
func (a *AppendCommand) commandNode()   {}
func (a *AppendCommand) String() string {
	return fmt.Sprintf("append [%s]", a.Subsearch.String())
}

// JoinType represents the type of join operation
type JoinType string

const (
	JoinTypeInner JoinType = "inner"
	JoinTypeLeft  JoinType = "left"
	JoinTypeRight JoinType = "right"
	JoinTypeOuter JoinType = "outer"
	JoinTypeFull  JoinType = "full"
)

// JoinCommand: join [type=TYPE] field [subsearch]
// Combines datasets with SQL-like joins
type JoinCommand struct {
	BaseNode
	JoinType   JoinType // Type of join (inner, left, right, outer, full)
	JoinField  string   // Field to join on (from both sides)
	Subsearch  *Query   // The right side query
}

func (j *JoinCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitJoinCommand(j)
}

func (j *JoinCommand) Type() NodeType { return NodeTypeJoinCommand }
func (j *JoinCommand) commandNode()   {}
func (j *JoinCommand) String() string {
	result := fmt.Sprintf("join")
	if j.JoinType != JoinTypeInner {
		result += fmt.Sprintf(" type=%s", j.JoinType)
	}
	result += fmt.Sprintf(" %s [%s]", j.JoinField, j.Subsearch.String())
	return result
}

// TableCommand: table <field1>, <field2>, ...
type TableCommand struct {
	BaseNode
	Fields []Expression
}

func (t *TableCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitTableCommand(t)
}

func (t *TableCommand) Type() NodeType { return NodeTypeFieldsCommand } // Reuse FieldsCommand type
func (t *TableCommand) commandNode()   {}
func (t *TableCommand) String() string {
	fields := make([]string, len(t.Fields))
	for i, f := range t.Fields {
		fields[i] = f.String()
	}
	return fmt.Sprintf("table %s", strings.Join(fields, ", "))
}

// EventstatsCommand: eventstats <aggregations> [by <fields>]
type EventstatsCommand struct {
	BaseNode
	Aggregations []*Aggregation
	GroupBy      []Expression
}

func (e *EventstatsCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitEventstatsCommand(e)
}

func (e *EventstatsCommand) Type() NodeType { return NodeTypeStatsCommand } // Similar to stats
func (e *EventstatsCommand) commandNode()   {}
func (e *EventstatsCommand) String() string {
	aggs := make([]string, len(e.Aggregations))
	for i, agg := range e.Aggregations {
		aggs[i] = agg.String()
	}
	result := fmt.Sprintf("eventstats %s", strings.Join(aggs, ", "))
	if len(e.GroupBy) > 0 {
		groups := make([]string, len(e.GroupBy))
		for i, g := range e.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(" by %s", strings.Join(groups, ", "))
	}
	return result
}

// StreamstatsCommand: streamstats [options] <aggregations> [by <fields>]
type StreamstatsCommand struct {
	BaseNode
	Aggregations []*Aggregation
	GroupBy      []Expression
	Window       int  // Window size for rolling aggregations
	Current      bool // Include current event
	Global       bool // Global vs per-group
}

func (s *StreamstatsCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitStreamstatsCommand(s)
}

func (s *StreamstatsCommand) Type() NodeType { return NodeTypeStatsCommand } // Similar to stats
func (s *StreamstatsCommand) commandNode()   {}
func (s *StreamstatsCommand) String() string {
	result := "streamstats"
	if s.Window > 0 {
		result += fmt.Sprintf(" window=%d", s.Window)
	}
	aggs := make([]string, len(s.Aggregations))
	for i, agg := range s.Aggregations {
		aggs[i] = agg.String()
	}
	result += " " + strings.Join(aggs, ", ")
	if len(s.GroupBy) > 0 {
		groups := make([]string, len(s.GroupBy))
		for i, g := range s.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(" by %s", strings.Join(groups, ", "))
	}
	return result
}

// ReverseCommand: reverse - Reverses the order of rows in the result set
type ReverseCommand struct {
	BaseNode
}

func (r *ReverseCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitReverseCommand(r)
}

func (r *ReverseCommand) Type() NodeType { return NodeTypeReverseCommand }
func (r *ReverseCommand) commandNode()   {}
func (r *ReverseCommand) String() string {
	return "reverse"
}

// FlattenCommand: flatten <field> - Flattens nested arrays/objects into separate rows
// Takes a field containing an array or nested object and creates multiple rows,
// one for each array element or nested value
type FlattenCommand struct {
	BaseNode
	Field Expression // Field to flatten (can be nested like "data.items")
}

func (f *FlattenCommand) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitFlattenCommand(f)
}

func (f *FlattenCommand) Type() NodeType { return NodeTypeFlattenCommand }
func (f *FlattenCommand) commandNode()   {}
func (f *FlattenCommand) String() string {
	return fmt.Sprintf("flatten %s", f.Field.String())
}
