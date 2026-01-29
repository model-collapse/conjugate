// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/parser/generated"
)

// ASTBuilder converts ANTLR4 parse trees to AST nodes
type ASTBuilder struct {
	generated.BasePPLParserVisitor
}

// NewASTBuilder creates a new AST builder
func NewASTBuilder() *ASTBuilder {
	return &ASTBuilder{}
}

// Helper function to get position from ANTLR4 context
func getPosition(ctx antlr.ParserRuleContext) ast.Position {
	if ctx == nil {
		return ast.NoPos
	}
	token := ctx.GetStart()
	return ast.NewPosition(
		token.GetLine(),
		token.GetColumn()+1, // ANTLR uses 0-indexed columns
		token.GetStart(),
	)
}

// VisitQuery builds a Query AST node
func (b *ASTBuilder) VisitQuery(ctx *generated.QueryContext) interface{} {
	if ctx == nil {
		return fmt.Errorf("nil query context")
	}

	commands := make([]ast.Command, 0)

	// Check if there's an EXPLAIN command wrapping the query
	hasExplain := ctx.ExplainCommand() != nil

	// Get the main commands
	if searchQueryCtx := ctx.SearchQuery(); searchQueryCtx != nil {
		// Search query: searchCommand followed by optional processing commands
		result := searchQueryCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		commands = result.([]ast.Command)
	} else if metadataCtx := ctx.MetadataCommand(); metadataCtx != nil {
		// Metadata command: single command
		result := metadataCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		commands = append(commands, result.(ast.Command))
	} else {
		return fmt.Errorf("query must be either a search query or metadata command")
	}

	// If EXPLAIN is present, wrap the first command with ExplainCommand
	if hasExplain {
		explainCmd := &ast.ExplainCommand{
			BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		}
		// Insert EXPLAIN as the first command
		commands = append([]ast.Command{explainCmd}, commands...)
	}

	return &ast.Query{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Commands: commands,
	}
}

// VisitSearchQuery builds command list from searchCommand + processingCommands
func (b *ASTBuilder) VisitSearchQuery(ctx *generated.SearchQueryContext) interface{} {
	if ctx == nil {
		return fmt.Errorf("nil search query context")
	}

	commands := make([]ast.Command, 0)

	// First command is always searchCommand
	searchCtx := ctx.SearchCommand()
	if searchCtx == nil {
		return fmt.Errorf("search query missing search command")
	}
	result := searchCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	commands = append(commands, result.(ast.Command))

	// Add processing commands
	for _, procCtx := range ctx.AllProcessingCommand() {
		result := procCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		commands = append(commands, result.(ast.Command))
	}

	return commands
}

// VisitProcessingCommand dispatches to specific processing command visitors
func (b *ASTBuilder) VisitProcessingCommand(ctx *generated.ProcessingCommandContext) interface{} {
	// Tier 0 commands
	if whereCtx := ctx.WhereCommand(); whereCtx != nil {
		return whereCtx.Accept(b)
	}
	if fieldsCtx := ctx.FieldsCommand(); fieldsCtx != nil {
		return fieldsCtx.Accept(b)
	}
	if statsCtx := ctx.StatsCommand(); statsCtx != nil {
		return statsCtx.Accept(b)
	}
	if sortCtx := ctx.SortCommand(); sortCtx != nil {
		return sortCtx.Accept(b)
	}
	if headCtx := ctx.HeadCommand(); headCtx != nil {
		return headCtx.Accept(b)
	}

	// Tier 1 commands
	if chartCtx := ctx.ChartCommand(); chartCtx != nil {
		return chartCtx.Accept(b)
	}
	if timechartCtx := ctx.TimechartCommand(); timechartCtx != nil {
		return timechartCtx.Accept(b)
	}
	if binCtx := ctx.BinCommand(); binCtx != nil {
		return binCtx.Accept(b)
	}
	if dedupCtx := ctx.DedupCommand(); dedupCtx != nil {
		return dedupCtx.Accept(b)
	}
	if topCtx := ctx.TopCommand(); topCtx != nil {
		return topCtx.Accept(b)
	}
	if rareCtx := ctx.RareCommand(); rareCtx != nil {
		return rareCtx.Accept(b)
	}
	if evalCtx := ctx.EvalCommand(); evalCtx != nil {
		return evalCtx.Accept(b)
	}
	if renameCtx := ctx.RenameCommand(); renameCtx != nil {
		return renameCtx.Accept(b)
	}
	if replaceCtx := ctx.ReplaceCommand(); replaceCtx != nil {
		return replaceCtx.Accept(b)
	}
	if fillnullCtx := ctx.FillnullCommand(); fillnullCtx != nil {
		return fillnullCtx.Accept(b)
	}
	if parseCtx := ctx.ParseCommand(); parseCtx != nil {
		return parseCtx.Accept(b)
	}
	if rexCtx := ctx.RexCommand(); rexCtx != nil {
		return rexCtx.Accept(b)
	}
	if lookupCtx := ctx.LookupCommand(); lookupCtx != nil {
		return lookupCtx.Accept(b)
	}
	if appendCtx := ctx.AppendCommand(); appendCtx != nil {
		return appendCtx.Accept(b)
	}
	if joinCtx := ctx.JoinCommand(); joinCtx != nil {
		return joinCtx.Accept(b)
	}
	if tableCtx := ctx.TableCommand(); tableCtx != nil {
		return tableCtx.Accept(b)
	}
	if eventstatsCtx := ctx.EventstatsCommand(); eventstatsCtx != nil {
		return eventstatsCtx.Accept(b)
	}
	if streamstatsCtx := ctx.StreamstatsCommand(); streamstatsCtx != nil {
		return streamstatsCtx.Accept(b)
	}

	return fmt.Errorf("unknown processing command type")
}

// VisitMetadataCommand dispatches to specific metadata command visitors
func (b *ASTBuilder) VisitMetadataCommand(ctx *generated.MetadataCommandContext) interface{} {
	if descCtx := ctx.DescribeCommand(); descCtx != nil {
		return descCtx.Accept(b)
	}
	if showCtx := ctx.ShowDatasourcesCommand(); showCtx != nil {
		return showCtx.Accept(b)
	}
	return fmt.Errorf("unknown metadata command type")
}

// VisitSearchWithSource builds a SearchCommand AST node from "source=logs"
func (b *ASTBuilder) VisitSearchWithSource(ctx *generated.SearchWithSourceContext) interface{} {
	// Get source name (the IDENTIFIER after =)
	idNode := ctx.IDENTIFIER()
	if idNode == nil {
		return fmt.Errorf("search command requires source parameter")
	}
	source := idNode.GetText()

	return &ast.SearchCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Source:   source,
	}
}

// VisitSearchWithKeyword builds a SearchCommand AST node from "search source=logs"
func (b *ASTBuilder) VisitSearchWithKeyword(ctx *generated.SearchWithKeywordContext) interface{} {
	// Get source name (the IDENTIFIER after =)
	idNode := ctx.IDENTIFIER()
	if idNode == nil {
		return fmt.Errorf("search command requires source parameter")
	}
	source := idNode.GetText()

	return &ast.SearchCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Source:   source,
	}
}

// VisitWhereCommand builds a WhereCommand AST node
func (b *ASTBuilder) VisitWhereCommand(ctx *generated.WhereCommandContext) interface{} {
	exprCtx := ctx.Expression()
	if exprCtx == nil {
		return fmt.Errorf("where command missing expression")
	}

	result := exprCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	return &ast.WhereCommand{
		BaseNode:  ast.BaseNode{Pos: getPosition(ctx)},
		Condition: result.(ast.Expression),
	}
}

