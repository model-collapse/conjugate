// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package physical

import (
	"fmt"
	"strings"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
)

// ExecutionLocation indicates where an operation executes
type ExecutionLocation int

const (
	// ExecuteOnDataNode means the operation is pushed down to OpenSearch
	ExecuteOnDataNode ExecutionLocation = iota
	// ExecuteOnCoordinator means the operation runs on the coordinator node
	ExecuteOnCoordinator
)

func (e ExecutionLocation) String() string {
	switch e {
	case ExecuteOnDataNode:
		return "DataNode"
	case ExecuteOnCoordinator:
		return "Coordinator"
	default:
		return "Unknown"
	}
}

// PhysicalPlan represents a physical query execution plan
type PhysicalPlan interface {
	// Schema returns the output schema
	Schema() *analyzer.Schema

	// Children returns child operators
	Children() []PhysicalPlan

	// Location returns where this operator executes
	Location() ExecutionLocation

	// String returns a string representation
	String() string
}

// PhysicalScan reads data from a source
type PhysicalScan struct {
	Source         string
	OutputSchema   *analyzer.Schema
	Filter         ast.Expression         // Optional pushed-down filter
	Fields         []string               // Optional pushed-down projection
	SortKeys       []*ast.SortKey         // Optional pushed-down sort
	Limit          int                    // Optional pushed-down limit (0 = no limit)
	ComputedFields []*ast.EvalAssignment  // Optional pushed-down computed fields (eval)
}

func (p *PhysicalScan) Schema() *analyzer.Schema        { return p.OutputSchema }
func (p *PhysicalScan) Children() []PhysicalPlan        { return nil }
func (p *PhysicalScan) Location() ExecutionLocation     { return ExecuteOnDataNode }
func (p *PhysicalScan) String() string {
	parts := []string{fmt.Sprintf("PhysicalScan(%s)", p.Source)}

	if p.Filter != nil {
		parts = append(parts, fmt.Sprintf("filter=%s", p.Filter.String()))
	}
	if len(p.Fields) > 0 {
		parts = append(parts, fmt.Sprintf("fields=%v", p.Fields))
	}
	if len(p.SortKeys) > 0 {
		keys := make([]string, len(p.SortKeys))
		for i, k := range p.SortKeys {
			order := "ASC"
			if k.Descending {
				order = "DESC"
			}
			keys[i] = fmt.Sprintf("%s %s", k.Field.String(), order)
		}
		parts = append(parts, fmt.Sprintf("sort=[%s]", strings.Join(keys, ", ")))
	}
	if p.Limit > 0 {
		parts = append(parts, fmt.Sprintf("limit=%d", p.Limit))
	}
	if len(p.ComputedFields) > 0 {
		computed := make([]string, len(p.ComputedFields))
		for i, cf := range p.ComputedFields {
			computed[i] = fmt.Sprintf("%s=%s", cf.Field, cf.Expression.String())
		}
		parts = append(parts, fmt.Sprintf("computed=[%s]", strings.Join(computed, ", ")))
	}

	return strings.Join(parts, ", ")
}

// PhysicalFilter filters rows on the coordinator
type PhysicalFilter struct {
	Condition ast.Expression
	Input     PhysicalPlan
}

func (p *PhysicalFilter) Schema() *analyzer.Schema    { return p.Input.Schema() }
func (p *PhysicalFilter) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalFilter) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalFilter) String() string {
	return fmt.Sprintf("PhysicalFilter(%s)", p.Condition.String())
}

// PhysicalProject selects fields on the coordinator
type PhysicalProject struct {
	Fields       []ast.Expression
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
	Exclude      bool
}

func (p *PhysicalProject) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalProject) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalProject) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalProject) String() string {
	fieldNames := make([]string, len(p.Fields))
	for i, f := range p.Fields {
		fieldNames[i] = f.String()
	}
	prefix := ""
	if p.Exclude {
		prefix = "-"
	}
	return fmt.Sprintf("PhysicalProject(%s%s)", prefix, strings.Join(fieldNames, ", "))
}

// PhysicalSort sorts rows on the coordinator
type PhysicalSort struct {
	SortKeys []*ast.SortKey
	Input    PhysicalPlan
}

func (p *PhysicalSort) Schema() *analyzer.Schema    { return p.Input.Schema() }
func (p *PhysicalSort) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalSort) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalSort) String() string {
	keys := make([]string, len(p.SortKeys))
	for i, k := range p.SortKeys {
		order := "ASC"
		if k.Descending {
			order = "DESC"
		}
		keys[i] = fmt.Sprintf("%s %s", k.Field.String(), order)
	}
	return fmt.Sprintf("PhysicalSort(%s)", strings.Join(keys, ", "))
}

