// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestReplaceOperator_Basic(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	// Create test data
	rows := []*Row{
		NewRow(map[string]interface{}{
			"level":   "error",
			"message": "Connection failed",
		}),
		NewRow(map[string]interface{}{
			"level":   "warn",
			"message": "Slow query",
		}),
		NewRow(map[string]interface{}{
			"level":   "info",
			"message": "Request processed",
		}),
	}

	iter := NewSliceIterator(rows)
	mappings := []*ast.ReplaceMapping{
		{
			OldValue: &ast.Literal{Value: "error"},
			NewValue: &ast.Literal{Value: "ERROR"},
		},
		{
			OldValue: &ast.Literal{Value: "warn"},
			NewValue: &ast.Literal{Value: "WARNING"},
		},
	}
	replaceOp := NewReplaceOperator(iter, mappings, "level", logger)

	err := replaceOp.Open(ctx)
	require.NoError(t, err)

	// First row: error -> ERROR
	row, err := replaceOp.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, row)
	level, _ := row.Get("level")
	assert.Equal(t, "ERROR", level)

	// Second row: warn -> WARNING
	row, err = replaceOp.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, row)
	level, _ = row.Get("level")
	assert.Equal(t, "WARNING", level)

	// Third row: info remains unchanged
	row, err = replaceOp.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, row)
	level, _ = row.Get("level")
	assert.Equal(t, "info", level)

	err = replaceOp.Close()
	assert.NoError(t, err)
}

func TestReplaceOperator_MultipleReplacements(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"text": "The quick brown fox",
		}),
		NewRow(map[string]interface{}{
			"text": "The lazy dog",
		}),
	}

	iter := NewSliceIterator(rows)
	mappings := []*ast.ReplaceMapping{
		{
			OldValue: &ast.Literal{Value: "quick"},
			NewValue: &ast.Literal{Value: "FAST"},
		},
		{
			OldValue: &ast.Literal{Value: "lazy"},
			NewValue: &ast.Literal{Value: "SLOW"},
		},
		{
			OldValue: &ast.Literal{Value: "The"},
			NewValue: &ast.Literal{Value: "A"},
		},
	}
	replaceOp := NewReplaceOperator(iter, mappings, "text", logger)

	err := replaceOp.Open(ctx)
	require.NoError(t, err)

	// First row: "The quick brown fox" -> "A FAST brown fox"
	row, err := replaceOp.Next(ctx)
	require.NoError(t, err)
	text, _ := row.Get("text")
	assert.Equal(t, "A FAST brown fox", text)

	// Second row: "The lazy dog" -> "A SLOW dog"
	row, err = replaceOp.Next(ctx)
	require.NoError(t, err)
	text, _ = row.Get("text")
	assert.Equal(t, "A SLOW dog", text)

	err = replaceOp.Close()
	assert.NoError(t, err)
}

func TestReplaceOperator_FieldNotExists(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"message": "Test message",
		}),
	}

	iter := NewSliceIterator(rows)
	mappings := []*ast.ReplaceMapping{
		{
			OldValue: &ast.Literal{Value: "old"},
			NewValue: &ast.Literal{Value: "new"},
		},
	}
	// Try to replace in non-existent field
	replaceOp := NewReplaceOperator(iter, mappings, "nonexistent", logger)

	err := replaceOp.Open(ctx)
	require.NoError(t, err)

	// Should return row unchanged when field doesn't exist
	row, err := replaceOp.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, row)

	// Field should still not exist
	_, exists := row.Get("nonexistent")
	assert.False(t, exists)

	err = replaceOp.Close()
	assert.NoError(t, err)
}

func TestReplaceOperator_NumericValues(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"status": 200,
		}),
		NewRow(map[string]interface{}{
			"status": 404,
		}),
		NewRow(map[string]interface{}{
			"status": 500,
		}),
	}

	iter := NewSliceIterator(rows)
	mappings := []*ast.ReplaceMapping{
		{
			OldValue: &ast.Literal{Value: int64(404)},
			NewValue: &ast.Literal{Value: "Not Found"},
		},
		{
			OldValue: &ast.Literal{Value: int64(500)},
			NewValue: &ast.Literal{Value: "Server Error"},
		},
	}
	replaceOp := NewReplaceOperator(iter, mappings, "status", logger)

	err := replaceOp.Open(ctx)
	require.NoError(t, err)

	// First row: 200 -> unchanged (converted to "200")
	row, err := replaceOp.Next(ctx)
	require.NoError(t, err)
	status, _ := row.Get("status")
	assert.Equal(t, "200", status)

	// Second row: 404 -> "Not Found"
	row, err = replaceOp.Next(ctx)
	require.NoError(t, err)
	status, _ = row.Get("status")
	assert.Equal(t, "Not Found", status)

	// Third row: 500 -> "Server Error"
	row, err = replaceOp.Next(ctx)
	require.NoError(t, err)
	status, _ = row.Get("status")
	assert.Equal(t, "Server Error", status)

	err = replaceOp.Close()
	assert.NoError(t, err)
}

func TestReplaceOperator_EmptyMappings(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"field": "value",
		}),
	}

	iter := NewSliceIterator(rows)
	mappings := []*ast.ReplaceMapping{} // Empty mappings
	replaceOp := NewReplaceOperator(iter, mappings, "field", logger)

	err := replaceOp.Open(ctx)
	require.NoError(t, err)

	// Should return row unchanged
	row, err := replaceOp.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, row)
	field, _ := row.Get("field")
	assert.Equal(t, "value", field)

	err = replaceOp.Close()
	assert.NoError(t, err)
}

func TestReplaceOperator_Stats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"field": "value1"}),
		NewRow(map[string]interface{}{"field": "value2"}),
		NewRow(map[string]interface{}{"field": "value3"}),
	}

	iter := NewSliceIterator(rows)
	mappings := []*ast.ReplaceMapping{
		{
			OldValue: &ast.Literal{Value: "value1"},
			NewValue: &ast.Literal{Value: "replaced"},
		},
	}
	replaceOp := NewReplaceOperator(iter, mappings, "field", logger)

	err := replaceOp.Open(ctx)
	require.NoError(t, err)

	// Read all rows
	for i := 0; i < 3; i++ {
		_, err := replaceOp.Next(ctx)
		require.NoError(t, err)
	}

	stats := replaceOp.Stats()
	assert.Equal(t, int64(3), stats.RowsRead)
	assert.Equal(t, int64(3), stats.RowsReturned)

	err = replaceOp.Close()
	assert.NoError(t, err)
}