// VisitFieldsInclude builds a FieldsCommand AST node (include mode)
func (b *ASTBuilder) VisitFieldsInclude(ctx *generated.FieldsIncludeContext) interface{} {
	fieldListCtx := ctx.FieldList()
	if fieldListCtx == nil {
		return fmt.Errorf("fields command missing field list")
	}

	result := fieldListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	fields := result.([]ast.Expression)
	return &ast.FieldsCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Fields:   fields,
		Includes: true,
	}
}

// VisitFieldsExclude builds a FieldsCommand AST node (exclude mode)
func (b *ASTBuilder) VisitFieldsExclude(ctx *generated.FieldsExcludeContext) interface{} {
	fieldListCtx := ctx.FieldList()
	if fieldListCtx == nil {
		return fmt.Errorf("fields command missing field list")
	}

	result := fieldListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	fields := result.([]ast.Expression)
	return &ast.FieldsCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Fields:   fields,
		Includes: false,
	}
}

// VisitFieldList processes a list of fields
func (b *ASTBuilder) VisitFieldList(ctx *generated.FieldListContext) interface{} {
	fields := make([]ast.Expression, 0)
	for _, exprCtx := range ctx.AllExpression() {
		result := exprCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		fields = append(fields, result.(ast.Expression))
	}
	return fields
}

// VisitStatsCommand builds a StatsCommand AST node
func (b *ASTBuilder) VisitStatsCommand(ctx *generated.StatsCommandContext) interface{} {
	// Process aggregations
	aggListCtx := ctx.AggregationList()
	if aggListCtx == nil {
		return fmt.Errorf("stats command missing aggregation list")
	}

	result := aggListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	aggregations := result.([]*ast.Aggregation)

	// Process group by (optional)
	var groupBy []ast.Expression
	if fieldListCtx := ctx.FieldList(); fieldListCtx != nil {
		result := fieldListCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		groupBy = result.([]ast.Expression)
	}

	return &ast.StatsCommand{
		BaseNode:     ast.BaseNode{Pos: getPosition(ctx)},
		Aggregations: aggregations,
		GroupBy:      groupBy,
	}
}

// VisitAggregationList processes a list of aggregations
func (b *ASTBuilder) VisitAggregationList(ctx *generated.AggregationListContext) interface{} {
	aggregations := make([]*ast.Aggregation, 0)
	for _, aggCtx := range ctx.AllAggregation() {
		result := aggCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		aggregations = append(aggregations, result.(*ast.Aggregation))
	}
	return aggregations
}

// VisitAggregation processes a single aggregation with optional alias
func (b *ASTBuilder) VisitAggregation(ctx *generated.AggregationContext) interface{} {
	// Get the expression (should be a function call)
	exprCtx := ctx.Expression()
	if exprCtx == nil {
		return fmt.Errorf("aggregation missing expression")
	}

	result := exprCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	expr := result.(ast.Expression)

	// Extract function call
	funcCall, ok := expr.(*ast.FunctionCall)
	if !ok {
		return fmt.Errorf("aggregation expression must be a function call")
	}

	// Get optional alias
	var alias string
	if idNode := ctx.IDENTIFIER(); idNode != nil {
		alias = idNode.GetText()
	}

	return &ast.Aggregation{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Func:     funcCall,
		Alias:    alias,
	}
}

// VisitSortCommand builds a SortCommand AST node
func (b *ASTBuilder) VisitSortCommand(ctx *generated.SortCommandContext) interface{} {
	sortFieldListCtx := ctx.SortFieldList()
	if sortFieldListCtx == nil {
		return fmt.Errorf("sort command missing sort field list")
	}

	result := sortFieldListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	return &ast.SortCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		SortKeys: result.([]*ast.SortKey),
	}
}

// VisitSortFieldList processes a list of sort fields
func (b *ASTBuilder) VisitSortFieldList(ctx *generated.SortFieldListContext) interface{} {
	sortKeys := make([]*ast.SortKey, 0)
	for _, sortFieldCtx := range ctx.AllSortField() {
		result := sortFieldCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		sortKeys = append(sortKeys, result.(*ast.SortKey))
	}
	return sortKeys
}

// VisitSortField processes a single sort field with optional order
func (b *ASTBuilder) VisitSortField(ctx *generated.SortFieldContext) interface{} {
	exprCtx := ctx.Expression()
	if exprCtx == nil {
		return fmt.Errorf("sort field missing expression")
	}

	result := exprCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	descending := ctx.DESC() != nil

	return &ast.SortKey{
		BaseNode:   ast.BaseNode{Pos: getPosition(ctx)},
		Field:      result.(ast.Expression),
		Descending: descending,
	}
}

// VisitHeadCommand builds a HeadCommand AST node
func (b *ASTBuilder) VisitHeadCommand(ctx *generated.HeadCommandContext) interface{} {
	intNode := ctx.INTEGER()
	if intNode == nil {
		return fmt.Errorf("head command missing count")
	}

	count, err := strconv.Atoi(intNode.GetText())
	if err != nil {
		return fmt.Errorf("invalid head count: %v", err)
	}

	return &ast.HeadCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Count:    count,
	}
}

// VisitDescribeCommand builds a DescribeCommand AST node
func (b *ASTBuilder) VisitDescribeCommand(ctx *generated.DescribeCommandContext) interface{} {
	idNode := ctx.IDENTIFIER()
	if idNode == nil {
		return fmt.Errorf("describe command missing source name")
	}

	return &ast.DescribeCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Source:   idNode.GetText(),
	}
}

// VisitShowDatasourcesCommand builds a ShowDatasourcesCommand AST node
func (b *ASTBuilder) VisitShowDatasourcesCommand(ctx *generated.ShowDatasourcesCommandContext) interface{} {
	return &ast.ShowDatasourcesCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
	}
}

// VisitExplainCommand builds an ExplainCommand AST node
func (b *ASTBuilder) VisitExplainCommand(ctx *generated.ExplainCommandContext) interface{} {
	return &ast.ExplainCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
	}
}

// VisitExpression processes an expression
func (b *ASTBuilder) VisitExpression(ctx *generated.ExpressionContext) interface{} {
	orExprCtx := ctx.OrExpression()
	if orExprCtx == nil {
		return fmt.Errorf("expression missing or-expression")
	}
	return orExprCtx.Accept(b)
}

// VisitOrExpression processes OR expressions
func (b *ASTBuilder) VisitOrExpression(ctx *generated.OrExpressionContext) interface{} {
	andExprs := ctx.AllAndExpression()
	if len(andExprs) == 0 {
		return fmt.Errorf("or-expression has no and-expressions")
	}

	// Start with the first and-expression
	result := andExprs[0].Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	left := result.(ast.Expression)

	// Chain OR operations
	for i := 1; i < len(andExprs); i++ {
		result := andExprs[i].Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		right := result.(ast.Expression)

		left = &ast.BinaryExpression{
			BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
			Left:     left,
			Operator: "OR",
			Right:    right,
		}
	}

	return left
}

// VisitAndExpression processes AND expressions
func (b *ASTBuilder) VisitAndExpression(ctx *generated.AndExpressionContext) interface{} {
	notExprs := ctx.AllNotExpression()
	if len(notExprs) == 0 {
		return fmt.Errorf("and-expression has no not-expressions")
	}

	// Start with the first not-expression
	result := notExprs[0].Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	left := result.(ast.Expression)

	// Chain AND operations
	for i := 1; i < len(notExprs); i++ {
		result := notExprs[i].Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		right := result.(ast.Expression)

		left = &ast.BinaryExpression{
			BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
			Left:     left,
			Operator: "AND",
			Right:    right,
		}
	}

	return left
}

