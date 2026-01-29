// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestEventstatsOperator_SimpleCount(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	// Create test data
	rows := []*Row{
		NewRow(map[string]interface{}{"host": "server1", "status": 200}),
		NewRow(map[string]interface{}{"host": "server2", "status": 404}),
		NewRow(map[string]interface{}{"host": "server3", "status": 200}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func:  &ast.FunctionCall{Name: "count"},
			Alias: "total",
		},
	}
	eventstatsOp := NewEventstatsOperator(iter, nil, aggregations, logger)

	err := eventstatsOp.Open(ctx)
	require.NoError(t, err)

	// All rows should have the same total count added
	for i := 0; i < 3; i++ {
		row, err := eventstatsOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		// Original fields should still exist
		_, exists := row.Get("host")
		assert.True(t, exists)

		// New aggregation field should be added
		total, exists := row.Get("total")
		assert.True(t, exists)
		assert.Equal(t, int64(3), total, "All rows should have total=3")
	}

	err = eventstatsOp.Close()
	assert.NoError(t, err)
}

func TestEventstatsOperator_AvgWithGroupBy(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	// Create test data with different hosts
	rows := []*Row{
		NewRow(map[string]interface{}{"host": "server1", "latency": 50}),
		NewRow(map[string]interface{}{"host": "server1", "latency": 100}),
		NewRow(map[string]interface{}{"host": "server2", "latency": 200}),
		NewRow(map[string]interface{}{"host": "server2", "latency": 300}),
	}

	iter := NewSliceIterator(rows)
	groupBy := []ast.Expression{
		&ast.FieldReference{Name: "host"},
	}
	aggregations := []*ast.Aggregation{
		{
			Func: &ast.FunctionCall{
				Name:      "avg",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
			},
			Alias: "avg_latency",
		},
	}
	eventstatsOp := NewEventstatsOperator(iter, groupBy, aggregations, logger)

	err := eventstatsOp.Open(ctx)
	require.NoError(t, err)

	// First two rows: server1 with avg_latency = 75
	for i := 0; i < 2; i++ {
		row, err := eventstatsOp.Next(ctx)
		require.NoError(t, err)

		host, _ := row.Get("host")
		assert.Equal(t, "server1", host)

		avgLatency, exists := row.Get("avg_latency")
		assert.True(t, exists)
		assert.Equal(t, float64(75), avgLatency, "server1 avg should be 75")
	}

	// Last two rows: server2 with avg_latency = 250
	for i := 0; i < 2; i++ {
		row, err := eventstatsOp.Next(ctx)
		require.NoError(t, err)

		host, _ := row.Get("host")
		assert.Equal(t, "server2", host)

		avgLatency, exists := row.Get("avg_latency")
		assert.True(t, exists)
		assert.Equal(t, float64(250), avgLatency, "server2 avg should be 250")
	}

	err = eventstatsOp.Close()
	assert.NoError(t, err)
}

func TestEventstatsOperator_MultipleAggregations(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"value": 10}),
		NewRow(map[string]interface{}{"value": 20}),
		NewRow(map[string]interface{}{"value": 30}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func: &ast.FunctionCall{
				Name:      "sum",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "total",
		},
		{
			Func: &ast.FunctionCall{
				Name:      "avg",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "average",
		},
		{
			Func: &ast.FunctionCall{
				Name:      "min",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "minimum",
		},
		{
			Func: &ast.FunctionCall{
				Name:      "max",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "maximum",
		},
	}
	eventstatsOp := NewEventstatsOperator(iter, nil, aggregations, logger)

	err := eventstatsOp.Open(ctx)
	require.NoError(t, err)

	// All rows should have all aggregation results
	for i := 0; i < 3; i++ {
		row, err := eventstatsOp.Next(ctx)
		require.NoError(t, err)

		total, _ := row.Get("total")
		assert.Equal(t, float64(60), total)

		average, _ := row.Get("average")
		assert.Equal(t, float64(20), average)

		minimum, _ := row.Get("minimum")
		assert.Equal(t, 10, minimum)

		maximum, _ := row.Get("maximum")
		assert.Equal(t, 30, maximum)
	}

	err = eventstatsOp.Close()
	assert.NoError(t, err)
}

func TestEventstatsOperator_CountByGroup(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"status": "error"}),
		NewRow(map[string]interface{}{"status": "ok"}),
		NewRow(map[string]interface{}{"status": "error"}),
		NewRow(map[string]interface{}{"status": "error"}),
		NewRow(map[string]interface{}{"status": "ok"}),
	}

	iter := NewSliceIterator(rows)
	groupBy := []ast.Expression{
		&ast.FieldReference{Name: "status"},
	}
	aggregations := []*ast.Aggregation{
		{
			Func:  &ast.FunctionCall{Name: "count"},
			Alias: "count",
		},
	}
	eventstatsOp := NewEventstatsOperator(iter, groupBy, aggregations, logger)

	err := eventstatsOp.Open(ctx)
	require.NoError(t, err)

	errorCount := 0
	okCount := 0

	for i := 0; i < 5; i++ {
		row, err := eventstatsOp.Next(ctx)
		require.NoError(t, err)

		status, _ := row.Get("status")
		count, _ := row.Get("count")

		if status == "error" {
			assert.Equal(t, int64(3), count, "error count should be 3")
			errorCount++
		} else if status == "ok" {
			assert.Equal(t, int64(2), count, "ok count should be 2")
			okCount++
		}
	}

	assert.Equal(t, 3, errorCount)
	assert.Equal(t, 2, okCount)

	err = eventstatsOp.Close()
	assert.NoError(t, err)
}

