// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/dsl"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// mockDataSource implements DataSource for testing
type mockDataSource struct {
	hits []*SearchHit
}

func (m *mockDataSource) Search(ctx context.Context, index string, queryDSL []byte, from, size int) (*SearchResult, error) {
	// Return mock hits
	end := from + size
	if end > len(m.hits) {
		end = len(m.hits)
	}
	if from > len(m.hits) {
		from = len(m.hits)
	}

	return &SearchResult{
		TookMillis: 10,
		TotalHits:  int64(len(m.hits)),
		MaxScore:   1.0,
		Hits:       m.hits[from:end],
	}, nil
}

func createMockDataSource() *mockDataSource {
	return &mockDataSource{
		hits: []*SearchHit{
			{ID: "1", Score: 1.0, Source: map[string]interface{}{"status": 200, "latency": 50.0, "host": "web1"}},
			{ID: "2", Score: 0.9, Source: map[string]interface{}{"status": 500, "latency": 100.0, "host": "web2"}},
			{ID: "3", Score: 0.8, Source: map[string]interface{}{"status": 200, "latency": 30.0, "host": "web1"}},
			{ID: "4", Score: 0.7, Source: map[string]interface{}{"status": 404, "latency": 20.0, "host": "web3"}},
			{ID: "5", Score: 0.6, Source: map[string]interface{}{"status": 200, "latency": 80.0, "host": "web2"}},
		},
	}
}

func createTestSchema() *analyzer.Schema {
	schema := analyzer.NewSchema("logs")
	schema.AddField("status", analyzer.FieldTypeInt)
	schema.AddField("latency", analyzer.FieldTypeDouble)
	schema.AddField("host", analyzer.FieldTypeString)
	return schema
}

func TestRow(t *testing.T) {
	t.Run("BasicOperations", func(t *testing.T) {
		row := NewRow(map[string]interface{}{
			"status":  200,
			"latency": 50.5,
			"host":    "web1",
			"active":  true,
		})

		// Get
		val, ok := row.Get("status")
		assert.True(t, ok)
		assert.Equal(t, 200, val)

		// GetString
		str, ok := row.GetString("host")
		assert.True(t, ok)
		assert.Equal(t, "web1", str)

		// GetFloat64
		f, ok := row.GetFloat64("latency")
		assert.True(t, ok)
		assert.Equal(t, 50.5, f)

		// GetBool
		b, ok := row.GetBool("active")
		assert.True(t, ok)
		assert.True(t, b)

		// Non-existent field
		_, ok = row.Get("nonexistent")
		assert.False(t, ok)
	})

	t.Run("SetAndDelete", func(t *testing.T) {
		row := NewRow(nil)

		row.Set("field1", "value1")
		val, ok := row.Get("field1")
		assert.True(t, ok)
		assert.Equal(t, "value1", val)

		row.Delete("field1")
		_, ok = row.Get("field1")
		assert.False(t, ok)
	})

	t.Run("Clone", func(t *testing.T) {
		original := NewRow(map[string]interface{}{"a": 1, "b": 2})
		clone := original.Clone()

		clone.Set("a", 100)

		origVal, _ := original.Get("a")
		cloneVal, _ := clone.Get("a")

		assert.Equal(t, 1, origVal)
		assert.Equal(t, 100, cloneVal)
	})
}

func TestSliceIterator(t *testing.T) {
	rows := []*Row{
		NewRow(map[string]interface{}{"id": 1}),
		NewRow(map[string]interface{}{"id": 2}),
		NewRow(map[string]interface{}{"id": 3}),
	}

	iter := NewSliceIterator(rows)
	ctx := context.Background()

	// Read all rows
	count := 0
	for {
		row, err := iter.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		require.NoError(t, err)

		id, _ := row.Get("id")
		assert.Equal(t, count+1, id)
		count++
	}

	assert.Equal(t, 3, count)

	// Stats
	stats := iter.Stats()
	assert.Equal(t, int64(3), stats.RowsRead)
	assert.Equal(t, int64(3), stats.RowsReturned)
}

func TestScanOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	dataSource := createMockDataSource()

	scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

	err := scanOp.Open(ctx)
	require.NoError(t, err)

	// Read all rows
	count := 0
	for {
		row, err := scanOp.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		require.NoError(t, err)
		require.NotNil(t, row)
		count++

		// Check _id is set
		_, ok := row.Get("_id")
		assert.True(t, ok)
	}

	assert.Equal(t, 5, count)
	scanOp.Close()
}

func TestFilterOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	dataSource := createMockDataSource()

	t.Run("EqualityFilter", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		// Filter: status = 200
		condition := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "status"},
			Operator: "=",
			Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
		}

		filterOp := NewFilterOperator(scanOp, condition, logger)

		err := filterOp.Open(ctx)
		require.NoError(t, err)

		// Should return only status=200 rows
		count := 0
		for {
			row, err := filterOp.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)

			status, _ := row.Get("status")
			assert.Equal(t, 200, status)
			count++
		}

		assert.Equal(t, 3, count) // 3 rows have status=200
		filterOp.Close()
	})

	t.Run("ComparisonFilter", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		// Filter: latency > 50
		condition := &ast.BinaryExpression{
			Left:     &ast.FieldReference{Name: "latency"},
			Operator: ">",
			Right:    &ast.Literal{Value: 50.0, LiteralTyp: ast.LiteralTypeFloat},
		}

		filterOp := NewFilterOperator(scanOp, condition, logger)

		err := filterOp.Open(ctx)
		require.NoError(t, err)

		count := 0
		for {
			row, err := filterOp.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)

			latency, _ := row.GetFloat64("latency")
			assert.Greater(t, latency, 50.0)
			count++
		}

		assert.Equal(t, 2, count) // 2 rows have latency > 50
		filterOp.Close()
	})

	t.Run("FunctionFilter", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		// Filter: abs(latency - 50) < 30 (latency between 20 and 80)
		condition := &ast.BinaryExpression{
			Left: &ast.FunctionCall{
				Name: "abs",
				Arguments: []ast.Expression{
					&ast.BinaryExpression{
						Left:     &ast.FieldReference{Name: "latency"},
						Operator: "-",
						Right:    &ast.Literal{Value: 50.0, LiteralTyp: ast.LiteralTypeFloat},
					},
				},
			},
			Operator: "<",
			Right:    &ast.Literal{Value: 31.0, LiteralTyp: ast.LiteralTypeFloat},
		}

		filterOp := NewFilterOperator(scanOp, condition, logger)

		err := filterOp.Open(ctx)
		require.NoError(t, err)

		count := 0
		for {
			_, err := filterOp.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)
			count++
		}

		assert.Equal(t, 4, count) // Rows with latency 20, 30, 50, 80
		filterOp.Close()
	})
}

func TestProjectOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	dataSource := createMockDataSource()

	t.Run("IncludeFields", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		// Project: status, host
		fields := []ast.Expression{
			&ast.FieldReference{Name: "status"},
			&ast.FieldReference{Name: "host"},
		}

		projectOp := NewProjectOperator(scanOp, fields, false, logger)

		err := projectOp.Open(ctx)
		require.NoError(t, err)

		row, err := projectOp.Next(ctx)
		require.NoError(t, err)

		// Should have status and host
		_, hasStatus := row.Get("status")
		_, hasHost := row.Get("host")
		_, hasLatency := row.Get("latency")

		assert.True(t, hasStatus)
		assert.True(t, hasHost)
		assert.False(t, hasLatency) // latency excluded

		projectOp.Close()
	})

	t.Run("ExcludeFields", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		// Project: -latency (exclude latency)
		fields := []ast.Expression{
			&ast.FieldReference{Name: "latency"},
		}

		projectOp := NewProjectOperator(scanOp, fields, true, logger)

		err := projectOp.Open(ctx)
		require.NoError(t, err)

		row, err := projectOp.Next(ctx)
		require.NoError(t, err)

		// Should have status and host but not latency
		_, hasStatus := row.Get("status")
		_, hasHost := row.Get("host")
		_, hasLatency := row.Get("latency")

		assert.True(t, hasStatus)
		assert.True(t, hasHost)
		assert.False(t, hasLatency)

		projectOp.Close()
	})
}

func TestSortOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	dataSource := createMockDataSource()

	t.Run("AscendingSort", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		sortKeys := []*ast.SortKey{
			{Field: &ast.FieldReference{Name: "latency"}, Descending: false},
		}

		sortOp := NewSortOperator(scanOp, sortKeys, logger)

		err := sortOp.Open(ctx)
		require.NoError(t, err)

		var prevLatency float64 = -1
		for {
			row, err := sortOp.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)

			latency, _ := row.GetFloat64("latency")
			assert.GreaterOrEqual(t, latency, prevLatency)
			prevLatency = latency
		}

		sortOp.Close()
	})

	t.Run("DescendingSort", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		sortKeys := []*ast.SortKey{
			{Field: &ast.FieldReference{Name: "latency"}, Descending: true},
		}

		sortOp := NewSortOperator(scanOp, sortKeys, logger)

		err := sortOp.Open(ctx)
		require.NoError(t, err)

		var prevLatency float64 = 1000
		for {
			row, err := sortOp.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)

			latency, _ := row.GetFloat64("latency")
			assert.LessOrEqual(t, latency, prevLatency)
			prevLatency = latency
		}

		sortOp.Close()
	})
}

func TestLimitOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	dataSource := createMockDataSource()

	scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

	limitOp := NewLimitOperator(scanOp, 3, logger)

	err := limitOp.Open(ctx)
	require.NoError(t, err)

	count := 0
	for {
		_, err := limitOp.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		require.NoError(t, err)
		count++
	}

	assert.Equal(t, 3, count) // Limited to 3 rows
	limitOp.Close()
}

func TestAggregationOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	dataSource := createMockDataSource()

	t.Run("GlobalCount", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		aggregations := []*ast.Aggregation{
			{
				Func:  &ast.FunctionCall{Name: "count", Arguments: nil},
				Alias: "total",
			},
		}

		aggOp := NewAggregationOperator(scanOp, nil, aggregations, physical.HashAggregation, logger)

		err := aggOp.Open(ctx)
		require.NoError(t, err)

		row, err := aggOp.Next(ctx)
		require.NoError(t, err)

		total, ok := row.Get("total")
		assert.True(t, ok)
		assert.Equal(t, int64(5), total)

		// No more rows
		_, err = aggOp.Next(ctx)
		assert.Equal(t, ErrNoMoreRows, err)

		aggOp.Close()
	})

	t.Run("GlobalAggregations", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		aggregations := []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name:      "sum",
					Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
				},
				Alias: "total_latency",
			},
			{
				Func: &ast.FunctionCall{
					Name:      "avg",
					Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
				},
				Alias: "avg_latency",
			},
			{
				Func: &ast.FunctionCall{
					Name:      "min",
					Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
				},
				Alias: "min_latency",
			},
			{
				Func: &ast.FunctionCall{
					Name:      "max",
					Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
				},
				Alias: "max_latency",
			},
		}

		aggOp := NewAggregationOperator(scanOp, nil, aggregations, physical.HashAggregation, logger)

		err := aggOp.Open(ctx)
		require.NoError(t, err)

		row, err := aggOp.Next(ctx)
		require.NoError(t, err)

		// Total: 50 + 100 + 30 + 20 + 80 = 280
		totalLatency, _ := row.Get("total_latency")
		assert.Equal(t, 280.0, totalLatency)

		// Avg: 280 / 5 = 56
		avgLatency, _ := row.Get("avg_latency")
		assert.Equal(t, 56.0, avgLatency)

		// Min: 20
		minLatency, _ := row.Get("min_latency")
		assert.Equal(t, 20.0, minLatency)

		// Max: 100
		maxLatency, _ := row.Get("max_latency")
		assert.Equal(t, 100.0, maxLatency)

		aggOp.Close()
	})

	t.Run("GroupedAggregation", func(t *testing.T) {
		scanOp := NewScanOperator(dataSource, "logs", []byte(`{}`), 0, 10, logger)

		groupBy := []ast.Expression{
			&ast.FieldReference{Name: "host"},
		}

		aggregations := []*ast.Aggregation{
			{
				Func: &ast.FunctionCall{
					Name:      "count",
					Arguments: nil,
				},
				Alias: "count",
			},
			{
				Func: &ast.FunctionCall{
					Name:      "avg",
					Arguments: []ast.Expression{&ast.FieldReference{Name: "latency"}},
				},
				Alias: "avg_latency",
			},
		}

		aggOp := NewAggregationOperator(scanOp, groupBy, aggregations, physical.HashAggregation, logger)

		err := aggOp.Open(ctx)
		require.NoError(t, err)

		results := make(map[string]int64)
		for {
			row, err := aggOp.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)

			host, _ := row.GetString("host")
			count, _ := row.GetInt64("count")
			results[host] = count
		}

		// web1: 2 docs, web2: 2 docs, web3: 1 doc
		assert.Equal(t, int64(2), results["web1"])
		assert.Equal(t, int64(2), results["web2"])
		assert.Equal(t, int64(1), results["web3"])

		aggOp.Close()
	})
}

// Tier 1 Operator Tests

func TestDedupOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	rows := []*Row{
		NewRow(map[string]interface{}{"host": "web1", "status": 200}),
		NewRow(map[string]interface{}{"host": "web1", "status": 500}),
		NewRow(map[string]interface{}{"host": "web2", "status": 200}),
		NewRow(map[string]interface{}{"host": "web1", "status": 404}),
		NewRow(map[string]interface{}{"host": "web2", "status": 500}),
	}

	t.Run("DedupCount1", func(t *testing.T) {
		iter := NewSliceIterator(rows)
		fields := []ast.Expression{&ast.FieldReference{Name: "host"}}
		dedupOp := NewDedupOperator(iter, fields, 1, false, logger)

		err := dedupOp.Open(ctx)
		require.NoError(t, err)

		// Should return 2 rows (one for web1, one for web2)
		count := 0
		for {
			_, err := dedupOp.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)
			count++
		}

		assert.Equal(t, 2, count) // Only 2 unique hosts
		dedupOp.Close()
	})

	t.Run("DedupCount2", func(t *testing.T) {
		iter := NewSliceIterator(rows)
		fields := []ast.Expression{&ast.FieldReference{Name: "host"}}
		dedupOp := NewDedupOperator(iter, fields, 2, false, logger)

		err := dedupOp.Open(ctx)
		require.NoError(t, err)

		// Should return 4 rows (two for web1, two for web2)
		count := 0
		for {
			_, err := dedupOp.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)
			count++
		}

		assert.Equal(t, 4, count) // 2 web1 + 2 web2
		dedupOp.Close()
	})
}

func TestTopOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	rows := []*Row{
		NewRow(map[string]interface{}{"status": 200}),
		NewRow(map[string]interface{}{"status": 200}),
		NewRow(map[string]interface{}{"status": 200}),
		NewRow(map[string]interface{}{"status": 500}),
		NewRow(map[string]interface{}{"status": 500}),
		NewRow(map[string]interface{}{"status": 404}),
	}

	iter := NewSliceIterator(rows)
	fields := []ast.Expression{&ast.FieldReference{Name: "status"}}
	topOp := NewTopOperator(iter, fields, 2, nil, true, true, logger)

	err := topOp.Open(ctx)
	require.NoError(t, err)

	// First result should be status=200 (3 occurrences)
	row1, err := topOp.Next(ctx)
	require.NoError(t, err)
	status1, _ := row1.Get("status")
	count1, _ := row1.GetInt64("count")
	assert.Equal(t, "200", status1) // String due to key parsing
	assert.Equal(t, int64(3), count1)

	// Second result should be status=500 (2 occurrences)
	row2, err := topOp.Next(ctx)
	require.NoError(t, err)
	status2, _ := row2.Get("status")
	count2, _ := row2.GetInt64("count")
	assert.Equal(t, "500", status2)
	assert.Equal(t, int64(2), count2)

	// No more results
	_, err = topOp.Next(ctx)
	assert.Equal(t, ErrNoMoreRows, err)

	topOp.Close()
}

func TestRareOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	rows := []*Row{
		NewRow(map[string]interface{}{"status": 200}),
		NewRow(map[string]interface{}{"status": 200}),
		NewRow(map[string]interface{}{"status": 200}),
		NewRow(map[string]interface{}{"status": 500}),
		NewRow(map[string]interface{}{"status": 500}),
		NewRow(map[string]interface{}{"status": 404}),
	}

	iter := NewSliceIterator(rows)
	fields := []ast.Expression{&ast.FieldReference{Name: "status"}}
	rareOp := NewRareOperator(iter, fields, 2, nil, true, false, logger)

	err := rareOp.Open(ctx)
	require.NoError(t, err)

	// First result should be status=404 (1 occurrence - rarest)
	row1, err := rareOp.Next(ctx)
	require.NoError(t, err)
	status1, _ := row1.Get("status")
	count1, _ := row1.GetInt64("count")
	assert.Equal(t, "404", status1)
	assert.Equal(t, int64(1), count1)

	// Second result should be status=500 (2 occurrences)
	row2, err := rareOp.Next(ctx)
	require.NoError(t, err)
	status2, _ := row2.Get("status")
	count2, _ := row2.GetInt64("count")
	assert.Equal(t, "500", status2)
	assert.Equal(t, int64(2), count2)

	rareOp.Close()
}

func TestEvalOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	rows := []*Row{
		NewRow(map[string]interface{}{"price": 100.0, "quantity": 5}),
		NewRow(map[string]interface{}{"price": 50.0, "quantity": 10}),
	}

	iter := NewSliceIterator(rows)
	assignments := []*ast.EvalAssignment{
		{
			Field: "total",
			Expression: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "price"},
				Operator: "*",
				Right:    &ast.FieldReference{Name: "quantity"},
			},
		},
	}
	evalOp := NewEvalOperator(iter, assignments, logger)

	err := evalOp.Open(ctx)
	require.NoError(t, err)

	// First row: total = 100 * 5 = 500
	row1, err := evalOp.Next(ctx)
	require.NoError(t, err)
	total1, _ := row1.GetFloat64("total")
	assert.Equal(t, 500.0, total1)

	// Second row: total = 50 * 10 = 500
	row2, err := evalOp.Next(ctx)
	require.NoError(t, err)
	total2, _ := row2.GetFloat64("total")
	assert.Equal(t, 500.0, total2)

	evalOp.Close()
}

func TestRenameOperator(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	rows := []*Row{
		NewRow(map[string]interface{}{"host": "web1", "status": 200}),
		NewRow(map[string]interface{}{"host": "web2", "status": 500}),
	}

	iter := NewSliceIterator(rows)
	assignments := []*ast.RenameAssignment{
		{OldName: "host", NewName: "server"},
	}
	renameOp := NewRenameOperator(iter, assignments, logger)

	err := renameOp.Open(ctx)
	require.NoError(t, err)

	// First row should have 'server' field instead of 'host'
	row1, err := renameOp.Next(ctx)
	require.NoError(t, err)

	server, ok := row1.GetString("server")
	assert.True(t, ok)
	assert.Equal(t, "web1", server)

	_, hasHost := row1.Get("host")
	assert.False(t, hasHost, "Old field 'host' should be removed")

	renameOp.Close()
}

func TestExecutor(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	dataSource := createMockDataSource()
	translator := dsl.NewTranslator()

	executor := NewExecutor(dataSource, translator, logger)
	schema := createTestSchema()

	t.Run("SimpleScan", func(t *testing.T) {
		plan := &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		}

		result, err := executor.Execute(ctx, plan)
		require.NoError(t, err)
		require.NotNil(t, result)

		count := 0
		for {
			_, err := result.Rows.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)
			count++
		}

		assert.Equal(t, 5, count)
		result.Rows.Close()
	})

	t.Run("ScanWithFilter", func(t *testing.T) {
		scanPlan := &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		}

		filterPlan := &physical.PhysicalFilter{
			Condition: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: "=",
				Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
			},
			Input: scanPlan,
		}

		result, err := executor.Execute(ctx, filterPlan)
		require.NoError(t, err)
		require.NotNil(t, result)

		count := 0
		for {
			row, err := result.Rows.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)

			status, _ := row.Get("status")
			assert.Equal(t, 200, status)
			count++
		}

		assert.Equal(t, 3, count)
		result.Rows.Close()
	})

	t.Run("ComplexPipeline", func(t *testing.T) {
		// source=logs | where status = 200 | sort latency | head 2
		scanPlan := &physical.PhysicalScan{
			Source:       "logs",
			OutputSchema: schema,
		}

		filterPlan := &physical.PhysicalFilter{
			Condition: &ast.BinaryExpression{
				Left:     &ast.FieldReference{Name: "status"},
				Operator: "=",
				Right:    &ast.Literal{Value: 200, LiteralTyp: ast.LiteralTypeInt},
			},
			Input: scanPlan,
		}

		sortPlan := &physical.PhysicalSort{
			SortKeys: []*ast.SortKey{
				{Field: &ast.FieldReference{Name: "latency"}, Descending: false},
			},
			Input: filterPlan,
		}

		limitPlan := &physical.PhysicalLimit{
			Count: 2,
			Input: sortPlan,
		}

		result, err := executor.Execute(ctx, limitPlan)
		require.NoError(t, err)
		require.NotNil(t, result)

		rows := make([]*Row, 0)
		for {
			row, err := result.Rows.Next(ctx)
			if err == ErrNoMoreRows {
				break
			}
			require.NoError(t, err)
			rows = append(rows, row)
		}

		// Should have 2 rows, sorted by latency ascending
		assert.Len(t, rows, 2)

		lat1, _ := rows[0].GetFloat64("latency")
		lat2, _ := rows[1].GetFloat64("latency")
		assert.Less(t, lat1, lat2)

		result.Rows.Close()
	})
}
