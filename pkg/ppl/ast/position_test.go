// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition_String(t *testing.T) {
	tests := []struct {
		name     string
		pos      Position
		expected string
	}{
		{
			name:     "valid position",
			pos:      Position{Line: 1, Column: 5, Offset: 4},
			expected: "line 1, column 5",
		},
		{
			name:     "zero position",
			pos:      Position{Line: 0, Column: 0, Offset: 0},
			expected: "line 0, column 0",
		},
		{
			name:     "large line number",
			pos:      Position{Line: 10000, Column: 500, Offset: 500000},
			expected: "line 10000, column 500",
		},
		{
			name:     "negative values (invalid but should handle)",
			pos:      Position{Line: -1, Column: -1, Offset: -1},
			expected: "line -1, column -1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPosition_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		pos      Position
		expected bool
	}{
		{
			name:     "valid position",
			pos:      Position{Line: 1, Column: 1, Offset: 0},
			expected: true,
		},
		{
			name:     "valid with offset",
			pos:      Position{Line: 5, Column: 10, Offset: 100},
			expected: true,
		},
		{
			name:     "zero line (invalid)",
			pos:      Position{Line: 0, Column: 1, Offset: 0},
			expected: false,
		},
		{
			name:     "zero column (invalid)",
			pos:      Position{Line: 1, Column: 0, Offset: 0},
			expected: false,
		},
		{
			name:     "negative offset (invalid)",
			pos:      Position{Line: 1, Column: 1, Offset: -1},
			expected: false,
		},
		{
			name:     "all zeros (invalid)",
			pos:      Position{Line: 0, Column: 0, Offset: 0},
			expected: false,
		},
		{
			name:     "negative line (invalid)",
			pos:      Position{Line: -1, Column: 1, Offset: 0},
			expected: false,
		},
		{
			name:     "negative column (invalid)",
			pos:      Position{Line: 1, Column: -1, Offset: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPosition_Before(t *testing.T) {
	tests := []struct {
		name     string
		pos1     Position
		pos2     Position
		expected bool
	}{
		{
			name:     "earlier line",
			pos1:     Position{Line: 1, Column: 5, Offset: 4},
			pos2:     Position{Line: 2, Column: 1, Offset: 10},
			expected: true,
		},
		{
			name:     "same line, earlier column",
			pos1:     Position{Line: 1, Column: 5, Offset: 4},
			pos2:     Position{Line: 1, Column: 10, Offset: 9},
			expected: true,
		},
		{
			name:     "later line",
			pos1:     Position{Line: 2, Column: 1, Offset: 10},
			pos2:     Position{Line: 1, Column: 5, Offset: 4},
			expected: false,
		},
		{
			name:     "same line, later column",
			pos1:     Position{Line: 1, Column: 10, Offset: 9},
			pos2:     Position{Line: 1, Column: 5, Offset: 4},
			expected: false,
		},
		{
			name:     "identical positions",
			pos1:     Position{Line: 1, Column: 5, Offset: 4},
			pos2:     Position{Line: 1, Column: 5, Offset: 4},
			expected: false,
		},
		{
			name:     "zero positions",
			pos1:     Position{Line: 0, Column: 0, Offset: 0},
			pos2:     Position{Line: 0, Column: 0, Offset: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos1.Before(tt.pos2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPosition_EdgeCases(t *testing.T) {
	t.Run("max int32 values", func(t *testing.T) {
		pos := Position{
			Line:   2147483647,
			Column: 2147483647,
			Offset: 2147483647,
		}
		assert.True(t, pos.IsValid())
		assert.Equal(t, "line 2147483647, column 2147483647", pos.String())
	})

	t.Run("boundary at line 1 column 1", func(t *testing.T) {
		pos := Position{Line: 1, Column: 1, Offset: 0}
		assert.True(t, pos.IsValid())
	})

	t.Run("just below valid threshold", func(t *testing.T) {
		pos1 := Position{Line: 0, Column: 1, Offset: 0}
		assert.False(t, pos1.IsValid())

		pos2 := Position{Line: 1, Column: 0, Offset: 0}
		assert.False(t, pos2.IsValid())
	})
}
