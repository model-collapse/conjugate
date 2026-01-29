// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"

	"go.uber.org/zap"
)

// scanOperator reads data from a data source
type scanOperator struct {
	dataSource DataSource
	index      string
	queryDSL   []byte
	from       int
	size       int
	logger     *zap.Logger

	// Runtime state
	ctx     context.Context
	rows    []*Row
	index_  int // Index into rows
	stats   *IteratorStats
	opened  bool
	closed  bool
}

// NewScanOperator creates a new scan operator
func NewScanOperator(dataSource DataSource, index string, queryDSL []byte, from, size int, logger *zap.Logger) *scanOperator {
	return &scanOperator{
		dataSource: dataSource,
		index:      index,
		queryDSL:   queryDSL,
		from:       from,
		size:       size,
		logger:     logger,
		stats:      &IteratorStats{},
	}
}

// Open executes the query and loads results
func (s *scanOperator) Open(ctx context.Context) error {
	if s.opened {
		return nil
	}

	s.ctx = ctx
	s.logger.Debug("Opening scan operator",
		zap.String("index", s.index),
		zap.Int("from", s.from),
		zap.Int("size", s.size))

	// Execute search
	result, err := s.dataSource.Search(ctx, s.index, s.queryDSL, s.from, s.size)
	if err != nil {
		return err
	}

	// Convert hits to rows
	s.rows = make([]*Row, len(result.Hits))
	for i, hit := range result.Hits {
		row := NewRow(hit.Source)
		// Add metadata fields
		row.Set("_id", hit.ID)
		row.Set("_score", hit.Score)
		s.rows[i] = row
	}

	s.stats.TookMillis = result.TookMillis
	s.index_ = 0
	s.opened = true

	s.logger.Debug("Scan operator opened",
		zap.Int("rows_loaded", len(s.rows)),
		zap.Int64("took_ms", result.TookMillis))

	return nil
}

// Next returns the next row
func (s *scanOperator) Next(ctx context.Context) (*Row, error) {
	if s.closed {
		return nil, ErrClosed
	}

	if !s.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if s.index_ >= len(s.rows) {
		return nil, ErrNoMoreRows
	}

	row := s.rows[s.index_]
	s.index_++
	s.stats.RowsRead++
	s.stats.RowsReturned++

	return row, nil
}

// Close releases resources
func (s *scanOperator) Close() error {
	s.closed = true
	s.rows = nil
	return nil
}

// Stats returns execution statistics
func (s *scanOperator) Stats() *IteratorStats {
	return s.stats
}
