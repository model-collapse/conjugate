// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

// spathOperator extracts and navigates JSON data from fields
// Supports JSONPath syntax for extracting nested values
//
// Examples:
//   spath path="response.user.id" output=user_id
//   spath input=data path="$.items[*].name" output=item_names
//   spath (auto-extract all JSON fields from _raw)
type spathOperator struct {
	input      Operator
	logger     *zap.Logger
	ctx        context.Context
	stats      *IteratorStats
	opened     bool
	closed     bool
	inputField string // Field to extract from (default: "_raw")
	path       string // JSONPath expression (empty = auto-extract)
	outputField string // Field to write to (default: last component of path)
	autoExtract bool   // If true, extract all JSON fields
}

// SpathConfig holds configuration for spath operator
type SpathConfig struct {
	InputField  string // Field containing JSON (default: "_raw")
	Path        string // JSONPath expression (empty = auto-extract all)
	OutputField string // Field to store result (default: derived from path)
}

// NewSpathOperator creates a new spath operator
func NewSpathOperator(input Operator, config SpathConfig, logger *zap.Logger) *spathOperator {
	// Set defaults
	if config.InputField == "" {
		config.InputField = "_raw"
	}

	autoExtract := config.Path == ""

	// If no output field specified and path given, derive from path
	if config.OutputField == "" && config.Path != "" {
		config.OutputField = deriveFieldName(config.Path)
	}

	return &spathOperator{
		input:       input,
		logger:      logger,
		stats:       &IteratorStats{},
		inputField:  config.InputField,
		path:        config.Path,
		outputField: config.OutputField,
		autoExtract: autoExtract,
	}
}

// deriveFieldName extracts a field name from a JSONPath
// Examples:
//   "$.user.name" -> "name"
//   "response.data.id" -> "id"
//   "items[0].title" -> "title"
func deriveFieldName(path string) string {
	// Remove leading $. if present
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")

	// Split by . and take last component
	parts := strings.Split(path, ".")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		// Remove array indexing like [0] or [*]
		lastPart = strings.Split(lastPart, "[")[0]
		if lastPart != "" {
			return lastPart
		}
	}

	return "result"
}

// Open initializes the operator
func (s *spathOperator) Open(ctx context.Context) error {
	if s.opened {
		return nil
	}

	s.ctx = ctx
	s.logger.Debug("Opening spath operator",
		zap.String("input_field", s.inputField),
		zap.String("path", s.path),
		zap.String("output_field", s.outputField),
		zap.Bool("auto_extract", s.autoExtract))

	if err := s.input.Open(ctx); err != nil {
		return fmt.Errorf("failed to open input: %w", err)
	}

	s.opened = true
	return nil
}

// Next processes the next row, extracting JSON data
func (s *spathOperator) Next(ctx context.Context) (*Row, error) {
	if s.closed {
		return nil, ErrClosed
	}

	if !s.opened {
		return nil, ErrClosed
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Read next row
	row, err := s.input.Next(ctx)
	if err != nil {
		return nil, err
	}

	s.stats.RowsRead++

	// Get the input field value
	inputValue, exists := row.Get(s.inputField)
	if !exists {
		// Input field doesn't exist, return row as-is
		s.stats.RowsReturned++
		return row, nil
	}

	// Convert to JSON string
	jsonStr, err := s.toJSONString(inputValue)
	if err != nil {
		// Not valid JSON, return row as-is
		s.logger.Debug("Failed to convert to JSON",
			zap.String("field", s.inputField),
			zap.Error(err))
		s.stats.RowsReturned++
		return row, nil
	}

	// Process based on mode
	if s.autoExtract {
		// Auto-extract all JSON fields
		s.autoExtractFields(row, jsonStr)
	} else {
		// Extract specific path
		s.extractPath(row, jsonStr)
	}

	s.stats.RowsReturned++
	return row, nil
}

// toJSONString converts a value to JSON string
func (s *spathOperator) toJSONString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case map[string]interface{}:
		// Marshal to JSON
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	default:
		// Try to marshal
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}

// autoExtractFields extracts all fields from JSON and adds them to the row
func (s *spathOperator) autoExtractFields(row *Row, jsonStr string) {
	// Parse JSON and extract all fields at the top level
	result := gjson.Parse(jsonStr)

	if !result.IsObject() {
		// Not an object, can't auto-extract
		return
	}

	// Iterate over all fields
	result.ForEach(func(key, value gjson.Result) bool {
		fieldName := key.String()
		fieldValue := s.convertGjsonValue(value)
		row.Set(fieldName, fieldValue)
		return true // continue iteration
	})
}

// extractPath extracts a specific JSONPath and sets the output field
func (s *spathOperator) extractPath(row *Row, jsonStr string) {
	// Normalize path for gjson
	path := s.path
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")

	// Query using gjson
	result := gjson.Get(jsonStr, path)

	if !result.Exists() {
		// Path doesn't exist, don't set field
		return
	}

	// Convert result to appropriate Go type
	value := s.convertGjsonValue(result)

	// Set output field
	row.Set(s.outputField, value)
}

// convertGjsonValue converts a gjson.Result to a Go value
// According to OpenSearch PPL specification, spath returns all fields as STRING type
func (s *spathOperator) convertGjsonValue(result gjson.Result) interface{} {
	switch result.Type {
	case gjson.Null:
		return nil
	case gjson.False:
		// Return as string per OpenSearch spec
		return "false"
	case gjson.True:
		// Return as string per OpenSearch spec
		return "true"
	case gjson.Number:
		// Return as string per OpenSearch spec
		return result.String()
	case gjson.String:
		return result.String()
	case gjson.JSON:
		// Complex type (object or array)
		if result.IsArray() {
			// Return as slice
			var arr []interface{}
			result.ForEach(func(_, value gjson.Result) bool {
				arr = append(arr, s.convertGjsonValue(value))
				return true
			})
			return arr
		} else if result.IsObject() {
			// Return as map
			obj := make(map[string]interface{})
			result.ForEach(func(key, value gjson.Result) bool {
				obj[key.String()] = s.convertGjsonValue(value)
				return true
			})
			return obj
		}
		// Fallback: return raw string
		return result.String()
	default:
		return result.String()
	}
}

// Close releases resources
func (s *spathOperator) Close() error {
	s.closed = true
	if s.input != nil {
		return s.input.Close()
	}
	return nil
}

// Stats returns execution statistics
func (s *spathOperator) Stats() *IteratorStats {
	return s.stats
}