// VisitNotExpression processes NOT expressions
func (b *ASTBuilder) VisitNotExpression(ctx *generated.NotExpressionContext) interface{} {
	// Check if this is a NOT expression
	if ctx.NOT() != nil {
		// Recursive NOT
		notExprCtx := ctx.NotExpression()
		if notExprCtx == nil {
			return fmt.Errorf("NOT missing operand")
		}

		result := notExprCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}

		return &ast.UnaryExpression{
			BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
			Operator: "NOT",
			Operand:  result.(ast.Expression),
		}
	}

	// Otherwise it's a comparison expression
	compCtx := ctx.ComparisonExpression()
	if compCtx == nil {
		return fmt.Errorf("not-expression missing comparison-expression")
	}
	return compCtx.Accept(b)
}

// VisitComparisonExpression processes comparison expressions
func (b *ASTBuilder) VisitComparisonExpression(ctx *generated.ComparisonExpressionContext) interface{} {
	addExprs := ctx.AllAdditiveExpression()
	if len(addExprs) == 0 {
		return fmt.Errorf("comparison-expression has no additive-expressions")
	}

	// Get left operand
	result := addExprs[0].Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	left := result.(ast.Expression)

	// Check for comparison operator
	if len(addExprs) > 1 {
		// Binary comparison
		result := addExprs[1].Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		right := result.(ast.Expression)

		var operator string
		if ctx.EQ() != nil {
			operator = "="
		} else if ctx.NEQ() != nil {
			operator = "!="
		} else if ctx.LT() != nil {
			operator = "<"
		} else if ctx.LTE() != nil {
			operator = "<="
		} else if ctx.GT() != nil {
			operator = ">"
		} else if ctx.GTE() != nil {
			operator = ">="
		} else if ctx.LIKE() != nil {
			operator = "LIKE"
		} else {
			return fmt.Errorf("unknown comparison operator")
		}

		return &ast.BinaryExpression{
			BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	// Check for IN operator
	if ctx.IN() != nil {
		exprListCtx := ctx.ExpressionList()
		if exprListCtx == nil {
			return fmt.Errorf("IN missing expression list")
		}

		result := exprListCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		values := result.([]ast.Expression)

		return &ast.BinaryExpression{
			BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
			Left:     left,
			Operator: "IN",
			Right: &ast.ListLiteral{
				BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
				Values:   values,
			},
		}
	}

	// No comparison, just return the additive expression
	return left
}

// VisitAdditiveExpression processes addition and subtraction
func (b *ASTBuilder) VisitAdditiveExpression(ctx *generated.AdditiveExpressionContext) interface{} {
	multExprs := ctx.AllMultiplicativeExpression()
	if len(multExprs) == 0 {
		return fmt.Errorf("additive-expression has no multiplicative-expressions")
	}

	// Start with the first multiplicative expression
	result := multExprs[0].Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	left := result.(ast.Expression)

	// Chain operations (need to parse operators from tokens)
	// This is simplified - in reality need to track operator positions
	for i := 1; i < len(multExprs); i++ {
		result := multExprs[i].Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		right := result.(ast.Expression)

		// Determine operator (simplified - assumes + for now)
		// TODO: Parse actual operator from token stream
		operator := "+"
		if i-1 < len(ctx.AllMINUS()) {
			operator = "-"
		}

		left = &ast.BinaryExpression{
			BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

// VisitMultiplicativeExpression processes multiplication, division, modulo
func (b *ASTBuilder) VisitMultiplicativeExpression(ctx *generated.MultiplicativeExpressionContext) interface{} {
	unaryExprs := ctx.AllUnaryExpression()
	if len(unaryExprs) == 0 {
		return fmt.Errorf("multiplicative-expression has no unary-expressions")
	}

	// Start with the first unary expression
	result := unaryExprs[0].Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	left := result.(ast.Expression)

	// Chain operations
	for i := 1; i < len(unaryExprs); i++ {
		result := unaryExprs[i].Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		right := result.(ast.Expression)

		// Determine operator (simplified)
		operator := "*"
		// TODO: Parse actual operator from token stream

		left = &ast.BinaryExpression{
			BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

// VisitUnaryExpression processes unary + and -
func (b *ASTBuilder) VisitUnaryExpression(ctx *generated.UnaryExpressionContext) interface{} {
	// Check for unary operator
	if ctx.PLUS() != nil || ctx.MINUS() != nil {
		unaryCtx := ctx.UnaryExpression()
		if unaryCtx == nil {
			return fmt.Errorf("unary operator missing operand")
		}

		result := unaryCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}

		operator := "+"
		if ctx.MINUS() != nil {
			operator = "-"
		}

		return &ast.UnaryExpression{
			BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
			Operator: operator,
			Operand:  result.(ast.Expression),
		}
	}

	// Otherwise it's a primary expression
	primaryCtx := ctx.PrimaryExpression()
	if primaryCtx == nil {
		return fmt.Errorf("unary-expression missing primary-expression")
	}
	return primaryCtx.Accept(b)
}

// VisitPrimaryExpression processes primary expressions
func (b *ASTBuilder) VisitPrimaryExpression(ctx *generated.PrimaryExpressionContext) interface{} {
	// Check each alternative
	if litCtx := ctx.Literal(); litCtx != nil {
		return litCtx.Accept(b)
	}
	if fieldCtx := ctx.FieldReference(); fieldCtx != nil {
		return fieldCtx.Accept(b)
	}
	if funcCtx := ctx.FunctionCall(); funcCtx != nil {
		return funcCtx.Accept(b)
	}
	if caseCtx := ctx.CaseExpression(); caseCtx != nil {
		return caseCtx.Accept(b)
	}
	if exprCtx := ctx.Expression(); exprCtx != nil {
		// Parenthesized expression
		return exprCtx.Accept(b)
	}
	return fmt.Errorf("unknown primary expression")
}

// VisitLiteral processes literal values
func (b *ASTBuilder) VisitLiteral(ctx *generated.LiteralContext) interface{} {
	if intNode := ctx.INTEGER(); intNode != nil {
		val, err := strconv.Atoi(intNode.GetText())
		if err != nil {
			return err
		}
		return &ast.Literal{
			BaseNode:   ast.BaseNode{Pos: getPosition(ctx)},
			Value:      val,
			LiteralTyp: ast.LiteralTypeInt,
		}
	}

	if decNode := ctx.DECIMAL(); decNode != nil {
		val, err := strconv.ParseFloat(decNode.GetText(), 64)
		if err != nil {
			return err
		}
		return &ast.Literal{
			BaseNode:   ast.BaseNode{Pos: getPosition(ctx)},
			Value:      val,
			LiteralTyp: ast.LiteralTypeFloat,
		}
	}

	if strNode := ctx.STRING(); strNode != nil {
		// Remove quotes
		text := strNode.GetText()
		if len(text) >= 2 {
			text = text[1 : len(text)-1]
		}
		// Handle escaped quotes
		text = strings.ReplaceAll(text, "''", "'")
		text = strings.ReplaceAll(text, "\"\"", "\"")

		return &ast.Literal{
			BaseNode:   ast.BaseNode{Pos: getPosition(ctx)},
			Value:      text,
			LiteralTyp: ast.LiteralTypeString,
		}
	}

	if ctx.TRUE() != nil {
		return &ast.Literal{
			BaseNode:   ast.BaseNode{Pos: getPosition(ctx)},
			Value:      true,
			LiteralTyp: ast.LiteralTypeBool,
		}
	}

	if ctx.FALSE() != nil {
		return &ast.Literal{
			BaseNode:   ast.BaseNode{Pos: getPosition(ctx)},
			Value:      false,
			LiteralTyp: ast.LiteralTypeBool,
		}
	}

	if ctx.NULL() != nil {
		return &ast.Literal{
			BaseNode:   ast.BaseNode{Pos: getPosition(ctx)},
			Value:      nil,
			LiteralTyp: ast.LiteralTypeNull,
		}
	}

	return fmt.Errorf("unknown literal type")
}

// VisitFieldReference processes field references
func (b *ASTBuilder) VisitFieldReference(ctx *generated.FieldReferenceContext) interface{} {
	// Build field name from identifiers (may include dots for nested fields)
	identifiers := ctx.AllIDENTIFIER()
	if len(identifiers) == 0 {
		return fmt.Errorf("field reference has no identifiers")
	}

	parts := make([]string, len(identifiers))
	for i, id := range identifiers {
		parts[i] = id.GetText()
	}
	fieldName := strings.Join(parts, ".")

	// Check for array indexing
	if intNode := ctx.INTEGER(); intNode != nil {
		fieldName += "[" + intNode.GetText() + "]"
	}

	return &ast.FieldReference{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Name:     fieldName,
	}
}

// VisitFunctionCallNoArgs processes function calls with no arguments
func (b *ASTBuilder) VisitFunctionCallNoArgs(ctx *generated.FunctionCallNoArgsContext) interface{} {
	idNode := ctx.IDENTIFIER()
	if idNode == nil {
		return fmt.Errorf("function call missing name")
	}

	return &ast.FunctionCall{
		BaseNode:  ast.BaseNode{Pos: getPosition(ctx)},
		Name:      idNode.GetText(),
		Arguments: []ast.Expression{},
		Distinct:  false,
	}
}

// VisitFunctionCallWithArgs processes function calls with arguments
func (b *ASTBuilder) VisitFunctionCallWithArgs(ctx *generated.FunctionCallWithArgsContext) interface{} {
	idNode := ctx.IDENTIFIER()
	if idNode == nil {
		return fmt.Errorf("function call missing name")
	}

	exprListCtx := ctx.ExpressionList()
	if exprListCtx == nil {
		return fmt.Errorf("function call missing expression list")
	}

	result := exprListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	args := result.([]ast.Expression)

	distinct := ctx.DISTINCT() != nil

	return &ast.FunctionCall{
		BaseNode:  ast.BaseNode{Pos: getPosition(ctx)},
		Name:      idNode.GetText(),
		Arguments: args,
		Distinct:  distinct,
	}
}

// VisitAggregationFunctionCallNoArgs processes aggregation function calls with no arguments (e.g., count())
func (b *ASTBuilder) VisitAggregationFunctionCallNoArgs(ctx *generated.AggregationFunctionCallNoArgsContext) interface{} {
	aggFuncCtx := ctx.AggregationFunction()
	if aggFuncCtx == nil {
		return fmt.Errorf("aggregation function call missing function")
	}

	// Get function name
	funcName := aggFuncCtx.GetText()

	return &ast.FunctionCall{
		BaseNode:  ast.BaseNode{Pos: getPosition(ctx)},
		Name:      funcName,
		Arguments: []ast.Expression{}, // No arguments
		Distinct:  false,
	}
}

// VisitAggregationFunctionCall processes aggregation function calls with arguments
func (b *ASTBuilder) VisitAggregationFunctionCall(ctx *generated.AggregationFunctionCallContext) interface{} {
	aggFuncCtx := ctx.AggregationFunction()
	if aggFuncCtx == nil {
		return fmt.Errorf("aggregation function call missing function")
	}

	// Get function name
	funcName := aggFuncCtx.GetText()

	exprListCtx := ctx.ExpressionList()
	if exprListCtx == nil {
		return fmt.Errorf("aggregation function call missing expression list")
	}

	result := exprListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	args := result.([]ast.Expression)

	distinct := ctx.DISTINCT() != nil

	return &ast.FunctionCall{
		BaseNode:  ast.BaseNode{Pos: getPosition(ctx)},
		Name:      funcName,
		Arguments: args,
		Distinct:  distinct,
	}
}

// VisitExpressionList processes a list of expressions
func (b *ASTBuilder) VisitExpressionList(ctx *generated.ExpressionListContext) interface{} {
	expressions := make([]ast.Expression, 0)
	for _, exprCtx := range ctx.AllExpression() {
		result := exprCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		expressions = append(expressions, result.(ast.Expression))
	}
	return expressions
}

// VisitCaseExpression processes CASE expressions
func (b *ASTBuilder) VisitCaseExpression(ctx *generated.CaseExpressionContext) interface{} {
	whenClauses := make([]*ast.WhenClause, 0)
	for _, whenCtx := range ctx.AllWhenClause() {
		result := whenCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		whenClauses = append(whenClauses, result.(*ast.WhenClause))
	}

	// Get optional ELSE clause
	var elseResult ast.Expression
	if exprCtx := ctx.Expression(); exprCtx != nil {
		result := exprCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		elseResult = result.(ast.Expression)
	}

	return &ast.CaseExpression{
		BaseNode:    ast.BaseNode{Pos: getPosition(ctx)},
		WhenClauses: whenClauses,
		ElseResult:  elseResult,
	}
}

// VisitWhenClause processes WHEN clauses
func (b *ASTBuilder) VisitWhenClause(ctx *generated.WhenClauseContext) interface{} {
	allExprs := ctx.AllExpression()
	if len(allExprs) != 2 {
		return fmt.Errorf("WHEN clause must have exactly 2 expressions (condition and result)")
	}

	// Condition
	result := allExprs[0].Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	condition := result.(ast.Expression)

	// Result
	result = allExprs[1].Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	resultExpr := result.(ast.Expression)

	return &ast.WhenClause{
		BaseNode:  ast.BaseNode{Pos: getPosition(ctx)},
		Condition: condition,
		Result:    resultExpr,
	}
}

// ============================================================================
// Tier 1 Command Visitors
// ============================================================================

// VisitChartCommand builds a ChartCommand AST node
func (b *ASTBuilder) VisitChartCommand(ctx *generated.ChartCommandContext) interface{} {
	// Process aggregations
	aggListCtx := ctx.AggregationList()
	if aggListCtx == nil {
		return fmt.Errorf("chart command missing aggregation list")
	}

	result := aggListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	aggregations := result.([]*ast.Aggregation)

	// Process group by (optional)
	var groupBy []ast.Expression
	if fieldListCtx := ctx.FieldList(); fieldListCtx != nil {
		result := fieldListCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		groupBy = result.([]ast.Expression)
	}

	cmd := &ast.ChartCommand{
		BaseNode:     ast.BaseNode{Pos: getPosition(ctx)},
		Aggregations: aggregations,
		GroupBy:      groupBy,
	}

	// Process options
	for _, optCtx := range ctx.AllChartOptions() {
		if opt, ok := optCtx.(*generated.ChartOptionsContext); ok {
			b.applyChartOption(cmd, opt)
		}
	}

	return cmd
}

// applyChartOption applies chart options to the command
func (b *ASTBuilder) applyChartOption(cmd *ast.ChartCommand, ctx *generated.ChartOptionsContext) {
	if spanCtx := ctx.TimeSpan(); spanCtx != nil {
		if span, ok := spanCtx.(*generated.TimeSpanContext); ok {
			cmd.Span = b.visitTimeSpan(span)
		}
	}
	if ctx.LIMIT() != nil {
		if intNode := ctx.INTEGER(); intNode != nil {
			limit, _ := strconv.Atoi(intNode.GetText())
			cmd.Limit = limit
		}
	}
	if ctx.USEOTHER() != nil {
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.UseOther = boolCtx.TRUE() != nil
		}
	}
	if ctx.OTHERSTR() != nil {
		if strNode := ctx.STRING(); strNode != nil {
			cmd.OtherStr = unquoteString(strNode.GetText())
		}
	}
	if ctx.NULLSTR() != nil {
		if strNode := ctx.STRING(); strNode != nil {
			cmd.NullStr = unquoteString(strNode.GetText())
		}
	}
}

// VisitTimechartCommand builds a TimechartCommand AST node
func (b *ASTBuilder) VisitTimechartCommand(ctx *generated.TimechartCommandContext) interface{} {
	// Process aggregations
	aggListCtx := ctx.AggregationList()
	if aggListCtx == nil {
		return fmt.Errorf("timechart command missing aggregation list")
	}

	result := aggListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	aggregations := result.([]*ast.Aggregation)

	// Process group by (optional)
	var groupBy []ast.Expression
	if fieldListCtx := ctx.FieldList(); fieldListCtx != nil {
		result := fieldListCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		groupBy = result.([]ast.Expression)
	}

	cmd := &ast.TimechartCommand{
		BaseNode:     ast.BaseNode{Pos: getPosition(ctx)},
		Aggregations: aggregations,
		GroupBy:      groupBy,
	}

	// Process options
	for _, optCtx := range ctx.AllTimechartOptions() {
		if opt, ok := optCtx.(*generated.TimechartOptionsContext); ok {
			b.applyTimechartOption(cmd, opt)
		}
	}

	return cmd
}

// applyTimechartOption applies timechart options to the command
func (b *ASTBuilder) applyTimechartOption(cmd *ast.TimechartCommand, ctx *generated.TimechartOptionsContext) {
	if spanCtx := ctx.TimeSpan(); spanCtx != nil {
		if span, ok := spanCtx.(*generated.TimeSpanContext); ok {
			cmd.Span = b.visitTimeSpan(span)
		}
	}
	if ctx.BINS() != nil {
		if intNode := ctx.INTEGER(); intNode != nil {
			bins, _ := strconv.Atoi(intNode.GetText())
			cmd.Bins = bins
		}
	}
	if ctx.LIMIT() != nil {
		if intNode := ctx.INTEGER(); intNode != nil {
			limit, _ := strconv.Atoi(intNode.GetText())
			cmd.Limit = limit
		}
	}
	if ctx.USEOTHER() != nil {
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.UseOther = boolCtx.TRUE() != nil
		}
	}
}

// visitTimeSpan parses a time span like "1h", "30m", "1d"
func (b *ASTBuilder) visitTimeSpan(ctx *generated.TimeSpanContext) *ast.TimeSpan {
	if ctx == nil {
		return nil
	}

	intNode := ctx.INTEGER()
	idNode := ctx.IDENTIFIER()

	if intNode == nil || idNode == nil {
		return nil
	}

	value, _ := strconv.Atoi(intNode.GetText())
	unit := idNode.GetText()

	return &ast.TimeSpan{
		Value: value,
		Unit:  unit,
	}
}

// VisitBinCommand builds a BinCommand AST node
func (b *ASTBuilder) VisitBinCommand(ctx *generated.BinCommandContext) interface{} {
	// Get field reference
	fieldCtx := ctx.FieldReference()
	if fieldCtx == nil {
		return fmt.Errorf("bin command missing field")
	}

	result := fieldCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	cmd := &ast.BinCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Field:    result.(ast.Expression),
	}

	// Process options
	for _, optCtx := range ctx.AllBinOptions() {
		if opt, ok := optCtx.(*generated.BinOptionsContext); ok {
			b.applyBinOption(cmd, opt)
		}
	}

	return cmd
}

// applyBinOption applies bin options to the command
func (b *ASTBuilder) applyBinOption(cmd *ast.BinCommand, ctx *generated.BinOptionsContext) {
	if ctx.SPAN() != nil {
		if spanCtx := ctx.TimeSpan(); spanCtx != nil {
			if span, ok := spanCtx.(*generated.TimeSpanContext); ok {
				cmd.Span = b.visitTimeSpan(span)
			}
		} else if idNode := ctx.IDENTIFIER(); idNode != nil {
			// Check for "auto"
			if strings.ToLower(idNode.GetText()) == "auto" {
				cmd.Auto = true
			}
		}
	}
	if ctx.BINS() != nil {
		if intNode := ctx.INTEGER(); intNode != nil {
			bins, _ := strconv.Atoi(intNode.GetText())
			cmd.Bins = bins
		}
	}
}

// VisitDedupCommand builds a DedupCommand AST node
func (b *ASTBuilder) VisitDedupCommand(ctx *generated.DedupCommandContext) interface{} {
	cmd := &ast.DedupCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Count:    1, // Default to keep first 1
	}

	// Check for optional count
	if intNode := ctx.INTEGER(); intNode != nil {
		count, _ := strconv.Atoi(intNode.GetText())
		cmd.Count = count
	}

	// Get field list
	fieldListCtx := ctx.FieldList()
	if fieldListCtx == nil {
		return fmt.Errorf("dedup command missing field list")
	}

	result := fieldListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	cmd.Fields = result.([]ast.Expression)

	// Process options
	for _, optCtx := range ctx.AllDedupOptions() {
		if opt, ok := optCtx.(*generated.DedupOptionsContext); ok {
			b.applyDedupOption(cmd, opt)
		}
	}

	return cmd
}

// applyDedupOption applies dedup options to the command
func (b *ASTBuilder) applyDedupOption(cmd *ast.DedupCommand, ctx *generated.DedupOptionsContext) {
	if ctx.KEEPEVENTS() != nil {
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.KeepEvents = boolCtx.TRUE() != nil
		}
	}
	if ctx.CONSECUTIVE() != nil {
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.Consecutive = boolCtx.TRUE() != nil
		}
	}
	if sortListCtx := ctx.SortFieldList(); sortListCtx != nil {
		result := sortListCtx.Accept(b)
		if sortKeys, ok := result.([]*ast.SortKey); ok {
			cmd.SortBy = sortKeys
		}
	}
}

// VisitTopCommand builds a TopCommand AST node
func (b *ASTBuilder) VisitTopCommand(ctx *generated.TopCommandContext) interface{} {
	cmd := &ast.TopCommand{
		BaseNode:    ast.BaseNode{Pos: getPosition(ctx)},
		Limit:       10, // Default to 10
		ShowCount:   true,
		ShowPercent: true,
	}

	// Check for optional limit
	if intNode := ctx.INTEGER(); intNode != nil {
		limit, _ := strconv.Atoi(intNode.GetText())
		cmd.Limit = limit
	}

	// Get field list(s)
	fieldLists := ctx.AllFieldList()
	if len(fieldLists) == 0 {
		return fmt.Errorf("top command missing field list")
	}

	// First field list is the fields to get top values for
	result := fieldLists[0].Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	cmd.Fields = result.([]ast.Expression)

	// Second field list is optional group by
	if len(fieldLists) > 1 {
		result := fieldLists[1].Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		cmd.GroupBy = result.([]ast.Expression)
	}

	// Process options
	for _, optCtx := range ctx.AllTopOptions() {
		if opt, ok := optCtx.(*generated.TopOptionsContext); ok {
			b.applyTopOption(cmd, opt)
		}
	}

	return cmd
}

// applyTopOption applies top options to the command
func (b *ASTBuilder) applyTopOption(cmd *ast.TopCommand, ctx *generated.TopOptionsContext) {
	if ctx.COUNTFIELD() != nil {
		if strNode := ctx.STRING(); strNode != nil {
			cmd.CountField = unquoteString(strNode.GetText())
		}
	}
	if ctx.PERCENTFIELD() != nil {
		if strNode := ctx.STRING(); strNode != nil {
			cmd.PercentField = unquoteString(strNode.GetText())
		}
	}
	if ctx.SHOWCOUNT() != nil {
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.ShowCount = boolCtx.TRUE() != nil
		}
	}
	if ctx.SHOWPERC() != nil {
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.ShowPercent = boolCtx.TRUE() != nil
		}
	}
	if ctx.LIMIT() != nil {
		if intNode := ctx.INTEGER(); intNode != nil {
			limit, _ := strconv.Atoi(intNode.GetText())
			cmd.Limit = limit
		}
	}
	if ctx.USEOTHER() != nil {
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.UseOther = boolCtx.TRUE() != nil
		}
	}
	if ctx.OTHERSTR() != nil {
		if strNode := ctx.STRING(); strNode != nil {
			cmd.OtherStr = unquoteString(strNode.GetText())
		}
	}
}

