// Code generated from PPLParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package generated // PPLParser
import "github.com/antlr4-go/antlr/v4"

// PPLParserListener is a complete listener for a parse tree produced by PPLParser.
type PPLParserListener interface {
	antlr.ParseTreeListener

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterSearchQuery is called when entering the searchQuery production.
	EnterSearchQuery(c *SearchQueryContext)

	// EnterProcessingCommand is called when entering the processingCommand production.
	EnterProcessingCommand(c *ProcessingCommandContext)

	// EnterMetadataCommand is called when entering the metadataCommand production.
	EnterMetadataCommand(c *MetadataCommandContext)

	// EnterSearchWithKeyword is called when entering the SearchWithKeyword production.
	EnterSearchWithKeyword(c *SearchWithKeywordContext)

	// EnterSearchWithSource is called when entering the SearchWithSource production.
	EnterSearchWithSource(c *SearchWithSourceContext)

	// EnterWhereCommand is called when entering the whereCommand production.
	EnterWhereCommand(c *WhereCommandContext)

	// EnterFieldsExclude is called when entering the FieldsExclude production.
	EnterFieldsExclude(c *FieldsExcludeContext)

	// EnterFieldsInclude is called when entering the FieldsInclude production.
	EnterFieldsInclude(c *FieldsIncludeContext)

	// EnterFieldList is called when entering the fieldList production.
	EnterFieldList(c *FieldListContext)

	// EnterStatsCommand is called when entering the statsCommand production.
	EnterStatsCommand(c *StatsCommandContext)

	// EnterAggregationList is called when entering the aggregationList production.
	EnterAggregationList(c *AggregationListContext)

	// EnterAggregation is called when entering the aggregation production.
	EnterAggregation(c *AggregationContext)

	// EnterSortCommand is called when entering the sortCommand production.
	EnterSortCommand(c *SortCommandContext)

	// EnterSortFieldList is called when entering the sortFieldList production.
	EnterSortFieldList(c *SortFieldListContext)

	// EnterSortField is called when entering the sortField production.
	EnterSortField(c *SortFieldContext)

	// EnterHeadCommand is called when entering the headCommand production.
	EnterHeadCommand(c *HeadCommandContext)

	// EnterChartCommand is called when entering the chartCommand production.
	EnterChartCommand(c *ChartCommandContext)

	// EnterChartOptions is called when entering the chartOptions production.
	EnterChartOptions(c *ChartOptionsContext)

	// EnterTimechartCommand is called when entering the timechartCommand production.
	EnterTimechartCommand(c *TimechartCommandContext)

	// EnterTimechartOptions is called when entering the timechartOptions production.
	EnterTimechartOptions(c *TimechartOptionsContext)

	// EnterTimeSpan is called when entering the timeSpan production.
	EnterTimeSpan(c *TimeSpanContext)

	// EnterBinCommand is called when entering the binCommand production.
	EnterBinCommand(c *BinCommandContext)

	// EnterBinOptions is called when entering the binOptions production.
	EnterBinOptions(c *BinOptionsContext)

	// EnterDedupCommand is called when entering the dedupCommand production.
	EnterDedupCommand(c *DedupCommandContext)

	// EnterDedupOptions is called when entering the dedupOptions production.
	EnterDedupOptions(c *DedupOptionsContext)

	// EnterTopCommand is called when entering the topCommand production.
	EnterTopCommand(c *TopCommandContext)

	// EnterTopOptions is called when entering the topOptions production.
	EnterTopOptions(c *TopOptionsContext)

	// EnterRareCommand is called when entering the rareCommand production.
	EnterRareCommand(c *RareCommandContext)

	// EnterEvalCommand is called when entering the evalCommand production.
	EnterEvalCommand(c *EvalCommandContext)

	// EnterEvalAssignment is called when entering the evalAssignment production.
	EnterEvalAssignment(c *EvalAssignmentContext)

	// EnterRenameCommand is called when entering the renameCommand production.
	EnterRenameCommand(c *RenameCommandContext)

	// EnterRenameAssignment is called when entering the renameAssignment production.
	EnterRenameAssignment(c *RenameAssignmentContext)

	// EnterReplaceCommand is called when entering the replaceCommand production.
	EnterReplaceCommand(c *ReplaceCommandContext)

	// EnterReplaceMapping is called when entering the replaceMapping production.
	EnterReplaceMapping(c *ReplaceMappingContext)

	// EnterFillnullWithDefault is called when entering the FillnullWithDefault production.
	EnterFillnullWithDefault(c *FillnullWithDefaultContext)

	// EnterFillnullWithAssignments is called when entering the FillnullWithAssignments production.
	EnterFillnullWithAssignments(c *FillnullWithAssignmentsContext)

	// EnterFillnullAssignment is called when entering the fillnullAssignment production.
	EnterFillnullAssignment(c *FillnullAssignmentContext)

	// EnterParseCommand is called when entering the parseCommand production.
	EnterParseCommand(c *ParseCommandContext)

	// EnterRexCommand is called when entering the rexCommand production.
	EnterRexCommand(c *RexCommandContext)

	// EnterLookupCommand is called when entering the lookupCommand production.
	EnterLookupCommand(c *LookupCommandContext)

	// EnterLookupOutputList is called when entering the lookupOutputList production.
	EnterLookupOutputList(c *LookupOutputListContext)

	// EnterLookupOutputField is called when entering the lookupOutputField production.
	EnterLookupOutputField(c *LookupOutputFieldContext)

	// EnterAppendCommand is called when entering the appendCommand production.
	EnterAppendCommand(c *AppendCommandContext)

	// EnterJoinCommand is called when entering the joinCommand production.
	EnterJoinCommand(c *JoinCommandContext)

	// EnterJoinType is called when entering the joinType production.
	EnterJoinType(c *JoinTypeContext)

	// EnterTableCommand is called when entering the tableCommand production.
	EnterTableCommand(c *TableCommandContext)

	// EnterEventstatsCommand is called when entering the eventstatsCommand production.
	EnterEventstatsCommand(c *EventstatsCommandContext)

	// EnterStreamstatsCommand is called when entering the streamstatsCommand production.
	EnterStreamstatsCommand(c *StreamstatsCommandContext)

	// EnterStreamstatsOptions is called when entering the streamstatsOptions production.
	EnterStreamstatsOptions(c *StreamstatsOptionsContext)

	// EnterReverseCommand is called when entering the reverseCommand production.
	EnterReverseCommand(c *ReverseCommandContext)

	// EnterFlattenCommand is called when entering the flattenCommand production.
	EnterFlattenCommand(c *FlattenCommandContext)

	// EnterBooleanValue is called when entering the booleanValue production.
	EnterBooleanValue(c *BooleanValueContext)

	// EnterDescribeCommand is called when entering the describeCommand production.
	EnterDescribeCommand(c *DescribeCommandContext)

	// EnterShowDatasourcesCommand is called when entering the showDatasourcesCommand production.
	EnterShowDatasourcesCommand(c *ShowDatasourcesCommandContext)

	// EnterExplainCommand is called when entering the explainCommand production.
	EnterExplainCommand(c *ExplainCommandContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterOrExpression is called when entering the orExpression production.
	EnterOrExpression(c *OrExpressionContext)

	// EnterAndExpression is called when entering the andExpression production.
	EnterAndExpression(c *AndExpressionContext)

	// EnterNotExpression is called when entering the notExpression production.
	EnterNotExpression(c *NotExpressionContext)

	// EnterComparisonExpression is called when entering the comparisonExpression production.
	EnterComparisonExpression(c *ComparisonExpressionContext)

	// EnterAdditiveExpression is called when entering the additiveExpression production.
	EnterAdditiveExpression(c *AdditiveExpressionContext)

	// EnterMultiplicativeExpression is called when entering the multiplicativeExpression production.
	EnterMultiplicativeExpression(c *MultiplicativeExpressionContext)

	// EnterUnaryExpression is called when entering the unaryExpression production.
	EnterUnaryExpression(c *UnaryExpressionContext)

	// EnterPrimaryExpression is called when entering the primaryExpression production.
	EnterPrimaryExpression(c *PrimaryExpressionContext)

	// EnterLiteral is called when entering the literal production.
	EnterLiteral(c *LiteralContext)

	// EnterFieldReference is called when entering the fieldReference production.
	EnterFieldReference(c *FieldReferenceContext)

	// EnterFunctionCallNoArgs is called when entering the FunctionCallNoArgs production.
	EnterFunctionCallNoArgs(c *FunctionCallNoArgsContext)

	// EnterFunctionCallWithArgs is called when entering the FunctionCallWithArgs production.
	EnterFunctionCallWithArgs(c *FunctionCallWithArgsContext)

	// EnterAggregationFunctionCallNoArgs is called when entering the AggregationFunctionCallNoArgs production.
	EnterAggregationFunctionCallNoArgs(c *AggregationFunctionCallNoArgsContext)

	// EnterAggregationFunctionCall is called when entering the AggregationFunctionCall production.
	EnterAggregationFunctionCall(c *AggregationFunctionCallContext)

	// EnterAggregationFunction is called when entering the aggregationFunction production.
	EnterAggregationFunction(c *AggregationFunctionContext)

	// EnterExpressionList is called when entering the expressionList production.
	EnterExpressionList(c *ExpressionListContext)

	// EnterCaseExpression is called when entering the caseExpression production.
	EnterCaseExpression(c *CaseExpressionContext)

	// EnterWhenClause is called when entering the whenClause production.
	EnterWhenClause(c *WhenClauseContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitSearchQuery is called when exiting the searchQuery production.
	ExitSearchQuery(c *SearchQueryContext)

	// ExitProcessingCommand is called when exiting the processingCommand production.
	ExitProcessingCommand(c *ProcessingCommandContext)

	// ExitMetadataCommand is called when exiting the metadataCommand production.
	ExitMetadataCommand(c *MetadataCommandContext)

	// ExitSearchWithKeyword is called when exiting the SearchWithKeyword production.
	ExitSearchWithKeyword(c *SearchWithKeywordContext)

	// ExitSearchWithSource is called when exiting the SearchWithSource production.
	ExitSearchWithSource(c *SearchWithSourceContext)

	// ExitWhereCommand is called when exiting the whereCommand production.
	ExitWhereCommand(c *WhereCommandContext)

	// ExitFieldsExclude is called when exiting the FieldsExclude production.
	ExitFieldsExclude(c *FieldsExcludeContext)

	// ExitFieldsInclude is called when exiting the FieldsInclude production.
	ExitFieldsInclude(c *FieldsIncludeContext)

	// ExitFieldList is called when exiting the fieldList production.
	ExitFieldList(c *FieldListContext)

	// ExitStatsCommand is called when exiting the statsCommand production.
	ExitStatsCommand(c *StatsCommandContext)

	// ExitAggregationList is called when exiting the aggregationList production.
	ExitAggregationList(c *AggregationListContext)

	// ExitAggregation is called when exiting the aggregation production.
	ExitAggregation(c *AggregationContext)

	// ExitSortCommand is called when exiting the sortCommand production.
	ExitSortCommand(c *SortCommandContext)

	// ExitSortFieldList is called when exiting the sortFieldList production.
	ExitSortFieldList(c *SortFieldListContext)

	// ExitSortField is called when exiting the sortField production.
	ExitSortField(c *SortFieldContext)

	// ExitHeadCommand is called when exiting the headCommand production.
	ExitHeadCommand(c *HeadCommandContext)

	// ExitChartCommand is called when exiting the chartCommand production.
	ExitChartCommand(c *ChartCommandContext)

	// ExitChartOptions is called when exiting the chartOptions production.
	ExitChartOptions(c *ChartOptionsContext)

	// ExitTimechartCommand is called when exiting the timechartCommand production.
	ExitTimechartCommand(c *TimechartCommandContext)

	// ExitTimechartOptions is called when exiting the timechartOptions production.
	ExitTimechartOptions(c *TimechartOptionsContext)

	// ExitTimeSpan is called when exiting the timeSpan production.
	ExitTimeSpan(c *TimeSpanContext)

	// ExitBinCommand is called when exiting the binCommand production.
	ExitBinCommand(c *BinCommandContext)

	// ExitBinOptions is called when exiting the binOptions production.
	ExitBinOptions(c *BinOptionsContext)

	// ExitDedupCommand is called when exiting the dedupCommand production.
	ExitDedupCommand(c *DedupCommandContext)

	// ExitDedupOptions is called when exiting the dedupOptions production.
	ExitDedupOptions(c *DedupOptionsContext)

	// ExitTopCommand is called when exiting the topCommand production.
	ExitTopCommand(c *TopCommandContext)

	// ExitTopOptions is called when exiting the topOptions production.
	ExitTopOptions(c *TopOptionsContext)

	// ExitRareCommand is called when exiting the rareCommand production.
	ExitRareCommand(c *RareCommandContext)

	// ExitEvalCommand is called when exiting the evalCommand production.
	ExitEvalCommand(c *EvalCommandContext)

	// ExitEvalAssignment is called when exiting the evalAssignment production.
	ExitEvalAssignment(c *EvalAssignmentContext)

	// ExitRenameCommand is called when exiting the renameCommand production.
	ExitRenameCommand(c *RenameCommandContext)

	// ExitRenameAssignment is called when exiting the renameAssignment production.
	ExitRenameAssignment(c *RenameAssignmentContext)

	// ExitReplaceCommand is called when exiting the replaceCommand production.
	ExitReplaceCommand(c *ReplaceCommandContext)

	// ExitReplaceMapping is called when exiting the replaceMapping production.
	ExitReplaceMapping(c *ReplaceMappingContext)

	// ExitFillnullWithDefault is called when exiting the FillnullWithDefault production.
	ExitFillnullWithDefault(c *FillnullWithDefaultContext)

	// ExitFillnullWithAssignments is called when exiting the FillnullWithAssignments production.
	ExitFillnullWithAssignments(c *FillnullWithAssignmentsContext)

	// ExitFillnullAssignment is called when exiting the fillnullAssignment production.
	ExitFillnullAssignment(c *FillnullAssignmentContext)

	// ExitParseCommand is called when exiting the parseCommand production.
	ExitParseCommand(c *ParseCommandContext)

	// ExitRexCommand is called when exiting the rexCommand production.
	ExitRexCommand(c *RexCommandContext)

	// ExitLookupCommand is called when exiting the lookupCommand production.
	ExitLookupCommand(c *LookupCommandContext)

	// ExitLookupOutputList is called when exiting the lookupOutputList production.
	ExitLookupOutputList(c *LookupOutputListContext)

	// ExitLookupOutputField is called when exiting the lookupOutputField production.
	ExitLookupOutputField(c *LookupOutputFieldContext)

	// ExitAppendCommand is called when exiting the appendCommand production.
	ExitAppendCommand(c *AppendCommandContext)

	// ExitJoinCommand is called when exiting the joinCommand production.
	ExitJoinCommand(c *JoinCommandContext)

	// ExitJoinType is called when exiting the joinType production.
	ExitJoinType(c *JoinTypeContext)

	// ExitTableCommand is called when exiting the tableCommand production.
	ExitTableCommand(c *TableCommandContext)

	// ExitEventstatsCommand is called when exiting the eventstatsCommand production.
	ExitEventstatsCommand(c *EventstatsCommandContext)

	// ExitStreamstatsCommand is called when exiting the streamstatsCommand production.
	ExitStreamstatsCommand(c *StreamstatsCommandContext)

	// ExitStreamstatsOptions is called when exiting the streamstatsOptions production.
	ExitStreamstatsOptions(c *StreamstatsOptionsContext)

	// ExitReverseCommand is called when exiting the reverseCommand production.
	ExitReverseCommand(c *ReverseCommandContext)

	// ExitFlattenCommand is called when exiting the flattenCommand production.
	ExitFlattenCommand(c *FlattenCommandContext)

	// ExitBooleanValue is called when exiting the booleanValue production.
	ExitBooleanValue(c *BooleanValueContext)

	// ExitDescribeCommand is called when exiting the describeCommand production.
	ExitDescribeCommand(c *DescribeCommandContext)

	// ExitShowDatasourcesCommand is called when exiting the showDatasourcesCommand production.
	ExitShowDatasourcesCommand(c *ShowDatasourcesCommandContext)

	// ExitExplainCommand is called when exiting the explainCommand production.
	ExitExplainCommand(c *ExplainCommandContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitOrExpression is called when exiting the orExpression production.
	ExitOrExpression(c *OrExpressionContext)

	// ExitAndExpression is called when exiting the andExpression production.
	ExitAndExpression(c *AndExpressionContext)

	// ExitNotExpression is called when exiting the notExpression production.
	ExitNotExpression(c *NotExpressionContext)

	// ExitComparisonExpression is called when exiting the comparisonExpression production.
	ExitComparisonExpression(c *ComparisonExpressionContext)

	// ExitAdditiveExpression is called when exiting the additiveExpression production.
	ExitAdditiveExpression(c *AdditiveExpressionContext)

	// ExitMultiplicativeExpression is called when exiting the multiplicativeExpression production.
	ExitMultiplicativeExpression(c *MultiplicativeExpressionContext)

	// ExitUnaryExpression is called when exiting the unaryExpression production.
	ExitUnaryExpression(c *UnaryExpressionContext)

	// ExitPrimaryExpression is called when exiting the primaryExpression production.
	ExitPrimaryExpression(c *PrimaryExpressionContext)

	// ExitLiteral is called when exiting the literal production.
	ExitLiteral(c *LiteralContext)

	// ExitFieldReference is called when exiting the fieldReference production.
	ExitFieldReference(c *FieldReferenceContext)

	// ExitFunctionCallNoArgs is called when exiting the FunctionCallNoArgs production.
	ExitFunctionCallNoArgs(c *FunctionCallNoArgsContext)

	// ExitFunctionCallWithArgs is called when exiting the FunctionCallWithArgs production.
	ExitFunctionCallWithArgs(c *FunctionCallWithArgsContext)

	// ExitAggregationFunctionCallNoArgs is called when exiting the AggregationFunctionCallNoArgs production.
	ExitAggregationFunctionCallNoArgs(c *AggregationFunctionCallNoArgsContext)

	// ExitAggregationFunctionCall is called when exiting the AggregationFunctionCall production.
	ExitAggregationFunctionCall(c *AggregationFunctionCallContext)

	// ExitAggregationFunction is called when exiting the aggregationFunction production.
	ExitAggregationFunction(c *AggregationFunctionContext)

	// ExitExpressionList is called when exiting the expressionList production.
	ExitExpressionList(c *ExpressionListContext)

	// ExitCaseExpression is called when exiting the caseExpression production.
	ExitCaseExpression(c *CaseExpressionContext)

	// ExitWhenClause is called when exiting the whenClause production.
	ExitWhenClause(c *WhenClauseContext)
}
