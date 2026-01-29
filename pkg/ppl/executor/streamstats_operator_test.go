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

func TestStreamstatsOperator_RunningCount(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	// Create test data
	rows := []*Row{
		NewRow(map[string]interface{}{"value": 10}),
		NewRow(map[string]interface{}{"value": 20}),
		NewRow(map[string]interface{}{"value": 30}),
		NewRow(map[string]interface{}{"value": 40}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func:  &ast.FunctionCall{Name: "count"},
			Alias: "running_count",
		},
	}
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 0, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// Running count should increase with each row
	for i := 1; i <= 4; i++ {
		row, err := streamstatsOp.Next(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)

		runningCount, exists := row.Get("running_count")
		assert.True(t, exists)
		assert.Equal(t, int64(i), runningCount, "Row %d should have running_count=%d", i, i)
	}

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_RunningSum(t *testing.T) {
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
			Alias: "running_sum",
		},
	}
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 0, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// Running sum: 10, 30, 60
	expectedSums := []float64{10, 30, 60}
	for i, expectedSum := range expectedSums {
		row, err := streamstatsOp.Next(ctx)
		require.NoError(t, err)

		runningSum, exists := row.Get("running_sum")
		assert.True(t, exists)
		assert.Equal(t, expectedSum, runningSum, "Row %d should have running_sum=%.0f", i+1, expectedSum)
	}

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_RunningAverage(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"value": 10}),
		NewRow(map[string]interface{}{"value": 20}),
		NewRow(map[string]interface{}{"value": 30}),
		NewRow(map[string]interface{}{"value": 40}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func: &ast.FunctionCall{
				Name:      "avg",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "running_avg",
		},
	}
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 0, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// Running average: 10, 15, 20, 25
	expectedAvgs := []float64{10, 15, 20, 25}
	for i, expectedAvg := range expectedAvgs {
		row, err := streamstatsOp.Next(ctx)
		require.NoError(t, err)

		runningAvg, exists := row.Get("running_avg")
		assert.True(t, exists)
		assert.Equal(t, expectedAvg, runningAvg, "Row %d should have running_avg=%.0f", i+1, expectedAvg)
	}

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_WindowedCount(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"value": 1}),
		NewRow(map[string]interface{}{"value": 2}),
		NewRow(map[string]interface{}{"value": 3}),
		NewRow(map[string]interface{}{"value": 4}),
		NewRow(map[string]interface{}{"value": 5}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func:  &ast.FunctionCall{Name: "count"},
			Alias: "windowed_count",
		},
	}
	// Window size of 3
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 3, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// Windowed count with window=3: 1, 2, 3, 3, 3
	expectedCounts := []int64{1, 2, 3, 3, 3}
	for i, expectedCount := range expectedCounts {
		row, err := streamstatsOp.Next(ctx)
		require.NoError(t, err)

		windowedCount, exists := row.Get("windowed_count")
		assert.True(t, exists)
		assert.Equal(t, expectedCount, windowedCount, "Row %d should have windowed_count=%d", i+1, expectedCount)
	}

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_WindowedSum(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"value": 10}),
		NewRow(map[string]interface{}{"value": 20}),
		NewRow(map[string]interface{}{"value": 30}),
		NewRow(map[string]interface{}{"value": 40}),
		NewRow(map[string]interface{}{"value": 50}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func: &ast.FunctionCall{
				Name:      "sum",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "windowed_sum",
		},
	}
	// Window size of 3
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 3, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// Windowed sum with window=3: 10, 30, 60, 90, 120
	expectedSums := []float64{10, 30, 60, 90, 120}
	for i, expectedSum := range expectedSums {
		row, err := streamstatsOp.Next(ctx)
		require.NoError(t, err)

		windowedSum, exists := row.Get("windowed_sum")
		assert.True(t, exists)
		assert.Equal(t, expectedSum, windowedSum, "Row %d should have windowed_sum=%.0f", i+1, expectedSum)
	}

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_WindowedAverage(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"value": 10}),
		NewRow(map[string]interface{}{"value": 20}),
		NewRow(map[string]interface{}{"value": 30}),
		NewRow(map[string]interface{}{"value": 40}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func: &ast.FunctionCall{
				Name:      "avg",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "rolling_avg",
		},
	}
	// Window size of 2
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 2, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// Rolling average with window=2: 10, 15, 25, 35
	expectedAvgs := []float64{10, 15, 25, 35}
	for i, expectedAvg := range expectedAvgs {
		row, err := streamstatsOp.Next(ctx)
		require.NoError(t, err)

		rollingAvg, exists := row.Get("rolling_avg")
		assert.True(t, exists)
		assert.Equal(t, expectedAvg, rollingAvg, "Row %d should have rolling_avg=%.0f", i+1, expectedAvg)
	}

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_GroupBy(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"host": "server1", "value": 10}),
		NewRow(map[string]interface{}{"host": "server2", "value": 20}),
		NewRow(map[string]interface{}{"host": "server1", "value": 30}),
		NewRow(map[string]interface{}{"host": "server2", "value": 40}),
		NewRow(map[string]interface{}{"host": "server1", "value": 50}),
	}

	iter := NewSliceIterator(rows)
	groupBy := []ast.Expression{
		&ast.FieldReference{Name: "host"},
	}
	aggregations := []*ast.Aggregation{
		{
			Func: &ast.FunctionCall{
				Name:      "sum",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "running_sum",
		},
	}
	streamstatsOp := NewStreamstatsOperator(iter, groupBy, aggregations, 0, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// Separate running sums per host
	// server1: 10, 40, 90
	// server2: 20, 60
	expectedResults := []struct {
		host       string
		runningSum float64
	}{
		{"server1", 10},
		{"server2", 20},
		{"server1", 40},
		{"server2", 60},
		{"server1", 90},
	}

	for i, expected := range expectedResults {
		row, err := streamstatsOp.Next(ctx)
		require.NoError(t, err)

		host, _ := row.Get("host")
		assert.Equal(t, expected.host, host, "Row %d host mismatch", i+1)

		runningSum, exists := row.Get("running_sum")
		assert.True(t, exists)
		assert.Equal(t, expected.runningSum, runningSum, "Row %d running_sum mismatch", i+1)
	}

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_MultipleAggregations(t *testing.T) {
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
			Alias: "running_sum",
		},
		{
			Func: &ast.FunctionCall{
				Name:      "avg",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "running_avg",
		},
		{
			Func:  &ast.FunctionCall{Name: "count"},
			Alias: "running_count",
		},
	}
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 0, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// First row
	row, err := streamstatsOp.Next(ctx)
	require.NoError(t, err)
	assert.Equal(t, float64(10), row.ToMap()["running_sum"])
	assert.Equal(t, float64(10), row.ToMap()["running_avg"])
	assert.Equal(t, int64(1), row.ToMap()["running_count"])

	// Second row
	row, err = streamstatsOp.Next(ctx)
	require.NoError(t, err)
	assert.Equal(t, float64(30), row.ToMap()["running_sum"])
	assert.Equal(t, float64(15), row.ToMap()["running_avg"])
	assert.Equal(t, int64(2), row.ToMap()["running_count"])

	// Third row
	row, err = streamstatsOp.Next(ctx)
	require.NoError(t, err)
	assert.Equal(t, float64(60), row.ToMap()["running_sum"])
	assert.Equal(t, float64(20), row.ToMap()["running_avg"])
	assert.Equal(t, int64(3), row.ToMap()["running_count"])

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_RunningMinMax(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{"value": 50}),
		NewRow(map[string]interface{}{"value": 30}),
		NewRow(map[string]interface{}{"value": 70}),
		NewRow(map[string]interface{}{"value": 20}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func: &ast.FunctionCall{
				Name:      "min",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "running_min",
		},
		{
			Func: &ast.FunctionCall{
				Name:      "max",
				Arguments: []ast.Expression{&ast.FieldReference{Name: "value"}},
			},
			Alias: "running_max",
		},
	}
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 0, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// Running min: 50, 30, 30, 20
	// Running max: 50, 50, 70, 70
	expected := []struct {
		min int
		max int
	}{
		{50, 50},
		{30, 50},
		{30, 70},
		{20, 70},
	}

	for i, exp := range expected {
		row, err := streamstatsOp.Next(ctx)
		require.NoError(t, err)

		runningMin, _ := row.Get("running_min")
		assert.Equal(t, exp.min, runningMin, "Row %d running_min mismatch", i+1)

		runningMax, _ := row.Get("running_max")
		assert.Equal(t, exp.max, runningMax, "Row %d running_max mismatch", i+1)
	}

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_PreservesOriginalFields(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	rows := []*Row{
		NewRow(map[string]interface{}{
			"field1": "value1",
			"field2": 100,
		}),
		NewRow(map[string]interface{}{
			"field1": "value2",
			"field2": 200,
		}),
	}

	iter := NewSliceIterator(rows)
	aggregations := []*ast.Aggregation{
		{
			Func:  &ast.FunctionCall{Name: "count"},
			Alias: "running_count",
		},
	}
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 0, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// First row
	row, err := streamstatsOp.Next(ctx)
	require.NoError(t, err)

	field1, _ := row.Get("field1")
	assert.Equal(t, "value1", field1)

	field2, _ := row.Get("field2")
	assert.Equal(t, 100, field2)

	runningCount, _ := row.Get("running_count")
	assert.Equal(t, int64(1), runningCount)

	// Second row
	row, err = streamstatsOp.Next(ctx)
	require.NoError(t, err)

	field1, _ = row.Get("field1")
	assert.Equal(t, "value2", field1)

	field2, _ = row.Get("field2")
	assert.Equal(t, 200, field2)

	runningCount, _ = row.Get("running_count")
	assert.Equal(t, int64(2), runningCount)

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}

func TestStreamstatsOperator_Stats(t *testing.T) {
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
			Alias: "running_count",
		},
	}
	streamstatsOp := NewStreamstatsOperator(iter, nil, aggregations, 0, logger)

	err := streamstatsOp.Open(ctx)
	require.NoError(t, err)

	// Read all rows
	for i := 0; i < 3; i++ {
		_, err := streamstatsOp.Next(ctx)
		require.NoError(t, err)
	}

	stats := streamstatsOp.Stats()
	assert.Equal(t, int64(3), stats.RowsRead)
	assert.Equal(t, int64(3), stats.RowsReturned)

	err = streamstatsOp.Close()
	assert.NoError(t, err)
}
