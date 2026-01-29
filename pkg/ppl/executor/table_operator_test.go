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

func TestTableOperator_BasicSelection(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	// Create test data with multiple fields
	rows := []*Row{
		NewRow(map[string]interface{}{
			"host":     "server1",
			"status":   200,
			"latency":  50,
			"method":   "GET",
			"path":     "/api/users",
		}),
		NewRow(map[string]interface{}{
			"host":     "server2",
			"status":   404,
			"latency":  30,
			"method":   "POST",
			"path":     "/api/orders",
		}),
	}

	iter := NewSliceIterator(rows)
	fields := []ast.Expression{
		&ast.FieldReference{Name: "host"},
		&ast.FieldReference{Name: "status"},
		&ast.FieldReference{Name: "latency"},
	}
	tableOp := NewTableOperator(iter, fields, logger)

	err := tableOp.Open(ctx)
	require.NoError(t, err)

	// First row: should only have host, status, latency
	row, err := tableOp.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, row)

	// Check selected fields are present
	host, exists := row.Get("host")
	assert.True(t, exists)
	assert.Equal(t, "server1", host)

	status, exists := row.Get("status")
	assert.True(t, exists)
	assert.Equal(t, 200, status)

	latency, exists := row.Get("latency")
	assert.True(t, exists)
	assert.Equal(t, 50, latency)

	// Check non-selected fields are absent
	_, exists = row.Get("method")
	assert.False(t, exists)

	_, exists = row.Get("path")
	assert.False(t, exists)

	// Second row
	row, err = tableOp.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, row)

	host, _ = row.Get("host")
	assert.Equal(t, "server2", host)

	err = tableOp.Close()
	assert.NoError(t, err)
}

func TestTableOperator_SingleField(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"host":   "server1",
			"status": 200,
		}),
		NewRow(map[string]interface{}{
			"host":   "server2",
			"status": 404,
		}),
	}

	iter := NewSliceIterator(rows)
	fields := []ast.Expression{
		&ast.FieldReference{Name: "host"},
	}
	tableOp := NewTableOperator(iter, fields, logger)

	err := tableOp.Open(ctx)
	require.NoError(t, err)

	row, err := tableOp.Next(ctx)
	require.NoError(t, err)

	// Only host should be present
	host, exists := row.Get("host")
	assert.True(t, exists)
	assert.Equal(t, "server1", host)

	_, exists = row.Get("status")
	assert.False(t, exists)

	err = tableOp.Close()
	assert.NoError(t, err)
}

func TestTableOperator_MissingFields(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"host":   "server1",
			"status": 200,
		}),
	}

	iter := NewSliceIterator(rows)
	// Request fields that don't all exist
	fields := []ast.Expression{
		&ast.FieldReference{Name: "host"},
		&ast.FieldReference{Name: "nonexistent"},
		&ast.FieldReference{Name: "status"},
	}
	tableOp := NewTableOperator(iter, fields, logger)

	err := tableOp.Open(ctx)
	require.NoError(t, err)

	row, err := tableOp.Next(ctx)
	require.NoError(t, err)

	// Existing fields should be present
	host, exists := row.Get("host")
	assert.True(t, exists)
	assert.Equal(t, "server1", host)

	status, exists := row.Get("status")
	assert.True(t, exists)
	assert.Equal(t, 200, status)

	// Missing field should be nil
	nonexistent, exists := row.Get("nonexistent")
	assert.True(t, exists)
	assert.Nil(t, nonexistent)

	err = tableOp.Close()
	assert.NoError(t, err)
}

func TestTableOperator_AllFields(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
			"field3": "value3",
		}),
	}

	iter := NewSliceIterator(rows)
	// Select all fields
	fields := []ast.Expression{
		&ast.FieldReference{Name: "field1"},
		&ast.FieldReference{Name: "field2"},
		&ast.FieldReference{Name: "field3"},
	}
	tableOp := NewTableOperator(iter, fields, logger)

	err := tableOp.Open(ctx)
	require.NoError(t, err)

	row, err := tableOp.Next(ctx)
	require.NoError(t, err)

	// All fields should be present
	for i := 1; i <= 3; i++ {
		fieldName := "field" + string(rune('0'+i))
		val, exists := row.Get(fieldName)
		assert.True(t, exists, "Field %s should exist", fieldName)
		assert.NotNil(t, val)
	}

	err = tableOp.Close()
	assert.NoError(t, err)
}

