// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =====================================================================
// Scope Extended Method Tests (for missing coverage)
// =====================================================================

func TestScope_DefineAlias(t *testing.T) {
	scope := NewScope(nil)
	scope.Define("original_field", FieldTypeString)

	t.Run("DefineValidAlias", func(t *testing.T) {
		err := scope.DefineAlias("alias_name", "original_field", FieldTypeString)
		require.NoError(t, err)

		// Alias should resolve to original field
		assert.True(t, scope.Has("alias_name"))
	})

	t.Run("DefineAliasConflict", func(t *testing.T) {
		scope.Define("existing", FieldTypeInt)
		err := scope.DefineAlias("existing", "original_field", FieldTypeString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already defined")
	})
}

func TestScope_Update(t *testing.T) {
	scope := NewScope(nil)

	t.Run("UpdateExistingSymbol", func(t *testing.T) {
		scope.Define("field1", FieldTypeInt)
		scope.Update("field1", FieldTypeString)

		fieldType, err := scope.GetType("field1")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, fieldType)
	})

	t.Run("UpdateNonexistentSymbol", func(t *testing.T) {
		// Update creates the symbol if it doesn't exist
		scope.Update("new_field", FieldTypeInt)
		// Field should now exist
		assert.True(t, scope.Has("new_field"))
		fieldType, err := scope.GetType("new_field")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, fieldType)
	})

	t.Run("UpdateWithParentScope", func(t *testing.T) {
		parent := NewScope(nil)
		parent.Define("parent_field", FieldTypeInt)

		child := NewScope(parent)
		child.Define("child_field", FieldTypeString)

		// Update child field
		child.Update("child_field", FieldTypeBool)
		fieldType, err := child.GetType("child_field")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeBool, fieldType)

		// Update doesn't modify parent scope fields
		// It creates a shadow in the child scope
		child.Update("parent_field", FieldTypeString)

		// Verify parent field was NOT changed
		parentFieldType, err := parent.GetType("parent_field")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeInt, parentFieldType)

		// But child scope sees the new value
		childFieldType, err := child.GetType("parent_field")
		require.NoError(t, err)
		assert.Equal(t, FieldTypeString, childFieldType)
	})
}

func TestScope_ResolveFieldName(t *testing.T) {
	scope := NewScope(nil)
	scope.Define("original_name", FieldTypeString)
	scope.DefineAlias("alias_name", "original_name", FieldTypeString)

	t.Run("ResolveRegularField", func(t *testing.T) {
		resolved := scope.ResolveFieldName("original_name")
		assert.Equal(t, "original_name", resolved)
	})

	t.Run("ResolveAlias", func(t *testing.T) {
		resolved := scope.ResolveFieldName("alias_name")
		assert.Equal(t, "original_name", resolved)
	})

	t.Run("ResolveNonexistent", func(t *testing.T) {
		resolved := scope.ResolveFieldName("nonexistent")
		assert.Equal(t, "nonexistent", resolved)
	})
}

func TestScope_Parent(t *testing.T) {
	t.Run("RootScopeHasNoParent", func(t *testing.T) {
		root := NewScope(nil)
		assert.Nil(t, root.Parent())
	})

	t.Run("ChildScopeHasParent", func(t *testing.T) {
		parent := NewScope(nil)
		child := NewScope(parent)
		assert.Equal(t, parent, child.Parent())
	})
}

func TestScope_Symbols(t *testing.T) {
	scope := NewScope(nil)
	scope.Define("field1", FieldTypeString)
	scope.Define("field2", FieldTypeInt)

	symbols := scope.Symbols()
	assert.Len(t, symbols, 2)
	assert.Contains(t, symbols, "field1")
	assert.Contains(t, symbols, "field2")
}