// VisitRareCommand builds a RareCommand AST node
func (b *ASTBuilder) VisitRareCommand(ctx *generated.RareCommandContext) interface{} {
	cmd := &ast.RareCommand{
		BaseNode:    ast.BaseNode{Pos: getPosition(ctx)},
		Limit:       10, // Default to 10
		ShowCount:   true,
		ShowPercent: true,
	}

	// Check for optional limit
	if intNode := ctx.INTEGER(); intNode != nil {
		limit, _ := strconv.Atoi(intNode.GetText())
		cmd.Limit = limit
	}

	// Get field list(s)
	fieldLists := ctx.AllFieldList()
	if len(fieldLists) == 0 {
		return fmt.Errorf("rare command missing field list")
	}

	// First field list is the fields to get rare values for
	result := fieldLists[0].Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	cmd.Fields = result.([]ast.Expression)

	// Second field list is optional group by
	if len(fieldLists) > 1 {
		result := fieldLists[1].Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		cmd.GroupBy = result.([]ast.Expression)
	}

	// Process options (reuses TopOptions)
	for _, optCtx := range ctx.AllTopOptions() {
		if opt, ok := optCtx.(*generated.TopOptionsContext); ok {
			b.applyRareOption(cmd, opt)
		}
	}

	return cmd
}

