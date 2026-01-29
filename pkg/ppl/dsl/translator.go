// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package dsl

import (
	"fmt"

	"github.com/conjugate/conjugate/pkg/ppl/ast"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
)

// DSL represents an OpenSearch Query DSL document
type DSL struct {
	Query        map[string]interface{}   `json:"query,omitempty"`
	Source       interface{}              `json:"_source,omitempty"` // bool or []string
	Sort         []map[string]interface{} `json:"sort,omitempty"`
	Size         *int                     `json:"size,omitempty"`
	From         *int                     `json:"from,omitempty"`
	Aggregations map[string]interface{}   `json:"aggs,omitempty"`
	ScriptFields map[string]interface{}   `json:"script_fields,omitempty"`
}

// Translator converts physical plans to OpenSearch DSL
type Translator struct {
	queryBuilder    *QueryBuilder
	aggBuilder      *AggregationBuilder
	functionBuilder FunctionBuilder
}

// FunctionBuilder interface for WASM UDF building
type FunctionBuilder interface {
	CanBuildUDF(expr ast.Expression) bool
	BuildUDF(expr ast.Expression) (*UDFReference, error)
	BuildComputedField(expr ast.Expression, alias string) (*UDFReference, error)
	BuildAggregationUDF(agg *ast.Aggregation) (*UDFReference, error)
}

// UDFReference represents a WASM UDF reference
type UDFReference struct {
	Name          string
	Version       string
	Parameters    map[string]interface{}
	FieldBindings map[string]string
}

// NewTranslator creates a new DSL translator
func NewTranslator() *Translator {
	return &Translator{
		queryBuilder: NewQueryBuilder(),
		aggBuilder:   NewAggregationBuilder(),
	}
}

// WithFunctionBuilder sets the function builder for WASM UDF support
func (t *Translator) WithFunctionBuilder(builder FunctionBuilder) *Translator {
	t.functionBuilder = builder
	t.queryBuilder.functionBuilder = builder
	t.aggBuilder.functionBuilder = builder
	return t
}

