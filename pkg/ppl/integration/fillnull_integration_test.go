// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package integration

import (
	"strings"
	"testing"

	"github.com/conjugate/conjugate/pkg/ppl/analyzer"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFillnullCommand_Integration(t *testing.T) {
	p := newTestPipeline()

	// Add test-specific fields
	p.schema.AddField("cpu", analyzer.FieldTypeInt)
	p.schema.AddField("memory", analyzer.FieldTypeInt)
	p.schema.AddField("disk", analyzer.FieldTypeInt)
	p.schema.AddField("active", analyzer.FieldTypeBool)

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "FillAllNullFieldsWithDefaultValue",
			query: "source=logs | fillnull value=0",
		},
		{
			name:  "FillSpecificFieldsWithDefaultValue",
			query: "source=logs | fillnull value=0 fields cpu",
		},
		{
			name:  "FillNullWithStringValue",
			query: "source=logs | fillnull value=\"N/A\" fields status",
		},
		{
			name:  "FillNullMultipleFields",
			query: "source=logs | fillnull value=-1 fields cpu, memory, disk",
		},
		{
			name:  "FillNullWithBooleanValue",
			query: "source=logs | fillnull value=false fields active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse and plan
			physicalPlan, err := p.parseAndPlan(tt.query)
			require.NoError(t, err)
			require.NotNil(t, physicalPlan)

			// Verify plan contains fillnull operator
			planStr := physicalPlan.String()
			assert.Contains(t, planStr, "PhysicalFillnull", "Plan should contain PhysicalFillnull operator")
		})
	}
}

// planTreeContains checks if any node in the physical plan tree contains the given substring
func planTreeContains(plan physical.PhysicalPlan, substr string) bool {
	if plan == nil {
		return false
	}
	if strings.Contains(plan.String(), substr) {
		return true
	}
	for _, child := range plan.Children() {
		if planTreeContains(child, substr) {
			return true
		}
	}
	return false
}

func TestFillnullCommand_WithFilter(t *testing.T) {
	p := newTestPipeline()
	p.schema.AddField("cpu", analyzer.FieldTypeInt)
	p.schema.AddField("memory", analyzer.FieldTypeInt)

	// Fillnull after filter
	query := "source=logs | where host=\"web1\" | fillnull value=0 fields cpu"

	// Parse and plan
	physicalPlan, err := p.parseAndPlan(query)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Verify plan tree contains fillnull (may not be the outermost node due to optimization)
	assert.True(t, planTreeContains(physicalPlan, "PhysicalFillnull"), "Plan tree should contain PhysicalFillnull")
}

func TestFillnullCommand_WithFields(t *testing.T) {
	p := newTestPipeline()
	p.schema.AddField("cpu", analyzer.FieldTypeInt)
	p.schema.AddField("memory", analyzer.FieldTypeInt)

	// Fillnull and then project specific fields
	query := "source=logs | fillnull value=0 fields cpu, memory | fields host, cpu"

	// Parse and plan
	physicalPlan, err := p.parseAndPlan(query)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Verify plan contains fillnull
	planStr := physicalPlan.String()
	assert.Contains(t, planStr, "PhysicalFillnull", "Plan should contain PhysicalFillnull")
}

func TestFillnullCommand_WithStats(t *testing.T) {
	p := newTestPipeline()
	p.schema.AddField("cpu", analyzer.FieldTypeInt)

	// Fillnull before stats
	query := "source=logs | fillnull value=0 fields cpu | stats avg(cpu) as avg_cpu by host"

	// Parse and plan
	physicalPlan, err := p.parseAndPlan(query)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// The plan should contain both fillnull and aggregation
	planStr := physicalPlan.String()
	assert.Contains(t, planStr, "PhysicalFillnull", "Plan should contain PhysicalFillnull")
}

func TestFillnullCommand_InPipeline(t *testing.T) {
	p := newTestPipeline()
	p.schema.AddField("cpu", analyzer.FieldTypeInt)
	p.schema.AddField("memory", analyzer.FieldTypeInt)

	// Fillnull in middle of pipeline
	query := "source=logs | where status=200 | fillnull value=0 fields cpu | sort cpu | head 10"

	// Parse and plan
	physicalPlan, err := p.parseAndPlan(query)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// The plan should succeed
	assert.NotNil(t, physicalPlan)
}

func TestFillnullCommand_PerFieldAssignments(t *testing.T) {
	p := newTestPipeline()

	// Add test-specific fields
	p.schema.AddField("cpu", analyzer.FieldTypeInt)
	p.schema.AddField("memory", analyzer.FieldTypeInt)
	p.schema.AddField("status", analyzer.FieldTypeString)

	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "SingleFieldAssignment",
			query: "source=logs | fillnull cpu=0",
		},
		{
			name:  "MultipleFieldAssignments",
			query: "source=logs | fillnull cpu=0, memory=0",
		},
		{
			name:  "MixedTypeAssignments",
			query: `source=logs | fillnull cpu=0, memory=0, status="unknown"`,
		},
		{
			name:  "AssignmentsWithDifferentValues",
			query: "source=logs | fillnull cpu=0, memory=100, status=\"N/A\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse and plan
			physicalPlan, err := p.parseAndPlan(tt.query)
			require.NoError(t, err, "Query: %s", tt.query)
			require.NotNil(t, physicalPlan, "Query: %s", tt.query)

			// Verify plan contains fillnull operator
			planStr := physicalPlan.String()
			assert.Contains(t, planStr, "PhysicalFillnull", "Plan should contain PhysicalFillnull operator for query: %s", tt.query)
		})
	}
}

func TestFillnullCommand_AssignmentsInPipeline(t *testing.T) {
	p := newTestPipeline()
	p.schema.AddField("cpu", analyzer.FieldTypeInt)
	p.schema.AddField("memory", analyzer.FieldTypeInt)

	// Fillnull with assignments in pipeline
	query := `source=logs | where host="web1" | fillnull cpu=0, memory=0 | sort cpu desc | head 10`

	// Parse and plan
	physicalPlan, err := p.parseAndPlan(query)
	require.NoError(t, err)
	require.NotNil(t, physicalPlan)

	// Verify plan tree contains fillnull (may not be the outermost node due to optimization)
	assert.True(t, planTreeContains(physicalPlan, "PhysicalFillnull"), "Plan tree should contain PhysicalFillnull")
}