// applyRareOption applies rare options to the command
func (b *ASTBuilder) applyRareOption(cmd *ast.RareCommand, ctx *generated.TopOptionsContext) {
	if ctx.COUNTFIELD() != nil {
		if strNode := ctx.STRING(); strNode != nil {
			cmd.CountField = unquoteString(strNode.GetText())
		}
	}
	if ctx.PERCENTFIELD() != nil {
		if strNode := ctx.STRING(); strNode != nil {
			cmd.PercentField = unquoteString(strNode.GetText())
		}
	}
	if ctx.SHOWCOUNT() != nil {
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.ShowCount = boolCtx.TRUE() != nil
		}
	}
	if ctx.SHOWPERC() != nil {
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.ShowPercent = boolCtx.TRUE() != nil
		}
	}
	if ctx.LIMIT() != nil {
		if intNode := ctx.INTEGER(); intNode != nil {
			limit, _ := strconv.Atoi(intNode.GetText())
			cmd.Limit = limit
		}
	}
}

// VisitEvalCommand builds an EvalCommand AST node
func (b *ASTBuilder) VisitEvalCommand(ctx *generated.EvalCommandContext) interface{} {
	assignments := make([]*ast.EvalAssignment, 0)

	for _, assignCtx := range ctx.AllEvalAssignment() {
		result := assignCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		assignments = append(assignments, result.(*ast.EvalAssignment))
	}

	return &ast.EvalCommand{
		BaseNode:    ast.BaseNode{Pos: getPosition(ctx)},
		Assignments: assignments,
	}
}

