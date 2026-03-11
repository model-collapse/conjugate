// Code generated from PPLParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package generated // PPLParser
import "github.com/antlr4-go/antlr/v4"

// BasePPLParserListener is a complete listener for a parse tree produced by PPLParser.
type BasePPLParserListener struct{}

var _ PPLParserListener = &BasePPLParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasePPLParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasePPLParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasePPLParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasePPLParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterQuery is called when production query is entered.
func (s *BasePPLParserListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BasePPLParserListener) ExitQuery(ctx *QueryContext) {}

// EnterSearchQuery is called when production searchQuery is entered.
func (s *BasePPLParserListener) EnterSearchQuery(ctx *SearchQueryContext) {}

// ExitSearchQuery is called when production searchQuery is exited.
func (s *BasePPLParserListener) ExitSearchQuery(ctx *SearchQueryContext) {}

// EnterProcessingCommand is called when production processingCommand is entered.
func (s *BasePPLParserListener) EnterProcessingCommand(ctx *ProcessingCommandContext) {}

// ExitProcessingCommand is called when production processingCommand is exited.
func (s *BasePPLParserListener) ExitProcessingCommand(ctx *ProcessingCommandContext) {}

// EnterMetadataCommand is called when production metadataCommand is entered.
func (s *BasePPLParserListener) EnterMetadataCommand(ctx *MetadataCommandContext) {}

// ExitMetadataCommand is called when production metadataCommand is exited.
func (s *BasePPLParserListener) ExitMetadataCommand(ctx *MetadataCommandContext) {}

// EnterSearchWithKeyword is called when production SearchWithKeyword is entered.
func (s *BasePPLParserListener) EnterSearchWithKeyword(ctx *SearchWithKeywordContext) {}

// ExitSearchWithKeyword is called when production SearchWithKeyword is exited.
func (s *BasePPLParserListener) ExitSearchWithKeyword(ctx *SearchWithKeywordContext) {}

// EnterSearchWithSource is called when production SearchWithSource is entered.
func (s *BasePPLParserListener) EnterSearchWithSource(ctx *SearchWithSourceContext) {}

// ExitSearchWithSource is called when production SearchWithSource is exited.
func (s *BasePPLParserListener) ExitSearchWithSource(ctx *SearchWithSourceContext) {}

// EnterWhereCommand is called when production whereCommand is entered.
func (s *BasePPLParserListener) EnterWhereCommand(ctx *WhereCommandContext) {}

// ExitWhereCommand is called when production whereCommand is exited.
func (s *BasePPLParserListener) ExitWhereCommand(ctx *WhereCommandContext) {}

// EnterFieldsExclude is called when production FieldsExclude is entered.
func (s *BasePPLParserListener) EnterFieldsExclude(ctx *FieldsExcludeContext) {}

// ExitFieldsExclude is called when production FieldsExclude is exited.
func (s *BasePPLParserListener) ExitFieldsExclude(ctx *FieldsExcludeContext) {}

// EnterFieldsInclude is called when production FieldsInclude is entered.
func (s *BasePPLParserListener) EnterFieldsInclude(ctx *FieldsIncludeContext) {}

// ExitFieldsInclude is called when production FieldsInclude is exited.
func (s *BasePPLParserListener) ExitFieldsInclude(ctx *FieldsIncludeContext) {}

// EnterFieldList is called when production fieldList is entered.
func (s *BasePPLParserListener) EnterFieldList(ctx *FieldListContext) {}

// ExitFieldList is called when production fieldList is exited.
func (s *BasePPLParserListener) ExitFieldList(ctx *FieldListContext) {}

// EnterStatsCommand is called when production statsCommand is entered.
func (s *BasePPLParserListener) EnterStatsCommand(ctx *StatsCommandContext) {}

// ExitStatsCommand is called when production statsCommand is exited.
func (s *BasePPLParserListener) ExitStatsCommand(ctx *StatsCommandContext) {}

// EnterAggregationList is called when production aggregationList is entered.
func (s *BasePPLParserListener) EnterAggregationList(ctx *AggregationListContext) {}

// ExitAggregationList is called when production aggregationList is exited.
func (s *BasePPLParserListener) ExitAggregationList(ctx *AggregationListContext) {}

