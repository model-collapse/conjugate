// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package analyzer

import (
	"fmt"
	"strings"
)

// FieldType represents the data type of a field
type FieldType int

const (
	FieldTypeUnknown FieldType = iota
	FieldTypeBool
	FieldTypeInt
	FieldTypeLong
	FieldTypeFloat
	FieldTypeDouble
	FieldTypeString
	FieldTypeText
	FieldTypeKeyword
	FieldTypeDate
	FieldTypeObject
	FieldTypeArray
	FieldTypeGeoPoint
	FieldTypeIP
)

// String returns the string representation of FieldType
func (ft FieldType) String() string {
	switch ft {
	case FieldTypeBool:
		return "boolean"
	case FieldTypeInt:
		return "integer"
	case FieldTypeLong:
		return "long"
	case FieldTypeFloat:
		return "float"
	case FieldTypeDouble:
		return "double"
	case FieldTypeString:
		return "string"
	case FieldTypeText:
		return "text"
	case FieldTypeKeyword:
		return "keyword"
	case FieldTypeDate:
		return "date"
	case FieldTypeObject:
		return "object"
	case FieldTypeArray:
		return "array"
	case FieldTypeGeoPoint:
		return "geo_point"
	case FieldTypeIP:
		return "ip"
	default:
		return "unknown"
	}
}

// IsNumeric returns true if the field type is numeric
func (ft FieldType) IsNumeric() bool {
	switch ft {
	case FieldTypeInt, FieldTypeLong, FieldTypeFloat, FieldTypeDouble:
		return true
	default:
		return false
	}
}

// IsString returns true if the field type is string-like
func (ft FieldType) IsString() bool {
	switch ft {
	case FieldTypeString, FieldTypeText, FieldTypeKeyword:
		return true
	default:
		return false
	}
}

// IsComparable returns true if values of this type can be compared
func (ft FieldType) IsComparable() bool {
	switch ft {
	case FieldTypeBool, FieldTypeInt, FieldTypeLong, FieldTypeFloat,
		FieldTypeDouble, FieldTypeString, FieldTypeText, FieldTypeKeyword,
		FieldTypeDate:
		return true
	default:
		return false
	}
}

// Field represents a field in a schema
type Field struct {
	Name     string
	Type     FieldType
	Nullable bool
	Array    bool
	Fields   map[string]*Field // For nested objects
}

// GetNestedField retrieves a nested field using dot notation
// e.g., "metadata.category" returns the category field within metadata
func (f *Field) GetNestedField(path string) (*Field, error) {
	if f.Type != FieldTypeObject {
		return nil, fmt.Errorf("field %s is not an object type", f.Name)
	}

	parts := strings.Split(path, ".")
	current := f

	for _, part := range parts {
		if current.Fields == nil {
			return nil, fmt.Errorf("field %s has no nested fields", current.Name)
		}

		next, ok := current.Fields[part]
		if !ok {
			return nil, fmt.Errorf("nested field %s not found in %s", part, current.Name)
		}

		current = next
	}

	return current, nil
}

// Schema represents the structure of an index or data source
type Schema struct {
	Source string
	Fields map[string]*Field
}

// NewSchema creates a new schema
func NewSchema(source string) *Schema {
	return &Schema{
		Source: source,
		Fields: make(map[string]*Field),
	}
}

// AddField adds a field to the schema
func (s *Schema) AddField(name string, fieldType FieldType) {
	s.Fields[name] = &Field{
		Name:     name,
		Type:     fieldType,
		Nullable: true, // Default: fields are nullable
		Array:    false,
	}
}

// AddObjectField adds a nested object field
func (s *Schema) AddObjectField(name string, fields map[string]*Field) {
	s.Fields[name] = &Field{
		Name:     name,
		Type:     FieldTypeObject,
		Nullable: true,
		Fields:   fields,
	}
}

// AddArrayField adds an array field
func (s *Schema) AddArrayField(name string, elementType FieldType) {
	s.Fields[name] = &Field{
		Name:     name,
		Type:     elementType,
		Nullable: true,
		Array:    true,
	}
}

// GetField retrieves a field by name (supports dot notation for nested fields)
func (s *Schema) GetField(name string) (*Field, error) {
	// Handle nested field references (e.g., "user.address.city")
	if strings.Contains(name, ".") {
		parts := strings.SplitN(name, ".", 2)
		rootField, ok := s.Fields[parts[0]]
		if !ok {
			return nil, fmt.Errorf("field %s not found in schema", parts[0])
		}
		return rootField.GetNestedField(parts[1])
	}

	// Simple field reference
	field, ok := s.Fields[name]
	if !ok {
		return nil, fmt.Errorf("field %s not found in schema", name)
	}

	return field, nil
}

// HasField checks if a field exists in the schema
func (s *Schema) HasField(name string) bool {
	_, err := s.GetField(name)
	return err == nil
}

// FieldType gets the type of a field
func (s *Schema) FieldType(name string) (FieldType, error) {
	field, err := s.GetField(name)
	if err != nil {
		return FieldTypeUnknown, err
	}
	return field.Type, nil
}

// Merge combines two schemas (for joins, unions, etc.)
func (s *Schema) Merge(other *Schema) *Schema {
	merged := NewSchema(fmt.Sprintf("%s_merged_%s", s.Source, other.Source))

	// Copy fields from first schema
	for name, field := range s.Fields {
		merged.Fields[name] = field
	}

	// Copy fields from second schema (conflicts use second schema's type)
	for name, field := range other.Fields {
		merged.Fields[name] = field
	}

	return merged
}

// Project creates a new schema with only the specified fields
func (s *Schema) Project(fieldNames []string) (*Schema, error) {
	projected := NewSchema(s.Source)

	for _, name := range fieldNames {
		field, err := s.GetField(name)
		if err != nil {
			return nil, err
		}
		projected.Fields[name] = field
	}

	return projected, nil
}

// Clone creates a deep copy of the schema
func (s *Schema) Clone() *Schema {
	cloned := NewSchema(s.Source)

	for name, field := range s.Fields {
		cloned.Fields[name] = &Field{
			Name:     field.Name,
			Type:     field.Type,
			Nullable: field.Nullable,
			Array:    field.Array,
			Fields:   field.Fields, // Shallow copy of nested fields (OK for now)
		}
	}

	return cloned
}

// String returns a string representation of the schema
func (s *Schema) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Schema(%s):\n", s.Source))
	for name, field := range s.Fields {
		sb.WriteString(fmt.Sprintf("  %s: %s", name, field.Type))
		if field.Array {
			sb.WriteString("[]")
		}
		if field.Nullable {
			sb.WriteString("?")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
