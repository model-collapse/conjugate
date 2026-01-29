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

func TestRexOperator(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("BasicRex", func(t *testing.T) {
		// Create test data with log messages
		rows := []*Row{
			NewRow(map[string]interface{}{
				"message": "error_code=404 error_msg=Not Found",
			}),
			NewRow(map[string]interface{}{
				"message": "error_code=500 error_msg=Internal Server Error",
			}),
		}

		input := NewSliceIterator(rows)
		pattern := `error_code=(?P<code>\d+) error_msg=(?P<msg>.*)`
		extractedFields := []string{"code", "msg"}

		op, err := NewRexOperator(input, "message", pattern, extractedFields, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		// First row
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		code, exists := row.Get("code")
		assert.True(t, exists)
		assert.Equal(t, "404", code)

		msg, exists := row.Get("msg")
		assert.True(t, exists)
		assert.Equal(t, "Not Found", msg)

		// Second row
		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		code, exists = row.Get("code")
		assert.True(t, exists)
		assert.Equal(t, "500", code)

		msg, exists = row.Get("msg")
		assert.True(t, exists)
		assert.Equal(t, "Internal Server Error", msg)

		// EOF
		row, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, row)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("RexWithDefaultField", func(t *testing.T) {
		// Test rex with empty sourceField (should default to _raw)
		rows := []*Row{
			NewRow(map[string]interface{}{
				"_raw": "user=john action=login",
			}),
		}

		input := NewSliceIterator(rows)
		pattern := `user=(?P<user>\w+) action=(?P<action>\w+)`
		extractedFields := []string{"user", "action"}

		// Pass empty string for sourceField (should default to _raw)
		op, err := NewRexOperator(input, "", pattern, extractedFields, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		user, exists := row.Get("user")
		assert.True(t, exists)
		assert.Equal(t, "john", user)

		action, exists := row.Get("action")
		assert.True(t, exists)
		assert.Equal(t, "login", action)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("NoMatch", func(t *testing.T) {
		// Create test data that doesn't match the pattern
		rows := []*Row{
			NewRow(map[string]interface{}{
				"message": "this message doesn't match the pattern",
			}),
		}

		input := NewSliceIterator(rows)
		pattern := `user=(?P<user>\w+)`
		extractedFields := []string{"user"}

		op, err := NewRexOperator(input, "message", pattern, extractedFields, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		// Should return row without adding new fields
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		_, exists := row.Get("user")
		assert.False(t, exists)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("MissingSourceField", func(t *testing.T) {
		// Create test data without the source field
		rows := []*Row{
			NewRow(map[string]interface{}{
				"other_field": "value",
			}),
		}

		input := NewSliceIterator(rows)
		pattern := `user=(?P<user>\w+)`
		extractedFields := []string{"user"}

		op, err := NewRexOperator(input, "message", pattern, extractedFields, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		// Should return row without adding new fields
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		_, exists := row.Get("user")
		assert.False(t, exists)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("MultipleExtractions", func(t *testing.T) {
		// Test multiple rex operations can be chained
		rows := []*Row{
			NewRow(map[string]interface{}{
				"message": "status=200 latency=150ms user=admin",
			}),
		}

		input := NewSliceIterator(rows)
		pattern := `status=(?P<status>\d+) latency=(?P<latency>\d+)ms user=(?P<user>\w+)`
		extractedFields := []string{"status", "latency", "user"}

		op, err := NewRexOperator(input, "message", pattern, extractedFields, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		status, _ := row.Get("status")
		assert.Equal(t, "200", status)

		latency, _ := row.Get("latency")
		assert.Equal(t, "150", latency)

		user, _ := row.Get("user")
		assert.Equal(t, "admin", user)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("InvalidPattern", func(t *testing.T) {
		rows := []*Row{
			NewRow(map[string]interface{}{
				"message": "test",
			}),
		}

		input := NewSliceIterator(rows)
		invalidPattern := `(?P<invalid[a-z` // Invalid regex

		_, err := NewRexOperator(input, "message", invalidPattern, []string{}, logger)
		assert.Error(t, err)
	})
}