// PhysicalLimit limits rows on the coordinator
type PhysicalLimit struct {
	Count int
	Input PhysicalPlan
}

func (p *PhysicalLimit) Schema() *analyzer.Schema    { return p.Input.Schema() }
func (p *PhysicalLimit) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalLimit) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalLimit) String() string {
	return fmt.Sprintf("PhysicalLimit(%d)", p.Count)
}

// AggregationAlgorithm specifies how aggregation is performed
type AggregationAlgorithm int

const (
	// HashAggregation uses a hash table (for unsorted input, high cardinality)
	HashAggregation AggregationAlgorithm = iota
	// StreamAggregation uses streaming (for sorted input, low cardinality)
	StreamAggregation
)

func (a AggregationAlgorithm) String() string {
	switch a {
	case HashAggregation:
		return "Hash"
	case StreamAggregation:
		return "Stream"
	default:
		return "Unknown"
	}
}

// PhysicalAggregate performs aggregation on the coordinator
type PhysicalAggregate struct {
	GroupBy      []ast.Expression
	Aggregations []*ast.Aggregation
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
	Algorithm    AggregationAlgorithm
}

func (p *PhysicalAggregate) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalAggregate) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalAggregate) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalAggregate) String() string {
	aggs := make([]string, len(p.Aggregations))
	for i, agg := range p.Aggregations {
		aggs[i] = agg.String()
	}

	if len(p.GroupBy) == 0 {
		return fmt.Sprintf("PhysicalAggregate[%s](%s)", p.Algorithm, strings.Join(aggs, ", "))
	}

	groups := make([]string, len(p.GroupBy))
	for i, g := range p.GroupBy {
		groups[i] = g.String()
	}
	return fmt.Sprintf("PhysicalAggregate[%s](%s, group_by=[%s])",
		p.Algorithm, strings.Join(aggs, ", "), strings.Join(groups, ", "))
}

// PrintPlan prints a physical plan tree
func PrintPlan(plan PhysicalPlan, indent int) string {
	prefix := strings.Repeat("  ", indent)
	result := prefix + plan.String()

	// Add location annotation
	result += fmt.Sprintf(" [%s]", plan.Location())
	result += "\n"

	for _, child := range plan.Children() {
		result += PrintPlan(child, indent+1)
	}

	return result
}

// GetLeafScans returns all PhysicalScan nodes in the plan
func GetLeafScans(plan PhysicalPlan) []*PhysicalScan {
	scans := make([]*PhysicalScan, 0)

	if scan, ok := plan.(*PhysicalScan); ok {
		scans = append(scans, scan)
		return scans
	}

	for _, child := range plan.Children() {
		scans = append(scans, GetLeafScans(child)...)
	}

	return scans
}

// IsPushedDown returns true if the plan has operations pushed down to data node
func IsPushedDown(plan PhysicalPlan) bool {
	scans := GetLeafScans(plan)
	for _, scan := range scans {
		if scan.Filter != nil || len(scan.Fields) > 0 ||
		   len(scan.SortKeys) > 0 || scan.Limit > 0 || len(scan.ComputedFields) > 0 {
			return true
		}
	}
	return false
}

// CountCoordinatorOps counts the number of coordinator-side operations
func CountCoordinatorOps(plan PhysicalPlan) int {
	count := 0
	if plan.Location() == ExecuteOnCoordinator {
		count++
	}

	for _, child := range plan.Children() {
		count += CountCoordinatorOps(child)
	}

	return count
}

// PhysicalDedup removes duplicate rows
type PhysicalDedup struct {
	Fields      []ast.Expression
	Count       int  // Number of duplicates to keep
	Consecutive bool // Only remove consecutive duplicates
	Input       PhysicalPlan
}

func (p *PhysicalDedup) Schema() *analyzer.Schema    { return p.Input.Schema() }
func (p *PhysicalDedup) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalDedup) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalDedup) String() string {
	fieldNames := make([]string, len(p.Fields))
	for i, f := range p.Fields {
		fieldNames[i] = f.String()
	}
	result := fmt.Sprintf("PhysicalDedup(%s", strings.Join(fieldNames, ", "))
	if p.Count > 1 {
		result += fmt.Sprintf(", count=%d", p.Count)
	}
	if p.Consecutive {
		result += ", consecutive"
	}
	result += ")"
	return result
}

