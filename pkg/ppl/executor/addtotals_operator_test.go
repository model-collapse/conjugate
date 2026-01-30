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

func TestAddtotalsOperator_Basic(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"category": "A", "revenue": 100, "count": 5}),
			NewRow(map[string]interface{}{"category": "B", "revenue": 200, "count": 10}),
			NewRow(map[string]interface{}{"category": "C", "revenue": 300, "count": 15}),
		},
	}

	// Create addtotals operator with col=true (add column totals summary row)
	op := NewAddtotalsOperator(input, nil, false, true, "category", "Total", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 4 rows (3 original + 1 total)
	assert.Equal(t, 4, len(rows))

	// Check original rows
	cat1, _ := rows[0].Get("category")
	assert.Equal(t, "A", cat1)
	rev1, _ := rows[0].Get("revenue")
	assert.Equal(t, 100, rev1)

	cat2, _ := rows[1].Get("category")
	assert.Equal(t, "B", cat2)
	rev2, _ := rows[1].Get("revenue")
	assert.Equal(t, 200, rev2)

	cat3, _ := rows[2].Get("category")
	assert.Equal(t, "C", cat3)
	rev3, _ := rows[2].Get("revenue")
	assert.Equal(t, 300, rev3)

	// Check totals row
	totalsRow := rows[3]
	totalCat, _ := totalsRow.Get("category")
	assert.Equal(t, "Total", totalCat) // Label in first non-numeric field
	totalRev, _ := totalsRow.Get("revenue")
	assert.Equal(t, float64(600), totalRev)
	totalCount, _ := totalsRow.Get("count")
	assert.Equal(t, float64(30), totalCount)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_CustomLabel(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"region": "North", "sales": 1000}),
			NewRow(map[string]interface{}{"region": "South", "sales": 2000}),
		},
	}

	// Create addtotals operator with custom label (col=true to add summary row)
	op := NewAddtotalsOperator(input, nil, false, true, "region", "Grand Total", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 3 rows (2 original + 1 total)
	assert.Equal(t, 3, len(rows))

	// Check totals row
	totalsRow := rows[2]
	region, _ := totalsRow.Get("region")
	assert.Equal(t, "Grand Total", region)
	sales, _ := totalsRow.Get("sales")
	assert.Equal(t, float64(3000), sales)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_EmptyInput(t *testing.T) {
	logger := zap.NewNop()

	// Create empty input
	input := &MockOperator{
		rows: []*Row{},
	}

	// Create addtotals operator with col=true
	op := NewAddtotalsOperator(input, nil, false, true, "", "Total", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 0 rows (no totals for empty input)
	assert.Equal(t, 0, len(rows))

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_SingleRow(t *testing.T) {
	logger := zap.NewNop()

	// Create single row input
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"product": "Widget", "price": 99.99, "quantity": 5}),
		},
	}

	// Create addtotals operator with col=true
	op := NewAddtotalsOperator(input, nil, false, true, "product", "Total", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 2 rows (1 original + 1 total)
	assert.Equal(t, 2, len(rows))

	// Check totals row
	totalsRow := rows[1]
	product, _ := totalsRow.Get("product")
	assert.Equal(t, "Total", product)
	price, _ := totalsRow.Get("price")
	assert.Equal(t, 99.99, price)
	quantity, _ := totalsRow.Get("quantity")
	assert.Equal(t, float64(5), quantity)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_MixedTypes(t *testing.T) {
	logger := zap.NewNop()

	// Create input with mixed types
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"name": "Alice", "age": 30, "score": 95.5, "city": "NYC"}),
			NewRow(map[string]interface{}{"name": "Bob", "age": 25, "score": 88.0, "city": "LA"}),
			NewRow(map[string]interface{}{"name": "Carol", "age": 35, "score": 92.3, "city": "SF"}),
		},
	}

	// Create addtotals operator with col=true
	op := NewAddtotalsOperator(input, nil, false, true, "name", "Total", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 4 rows
	assert.Equal(t, 4, len(rows))

	// Check totals row - should only total numeric fields
	totalsRow := rows[3]
	name, _ := totalsRow.Get("name")
	assert.Equal(t, "Total", name)
	age, _ := totalsRow.Get("age")
	assert.Equal(t, float64(90), age) // 30 + 25 + 35
	score, _ := totalsRow.Get("score")
	assert.Equal(t, 275.8, score)      // 95.5 + 88.0 + 92.3
	_, hasCity := totalsRow.Get("city")
	assert.False(t, hasCity) // String field should not be totaled

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_IntegerTypes(t *testing.T) {
	logger := zap.NewNop()

	// Create input with different integer types
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"id": int32(1), "count": int64(100), "value": uint16(50)}),
			NewRow(map[string]interface{}{"id": int32(2), "count": int64(200), "value": uint16(75)}),
		},
	}

	// Create addtotals operator with col=true
	op := NewAddtotalsOperator(input, nil, false, true, "", "Total", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals row
	totalsRow := rows[2]
	id, _ := totalsRow.Get("id")
	assert.Equal(t, float64(3), id) // 1 + 2
	count, _ := totalsRow.Get("count")
	assert.Equal(t, float64(300), count) // 100 + 200
	value, _ := totalsRow.Get("value")
	assert.Equal(t, float64(125), value) // 50 + 75

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_NegativeNumbers(t *testing.T) {
	logger := zap.NewNop()

	// Create input with negative numbers
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"account": "A", "balance": 100}),
			NewRow(map[string]interface{}{"account": "B", "balance": -50}),
			NewRow(map[string]interface{}{"account": "C", "balance": 75}),
		},
	}

	// Create addtotals operator with col=true
	op := NewAddtotalsOperator(input, nil, false, true, "account", "Total", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals row
	totalsRow := rows[3]
	account, _ := totalsRow.Get("account")
	assert.Equal(t, "Total", account)
	balance, _ := totalsRow.Get("balance")
	assert.Equal(t, float64(125), balance) // 100 + (-50) + 75

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_WithFieldName(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"value1": 10, "value2": 20}),
			NewRow(map[string]interface{}{"value1": 30, "value2": 40}),
		},
	}

	// Create addtotals operator with fieldName (col=true to test labelfield in summary row)
	op := NewAddtotalsOperator(input, nil, false, true, "row_type", "TOTAL", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals row
	totalsRow := rows[2]
	rowType, _ := totalsRow.Get("row_type")
	assert.Equal(t, "TOTAL", rowType)
	value1, _ := totalsRow.Get("value1")
	assert.Equal(t, float64(40), value1)
	value2, _ := totalsRow.Get("value2")
	assert.Equal(t, float64(60), value2)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_ZeroValues(t *testing.T) {
	logger := zap.NewNop()

	// Create input with zero values
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"metric": "m1", "value": 0}),
			NewRow(map[string]interface{}{"metric": "m2", "value": 0}),
			NewRow(map[string]interface{}{"metric": "m3", "value": 0}),
		},
	}

	// Create addtotals operator with col=true
	op := NewAddtotalsOperator(input, nil, false, true, "metric", "Total", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check totals row
	totalsRow := rows[3]
	metric, _ := totalsRow.Get("metric")
	assert.Equal(t, "Total", metric)
	value, _ := totalsRow.Get("value")
	assert.Equal(t, float64(0), value) // 0 + 0 + 0

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_RowTotalsDefault(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"category": "A", "revenue": 100, "cost": 40}),
			NewRow(map[string]interface{}{"category": "B", "revenue": 200, "cost": 80}),
			NewRow(map[string]interface{}{"category": "C", "revenue": 300, "cost": 120}),
		},
	}

	// Create addtotals operator with default behavior: row=true, col=false
	op := NewAddtotalsOperator(input, nil, true, false, "", "Total", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 3 rows (no summary row added)
	assert.Equal(t, 3, len(rows))

	// Check row 1 - should have "total" field with row sum
	cat1, _ := rows[0].Get("category")
	assert.Equal(t, "A", cat1)
	total1, _ := rows[0].Get("total")
	assert.Equal(t, float64(140), total1) // 100 + 40

	// Check row 2
	total2, _ := rows[1].Get("total")
	assert.Equal(t, float64(280), total2) // 200 + 80

	// Check row 3
	total3, _ := rows[2].Get("total")
	assert.Equal(t, float64(420), total3) // 300 + 120

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_BothRowAndCol(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"product": "Widget", "q1": 100, "q2": 150}),
			NewRow(map[string]interface{}{"product": "Gadget", "q1": 200, "q2": 250}),
		},
	}

	// Create addtotals operator with both row=true and col=true
	op := NewAddtotalsOperator(input, nil, true, true, "product", "Total", "row_total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 3 rows (2 data rows with row totals + 1 column totals summary row)
	assert.Equal(t, 3, len(rows))

	// Check row 1 - has row total
	rowTotal1, _ := rows[0].Get("row_total")
	assert.Equal(t, float64(250), rowTotal1) // 100 + 150

	// Check row 2 - has row total
	rowTotal2, _ := rows[1].Get("row_total")
	assert.Equal(t, float64(450), rowTotal2) // 200 + 250

	// Check summary row - has column totals
	summaryRow := rows[2]
	product, _ := summaryRow.Get("product")
	assert.Equal(t, "Total", product)
	q1Total, _ := summaryRow.Get("q1")
	assert.Equal(t, float64(300), q1Total) // 100 + 200
	q2Total, _ := summaryRow.Get("q2")
	assert.Equal(t, float64(400), q2Total) // 150 + 250
	// Summary row should also have the sum of row totals
	summaryRowTotal, _ := summaryRow.Get("row_total")
	assert.Equal(t, float64(700), summaryRowTotal) // 250 + 450

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_NoModifications(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"id": 1, "value": 100}),
			NewRow(map[string]interface{}{"id": 2, "value": 200}),
		},
	}

	// Create addtotals operator with row=false, col=false (pass-through mode)
	op := NewAddtotalsOperator(input, nil, false, false, "", "", "total", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Should have 2 rows (no modifications)
	assert.Equal(t, 2, len(rows))

	// Check that rows are unchanged (no "total" field added)
	_, hasTotal1 := rows[0].Get("total")
	assert.False(t, hasTotal1)
	_, hasTotal2 := rows[1].Get("total")
	assert.False(t, hasTotal2)

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

func TestAddtotalsOperator_CustomFieldNameForRowTotals(t *testing.T) {
	logger := zap.NewNop()

	// Create input rows
	input := &MockOperator{
		rows: []*Row{
			NewRow(map[string]interface{}{"a": 10, "b": 20, "c": 30}),
			NewRow(map[string]interface{}{"a": 5, "b": 15, "c": 25}),
		},
	}

	// Create addtotals operator with custom field name for row totals
	op := NewAddtotalsOperator(input, nil, true, false, "", "Total", "sum_all", logger)

	// Open operator
	err := op.Open(context.Background())
	require.NoError(t, err)

	// Collect all rows
	rows, err := collectRows(context.Background(), op)
	require.NoError(t, err)

	// Check that custom field name is used
	sumAll1, _ := rows[0].Get("sum_all")
	assert.Equal(t, float64(60), sumAll1) // 10 + 20 + 30

	sumAll2, _ := rows[1].Get("sum_all")
	assert.Equal(t, float64(45), sumAll2) // 5 + 15 + 25

	// Close operator
	err = op.Close()
	require.NoError(t, err)
}

// Helper function to collect all rows from an operator
func collectRows(ctx context.Context, op Operator) ([]*Row, error) {
	var rows []*Row
	for {
		row, err := op.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}
