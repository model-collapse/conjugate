// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package parser

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: These tests require ANTLR4-generated code.
// Run 'make' in this directory to generate parser code before running tests.

func TestParser_SearchCommand(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	tests := []struct {
		name   string
		query  string
		source string
	}{
		{
			name:   "simple search",
			query:  "source=logs",
			source: "logs",
		},
		{
			name:   "search with keyword",
			query:  "search source=logs",
			source: "logs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)
			require.Len(t, astQuery.Commands, 1)

			searchCmd, ok := astQuery.Commands[0].(*ast.SearchCommand)
			require.True(t, ok)
			assert.Equal(t, tt.source, searchCmd.Source)
		})
	}
}

func TestParser_WhereCommand(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "simple where",
			query: "source=logs | where status = 200",
		},
		{
			name:  "where with and",
			query: "source=logs | where status = 200 and method = 'GET'",
		},
		{
			name:  "where with or",
			query: "source=logs | where status = 200 or status = 404",
		},
		{
			name:  "where with not",
			query: "source=logs | where not status = 500",
		},
		{
			name:  "where with comparison",
			query: "source=logs | where response_time > 1000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)
			assert.GreaterOrEqual(t, len(astQuery.Commands), 2)

			whereCmd, ok := astQuery.Commands[1].(*ast.WhereCommand)
			require.True(t, ok)
			assert.NotNil(t, whereCmd.Condition)
		})
	}
}

func TestParser_FieldsCommand(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	tests := []struct {
		name     string
		query    string
		includes bool
		numFields int
	}{
		{
			name:      "include fields",
			query:     "source=logs | fields timestamp, message, level",
			includes:  true,
			numFields: 3,
		},
		{
			name:      "exclude fields",
			query:     "source=logs | fields - internal_id, debug_info",
			includes:  false,
			numFields: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			fieldsCmd, ok := astQuery.Commands[1].(*ast.FieldsCommand)
			require.True(t, ok)
			assert.Equal(t, tt.includes, fieldsCmd.Includes)
			assert.Len(t, fieldsCmd.Fields, tt.numFields)
		})
	}
}

func TestParser_StatsCommand(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	tests := []struct {
		name    string
		query   string
		hasGroupBy bool
	}{
		{
			name:       "simple stats",
			query:      "source=logs | stats count()",
			hasGroupBy: false,
		},
		{
			name:       "stats with group by",
			query:      "source=logs | stats count() by status",
			hasGroupBy: true,
		},
		{
			name:       "stats with multiple aggregations",
			query:      "source=logs | stats count() as total, avg(response_time) as avg_time by status",
			hasGroupBy: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			statsCmd, ok := astQuery.Commands[1].(*ast.StatsCommand)
			require.True(t, ok)
			assert.Greater(t, len(statsCmd.Aggregations), 0)

			if tt.hasGroupBy {
				assert.Greater(t, len(statsCmd.GroupBy), 0)
			} else {
				assert.Len(t, statsCmd.GroupBy, 0)
			}
		})
	}
}

func TestParser_SortCommand(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "sort ascending",
			query: "source=logs | sort timestamp",
		},
		{
			name:  "sort descending",
			query: "source=logs | sort timestamp desc",
		},
		{
			name:  "multi-field sort",
			query: "source=logs | sort status desc, timestamp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			sortCmd, ok := astQuery.Commands[1].(*ast.SortCommand)
			require.True(t, ok)
			assert.Greater(t, len(sortCmd.SortKeys), 0)
		})
	}
}

func TestParser_HeadCommand(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	tests := []struct {
		name  string
		query string
		count int
	}{
		{
			name:  "head 10",
			query: "source=logs | head 10",
			count: 10,
		},
		{
			name:  "head 100",
			query: "source=logs | head 100",
			count: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			headCmd, ok := astQuery.Commands[1].(*ast.HeadCommand)
			require.True(t, ok)
			assert.Equal(t, tt.count, headCmd.Count)
		})
	}
}

func TestParser_DescribeCommand(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	query := "describe logs"
	astQuery, err := parser.Parse(query)
	require.NoError(t, err)
	require.NotNil(t, astQuery)
	require.Len(t, astQuery.Commands, 1)

	describeCmd, ok := astQuery.Commands[0].(*ast.DescribeCommand)
	require.True(t, ok)
	assert.Equal(t, "logs", describeCmd.Source)
}

