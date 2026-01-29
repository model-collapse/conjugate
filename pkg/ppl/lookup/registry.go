// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package lookup

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// Registry manages all lookup tables
type Registry struct {
	tables map[string]*LookupTable
	mu     sync.RWMutex
	logger *zap.Logger
}

// NewRegistry creates a new lookup table registry
func NewRegistry(logger *zap.Logger) *Registry {
	return &Registry{
		tables: make(map[string]*LookupTable),
		logger: logger,
	}
}

// Register adds a new lookup table to the registry
func (r *Registry) Register(table *LookupTable) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tables[table.Name]; exists {
		return fmt.Errorf("lookup table '%s' already registered", table.Name)
	}

	r.tables[table.Name] = table
	r.logger.Info("Registered lookup table",
		zap.String("table", table.Name),
		zap.Int("rows", table.Size()))

	return nil
}

// Get retrieves a lookup table by name
func (r *Registry) Get(name string) (*LookupTable, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	table, exists := r.tables[name]
	if !exists {
		return nil, fmt.Errorf("lookup table '%s' not found", name)
	}

	return table, nil
}

// Unregister removes a lookup table from the registry
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tables[name]; !exists {
		return fmt.Errorf("lookup table '%s' not found", name)
	}

	delete(r.tables, name)
	r.logger.Info("Unregistered lookup table", zap.String("table", name))

	return nil
}

// List returns all registered lookup table names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.tables))
	for name := range r.tables {
		names = append(names, name)
	}

	return names
}

// Clear removes all lookup tables
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tables = make(map[string]*LookupTable)
	r.logger.Info("Cleared all lookup tables")
}

// LoadFromCSV creates and registers a lookup table from a CSV file
func (r *Registry) LoadFromCSV(name string, filepath string, keyField string) error {
	table := NewLookupTable(name, r.logger)

	if err := table.LoadFromCSV(filepath, keyField); err != nil {
		return fmt.Errorf("failed to load lookup table: %w", err)
	}

	return r.Register(table)
}
