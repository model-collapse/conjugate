// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package lookup

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
)

// LookupTable represents an in-memory lookup table with hash-based indexing
type LookupTable struct {
	Name   string
	Fields []string                       // Column names
	Index  map[string]map[string]interface{} // key -> {field -> value}
	mu     sync.RWMutex                   // Protect concurrent access
	logger *zap.Logger
}

// NewLookupTable creates a new empty lookup table
func NewLookupTable(name string, logger *zap.Logger) *LookupTable {
	return &LookupTable{
		Name:   name,
		Fields: []string{},
		Index:  make(map[string]map[string]interface{}),
		logger: logger,
	}
}

// LoadFromCSV loads lookup table data from a CSV file
func (t *LookupTable) LoadFromCSV(filepath string, keyField string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header row
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	t.Fields = header

	// Find key field index
	keyIdx := -1
	for i, field := range header {
		if field == keyField {
			keyIdx = i
			break
		}
	}

	if keyIdx == -1 {
		return fmt.Errorf("key field '%s' not found in CSV header", keyField)
	}

	// Read data rows
	rowCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV row: %w", err)
		}

		if len(record) != len(header) {
			t.logger.Warn("Skipping malformed CSV row",
				zap.Int("expected", len(header)),
				zap.Int("got", len(record)))
			continue
		}

		// Build row map
		rowData := make(map[string]interface{})
		for i, value := range record {
			rowData[header[i]] = value
		}

		// Index by key field
		key := record[keyIdx]
		t.Index[key] = rowData
		rowCount++
	}

	t.logger.Info("Loaded lookup table from CSV",
		zap.String("table", t.Name),
		zap.String("file", filepath),
		zap.Int("rows", rowCount),
		zap.Strings("fields", t.Fields))

	return nil
}

// AddRow adds a single row to the lookup table
func (t *LookupTable) AddRow(key string, data map[string]interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Index[key] = data
}

// Lookup retrieves a row by key
func (t *LookupTable) Lookup(key string) (map[string]interface{}, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	row, exists := t.Index[key]
	return row, exists
}

// GetFields returns the list of available fields
func (t *LookupTable) GetFields() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.Fields
}

// Size returns the number of rows in the lookup table
func (t *LookupTable) Size() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return len(t.Index)
}

// Clear removes all data from the lookup table
func (t *LookupTable) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Index = make(map[string]map[string]interface{})
	t.logger.Info("Cleared lookup table", zap.String("table", t.Name))
}
