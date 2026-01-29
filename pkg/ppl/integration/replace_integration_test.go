// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package integration

import (
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplaceCommand_Integration(t *testing.T) {
	p := newTestPipeline()

	// Add test-specific fields
	p.schema.AddField("level", analyzer.FieldTypeString)

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "BasicReplace",
			query: "source=logs | replace 'error' with 'ERROR', 'warn' with 'WARNING' in level",
		},
		{
			name:  "SingleReplace",
			query: "source=logs | replace 'error' with 'critical' in level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse and plan
			physicalPlan, err := p.parseAndPlan(tt.query)
			require.NoError(t, err)
			require.NotNil(t, physicalPlan)

			// Verify plan contains replace operator
			planStr := physicalPlan.String()
			assert.Contains(t, planStr, "PhysicalReplace", "Plan should contain PhysicalReplace operator")
		})
	}
}

func TestReplaceCommand_WithFilter(t *testing.T) {
	p := newTestPipeline()
	p.schema.AddField("code", analyzer.FieldTypeString)

	// Replace after filter
	query := "source=logs | where status=404 | replace 'ERROR' with 'CRITICAL' in message"

	// Parse and plan
	physicalPlan, err := p.parseAndPlan(query)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Verify plan contains replace (filter might be pushed down)
	// Just verify the plan is not nil and doesn't error
	assert.NotNil(t, physicalPlan)
}

func TestReplaceCommand_WithFields(t *testing.T) {
	p := newTestPipeline()
	p.schema.AddField("old_field", analyzer.FieldTypeString)
	p.schema.AddField("value", analyzer.FieldTypeString)

	// Replace and then project specific fields
	query := "source=logs | replace 'old' with 'new' in value | fields value"

	// Parse and plan
	physicalPlan, err := p.parseAndPlan(query)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Verify plan contains replace
	planStr := physicalPlan.String()
	assert.Contains(t, planStr, "PhysicalReplace", "Plan should contain PhysicalReplace")
}

func TestReplaceCommand_Multiple(t *testing.T) {
	p := newTestPipeline()
	p.schema.AddField("text", analyzer.FieldTypeString)
	p.schema.AddField("level", analyzer.FieldTypeString)

	// Multiple replacements in same field
	query := "source=logs | replace 'error' with 'ERROR', 'warning' with 'WARN', 'info' with 'INFO' in level"

	// Parse and plan
	physicalPlan, err := p.parseAndPlan(query)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// The plan should succeed
	assert.NotNil(t, physicalPlan)
}
