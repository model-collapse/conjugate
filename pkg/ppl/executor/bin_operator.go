// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"math"
	"time"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// binOperator bins field values into buckets
type binOperator struct {
	input  Operator
	field  ast.Expression
	span   *ast.TimeSpan
	bins   int
	logger *zap.Logger

	ctx       context.Context
	fieldName string
	stats     *IteratorStats
	opened    bool
	closed    bool
}

// NewBinOperator creates a new bin operator
func NewBinOperator(
	input Operator,
	field ast.Expression,
	span *ast.TimeSpan,
	bins int,
	logger *zap.Logger,
) *binOperator {
	fieldName := ""
	if ref, ok := field.(*ast.FieldReference); ok {
		fieldName = ref.Name
	}

	return &binOperator{
		input:     input,
		field:     field,
		span:      span,
		bins:      bins,
		logger:    logger,
		fieldName: fieldName,
		stats:     &IteratorStats{},
	}
}

// Open initializes the operator
func (b *binOperator) Open(ctx context.Context) error {
	if b.opened {
		return nil
	}

	b.ctx = ctx
	b.logger.Debug("Opening bin operator",
		zap.String("field", b.fieldName),
		zap.Int("bins", b.bins))

	if err := b.input.Open(ctx); err != nil {
		return err
	}

	b.opened = true
	return nil
}

// Next returns the next row with binned field value
func (b *binOperator) Next(ctx context.Context) (*Row, error) {
	if b.closed {
		return nil, ErrClosed
	}

	if !b.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	row, err := b.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	b.stats.RowsRead++

	// Get the field value
	val, ok := row.Get(b.fieldName)
	if !ok {
		// Field doesn't exist, return row unchanged
		b.stats.RowsReturned++
		return row, nil
	}

	// Bin the value
	binnedValue := b.binValue(val)
	row.Set(b.fieldName, binnedValue)

	b.stats.RowsReturned++
	return row, nil
}

// binValue bins a value based on span or bins configuration
func (b *binOperator) binValue(val interface{}) interface{} {
	if b.span != nil {
		return b.binByTimeSpan(val)
	} else if b.bins > 0 {
		return b.binByNumericRange(val)
	}
	return val
}

// binByTimeSpan bins a time value by the specified span
func (b *binOperator) binByTimeSpan(val interface{}) interface{} {
	// Try to parse as time
	var t time.Time
	switch v := val.(type) {
	case time.Time:
		t = v
	case string:
		// Try common time formats
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		var err error
		for _, format := range formats {
			t, err = time.Parse(format, v)
			if err == nil {
				break
			}
		}
		if err != nil {
			return val // Can't parse, return unchanged
		}
	default:
		return val // Not a time value
	}

	// Compute bucket based on span
	spanDuration := b.spanToDuration()
	if spanDuration == 0 {
		return val
	}

	// Truncate to bucket
	unixNanos := t.UnixNano()
	bucketNanos := (unixNanos / int64(spanDuration)) * int64(spanDuration)
	return time.Unix(0, bucketNanos).UTC()
}

// spanToDuration converts TimeSpan to time.Duration
func (b *binOperator) spanToDuration() time.Duration {
	if b.span == nil {
		return 0
	}

	switch b.span.Unit {
	case "s", "sec", "second", "seconds":
		return time.Duration(b.span.Value) * time.Second
	case "m", "min", "minute", "minutes":
		return time.Duration(b.span.Value) * time.Minute
	case "h", "hr", "hour", "hours":
		return time.Duration(b.span.Value) * time.Hour
	case "d", "day", "days":
		return time.Duration(b.span.Value) * 24 * time.Hour
	case "w", "week", "weeks":
		return time.Duration(b.span.Value) * 7 * 24 * time.Hour
	case "mon", "month", "months":
		return time.Duration(b.span.Value) * 30 * 24 * time.Hour // Approximate
	case "y", "year", "years":
		return time.Duration(b.span.Value) * 365 * 24 * time.Hour // Approximate
	default:
		return 0
	}
}

// binByNumericRange bins a numeric value into one of n buckets
func (b *binOperator) binByNumericRange(val interface{}) interface{} {
	num, ok := toNumber(val)
	if !ok {
		return val // Not numeric
	}

	// Simple binning: round down to nearest bucket
	// For proper implementation, would need to scan data first to find min/max
	bucketSize := 1.0 // Default bucket size
	if b.bins > 0 {
		// Assume data range 0-100 for now (would need stats for actual range)
		bucketSize = 100.0 / float64(b.bins)
	}

	bucket := math.Floor(num/bucketSize) * bucketSize
	return bucket
}

// Close releases resources
func (b *binOperator) Close() error {
	b.closed = true
	return b.input.Close()
}

// Stats returns execution statistics
func (b *binOperator) Stats() *IteratorStats {
	return b.stats
}