// PhysicalBin bins values into buckets
type PhysicalBin struct {
	Field        ast.Expression
	Span         *ast.TimeSpan
	Bins         int
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
}

func (p *PhysicalBin) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalBin) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalBin) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalBin) String() string {
	result := fmt.Sprintf("PhysicalBin(%s", p.Field.String())
	if p.Span != nil {
		result += fmt.Sprintf(", span=%s", p.Span.String())
	} else if p.Bins > 0 {
		result += fmt.Sprintf(", bins=%d", p.Bins)
	}
	result += ")"
	return result
}

// TopRareAlgorithm specifies the algorithm for top/rare computation
type TopRareAlgorithm int

const (
	// TopRareHash uses hash-based counting
	TopRareHash TopRareAlgorithm = iota
	// TopRareHeap uses a min/max heap for streaming
	TopRareHeap
)

func (a TopRareAlgorithm) String() string {
	switch a {
	case TopRareHash:
		return "Hash"
	case TopRareHeap:
		return "Heap"
	default:
		return "Unknown"
	}
}

// PhysicalTop returns most frequent values
type PhysicalTop struct {
	Fields       []ast.Expression
	Limit        int
	GroupBy      []ast.Expression
	ShowCount    bool
	ShowPercent  bool
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
	Algorithm    TopRareAlgorithm
}

func (p *PhysicalTop) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalTop) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalTop) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalTop) String() string {
	fieldNames := make([]string, len(p.Fields))
	for i, f := range p.Fields {
		fieldNames[i] = f.String()
	}
	result := fmt.Sprintf("PhysicalTop[%s](%d, %s", p.Algorithm, p.Limit, strings.Join(fieldNames, ", "))
	if len(p.GroupBy) > 0 {
		groups := make([]string, len(p.GroupBy))
		for i, g := range p.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(", by=[%s]", strings.Join(groups, ", "))
	}
	result += ")"
	return result
}

// PhysicalRare returns least frequent values
type PhysicalRare struct {
	Fields       []ast.Expression
	Limit        int
	GroupBy      []ast.Expression
	ShowCount    bool
	ShowPercent  bool
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
	Algorithm    TopRareAlgorithm
}

func (p *PhysicalRare) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalRare) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalRare) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalRare) String() string {
	fieldNames := make([]string, len(p.Fields))
	for i, f := range p.Fields {
		fieldNames[i] = f.String()
	}
	result := fmt.Sprintf("PhysicalRare[%s](%d, %s", p.Algorithm, p.Limit, strings.Join(fieldNames, ", "))
	if len(p.GroupBy) > 0 {
		groups := make([]string, len(p.GroupBy))
		for i, g := range p.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(", by=[%s]", strings.Join(groups, ", "))
	}
	result += ")"
	return result
}

// PhysicalEval computes new fields from expressions
type PhysicalEval struct {
	Assignments  []*ast.EvalAssignment
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
}

func (p *PhysicalEval) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalEval) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalEval) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalEval) String() string {
	assignments := make([]string, len(p.Assignments))
	for i, a := range p.Assignments {
		assignments[i] = fmt.Sprintf("%s=%s", a.Field, a.Expression.String())
	}
	return fmt.Sprintf("PhysicalEval(%s)", strings.Join(assignments, ", "))
}

// PhysicalRename renames fields
type PhysicalRename struct {
	Assignments  []*ast.RenameAssignment
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
}

func (p *PhysicalRename) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalRename) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalRename) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalRename) String() string {
	assignments := make([]string, len(p.Assignments))
	for i, a := range p.Assignments {
		assignments[i] = fmt.Sprintf("%s→%s", a.OldName, a.NewName)
	}
	return fmt.Sprintf("PhysicalRename(%s)", strings.Join(assignments, ", "))
}

// PhysicalReplace replaces values in a field
type PhysicalReplace struct {
	Mappings []*ast.ReplaceMapping
	Field    string
	Input    PhysicalPlan
}

func (p *PhysicalReplace) Schema() *analyzer.Schema    { return p.Input.Schema() }
func (p *PhysicalReplace) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalReplace) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalReplace) String() string {
	mappings := make([]string, len(p.Mappings))
	for i, m := range p.Mappings {
		mappings[i] = fmt.Sprintf("%s→%s", m.OldValue.String(), m.NewValue.String())
	}
	return fmt.Sprintf("PhysicalReplace(%s in %s)", strings.Join(mappings, ", "), p.Field)
}