func TestParser_ShowDatasourcesCommand(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	query := "showdatasources"
	astQuery, err := parser.Parse(query)
	require.NoError(t, err)
	require.NotNil(t, astQuery)
	require.Len(t, astQuery.Commands, 1)

	_, ok := astQuery.Commands[0].(*ast.ShowDatasourcesCommand)
	require.True(t, ok)
}

func TestParser_ExplainCommand(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	query := "explain source=logs | where status = 200"
	astQuery, err := parser.Parse(query)
	require.NoError(t, err)
	require.NotNil(t, astQuery)

	// Explain should be the first command
	_, ok := astQuery.Commands[0].(*ast.ExplainCommand)
	require.True(t, ok)
}

func TestParser_ComplexQuery(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	query := `source=logs
		| where status = 200 and method = 'GET'
		| stats count() as total, avg(response_time) as avg_time by endpoint
		| sort total desc
		| head 10`

	astQuery, err := parser.Parse(query)
	require.NoError(t, err)
	require.NotNil(t, astQuery)
	assert.Len(t, astQuery.Commands, 5)

	// Verify command types
	_, ok := astQuery.Commands[0].(*ast.SearchCommand)
	assert.True(t, ok)

	_, ok = astQuery.Commands[1].(*ast.WhereCommand)
	assert.True(t, ok)

	_, ok = astQuery.Commands[2].(*ast.StatsCommand)
	assert.True(t, ok)

	_, ok = astQuery.Commands[3].(*ast.SortCommand)
	assert.True(t, ok)

	_, ok = astQuery.Commands[4].(*ast.HeadCommand)
	assert.True(t, ok)
}

func TestParser_SyntaxErrors(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "missing source",
			query: "where status = 200",
		},
		{
			name:  "invalid operator",
			query: "source=logs | where status === 200",
		},
		{
			name:  "unclosed parenthesis",
			query: "source=logs | stats count(",
		},
		{
			name:  "missing pipe",
			query: "source=logs where status = 200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.query)
			assert.Error(t, err)
		})
	}
}

func TestParser_Expressions(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "arithmetic expression",
			query: "source=logs | where response_time > 1000 + 500",
		},
		{
			name:  "nested expression",
			query: "source=logs | where (status = 200 or status = 404) and method = 'GET'",
		},
		{
			name:  "function call",
			query: "source=logs | stats max(response_time)",
		},
		{
			name:  "case expression",
			query: "source=logs | stats count() as total by case when status < 300 then 'success' when status < 500 then 'client_error' else 'server_error' end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			assert.NotNil(t, astQuery)
		})
	}
}

func TestParser_ValidateSyntax(t *testing.T) {
	// ANTLR4 code is now generated

	parser := NewParser()

	// Valid query
	err := parser.ValidateSyntax("source=logs | where status = 200")
	assert.NoError(t, err)

	// Invalid query
	err = parser.ValidateSyntax("source=logs | where")
	assert.Error(t, err)
}

func TestErrorListener(t *testing.T) {
	listener := NewErrorListener()

	// Initially no errors
	assert.False(t, listener.HasErrors())
	assert.Nil(t, listener.GetError())

	// Add error
	listener.SyntaxError(nil, nil, 1, 5, "test error", nil)

	// Should have error now
	assert.True(t, listener.HasErrors())
	assert.NotNil(t, listener.GetError())

	errors := listener.GetErrors()
	assert.Len(t, errors, 1)
	assert.Equal(t, 1, errors[0].Line)
	assert.Equal(t, 5, errors[0].Column)

	// Reset
	listener.Reset()
	assert.False(t, listener.HasErrors())
}

// ============================================================================
// Tier 1 Command Tests
// ============================================================================

func TestParser_Tier1_TopCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple top", "source=logs | top host"},
		{"top with limit", "source=logs | top 10 host, status"},
		{"top with by", "source=logs | top 5 host by datacenter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			_, ok := astQuery.Commands[1].(*ast.TopCommand)
			assert.True(t, ok, "Expected TopCommand")
		})
	}
}

