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

func TestFillnullOperator(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("FillAllFields", func(t *testing.T) {
		// Create test data with NULL values
		rows := []*Row{
			NewRow(map[string]interface{}{"id": 1, "status": nil, "name": "test"}),
			NewRow(map[string]interface{}{"id": 2, "status": "ok", "name": nil}),
			NewRow(map[string]interface{}{"id": 3, "status": nil, "name": nil}),
		}

		input := NewSliceIterator(rows)
		fillValue := &ast.Literal{Value: "N/A"}
		op := NewFillnullOperator(input, fillValue, nil, logger) // nil fields = fill all

		err := op.Open(ctx)
		require.NoError(t, err)

		// Row 1: status=nil -> "N/A", name="test"
		result1, err := op.Next(ctx)
		require.NoError(t, err)
		status, _ := result1.Get("status")
		assert.Equal(t, "N/A", status)
		name, _ := result1.Get("name")
		assert.Equal(t, "test", name)

		// Row 2: status="ok", name=nil -> "N/A"
		result2, err := op.Next(ctx)
		require.NoError(t, err)
		status, _ = result2.Get("status")
		assert.Equal(t, "ok", status)
		name, _ = result2.Get("name")
		assert.Equal(t, "N/A", name)

		// Row 3: both NULL -> both filled
		result3, err := op.Next(ctx)
		require.NoError(t, err)
		status, _ = result3.Get("status")
		assert.Equal(t, "N/A", status)
		name, _ = result3.Get("name")
		assert.Equal(t, "N/A", name)

		// EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("FillSpecificFields", func(t *testing.T) {
		// Create test data
		rows := []*Row{
			NewRow(map[string]interface{}{"id": 1, "cpu": nil, "memory": nil, "disk": nil}),
			NewRow(map[string]interface{}{"id": 2, "cpu": 50.0, "memory": nil, "disk": nil}),
		}

		input := NewSliceIterator(rows)
		fillValue := &ast.Literal{Value: 0.0}
		fields := []ast.Expression{
			&ast.FieldReference{Name: "cpu"},
			&ast.FieldReference{Name: "memory"},
		}
		op := NewFillnullOperator(input, fillValue, fields, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Row 1: cpu=nil -> 0.0, memory=nil -> 0.0, disk=nil (not filled)
		result1, err := op.Next(ctx)
		require.NoError(t, err)
		cpu, _ := result1.Get("cpu")
		assert.Equal(t, 0.0, cpu)
		memory, _ := result1.Get("memory")
		assert.Equal(t, 0.0, memory)
		disk, _ := result1.Get("disk")
		assert.Nil(t, disk) // Not in fields list, stays NULL

		// Row 2: cpu=50.0 (unchanged), memory=nil -> 0.0, disk=nil (not filled)
		result2, err := op.Next(ctx)
		require.NoError(t, err)
		cpu, _ = result2.Get("cpu")
		assert.Equal(t, 50.0, cpu)
		memory, _ = result2.Get("memory")
		assert.Equal(t, 0.0, memory)
		disk, _ = result2.Get("disk")
		assert.Nil(t, disk)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("NumericFillValue", func(t *testing.T) {
		rows := []*Row{
			NewRow(map[string]interface{}{"count": nil, "total": 100}),
			NewRow(map[string]interface{}{"count": 5, "total": nil}),
		}

		input := NewSliceIterator(rows)
		fillValue := &ast.Literal{Value: int64(0)}
		op := NewFillnullOperator(input, fillValue, nil, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Row 1
		result1, err := op.Next(ctx)
		require.NoError(t, err)
		count, _ := result1.Get("count")
		assert.Equal(t, int64(0), count)
		total, _ := result1.Get("total")
		assert.Equal(t, 100, total)

		// Row 2
		result2, err := op.Next(ctx)
		require.NoError(t, err)
		count, _ = result2.Get("count")
		assert.Equal(t, 5, count)
		total, _ = result2.Get("total")
		assert.Equal(t, int64(0), total)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("BooleanFillValue", func(t *testing.T) {
		rows := []*Row{
			NewRow(map[string]interface{}{"enabled": nil, "active": true}),
			NewRow(map[string]interface{}{"enabled": true, "active": nil}),
		}

		input := NewSliceIterator(rows)
		fillValue := &ast.Literal{Value: false}
		op := NewFillnullOperator(input, fillValue, nil, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Row 1
		result1, err := op.Next(ctx)
		require.NoError(t, err)
		enabled, _ := result1.Get("enabled")
		assert.Equal(t, false, enabled)
		active, _ := result1.Get("active")
		assert.Equal(t, true, active)

		// Row 2
		result2, err := op.Next(ctx)
		require.NoError(t, err)
		enabled, _ = result2.Get("enabled")
		assert.Equal(t, true, enabled)
		active, _ = result2.Get("active")
		assert.Equal(t, false, active)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("EmptyInput", func(t *testing.T) {
		rows := []*Row{}
		input := NewSliceIterator(rows)
		fillValue := &ast.Literal{Value: "default"}
		op := NewFillnullOperator(input, fillValue, nil, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Should immediately return EOF
		_, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("NoNullValues", func(t *testing.T) {
		rows := []*Row{
			NewRow(map[string]interface{}{"id": 1, "name": "test1"}),
			NewRow(map[string]interface{}{"id": 2, "name": "test2"}),
		}

		input := NewSliceIterator(rows)
		fillValue := &ast.Literal{Value: "N/A"}
		op := NewFillnullOperator(input, fillValue, nil, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Rows should pass through unchanged
		result1, err := op.Next(ctx)
		require.NoError(t, err)
		name, _ := result1.Get("name")
		assert.Equal(t, "test1", name)

		result2, err := op.Next(ctx)
		require.NoError(t, err)
		name, _ = result2.Get("name")
		assert.Equal(t, "test2", name)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("CreateMissingField", func(t *testing.T) {
		// If a field doesn't exist and is in the fields list, create it
		rows := []*Row{
			NewRow(map[string]interface{}{"id": 1}), // status field missing
			NewRow(map[string]interface{}{"id": 2, "status": "ok"}),
		}

		input := NewSliceIterator(rows)
		fillValue := &ast.Literal{Value: "unknown"}
		fields := []ast.Expression{
			&ast.FieldReference{Name: "status"},
		}
		op := NewFillnullOperator(input, fillValue, fields, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Row 1: status doesn't exist, should be created with fill value
		result1, err := op.Next(ctx)
		require.NoError(t, err)
		status, exists := result1.Get("status")
		assert.True(t, exists)
		assert.Equal(t, "unknown", status)

		// Row 2: status exists, should remain unchanged
		result2, err := op.Next(ctx)
		require.NoError(t, err)
		status, _ = result2.Get("status")
		assert.Equal(t, "ok", status)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("Stats", func(t *testing.T) {
		rows := []*Row{
			NewRow(map[string]interface{}{"field": nil}),
			NewRow(map[string]interface{}{"field": "value"}),
			NewRow(map[string]interface{}{"field": nil}),
		}

		input := NewSliceIterator(rows)
		fillValue := &ast.Literal{Value: "filled"}
		op := NewFillnullOperator(input, fillValue, nil, logger)

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
		// Test with a larger dataset to ensure performance
		numRows := 1000
		rows := make([]*Row, numRows)
		for i := 0; i < numRows; i++ {
			var value interface{}
			if i%3 == 0 {
				value = nil // Every 3rd row has NULL
			} else {
				value = i
			}
			rows[i] = NewRow(map[string]interface{}{"index": i, "value": value})
		}

		input := NewSliceIterator(rows)
		fillValue := &ast.Literal{Value: -1}
		op := NewFillnullOperator(input, fillValue, nil, logger)

		err := op.Open(ctx)
		require.NoError(t, err)

		// Read and verify all rows
		count := 0
		for {
			row, err := op.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)

			value, _ := row.Get("value")
			if count%3 == 0 {
				assert.Equal(t, -1, value, "Row %d should have filled value", count)
			} else {
				assert.Equal(t, count, value, "Row %d should have original value", count)
			}
			count++
		}

		assert.Equal(t, numRows, count, "Should have processed all rows")

		err = op.Close()
		require.NoError(t, err)
	})
}
