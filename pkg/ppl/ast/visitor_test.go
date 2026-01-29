// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CountingVisitor counts the number of nodes visited
type CountingVisitor struct {
	BaseVisitor
	QueryCount                  int
	SearchCommandCount          int
	WhereCommandCount           int
	FieldsCommandCount          int
	StatsCommandCount           int
	SortCommandCount            int
	HeadCommandCount            int
	DescribeCommandCount        int
	ShowDatasourcesCommandCount int
	ExplainCommandCount         int
	BinaryExpressionCount       int
	UnaryExpressionCount        int
	FunctionCallCount           int
	FieldReferenceCount         int
	LiteralCount                int
	ListLiteralCount            int
	CaseExpressionCount         int
	WhenClauseCount             int
	AggregationCount            int
	SortKeyCount                int
}

func (v *CountingVisitor) VisitQuery(q *Query) (interface{}, error) {
	v.QueryCount++
	for _, cmd := range q.Commands {
		if _, err := cmd.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitSearchCommand(s *SearchCommand) (interface{}, error) {
	v.SearchCommandCount++
	return nil, nil
}

func (v *CountingVisitor) VisitWhereCommand(w *WhereCommand) (interface{}, error) {
	v.WhereCommandCount++
	if w.Condition != nil {
		if _, err := w.Condition.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitFieldsCommand(f *FieldsCommand) (interface{}, error) {
	v.FieldsCommandCount++
	for _, field := range f.Fields {
		if _, err := field.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitStatsCommand(s *StatsCommand) (interface{}, error) {
	v.StatsCommandCount++
	for _, agg := range s.Aggregations {
		if _, err := agg.Accept(v); err != nil {
			return nil, err
		}
	}
	for _, group := range s.GroupBy {
		if _, err := group.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitSortCommand(s *SortCommand) (interface{}, error) {
	v.SortCommandCount++
	for _, key := range s.SortKeys {
		if _, err := key.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitHeadCommand(h *HeadCommand) (interface{}, error) {
	v.HeadCommandCount++
	return nil, nil
}

func (v *CountingVisitor) VisitDescribeCommand(d *DescribeCommand) (interface{}, error) {
	v.DescribeCommandCount++
	return nil, nil
}

func (v *CountingVisitor) VisitShowDatasourcesCommand(s *ShowDatasourcesCommand) (interface{}, error) {
	v.ShowDatasourcesCommandCount++
	return nil, nil
}

func (v *CountingVisitor) VisitExplainCommand(e *ExplainCommand) (interface{}, error) {
	v.ExplainCommandCount++
	return nil, nil
}

func (v *CountingVisitor) VisitBinaryExpression(b *BinaryExpression) (interface{}, error) {
	v.BinaryExpressionCount++
	if b.Left != nil {
		if _, err := b.Left.Accept(v); err != nil {
			return nil, err
		}
	}
	if b.Right != nil {
		if _, err := b.Right.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitUnaryExpression(u *UnaryExpression) (interface{}, error) {
	v.UnaryExpressionCount++
	if u.Operand != nil {
		if _, err := u.Operand.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitFunctionCall(f *FunctionCall) (interface{}, error) {
	v.FunctionCallCount++
	for _, arg := range f.Arguments {
		if _, err := arg.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitFieldReference(f *FieldReference) (interface{}, error) {
	v.FieldReferenceCount++
	return nil, nil
}

func (v *CountingVisitor) VisitLiteral(l *Literal) (interface{}, error) {
	v.LiteralCount++
	return nil, nil
}

func (v *CountingVisitor) VisitListLiteral(l *ListLiteral) (interface{}, error) {
	v.ListLiteralCount++
	for _, val := range l.Values {
		if _, err := val.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitCaseExpression(c *CaseExpression) (interface{}, error) {
	v.CaseExpressionCount++
	for _, when := range c.WhenClauses {
		if _, err := when.Accept(v); err != nil {
			return nil, err
		}
	}
	if c.ElseResult != nil {
		if _, err := c.ElseResult.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitWhenClause(w *WhenClause) (interface{}, error) {
	v.WhenClauseCount++
	if w.Condition != nil {
		if _, err := w.Condition.Accept(v); err != nil {
			return nil, err
		}
	}
	if w.Result != nil {
		if _, err := w.Result.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitAggregation(a *Aggregation) (interface{}, error) {
	v.AggregationCount++
	if a.Func != nil {
		if _, err := a.Func.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CountingVisitor) VisitSortKey(s *SortKey) (interface{}, error) {
	v.SortKeyCount++
	if s.Field != nil {
		if _, err := s.Field.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func TestBaseVisitor_AllMethodsReturnNil(t *testing.T) {
	visitor := &BaseVisitor{}

	tests := []struct {
		name string
		node Node
	}{
		{"Query", &Query{}},
		{"SearchCommand", &SearchCommand{}},
		{"WhereCommand", &WhereCommand{}},
		{"FieldsCommand", &FieldsCommand{}},
		{"StatsCommand", &StatsCommand{}},
		{"SortCommand", &SortCommand{}},
		{"HeadCommand", &HeadCommand{}},
		{"DescribeCommand", &DescribeCommand{}},
		{"ShowDatasourcesCommand", &ShowDatasourcesCommand{}},
		{"ExplainCommand", &ExplainCommand{}},
		{"BinaryExpression", &BinaryExpression{}},
		{"UnaryExpression", &UnaryExpression{}},
		{"FunctionCall", &FunctionCall{}},
		{"FieldReference", &FieldReference{}},
		{"Literal", &Literal{}},
		{"ListLiteral", &ListLiteral{}},
		{"CaseExpression", &CaseExpression{}},
		{"WhenClause", &WhenClause{}},
		{"Aggregation", &Aggregation{}},
		{"SortKey", &SortKey{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.node.Accept(visitor)
			require.NoError(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestCountingVisitor_SimpleQuery(t *testing.T) {
	// Create query: source=logs | head 10
	query := &Query{
		Commands: []Command{
			&SearchCommand{Source: "logs"},
			&HeadCommand{Count: 10},
		},
	}

	visitor := &CountingVisitor{}
	_, err := query.Accept(visitor)
	require.NoError(t, err)

	assert.Equal(t, 1, visitor.QueryCount)
	assert.Equal(t, 1, visitor.SearchCommandCount)
	assert.Equal(t, 1, visitor.HeadCommandCount)
	assert.Equal(t, 0, visitor.WhereCommandCount)
}

func TestCountingVisitor_ComplexQuery(t *testing.T) {
	// Create query: source=logs | where status = 200 AND method = 'GET' | stats count() by endpoint | sort count desc | head 10
	query := &Query{
		Commands: []Command{
			&SearchCommand{Source: "logs"},
			&WhereCommand{
				Condition: &BinaryExpression{
					Left: &BinaryExpression{
						Left:     &FieldReference{Name: "status"},
						Operator: "=",
						Right:    &Literal{Value: 200, LiteralTyp: LiteralTypeInt},
					},
					Operator: "AND",
					Right: &BinaryExpression{
						Left:     &FieldReference{Name: "method"},
						Operator: "=",
						Right:    &Literal{Value: "GET", LiteralTyp: LiteralTypeString},
					},
				},
			},
			&StatsCommand{
				Aggregations: []*Aggregation{
					{
						Func: &FunctionCall{
							Name:      "count",
							Arguments: []Expression{},
						},
					},
				},
				GroupBy: []Expression{
					&FieldReference{Name: "endpoint"},
				},
			},
			&SortCommand{
				SortKeys: []*SortKey{
					{
						Field:      &FieldReference{Name: "count"},
						Descending: true,
					},
				},
			},
			&HeadCommand{Count: 10},
		},
	}

	visitor := &CountingVisitor{}
	_, err := query.Accept(visitor)
	require.NoError(t, err)

	// Verify counts
	assert.Equal(t, 1, visitor.QueryCount)
	assert.Equal(t, 1, visitor.SearchCommandCount)
	assert.Equal(t, 1, visitor.WhereCommandCount)
	assert.Equal(t, 1, visitor.StatsCommandCount)
	assert.Equal(t, 1, visitor.SortCommandCount)
	assert.Equal(t, 1, visitor.HeadCommandCount)

	// Expression counts
	assert.Equal(t, 3, visitor.BinaryExpressionCount) // status=200, method='GET', AND
	assert.Equal(t, 4, visitor.FieldReferenceCount)   // status, method, endpoint, count
	assert.Equal(t, 2, visitor.LiteralCount)          // 200, 'GET'
	assert.Equal(t, 1, visitor.FunctionCallCount)     // count()
	assert.Equal(t, 1, visitor.AggregationCount)
	assert.Equal(t, 1, visitor.SortKeyCount)
}

func TestCountingVisitor_CaseExpression(t *testing.T) {
	// Create CASE expression
	caseExpr := &CaseExpression{
		WhenClauses: []*WhenClause{
			{
				Condition: &BinaryExpression{
					Left:     &FieldReference{Name: "status"},
					Operator: "<",
					Right:    &Literal{Value: 300, LiteralTyp: LiteralTypeInt},
				},
				Result: &Literal{Value: "success", LiteralTyp: LiteralTypeString},
			},
			{
				Condition: &BinaryExpression{
					Left:     &FieldReference{Name: "status"},
					Operator: "<",
					Right:    &Literal{Value: 500, LiteralTyp: LiteralTypeInt},
				},
				Result: &Literal{Value: "client_error", LiteralTyp: LiteralTypeString},
			},
		},
		ElseResult: &Literal{Value: "server_error", LiteralTyp: LiteralTypeString},
	}

	visitor := &CountingVisitor{}
	_, err := caseExpr.Accept(visitor)
	require.NoError(t, err)

	assert.Equal(t, 1, visitor.CaseExpressionCount)
	assert.Equal(t, 2, visitor.WhenClauseCount)
	assert.Equal(t, 2, visitor.BinaryExpressionCount) // status < 300, status < 500
	assert.Equal(t, 2, visitor.FieldReferenceCount)   // status (twice)
	assert.Equal(t, 5, visitor.LiteralCount)          // 300, "success", 500, "client_error", "server_error"
}

func TestCountingVisitor_ListLiteral(t *testing.T) {
	// Create list literal
	list := &ListLiteral{
		Values: []Expression{
			&Literal{Value: 1, LiteralTyp: LiteralTypeInt},
			&Literal{Value: 2, LiteralTyp: LiteralTypeInt},
			&Literal{Value: 3, LiteralTyp: LiteralTypeInt},
		},
	}

	visitor := &CountingVisitor{}
	_, err := list.Accept(visitor)
	require.NoError(t, err)

	assert.Equal(t, 1, visitor.ListLiteralCount)
	assert.Equal(t, 3, visitor.LiteralCount)
}

func TestCountingVisitor_UnaryExpression(t *testing.T) {
	// Create NOT expression
	expr := &UnaryExpression{
		Operator: "NOT",
		Operand: &BinaryExpression{
			Left:     &FieldReference{Name: "active"},
			Operator: "=",
			Right:    &Literal{Value: true, LiteralTyp: LiteralTypeBool},
		},
	}

	visitor := &CountingVisitor{}
	_, err := expr.Accept(visitor)
	require.NoError(t, err)

	assert.Equal(t, 1, visitor.UnaryExpressionCount)
	assert.Equal(t, 1, visitor.BinaryExpressionCount)
	assert.Equal(t, 1, visitor.FieldReferenceCount)
	assert.Equal(t, 1, visitor.LiteralCount)
}

func TestWalk_Function(t *testing.T) {
	// Test the Walk helper function
	query := &Query{
		Commands: []Command{
			&SearchCommand{Source: "logs"},
			&HeadCommand{Count: 10},
		},
	}

	visitor := &CountingVisitor{}
	result, err := Walk(query, visitor)
	require.NoError(t, err)
	assert.Nil(t, result)

	assert.Equal(t, 1, visitor.QueryCount)
	assert.Equal(t, 1, visitor.SearchCommandCount)
	assert.Equal(t, 1, visitor.HeadCommandCount)
}

// ErrorReturningVisitor returns an error on a specific node type
type ErrorReturningVisitor struct {
	BaseVisitor
	ErrorOnType NodeType
}

func (v *ErrorReturningVisitor) VisitWhereCommand(w *WhereCommand) (interface{}, error) {
	if v.ErrorOnType == NodeTypeWhereCommand {
		return nil, fmt.Errorf("test error on WhereCommand")
	}
	return nil, nil
}

func (v *ErrorReturningVisitor) VisitBinaryExpression(b *BinaryExpression) (interface{}, error) {
	if v.ErrorOnType == NodeTypeBinaryExpression {
		return nil, fmt.Errorf("test error on BinaryExpression")
	}
	if b.Left != nil {
		if _, err := b.Left.Accept(v); err != nil {
			return nil, err
		}
	}
	if b.Right != nil {
		if _, err := b.Right.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func TestVisitor_ErrorHandling(t *testing.T) {
	query := &Query{
		Commands: []Command{
			&SearchCommand{Source: "logs"},
			&WhereCommand{
				Condition: &BinaryExpression{
					Left:     &FieldReference{Name: "status"},
					Operator: "=",
					Right:    &Literal{Value: 200, LiteralTyp: LiteralTypeInt},
				},
			},
		},
	}

	t.Run("error on WhereCommand", func(t *testing.T) {
		visitor := &ErrorReturningVisitor{ErrorOnType: NodeTypeWhereCommand}
		_, err := Walk(query.Commands[1], visitor)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "test error on WhereCommand")
	})

	t.Run("error propagation from nested expression", func(t *testing.T) {
		visitor := &ErrorReturningVisitor{ErrorOnType: NodeTypeBinaryExpression}
		whereCmd := query.Commands[1].(*WhereCommand)

		// Visit where command, which will visit the condition
		_, err := whereCmd.Condition.Accept(visitor)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "test error on BinaryExpression")
	})
}

// CollectingVisitor collects all nodes of a specific type
type CollectingVisitor struct {
	BaseVisitor
	FieldReferences []*FieldReference
}

func (v *CollectingVisitor) VisitFieldReference(f *FieldReference) (interface{}, error) {
	v.FieldReferences = append(v.FieldReferences, f)
	return nil, nil
}

func (v *CollectingVisitor) VisitBinaryExpression(b *BinaryExpression) (interface{}, error) {
	if b.Left != nil {
		if _, err := b.Left.Accept(v); err != nil {
			return nil, err
		}
	}
	if b.Right != nil {
		if _, err := b.Right.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *CollectingVisitor) VisitWhereCommand(w *WhereCommand) (interface{}, error) {
	if w.Condition != nil {
		if _, err := w.Condition.Accept(v); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func TestCollectingVisitor_GatherFieldReferences(t *testing.T) {
	// Create where condition with multiple field references
	where := &WhereCommand{
		Condition: &BinaryExpression{
			Left: &BinaryExpression{
				Left:     &FieldReference{Name: "status"},
				Operator: "=",
				Right:    &Literal{Value: 200, LiteralTyp: LiteralTypeInt},
			},
			Operator: "AND",
			Right: &BinaryExpression{
				Left:     &FieldReference{Name: "method"},
				Operator: "=",
				Right:    &Literal{Value: "GET", LiteralTyp: LiteralTypeString},
			},
		},
	}

	visitor := &CollectingVisitor{}
	_, err := where.Accept(visitor)
	require.NoError(t, err)

	// Should have collected 2 field references
	require.Len(t, visitor.FieldReferences, 2)
	assert.Equal(t, "status", visitor.FieldReferences[0].Name)
	assert.Equal(t, "method", visitor.FieldReferences[1].Name)
}

func TestVisitor_NilHandling(t *testing.T) {
	t.Run("where command with nil condition", func(t *testing.T) {
		where := &WhereCommand{Condition: nil}
		visitor := &CountingVisitor{}

		// Should not panic
		_, err := where.Accept(visitor)
		require.NoError(t, err)
		assert.Equal(t, 1, visitor.WhereCommandCount)
	})

	t.Run("binary expression with nil operands", func(t *testing.T) {
		expr := &BinaryExpression{
			Left:     nil,
			Operator: "=",
			Right:    nil,
		}
		visitor := &CountingVisitor{}

		// Should not panic
		_, err := expr.Accept(visitor)
		require.NoError(t, err)
		assert.Equal(t, 1, visitor.BinaryExpressionCount)
	})

	t.Run("function call with nil arguments", func(t *testing.T) {
		fn := &FunctionCall{
			Name:      "count",
			Arguments: nil,
		}
		visitor := &CountingVisitor{}

		// Should not panic
		_, err := fn.Accept(visitor)
		require.NoError(t, err)
		assert.Equal(t, 1, visitor.FunctionCallCount)
	})
}
