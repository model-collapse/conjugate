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

func TestSpathOperator_BasicExtraction(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{"user": {"name": "Alice", "age": 30, "email": "alice@example.com"}}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonData,
			"id":   1,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "user.name",
		OutputField: "user_name",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Should have extracted user.name
	userName, exists := row.Get("user_name")
	assert.True(t, exists)
	assert.Equal(t, "Alice", userName)

	// Original fields should remain
	id, _ := row.Get("id")
	assert.Equal(t, 1, id)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_NestedPath(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{"response": {"data": {"user": {"id": 12345, "active": true}}}}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonData,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "response.data.user.id",
		OutputField: "user_id",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	userId, exists := row.Get("user_id")
	assert.True(t, exists)
	// OpenSearch spath returns all values as strings
	assert.Equal(t, "12345", userId)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_ArrayAccess(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{"items": [{"name": "Item1"}, {"name": "Item2"}, {"name": "Item3"}]}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonData,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "items.0.name",
		OutputField: "first_item",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	firstItem, exists := row.Get("first_item")
	assert.True(t, exists)
	assert.Equal(t, "Item1", firstItem)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_ArrayWildcard(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{"items": [{"name": "Item1"}, {"name": "Item2"}, {"name": "Item3"}]}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonData,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "items.#.name",
		OutputField: "all_names",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	allNames, exists := row.Get("all_names")
	assert.True(t, exists)

	// Should be an array of names
	namesArray, ok := allNames.([]interface{})
	assert.True(t, ok)
	assert.Len(t, namesArray, 3)
	assert.Equal(t, "Item1", namesArray[0])
	assert.Equal(t, "Item2", namesArray[1])
	assert.Equal(t, "Item3", namesArray[2])

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_AutoExtract(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{"user": "Alice", "age": 30, "city": "NYC", "active": true}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonData,
			"id":   1,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path: "", // Empty path = auto-extract
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// All JSON fields should be extracted
	user, exists := row.Get("user")
	assert.True(t, exists)
	assert.Equal(t, "Alice", user)

	age, exists := row.Get("age")
	assert.True(t, exists)
	// OpenSearch spath returns all values as strings
	assert.Equal(t, "30", age)

	city, exists := row.Get("city")
	assert.True(t, exists)
	assert.Equal(t, "NYC", city)

	active, exists := row.Get("active")
	assert.True(t, exists)
	// OpenSearch spath returns all values as strings
	assert.Equal(t, "true", active)

	// Original id should remain
	id, exists := row.Get("id")
	assert.True(t, exists)
	assert.Equal(t, 1, id)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_MissingInputField(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"id": 1,
			// No _raw field
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "user.name",
		OutputField: "user_name",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Should return row as-is (no extraction)
	_, exists := row.Get("user_name")
	assert.False(t, exists)

	id, exists := row.Get("id")
	assert.True(t, exists)
	assert.Equal(t, 1, id)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_InvalidJSON(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": "not valid json {{{",
			"id":   1,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "user.name",
		OutputField: "user_name",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Should return row as-is (invalid JSON)
	_, exists := row.Get("user_name")
	assert.False(t, exists)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_NonExistentPath(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{"user": {"name": "Alice"}}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonData,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "user.email", // Doesn't exist
		OutputField: "user_email",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Path doesn't exist, field shouldn't be set
	_, exists := row.Get("user_email")
	assert.False(t, exists)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_CustomInputField(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{"status": "success", "code": 200}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"response_data": jsonData,
			"id":            1,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		InputField:  "response_data", // Custom input field
		Path:        "status",
		OutputField: "response_status",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	status, exists := row.Get("response_status")
	assert.True(t, exists)
	assert.Equal(t, "success", status)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_TypePreservation(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{
		"string_val": "hello",
		"int_val": 42,
		"float_val": 3.14,
		"bool_val": true,
		"null_val": null,
		"array_val": [1, 2, 3],
		"object_val": {"nested": "value"}
	}`

	// OpenSearch spath returns all scalar values as strings
	tests := []struct {
		path     string
		output   string
		expected interface{}
	}{
		{"string_val", "str", "hello"},
		{"int_val", "int", "42"},         // String per OpenSearch spec
		{"float_val", "float", "3.14"},   // String per OpenSearch spec
		{"bool_val", "bool", "true"},     // String per OpenSearch spec
		{"null_val", "null", nil},        // null remains nil
	}

	for _, tt := range tests {
		input := NewSliceIterator([]*Row{NewRow(map[string]interface{}{"_raw": jsonData})})
		config := SpathConfig{
			Path:        tt.path,
			OutputField: tt.output,
		}

		op := NewSpathOperator(input, config, logger)
		err := op.Open(ctx)
		require.NoError(t, err)

		row, err := op.Next(ctx)
		require.NoError(t, err)

		val, exists := row.Get(tt.output)
		assert.True(t, exists, "Field %s should exist", tt.output)
		assert.Equal(t, tt.expected, val, "Field %s value mismatch", tt.output)

		op.Close()
	}
}

func TestSpathOperator_ComplexObject(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{
		"user": {
			"id": 12345,
			"profile": {
				"name": "Alice",
				"location": "NYC"
			}
		}
	}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonData,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "user.profile",
		OutputField: "profile",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	profile, exists := row.Get("profile")
	assert.True(t, exists)

	// Should be a map
	profileMap, ok := profile.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Alice", profileMap["name"])
	assert.Equal(t, "NYC", profileMap["location"])

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_DollarPrefix(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{"user": {"name": "Bob"}}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonData,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "$.user.name", // JSONPath with $ prefix
		OutputField: "user_name",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	userName, exists := row.Get("user_name")
	assert.True(t, exists)
	assert.Equal(t, "Bob", userName)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_DerivedFieldName(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	jsonData := `{"response": {"data": {"user_id": 999}}}`

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonData,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path: "response.data.user_id",
		// No OutputField specified - should derive "user_id"
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	// Should auto-derive field name from path
	userId, exists := row.Get("user_id")
	assert.True(t, exists)
	// OpenSearch spath returns all values as strings
	assert.Equal(t, "999", userId)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_MultipleRows(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": `{"user": "Alice", "score": 100}`,
		}),
		NewRow(map[string]interface{}{
			"_raw": `{"user": "Bob", "score": 200}`,
		}),
		NewRow(map[string]interface{}{
			"_raw": `{"user": "Charlie", "score": 300}`,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "user",
		OutputField: "username",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	// First row
	row1, err := op.Next(ctx)
	require.NoError(t, err)
	user1, _ := row1.Get("username")
	assert.Equal(t, "Alice", user1)

	// Second row
	row2, err := op.Next(ctx)
	require.NoError(t, err)
	user2, _ := row2.Get("username")
	assert.Equal(t, "Bob", user2)

	// Third row
	row3, err := op.Next(ctx)
	require.NoError(t, err)
	user3, _ := row3.Get("username")
	assert.Equal(t, "Charlie", user3)

	// EOF
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = op.Close()
	require.NoError(t, err)
}

func TestSpathOperator_MapInput(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Test with map[string]interface{} as input (not string)
	jsonMap := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}

	rows := []*Row{
		NewRow(map[string]interface{}{
			"_raw": jsonMap,
		}),
	}

	input := NewSliceIterator(rows)
	config := SpathConfig{
		Path:        "user.name",
		OutputField: "user_name",
	}

	op := NewSpathOperator(input, config, logger)
	err := op.Open(ctx)
	require.NoError(t, err)

	row, err := op.Next(ctx)
	require.NoError(t, err)

	userName, exists := row.Get("user_name")
	assert.True(t, exists)
	assert.Equal(t, "Alice", userName)

	err = op.Close()
	require.NoError(t, err)
}
