// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================================
// Schema Extended Method Tests (matching actual API)
// =====================================================================

func TestSchema_AddObjectField(t *testing.T) {
	schema := NewSchema("test")

	t.Run("AddSimpleObjectField", func(t *testing.T) {
		schema.AddObjectField("config", map[string]*Field{
			"enabled": {Name: "enabled", Type: FieldTypeBool},
			"timeout": {Name: "timeout", Type: FieldTypeInt},
		})

		field, err := schema.GetField("config")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeObject, field.Type)
	})

	t.Run("AddEmptyObjectField", func(t *testing.T) {
		schema.AddObjectField("empty", map[string]*Field{})

		field, err := schema.GetField("empty")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeObject, field.Type)
	})
}

func TestSchema_AddArrayField(t *testing.T) {
	schema := NewSchema("test")

	t.Run("AddSimpleArrayField", func(t *testing.T) {
		schema.AddArrayField("tags", FieldTypeString)

		field, err := schema.GetField("tags")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, field.Type)
		assert.True(t, field.Array)
	})

	t.Run("AddObjectArrayField", func(t *testing.T) {
		schema.AddArrayField("items", FieldTypeObject)

		field, err := schema.GetField("items")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeObject, field.Type)
		assert.True(t, field.Array)
	})
}

func TestSchema_Merge(t *testing.T) {
	t.Run("MergeDisjointSchemas", func(t *testing.T) {
		schema1 := NewSchema("logs")
		schema1.AddField("field1", FieldTypeString)

		schema2 := NewSchema("logs")
		schema2.AddField("field2", FieldTypeInt)

		merged := schema1.Merge(schema2)

		assert.True(t, merged.HasField("field1"))
		assert.True(t, merged.HasField("field2"))
	})

	t.Run("MergeOverlappingSchemas", func(t *testing.T) {
		schema1 := NewSchema("logs")
		schema1.AddField("shared", FieldTypeString)
		schema1.AddField("unique1", FieldTypeInt)

		schema2 := NewSchema("logs")
		schema2.AddField("shared", FieldTypeString)
		schema2.AddField("unique2", FieldTypeBool)

		merged := schema1.Merge(schema2)

		assert.True(t, merged.HasField("shared"))
		assert.True(t, merged.HasField("unique1"))
		assert.True(t, merged.HasField("unique2"))
		assert.Len(t, merged.Fields, 3)
	})

	t.Run("MergeEmptySchema", func(t *testing.T) {
		schema1 := NewSchema("logs")
		schema1.AddField("field1", FieldTypeString)

		schema2 := NewSchema("logs")

		merged := schema1.Merge(schema2)

		assert.True(t, merged.HasField("field1"))
		assert.Len(t, merged.Fields, 1)
	})
}

func TestSchema_Project(t *testing.T) {
	t.Run("ProjectSubsetOfFields", func(t *testing.T) {
		schema := NewSchema("logs")
		schema.AddField("field1", FieldTypeString)
		schema.AddField("field2", FieldTypeInt)
		schema.AddField("field3", FieldTypeBool)

		projected, err := schema.Project([]string{"field1", "field3"})
		require.NoError(t, err)

		assert.True(t, projected.HasField("field1"))
		assert.False(t, projected.HasField("field2"))
		assert.True(t, projected.HasField("field3"))
		assert.Len(t, projected.Fields, 2)
	})

	t.Run("ProjectNonexistentFields", func(t *testing.T) {
		schema := NewSchema("logs")
		schema.AddField("field1", FieldTypeString)

		_, err := schema.Project([]string{"field1", "nonexistent"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("ProjectEmptyFieldList", func(t *testing.T) {
		schema := NewSchema("logs")
		schema.AddField("field1", FieldTypeString)

		projected, err := schema.Project([]string{})
		require.NoError(t, err)
		assert.Len(t, projected.Fields, 0)
	})

	t.Run("ProjectAllFields", func(t *testing.T) {
		schema := NewSchema("logs")
		schema.AddField("field1", FieldTypeString)
		schema.AddField("field2", FieldTypeInt)

		projected, err := schema.Project([]string{"field1", "field2"})
		require.NoError(t, err)
		assert.Len(t, projected.Fields, 2)
	})
}

func TestSchema_Clone(t *testing.T) {
	t.Run("CloneSimpleSchema", func(t *testing.T) {
		schema := NewSchema("logs")
		schema.AddField("field1", FieldTypeString)
		schema.AddField("field2", FieldTypeInt)

		cloned := schema.Clone()

		assert.Equal(t, schema.Source, cloned.Source)
		assert.Len(t, cloned.Fields, 2)
		assert.True(t, cloned.HasField("field1"))
		assert.True(t, cloned.HasField("field2"))

		// Modify clone shouldn't affect original
		cloned.AddField("field3", FieldTypeBool)
		assert.False(t, schema.HasField("field3"))
		assert.True(t, cloned.HasField("field3"))
	})

	t.Run("CloneObjectField", func(t *testing.T) {
		schema := NewSchema("logs")
		schema.AddObjectField("metadata", map[string]*Field{
			"category": {Name: "category", Type: FieldTypeString},
		})

		cloned := schema.Clone()

		field, err := cloned.GetField("metadata")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeObject, field.Type)
	})

	t.Run("CloneArrayField", func(t *testing.T) {
		schema := NewSchema("logs")
		schema.AddArrayField("tags", FieldTypeString)

		cloned := schema.Clone()

		field, err := cloned.GetField("tags")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, field.Type)
		assert.True(t, field.Array)
	})
}

func TestSchema_String(t *testing.T) {
	schema := NewSchema("logs")
	schema.AddField("field1", FieldTypeString)
	schema.AddField("field2", FieldTypeInt)

	str := schema.String()
	assert.Contains(t, str, "logs")
	assert.Contains(t, str, "field1")
	assert.Contains(t, str, "field2")
}

func TestSchema_FieldType(t *testing.T) {
	schema := NewSchema("test")
	schema.AddField("status", FieldTypeInt)

	t.Run("GetExistingFieldType", func(t *testing.T) {
		fieldType, err := schema.FieldType("status")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("GetNonexistentFieldType", func(t *testing.T) {
		_, err := schema.FieldType("nonexistent")
		assert.Error(t, err)
	})
}

func TestSchema_ComplexOperations(t *testing.T) {
	schema := NewSchema("logs")

	// Add various field types
	schema.AddField("simple", FieldTypeString)
	schema.AddObjectField("nested", map[string]*Field{
		"inner": {Name: "inner", Type: FieldTypeInt},
	})
	schema.AddArrayField("list", FieldTypeString)

	// Project
	projected, err := schema.Project([]string{"simple", "nested"})
	require.NoError(t, err)
	assert.True(t, projected.HasField("simple"))
	assert.True(t, projected.HasField("nested"))
	assert.False(t, projected.HasField("list"))

	// Clone
	cloned := projected.Clone()
	assert.Equal(t, projected.Source, cloned.Source)
	assert.Len(t, cloned.Fields, 2)

	// Merge
	schema2 := NewSchema("logs")
	schema2.AddField("additional", FieldTypeBool)

	merged := cloned.Merge(schema2)
	assert.True(t, merged.HasField("simple"))
	assert.True(t, merged.HasField("nested"))
	assert.True(t, merged.HasField("additional"))
}

func TestField_Nullable(t *testing.T) {
	field := &Field{
		Name:     "test",
		Type:     FieldTypeString,
		Nullable: true,
	}

	assert.True(t, field.Nullable)
}
