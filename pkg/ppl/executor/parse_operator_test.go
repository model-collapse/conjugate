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

func TestParseOperator(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("BasicParse", func(t *testing.T) {
		// Create test data with log messages
		rows := []*Row{
			NewRow(map[string]interface{}{
				"message": "user john logged in from 192.168.1.100",
			}),
			NewRow(map[string]interface{}{
				"message": "user jane logged in from 10.0.0.5",
			}),
		}

		input := NewSliceIterator(rows)
		pattern := `user (?P<username>\w+) logged in from (?P<ip>\d+\.\d+\.\d+\.\d+)`
		extractedFields := []string{"username", "ip"}

		op, err := NewParseOperator(input, "message", pattern, extractedFields, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		// First row
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		username, exists := row.Get("username")
		assert.True(t, exists)
		assert.Equal(t, "john", username)

		ip, exists := row.Get("ip")
		assert.True(t, exists)
		assert.Equal(t, "192.168.1.100", ip)

		// Second row
		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		username, exists = row.Get("username")
		assert.True(t, exists)
		assert.Equal(t, "jane", username)

		ip, exists = row.Get("ip")
		assert.True(t, exists)
		assert.Equal(t, "10.0.0.5", ip)

		// EOF
		row, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, row)

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
		pattern := `user (?P<username>\w+) logged in from (?P<ip>\d+\.\d+\.\d+\.\d+)`
		extractedFields := []string{"username", "ip"}

		op, err := NewParseOperator(input, "message", pattern, extractedFields, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		// Should return row without adding new fields
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		_, exists := row.Get("username")
		assert.False(t, exists)

		_, exists = row.Get("ip")
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
		pattern := `user (?P<username>\w+)`
		extractedFields := []string{"username"}

		op, err := NewParseOperator(input, "message", pattern, extractedFields, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		// Should return row without adding new fields
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		_, exists := row.Get("username")
		assert.False(t, exists)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("ComplexPattern", func(t *testing.T) {
		// Test with more complex pattern
		rows := []*Row{
			NewRow(map[string]interface{}{
				"log": "[2024-01-15 10:30:45] ERROR in module auth: Authentication failed for user admin",
			}),
		}

		input := NewSliceIterator(rows)
		pattern := `\[(?P<timestamp>[^\]]+)\] (?P<level>\w+) in module (?P<module>\w+): (?P<message>.+)`
		extractedFields := []string{"timestamp", "level", "module", "message"}

		op, err := NewParseOperator(input, "log", pattern, extractedFields, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		timestamp, _ := row.Get("timestamp")
		assert.Equal(t, "2024-01-15 10:30:45", timestamp)

		level, _ := row.Get("level")
		assert.Equal(t, "ERROR", level)

		module, _ := row.Get("module")
		assert.Equal(t, "auth", module)

		message, _ := row.Get("message")
		assert.Equal(t, "Authentication failed for user admin", message)

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

		_, err := NewParseOperator(input, "message", invalidPattern, []string{}, logger)
		assert.Error(t, err)
	})
}