// Translate converts a physical plan to OpenSearch DSL
func (t *Translator) Translate(plan physical.PhysicalPlan) (*DSL, error) {
	dsl := &DSL{
		Query: make(map[string]interface{}),
	}

	// Extract the scan node (should be at the leaf)
	scans := physical.GetLeafScans(plan)
	if len(scans) == 0 {
		return nil, fmt.Errorf("no scan operation found in physical plan")
	}
	if len(scans) > 1 {
		return nil, fmt.Errorf("multiple scans not yet supported")
	}
	scan := scans[0]

	// Build query from pushed-down filter
	if scan.Filter != nil {
		query, err := t.queryBuilder.BuildFilter(scan.Filter)
		if err != nil {
			return nil, fmt.Errorf("failed to build filter: %w", err)
		}
		dsl.Query = query
	} else {
		// No filter - match all
		dsl.Query = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	// Build _source from pushed-down fields
	if len(scan.Fields) > 0 {
		dsl.Source = scan.Fields
	}

	// Build sort from pushed-down sort keys
	if len(scan.SortKeys) > 0 {
		dsl.Sort = t.buildSort(scan.SortKeys)
	}

	// Build size from pushed-down limit
	if scan.Limit > 0 {
		dsl.Size = &scan.Limit
	}

	// Build script_fields from pushed-down computed fields (eval)
	if len(scan.ComputedFields) > 0 && t.functionBuilder != nil {
		scriptFields, err := t.buildScriptFields(scan.ComputedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to build script fields: %w", err)
		}
		if len(scriptFields) > 0 {
			dsl.ScriptFields = scriptFields
		}
	}

	// Check for Tier 1 operators that can be translated to aggregations

	// Check if there's a top operator in the plan
	if top := t.findTop(plan); top != nil {
		aggs, err := t.aggBuilder.BuildTopAggregations(top)
		if err != nil {
			return nil, fmt.Errorf("failed to build top aggregations: %w", err)
		}
		dsl.Aggregations = aggs
		zero := 0
		dsl.Size = &zero
		return dsl, nil
	}

	// Check if there's a rare operator in the plan
	if rare := t.findRare(plan); rare != nil {
		aggs, err := t.aggBuilder.BuildRareAggregations(rare)
		if err != nil {
			return nil, fmt.Errorf("failed to build rare aggregations: %w", err)
		}
		dsl.Aggregations = aggs
		zero := 0
		dsl.Size = &zero
		return dsl, nil
	}

	// Check if there's a bin operator in the plan (for date_histogram)
	if bin := t.findBin(plan); bin != nil {
		aggs, err := t.aggBuilder.BuildBinAggregations(bin)
		if err != nil {
			return nil, fmt.Errorf("failed to build bin aggregations: %w", err)
		}
		dsl.Aggregations = aggs
		zero := 0
		dsl.Size = &zero
		return dsl, nil
	}

	// Check if there's an aggregation in the plan
	if agg := t.findAggregation(plan); agg != nil {
		// Build aggregations
		aggs, err := t.aggBuilder.BuildAggregations(agg)
		if err != nil {
			return nil, fmt.Errorf("failed to build aggregations: %w", err)
		}
		dsl.Aggregations = aggs

		// For aggregation queries, set size to 0 (we only want aggregation results)
		zero := 0
		dsl.Size = &zero
	}

	return dsl, nil
}

// buildSort converts sort keys to OpenSearch sort DSL
func (t *Translator) buildSort(sortKeys []*ast.SortKey) []map[string]interface{} {
	sorts := make([]map[string]interface{}, 0, len(sortKeys))

	for _, key := range sortKeys {
		// Extract field name
		fieldRef, ok := key.Field.(*ast.FieldReference)
		if !ok {
			// Skip complex expressions (shouldn't happen with push-down)
			continue
		}

		order := "asc"
		if key.Descending {
			order = "desc"
		}

		sorts = append(sorts, map[string]interface{}{
			fieldRef.Name: map[string]interface{}{
				"order": order,
			},
		})
	}

	return sorts
}

// buildScriptFields builds script_fields from eval assignments using WASM UDFs
func (t *Translator) buildScriptFields(computedFields []*ast.EvalAssignment) (map[string]interface{}, error) {
	scriptFields := make(map[string]interface{})

	for _, assignment := range computedFields {
		// Build WASM UDF reference for this computed field
		udfRef, err := t.functionBuilder.BuildComputedField(assignment.Expression, assignment.Field)
		if err != nil {
			return nil, fmt.Errorf("failed to build UDF for field %s: %w", assignment.Field, err)
		}

		// Create script_field entry with wasm_udf reference
		scriptFields[assignment.Field] = map[string]interface{}{
			"wasm_udf": map[string]interface{}{
				"name":           udfRef.Name,
				"version":        udfRef.Version,
				"parameters":     udfRef.Parameters,
				"field_bindings": udfRef.FieldBindings,
			},
		}
	}

	return scriptFields, nil
}

// findAggregation finds the aggregation node in the plan (if any)
func (t *Translator) findAggregation(plan physical.PhysicalPlan) *physical.PhysicalAggregate {
	if agg, ok := plan.(*physical.PhysicalAggregate); ok {
		return agg
	}

	for _, child := range plan.Children() {
		if agg := t.findAggregation(child); agg != nil {
			return agg
		}
	}

	return nil
}

// findTop finds the top operator in the plan (if any)
func (t *Translator) findTop(plan physical.PhysicalPlan) *physical.PhysicalTop {
	if top, ok := plan.(*physical.PhysicalTop); ok {
		return top
	}

	for _, child := range plan.Children() {
		if top := t.findTop(child); top != nil {
			return top
		}
	}

	return nil
}

// findRare finds the rare operator in the plan (if any)
func (t *Translator) findRare(plan physical.PhysicalPlan) *physical.PhysicalRare {
	if rare, ok := plan.(*physical.PhysicalRare); ok {
		return rare
	}

	for _, child := range plan.Children() {
		if rare := t.findRare(child); rare != nil {
			return rare
		}
	}

	return nil
}

// findBin finds the bin operator in the plan (if any)
func (t *Translator) findBin(plan physical.PhysicalPlan) *physical.PhysicalBin {
	if bin, ok := plan.(*physical.PhysicalBin); ok {
		return bin
	}

	for _, child := range plan.Children() {
		if bin := t.findBin(child); bin != nil {
			return bin
		}
	}

	return nil
}

// TranslateToJSON converts a physical plan to JSON-serializable DSL
func (t *Translator) TranslateToJSON(plan physical.PhysicalPlan) (map[string]interface{}, error) {
	dsl, err := t.Translate(plan)
	if err != nil {
		return nil, err
	}

	// Convert to map
	result := make(map[string]interface{})

	if dsl.Query != nil {
		result["query"] = dsl.Query
	}

	if dsl.Source != nil {
		result["_source"] = dsl.Source
	}

	if len(dsl.Sort) > 0 {
		result["sort"] = dsl.Sort
	}

	if dsl.Size != nil {
		result["size"] = *dsl.Size
	}

	if dsl.From != nil {
		result["from"] = *dsl.From
	}

	if len(dsl.Aggregations) > 0 {
		result["aggs"] = dsl.Aggregations
	}

	return result, nil
}
