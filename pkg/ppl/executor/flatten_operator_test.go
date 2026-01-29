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

func TestFlattenOperator(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("BasicArrayFlatten", func(t *testing.T) {
		// Create test data: single row with array field
		rows := []*Row{
			NewRow(map[string]interface{}{
				"host": "server1",
				"tags": []interface{}{"red", "blue", "green"},
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "tags"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get 3 rows, one for each tag
		row1, err := op.Next(ctx)
		require.NoError(t, err)
		host, _ := row1.Get("host")
		assert.Equal(t, "server1", host)
		tag, _ := row1.Get("tags")
		assert.Equal(t, "red", tag)

		row2, err := op.Next(ctx)
		require.NoError(t, err)
		tag, _ = row2.Get("tags")
		assert.Equal(t, "blue", tag)

		row3, err := op.Next(ctx)
		require.NoError(t, err)
		tag, _ = row3.Get("tags")
		assert.Equal(t, "green", tag)

		// EOF
		result, err := op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, result)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("EmptyArray", func(t *testing.T) {
		// Empty array should produce one row with nil value
		rows := []*Row{
			NewRow(map[string]interface{}{
				"host": "server1",
				"tags": []interface{}{},
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "tags"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get one row with tags=nil
		row, err := op.Next(ctx)
		require.NoError(t, err)
		host, _ := row.Get("host")
		assert.Equal(t, "server1", host)
		tag, exists := row.Get("tags")
		assert.True(t, exists)
		assert.Nil(t, tag)

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("MissingField", func(t *testing.T) {
		// Field doesn't exist - should return row as-is
		rows := []*Row{
			NewRow(map[string]interface{}{
				"host": "server1",
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "tags"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get original row
		row, err := op.Next(ctx)
		require.NoError(t, err)
		host, _ := row.Get("host")
		assert.Equal(t, "server1", host)
		_, exists := row.Get("tags")
		assert.False(t, exists)

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("NullField", func(t *testing.T) {
		// Null field - should return row as-is
		rows := []*Row{
			NewRow(map[string]interface{}{
				"host": "server1",
				"tags": nil,
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "tags"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get original row
		row, err := op.Next(ctx)
		require.NoError(t, err)
		host, _ := row.Get("host")
		assert.Equal(t, "server1", host)
		tag, exists := row.Get("tags")
		assert.True(t, exists)
		assert.Nil(t, tag)

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("NonArrayField", func(t *testing.T) {
		// Non-array field - should return row as-is
		rows := []*Row{
			NewRow(map[string]interface{}{
				"host":  "server1",
				"count": 42,
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "count"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get original row
		row, err := op.Next(ctx)
		require.NoError(t, err)
		count, _ := row.Get("count")
		assert.Equal(t, 42, count)

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("MultipleRows", func(t *testing.T) {
		// Multiple input rows, each with array
		rows := []*Row{
			NewRow(map[string]interface{}{
				"host": "server1",
				"tags": []interface{}{"red", "blue"},
			}),
			NewRow(map[string]interface{}{
				"host": "server2",
				"tags": []interface{}{"green"},
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "tags"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// First row, first tag
		row1, err := op.Next(ctx)
		require.NoError(t, err)
		host, _ := row1.Get("host")
		assert.Equal(t, "server1", host)
		tag, _ := row1.Get("tags")
		assert.Equal(t, "red", tag)

		// First row, second tag
		row2, err := op.Next(ctx)
		require.NoError(t, err)
		host, _ = row2.Get("host")
		assert.Equal(t, "server1", host)
		tag, _ = row2.Get("tags")
		assert.Equal(t, "blue", tag)

		// Second row, only tag
		row3, err := op.Next(ctx)
		require.NoError(t, err)
		host, _ = row3.Get("host")
		assert.Equal(t, "server2", host)
		tag, _ = row3.Get("tags")
		assert.Equal(t, "green", tag)

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("SingleElement", func(t *testing.T) {
		// Single-element array
		rows := []*Row{
			NewRow(map[string]interface{}{
				"host": "server1",
				"tags": []interface{}{"only"},
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "tags"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get one row
		row, err := op.Next(ctx)
		require.NoError(t, err)
		tag, _ := row.Get("tags")
		assert.Equal(t, "only", tag)

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("NestedObjects", func(t *testing.T) {
		// Array of objects (common JSON pattern)
		rows := []*Row{
			NewRow(map[string]interface{}{
				"user": "alice",
				"items": []interface{}{
					map[string]interface{}{"id": 1, "name": "apple"},
					map[string]interface{}{"id": 2, "name": "banana"},
				},
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "items"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// First item
		row1, err := op.Next(ctx)
		require.NoError(t, err)
		user, _ := row1.Get("user")
		assert.Equal(t, "alice", user)
		item, _ := row1.Get("items")
		itemMap, ok := item.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, 1, itemMap["id"])
		assert.Equal(t, "apple", itemMap["name"])

		// Second item
		row2, err := op.Next(ctx)
		require.NoError(t, err)
		item, _ = row2.Get("items")
		itemMap, ok = item.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, 2, itemMap["id"])
		assert.Equal(t, "banana", itemMap["name"])

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("Stats", func(t *testing.T) {
		// Test statistics tracking
		rows := []*Row{
			NewRow(map[string]interface{}{
				"tags": []interface{}{"a", "b", "c"},
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "tags"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Read all rows
		for {
			_, err := op.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)
		}

		stats := op.Stats()
		assert.Equal(t, int64(1), stats.RowsRead, "Should have read 1 input row")
		assert.Equal(t, int64(3), stats.RowsReturned, "Should have returned 3 flattened rows")

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("DifferentTypes", func(t *testing.T) {
		// Array with different types
		rows := []*Row{
			NewRow(map[string]interface{}{
				"host": "server1",
				"values": []interface{}{
					"string",
					42,
					3.14,
					true,
					nil,
				},
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "values"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Check each type
		row1, err := op.Next(ctx)
		require.NoError(t, err)
		val, _ := row1.Get("values")
		assert.Equal(t, "string", val)

		row2, err := op.Next(ctx)
		require.NoError(t, err)
		val, _ = row2.Get("values")
		assert.Equal(t, 42, val)

		row3, err := op.Next(ctx)
		require.NoError(t, err)
		val, _ = row3.Get("values")
		assert.Equal(t, 3.14, val)

		row4, err := op.Next(ctx)
		require.NoError(t, err)
		val, _ = row4.Get("values")
		assert.Equal(t, true, val)

		row5, err := op.Next(ctx)
		require.NoError(t, err)
		val, _ = row5.Get("values")
		assert.Nil(t, val)

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("LargeArray", func(t *testing.T) {
		// Test with larger array
		largeArray := make([]interface{}, 100)
		for i := 0; i < 100; i++ {
			largeArray[i] = i
		}

		rows := []*Row{
			NewRow(map[string]interface{}{
				"numbers": largeArray,
			}),
		}

		input := NewSliceIterator(rows)
		field := &ast.FieldReference{Name: "numbers"}
		op := NewFlattenOperator(input, field, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Read all rows and verify count
		count := 0
		for {
			_, err := op.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)
			count++
		}

		assert.Equal(t, 100, count, "Should have 100 flattened rows")

		stats := op.Stats()
		assert.Equal(t, int64(1), stats.RowsRead)
		assert.Equal(t, int64(100), stats.RowsReturned)

		err = op.Close()
		require.NoError(t, err)
	})
}
