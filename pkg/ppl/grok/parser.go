// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package grok

import (
	"fmt"
	"regexp"
	"strings"
)

// Grok represents a compiled grok pattern
type Grok struct {
	pattern      string                 // Original grok pattern
	regexp       *regexp.Regexp         // Compiled regex
	fieldNames   []string               // Ordered list of captured field names
	fieldTypes   map[string]string      // Field name -> type (int, float, string)
	subPatterns  map[string]string      // Custom sub-patterns
}

// Field represents a captured field from a grok match
type Field struct {
	Name  string
	Value string
	Type  string // "string", "int", "float"
}

// Match represents a successful grok match
type Match struct {
	Fields map[string]interface{} // Field name -> typed value
}

// grokPattern matches %{PATTERN:field:type} syntax
var grokPattern = regexp.MustCompile(`%{([A-Za-z0-9_]+)(?::([A-Za-z0-9_]+))?(?::([A-Za-z]+))?}`)

// NewGrok creates a new Grok parser
func NewGrok(pattern string) (*Grok, error) {
	g := &Grok{
		pattern:     pattern,
		fieldNames:  make([]string, 0),
		fieldTypes:  make(map[string]string),
		subPatterns: make(map[string]string),
	}

	// Compile the pattern
	if err := g.compile(); err != nil {
		return nil, err
	}

	return g, nil
}

// NewGrokWithPatterns creates a Grok parser with custom patterns
func NewGrokWithPatterns(pattern string, customPatterns map[string]string) (*Grok, error) {
	g := &Grok{
		pattern:     pattern,
		fieldNames:  make([]string, 0),
		fieldTypes:  make(map[string]string),
		subPatterns: customPatterns,
	}

	if err := g.compile(); err != nil {
		return nil, err
	}

	return g, nil
}

// compile converts the grok pattern to a regular expression
func (g *Grok) compile() error {
	expandedPattern := g.pattern

	// Track seen fields to avoid duplicates
	seenFields := make(map[string]bool)

	// Iteratively expand all grok patterns
	maxIterations := 100 // Prevent infinite loops
	for i := 0; i < maxIterations; i++ {
		matches := grokPattern.FindAllStringSubmatch(expandedPattern, -1)
		if len(matches) == 0 {
			break // No more patterns to expand
		}

		replaced := false
		for _, match := range matches {
			fullMatch := match[0]     // %{PATTERN:field:type}
			patternName := match[1]   // PATTERN
			fieldName := ""
			fieldType := "string"     // default type

			if len(match) > 2 && match[2] != "" {
				fieldName = match[2]  // field
			}
			if len(match) > 3 && match[3] != "" {
				fieldType = match[3]  // type
			}

			// Look up the pattern (first in custom, then in built-in)
			var patternRegex string
			var ok bool

			if g.subPatterns != nil {
				patternRegex, ok = g.subPatterns[patternName]
			}
			if !ok {
				patternRegex, ok = AllPatterns[patternName]
			}
			if !ok {
				return fmt.Errorf("unknown pattern: %s", patternName)
			}

			// If field name is specified, create a named capture group
			var replacement string
			if fieldName != "" {
				// Store field info
				if !seenFields[fieldName] {
					g.fieldNames = append(g.fieldNames, fieldName)
					g.fieldTypes[fieldName] = fieldType
					seenFields[fieldName] = true
				}

				// Wrap in named capture group
				replacement = fmt.Sprintf("(?P<%s>%s)", fieldName, patternRegex)
			} else {
				// No field name, just expand the pattern without capturing
				replacement = fmt.Sprintf("(?:%s)", patternRegex)
			}

			// Replace first occurrence
			expandedPattern = strings.Replace(expandedPattern, fullMatch, replacement, 1)
			replaced = true
		}

		if !replaced {
			break
		}
	}

	// Compile the final regex
	re, err := regexp.Compile(expandedPattern)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w (pattern: %s)", err, expandedPattern)
	}

	g.regexp = re
	return nil
}

// Match attempts to match the pattern against the input text
func (g *Grok) Match(text string) (*Match, bool) {
	matches := g.regexp.FindStringSubmatch(text)
	if matches == nil {
		return nil, false
	}

	result := &Match{
		Fields: make(map[string]interface{}),
	}

	// Extract named groups
	for i, name := range g.regexp.SubexpNames() {
		if i == 0 || name == "" {
			continue // Skip the full match and unnamed groups
		}

		if i >= len(matches) {
			continue
		}

		value := matches[i]

		// Convert to appropriate type
		fieldType := g.fieldTypes[name]
		typedValue := convertType(value, fieldType)

		result.Fields[name] = typedValue
	}

	return result, true
}

// ParseAll finds all matches in the text
func (g *Grok) ParseAll(text string) []*Match {
	allMatches := g.regexp.FindAllStringSubmatch(text, -1)
	if allMatches == nil {
		return nil
	}

	results := make([]*Match, 0, len(allMatches))

	for _, matches := range allMatches {
		result := &Match{
			Fields: make(map[string]interface{}),
		}

		for i, name := range g.regexp.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}

			if i >= len(matches) {
				continue
			}

			value := matches[i]
			fieldType := g.fieldTypes[name]
			typedValue := convertType(value, fieldType)

			result.Fields[name] = typedValue
		}

		results = append(results, result)
	}

	return results
}

// convertType converts a string value to the specified type
func convertType(value, typeName string) interface{} {
	switch strings.ToLower(typeName) {
	case "int", "integer":
		// Try to parse as int64
		var intVal int64
		fmt.Sscanf(value, "%d", &intVal)
		return intVal

	case "float", "double", "number":
		// Try to parse as float64
		var floatVal float64
		fmt.Sscanf(value, "%f", &floatVal)
		return floatVal

	case "bool", "boolean":
		// Parse boolean
		switch strings.ToLower(value) {
		case "true", "yes", "1", "on":
			return true
		case "false", "no", "0", "off":
			return false
		}
		return value

	default:
		// Return as string
		return value
	}
}

// Pattern returns the original grok pattern
func (g *Grok) Pattern() string {
	return g.pattern
}

// FieldNames returns the list of captured field names
func (g *Grok) FieldNames() []string {
	return g.fieldNames
}

// FieldTypes returns the map of field types
func (g *Grok) FieldTypes() map[string]string {
	return g.fieldTypes
}
