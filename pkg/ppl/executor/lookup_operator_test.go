// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/lookup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLookupOperator(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("BasicLookup", func(t *testing.T) {
		// Create lookup table
		registry := lookup.NewRegistry(logger)
		table := lookup.NewLookupTable("products", logger)

		// Add test data to lookup table
		table.AddRow("101", map[string]interface{}{
			"product_id": "101",
			"name":       "Laptop",
			"price":      "999.99",
		})
		table.AddRow("102", map[string]interface{}{
			"product_id": "102",
			"name":       "Mouse",
			"price":      "29.99",
		})

		err := registry.Register(table)
		require.NoError(t, err)

		// Create test input data
		rows := []*Row{
			NewRow(map[string]interface{}{
				"order_id":   "1",
				"product_id": "101",
				"quantity":   "2",
			}),
			NewRow(map[string]interface{}{
				"order_id":   "2",
				"product_id": "102",
				"quantity":   "5",
			}),
		}

		input := NewSliceIterator(rows)
		outputFields := []string{"name", "price"}
		outputAliases := []string{"", ""} // No aliases

		op, err := NewLookupOperator(input, registry, "products", "product_id", "", outputFields, outputAliases, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		// First row - should have product name and price added
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		name, exists := row.Get("name")
		assert.True(t, exists)
		assert.Equal(t, "Laptop", name)

		price, exists := row.Get("price")
		assert.True(t, exists)
		assert.Equal(t, "999.99", price)

		// Original fields should still be present
		orderID, exists := row.Get("order_id")
		assert.True(t, exists)
		assert.Equal(t, "1", orderID)

		// Second row
		row, err = op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		name, _ = row.Get("name")
		assert.Equal(t, "Mouse", name)

		price, _ = row.Get("price")
		assert.Equal(t, "29.99", price)

		// EOF
		row, err = op.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)
		assert.Nil(t, row)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("LookupWithAliases", func(t *testing.T) {
		// Create lookup table
		registry := lookup.NewRegistry(logger)
		table := lookup.NewLookupTable("users", logger)

		table.AddRow("u1", map[string]interface{}{
			"user_id":  "u1",
			"username": "john_doe",
			"email":    "john@example.com",
		})

		err := registry.Register(table)
		require.NoError(t, err)

		// Create test input data
		rows := []*Row{
			NewRow(map[string]interface{}{
				"event_id": "e1",
				"user_id":  "u1",
			}),
		}

		input := NewSliceIterator(rows)
		outputFields := []string{"username", "email"}
		outputAliases := []string{"user", "contact_email"}

		op, err := NewLookupOperator(input, registry, "users", "user_id", "", outputFields, outputAliases, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		// Check aliased fields
		user, exists := row.Get("user")
		assert.True(t, exists)
		assert.Equal(t, "john_doe", user)

		email, exists := row.Get("contact_email")
		assert.True(t, exists)
		assert.Equal(t, "john@example.com", email)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("NoLookupMatch", func(t *testing.T) {
		// Create lookup table
		registry := lookup.NewRegistry(logger)
		table := lookup.NewLookupTable("products", logger)

		table.AddRow("101", map[string]interface{}{
			"product_id": "101",
			"name":       "Laptop",
		})

		err := registry.Register(table)
		require.NoError(t, err)

		// Create test input with non-existent product_id
		rows := []*Row{
			NewRow(map[string]interface{}{
				"order_id":   "1",
				"product_id": "999", // Doesn't exist in lookup table
			}),
		}

		input := NewSliceIterator(rows)
		outputFields := []string{"name"}
		outputAliases := []string{""}

		op, err := NewLookupOperator(input, registry, "products", "product_id", "", outputFields, outputAliases, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		// Should return row without adding lookup fields
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		_, exists := row.Get("name")
		assert.False(t, exists)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("MissingJoinField", func(t *testing.T) {
		// Create lookup table
		registry := lookup.NewRegistry(logger)
		table := lookup.NewLookupTable("products", logger)

		table.AddRow("101", map[string]interface{}{
			"product_id": "101",
			"name":       "Laptop",
		})

		err := registry.Register(table)
		require.NoError(t, err)

		// Create test input without product_id field
		rows := []*Row{
			NewRow(map[string]interface{}{
				"order_id": "1",
				"quantity": "2",
			}),
		}

		input := NewSliceIterator(rows)
		outputFields := []string{"name"}
		outputAliases := []string{""}

		op, err := NewLookupOperator(input, registry, "products", "product_id", "", outputFields, outputAliases, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		// Should return row without adding lookup fields
		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		_, exists := row.Get("name")
		assert.False(t, exists)

		err = op.Close()
		require.NoError(t, err)
	})

	t.Run("InvalidLookupTable", func(t *testing.T) {
		// Create empty registry
		registry := lookup.NewRegistry(logger)

		rows := []*Row{
			NewRow(map[string]interface{}{
				"product_id": "101",
			}),
		}

		input := NewSliceIterator(rows)
		outputFields := []string{"name"}
		outputAliases := []string{""}

		// Should fail because table doesn't exist
		_, err := NewLookupOperator(input, registry, "nonexistent", "product_id", "", outputFields, outputAliases, logger)
		assert.Error(t, err)
	})

	t.Run("MultipleFields", func(t *testing.T) {
		// Create lookup table with many fields
		registry := lookup.NewRegistry(logger)
		table := lookup.NewLookupTable("products", logger)

		table.AddRow("101", map[string]interface{}{
			"product_id":  "101",
			"name":        "Laptop",
			"price":       "999.99",
			"category":    "Electronics",
			"description": "High-performance laptop",
		})

		err := registry.Register(table)
		require.NoError(t, err)

		rows := []*Row{
			NewRow(map[string]interface{}{
				"order_id":   "1",
				"product_id": "101",
			}),
		}

		input := NewSliceIterator(rows)
		outputFields := []string{"name", "price", "category"}
		outputAliases := []string{"product_name", "product_price", ""}

		op, err := NewLookupOperator(input, registry, "products", "product_id", "", outputFields, outputAliases, logger)
		require.NoError(t, err)

		err = op.Open(ctx)
		require.NoError(t, err)

		row, err := op.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		// Check all fields
		productName, _ := row.Get("product_name")
		assert.Equal(t, "Laptop", productName)

		productPrice, _ := row.Get("product_price")
		assert.Equal(t, "999.99", productPrice)

		category, _ := row.Get("category")
		assert.Equal(t, "Electronics", category)

		err = op.Close()
		require.NoError(t, err)
	})
}
