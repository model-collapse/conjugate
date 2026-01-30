// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
)

// MockOperator is a test mock that returns a fixed set of rows
type MockOperator struct {
	rows   []*Row
	index  int
	opened bool
	closed bool
	stats  *IteratorStats
}

// Open initializes the mock operator
func (m *MockOperator) Open(ctx context.Context) error {
	m.opened = true
	m.index = 0
	if m.stats == nil {
		m.stats = &IteratorStats{}
	}
	return nil
}

// Next returns the next row from the mock data
func (m *MockOperator) Next(ctx context.Context) (*Row, error) {
	if m.closed {
		return nil, ErrClosed
	}

	if !m.opened {
		return nil, ErrClosed
	}

	if m.index >= len(m.rows) {
		return nil, ErrNoMoreRows
	}

	row := m.rows[m.index]
	m.index++
	m.stats.RowsRead++
	m.stats.RowsReturned++

	return row, nil
}

// Close closes the mock operator
func (m *MockOperator) Close() error {
	m.closed = true
	return nil
}

// Stats returns execution statistics
func (m *MockOperator) Stats() *IteratorStats {
	return m.stats
}