// VisitEvalAssignment builds an EvalAssignment AST node
func (b *ASTBuilder) VisitEvalAssignment(ctx *generated.EvalAssignmentContext) interface{} {
	idNode := ctx.IDENTIFIER()
	if idNode == nil {
		return fmt.Errorf("eval assignment missing field name")
	}

	exprCtx := ctx.Expression()
	if exprCtx == nil {
		return fmt.Errorf("eval assignment missing expression")
	}

	result := exprCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	return &ast.EvalAssignment{
		BaseNode:   ast.BaseNode{Pos: getPosition(ctx)},
		Field:      idNode.GetText(),
		Expression: result.(ast.Expression),
	}
}

// VisitRenameCommand builds a RenameCommand AST node
func (b *ASTBuilder) VisitRenameCommand(ctx *generated.RenameCommandContext) interface{} {
	assignments := make([]*ast.RenameAssignment, 0)

	for _, assignCtx := range ctx.AllRenameAssignment() {
		result := assignCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		assignments = append(assignments, result.(*ast.RenameAssignment))
	}

	return &ast.RenameCommand{
		BaseNode:    ast.BaseNode{Pos: getPosition(ctx)},
		Assignments: assignments,
	}
}

// VisitRenameAssignment builds a RenameAssignment AST node
func (b *ASTBuilder) VisitRenameAssignment(ctx *generated.RenameAssignmentContext) interface{} {
	ids := ctx.AllIDENTIFIER()
	if len(ids) != 2 {
		return fmt.Errorf("rename assignment requires old and new field names")
	}

	return &ast.RenameAssignment{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		OldName:  ids[0].GetText(),
		NewName:  ids[1].GetText(),
	}
}

// VisitReplaceCommand builds a ReplaceCommand AST node
func (b *ASTBuilder) VisitReplaceCommand(ctx *generated.ReplaceCommandContext) interface{} {
	mappings := make([]*ast.ReplaceMapping, 0)

	for _, mappingCtx := range ctx.AllReplaceMapping() {
		result := mappingCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		mappings = append(mappings, result.(*ast.ReplaceMapping))
	}

	// Get the target field name
	idCtx := ctx.IDENTIFIER()
	if idCtx == nil {
		return fmt.Errorf("replace command requires a field name after 'in'")
	}

	return &ast.ReplaceCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Mappings: mappings,
		Field:    idCtx.GetText(),
	}
}

// VisitReplaceMapping builds a ReplaceMapping AST node
func (b *ASTBuilder) VisitReplaceMapping(ctx *generated.ReplaceMappingContext) interface{} {
	exprs := ctx.AllExpression()
	if len(exprs) != 2 {
		return fmt.Errorf("replace mapping requires old and new values")
	}

	oldResult := exprs[0].Accept(b)
	if err, ok := oldResult.(error); ok {
		return err
	}

	newResult := exprs[1].Accept(b)
	if err, ok := newResult.(error); ok {
		return err
	}

	return &ast.ReplaceMapping{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		OldValue: oldResult.(ast.Expression),
		NewValue: newResult.(ast.Expression),
	}
}

// VisitFillnullCommand builds a FillnullCommand AST node
func (b *ASTBuilder) VisitFillnullCommand(ctx *generated.FillnullCommandContext) interface{} {
	assignments := make([]*ast.FillnullAssignment, 0)

	for _, assignCtx := range ctx.AllFillnullAssignment() {
		result := assignCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		assignments = append(assignments, result.(*ast.FillnullAssignment))
	}

	// Check if this is "value=<default>" syntax
	var defaultValue ast.Expression
	var fields []ast.Expression

	// Look for "value" assignment which means default value for all/specified fields
	for i, assignment := range assignments {
		if assignment.Field == "value" {
			defaultValue = assignment.Value
			// Remove this assignment as it's not a field assignment
			assignments = append(assignments[:i], assignments[i+1:]...)
			break
		}
	}

	// Look for "fields" pseudo-assignment (this would be handled differently in practice)
	// For now, if we have a default value but no other assignments, apply to all fields
	if defaultValue != nil && len(assignments) == 0 {
		// Default value applies to all fields
		return &ast.FillnullCommand{
			BaseNode:     ast.BaseNode{Pos: getPosition(ctx)},
			DefaultValue: defaultValue,
			Fields:       fields,
		}
	}

	return &ast.FillnullCommand{
		BaseNode:    ast.BaseNode{Pos: getPosition(ctx)},
		Assignments: assignments,
	}
}