func TestParser_Tier1_RareCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple rare", "source=logs | rare host"},
		{"rare with limit", "source=logs | rare 10 host, status"},
		{"rare with by", "source=logs | rare 5 host by datacenter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			_, ok := astQuery.Commands[1].(*ast.RareCommand)
			assert.True(t, ok, "Expected RareCommand")
		})
	}
}

func TestParser_Tier1_DedupCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple dedup", "source=logs | dedup host"},
		{"dedup with count", "source=logs | dedup 5 host, status"},
		{"dedup with options", "source=logs | dedup host keepevents=true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			_, ok := astQuery.Commands[1].(*ast.DedupCommand)
			assert.True(t, ok, "Expected DedupCommand")
		})
	}
}

func TestParser_Tier1_EvalCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple eval", "source=logs | eval duration_sec = duration / 1000"},
		{"eval with function", "source=logs | eval is_error = status >= 400"},
		{"multi eval", "source=logs | eval a = b + 1, c = d * 2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			evalCmd, ok := astQuery.Commands[1].(*ast.EvalCommand)
			assert.True(t, ok, "Expected EvalCommand")
			assert.NotEmpty(t, evalCmd.Assignments)
		})
	}
}

func TestParser_Tier1_RenameCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple rename", "source=logs | rename old as new"},
		{"multi rename", "source=logs | rename a as b, c as d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			renameCmd, ok := astQuery.Commands[1].(*ast.RenameCommand)
			assert.True(t, ok, "Expected RenameCommand")
			assert.NotEmpty(t, renameCmd.Assignments)
		})
	}
}

func TestParser_Tier1_TableCommand(t *testing.T) {
	parser := NewParser()

	query := "source=logs | table host, status, latency"
	astQuery, err := parser.Parse(query)
	require.NoError(t, err)
	require.NotNil(t, astQuery)

	tableCmd, ok := astQuery.Commands[1].(*ast.TableCommand)
	assert.True(t, ok, "Expected TableCommand")
	assert.Len(t, tableCmd.Fields, 3)
}

func TestParser_Tier1_EventstatsCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple eventstats", "source=logs | eventstats avg(latency) as avg_lat"},
		{"eventstats with by", "source=logs | eventstats count() as total by host"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			_, ok := astQuery.Commands[1].(*ast.EventstatsCommand)
			assert.True(t, ok, "Expected EventstatsCommand")
		})
	}
}

func TestParser_Tier1_ChartCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple chart", "source=logs | chart count() by host"},
		{"chart with agg", "source=logs | chart avg(latency), max(latency) by status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			_, ok := astQuery.Commands[1].(*ast.ChartCommand)
			assert.True(t, ok, "Expected ChartCommand")
		})
	}
}

func TestParser_Tier1_TimechartCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple timechart", "source=logs | timechart count()"},
		{"timechart with span", "source=logs | timechart span=1h count() by status"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			_, ok := astQuery.Commands[1].(*ast.TimechartCommand)
			assert.True(t, ok, "Expected TimechartCommand")
		})
	}
}

func TestParser_Tier1_BinCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple bin", "source=logs | bin timestamp"},
		{"bin with span", "source=logs | bin timestamp span=1h"},
		{"bin with bins", "source=logs | bin latency bins=10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)

			_, ok := astQuery.Commands[1].(*ast.BinCommand)
			assert.True(t, ok, "Expected BinCommand")
		})
	}
}

func TestParser_Tier2_FillnullCommand(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name  string
		query string
	}{
		{"simple fillnull", "source=logs | fillnull status='unknown'"},
		{"multi field fillnull", "source=logs | fillnull host='localhost', status='ok'"},
		{"fillnull with value", "source=logs | fillnull value=0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			astQuery, err := parser.Parse(tt.query)
			require.NoError(t, err)
			require.NotNil(t, astQuery)
			require.Len(t, astQuery.Commands, 2, "Expected 2 commands (search + fillnull)")

			fillnullCmd, ok := astQuery.Commands[1].(*ast.FillnullCommand)
			assert.True(t, ok, "Expected FillnullCommand, got %T", astQuery.Commands[1])

			if tt.name == "fillnull with value" {
				// Should have default value
				assert.NotNil(t, fillnullCmd.DefaultValue, "Expected default value")
			} else {
				// Should have assignments
				assert.NotEmpty(t, fillnullCmd.Assignments, "Expected assignments")
			}
		})
	}
}
