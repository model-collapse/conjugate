// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery_String(t *testing.T) {
	tests := []struct {
		name     string
		query    *Query
		expected string
	}{
		{
			name: "single command",
			query: &Query{
				Commands: []Command{
					&SearchCommand{Source: "logs"},
				},
			},
			expected: "source=logs",
		},
		{
			name: "two commands",
			query: &Query{
				Commands: []Command{
					&SearchCommand{Source: "logs"},
					&HeadCommand{Count: 10},
				},
			},
			expected: "source=logs | head 10",
		},
		{
			name: "complex pipeline",
			query: &Query{
				Commands: []Command{
					&SearchCommand{Source: "logs"},
					&WhereCommand{
						Condition: &BinaryExpression{
							Left:     &FieldReference{Name: "status"},
							Operator: "=",
							Right:    &Literal{Value: 200, LiteralTyp: LiteralTypeInt},
						},
					},
					&HeadCommand{Count: 10},
				},
			},
			expected: "source=logs | where (status = 200) | head 10",
		},
		{
			name:     "empty query",
			query:    &Query{Commands: []Command{}},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.query.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQuery_Type(t *testing.T) {
	query := &Query{}
	assert.Equal(t, NodeTypeQuery, query.Type())
}

func TestSearchCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *SearchCommand
		expected string
	}{
		{
			name:     "simple source",
			cmd:      &SearchCommand{Source: "logs"},
			expected: "source=logs",
		},
		{
			name:     "source with dots",
			cmd:      &SearchCommand{Source: "app.logs.prod"},
			expected: "source=app.logs.prod",
		},
		{
			name:     "empty source",
			cmd:      &SearchCommand{Source: ""},
			expected: "source=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchCommand_Type(t *testing.T) {
	cmd := &SearchCommand{}
	assert.Equal(t, NodeTypeSearchCommand, cmd.Type())
}

func TestWhereCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *WhereCommand
		expected string
	}{
		{
			name: "simple condition",
			cmd: &WhereCommand{
				Condition: &BinaryExpression{
					Left:     &FieldReference{Name: "status"},
					Operator: "=",
					Right:    &Literal{Value: 200, LiteralTyp: LiteralTypeInt},
				},
			},
			expected: "where (status = 200)",
		},
		{
			name: "AND condition",
			cmd: &WhereCommand{
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
			expected: "where ((status = 200) AND (method = \"GET\"))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWhereCommand_Type(t *testing.T) {
	cmd := &WhereCommand{}
	assert.Equal(t, NodeTypeWhereCommand, cmd.Type())
}

func TestFieldsCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *FieldsCommand
		expected string
	}{
		{
			name: "include fields",
			cmd: &FieldsCommand{
				Fields: []Expression{
					&FieldReference{Name: "timestamp"},
					&FieldReference{Name: "message"},
				},
				Includes: true,
			},
			expected: "fields timestamp, message",
		},
		{
			name: "exclude fields",
			cmd: &FieldsCommand{
				Fields: []Expression{
					&FieldReference{Name: "internal_id"},
					&FieldReference{Name: "debug_info"},
				},
				Includes: false,
			},
			expected: "fields - internal_id, debug_info",
		},
		{
			name: "single field",
			cmd: &FieldsCommand{
				Fields: []Expression{
					&FieldReference{Name: "status"},
				},
				Includes: true,
			},
			expected: "fields status",
		},
		{
			name: "empty fields",
			cmd: &FieldsCommand{
				Fields:   []Expression{},
				Includes: true,
			},
			expected: "fields ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFieldsCommand_Type(t *testing.T) {
	cmd := &FieldsCommand{}
	assert.Equal(t, NodeTypeFieldsCommand, cmd.Type())
}

func TestAggregation_String(t *testing.T) {
	tests := []struct {
		name     string
		agg      *Aggregation
		expected string
	}{
		{
			name: "no alias",
			agg: &Aggregation{
				Func: &FunctionCall{
					Name:      "count",
					Arguments: []Expression{},
				},
			},
			expected: "count()",
		},
		{
			name: "with alias",
			agg: &Aggregation{
				Func: &FunctionCall{
					Name: "avg",
					Arguments: []Expression{
						&FieldReference{Name: "response_time"},
					},
				},
				Alias: "avg_time",
			},
			expected: "avg(response_time) as avg_time",
		},
		{
			name: "distinct count",
			agg: &Aggregation{
				Func: &FunctionCall{
					Name: "count",
					Arguments: []Expression{
						&FieldReference{Name: "user_id"},
					},
					Distinct: true,
				},
				Alias: "unique_users",
			},
			expected: "count(DISTINCT user_id) as unique_users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.agg.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAggregation_Type(t *testing.T) {
	agg := &Aggregation{}
	assert.Equal(t, NodeTypeAggregation, agg.Type())
}

func TestStatsCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *StatsCommand
		expected string
	}{
		{
			name: "simple count",
			cmd: &StatsCommand{
				Aggregations: []*Aggregation{
					{
						Func: &FunctionCall{
							Name:      "count",
							Arguments: []Expression{},
						},
					},
				},
			},
			expected: "stats count()",
		},
		{
			name: "with group by",
			cmd: &StatsCommand{
				Aggregations: []*Aggregation{
					{
						Func: &FunctionCall{
							Name:      "count",
							Arguments: []Expression{},
						},
					},
				},
				GroupBy: []Expression{
					&FieldReference{Name: "status"},
				},
			},
			expected: "stats count() by status",
		},
		{
			name: "multiple aggregations",
			cmd: &StatsCommand{
				Aggregations: []*Aggregation{
					{
						Func: &FunctionCall{
							Name:      "count",
							Arguments: []Expression{},
						},
						Alias: "total",
					},
					{
						Func: &FunctionCall{
							Name: "avg",
							Arguments: []Expression{
								&FieldReference{Name: "response_time"},
							},
						},
						Alias: "avg_time",
					},
				},
				GroupBy: []Expression{
					&FieldReference{Name: "status"},
					&FieldReference{Name: "method"},
				},
			},
			expected: "stats count() as total, avg(response_time) as avg_time by status, method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStatsCommand_Type(t *testing.T) {
	cmd := &StatsCommand{}
	assert.Equal(t, NodeTypeStatsCommand, cmd.Type())
}

func TestSortKey_String(t *testing.T) {
	tests := []struct {
		name     string
		key      *SortKey
		expected string
	}{
		{
			name: "ascending",
			key: &SortKey{
				Field:      &FieldReference{Name: "timestamp"},
				Descending: false,
			},
			expected: "timestamp",
		},
		{
			name: "descending",
			key: &SortKey{
				Field:      &FieldReference{Name: "count"},
				Descending: true,
			},
			expected: "count DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.key.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSortKey_Type(t *testing.T) {
	key := &SortKey{}
	assert.Equal(t, NodeTypeSortKey, key.Type())
}

func TestSortCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *SortCommand
		expected string
	}{
		{
			name: "single field ascending",
			cmd: &SortCommand{
				SortKeys: []*SortKey{
					{
						Field:      &FieldReference{Name: "timestamp"},
						Descending: false,
					},
				},
			},
			expected: "sort timestamp",
		},
		{
			name: "single field descending",
			cmd: &SortCommand{
				SortKeys: []*SortKey{
					{
						Field:      &FieldReference{Name: "count"},
						Descending: true,
					},
				},
			},
			expected: "sort count DESC",
		},
		{
			name: "multiple fields",
			cmd: &SortCommand{
				SortKeys: []*SortKey{
					{
						Field:      &FieldReference{Name: "status"},
						Descending: true,
					},
					{
						Field:      &FieldReference{Name: "timestamp"},
						Descending: false,
					},
				},
			},
			expected: "sort status DESC, timestamp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSortCommand_Type(t *testing.T) {
	cmd := &SortCommand{}
	assert.Equal(t, NodeTypeSortCommand, cmd.Type())
}

func TestHeadCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *HeadCommand
		expected string
	}{
		{
			name:     "head 10",
			cmd:      &HeadCommand{Count: 10},
			expected: "head 10",
		},
		{
			name:     "head 1",
			cmd:      &HeadCommand{Count: 1},
			expected: "head 1",
		},
		{
			name:     "head 1000",
			cmd:      &HeadCommand{Count: 1000},
			expected: "head 1000",
		},
		{
			name:     "head 0",
			cmd:      &HeadCommand{Count: 0},
			expected: "head 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHeadCommand_Type(t *testing.T) {
	cmd := &HeadCommand{}
	assert.Equal(t, NodeTypeHeadCommand, cmd.Type())
}

func TestDescribeCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *DescribeCommand
		expected string
	}{
		{
			name:     "simple source",
			cmd:      &DescribeCommand{Source: "logs"},
			expected: "describe logs",
		},
		{
			name:     "nested source",
			cmd:      &DescribeCommand{Source: "app.logs.prod"},
			expected: "describe app.logs.prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDescribeCommand_Type(t *testing.T) {
	cmd := &DescribeCommand{}
	assert.Equal(t, NodeTypeDescribeCommand, cmd.Type())
}

func TestShowDatasourcesCommand_String(t *testing.T) {
	cmd := &ShowDatasourcesCommand{}
	assert.Equal(t, "showdatasources", cmd.String())
}

func TestShowDatasourcesCommand_Type(t *testing.T) {
	cmd := &ShowDatasourcesCommand{}
	assert.Equal(t, NodeTypeShowDatasourcesCommand, cmd.Type())
}

func TestExplainCommand_String(t *testing.T) {
	cmd := &ExplainCommand{}
	assert.Equal(t, "explain", cmd.String())
}

func TestExplainCommand_Type(t *testing.T) {
	cmd := &ExplainCommand{}
	assert.Equal(t, NodeTypeExplainCommand, cmd.Type())
}

func TestCommand_EdgeCases(t *testing.T) {
	t.Run("query with nil commands", func(t *testing.T) {
		query := &Query{Commands: nil}
		// Should not panic
		result := query.String()
		assert.Equal(t, "", result)
	})

	t.Run("fields with nil expressions", func(t *testing.T) {
		cmd := &FieldsCommand{
			Fields:   nil,
			Includes: true,
		}
		// Should not panic
		result := cmd.String()
		assert.Equal(t, "fields ", result)
	})

	t.Run("stats with nil aggregations", func(t *testing.T) {
		cmd := &StatsCommand{
			Aggregations: nil,
			GroupBy:      nil,
		}
		// Should not panic
		result := cmd.String()
		assert.Equal(t, "stats ", result)
	})

	t.Run("sort with nil keys", func(t *testing.T) {
		cmd := &SortCommand{SortKeys: nil}
		// Should not panic
		result := cmd.String()
		assert.Equal(t, "sort ", result)
	})

	t.Run("negative head count", func(t *testing.T) {
		cmd := &HeadCommand{Count: -1}
		result := cmd.String()
		assert.Equal(t, "head -1", result)
	})
}

func TestCommand_Accept(t *testing.T) {
	// Create a simple visitor to test Accept
	visitor := &BaseVisitor{}

	tests := []struct {
		name string
		cmd  Command
	}{
		{"SearchCommand", &SearchCommand{}},
		{"WhereCommand", &WhereCommand{}},
		{"FieldsCommand", &FieldsCommand{}},
		{"StatsCommand", &StatsCommand{}},
		{"SortCommand", &SortCommand{}},
		{"HeadCommand", &HeadCommand{}},
		{"DescribeCommand", &DescribeCommand{}},
		{"ShowDatasourcesCommand", &ShowDatasourcesCommand{}},
		{"ExplainCommand", &ExplainCommand{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result, err := tt.cmd.Accept(visitor)
			require.NoError(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestComplexQueryConstruction(t *testing.T) {
	// Test building a complex query programmatically
	query := &Query{
		Commands: []Command{
			&SearchCommand{
				BaseNode: BaseNode{Pos: Position{Line: 1, Column: 1, Offset: 0}},
				Source:   "logs",
			},
			&WhereCommand{
				BaseNode: BaseNode{Pos: Position{Line: 1, Column: 16, Offset: 15}},
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
				BaseNode: BaseNode{Pos: Position{Line: 2, Column: 1, Offset: 50}},
				Aggregations: []*Aggregation{
					{
						Func: &FunctionCall{
							Name:      "count",
							Arguments: []Expression{},
						},
						Alias: "total",
					},
					{
						Func: &FunctionCall{
							Name: "avg",
							Arguments: []Expression{
								&FieldReference{Name: "response_time"},
							},
						},
						Alias: "avg_time",
					},
				},
				GroupBy: []Expression{
					&FieldReference{Name: "endpoint"},
				},
			},
			&SortCommand{
				BaseNode: BaseNode{Pos: Position{Line: 3, Column: 1, Offset: 100}},
				SortKeys: []*SortKey{
					{
						Field:      &FieldReference{Name: "total"},
						Descending: true,
					},
				},
			},
			&HeadCommand{
				BaseNode: BaseNode{Pos: Position{Line: 4, Column: 1, Offset: 120}},
				Count:    10,
			},
		},
	}

	// Verify structure
	require.Len(t, query.Commands, 5)

	// Verify types
	assert.IsType(t, &SearchCommand{}, query.Commands[0])
	assert.IsType(t, &WhereCommand{}, query.Commands[1])
	assert.IsType(t, &StatsCommand{}, query.Commands[2])
	assert.IsType(t, &SortCommand{}, query.Commands[3])
	assert.IsType(t, &HeadCommand{}, query.Commands[4])

	// Verify String() output
	result := query.String()
	assert.Contains(t, result, "source=logs")
	assert.Contains(t, result, "where")
	assert.Contains(t, result, "stats")
	assert.Contains(t, result, "sort")
	assert.Contains(t, result, "head 10")

	// Verify positions
	assert.Equal(t, 1, query.Commands[0].Position().Line)
	assert.Equal(t, 1, query.Commands[1].Position().Line)
	assert.Equal(t, 2, query.Commands[2].Position().Line)
	assert.Equal(t, 3, query.Commands[3].Position().Line)
	assert.Equal(t, 4, query.Commands[4].Position().Line)
}