// VisitFillnullAssignment builds a FillnullAssignment AST node
func (b *ASTBuilder) VisitFillnullAssignment(ctx *generated.FillnullAssignmentContext) interface{} {
	idCtx := ctx.IDENTIFIER()
	if idCtx == nil {
		return fmt.Errorf("fillnull assignment requires a field name")
	}

	exprCtx := ctx.Expression()
	if exprCtx == nil {
		return fmt.Errorf("fillnull assignment requires a value expression")
	}

	result := exprCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	return &ast.FillnullAssignment{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Field:    idCtx.GetText(),
		Value:    result.(ast.Expression),
	}
}

// VisitParseCommand builds a ParseCommand AST node
// Syntax: parse [field=]<source_field> "<pattern>"
func (b *ASTBuilder) VisitParseCommand(ctx *generated.ParseCommandContext) interface{} {
	var sourceField string
	var fieldParam string
	var pattern string

	// Get identifiers (could be 1 or 2 depending on syntax)
	identifiers := ctx.AllIDENTIFIER()

	if len(identifiers) == 2 {
		// Syntax: parse field=message "pattern"
		fieldParam = identifiers[0].GetText()
		sourceField = identifiers[1].GetText()
	} else if len(identifiers) == 1 {
		// Syntax: parse message "pattern"
		sourceField = identifiers[0].GetText()
	}

	// Get pattern string (remove quotes)
	if ctx.STRING() != nil {
		pattern = unquoteString(ctx.STRING().GetText())
	}

	return &ast.ParseCommand{
		BaseNode:    ast.BaseNode{Pos: getPosition(ctx)},
		SourceField: sourceField,
		Pattern:     pattern,
		FieldParam:  fieldParam,
	}
}

// VisitRexCommand builds a RexCommand AST node
// Syntax: rex [field=<source_field>] "<pattern>"
func (b *ASTBuilder) VisitRexCommand(ctx *generated.RexCommandContext) interface{} {
	var sourceField string
	var fieldParam string
	var pattern string

	// Get identifiers (could be 0 or 2 depending on syntax)
	identifiers := ctx.AllIDENTIFIER()

	if len(identifiers) == 2 {
		// Syntax: rex field=message "pattern"
		fieldParam = identifiers[0].GetText()
		sourceField = identifiers[1].GetText()
	}
	// If len(identifiers) == 0: rex "pattern" - sourceField will remain empty

	// Get pattern string (remove quotes)
	if ctx.STRING() != nil {
		pattern = unquoteString(ctx.STRING().GetText())
	}

	return &ast.RexCommand{
		BaseNode:    ast.BaseNode{Pos: getPosition(ctx)},
		SourceField: sourceField,
		Pattern:     pattern,
		FieldParam:  fieldParam,
	}
}

// VisitLookupCommand builds a LookupCommand AST node
// Syntax: lookup <table> <join_field> [AS <alias>] OUTPUT <output_fields>
func (b *ASTBuilder) VisitLookupCommand(ctx *generated.LookupCommandContext) interface{} {
	// Get identifiers: table name, join field, optional alias
	identifiers := ctx.AllIDENTIFIER()
	if len(identifiers) < 2 {
		return fmt.Errorf("lookup command requires at least table name and join field")
	}

	tableName := identifiers[0].GetText()
	joinField := identifiers[1].GetText()
	var joinFieldAlias string

	// Check for AS clause (if len == 3, we have: table, joinField, alias)
	if len(identifiers) == 3 {
		joinFieldAlias = identifiers[2].GetText()
	}

	// Get output fields
	outputListCtx := ctx.LookupOutputList()
	if outputListCtx == nil {
		return fmt.Errorf("lookup command requires OUTPUT clause")
	}

	result := outputListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	outputFields, ok := result.([]*ast.LookupOutputField)
	if !ok {
		return fmt.Errorf("invalid output fields for lookup command")
	}

	return &ast.LookupCommand{
		BaseNode:       ast.BaseNode{Pos: getPosition(ctx)},
		TableName:      tableName,
		JoinField:      joinField,
		JoinFieldAlias: joinFieldAlias,
		OutputFields:   outputFields,
	}
}

// VisitLookupOutputList builds the list of output fields
func (b *ASTBuilder) VisitLookupOutputList(ctx *generated.LookupOutputListContext) interface{} {
	outputFieldCtxs := ctx.AllLookupOutputField()
	if len(outputFieldCtxs) == 0 {
		return fmt.Errorf("lookup command requires at least one output field")
	}

	outputFields := make([]*ast.LookupOutputField, 0, len(outputFieldCtxs))
	for _, fieldCtx := range outputFieldCtxs {
		result := fieldCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		if field, ok := result.(*ast.LookupOutputField); ok {
			outputFields = append(outputFields, field)
		}
	}

	return outputFields
}

// VisitLookupOutputField builds a single output field
func (b *ASTBuilder) VisitLookupOutputField(ctx *generated.LookupOutputFieldContext) interface{} {
	identifiers := ctx.AllIDENTIFIER()
	if len(identifiers) == 0 {
		return fmt.Errorf("lookup output field requires a field name")
	}

	field := identifiers[0].GetText()
	var alias string

	// Check for AS clause
	if len(identifiers) == 2 {
		alias = identifiers[1].GetText()
	}

	return &ast.LookupOutputField{
		Field: field,
		Alias: alias,
	}
}

// VisitAppendCommand builds an AppendCommand AST node
// Syntax: append [subsearch]
func (b *ASTBuilder) VisitAppendCommand(ctx *generated.AppendCommandContext) interface{} {
	// Get the subsearch query context
	subsearchCtx := ctx.SearchQuery()
	if subsearchCtx == nil {
		return fmt.Errorf("append command requires a subsearch query")
	}

	// Visit the subsearch to build its AST
	result := subsearchCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	// VisitSearchQuery returns []ast.Command, so we need to wrap it in a Query
	commands, ok := result.([]ast.Command)
	if !ok {
		return fmt.Errorf("invalid subsearch in append command")
	}

	subsearch := &ast.Query{
		BaseNode: ast.BaseNode{Pos: getPosition(subsearchCtx)},
		Commands: commands,
	}

	return &ast.AppendCommand{
		BaseNode:  ast.BaseNode{Pos: getPosition(ctx)},
		Subsearch: subsearch,
	}
}

