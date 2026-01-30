// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/conjugate/conjugate/pkg/ppl/dsl"
	"github.com/conjugate/conjugate/pkg/ppl/lookup"
	"github.com/conjugate/conjugate/pkg/ppl/physical"
	"go.uber.org/zap"
)

// DataSource provides search capabilities
type DataSource interface {
	// Search executes a search query and returns results
	Search(ctx context.Context, index string, queryDSL []byte, from, size int) (*SearchResult, error)
}

// SearchResult from data source
type SearchResult struct {
	TookMillis   int64
	TotalHits    int64
	MaxScore     float64
	Hits         []*SearchHit
	Aggregations map[string]*AggregationResult
}

// SearchHit represents a single search result
type SearchHit struct {
	ID     string
	Score  float64
	Source map[string]interface{}
}

// AggregationResult from data source
type AggregationResult struct {
	Type    string
	Value   interface{}
	Buckets []*AggregationBucket
}

// AggregationBucket for bucket aggregations
type AggregationBucket struct {
	Key      string
	DocCount int64
	SubAggs  map[string]*AggregationResult
}

// Executor executes PPL physical plans
type Executor struct {
	logger         *zap.Logger
	dataSource     DataSource
	translator     *dsl.Translator
	lookupRegistry *lookup.Registry
}

// NewExecutor creates a new executor
func NewExecutor(dataSource DataSource, translator *dsl.Translator, logger *zap.Logger) *Executor {
	return &Executor{
		logger:         logger.With(zap.String("component", "ppl_executor")),
		dataSource:     dataSource,
		translator:     translator,
		lookupRegistry: lookup.NewRegistry(logger),
	}
}

// SetLookupRegistry sets the lookup table registry
func (e *Executor) SetLookupRegistry(registry *lookup.Registry) {
	e.lookupRegistry = registry
}

// GetLookupRegistry returns the lookup table registry
func (e *Executor) GetLookupRegistry() *lookup.Registry {
	return e.lookupRegistry
}

// Execute executes a physical plan and returns a streaming result
func (e *Executor) Execute(ctx context.Context, plan physical.PhysicalPlan) (*Result, error) {
	startTime := time.Now()

	e.logger.Debug("Executing physical plan",
		zap.String("plan", physical.PrintPlan(plan, 0)))

	// Build the operator tree
	operator, err := e.buildOperator(plan)
	if err != nil {
		return nil, fmt.Errorf("failed to build operator: %w", err)
	}

	// Open the operator tree
	if err := operator.Open(ctx); err != nil {
		return nil, fmt.Errorf("failed to open operator: %w", err)
	}

	// Get aggregations if present
	var aggregations map[string]*AggregationValue
	if aggOp, ok := operator.(*aggregationOperator); ok {
		aggregations = aggOp.GetAggregations()
	}

	took := time.Since(startTime).Milliseconds()

	return &Result{
		Rows:         operator,
		Aggregations: aggregations,
		TookMillis:   took,
	}, nil
}