// EnterAggregation is called when production aggregation is entered.
func (s *BasePPLParserListener) EnterAggregation(ctx *AggregationContext) {}

// ExitAggregation is called when production aggregation is exited.
func (s *BasePPLParserListener) ExitAggregation(ctx *AggregationContext) {}

// EnterSortCommand is called when production sortCommand is entered.
func (s *BasePPLParserListener) EnterSortCommand(ctx *SortCommandContext) {}

// ExitSortCommand is called when production sortCommand is exited.
func (s *BasePPLParserListener) ExitSortCommand(ctx *SortCommandContext) {}

// EnterSortFieldList is called when production sortFieldList is entered.
func (s *BasePPLParserListener) EnterSortFieldList(ctx *SortFieldListContext) {}

// ExitSortFieldList is called when production sortFieldList is exited.
func (s *BasePPLParserListener) ExitSortFieldList(ctx *SortFieldListContext) {}

// EnterSortField is called when production sortField is entered.
func (s *BasePPLParserListener) EnterSortField(ctx *SortFieldContext) {}

// ExitSortField is called when production sortField is exited.
func (s *BasePPLParserListener) ExitSortField(ctx *SortFieldContext) {}

// EnterHeadCommand is called when production headCommand is entered.
func (s *BasePPLParserListener) EnterHeadCommand(ctx *HeadCommandContext) {}

// ExitHeadCommand is called when production headCommand is exited.
func (s *BasePPLParserListener) ExitHeadCommand(ctx *HeadCommandContext) {}

// EnterChartCommand is called when production chartCommand is entered.
func (s *BasePPLParserListener) EnterChartCommand(ctx *ChartCommandContext) {}

// ExitChartCommand is called when production chartCommand is exited.
func (s *BasePPLParserListener) ExitChartCommand(ctx *ChartCommandContext) {}

// EnterChartOptions is called when production chartOptions is entered.
func (s *BasePPLParserListener) EnterChartOptions(ctx *ChartOptionsContext) {}

// ExitChartOptions is called when production chartOptions is exited.
func (s *BasePPLParserListener) ExitChartOptions(ctx *ChartOptionsContext) {}

// EnterTimechartCommand is called when production timechartCommand is entered.
func (s *BasePPLParserListener) EnterTimechartCommand(ctx *TimechartCommandContext) {}

// ExitTimechartCommand is called when production timechartCommand is exited.
func (s *BasePPLParserListener) ExitTimechartCommand(ctx *TimechartCommandContext) {}

// EnterTimechartOptions is called when production timechartOptions is entered.
func (s *BasePPLParserListener) EnterTimechartOptions(ctx *TimechartOptionsContext) {}

// ExitTimechartOptions is called when production timechartOptions is exited.
func (s *BasePPLParserListener) ExitTimechartOptions(ctx *TimechartOptionsContext) {}

// EnterTimeSpan is called when production timeSpan is entered.
func (s *BasePPLParserListener) EnterTimeSpan(ctx *TimeSpanContext) {}

// ExitTimeSpan is called when production timeSpan is exited.
func (s *BasePPLParserListener) ExitTimeSpan(ctx *TimeSpanContext) {}

// EnterBinCommand is called when production binCommand is entered.
func (s *BasePPLParserListener) EnterBinCommand(ctx *BinCommandContext) {}

// ExitBinCommand is called when production binCommand is exited.
func (s *BasePPLParserListener) ExitBinCommand(ctx *BinCommandContext) {}

// EnterBinOptions is called when production binOptions is entered.
func (s *BasePPLParserListener) EnterBinOptions(ctx *BinOptionsContext) {}

// ExitBinOptions is called when production binOptions is exited.
func (s *BasePPLParserListener) ExitBinOptions(ctx *BinOptionsContext) {}

// EnterDedupCommand is called when production dedupCommand is entered.
func (s *BasePPLParserListener) EnterDedupCommand(ctx *DedupCommandContext) {}

// ExitDedupCommand is called when production dedupCommand is exited.
func (s *BasePPLParserListener) ExitDedupCommand(ctx *DedupCommandContext) {}

// EnterDedupOptions is called when production dedupOptions is entered.
func (s *BasePPLParserListener) EnterDedupOptions(ctx *DedupOptionsContext) {}