// VisitJoinCommand builds a JoinCommand AST node
// Syntax: join [type=TYPE] field [subsearch]
func (b *ASTBuilder) VisitJoinCommand(ctx *generated.JoinCommandContext) interface{} {
	// Default join type is inner
	joinType := ast.JoinTypeInner

	// Check for explicit join type
	if joinTypeCtx := ctx.JoinType(); joinTypeCtx != nil {
		result := joinTypeCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		joinType = result.(ast.JoinType)
	}

	// Get join field name
	identifier := ctx.IDENTIFIER()
	if identifier == nil {
		return fmt.Errorf("join command requires a join field")
	}
	joinField := identifier.GetText()

	// Get the subsearch query context
	subsearchCtx := ctx.SearchQuery()
	if subsearchCtx == nil {
		return fmt.Errorf("join command requires a subsearch query")
	}

	// Visit the subsearch to build its AST
	result := subsearchCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	// VisitSearchQuery returns []ast.Command, wrap it in a Query
	commands, ok := result.([]ast.Command)
	if !ok {
		return fmt.Errorf("invalid subsearch in join command")
	}

	subsearch := &ast.Query{
		BaseNode: ast.BaseNode{Pos: getPosition(subsearchCtx)},
		Commands: commands,
	}

	return &ast.JoinCommand{
		BaseNode:  ast.BaseNode{Pos: getPosition(ctx)},
		JoinType:  joinType,
		JoinField: joinField,
		Subsearch: subsearch,
	}
}

// VisitJoinType builds a JoinType enum value
func (b *ASTBuilder) VisitJoinType(ctx *generated.JoinTypeContext) interface{} {
	switch {
	case ctx.INNER() != nil:
		return ast.JoinTypeInner
	case ctx.LEFT() != nil:
		return ast.JoinTypeLeft
	case ctx.RIGHT() != nil:
		return ast.JoinTypeRight
	case ctx.OUTER() != nil:
		return ast.JoinTypeOuter
	case ctx.FULL() != nil:
		return ast.JoinTypeFull
	default:
		return fmt.Errorf("unknown join type")
	}
}

// VisitTableCommand builds a TableCommand AST node
func (b *ASTBuilder) VisitTableCommand(ctx *generated.TableCommandContext) interface{} {
	fieldListCtx := ctx.FieldList()
	if fieldListCtx == nil {
		return fmt.Errorf("table command missing field list")
	}

	result := fieldListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	return &ast.TableCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Fields:   result.([]ast.Expression),
	}
}

// VisitEventstatsCommand builds an EventstatsCommand AST node
func (b *ASTBuilder) VisitEventstatsCommand(ctx *generated.EventstatsCommandContext) interface{} {
	// Process aggregations
	aggListCtx := ctx.AggregationList()
	if aggListCtx == nil {
		return fmt.Errorf("eventstats command missing aggregation list")
	}

	result := aggListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	aggregations := result.([]*ast.Aggregation)

	// Process group by (optional)
	var groupBy []ast.Expression
	if fieldListCtx := ctx.FieldList(); fieldListCtx != nil {
		result := fieldListCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		groupBy = result.([]ast.Expression)
	}

	return &ast.EventstatsCommand{
		BaseNode:     ast.BaseNode{Pos: getPosition(ctx)},
		Aggregations: aggregations,
		GroupBy:      groupBy,
	}
}

// VisitStreamstatsCommand builds a StreamstatsCommand AST node
func (b *ASTBuilder) VisitStreamstatsCommand(ctx *generated.StreamstatsCommandContext) interface{} {
	cmd := &ast.StreamstatsCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Current:  true, // Default: include current event
	}

	// Process options first
	for _, optCtx := range ctx.AllStreamstatsOptions() {
		if opt, ok := optCtx.(*generated.StreamstatsOptionsContext); ok {
			b.applyStreamstatsOption(cmd, opt)
		}
	}

	// Process aggregations
	aggListCtx := ctx.AggregationList()
	if aggListCtx == nil {
		return fmt.Errorf("streamstats command missing aggregation list")
	}

	result := aggListCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}
	cmd.Aggregations = result.([]*ast.Aggregation)

	// Process group by (optional)
	if fieldListCtx := ctx.FieldList(); fieldListCtx != nil {
		result := fieldListCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		cmd.GroupBy = result.([]ast.Expression)
	}

	return cmd
}

// applyStreamstatsOption applies streamstats options to the command
func (b *ASTBuilder) applyStreamstatsOption(cmd *ast.StreamstatsCommand, ctx *generated.StreamstatsOptionsContext) {
	idNode := ctx.IDENTIFIER()
	if idNode == nil {
		return
	}

	optName := strings.ToLower(idNode.GetText())

	switch optName {
	case "window":
		if intNode := ctx.INTEGER(); intNode != nil {
			window, _ := strconv.Atoi(intNode.GetText())
			cmd.Window = window
		}
	case "current":
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.Current = boolCtx.TRUE() != nil
		}
	case "global":
		if boolCtx := ctx.BooleanValue(); boolCtx != nil {
			cmd.Global = boolCtx.TRUE() != nil
		}
	}
}

// VisitReverseCommand builds a ReverseCommand AST node
// NOTE: Requires parser regeneration after grammar changes
// Uncomment after running: antlr4 -Dlanguage=Go -package generated -o pkg/ppl/parser/generated PPLLexer.g4 PPLParser.g4
/* func (b *ASTBuilder) VisitReverseCommand(ctx *generated.ReverseCommandContext) interface{} {
	return &ast.ReverseCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
	}
} */

// VisitFlattenCommand builds a FlattenCommand AST node
// NOTE: Requires parser regeneration after grammar changes
// Uncomment after running: antlr4 -Dlanguage=Go -package generated -o pkg/ppl/parser/generated PPLLexer.g4 PPLParser.g4
/* func (b *ASTBuilder) VisitFlattenCommand(ctx *generated.FlattenCommandContext) interface{} {
	fieldRefCtx := ctx.FieldReference()
	if fieldRefCtx == nil {
		return fmt.Errorf("flatten command missing field reference")
	}

	result := fieldRefCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	fieldExpr, ok := result.(ast.Expression)
	if !ok {
		return fmt.Errorf("flatten field is not a valid expression")
	}

	return &ast.FlattenCommand{
		BaseNode: ast.BaseNode{Pos: getPosition(ctx)},
		Field:    fieldExpr,
	}
} */

// VisitFillnullCommand builds a FillnullCommand AST node
// NOTE: Requires parser regeneration after grammar changes
// Uncomment after running: antlr4 -Dlanguage=Go -package generated -o pkg/ppl/parser/generated PPLLexer.g4 PPLParser.g4
/* func (b *ASTBuilder) VisitFillnullCommand(ctx *generated.FillnullCommandContext) interface{} {
	// Get the fill value
	valueCtx := ctx.LiteralValue()
	if valueCtx == nil {
		return fmt.Errorf("fillnull command missing value")
	}

	result := valueCtx.Accept(b)
	if err, ok := result.(error); ok {
		return err
	}

	valueExpr, ok := result.(ast.Expression)
	if !ok {
		return fmt.Errorf("fillnull value is not a valid expression")
	}

	// Get the optional fields list
	var fields []ast.Expression
	fieldListCtx := ctx.FieldList()
	if fieldListCtx != nil {
		result := fieldListCtx.Accept(b)
		if err, ok := result.(error); ok {
			return err
		}
		if fieldList, ok := result.([]string); ok {
			// Convert strings to field references
			fields = make([]ast.Expression, len(fieldList))
			for i, fieldName := range fieldList {
				fields[i] = &ast.FieldReference{
					BaseNode: ast.BaseNode{},
					Name:     fieldName,
				}
			}
		}
	}

	return &ast.FillnullCommand{
		BaseNode:     ast.BaseNode{Pos: getPosition(ctx)},
		DefaultValue: valueExpr,
		Fields:       fields,
	}
} */

// unquoteString removes quotes from a string literal
func unquoteString(s string) string {
	if len(s) < 2 {
		return s
	}
	// Remove surrounding quotes
	if (s[0] == '\'' && s[len(s)-1] == '\'') ||
		(s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '`' && s[len(s)-1] == '`') {
		return s[1 : len(s)-1]
	}
	return s
}
