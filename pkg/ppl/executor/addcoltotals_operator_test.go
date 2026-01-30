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

func TestAddcoltotalsOperator_Basic(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"category": "A", "revenue": 100, "count": 5}),
			NewRow(map[string]interface{}{"category": "B", "revenue": 200, "count": 10}),
			NewRow(map[string]interface{}{"category": "C", "revenue": 300, "count": 15}),
		},
	}

	// Create addcoltotals operator (adds "Total" column with row sums)
	op := NewAddcoltotalsOperator(input, nil, "", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 3 rows (same count, but each has Total column)
	assert.Equal(t, 3, len(rows))

	// Check first row: category="A", revenue=100, count=5, Total=105
	cat1, _ := rows[0].Get("category")
	assert.Equal(t, "A", cat1)
	rev1, _ := rows[0].Get("revenue")
	assert.Equal(t, 100, rev1)
	count1, _ := rows[0].Get("count")
	assert.Equal(t, 5, count1)
	total1, _ := rows[0].Get("Total")
	assert.Equal(t, float64(105), total1) // 100 + 5

	// Check second row: category="B", revenue=200, count=10, Total=210
	cat2, _ := rows[1].Get("category")
	assert.Equal(t, "B", cat2)
	total2, _ := rows[1].Get("Total")
	assert.Equal(t, float64(210), total2) // 200 + 10

	// Check third row: category="C", revenue=300, count=15, Total=315
	cat3, _ := rows[2].Get("category")
	assert.Equal(t, "C", cat3)
	total3, _ := rows[2].Get("Total")
	assert.Equal(t, float64(315), total3) // 300 + 15

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddcoltotalsOperator_CustomColumnName(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"region": "North", "sales": 1000, "expenses": 200}),
			NewRow(map[string]interface{}{"region": "South", "sales": 2000, "expenses": 500}),
		},
	}

	// Create addcoltotals operator with custom column name
	op := NewAddcoltotalsOperator(input, nil, "RowSum", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 2 rows with RowSum column
	assert.Equal(t, 2, len(rows))

	// Check first row
	rowSum1, exists := rows[0].Get("RowSum")
	assert.True(t, exists)
	assert.Equal(t, float64(1200), rowSum1) // 1000 + 200

	// Check second row
	rowSum2, _ := rows[1].Get("RowSum")
	assert.Equal(t, float64(2500), rowSum2) // 2000 + 500

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddcoltotalsOperator_EmptyInput(t *testing.T) {
	logger := zap.NewNop()

	// Create empty input
	input := &MockOperator{
		rows: []*Row{},
	}

	// Create addcoltotals operator
	op := NewAddcoltotalsOperator(input, nil, "", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 0 rows
	assert.Equal(t, 0, len(rows))

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddcoltotalsOperator_SingleNumericField(t *testing.T) {
	logger := zap.NewNop()

	// Create input with single numeric field
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"product": "Widget", "price": 99.99}),
			NewRow(map[string]interface{}{"product": "Gadget", "price": 149.99}),
		},
	}

	// Create addcoltotals operator
	op := NewAddcoltotalsOperator(input, nil, "", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals (should just be the price since it's the only numeric field)
	total1, _ := rows[0].Get("Total")
	assert.Equal(t, 99.99, total1)

	total2, _ := rows[1].Get("Total")
	assert.Equal(t, 149.99, total2)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddcoltotalsOperator_MixedTypes(t *testing.T) {
	logger := zap.NewNop()

	// Create input with mixed types (numeric and string)
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"name": "Alice", "age": 30, "score": 95.5, "city": "NYC"}),
			NewRow(map[string]interface{}{"name": "Bob", "age": 25, "score": 88.0, "city": "LA"}),
		},
	}

	// Create addcoltotals operator
	op := NewAddcoltotalsOperator(input, nil, "", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals (should only sum numeric fields: age + score)
	total1, _ := rows[0].Get("Total")
	assert.Equal(t, 125.5, total1) // 30 + 95.5

	total2, _ := rows[1].Get("Total")
	assert.Equal(t, 113.0, total2) // 25 + 88.0

	// Verify string fields are preserved
	name1, _ := rows[0].Get("name")
	assert.Equal(t, "Alice", name1)
	city1, _ := rows[0].Get("city")
	assert.Equal(t, "NYC", city1)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddcoltotalsOperator_NegativeNumbers(t *testing.T) {
	logger := zap.NewNop()

	// Create input with negative numbers
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"account": "A", "credits": 100, "debits": -30}),
			NewRow(map[string]interface{}{"account": "B", "credits": 200, "debits": -75}),
		},
	}

	// Create addcoltotals operator
	op := NewAddcoltotalsOperator(input, nil, "", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals (should handle negative numbers)
	total1, _ := rows[0].Get("Total")
	assert.Equal(t, float64(70), total1) // 100 + (-30)

	total2, _ := rows[1].Get("Total")
	assert.Equal(t, float64(125), total2) // 200 + (-75)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddcoltotalsOperator_ZeroValues(t *testing.T) {
	logger := zap.NewNop()

	// Create input with zero values
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"metric": "m1", "value1": 0, "value2": 0}),
			NewRow(map[string]interface{}{"metric": "m2", "value1": 10, "value2": 0}),
		},
	}

	// Create addcoltotals operator
	op := NewAddcoltotalsOperator(input, nil, "", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals
	total1, _ := rows[0].Get("Total")
	assert.Equal(t, float64(0), total1) // 0 + 0

	total2, _ := rows[1].Get("Total")
	assert.Equal(t, float64(10), total2) // 10 + 0

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddcoltotalsOperator_IntegerTypes(t *testing.T) {
	logger := zap.NewNop()

	// Create input with different integer types
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"id": int32(1), "count": int64(100), "value": uint16(50)}),
			NewRow(map[string]interface{}{"id": int32(2), "count": int64(200), "value": uint16(75)}),
		},
	}

	// Create addcoltotals operator
	op := NewAddcoltotalsOperator(input, nil, "", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals (all types should be converted to float64)
	total1, _ := rows[0].Get("Total")
	assert.Equal(t, float64(151), total1) // 1 + 100 + 50

	total2, _ := rows[1].Get("Total")
	assert.Equal(t, float64(277), total2) // 2 + 200 + 75

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddcoltotalsOperator_OnlyStrings(t *testing.T) {
	logger := zap.NewNop()

	// Create input with only string fields
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"name": "Alice", "city": "NYC", "country": "USA"}),
			NewRow(map[string]interface{}{"name": "Bob", "city": "LA", "country": "USA"}),
		},
	}

	// Create addcoltotals operator
	op := NewAddcoltotalsOperator(input, nil, "", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals (should be 0 since no numeric fields)
	total1, _ := rows[0].Get("Total")
	assert.Equal(t, float64(0), total1)

	total2, _ := rows[1].Get("Total")
	assert.Equal(t, float64(0), total2)

	// Verify string fields are preserved
	name1, _ := rows[0].Get("name")
	assert.Equal(t, "Alice", name1)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddcoltotalsOperator_Streaming(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"value": 10}),
			NewRow(map[string]interface{}{"value": 20}),
			NewRow(map[string]interface{}{"value": 30}),
		},
	}

	// Create addcoltotals operator
	op := NewAddcoltotalsOperator(input, nil, "", "", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Read rows one at a time (streaming)
	ctx := context.Background()

	row1, err := op.Next(ctx)
	require.NoError(t, err)
	total1, _ := row1.Get("Total")
	assert.Equal(t, float64(10), total1)

	row2, err := op.Next(ctx)
	require.NoError(t, err)
	total2, _ := row2.Get("Total")
	assert.Equal(t, float64(20), total2)

	row3, err := op.Next(ctx)
	require.NoError(t, err)
	total3, _ := row3.Get("Total")
	assert.Equal(t, float64(30), total3)

	// No more rows
	_, err = op.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}