// ExitDedupOptions is called when production dedupOptions is exited.
func (s *BasePPLParserListener) ExitDedupOptions(ctx *DedupOptionsContext) {}

// EnterTopCommand is called when production topCommand is entered.
func (s *BasePPLParserListener) EnterTopCommand(ctx *TopCommandContext) {}

// ExitTopCommand is called when production topCommand is exited.
func (s *BasePPLParserListener) ExitTopCommand(ctx *TopCommandContext) {}

// EnterTopOptions is called when production topOptions is entered.
func (s *BasePPLParserListener) EnterTopOptions(ctx *TopOptionsContext) {}

// ExitTopOptions is called when production topOptions is exited.
func (s *BasePPLParserListener) ExitTopOptions(ctx *TopOptionsContext) {}

// EnterRareCommand is called when production rareCommand is entered.
func (s *BasePPLParserListener) EnterRareCommand(ctx *RareCommandContext) {}

// ExitRareCommand is called when production rareCommand is exited.
func (s *BasePPLParserListener) ExitRareCommand(ctx *RareCommandContext) {}

// EnterEvalCommand is called when production evalCommand is entered.
func (s *BasePPLParserListener) EnterEvalCommand(ctx *EvalCommandContext) {}

// ExitEvalCommand is called when production evalCommand is exited.
func (s *BasePPLParserListener) ExitEvalCommand(ctx *EvalCommandContext) {}

// EnterEvalAssignment is called when production evalAssignment is entered.
func (s *BasePPLParserListener) EnterEvalAssignment(ctx *EvalAssignmentContext) {}

// ExitEvalAssignment is called when production evalAssignment is exited.
func (s *BasePPLParserListener) ExitEvalAssignment(ctx *EvalAssignmentContext) {}

// EnterRenameCommand is called when production renameCommand is entered.
func (s *BasePPLParserListener) EnterRenameCommand(ctx *RenameCommandContext) {}

// ExitRenameCommand is called when production renameCommand is exited.
func (s *BasePPLParserListener) ExitRenameCommand(ctx *RenameCommandContext) {}

// EnterRenameAssignment is called when production renameAssignment is entered.
func (s *BasePPLParserListener) EnterRenameAssignment(ctx *RenameAssignmentContext) {}

// ExitRenameAssignment is called when production renameAssignment is exited.
func (s *BasePPLParserListener) ExitRenameAssignment(ctx *RenameAssignmentContext) {}

// EnterReplaceCommand is called when production replaceCommand is entered.
func (s *BasePPLParserListener) EnterReplaceCommand(ctx *ReplaceCommandContext) {}

// ExitReplaceCommand is called when production replaceCommand is exited.
func (s *BasePPLParserListener) ExitReplaceCommand(ctx *ReplaceCommandContext) {}

// EnterReplaceMapping is called when production replaceMapping is entered.
func (s *BasePPLParserListener) EnterReplaceMapping(ctx *ReplaceMappingContext) {}

// ExitReplaceMapping is called when production replaceMapping is exited.
func (s *BasePPLParserListener) ExitReplaceMapping(ctx *ReplaceMappingContext) {}

// EnterFillnullWithDefault is called when production FillnullWithDefault is entered.
func (s *BasePPLParserListener) EnterFillnullWithDefault(ctx *FillnullWithDefaultContext) {}

// ExitFillnullWithDefault is called when production FillnullWithDefault is exited.
func (s *BasePPLParserListener) ExitFillnullWithDefault(ctx *FillnullWithDefaultContext) {}

// EnterFillnullWithAssignments is called when production FillnullWithAssignments is entered.
func (s *BasePPLParserListener) EnterFillnullWithAssignments(ctx *FillnullWithAssignmentsContext) {}

// ExitFillnullWithAssignments is called when production FillnullWithAssignments is exited.
func (s *BasePPLParserListener) ExitFillnullWithAssignments(ctx *FillnullWithAssignmentsContext) {}

// EnterFillnullAssignment is called when production fillnullAssignment is entered.
func (s *BasePPLParserListener) EnterFillnullAssignment(ctx *FillnullAssignmentContext) {}

// ExitFillnullAssignment is called when production fillnullAssignment is exited.
func (s *BasePPLParserListener) ExitFillnullAssignment(ctx *FillnullAssignmentContext) {}