// buildOperator recursively builds the operator tree from a physical plan
func (e *Executor) buildOperator(plan physical.PhysicalPlan) (Operator, error) {
	switch p := plan.(type) {
	case *physical.PhysicalScan:
		return e.buildScanOperator(p)

	case *physical.PhysicalFilter:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewFilterOperator(input, p.Condition, e.logger), nil

	case *physical.PhysicalProject:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewProjectOperator(input, p.Fields, p.Exclude, e.logger), nil

	case *physical.PhysicalSort:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewSortOperator(input, p.SortKeys, e.logger), nil

	case *physical.PhysicalLimit:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewLimitOperator(input, p.Count, e.logger), nil

	case *physical.PhysicalAggregate:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewAggregationOperator(input, p.GroupBy, p.Aggregations, p.Algorithm, e.logger), nil

	// Tier 1 operators
	case *physical.PhysicalDedup:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewDedupOperator(input, p.Fields, p.Count, p.Consecutive, e.logger), nil

	case *physical.PhysicalBin:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewBinOperator(input, p.Field, p.Span, p.Bins, e.logger), nil

	case *physical.PhysicalTop:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewTopOperator(input, p.Fields, p.Limit, p.GroupBy, p.ShowCount, p.ShowPercent, e.logger), nil

	case *physical.PhysicalRare:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewRareOperator(input, p.Fields, p.Limit, p.GroupBy, p.ShowCount, p.ShowPercent, e.logger), nil

	case *physical.PhysicalEval:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewEvalOperator(input, p.Assignments, e.logger), nil

	case *physical.PhysicalRename:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewRenameOperator(input, p.Assignments, e.logger), nil

	case *physical.PhysicalReplace:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewReplaceOperator(input, p.Mappings, p.Field, e.logger), nil

	case *physical.PhysicalFillnull:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		// Use DefaultValue as the fill value, pass Fields as expressions
		return NewFillnullOperator(input, p.DefaultValue, p.Fields, e.logger), nil

	case *physical.PhysicalParse:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewParseOperator(input, p.SourceField, p.Pattern, p.ExtractedFields, e.logger)

	case *physical.PhysicalRex:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewRexOperator(input, p.SourceField, p.Pattern, p.ExtractedFields, e.logger)

	case *physical.PhysicalLookup:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewLookupOperator(
			input,
			e.lookupRegistry,
			p.TableName,
			p.JoinField,
			p.JoinFieldAlias,
			p.OutputFields,
			p.OutputAliases,
			e.logger,
		)

	case *physical.PhysicalAppend:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		subsearch, err := e.buildOperator(p.Subsearch)
		if err != nil {
			return nil, err
		}
		return NewAppendOperator(input, subsearch, e.logger), nil

	case *physical.PhysicalJoin:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		right, err := e.buildOperator(p.Right)
		if err != nil {
			return nil, err
		}
		return NewJoinOperator(input, right, p.JoinType, p.JoinField, p.RightField, e.logger), nil

	case *physical.PhysicalReverse:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewReverseOperator(input, e.logger), nil

	case *physical.PhysicalFlatten:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewFlattenOperator(input, p.Field, e.logger), nil

	case *physical.PhysicalTable:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewTableOperator(input, p.Fields, e.logger), nil

	case *physical.PhysicalEventstats:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewEventstatsOperator(input, p.GroupBy, p.Aggregations, p.BucketNullable, e.logger), nil

	case *physical.PhysicalStreamstats:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewStreamstatsOperator(input, p.GroupBy, p.Aggregations, p.Window, p.Global, p.ResetBefore, p.ResetAfter, e.logger), nil

	case *physical.PhysicalAddtotals:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewAddtotalsOperator(input, p.Fields, p.Row, p.Col, p.LabelField, p.Label, p.FieldName, e.logger), nil

	case *physical.PhysicalAddcoltotals:
		input, err := e.buildOperator(p.Input)
		if err != nil {
			return nil, err
		}
		return NewAddcoltotalsOperator(input, p.Fields, p.LabelField, p.Label, e.logger), nil

	default:
		return nil, fmt.Errorf("unsupported physical plan type: %T", plan)
	}
}

// buildScanOperator creates a scan operator
func (e *Executor) buildScanOperator(scan *physical.PhysicalScan) (Operator, error) {
	// Build a temporary physical plan for the scan only
	scanPlan := &physical.PhysicalScan{
		Source:       scan.Source,
		OutputSchema: scan.OutputSchema,
		Filter:       scan.Filter,
		Fields:       scan.Fields,
		SortKeys:     scan.SortKeys,
		Limit:        scan.Limit,
	}

	// Translate to DSL
	dslJSON, err := e.translator.TranslateToJSON(scanPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to translate scan to DSL: %w", err)
	}

	// Serialize DSL to JSON bytes
	queryBytes, err := json.Marshal(dslJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize DSL: %w", err)
	}

	// Determine size
	size := scan.Limit
	if size == 0 {
		size = 10000 // Default max size
	}

	return NewScanOperator(e.dataSource, scan.Source, queryBytes, 0, size, e.logger), nil
}

// Operator is the base interface for all operators
type Operator interface {
	RowIterator

	// Open initializes the operator
	Open(ctx context.Context) error
}

// baseOperator provides common functionality
type baseOperator struct {
	input  Operator
	logger *zap.Logger
	stats  *IteratorStats
	ctx    context.Context
	closed bool
}

func (b *baseOperator) Stats() *IteratorStats {
	return b.stats
}

func (b *baseOperator) Close() error {
	b.closed = true
	if b.input != nil {
		return b.input.Close()
	}
	return nil
}
