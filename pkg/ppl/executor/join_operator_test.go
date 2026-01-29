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

// Helper function to get field value
func mustGet(row *Row, field string) interface{} {
	val, ok := row.Get(field)
	if !ok {
		panic("field not found: " + field)
	}
	return val
}

func TestJoinOperator(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("InnerJoin", func(t *testing.T) {
		// Left side (orders)
		leftRows := []*Row{
			NewRow(map[string]interface{}{
				"order_id": "1",
				"user_id":  "u1",
				"amount":   "100",
			}),
			NewRow(map[string]interface{}{
				"order_id": "2",
				"user_id":  "u2",
				"amount":   "200",
			}),
			NewRow(map[string]interface{}{
				"order_id": "3",
				"user_id":  "u999", // No matching user
				"amount":   "300",
			}),
		}

		// Right side (users)
		rightRows := []*Row{
			NewRow(map[string]interface{}{
				"user_id": "u1",
				"name":    "Alice",
			}),
			NewRow(map[string]interface{}{
				"user_id": "u2",
				"name":    "Bob",
			}),
		}

		left := NewSliceIterator(leftRows)
		right := NewSliceIterator(rightRows)

		op := NewJoinOperator(left, right, ast.JoinTypeInner, "user_id", "user_id", logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		// First joined row: order 1 + user 1
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)
		assert.Equal(t, "1", mustGet(row, "order_id"))
		assert.Equal(t, "u1", mustGet(row, "user_id"))
		assert.Equal(t, "Alice", mustGet(row, "name"))

		// Second joined row: order 2 + user 2
		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)
		assert.Equal(t, "2", mustGet(row, "order_id"))
		assert.Equal(t, "u2", mustGet(row, "user_id"))
		assert.Equal(t, "Bob", mustGet(row, "name"))

		// Order 3 has no matching user, so it's not returned (inner join)
		// EOF
		row, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, row)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("LeftJoin", func(t *testing.T) {
		// Left side (orders)
		leftRows := []*Row{
			NewRow(map[string]interface{}{
				"order_id": "1",
				"user_id":  "u1",
			}),
			NewRow(map[string]interface{}{
				"order_id": "2",
				"user_id":  "u999", // No matching user
			}),
		}

		// Right side (users)
		rightRows := []*Row{
			NewRow(map[string]interface{}{
				"user_id": "u1",
				"name":    "Alice",
			}),
		}

		left := NewSliceIterator(leftRows)
		right := NewSliceIterator(rightRows)

		op := NewJoinOperator(left, right, ast.JoinTypeLeft, "user_id", "user_id", logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		// First joined row: order 1 + user 1
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)
		assert.Equal(t, "1", mustGet(row, "order_id"))
		assert.Equal(t, "Alice", mustGet(row, "name"))

		// Second row: order 2 with no matching user (left join keeps it)
		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)
		assert.Equal(t, "2", mustGet(row, "order_id"))
		assert.Equal(t, "u999", mustGet(row, "user_id"))
		_, exists := row.Get("name")
		assert.False(t, exists) // No name field (no match)

		// EOF
		row, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("MultipleMatches", func(t *testing.T) {
		// Left side
		leftRows := []*Row{
			NewRow(map[string]interface{}{
				"id":       "1",
				"category": "A",
			}),
		}

		// Right side has multiple rows with category A
		rightRows := []*Row{
			NewRow(map[string]interface{}{
				"category": "A",
				"product":  "P1",
			}),
			NewRow(map[string]interface{}{
				"category": "A",
				"product":  "P2",
			}),
			NewRow(map[string]interface{}{
				"category": "A",
				"product":  "P3",
			}),
		}

		left := NewSliceIterator(leftRows)
		right := NewSliceIterator(rightRows)

		op := NewJoinOperator(left, right, ast.JoinTypeInner, "category", "category", logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get 3 joined rows (1 left * 3 right)
		var products []string
		for i := 0; i < 3; i++ {
			row, err := op.Next(ctx)
			require.NoError(t, err)
			require.NotNil(t, row)
			assert.Equal(t, "1", mustGet(row, "id"))
			product, _ := row.Get("product")
			products = append(products, product.(string))
		}

		assert.ElementsMatch(t, []string{"P1", "P2", "P3"}, products)

		// EOF
		row, err := op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, row)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("FieldNameConflict", func(t *testing.T) {
		// Both sides have "status" field
		leftRows := []*Row{
			NewRow(map[string]interface{}{
				"id":     "1",
				"status": "pending",
			}),
		}

		rightRows := []*Row{
			NewRow(map[string]interface{}{
				"id":     "1",
				"status": "active",
			}),
		}

		left := NewSliceIterator(leftRows)
		right := NewSliceIterator(rightRows)

		op := NewJoinOperator(left, right, ast.JoinTypeInner, "id", "id", logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		// Left side status kept as-is
		assert.Equal(t, "pending", mustGet(row, "status"))
		// Right side status gets "_right" suffix
		assert.Equal(t, "active", mustGet(row, "status_right"))

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("EmptyRight", func(t *testing.T) {
		leftRows := []*Row{
			NewRow(map[string]interface{}{"id": "1"}),
		}
		rightRows := []*Row{}

		left := NewSliceIterator(leftRows)
		right := NewSliceIterator(rightRows)

		// Inner join with empty right = no results
		op := NewJoinOperator(left, right, ast.JoinTypeInner, "id", "id", logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		row, err := op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, row)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("EmptyLeft", func(t *testing.T) {
		leftRows := []*Row{}
		rightRows := []*Row{
			NewRow(map[string]interface{}{"id": "1"}),
		}

		left := NewSliceIterator(leftRows)
		right := NewSliceIterator(rightRows)

		op := NewJoinOperator(left, right, ast.JoinTypeInner, "id", "id", logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		// Empty left = no results
		row, err := op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, row)

		err = op.Close()
		require.NoError(t, err)
	})
}
