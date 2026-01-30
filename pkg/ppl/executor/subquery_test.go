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

// Test IN Subquery

func TestSubqueryIn_Basic(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Main data
	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "status": "active"}),
		NewRow(map[string]interface{}{"id": 2, "status": "pending"}),
		NewRow(map[string]interface{}{"id": 3, "status": "active"}),
		NewRow(map[string]interface{}{"id": 4, "status": "inactive"}),
	}

	// Subquery: valid statuses
	subqueryRows := []*Row{
		NewRow(map[string]interface{}{"status": "active"}),
		NewRow(map[string]interface{}{"status": "pending"}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryInConfig{
		FieldName: "status",
		Subsearch: subsearch,
		Negate:    false,
	}

	op := NewSubqueryInOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get rows with status in (active, pending)
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id1, _ := row1.Get("id")
	assert.Equal(t, 1, id1)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id2, _ := row2.Get("id")
	assert.Equal(t, 2, id2)

	row3, err := op.Next(ctx)
	require.NoError(t, err)
	id3, _ := row3.Get("id")
	assert.Equal(t, 3, id3)

	// No more rows (id=4 with "inactive" filtered out)
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSubqueryIn_NotIn(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "status": "active"}),
		NewRow(map[string]interface{}{"id": 2, "status": "pending"}),
		NewRow(map[string]interface{}{"id": 3, "status": "inactive"}),
	}

	subqueryRows := []*Row{
		NewRow(map[string]interface{}{"status": "active"}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryInConfig{
		FieldName: "status",
		Subsearch: subsearch,
		Negate:    true, // NOT IN
	}

	op := NewSubqueryInOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get rows NOT in (active) = pending, inactive
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id1, _ := row1.Get("id")
	assert.Equal(t, 2, id1)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id2, _ := row2.Get("id")
	assert.Equal(t, 3, id2)

	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSubqueryIn_EmptySubquery(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "status": "active"}),
		NewRow(map[string]interface{}{"id": 2, "status": "pending"}),
	}

	// Empty subquery
	subqueryRows := []*Row{}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryInConfig{
		FieldName: "status",
		Subsearch: subsearch,
	}

	op := NewSubqueryInOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Empty subquery means nothing matches
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSubqueryIn_NumericValues(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "code": 200}),
		NewRow(map[string]interface{}{"id": 2, "code": 404}),
		NewRow(map[string]interface{}{"id": 3, "code": 500}),
		NewRow(map[string]interface{}{"id": 4, "code": 200}),
	}

	subqueryRows := []*Row{
		NewRow(map[string]interface{}{"code": 200}),
		NewRow(map[string]interface{}{"code": 201}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryInConfig{
		FieldName: "code",
		Subsearch: subsearch,
	}

	op := NewSubqueryInOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get rows with code in (200, 201)
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id1, _ := row1.Get("id")
	assert.Equal(t, 1, id1)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id2, _ := row2.Get("id")
	assert.Equal(t, 4, id2)

	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

// Test EXISTS Subquery

func TestSubqueryExists_HasResults(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "name": "Alice"}),
		NewRow(map[string]interface{}{"id": 2, "name": "Bob"}),
	}

	// Subquery returns results
	subqueryRows := []*Row{
		NewRow(map[string]interface{}{"value": "something"}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryExistsConfig{
		Subsearch: subsearch,
	}

	op := NewSubqueryExistsOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// EXISTS is true, should return all rows
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id1, _ := row1.Get("id")
	assert.Equal(t, 1, id1)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id2, _ := row2.Get("id")
	assert.Equal(t, 2, id2)

	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSubqueryExists_NoResults(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "name": "Alice"}),
		NewRow(map[string]interface{}{"id": 2, "name": "Bob"}),
	}

	// Subquery returns no results
	subqueryRows := []*Row{}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryExistsConfig{
		Subsearch: subsearch,
	}

	op := NewSubqueryExistsOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// EXISTS is false, should return no rows
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSubqueryExists_NotExists(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "name": "Alice"}),
	}

	// Subquery returns results
	subqueryRows := []*Row{
		NewRow(map[string]interface{}{"value": "something"}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryExistsConfig{
		Subsearch: subsearch,
		Negate:    true, // NOT EXISTS
	}

	op := NewSubqueryExistsOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// NOT EXISTS is false (subquery has results), should return no rows
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

// Test Scalar Subquery

func TestSubqueryScalar_GreaterThan(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "revenue": 100}),
		NewRow(map[string]interface{}{"id": 2, "revenue": 200}),
		NewRow(map[string]interface{}{"id": 3, "revenue": 50}),
		NewRow(map[string]interface{}{"id": 4, "revenue": 300}),
	}

	// Subquery returns threshold: 150
	subqueryRows := []*Row{
		NewRow(map[string]interface{}{"threshold": 150}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryScalarConfig{
		FieldName:    "revenue",
		Subsearch:    subsearch,
		ComparisonOp: ">",
	}

	op := NewScalarSubqueryOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get rows where revenue > 150
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id1, _ := row1.Get("id")
	assert.Equal(t, 2, id1) // revenue=200

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id2, _ := row2.Get("id")
	assert.Equal(t, 4, id2) // revenue=300

	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSubqueryScalar_Equals(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "status": "A"}),
		NewRow(map[string]interface{}{"id": 2, "status": "B"}),
		NewRow(map[string]interface{}{"id": 3, "status": "A"}),
	}

	subqueryRows := []*Row{
		NewRow(map[string]interface{}{"target": "A"}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryScalarConfig{
		FieldName:    "status",
		Subsearch:    subsearch,
		ComparisonOp: "=",
	}

	op := NewScalarSubqueryOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get rows where status = "A"
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id1, _ := row1.Get("id")
	assert.Equal(t, 1, id1)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id2, _ := row2.Get("id")
	assert.Equal(t, 3, id2)

	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSubqueryScalar_LessThanOrEqual(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "score": 85}),
		NewRow(map[string]interface{}{"id": 2, "score": 90}),
		NewRow(map[string]interface{}{"id": 3, "score": 78}),
	}

	subqueryRows := []*Row{
		NewRow(map[string]interface{}{"max_score": 85}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryScalarConfig{
		FieldName:    "score",
		Subsearch:    subsearch,
		ComparisonOp: "<=",
	}

	op := NewScalarSubqueryOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get rows where score <= 85
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id1, _ := row1.Get("id")
	assert.Equal(t, 1, id1) // score=85

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id2, _ := row2.Get("id")
	assert.Equal(t, 3, id2) // score=78

	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSubqueryScalar_EmptySubquery(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "value": 100}),
	}

	// Empty subquery
	subqueryRows := []*Row{}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryScalarConfig{
		FieldName:    "value",
		Subsearch:    subsearch,
		ComparisonOp: ">",
	}

	op := NewScalarSubqueryOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Scalar subquery with no results is invalid, should return no rows
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSubqueryScalar_TypeConversion(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1, "value": int32(100)}),
		NewRow(map[string]interface{}{"id": 2, "value": float64(200.5)}),
		NewRow(map[string]interface{}{"id": 3, "value": int64(50)}),
	}

	subqueryRows := []*Row{
		NewRow(map[string]interface{}{"threshold": float64(150.0)}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearch := NewSliceIterator(subqueryRows)

	config := SubqueryScalarConfig{
		FieldName:    "value",
		Subsearch:    subsearch,
		ComparisonOp: ">",
	}

	op := NewScalarSubqueryOperator(mainInput, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should handle type conversion (int32, float64, int64 all compared as float64)
	row, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ := row.Get("id")
	assert.Equal(t, 2, id) // value=200.5

	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

// Test SubqueryExecutor directly

func TestSubqueryExecutor_MaxRows(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Create more rows than the limit
	rows := make([]*Row, 150)
	for i := 0; i < 150; i++ {
		rows[i] = NewRow(map[string]interface{}{"id": i})
	}

	subsearch := NewSliceIterator(rows)
	executor := NewSubqueryExecutor(subsearch, SubqueryTypeIN, logger)
	executor.SetMaxRows(100)

	err := executor.Execute(ctx)
	require.NoError(t, err)

	// Should only materialize 100 rows
	results := executor.GetResults()
	assert.Len(t, results, 100)
}

func TestSubqueryExecutor_GetFieldValues(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"status": "active", "other": "data"}),
		NewRow(map[string]interface{}{"status": "pending", "other": "data"}),
		NewRow(map[string]interface{}{"status": "active", "other": "data"}),
	}

	subsearch := NewSliceIterator(rows)
	executor := NewSubqueryExecutor(subsearch, SubqueryTypeIN, logger)

	err := executor.Execute(ctx)
	require.NoError(t, err)

	// Get field values
	values := executor.GetFieldValues("status")
	assert.Len(t, values, 3)
	assert.Contains(t, values, "active")
	assert.Contains(t, values, "pending")
}
