// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"go.uber.org/zap"
)

// joinOperator implements hash join algorithm
// Build phase: builds hash table from right side
// Probe phase: probes hash table with left side
type joinOperator struct {
	input      Operator
	right      Operator
	joinType   ast.JoinType
	joinField  string
	rightField string
	logger     *zap.Logger

	// Hash table built from right side
	// Key: join field value (as string)
	// Value: list of rows with that key value
	hashTable map[string][]*Row

	// State
	ctx            context.Context
	stats          *IteratorStats
	opened         bool
	closed         bool
	buildCompleted bool
	currentLeft    *Row           // Current left row being processed
	rightMatches   []*Row         // Matching right rows for current left
	matchIndex     int            // Current index in rightMatches
	leftDone       bool           // True when left side exhausted
}

// NewJoinOperator creates a new join operator
func NewJoinOperator(
	input Operator,
	right Operator,
	joinType ast.JoinType,
	joinField string,
	rightField string,
	logger *zap.Logger,
) *joinOperator {
	return &joinOperator{
		input:      input,
		right:      right,
		joinType:   joinType,
		joinField:  joinField,
		rightField: rightField,
		logger:     logger,
		stats:      &IteratorStats{},
		hashTable:  make(map[string][]*Row),
	}
}

// Open initializes the operator and builds the hash table
func (j *joinOperator) Open(ctx context.Context) error {
	if j.opened {
		return nil
	}

	j.ctx = ctx
	j.logger.Debug("Opening join operator",
		zap.String("join_type", string(j.joinType)),
		zap.String("join_field", j.joinField))

	// Open both sides
	if err := j.input.Open(ctx); err != nil {
		return fmt.Errorf("failed to open left side: %w", err)
	}

	if err := j.right.Open(ctx); err != nil {
		return fmt.Errorf("failed to open right side: %w", err)
	}

	// Build phase: Read all rows from right side and build hash table
	j.logger.Debug("Building hash table from right side")
	buildCount := 0
	for {
		row, err := j.right.Next(ctx)
		if err == ErrNoMoreRows {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading right side: %w", err)
		}

		// Get join key from right side
		keyValue, exists := row.Get(j.rightField)
		if !exists {
			// Right row doesn't have join field, skip it
			j.logger.Debug("Right row missing join field",
				zap.String("field", j.rightField))
			continue
		}

		// Convert key to string for hashing
		key := fmt.Sprintf("%v", keyValue)

		// Add to hash table
		j.hashTable[key] = append(j.hashTable[key], row)
		buildCount++
	}

	j.buildCompleted = true
	j.logger.Debug("Hash table built",
		zap.Int("right_rows", buildCount),
		zap.Int("unique_keys", len(j.hashTable)))

	j.opened = true
	return nil
}

// Next returns the next joined row
func (j *joinOperator) Next(ctx context.Context) (*Row, error) {
	if j.closed {
		return nil, ErrClosed
	}

	if !j.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	for {
		// If we have matches for the current left row, return them
		if j.currentLeft != nil && j.matchIndex < len(j.rightMatches) {
			joinedRow := j.joinRows(j.currentLeft, j.rightMatches[j.matchIndex])
			j.matchIndex++
			j.stats.RowsReturned++
			return joinedRow, nil
		}

		// If we've exhausted matches for current left row
		// Check if this was a left join with no matches
		if j.currentLeft != nil && len(j.rightMatches) == 0 && j.joinType == ast.JoinTypeLeft {
			// Left join: emit left row with NULLs for right side
			joinedRow := j.joinRows(j.currentLeft, nil)
			j.currentLeft = nil // Move to next left row
			j.stats.RowsReturned++
			return joinedRow, nil
		}

		// Need next left row
		leftRow, err := j.input.Next(ctx)
		if err == ErrNoMoreRows {
			j.leftDone = true
			return nil, ErrNoMoreRows
		}
		if err != nil {
			return nil, err
		}

		j.stats.RowsRead++

		// Get join key from left side
		keyValue, exists := leftRow.Get(j.joinField)
		if !exists {
			// Left row doesn't have join field
			if j.joinType == ast.JoinTypeLeft {
				// For left join, still emit the row with NULLs
				joinedRow := j.joinRows(leftRow, nil)
				j.stats.RowsReturned++
				return joinedRow, nil
			}
			// For inner join, skip this row
			j.logger.Debug("Left row missing join field",
				zap.String("field", j.joinField))
			continue
		}

		// Convert key to string
		key := fmt.Sprintf("%v", keyValue)

		// Probe hash table
		matches := j.hashTable[key]

		j.currentLeft = leftRow
		j.rightMatches = matches
		j.matchIndex = 0

		// If inner join and no matches, skip this left row
		if len(matches) == 0 && j.joinType == ast.JoinTypeInner {
			j.currentLeft = nil
			continue
		}

		// Continue loop to emit first joined row
	}
}

// joinRows combines a left row and a right row into a single row
// If rightRow is nil (left join with no match), right fields are omitted
func (j *joinOperator) joinRows(leftRow *Row, rightRow *Row) *Row {
	// Start with all fields from left row
	joined := leftRow.ToMap()

	// Add fields from right row (if present)
	if rightRow != nil {
		for _, fieldName := range rightRow.Fields() {
			// Skip the join field from right side (already in left)
			if fieldName == j.rightField {
				continue
			}

			value, _ := rightRow.Get(fieldName)

			// Check for field name conflicts
			if _, exists := joined[fieldName]; exists {
				// Conflict: add "_right" suffix
				joined[fieldName+"_right"] = value
			} else {
				joined[fieldName] = value
			}
		}
	}

	return NewRow(joined)
}

// Close releases resources
func (j *joinOperator) Close() error {
	j.closed = true

	// Close both sides
	var leftErr, rightErr error
	if j.input != nil {
		leftErr = j.input.Close()
	}
	if j.right != nil {
		rightErr = j.right.Close()
	}

	// Clear hash table
	j.hashTable = nil

	if leftErr != nil {
		return leftErr
	}
	return rightErr
}

// Stats returns execution statistics
func (j *joinOperator) Stats() *IteratorStats {
	return j.stats
}
