// Code generated from PPLParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package generated // PPLParser
import "github.com/antlr4-go/antlr/v4"

type BasePPLParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BasePPLParserVisitor) VisitQuery(ctx *QueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitSearchQuery(ctx *SearchQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitProcessingCommand(ctx *ProcessingCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitMetadataCommand(ctx *MetadataCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitSearchWithKeyword(ctx *SearchWithKeywordContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitSearchWithSource(ctx *SearchWithSourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitWhereCommand(ctx *WhereCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFieldsExclude(ctx *FieldsExcludeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFieldsInclude(ctx *FieldsIncludeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFieldList(ctx *FieldListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitStatsCommand(ctx *StatsCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitAggregationList(ctx *AggregationListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitAggregation(ctx *AggregationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitSortCommand(ctx *SortCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitSortFieldList(ctx *SortFieldListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitSortField(ctx *SortFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitHeadCommand(ctx *HeadCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitChartCommand(ctx *ChartCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitChartOptions(ctx *ChartOptionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitTimechartCommand(ctx *TimechartCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitTimechartOptions(ctx *TimechartOptionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitTimeSpan(ctx *TimeSpanContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitBinCommand(ctx *BinCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitBinOptions(ctx *BinOptionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitDedupCommand(ctx *DedupCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitDedupOptions(ctx *DedupOptionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitTopCommand(ctx *TopCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitTopOptions(ctx *TopOptionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitRareCommand(ctx *RareCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitEvalCommand(ctx *EvalCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitEvalAssignment(ctx *EvalAssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitRenameCommand(ctx *RenameCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitRenameAssignment(ctx *RenameAssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitReplaceCommand(ctx *ReplaceCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitReplaceMapping(ctx *ReplaceMappingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFillnullWithDefault(ctx *FillnullWithDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFillnullWithAssignments(ctx *FillnullWithAssignmentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFillnullAssignment(ctx *FillnullAssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitParseCommand(ctx *ParseCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitRexCommand(ctx *RexCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitLookupCommand(ctx *LookupCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitLookupOutputList(ctx *LookupOutputListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitLookupOutputField(ctx *LookupOutputFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitAppendCommand(ctx *AppendCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitJoinCommand(ctx *JoinCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitJoinType(ctx *JoinTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitTableCommand(ctx *TableCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitEventstatsCommand(ctx *EventstatsCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitStreamstatsCommand(ctx *StreamstatsCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitStreamstatsOptions(ctx *StreamstatsOptionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitReverseCommand(ctx *ReverseCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFlattenCommand(ctx *FlattenCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitBooleanValue(ctx *BooleanValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitDescribeCommand(ctx *DescribeCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitShowDatasourcesCommand(ctx *ShowDatasourcesCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitExplainCommand(ctx *ExplainCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitOrExpression(ctx *OrExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitAndExpression(ctx *AndExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitNotExpression(ctx *NotExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitComparisonExpression(ctx *ComparisonExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitAdditiveExpression(ctx *AdditiveExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitUnaryExpression(ctx *UnaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitPrimaryExpression(ctx *PrimaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFieldReference(ctx *FieldReferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFunctionCallNoArgs(ctx *FunctionCallNoArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitFunctionCallWithArgs(ctx *FunctionCallWithArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitAggregationFunctionCallNoArgs(ctx *AggregationFunctionCallNoArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitAggregationFunctionCall(ctx *AggregationFunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitAggregationFunction(ctx *AggregationFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitExpressionList(ctx *ExpressionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitCaseExpression(ctx *CaseExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasePPLParserVisitor) VisitWhenClause(ctx *WhenClauseContext) interface{} {
	return v.VisitChildren(ctx)
}
