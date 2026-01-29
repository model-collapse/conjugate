// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package ast

import "fmt"

// Position represents a location in the source query
type Position struct {
	Line   int // Line number (1-indexed)
	Column int // Column number (1-indexed)
	Offset int // Byte offset (0-indexed)
}

// NewPosition creates a new position
func NewPosition(line, column, offset int) Position {
	return Position{
		Line:   line,
		Column: column,
		Offset: offset,
	}
}

// String returns a string representation of the position
func (p Position) String() string {
	return fmt.Sprintf("line %d, column %d", p.Line, p.Column)
}

// IsValid returns true if the position is valid
func (p Position) IsValid() bool {
	return p.Line > 0 && p.Column > 0 && p.Offset >= 0
}

// Before returns true if this position comes before the other position
func (p Position) Before(other Position) bool {
	return p.Offset < other.Offset
}

// NoPos represents an invalid position
var NoPos = Position{Line: 0, Column: 0, Offset: -1}
