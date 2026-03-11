// Code generated from PPLParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package generated // PPLParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by PPLParser.
type PPLParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by PPLParser#query.
	VisitQuery(ctx *QueryContext) interface{}

	// Visit a parse tree produced by PPLParser#searchQuery.
	VisitSearchQuery(ctx *SearchQueryContext) interface{}

	// Visit a parse tree produced by PPLParser#processingCommand.
	VisitProcessingCommand(ctx *ProcessingCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#metadataCommand.
	VisitMetadataCommand(ctx *MetadataCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#SearchWithKeyword.
	VisitSearchWithKeyword(ctx *SearchWithKeywordContext) interface{}

	// Visit a parse tree produced by PPLParser#SearchWithSource.
	VisitSearchWithSource(ctx *SearchWithSourceContext) interface{}

	// Visit a parse tree produced by PPLParser#whereCommand.
	VisitWhereCommand(ctx *WhereCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#FieldsExclude.
	VisitFieldsExclude(ctx *FieldsExcludeContext) interface{}

	// Visit a parse tree produced by PPLParser#FieldsInclude.
	VisitFieldsInclude(ctx *FieldsIncludeContext) interface{}

	// Visit a parse tree produced by PPLParser#fieldList.
	VisitFieldList(ctx *FieldListContext) interface{}

	// Visit a parse tree produced by PPLParser#statsCommand.
	VisitStatsCommand(ctx *StatsCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#aggregationList.
	VisitAggregationList(ctx *AggregationListContext) interface{}

	// Visit a parse tree produced by PPLParser#aggregation.
	VisitAggregation(ctx *AggregationContext) interface{}

	// Visit a parse tree produced by PPLParser#sortCommand.
	VisitSortCommand(ctx *SortCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#sortFieldList.
	VisitSortFieldList(ctx *SortFieldListContext) interface{}

	// Visit a parse tree produced by PPLParser#sortField.
	VisitSortField(ctx *SortFieldContext) interface{}

	// Visit a parse tree produced by PPLParser#headCommand.
	VisitHeadCommand(ctx *HeadCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#chartCommand.
	VisitChartCommand(ctx *ChartCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#chartOptions.
	VisitChartOptions(ctx *ChartOptionsContext) interface{}

	// Visit a parse tree produced by PPLParser#timechartCommand.
	VisitTimechartCommand(ctx *TimechartCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#timechartOptions.
	VisitTimechartOptions(ctx *TimechartOptionsContext) interface{}

	// Visit a parse tree produced by PPLParser#timeSpan.
	VisitTimeSpan(ctx *TimeSpanContext) interface{}

	// Visit a parse tree produced by PPLParser#binCommand.
	VisitBinCommand(ctx *BinCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#binOptions.
	VisitBinOptions(ctx *BinOptionsContext) interface{}

	// Visit a parse tree produced by PPLParser#dedupCommand.
	VisitDedupCommand(ctx *DedupCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#dedupOptions.
	VisitDedupOptions(ctx *DedupOptionsContext) interface{}

	// Visit a parse tree produced by PPLParser#topCommand.
	VisitTopCommand(ctx *TopCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#topOptions.
	VisitTopOptions(ctx *TopOptionsContext) interface{}

	// Visit a parse tree produced by PPLParser#rareCommand.
	VisitRareCommand(ctx *RareCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#evalCommand.
	VisitEvalCommand(ctx *EvalCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#evalAssignment.
	VisitEvalAssignment(ctx *EvalAssignmentContext) interface{}

	// Visit a parse tree produced by PPLParser#renameCommand.
	VisitRenameCommand(ctx *RenameCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#renameAssignment.
	VisitRenameAssignment(ctx *RenameAssignmentContext) interface{}

	// Visit a parse tree produced by PPLParser#replaceCommand.
	VisitReplaceCommand(ctx *ReplaceCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#replaceMapping.
	VisitReplaceMapping(ctx *ReplaceMappingContext) interface{}

	// Visit a parse tree produced by PPLParser#FillnullWithDefault.
	VisitFillnullWithDefault(ctx *FillnullWithDefaultContext) interface{}

	// Visit a parse tree produced by PPLParser#FillnullWithAssignments.
	VisitFillnullWithAssignments(ctx *FillnullWithAssignmentsContext) interface{}

	// Visit a parse tree produced by PPLParser#fillnullAssignment.
	VisitFillnullAssignment(ctx *FillnullAssignmentContext) interface{}

	// Visit a parse tree produced by PPLParser#parseCommand.
	VisitParseCommand(ctx *ParseCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#rexCommand.
	VisitRexCommand(ctx *RexCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#lookupCommand.
	VisitLookupCommand(ctx *LookupCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#lookupOutputList.
	VisitLookupOutputList(ctx *LookupOutputListContext) interface{}

	// Visit a parse tree produced by PPLParser#lookupOutputField.
	VisitLookupOutputField(ctx *LookupOutputFieldContext) interface{}

	// Visit a parse tree produced by PPLParser#appendCommand.
	VisitAppendCommand(ctx *AppendCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#joinCommand.
	VisitJoinCommand(ctx *JoinCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#joinType.
	VisitJoinType(ctx *JoinTypeContext) interface{}

	// Visit a parse tree produced by PPLParser#tableCommand.
	VisitTableCommand(ctx *TableCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#eventstatsCommand.
	VisitEventstatsCommand(ctx *EventstatsCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#streamstatsCommand.
	VisitStreamstatsCommand(ctx *StreamstatsCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#streamstatsOptions.
	VisitStreamstatsOptions(ctx *StreamstatsOptionsContext) interface{}

	// Visit a parse tree produced by PPLParser#reverseCommand.
	VisitReverseCommand(ctx *ReverseCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#flattenCommand.
	VisitFlattenCommand(ctx *FlattenCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#booleanValue.
	VisitBooleanValue(ctx *BooleanValueContext) interface{}

	// Visit a parse tree produced by PPLParser#describeCommand.
	VisitDescribeCommand(ctx *DescribeCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#showDatasourcesCommand.
	VisitShowDatasourcesCommand(ctx *ShowDatasourcesCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#explainCommand.
	VisitExplainCommand(ctx *ExplainCommandContext) interface{}

	// Visit a parse tree produced by PPLParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#orExpression.
	VisitOrExpression(ctx *OrExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#andExpression.
	VisitAndExpression(ctx *AndExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#notExpression.
	VisitNotExpression(ctx *NotExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#comparisonExpression.
	VisitComparisonExpression(ctx *ComparisonExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#additiveExpression.
	VisitAdditiveExpression(ctx *AdditiveExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#multiplicativeExpression.
	VisitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#unaryExpression.
	VisitUnaryExpression(ctx *UnaryExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#primaryExpression.
	VisitPrimaryExpression(ctx *PrimaryExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}

	// Visit a parse tree produced by PPLParser#fieldReference.
	VisitFieldReference(ctx *FieldReferenceContext) interface{}

	// Visit a parse tree produced by PPLParser#FunctionCallNoArgs.
	VisitFunctionCallNoArgs(ctx *FunctionCallNoArgsContext) interface{}

	// Visit a parse tree produced by PPLParser#FunctionCallWithArgs.
	VisitFunctionCallWithArgs(ctx *FunctionCallWithArgsContext) interface{}

	// Visit a parse tree produced by PPLParser#AggregationFunctionCallNoArgs.
	VisitAggregationFunctionCallNoArgs(ctx *AggregationFunctionCallNoArgsContext) interface{}

	// Visit a parse tree produced by PPLParser#AggregationFunctionCall.
	VisitAggregationFunctionCall(ctx *AggregationFunctionCallContext) interface{}

	// Visit a parse tree produced by PPLParser#aggregationFunction.
	VisitAggregationFunction(ctx *AggregationFunctionContext) interface{}

	// Visit a parse tree produced by PPLParser#expressionList.
	VisitExpressionList(ctx *ExpressionListContext) interface{}

	// Visit a parse tree produced by PPLParser#caseExpression.
	VisitCaseExpression(ctx *CaseExpressionContext) interface{}

	// Visit a parse tree produced by PPLParser#whenClause.
	VisitWhenClause(ctx *WhenClauseContext) interface{}
}
