// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package analyzer

import (
	"fmt"
)

// Symbol represents a named entity in the scope (field, alias, etc.)
type Symbol struct {
	Name      string
	Type      FieldType
	Alias     string // If this symbol is an alias for another field
	SourceCol string // Original field name if aliased
}

// Scope manages symbols (fields, aliases) in a lexical scope
type Scope struct {
	parent  *Scope
	symbols map[string]*Symbol
}

// NewScope creates a new scope
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		symbols: make(map[string]*Symbol),
	}
}

// Define adds a symbol to the current scope
func (s *Scope) Define(name string, fieldType FieldType) error {
	if _, exists := s.symbols[name]; exists {
		return fmt.Errorf("symbol %s already defined in current scope", name)
	}

	s.symbols[name] = &Symbol{
		Name: name,
		Type: fieldType,
	}

	return nil
}

// DefineAlias adds an aliased symbol to the current scope
func (s *Scope) DefineAlias(alias string, sourceCol string, fieldType FieldType) error {
	if _, exists := s.symbols[alias]; exists {
		return fmt.Errorf("symbol %s already defined in current scope", alias)
	}

	s.symbols[alias] = &Symbol{
		Name:      alias,
		Type:      fieldType,
		Alias:     alias,
		SourceCol: sourceCol,
	}

	return nil
}

// Resolve looks up a symbol by name, searching parent scopes if needed
func (s *Scope) Resolve(name string) (*Symbol, error) {
	// Look in current scope
	if symbol, ok := s.symbols[name]; ok {
		return symbol, nil
	}

	// Look in parent scope
	if s.parent != nil {
		return s.parent.Resolve(name)
	}

	return nil, fmt.Errorf("symbol %s not found", name)
}

// Has checks if a symbol exists in this scope or any parent scope
func (s *Scope) Has(name string) bool {
	_, err := s.Resolve(name)
	return err == nil
}

// Lookup returns the type of a symbol if it exists, or nil if not found
func (s *Scope) Lookup(name string) *FieldType {
	symbol, err := s.Resolve(name)
	if err != nil {
		return nil
	}
	return &symbol.Type
}

// Update updates an existing symbol's type
func (s *Scope) Update(name string, fieldType FieldType) {
	// Update in current scope if it exists
	if symbol, ok := s.symbols[name]; ok {
		symbol.Type = fieldType
		return
	}

	// If not in current scope, add it
	s.symbols[name] = &Symbol{
		Name: name,
		Type: fieldType,
	}
}

// GetType returns the type of a symbol
func (s *Scope) GetType(name string) (FieldType, error) {
	symbol, err := s.Resolve(name)
	if err != nil {
		return FieldTypeUnknown, err
	}
	return symbol.Type, nil
}

// ResolveFieldName resolves an alias to its source column name
// If the name is not an alias, returns the name itself
func (s *Scope) ResolveFieldName(name string) string {
	symbol, err := s.Resolve(name)
	if err != nil {
		return name // Not found, return as-is
	}

	if symbol.SourceCol != "" {
		return symbol.SourceCol
	}

	return symbol.Name
}

// Parent returns the parent scope
func (s *Scope) Parent() *Scope {
	return s.parent
}

// Symbols returns all symbols defined in this scope (not including parents)
func (s *Scope) Symbols() map[string]*Symbol {
	return s.symbols
}

// AllSymbols returns all symbols including parent scopes
func (s *Scope) AllSymbols() map[string]*Symbol {
	all := make(map[string]*Symbol)

	// Walk up the scope chain
	current := s
	for current != nil {
		for name, symbol := range current.symbols {
			// Only add if not already present (inner scopes shadow outer scopes)
			if _, exists := all[name]; !exists {
				all[name] = symbol
			}
		}
		current = current.parent
	}

	return all
}

// Clone creates a shallow copy of the scope (same parent, new symbol map)
func (s *Scope) Clone() *Scope {
	cloned := NewScope(s.parent)
	for name, symbol := range s.symbols {
		cloned.symbols[name] = symbol
	}
	return cloned
}