func TestEventstatsOperator_EmptyInput(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{}
	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func:  &ast.FunctionCall{Name: "count"},
			Alias: "total",
		},
	}
	eventstatsOp := NewEventstatsOperator(iter, nil, aggregations, logger)

	err := eventstatsOp.Open(ctx)
	require.NoError(t, err)

	// Should return no rows
	_, err = eventstatsOp.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	err = eventstatsOp.Close()
	assert.NoError(t, err)
}

func TestEventstatsOperator_PreservesOriginalFields(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"field1": "value1",
			"field2": 100,
			"field3": true,
		}),
		NewRow(map[string]interface{}{
			"field1": "value2",
			"field2": 200,
			"field3": false,
		}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func:  &ast.FunctionCall{Name: "count"},
			Alias: "total",
		},
	}
	eventstatsOp := NewEventstatsOperator(iter, nil, aggregations, logger)

	err := eventstatsOp.Open(ctx)
	require.NoError(t, err)

	// First row
	row, err := eventstatsOp.Next(ctx)
	require.NoError(t, err)

	field1, _ := row.Get("field1")
	assert.Equal(t, "value1", field1)

	field2, _ := row.Get("field2")
	assert.Equal(t, 100, field2)

	field3, _ := row.Get("field3")
	assert.Equal(t, true, field3)

	total, _ := row.Get("total")
	assert.Equal(t, int64(2), total)

	// Second row
	row, err = eventstatsOp.Next(ctx)
	require.NoError(t, err)

	field1, _ = row.Get("field1")
	assert.Equal(t, "value2", field1)

	field2, _ = row.Get("field2")
	assert.Equal(t, 200, field2)

	field3, _ = row.Get("field3")
	assert.Equal(t, false, field3)

	total, _ = row.Get("total")
	assert.Equal(t, int64(2), total)

	err = eventstatsOp.Close()
	assert.NoError(t, err)
}

func TestEventstatsOperator_Stats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"value": 1}),
		NewRow(map[string]interface{}{"value": 2}),
		NewRow(map[string]interface{}{"value": 3}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func:  &ast.FunctionCall{Name: "count"},
			Alias: "total",
		},
	}
	eventstatsOp := NewEventstatsOperator(iter, nil, aggregations, logger)

	err := eventstatsOp.Open(ctx)
	require.NoError(t, err)

	// Read all rows
	for i := 0; i < 3; i++ {
		_, err := eventstatsOp.Next(ctx)
		require.NoError(t, err)
	}

	stats := eventstatsOp.Stats()
	assert.Equal(t, int64(3), stats.RowsRead)
	assert.Equal(t, int64(3), stats.RowsReturned)

	err = eventstatsOp.Close()
	assert.NoError(t, err)
}

func TestEventstatsOperator_MultiFieldGroupBy(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"host": "server1", "status": "ok", "value": 10}),
		NewRow(map[string]interface{}{"host": "server1", "status": "error", "value": 20}),
		NewRow(map[string]interface{}{"host": "server1", "status": "ok", "value": 30}),
		NewRow(map[string]interface{}{"host": "server2", "status": "ok", "value": 40}),
	}

	iter := NewSliceIterator(rows)
	groupBy := []ast.Expression{
		&ast.FieldReference{Name: "host"},
		&ast.FieldReference{Name: "status"},
	}
	aggregations := []*ast.Aggregation{
		{
			Func: &ast.FunctionCall{
				Name:      "sum",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "total",
		},
	}
	eventstatsOp := NewEventstatsOperator(iter, groupBy, aggregations, logger)

	err := eventstatsOp.Open(ctx)
	require.NoError(t, err)

	results := make(map[string]float64)
	for i := 0; i < 4; i++ {
		row, err := eventstatsOp.Next(ctx)
		require.NoError(t, err)

		host, _ := row.Get("host")
		status, _ := row.Get("status")
		total, _ := row.Get("total")

		key := fmt.Sprintf("%s-%s", host, status)
		results[key] = total.(float64)
	}

	// Verify group totals
	assert.Equal(t, float64(40), results["server1-ok"], "server1-ok should sum to 40")
	assert.Equal(t, float64(20), results["server1-error"], "server1-error should sum to 20")
	assert.Equal(t, float64(40), results["server2-ok"], "server2-ok should sum to 40")

	err = eventstatsOp.Close()
	assert.NoError(t, err)
}
