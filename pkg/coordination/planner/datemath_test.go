package planner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateMathParser_Parse(t *testing.T) {
	// Fixed time for testing: 2024-06-15 12:30:45 UTC
	fixedTime := time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC)
	parser := &DateMathParser{
		nowFunc: func() time.Time { return fixedTime },
	}

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		// Non-string inputs (pass through)
		{
			name:     "integer value",
			input:    12345,
			expected: "",
		},
		{
			name:     "nil value",
			input:    nil,
			expected: "",
		},

		// Non-date math strings (pass through)
		{
			name:     "regular string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "number string",
			input:    "12345",
			expected: "12345",
		},

		// Simple "now"
		{
			name:     "now",
			input:    "now",
			expected: "2024-06-15T12:30:45.000Z",
		},

		// Subtraction
		{
			name:     "now-1d",
			input:    "now-1d",
			expected: "2024-06-14T12:30:45.000Z",
		},
		{
			name:     "now-7d",
			input:    "now-7d",
			expected: "2024-06-08T12:30:45.000Z",
		},
		{
			name:     "now-1h",
			input:    "now-1h",
			expected: "2024-06-15T11:30:45.000Z",
		},
		{
			name:     "now-30m",
			input:    "now-30m",
			expected: "2024-06-15T12:00:45.000Z",
		},
		{
			name:     "now-1M",
			input:    "now-1M",
			expected: "2024-05-15T12:30:45.000Z",
		},
		{
			name:     "now-1y",
			input:    "now-1y",
			expected: "2023-06-15T12:30:45.000Z",
		},

		// Addition
		{
			name:     "now+1d",
			input:    "now+1d",
			expected: "2024-06-16T12:30:45.000Z",
		},
		{
			name:     "now+3h",
			input:    "now+3h",
			expected: "2024-06-15T15:30:45.000Z",
		},

		// Rounding
		{
			name:     "now/d",
			input:    "now/d",
			expected: "2024-06-15T00:00:00.000Z",
		},
		{
			name:     "now/h",
			input:    "now/h",
			expected: "2024-06-15T12:00:00.000Z",
		},
		{
			name:     "now/M",
			input:    "now/M",
			expected: "2024-06-01T00:00:00.000Z",
		},

		// Complex expressions
		{
			name:     "now-1d/d",
			input:    "now-1d/d",
			expected: "2024-06-14T00:00:00.000Z",
		},
		{
			name:     "now-7d/d",
			input:    "now-7d/d",
			expected: "2024-06-08T00:00:00.000Z",
		},
		{
			name:     "now+1M/M",
			input:    "now+1M/M",
			expected: "2024-07-01T00:00:00.000Z",
		},
		{
			name:     "now/M+1d",
			input:    "now/M+1d",
			expected: "2024-06-02T00:00:00.000Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.Parse(tt.input)

			if tt.expected == "" {
				// Non-date math should pass through unchanged
				assert.Equal(t, tt.input, result)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDateMathParser_ParseRangeParams(t *testing.T) {
	// Fixed time: 2024-06-15 12:00:00 UTC
	fixedTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	parser := &DateMathParser{
		nowFunc: func() time.Time { return fixedTime },
	}

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "gte with date math",
			input: map[string]interface{}{
				"gte": "now-7d",
			},
			expected: map[string]interface{}{
				"gte": "2024-06-08T12:00:00.000Z",
			},
		},
		{
			name: "gte and lte with date math",
			input: map[string]interface{}{
				"gte": "now-30d",
				"lte": "now",
			},
			expected: map[string]interface{}{
				"gte": "2024-05-16T12:00:00.000Z",
				"lte": "2024-06-15T12:00:00.000Z",
			},
		},
		{
			name: "mixed date math and literal",
			input: map[string]interface{}{
				"gte": "now-1d",
				"lte": 100,
			},
			expected: map[string]interface{}{
				"gte": "2024-06-14T12:00:00.000Z",
				"lte": 100,
			},
		},
		{
			name: "no date math",
			input: map[string]interface{}{
				"gte": 10,
				"lte": 100,
			},
			expected: map[string]interface{}{
				"gte": 10,
				"lte": 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.ParseRangeParams(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDateMathParser_RoundDown(t *testing.T) {
	parser := NewDateMathParser()

	// Test time: 2024-06-15 14:35:42 (Saturday)
	testTime := time.Date(2024, 6, 15, 14, 35, 42, 0, time.UTC)

	tests := []struct {
		unit     string
		expected time.Time
	}{
		{
			unit:     "s",
			expected: time.Date(2024, 6, 15, 14, 35, 42, 0, time.UTC),
		},
		{
			unit:     "m",
			expected: time.Date(2024, 6, 15, 14, 35, 0, 0, time.UTC),
		},
		{
			unit:     "h",
			expected: time.Date(2024, 6, 15, 14, 0, 0, 0, time.UTC),
		},
		{
			unit:     "d",
			expected: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			unit:     "w",
			expected: time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC), // Monday
		},
		{
			unit:     "M",
			expected: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			unit:     "y",
			expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run("round_"+tt.unit, func(t *testing.T) {
			result := parser.roundDown(testTime, tt.unit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDateMathParser_AddDuration(t *testing.T) {
	parser := NewDateMathParser()

	// Base time: 2024-06-15 12:00:00
	baseTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		num      int
		unit     string
		expected time.Time
	}{
		{
			name:     "add 1 second",
			num:      1,
			unit:     "s",
			expected: time.Date(2024, 6, 15, 12, 0, 1, 0, time.UTC),
		},
		{
			name:     "add 30 minutes",
			num:      30,
			unit:     "m",
			expected: time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC),
		},
		{
			name:     "add 3 hours",
			num:      3,
			unit:     "h",
			expected: time.Date(2024, 6, 15, 15, 0, 0, 0, time.UTC),
		},
		{
			name:     "add 7 days",
			num:      7,
			unit:     "d",
			expected: time.Date(2024, 6, 22, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "add 2 weeks",
			num:      2,
			unit:     "w",
			expected: time.Date(2024, 6, 29, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "add 1 month",
			num:      1,
			unit:     "M",
			expected: time.Date(2024, 7, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "add 1 year",
			num:      1,
			unit:     "y",
			expected: time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "subtract 1 day",
			num:      -1,
			unit:     "d",
			expected: time.Date(2024, 6, 14, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.addDuration(baseTime, tt.num, tt.unit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDateMathParser_EvaluateNowExpression(t *testing.T) {
	// Fixed time: 2024-06-15 12:30:00 UTC
	fixedTime := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)
	parser := &DateMathParser{
		nowFunc: func() time.Time { return fixedTime },
	}

	tests := []struct {
		name     string
		expr     string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "simple now",
			expr:     "now",
			expected: fixedTime,
		},
		{
			name:     "now-1d",
			expr:     "now-1d",
			expected: time.Date(2024, 6, 14, 12, 30, 0, 0, time.UTC),
		},
		{
			name:     "now-7d/d",
			expr:     "now-7d/d",
			expected: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "now+1M/M",
			expr:     "now+1M/M",
			expected: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.evaluateNowExpression(tt.expr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDateMathParser_Benchmark_Scenario(t *testing.T) {
	// This test simulates the big5 benchmark range query:
	// "gte": "now-7d" on @timestamp field

	// Fixed time: 2024-06-15 12:00:00 UTC
	fixedTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	parser := &DateMathParser{
		nowFunc: func() time.Time { return fixedTime },
	}

	// Simulate range query params from big5 benchmark
	input := map[string]interface{}{
		"gte": "now-7d",
	}

	result := parser.ParseRangeParams(input)

	expected := map[string]interface{}{
		"gte": "2024-06-08T12:00:00.000Z",
	}

	assert.Equal(t, expected, result)

	// Verify it's a proper timestamp that Diagon can use
	resultTimestamp := result["gte"].(string)
	_, err := time.Parse("2006-01-02T15:04:05.000Z", resultTimestamp)
	require.NoError(t, err, "Result should be valid ISO8601 timestamp")
}

func BenchmarkDateMathParser_Parse(b *testing.B) {
	parser := NewDateMathParser()

	b.Run("now-7d", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parser.Parse("now-7d")
		}
	})

	b.Run("now-30d/d", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parser.Parse("now-30d/d")
		}
	})

	b.Run("passthrough_integer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parser.Parse(12345)
		}
	})
}

func BenchmarkDateMathParser_ParseRangeParams(b *testing.B) {
	parser := NewDateMathParser()

	params := map[string]interface{}{
		"gte": "now-7d",
		"lte": "now",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.ParseRangeParams(params)
	}
}