// EnterParseCommand is called when production parseCommand is entered.
func (s *BasePPLParserListener) EnterParseCommand(ctx *ParseCommandContext) {}

// ExitParseCommand is called when production parseCommand is exited.
func (s *BasePPLParserListener) ExitParseCommand(ctx *ParseCommandContext) {}

// EnterRexCommand is called when production rexCommand is entered.
func (s *BasePPLParserListener) EnterRexCommand(ctx *RexCommandContext) {}

// ExitRexCommand is called when production rexCommand is exited.
func (s *BasePPLParserListener) ExitRexCommand(ctx *RexCommandContext) {}

// EnterLookupCommand is called when production lookupCommand is entered.
func (s *BasePPLParserListener) EnterLookupCommand(ctx *LookupCommandContext) {}

// ExitLookupCommand is called when production lookupCommand is exited.
func (s *BasePPLParserListener) ExitLookupCommand(ctx *LookupCommandContext) {}

// EnterLookupOutputList is called when production lookupOutputList is entered.
func (s *BasePPLParserListener) EnterLookupOutputList(ctx *LookupOutputListContext) {}

// ExitLookupOutputList is called when production lookupOutputList is exited.
func (s *BasePPLParserListener) ExitLookupOutputList(ctx *LookupOutputListContext) {}

// EnterLookupOutputField is called when production lookupOutputField is entered.
func (s *BasePPLParserListener) EnterLookupOutputField(ctx *LookupOutputFieldContext) {}

// ExitLookupOutputField is called when production lookupOutputField is exited.
func (s *BasePPLParserListener) ExitLookupOutputField(ctx *LookupOutputFieldContext) {}

// EnterAppendCommand is called when production appendCommand is entered.
func (s *BasePPLParserListener) EnterAppendCommand(ctx *AppendCommandContext) {}

// ExitAppendCommand is called when production appendCommand is exited.
func (s *BasePPLParserListener) ExitAppendCommand(ctx *AppendCommandContext) {}

// EnterJoinCommand is called when production joinCommand is entered.
func (s *BasePPLParserListener) EnterJoinCommand(ctx *JoinCommandContext) {}

// ExitJoinCommand is called when production joinCommand is exited.
func (s *BasePPLParserListener) ExitJoinCommand(ctx *JoinCommandContext) {}

// EnterJoinType is called when production joinType is entered.
func (s *BasePPLParserListener) EnterJoinType(ctx *JoinTypeContext) {}

// ExitJoinType is called when production joinType is exited.
func (s *BasePPLParserListener) ExitJoinType(ctx *JoinTypeContext) {}

// EnterTableCommand is called when production tableCommand is entered.
func (s *BasePPLParserListener) EnterTableCommand(ctx *TableCommandContext) {}

// ExitTableCommand is called when production tableCommand is exited.
func (s *BasePPLParserListener) ExitTableCommand(ctx *TableCommandContext) {}

// EnterEventstatsCommand is called when production eventstatsCommand is entered.
func (s *BasePPLParserListener) EnterEventstatsCommand(ctx *EventstatsCommandContext) {}

// ExitEventstatsCommand is called when production eventstatsCommand is exited.
func (s *BasePPLParserListener) ExitEventstatsCommand(ctx *EventstatsCommandContext) {}

// EnterStreamstatsCommand is called when production streamstatsCommand is entered.
func (s *BasePPLParserListener) EnterStreamstatsCommand(ctx *StreamstatsCommandContext) {}

// ExitStreamstatsCommand is called when production streamstatsCommand is exited.
func (s *BasePPLParserListener) ExitStreamstatsCommand(ctx *StreamstatsCommandContext) {}

// EnterStreamstatsOptions is called when production streamstatsOptions is entered.
func (s *BasePPLParserListener) EnterStreamstatsOptions(ctx *StreamstatsOptionsContext) {}

// ExitStreamstatsOptions is called when production streamstatsOptions is exited.
func (s *BasePPLParserListener) ExitStreamstatsOptions(ctx *StreamstatsOptionsContext) {}

// EnterReverseCommand is called when production reverseCommand is entered.
func (s *BasePPLParserListener) EnterReverseCommand(ctx *ReverseCommandContext) {}