// PhysicalFillnull fills null/missing values in fields
type PhysicalFillnull struct {
	Assignments  []*ast.FillnullAssignment
	DefaultValue ast.Expression
	Fields       []ast.Expression
	Input        PhysicalPlan
}

func (p *PhysicalFillnull) Schema() *analyzer.Schema    { return p.Input.Schema() }
func (p *PhysicalFillnull) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalFillnull) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalFillnull) String() string {
	if len(p.Assignments) > 0 {
		assignments := make([]string, len(p.Assignments))
		for i, a := range p.Assignments {
			assignments[i] = fmt.Sprintf("%s=%s", a.Field, a.Value.String())
		}
		return fmt.Sprintf("PhysicalFillnull(%s)", strings.Join(assignments, ", "))
	}
	if p.DefaultValue != nil {
		result := fmt.Sprintf("PhysicalFillnull(value=%s", p.DefaultValue.String())
		if len(p.Fields) > 0 {
			fieldStrs := make([]string, len(p.Fields))
			for i, f := range p.Fields {
				fieldStrs[i] = f.String()
			}
			result += fmt.Sprintf(", fields=[%s]", strings.Join(fieldStrs, ", "))
		}
		result += ")"
		return result
	}
	return "PhysicalFillnull()"
}

// PhysicalTable selects specific columns for display
type PhysicalTable struct {
	Fields       []ast.Expression
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
}

func (p *PhysicalTable) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalTable) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalTable) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalTable) String() string {
	fieldNames := make([]string, len(p.Fields))
	for i, f := range p.Fields {
		fieldNames[i] = f.String()
	}
	return fmt.Sprintf("PhysicalTable(%s)", strings.Join(fieldNames, ", "))
}

// PhysicalEventstats computes running statistics across all events
type PhysicalEventstats struct {
	GroupBy      []ast.Expression
	Aggregations []*ast.Aggregation
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
}

func (p *PhysicalEventstats) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalEventstats) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalEventstats) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalEventstats) String() string {
	aggs := make([]string, len(p.Aggregations))
	for i, agg := range p.Aggregations {
		aggs[i] = agg.String()
	}

	if len(p.GroupBy) == 0 {
		return fmt.Sprintf("PhysicalEventstats(%s)", strings.Join(aggs, ", "))
	}

	groups := make([]string, len(p.GroupBy))
	for i, g := range p.GroupBy {
		groups[i] = g.String()
	}
	return fmt.Sprintf("PhysicalEventstats(%s, by=[%s])", strings.Join(aggs, ", "), strings.Join(groups, ", "))
}

// PhysicalStreamstats computes running statistics in streaming fashion
type PhysicalStreamstats struct {
	GroupBy      []ast.Expression
	Aggregations []*ast.Aggregation
	Window       int
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
}

func (p *PhysicalStreamstats) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalStreamstats) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalStreamstats) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalStreamstats) String() string {
	result := "PhysicalStreamstats("
	if p.Window > 0 {
		result += fmt.Sprintf("window=%d, ", p.Window)
	}

	aggs := make([]string, len(p.Aggregations))
	for i, agg := range p.Aggregations {
		aggs[i] = agg.String()
	}
	result += strings.Join(aggs, ", ")

	if len(p.GroupBy) > 0 {
		groups := make([]string, len(p.GroupBy))
		for i, g := range p.GroupBy {
			groups[i] = g.String()
		}
		result += fmt.Sprintf(", by=[%s]", strings.Join(groups, ", "))
	}
	result += ")"
	return result
}

// PhysicalParse extracts fields from text using regex patterns
type PhysicalParse struct {
	SourceField     string
	Pattern         string
	ExtractedFields []string
	OutputSchema    *analyzer.Schema
	Input           PhysicalPlan
}

func (p *PhysicalParse) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalParse) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalParse) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalParse) String() string {
	if len(p.ExtractedFields) > 0 {
		return fmt.Sprintf("PhysicalParse(%s: %q -> [%s])",
			p.SourceField, p.Pattern, strings.Join(p.ExtractedFields, ", "))
	}
	return fmt.Sprintf("PhysicalParse(%s: %q)", p.SourceField, p.Pattern)
}

// PhysicalRex extracts fields using regular expressions
type PhysicalRex struct {
	SourceField     string   // Field to apply regex to (empty = _raw)
	Pattern         string   // Regex pattern with named captures
	ExtractedFields []string // Field names extracted
	OutputSchema    *analyzer.Schema
	Input           PhysicalPlan
}