func TestScope_AllSymbols(t *testing.T) {
	t.Run("RootScopeAllSymbols", func(t *testing.T) {
		root := NewScope(nil)
		root.Define("field1", FieldTypeString)
		root.Define("field2", FieldTypeInt)

		allSymbols := root.AllSymbols()
		assert.Len(t, allSymbols, 2)
		assert.Contains(t, allSymbols, "field1")
		assert.Contains(t, allSymbols, "field2")
	})

	t.Run("NestedScopeAllSymbols", func(t *testing.T) {
		parent := NewScope(nil)
		parent.Define("parent_field", FieldTypeString)

		child := NewScope(parent)
		child.Define("child_field", FieldTypeInt)

		allSymbols := child.AllSymbols()
		assert.Len(t, allSymbols, 2)
		assert.Contains(t, allSymbols, "parent_field")
		assert.Contains(t, allSymbols, "child_field")
	})

	t.Run("MultiLevelScopeAllSymbols", func(t *testing.T) {
		root := NewScope(nil)
		root.Define("root_field", FieldTypeString)

		middle := NewScope(root)
		middle.Define("middle_field", FieldTypeInt)

		leaf := NewScope(middle)
		leaf.Define("leaf_field", FieldTypeBool)

		allSymbols := leaf.AllSymbols()
		assert.Len(t, allSymbols, 3)
		assert.Contains(t, allSymbols, "root_field")
		assert.Contains(t, allSymbols, "middle_field")
		assert.Contains(t, allSymbols, "leaf_field")
	})
}

func TestScope_Clone(t *testing.T) {
	t.Run("CloneSimpleScope", func(t *testing.T) {
		original := NewScope(nil)
		original.Define("field1", FieldTypeString)
		original.Define("field2", FieldTypeInt)

		cloned := original.Clone()

		// Should have same symbols
		assert.True(t, cloned.Has("field1"))
		assert.True(t, cloned.Has("field2"))

		// Modifying clone shouldn't affect original
		cloned.Define("field3", FieldTypeBool)
		assert.True(t, cloned.Has("field3"))
		assert.False(t, original.Has("field3"))
	})

	t.Run("CloneScopeWithParent", func(t *testing.T) {
		parent := NewScope(nil)
		parent.Define("parent_field", FieldTypeString)

		original := NewScope(parent)
		original.Define("child_field", FieldTypeInt)

		cloned := original.Clone()

		// Should have same parent
		assert.Equal(t, parent, cloned.Parent())

		// Should have child field
		assert.True(t, cloned.Has("child_field"))

		// Should access parent field
		assert.True(t, cloned.Has("parent_field"))
	})

	t.Run("CloneScopeWithAliases", func(t *testing.T) {
		original := NewScope(nil)
		original.Define("original_field", FieldTypeString)
		original.DefineAlias("alias_field", "original_field", FieldTypeString)

		cloned := original.Clone()

		// Alias should work in cloned scope
		assert.True(t, cloned.Has("alias_field"))
		resolved := cloned.ResolveFieldName("alias_field")
		assert.Equal(t, "original_field", resolved)
	})
}

func TestScope_ComplexOperations(t *testing.T) {
	// Test complex combination of scope operations
	root := NewScope(nil)
	root.Define("root_field", FieldTypeString)

	child := NewScope(root)
	child.Define("child_field", FieldTypeInt)
	child.DefineAlias("alias1", "root_field", FieldTypeString)

	grandchild := NewScope(child)
	grandchild.Define("grandchild_field", FieldTypeBool)
	grandchild.DefineAlias("alias2", "child_field", FieldTypeInt)

	// Test lookup through nested scopes
	assert.True(t, grandchild.Has("root_field"))
	assert.True(t, grandchild.Has("child_field"))
	assert.True(t, grandchild.Has("grandchild_field"))
	assert.True(t, grandchild.Has("alias1"))
	assert.True(t, grandchild.Has("alias2"))

	// Test AllSymbols
	allSymbols := grandchild.AllSymbols()
	assert.Len(t, allSymbols, 5)

	// Test Clone preserves hierarchy
	cloned := grandchild.Clone()
	assert.Equal(t, child, cloned.Parent())
	assert.True(t, cloned.Has("grandchild_field"))
	assert.True(t, cloned.Has("alias2"))
}