// ExitReverseCommand is called when production reverseCommand is exited.
func (s *BasePPLParserListener) ExitReverseCommand(ctx *ReverseCommandContext) {}

// EnterFlattenCommand is called when production flattenCommand is entered.
func (s *BasePPLParserListener) EnterFlattenCommand(ctx *FlattenCommandContext) {}

// ExitFlattenCommand is called when production flattenCommand is exited.
func (s *BasePPLParserListener) ExitFlattenCommand(ctx *FlattenCommandContext) {}

// EnterBooleanValue is called when production booleanValue is entered.
func (s *BasePPLParserListener) EnterBooleanValue(ctx *BooleanValueContext) {}

// ExitBooleanValue is called when production booleanValue is exited.
func (s *BasePPLParserListener) ExitBooleanValue(ctx *BooleanValueContext) {}

// EnterDescribeCommand is called when production describeCommand is entered.
func (s *BasePPLParserListener) EnterDescribeCommand(ctx *DescribeCommandContext) {}

// ExitDescribeCommand is called when production describeCommand is exited.
func (s *BasePPLParserListener) ExitDescribeCommand(ctx *DescribeCommandContext) {}

// EnterShowDatasourcesCommand is called when production showDatasourcesCommand is entered.
func (s *BasePPLParserListener) EnterShowDatasourcesCommand(ctx *ShowDatasourcesCommandContext) {}

// ExitShowDatasourcesCommand is called when production showDatasourcesCommand is exited.
func (s *BasePPLParserListener) ExitShowDatasourcesCommand(ctx *ShowDatasourcesCommandContext) {}

// EnterExplainCommand is called when production explainCommand is entered.
func (s *BasePPLParserListener) EnterExplainCommand(ctx *ExplainCommandContext) {}

// ExitExplainCommand is called when production explainCommand is exited.
func (s *BasePPLParserListener) ExitExplainCommand(ctx *ExplainCommandContext) {}

// EnterExpression is called when production expression is entered.
func (s *BasePPLParserListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BasePPLParserListener) ExitExpression(ctx *ExpressionContext) {}

// EnterOrExpression is called when production orExpression is entered.
func (s *BasePPLParserListener) EnterOrExpression(ctx *OrExpressionContext) {}

// ExitOrExpression is called when production orExpression is exited.
func (s *BasePPLParserListener) ExitOrExpression(ctx *OrExpressionContext) {}

// EnterAndExpression is called when production andExpression is entered.
func (s *BasePPLParserListener) EnterAndExpression(ctx *AndExpressionContext) {}

// ExitAndExpression is called when production andExpression is exited.
func (s *BasePPLParserListener) ExitAndExpression(ctx *AndExpressionContext) {}

// EnterNotExpression is called when production notExpression is entered.
func (s *BasePPLParserListener) EnterNotExpression(ctx *NotExpressionContext) {}

// ExitNotExpression is called when production notExpression is exited.
func (s *BasePPLParserListener) ExitNotExpression(ctx *NotExpressionContext) {}

// EnterComparisonExpression is called when production comparisonExpression is entered.
func (s *BasePPLParserListener) EnterComparisonExpression(ctx *ComparisonExpressionContext) {}

// ExitComparisonExpression is called when production comparisonExpression is exited.
func (s *BasePPLParserListener) ExitComparisonExpression(ctx *ComparisonExpressionContext) {}

// EnterAdditiveExpression is called when production additiveExpression is entered.
func (s *BasePPLParserListener) EnterAdditiveExpression(ctx *AdditiveExpressionContext) {}

// ExitAdditiveExpression is called when production additiveExpression is exited.
func (s *BasePPLParserListener) ExitAdditiveExpression(ctx *AdditiveExpressionContext) {}

// EnterMultiplicativeExpression is called when production multiplicativeExpression is entered.
func (s *BasePPLParserListener) EnterMultiplicativeExpression(ctx *MultiplicativeExpressionContext) {}

// ExitMultiplicativeExpression is called when production multiplicativeExpression is exited.
func (s *BasePPLParserListener) ExitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) {}

// EnterUnaryExpression is called when production unaryExpression is entered.
func (s *BasePPLParserListener) EnterUnaryExpression(ctx *UnaryExpressionContext) {}

