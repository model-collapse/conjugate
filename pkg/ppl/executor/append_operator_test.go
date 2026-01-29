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

func TestAppendOperator(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("BasicAppend", func(t *testing.T) {
		// Create main input data
		mainRows := []*Row{
			NewRow(map[string]interface{}{
				"id":   "1",
				"type": "main",
			}),
			NewRow(map[string]interface{}{
				"id":   "2",
				"type": "main",
			}),
		}

		// Create subsearch data
		subsearchRows := []*Row{
			NewRow(map[string]interface{}{
				"id":   "3",
				"type": "sub",
			}),
			NewRow(map[string]interface{}{
				"id":   "4",
				"type": "sub",
			}),
		}

		mainInput := NewSliceIterator(mainRows)
		subsearchInput := NewSliceIterator(subsearchRows)

		op := NewAppendOperator(mainInput, subsearchInput, logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		// First two rows should be from main input
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)
		rowType, _ := row.Get("type")
		assert.Equal(t, "main", rowType)

		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)
		rowType, _ = row.Get("type")
		assert.Equal(t, "main", rowType)

		// Next two rows should be from subsearch
		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)
		rowType, _ = row.Get("type")
		assert.Equal(t, "sub", rowType)

		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)
		rowType, _ = row.Get("type")
		assert.Equal(t, "sub", rowType)

		// EOF
		row, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, row)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("EmptyMainInput", func(t *testing.T) {
		// Empty main input
		mainRows := []*Row{}

		// Subsearch has data
		subsearchRows := []*Row{
			NewRow(map[string]interface{}{
				"id": "1",
			}),
			NewRow(map[string]interface{}{
				"id": "2",
			}),
		}

		mainInput := NewSliceIterator(mainRows)
		subsearchInput := NewSliceIterator(subsearchRows)

		op := NewAppendOperator(mainInput, subsearchInput, logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get rows from subsearch
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		// EOF
		row, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("EmptySubsearch", func(t *testing.T) {
		// Main input has data
		mainRows := []*Row{
			NewRow(map[string]interface{}{
				"id": "1",
			}),
			NewRow(map[string]interface{}{
				"id": "2",
			}),
		}

		// Empty subsearch
		subsearchRows := []*Row{}

		mainInput := NewSliceIterator(mainRows)
		subsearchInput := NewSliceIterator(subsearchRows)

		op := NewAppendOperator(mainInput, subsearchInput, logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get rows from main input
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		// EOF after main input exhausted (subsearch is empty)
		row, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("BothEmpty", func(t *testing.T) {
		mainInput := NewSliceIterator([]*Row{})
		subsearchInput := NewSliceIterator([]*Row{})

		op := NewAppendOperator(mainInput, subsearchInput, logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		// Should get EOF immediately
		row, err := op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, row)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("DifferentSchemas", func(t *testing.T) {
		// Main input has fields A and B
		mainRows := []*Row{
			NewRow(map[string]interface{}{
				"fieldA": "value1",
				"fieldB": "value2",
			}),
		}

		// Subsearch has fields B and C (B overlaps, C is new)
		subsearchRows := []*Row{
			NewRow(map[string]interface{}{
				"fieldB": "value3",
				"fieldC": "value4",
			}),
		}

		mainInput := NewSliceIterator(mainRows)
		subsearchInput := NewSliceIterator(subsearchRows)

		op := NewAppendOperator(mainInput, subsearchInput, logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		// First row from main input
		row, err := op.Next(ctx)
		require.NoError(t, err)
		valA, exists := row.Get("fieldA")
		assert.True(t, exists)
		assert.Equal(t, "value1", valA)
		valB, exists := row.Get("fieldB")
		assert.True(t, exists)
		assert.Equal(t, "value2", valB)
		_, exists = row.Get("fieldC")
		assert.False(t, exists) // fieldC doesn't exist in main input

		// Second row from subsearch
		row, err = op.Next(ctx)
		require.NoError(t, err)
		_, exists = row.Get("fieldA")
		assert.False(t, exists) // fieldA doesn't exist in subsearch
		valB, exists = row.Get("fieldB")
		assert.True(t, exists)
		assert.Equal(t, "value3", valB)
		valC, exists := row.Get("fieldC")
		assert.True(t, exists)
		assert.Equal(t, "value4", valC)

		err = op.Close()
		require.NoError(t, err)
	})
}