func (p *PhysicalRex) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalRex) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalRex) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalRex) String() string {
	sourceField := p.SourceField
	if sourceField == "" {
		sourceField = "_raw"
	}
	if len(p.ExtractedFields) > 0 {
		return fmt.Sprintf("PhysicalRex(%s: %q -> [%s])",
			sourceField, p.Pattern, strings.Join(p.ExtractedFields, ", "))
	}
	return fmt.Sprintf("PhysicalRex(%s: %q)", sourceField, p.Pattern)
}

// PhysicalLookup enriches data with external lookup tables
type PhysicalLookup struct {
	TableName      string   // Name of the lookup table
	JoinField      string   // Field from input data to join on
	JoinFieldAlias string   // Optional alias for join field
	OutputFields   []string // Fields to extract from lookup table
	OutputAliases  []string // Optional aliases for output fields
	OutputSchema   *analyzer.Schema
	Input          PhysicalPlan
}

func (p *PhysicalLookup) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalLookup) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input} }
func (p *PhysicalLookup) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalLookup) String() string {
	result := fmt.Sprintf("PhysicalLookup(table=%s, key=%s", p.TableName, p.JoinField)
	if p.JoinFieldAlias != "" {
		result += fmt.Sprintf(" AS %s", p.JoinFieldAlias)
	}
	result += " -> ["
	for i, field := range p.OutputFields {
		if i > 0 {
			result += ", "
		}
		result += field
		if i < len(p.OutputAliases) && p.OutputAliases[i] != "" {
			result += " AS " + p.OutputAliases[i]
		}
	}
	result += "])"
	return result
}

// PhysicalAppend concatenates results from a subsearch
type PhysicalAppend struct {
	Subsearch    PhysicalPlan // Physical plan for the subsearch
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan
}

func (p *PhysicalAppend) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalAppend) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input, p.Subsearch} }
func (p *PhysicalAppend) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalAppend) String() string {
	return fmt.Sprintf("PhysicalAppend(subsearch=%s)", p.Subsearch.String())
}

// PhysicalJoin combines datasets with SQL-like joins using hash join algorithm
type PhysicalJoin struct {
	JoinType     ast.JoinType // Type of join (inner, left, right, outer, full)
	JoinField    string       // Field to join on from left side
	RightField   string       // Field to join on from right side
	Right        PhysicalPlan // Right side (build side) physical plan
	OutputSchema *analyzer.Schema
	Input        PhysicalPlan // Left side (probe side) physical plan
}

func (p *PhysicalJoin) Schema() *analyzer.Schema    { return p.OutputSchema }
func (p *PhysicalJoin) Children() []PhysicalPlan    { return []PhysicalPlan{p.Input, p.Right} }
func (p *PhysicalJoin) Location() ExecutionLocation { return ExecuteOnCoordinator }
func (p *PhysicalJoin) String() string {
	return fmt.Sprintf("PhysicalJoin(type=%s, field=%s, right=%s)",
		p.JoinType, p.JoinField, p.Right.String())
}

// PhysicalReverse reverses the order of rows
// Must execute on coordinator since it needs to buffer all results
type PhysicalReverse struct {
	Input        PhysicalPlan
	OutputSchema *analyzer.Schema
}

func (p *PhysicalReverse) Schema() *analyzer.Schema         { return p.OutputSchema }
func (p *PhysicalReverse) Children() []PhysicalPlan         { return []PhysicalPlan{p.Input} }
func (p *PhysicalReverse) Location() ExecutionLocation      { return ExecuteOnCoordinator }
func (p *PhysicalReverse) String() string {
	return "Reverse()"
}

// PhysicalFlatten flattens nested arrays/objects into separate rows
// Must execute on coordinator since it generates multiple rows per input row
type PhysicalFlatten struct {
	Input        PhysicalPlan
	Field        ast.Expression
	OutputSchema *analyzer.Schema
}

func (p *PhysicalFlatten) Schema() *analyzer.Schema         { return p.OutputSchema }
func (p *PhysicalFlatten) Children() []PhysicalPlan         { return []PhysicalPlan{p.Input} }
func (p *PhysicalFlatten) Location() ExecutionLocation      { return ExecuteOnCoordinator }
func (p *PhysicalFlatten) String() string {
	return fmt.Sprintf("Flatten(%s)", p.Field.String())
}
