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

func TestReverseOperator(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("BasicReverse", func(t *testing.T) {
		// Create test data: rows with sequential IDs
		rows := []*Row{
			NewRow(map[string]interface{}{"id": 1, "value": "first"}),
			NewRow(map[string]interface{}{"id": 2, "value": "second"}),
			NewRow(map[string]interface{}{"id": 3, "value": "third"}),
			NewRow(map[string]interface{}{"id": 4, "value": "fourth"}),
			NewRow(map[string]interface{}{"id": 5, "value": "fifth"}),
		}

		input := NewSliceIterator(rows)
		op := NewReverseOperator(input, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get rows in reverse order
		result1, err := op.Next(ctx)
		require.NoError(t, err)
		id, _ := result1.Get("id")
		assert.Equal(t, 5, id) // Fifth row first

		result2, err := op.Next(ctx)
		require.NoError(t, err)
		id, _ = result2.Get("id")
		assert.Equal(t, 4, id) // Fourth row second

		result3, err := op.Next(ctx)
		require.NoError(t, err)
		id, _ = result3.Get("id")
		assert.Equal(t, 3, id) // Third row third

		result4, err := op.Next(ctx)
		require.NoError(t, err)
		id, _ = result4.Get("id")
		assert.Equal(t, 2, id) // Second row fourth

		result5, err := op.Next(ctx)
		require.NoError(t, err)
		id, _ = result5.Get("id")
		assert.Equal(t, 1, id) // First row last

		// EOF
		result, err := op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, result)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("EmptyInput", func(t *testing.T) {
		// Empty input
		rows := []*Row{}
		input := NewSliceIterator(rows)
		op := NewReverseOperator(input, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should immediately return EOF
		result, err := op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, result)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("SingleRow", func(t *testing.T) {
		// Single row
		rows := []*Row{
			NewRow(map[string]interface{}{"id": 42, "name": "answer"}),
		}

		input := NewSliceIterator(rows)
		op := NewReverseOperator(input, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get the single row
		result, err := op.Next(ctx)
		require.NoError(t, err)
		id, _ := result.Get("id")
		assert.Equal(t, 42, id)

		// EOF
		result, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, result)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("TwoRows", func(t *testing.T) {
		// Two rows
		rows := []*Row{
			NewRow(map[string]interface{}{"order": 1}),
			NewRow(map[string]interface{}{"order": 2}),
		}

		input := NewSliceIterator(rows)
		op := NewReverseOperator(input, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// First should be second row
		result1, err := op.Next(ctx)
		require.NoError(t, err)
		order, _ := result1.Get("order")
		assert.Equal(t, 2, order)

		// Second should be first row
		result2, err := op.Next(ctx)
		require.NoError(t, err)
		order, _ = result2.Get("order")
		assert.Equal(t, 1, order)

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("Stats", func(t *testing.T) {
		// Create test data
		rows := []*Row{
			NewRow(map[string]interface{}{"id": 1}),
			NewRow(map[string]interface{}{"id": 2}),
			NewRow(map[string]interface{}{"id": 3}),
		}

		input := NewSliceIterator(rows)
		op := NewReverseOperator(input, logger)

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

		// Check stats
		stats := op.Stats()
		assert.Equal(t, int64(3), stats.RowsRead, "Should have read 3 rows")
		assert.Equal(t, int64(3), stats.RowsReturned, "Should have returned 3 rows")

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("LargeDataset", func(t *testing.T) {
		// Test with a larger dataset to ensure buffering works correctly
		numRows := 1000
		rows := make([]*Row, numRows)
		for i := 0; i < numRows; i++ {
			rows[i] = NewRow(map[string]interface{}{"index": i})
		}

		input := NewSliceIterator(rows)
		op := NewReverseOperator(input, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Check first few reversed rows
		for expected := numRows - 1; expected >= numRows-5; expected-- {
			result, err := op.Next(ctx)
			require.NoError(t, err)
			index, _ := result.Get("index")
			assert.Equal(t, expected, index)
		}

		// Skip to the end
		count := 5 // Already read 5
		for {
			_, err := op.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)
			count++
		}

		assert.Equal(t, numRows, count, "Should have read all rows")

		err = op.Close()
		require.NoError(t, err)
	})
}
