// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"errors"
)

// Common errors
var (
	ErrClosed       = errors.New("iterator closed")
	ErrNoMoreRows   = errors.New("no more rows")
	ErrInvalidField = errors.New("invalid field")
)

// Row represents a single row of data with named fields
type Row struct {
	fields map[string]interface{}
}

// NewRow creates a new row from a map
func NewRow(fields map[string]interface{}) *Row {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	return &Row{fields: fields}
}

// Get returns a field value
func (r *Row) Get(name string) (interface{}, bool) {
	v, ok := r.fields[name]
	return v, ok
}

// GetString returns a field as string
func (r *Row) GetString(name string) (string, bool) {
	v, ok := r.fields[name]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// GetInt64 returns a field as int64
func (r *Row) GetInt64(name string) (int64, bool) {
	v, ok := r.fields[name]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case int64:
		return n, true
	case int:
		return int64(n), true
	case int32:
		return int64(n), true
	case float64:
		return int64(n), true
	default:
		return 0, false
	}
}

// GetFloat64 returns a field as float64
func (r *Row) GetFloat64(name string) (float64, bool) {
	v, ok := r.fields[name]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int64:
		return float64(n), true
	case int:
		return float64(n), true
	default:
		return 0, false
	}
}

// GetBool returns a field as bool
func (r *Row) GetBool(name string) (bool, bool) {
	v, ok := r.fields[name]
	if !ok {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

// Set sets a field value
func (r *Row) Set(name string, value interface{}) {
	r.fields[name] = value
}

// Delete removes a field
func (r *Row) Delete(name string) {
	delete(r.fields, name)
}

// Fields returns all field names
func (r *Row) Fields() []string {
	names := make([]string, 0, len(r.fields))
	for name := range r.fields {
		names = append(names, name)
	}
	return names
}

// ToMap returns the underlying map
func (r *Row) ToMap() map[string]interface{} {
	result := make(map[string]interface{}, len(r.fields))
	for k, v := range r.fields {
		result[k] = v
	}
	return result
}

// Clone creates a copy of the row
func (r *Row) Clone() *Row {
	return NewRow(r.ToMap())
}

// RowIterator provides streaming access to rows
type RowIterator interface {
	// Next advances to the next row and returns it
	// Returns nil, ErrNoMoreRows when iteration is complete
	Next(ctx context.Context) (*Row, error)

	// Close releases resources
	Close() error

	// Stats returns execution statistics (optional)
	Stats() *IteratorStats
}

// IteratorStats contains execution statistics
type IteratorStats struct {
	RowsRead     int64
	RowsReturned int64
	TookMillis   int64
}

// SliceIterator wraps a slice of rows as an iterator
type SliceIterator struct {
	rows   []*Row
	index  int
	closed bool
	stats  *IteratorStats
}

// NewSliceIterator creates an iterator from a slice of rows
func NewSliceIterator(rows []*Row) *SliceIterator {
	return &SliceIterator{
		rows:  rows,
		index: 0,
		stats: &IteratorStats{},
	}
}

// Open initializes the iterator (no-op for slice iterator)
func (s *SliceIterator) Open(ctx context.Context) error {
	return nil
}

// Next returns the next row
func (s *SliceIterator) Next(ctx context.Context) (*Row, error) {
	if s.closed {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if s.index >= len(s.rows) {
		return nil, ErrNoMoreRows
	}

	row := s.rows[s.index]
	s.index++
	s.stats.RowsRead++
	s.stats.RowsReturned++

	return row, nil
}

// Close releases resources
func (s *SliceIterator) Close() error {
	s.closed = true
	return nil
}

// Stats returns execution statistics
func (s *SliceIterator) Stats() *IteratorStats {
	return s.stats
}

// EmptyIterator returns no rows
type EmptyIterator struct {
	stats *IteratorStats
}

// NewEmptyIterator creates an empty iterator
func NewEmptyIterator() *EmptyIterator {
	return &EmptyIterator{stats: &IteratorStats{}}
}

// Next always returns ErrNoMoreRows
func (e *EmptyIterator) Next(ctx context.Context) (*Row, error) {
	return nil, ErrNoMoreRows
}

// Close is a no-op
func (e *EmptyIterator) Close() error {
	return nil
}

// Stats returns execution statistics
func (e *EmptyIterator) Stats() *IteratorStats {
	return e.stats
}

// Result represents the complete result of a query execution
type Result struct {
	// Rows is the iterator for streaming results
	Rows RowIterator

	// Aggregations contains aggregation results (if any)
	Aggregations map[string]*AggregationValue

	// Metadata
	TookMillis int64
	TotalHits  int64
	MaxScore   float64
}

// AggregationValue represents an aggregation result
type AggregationValue struct {
	Type string

	// Metric aggregation values
	Value      float64
	Count      int64
	Min        float64
	Max        float64
	Avg        float64
	Sum        float64
	Variance   float64
	StdDev     float64

	// Bucket aggregation values
	Buckets []*BucketValue

	// Percentile values
	Percentiles map[string]float64
}

// BucketValue represents a bucket in bucket aggregations
type BucketValue struct {
	Key        string
	NumericKey float64
	DocCount   int64
	SubAggs    map[string]*AggregationValue
}