// ExitUnaryExpression is called when production unaryExpression is exited.
func (s *BasePPLParserListener) ExitUnaryExpression(ctx *UnaryExpressionContext) {}

// EnterPrimaryExpression is called when production primaryExpression is entered.
func (s *BasePPLParserListener) EnterPrimaryExpression(ctx *PrimaryExpressionContext) {}

// ExitPrimaryExpression is called when production primaryExpression is exited.
func (s *BasePPLParserListener) ExitPrimaryExpression(ctx *PrimaryExpressionContext) {}

// EnterLiteral is called when production literal is entered.
func (s *BasePPLParserListener) EnterLiteral(ctx *LiteralContext) {}

// ExitLiteral is called when production literal is exited.
func (s *BasePPLParserListener) ExitLiteral(ctx *LiteralContext) {}

// EnterFieldReference is called when production fieldReference is entered.
func (s *BasePPLParserListener) EnterFieldReference(ctx *FieldReferenceContext) {}

// ExitFieldReference is called when production fieldReference is exited.
func (s *BasePPLParserListener) ExitFieldReference(ctx *FieldReferenceContext) {}

// EnterFunctionCallNoArgs is called when production FunctionCallNoArgs is entered.
func (s *BasePPLParserListener) EnterFunctionCallNoArgs(ctx *FunctionCallNoArgsContext) {}

// ExitFunctionCallNoArgs is called when production FunctionCallNoArgs is exited.
func (s *BasePPLParserListener) ExitFunctionCallNoArgs(ctx *FunctionCallNoArgsContext) {}

// EnterFunctionCallWithArgs is called when production FunctionCallWithArgs is entered.
func (s *BasePPLParserListener) EnterFunctionCallWithArgs(ctx *FunctionCallWithArgsContext) {}

// ExitFunctionCallWithArgs is called when production FunctionCallWithArgs is exited.
func (s *BasePPLParserListener) ExitFunctionCallWithArgs(ctx *FunctionCallWithArgsContext) {}

// EnterAggregationFunctionCallNoArgs is called when production AggregationFunctionCallNoArgs is entered.
func (s *BasePPLParserListener) EnterAggregationFunctionCallNoArgs(ctx *AggregationFunctionCallNoArgsContext) {
}

// ExitAggregationFunctionCallNoArgs is called when production AggregationFunctionCallNoArgs is exited.
func (s *BasePPLParserListener) ExitAggregationFunctionCallNoArgs(ctx *AggregationFunctionCallNoArgsContext) {
}

// EnterAggregationFunctionCall is called when production AggregationFunctionCall is entered.
func (s *BasePPLParserListener) EnterAggregationFunctionCall(ctx *AggregationFunctionCallContext) {}

// ExitAggregationFunctionCall is called when production AggregationFunctionCall is exited.
func (s *BasePPLParserListener) ExitAggregationFunctionCall(ctx *AggregationFunctionCallContext) {}

// EnterAggregationFunction is called when production aggregationFunction is entered.
func (s *BasePPLParserListener) EnterAggregationFunction(ctx *AggregationFunctionContext) {}

// ExitAggregationFunction is called when production aggregationFunction is exited.
func (s *BasePPLParserListener) ExitAggregationFunction(ctx *AggregationFunctionContext) {}

// EnterExpressionList is called when production expressionList is entered.
func (s *BasePPLParserListener) EnterExpressionList(ctx *ExpressionListContext) {}

// ExitExpressionList is called when production expressionList is exited.
func (s *BasePPLParserListener) ExitExpressionList(ctx *ExpressionListContext) {}

// EnterCaseExpression is called when production caseExpression is entered.
func (s *BasePPLParserListener) EnterCaseExpression(ctx *CaseExpressionContext) {}

// ExitCaseExpression is called when production caseExpression is exited.
func (s *BasePPLParserListener) ExitCaseExpression(ctx *CaseExpressionContext) {}

// EnterWhenClause is called when production whenClause is entered.
func (s *BasePPLParserListener) EnterWhenClause(ctx *WhenClauseContext) {}

// ExitWhenClause is called when production whenClause is exited.
func (s *BasePPLParserListener) ExitWhenClause(ctx *WhenClauseContext) {}
