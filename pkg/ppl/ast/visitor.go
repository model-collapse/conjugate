// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

// Visitor defines the visitor pattern for traversing the AST
type Visitor interface {
	// Query and commands (Tier 0)
	VisitQuery(*Query) (interface{}, error)
	VisitSearchCommand(*SearchCommand) (interface{}, error)
	VisitWhereCommand(*WhereCommand) (interface{}, error)
	VisitFieldsCommand(*FieldsCommand) (interface{}, error)
	VisitStatsCommand(*StatsCommand) (interface{}, error)
	VisitSortCommand(*SortCommand) (interface{}, error)
	VisitHeadCommand(*HeadCommand) (interface{}, error)
	VisitDescribeCommand(*DescribeCommand) (interface{}, error)
	VisitShowDatasourcesCommand(*ShowDatasourcesCommand) (interface{}, error)
	VisitExplainCommand(*ExplainCommand) (interface{}, error)

	// Tier 1 commands
	VisitChartCommand(*ChartCommand) (interface{}, error)
	VisitTimechartCommand(*TimechartCommand) (interface{}, error)
	VisitBinCommand(*BinCommand) (interface{}, error)
	VisitDedupCommand(*DedupCommand) (interface{}, error)
	VisitTopCommand(*TopCommand) (interface{}, error)
	VisitRareCommand(*RareCommand) (interface{}, error)
	VisitEvalCommand(*EvalCommand) (interface{}, error)
	VisitRenameCommand(*RenameCommand) (interface{}, error)
	VisitReplaceCommand(*ReplaceCommand) (interface{}, error)
	VisitFillnullCommand(*FillnullCommand) (interface{}, error)
	VisitParseCommand(*ParseCommand) (interface{}, error)
	VisitRexCommand(*RexCommand) (interface{}, error)
	VisitLookupCommand(*LookupCommand) (interface{}, error)
	VisitAppendCommand(*AppendCommand) (interface{}, error)
	VisitJoinCommand(*JoinCommand) (interface{}, error)
	VisitTableCommand(*TableCommand) (interface{}, error)
	VisitEventstatsCommand(*EventstatsCommand) (interface{}, error)
	VisitStreamstatsCommand(*StreamstatsCommand) (interface{}, error)
	VisitReverseCommand(*ReverseCommand) (interface{}, error)
	VisitFlattenCommand(*FlattenCommand) (interface{}, error)
	VisitAddtotalsCommand(*AddtotalsCommand) (interface{}, error)
	VisitAddcoltotalsCommand(*AddcoltotalsCommand) (interface{}, error)

	// Expressions
	VisitBinaryExpression(*BinaryExpression) (interface{}, error)
	VisitUnaryExpression(*UnaryExpression) (interface{}, error)
	VisitFunctionCall(*FunctionCall) (interface{}, error)
	VisitFieldReference(*FieldReference) (interface{}, error)
	VisitLiteral(*Literal) (interface{}, error)
	VisitListLiteral(*ListLiteral) (interface{}, error)
	VisitCaseExpression(*CaseExpression) (interface{}, error)
	VisitWhenClause(*WhenClause) (interface{}, error)

	// Other nodes
	VisitAggregation(*Aggregation) (interface{}, error)
	VisitSortKey(*SortKey) (interface{}, error)
}

// BaseVisitor provides default implementations for all visitor methods
// Embed this in your visitor to only implement the methods you need
type BaseVisitor struct{}

func (b *BaseVisitor) VisitQuery(q *Query) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitSearchCommand(s *SearchCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitWhereCommand(w *WhereCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitFieldsCommand(f *FieldsCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitStatsCommand(s *StatsCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitSortCommand(s *SortCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitHeadCommand(h *HeadCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitDescribeCommand(d *DescribeCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitShowDatasourcesCommand(s *ShowDatasourcesCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitExplainCommand(e *ExplainCommand) (interface{}, error) {
	return nil, nil
}

// Tier 1 commands
func (b *BaseVisitor) VisitChartCommand(c *ChartCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitTimechartCommand(t *TimechartCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitBinCommand(bin *BinCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitDedupCommand(d *DedupCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitTopCommand(t *TopCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitRareCommand(r *RareCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitEvalCommand(e *EvalCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitRenameCommand(r *RenameCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitReplaceCommand(r *ReplaceCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitFillnullCommand(f *FillnullCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitParseCommand(p *ParseCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitRexCommand(r *RexCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitLookupCommand(l *LookupCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitAppendCommand(a *AppendCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitJoinCommand(j *JoinCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitTableCommand(t *TableCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitEventstatsCommand(e *EventstatsCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitStreamstatsCommand(s *StreamstatsCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitReverseCommand(r *ReverseCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitFlattenCommand(f *FlattenCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitAddtotalsCommand(a *AddtotalsCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitAddcoltotalsCommand(a *AddcoltotalsCommand) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitBinaryExpression(be *BinaryExpression) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitUnaryExpression(u *UnaryExpression) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitFunctionCall(f *FunctionCall) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitFieldReference(f *FieldReference) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitLiteral(l *Literal) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitListLiteral(l *ListLiteral) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitCaseExpression(c *CaseExpression) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitWhenClause(w *WhenClause) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitAggregation(a *Aggregation) (interface{}, error) {
	return nil, nil
}

func (b *BaseVisitor) VisitSortKey(s *SortKey) (interface{}, error) {
	return nil, nil
}

// Walk traverses the AST starting from the given node
func Walk(node Node, visitor Visitor) (interface{}, error) {
	return node.Accept(visitor)
}
