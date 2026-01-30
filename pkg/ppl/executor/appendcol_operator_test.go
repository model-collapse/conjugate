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

func TestAppendcolOperator_Basic(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Main input
	mainRows := []*Row{
		NewRow(map[string]interface{}{
			"id":   1,
			"name": "Alice",
		}),
		NewRow(map[string]interface{}{
			"id":   2,
			"name": "Bob",
		}),
	}

	// Subsearch results
	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"email": "alice@example.com",
			"dept":  "Engineering",
		}),
		NewRow(map[string]interface{}{
			"email": "bob@example.com",
			"dept":  "Sales",
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchInput := NewSliceIterator(subsearchRows)

	op := NewAppendcolOperator(mainInput, subsearchInput, false, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// First row: id=1, name=Alice, email=alice@example.com, dept=Engineering
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, row1)
	id, _ := row1.Get("id")
	assert.Equal(t, 1, id)
	name, _ := row1.Get("name")
	assert.Equal(t, "Alice", name)
	email, _ := row1.Get("email")
	assert.Equal(t, "alice@example.com", email)
	dept, _ := row1.Get("dept")
	assert.Equal(t, "Engineering", dept)

	// Second row: id=2, name=Bob, email=bob@example.com, dept=Sales
	row2, err := op.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, row2)
	id, _ = row2.Get("id")
	assert.Equal(t, 2, id)
	name, _ = row2.Get("name")
	assert.Equal(t, "Bob", name)
	email, _ = row2.Get("email")
	assert.Equal(t, "bob@example.com", email)
	dept, _ = row2.Get("dept")
	assert.Equal(t, "Sales", dept)

	// EOF
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendcolOperator_EmptySubsearch(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{
			"id": 1,
		}),
		NewRow(map[string]interface{}{
			"id": 2,
		}),
	}

	// Empty subsearch
	subsearchRows := []*Row{}

	mainInput := NewSliceIterator(mainRows)
	subsearchInput := NewSliceIterator(subsearchRows)

	op := NewAppendcolOperator(mainInput, subsearchInput, false, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should return main rows without any additional columns
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ := row1.Get("id")
	assert.Equal(t, 1, id)
	assert.Len(t, row1.Fields(), 1)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ = row2.Get("id")
	assert.Equal(t, 2, id)
	assert.Len(t, row2.Fields(), 1)

	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendcolOperator_EmptyMainInput(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Empty main input
	mainRows := []*Row{}

	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"email": "test@example.com",
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchInput := NewSliceIterator(subsearchRows)

	op := NewAppendcolOperator(mainInput, subsearchInput, false, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get EOF immediately since main input is empty
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendcolOperator_MoreMainRowsThanSubsearch(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// 3 main rows
	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1}),
		NewRow(map[string]interface{}{"id": 2}),
		NewRow(map[string]interface{}{"id": 3}),
	}

	// Only 2 subsearch rows
	subsearchRows := []*Row{
		NewRow(map[string]interface{}{"email": "a@x.com"}),
		NewRow(map[string]interface{}{"email": "b@x.com"}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchInput := NewSliceIterator(subsearchRows)

	op := NewAppendcolOperator(mainInput, subsearchInput, false, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// First row: merged
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ := row1.Get("id")
	assert.Equal(t, 1, id)
	email, _ := row1.Get("email")
	assert.Equal(t, "a@x.com", email)

	// Second row: merged
	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ = row2.Get("id")
	assert.Equal(t, 2, id)
	email, _ = row2.Get("email")
	assert.Equal(t, "b@x.com", email)

	// Third row: no subsearch data to merge
	row3, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ = row3.Get("id")
	assert.Equal(t, 3, id)
	_, hasEmail := row3.Get("email")
	assert.False(t, hasEmail, "Third row should not have email field")

	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendcolOperator_MoreSubsearchRowsThanMain(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Only 2 main rows
	mainRows := []*Row{
		NewRow(map[string]interface{}{"id": 1}),
		NewRow(map[string]interface{}{"id": 2}),
	}

	// 3 subsearch rows (extra one is ignored)
	subsearchRows := []*Row{
		NewRow(map[string]interface{}{"email": "a@x.com"}),
		NewRow(map[string]interface{}{"email": "b@x.com"}),
		NewRow(map[string]interface{}{"email": "c@x.com"}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchInput := NewSliceIterator(subsearchRows)

	op := NewAppendcolOperator(mainInput, subsearchInput, false, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// First row: merged
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ := row1.Get("id")
	assert.Equal(t, 1, id)
	email, _ := row1.Get("email")
	assert.Equal(t, "a@x.com", email)

	// Second row: merged
	row2, err := op.Next(ctx)
	require.NoError(t, err)
	id, _ = row2.Get("id")
	assert.Equal(t, 2, id)
	email, _ = row2.Get("email")
	assert.Equal(t, "b@x.com", email)

	// Third subsearch row is unused (main input exhausted)
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendcolOperator_ColumnConflict_NoOverride(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Main and subsearch have overlapping "status" field
	mainRows := []*Row{
		NewRow(map[string]interface{}{
			"id":     1,
			"status": "active",
		}),
	}

	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"status": "pending",
			"email":  "test@example.com",
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchInput := NewSliceIterator(subsearchRows)

	// override = false: main wins
	op := NewAppendcolOperator(mainInput, subsearchInput, false, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Main's "status" should be preserved
	status, _ := row.Get("status")
	assert.Equal(t, "active", status)
	id, _ := row.Get("id")
	assert.Equal(t, 1, id)
	email, _ := row.Get("email")
	assert.Equal(t, "test@example.com", email)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendcolOperator_ColumnConflict_WithOverride(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Main and subsearch have overlapping "status" field
	mainRows := []*Row{
		NewRow(map[string]interface{}{
			"id":     1,
			"status": "active",
		}),
	}

	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"status": "pending",
			"email":  "test@example.com",
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchInput := NewSliceIterator(subsearchRows)

	// override = true: subsearch wins
	op := NewAppendcolOperator(mainInput, subsearchInput, true, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Subsearch's "status" should override
	status, _ := row.Get("status")
	assert.Equal(t, "pending", status)
	id, _ := row.Get("id")
	assert.Equal(t, 1, id)
	email, _ := row.Get("email")
	assert.Equal(t, "test@example.com", email)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendcolOperator_MultipleColumns(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{
			"id":   1,
			"name": "Alice",
			"age":  30,
		}),
	}

	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"email":   "alice@example.com",
			"dept":    "Engineering",
			"level":   "Senior",
			"manager": "Bob",
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchInput := NewSliceIterator(subsearchRows)

	op := NewAppendcolOperator(mainInput, subsearchInput, false, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Should have all 7 fields
	assert.Len(t, row.Fields(), 7)
	id, _ := row.Get("id")
	assert.Equal(t, 1, id)
	name, _ := row.Get("name")
	assert.Equal(t, "Alice", name)
	age, _ := row.Get("age")
	assert.Equal(t, 30, age)
	email, _ := row.Get("email")
	assert.Equal(t, "alice@example.com", email)
	dept, _ := row.Get("dept")
	assert.Equal(t, "Engineering", dept)
	level, _ := row.Get("level")
	assert.Equal(t, "Senior", level)
	manager, _ := row.Get("manager")
	assert.Equal(t, "Bob", manager)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendcolOperator_TypePreservation(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainRows := []*Row{
		NewRow(map[string]interface{}{
			"id":       1,
			"count":    int64(100),
			"ratio":    0.75,
			"active":   true,
			"metadata": map[string]interface{}{"key": "value"},
		}),
	}

	subsearchRows := []*Row{
		NewRow(map[string]interface{}{
			"score":    float64(98.5),
			"attempts": int32(5),
			"verified": false,
		}),
	}

	mainInput := NewSliceIterator(mainRows)
	subsearchInput := NewSliceIterator(subsearchRows)

	op := NewAppendcolOperator(mainInput, subsearchInput, false, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Verify types are preserved
	id, _ := row.Get("id")
	assert.IsType(t, int(0), id)
	count, _ := row.Get("count")
	assert.IsType(t, int64(0), count)
	ratio, _ := row.Get("ratio")
	assert.IsType(t, float64(0), ratio)
	active, _ := row.Get("active")
	assert.IsType(t, true, active)
	score, _ := row.Get("score")
	assert.IsType(t, float64(0), score)
	attempts, _ := row.Get("attempts")
	assert.IsType(t, int32(0), attempts)
	verified, _ := row.Get("verified")
	assert.IsType(t, false, verified)

	err = op.Close()
	require.NoError(t, err)
}

func TestAppendcolOperator_BothEmpty(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	mainInput := NewSliceIterator([]*Row{})
	subsearchInput := NewSliceIterator([]*Row{})

	op := NewAppendcolOperator(mainInput, subsearchInput, false, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// Should get EOF immediately
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}
