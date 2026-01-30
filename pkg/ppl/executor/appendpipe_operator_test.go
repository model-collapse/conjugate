// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAppendpipeOperator_Basic(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Main input: 3 rows
	mainRows := []*Row{
		NewRow(map[string]interface{}{
			"category": "A",
			"count":    10,
		}),
		NewRow(map[string]interface{}{
			"category": "B",
			"count":    20,
		}),
		NewRow(map[string]interface{}{
			"category": "C",
			"count":    30,
		}),
	}

	// Subsearch will add a summary row (in real use, it would compute from mainRows)
	// For this test, we simulate that the subsearch produces a total row
	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"category": "Total",
			"count":    60,
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchPlan := NewSliceIterator(subsearchRows)

	op := NewAppendpipeOperator(mainInput, subsearchPlan, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// First 3 rows should be from main input
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	cat, _ := row1.Get("category")
	assert.Equal(t, "A", cat)
	count, _ := row1.Get("count")
	assert.Equal(t, 10, count)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	cat, _ = row2.Get("category")
	assert.Equal(t, "B", cat)
	count, _ = row2.Get("count")
	assert.Equal(t, 20, count)

	row3, err := op.Next(ctx)
	require.NoError(t, err)
	cat, _ = row3.Get("category")
	assert.Equal(t, "C", cat)
	count, _ = row3.Get("count")
	assert.Equal(t, 30, count)

	// Next row should be from subsearch (summary)
	row4, err := op.Next(ctx)
	require.NoError(t, err)
	cat, _ = row4.Get("category")
	assert.Equal(t, "Total", cat)
	count, _ = row4.Get("count")
	assert.Equal(t, 60, count)

	// EOF
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendpipeOperator_EmptyInput(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Empty main input
	mainInput := NewSliceIterator([]*Row{})

	// Subsearch also empty (no input to process)
	subsearchPlan := NewSliceIterator([]*Row{})

	op := NewAppendpipeOperator(mainInput, subsearchPlan, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get EOF immediately
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendpipeOperator_EmptySubsearch(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1}),
		NewRow(map[string]interface{}{"id": 2}),
	}

	mainInput := NewSliceIterator(mainRows)
	// Empty subsearch (produces no additional rows)
	subsearchPlan := NewSliceIterator([]*Row{})

	op := NewAppendpipeOperator(mainInput, subsearchPlan, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get original 2 rows
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ := row1.Get("id")
	assert.Equal(t, 1, id)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ = row2.Get("id")
	assert.Equal(t, 2, id)

	// Then EOF (no subsearch results)
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendpipeOperator_MultipleSubsearchRows(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"type": "data", "value": 100}),
		NewRow(map[string]interface{}{"type": "data", "value": 200}),
	}

	// Subsearch adds multiple summary rows
	subsearchRows := []*Row{
		NewRow(map[string]interface{}{"type": "summary", "total": 300}),
		NewRow(map[string]interface{}{"type": "summary", "avg": 150}),
		NewRow(map[string]interface{}{"type": "summary", "count": 2}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchPlan := NewSliceIterator(subsearchRows)

	op := NewAppendpipeOperator(mainInput, subsearchPlan, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// First 2 rows from main
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	typ, _ := row1.Get("type")
	assert.Equal(t, "data", typ)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	typ, _ = row2.Get("type")
	assert.Equal(t, "data", typ)

	// Next 3 rows from subsearch
	row3, err := op.Next(ctx)
	require.NoError(t, err)
	typ, _ = row3.Get("type")
	assert.Equal(t, "summary", typ)
	total, _ := row3.Get("total")
	assert.Equal(t, 300, total)

	row4, err := op.Next(ctx)
	require.NoError(t, err)
	avg, _ := row4.Get("avg")
	assert.Equal(t, 150, avg)

	row5, err := op.Next(ctx)
	require.NoError(t, err)
	count, _ := row5.Get("count")
	assert.Equal(t, 2, count)

	// EOF
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendpipeOperator_DifferentSchemas(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Main rows have fields A and B
	mainRows := []*Row{
		NewRow(map[string]interface{}{
			"fieldA": "value1",
			"fieldB": 100,
		}),
	}

	// Subsearch rows have different fields (C and D)
	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"fieldC": "summary",
			"fieldD": 9999,
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchPlan := NewSliceIterator(subsearchRows)

	op := NewAppendpipeOperator(mainInput, subsearchPlan, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// First row has A and B
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	valA, exists := row1.Get("fieldA")
	assert.True(t, exists)
	assert.Equal(t, "value1", valA)
	valB, exists := row1.Get("fieldB")
	assert.True(t, exists)
	assert.Equal(t, 100, valB)

	// Second row has C and D (but not A and B)
	row2, err := op.Next(ctx)
	require.NoError(t, err)
	valC, exists := row2.Get("fieldC")
	assert.True(t, exists)
	assert.Equal(t, "summary", valC)
	valD, exists := row2.Get("fieldD")
	assert.True(t, exists)
	assert.Equal(t, 9999, valD)
	_, exists = row2.Get("fieldA")
	assert.False(t, exists)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendpipeOperator_LargeInput(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Create 100 main rows
	mainRows := make([]*Row, 100)
	for i := 0; i < 100; i++ {
		mainRows[i] = NewRow(map[string]interface{}{
			"id":    i,
			"value": i * 10,
		})
	}

	// Subsearch adds 1 summary row
	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"id":    -1,
			"total": 49500, // sum of 0+10+20+...+990
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchPlan := NewSliceIterator(subsearchRows)

	op := NewAppendpipeOperator(mainInput, subsearchPlan, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Read all 100 main rows
	for i := 0; i < 100; i++ {
		row, err := op.Next(ctx)
		require.NoError(t, err)
		id, _ := row.Get("id")
		assert.Equal(t, i, id)
	}

	// Then get the summary row
	summaryRow, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ := summaryRow.Get("id")
	assert.Equal(t, -1, id)
	total, _ := summaryRow.Get("total")
	assert.Equal(t, 49500, total)

	// EOF
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendpipeOperator_TypePreservation(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{
			"int_val":   int32(42),
			"float_val": float64(3.14),
			"bool_val":  true,
			"str_val":   "hello",
		}),
	}

	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"int_val":   int64(999),
			"float_val": float32(2.71),
			"bool_val":  false,
			"str_val":   "world",
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchPlan := NewSliceIterator(subsearchRows)

	op := NewAppendpipeOperator(mainInput, subsearchPlan, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Main row types
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	intVal, _ := row1.Get("int_val")
	assert.IsType(t, int32(0), intVal)
	floatVal, _ := row1.Get("float_val")
	assert.IsType(t, float64(0), floatVal)
	boolVal, _ := row1.Get("bool_val")
	assert.IsType(t, true, boolVal)
	strVal, _ := row1.Get("str_val")
	assert.IsType(t, "", strVal)

	// Subsearch row types
	row2, err := op.Next(ctx)
	require.NoError(t, err)
	intVal, _ = row2.Get("int_val")
	assert.IsType(t, int64(0), intVal)
	floatVal, _ = row2.Get("float_val")
	assert.IsType(t, float32(0), floatVal)

	err = op.Close()
	require.NoError(t, err)
}