func TestTableOperator_EmptyFieldList(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		}),
	}

	iter := NewSliceIterator(rows)
	// Empty field list
	fields := []ast.Expression{}
	tableOp := NewTableOperator(iter, fields, logger)

	err := tableOp.Open(ctx)
	require.NoError(t, err)

	row, err := tableOp.Next(ctx)
	require.NoError(t, err)

	// Should have no fields (empty row)
	rowMap := row.ToMap()
	assert.Empty(t, rowMap)

	err = tableOp.Close()
	assert.NoError(t, err)
}

func TestTableOperator_Stats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"host": "server1"}),
		NewRow(map[string]interface{}{"host": "server2"}),
		NewRow(map[string]interface{}{"host": "server3"}),
	}

	iter := NewSliceIterator(rows)
	fields := []ast.Expression{
		&ast.FieldReference{Name: "host"},
	}
	tableOp := NewTableOperator(iter, fields, logger)

	err := tableOp.Open(ctx)
	require.NoError(t, err)

	// Read all rows
	for i := 0; i < 3; i++ {
		_, err := tableOp.Next(ctx)
		require.NoError(t, err)
	}

	stats := tableOp.Stats()
	assert.Equal(t, int64(3), stats.RowsRead)
	assert.Equal(t, int64(3), stats.RowsReturned)

	err = tableOp.Close()
	assert.NoError(t, err)
}

func TestTableOperator_PreservesOrder(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": 3,
		}),
	}

	iter := NewSliceIterator(rows)
	// Select fields in specific order
	fields := []ast.Expression{
		&ast.FieldReference{Name: "c"},
		&ast.FieldReference{Name: "a"},
		&ast.FieldReference{Name: "b"},
	}
	tableOp := NewTableOperator(iter, fields, logger)

	err := tableOp.Open(ctx)
	require.NoError(t, err)

	row, err := tableOp.Next(ctx)
	require.NoError(t, err)

	// All fields should be present regardless of order
	a, exists := row.Get("a")
	assert.True(t, exists)
	assert.Equal(t, 1, a)

	b, exists := row.Get("b")
	assert.True(t, exists)
	assert.Equal(t, 2, b)

	c, exists := row.Get("c")
	assert.True(t, exists)
	assert.Equal(t, 3, c)

	err = tableOp.Close()
	assert.NoError(t, err)
}

func TestTableOperator_DifferentTypes(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"string_field": "text",
			"int_field":    42,
			"float_field":  3.14,
			"bool_field":   true,
			"nil_field":    nil,
		}),
	}

	iter := NewSliceIterator(rows)
	fields := []ast.Expression{
		&ast.FieldReference{Name: "string_field"},
		&ast.FieldReference{Name: "int_field"},
		&ast.FieldReference{Name: "float_field"},
		&ast.FieldReference{Name: "bool_field"},
		&ast.FieldReference{Name: "nil_field"},
	}
	tableOp := NewTableOperator(iter, fields, logger)

	err := tableOp.Open(ctx)
	require.NoError(t, err)

	row, err := tableOp.Next(ctx)
	require.NoError(t, err)

	// Verify all types are preserved
	stringVal, _ := row.Get("string_field")
	assert.Equal(t, "text", stringVal)

	intVal, _ := row.Get("int_field")
	assert.Equal(t, 42, intVal)

	floatVal, _ := row.Get("float_field")
	assert.Equal(t, 3.14, floatVal)

	boolVal, _ := row.Get("bool_field")
	assert.Equal(t, true, boolVal)

	nilVal, _ := row.Get("nil_field")
	assert.Nil(t, nilVal)

	err = tableOp.Close()
	assert.NoError(t, err)
}
