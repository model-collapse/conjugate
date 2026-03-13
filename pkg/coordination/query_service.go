package coordination

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
	"github.com/conjugate/conjugate/pkg/coordination/cache"
	"github.com/conjugate/conjugate/pkg/coordination/executor"
	"github.com/conjugate/conjugate/pkg/coordination/parser"
	"github.com/conjugate/conjugate/pkg/coordination/pipeline"
	"github.com/conjugate/conjugate/pkg/coordination/planner"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

// aggResponseCacheEntry stores a cached aggregation response
type aggResponseCacheEntry struct {
	result   *SearchResult
	cachedAt time.Time
}

const aggResponseCacheTTL = 30 * time.Second

// Prometheus metrics for query service
var (
	queryPlanningTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "conjugate_query_planning_seconds",
			Help:    "Query planning time in seconds",
			Buckets: prometheus.ExponentialBuckets(0.0001, 2, 12), // 0.1ms to ~400ms
		},
		[]string{"index", "stage"},
	)

	queryOptimizationTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "conjugate_query_optimization_seconds",
			Help:    "Query optimization time in seconds",
			Buckets: prometheus.ExponentialBuckets(0.0001, 2, 12), // 0.1ms to ~400ms
		},
		[]string{"index"},
	)

	queryExecutionTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "conjugate_query_execution_seconds",
			Help:    "Query execution time in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 14), // 1ms to ~16s
		},
		[]string{"index", "status"},
	)

	logicalPlanComplexity = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "conjugate_logical_plan_complexity",
			Help:    "Logical plan complexity (estimated cardinality)",
			Buckets: prometheus.ExponentialBuckets(1, 2, 20), // 1 to ~1M
		},
		[]string{"index"},
	)

	optimizationPassCount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "conjugate_optimization_passes",
			Help:    "Number of optimization passes applied",
			Buckets: prometheus.LinearBuckets(0, 1, 10), // 0 to 10 passes
		},
		[]string{"index"},
	)
)

// QueryService provides high-level query execution with the complete planner pipeline
type QueryService struct {
	logger           *zap.Logger
	queryParser      *parser.QueryParser
	converter        *planner.Converter
	optimizer        *planner.Optimizer
	costModel        *planner.CostModel
	physicalPlanner  *planner.Planner
	queryExecutor    queryExecutorInterface
	masterClient     masterClientInterface
	queryCache       *cache.QueryCache
	pipelineRegistry *pipeline.Registry
	pipelineExecutor *pipeline.Executor
	aggResponseCache sync.Map // key: string (sha256 of index+query), value: *aggResponseCacheEntry
}

// queryExecutorInterface defines the methods needed from query executor
type queryExecutorInterface interface {
	ExecuteSearch(ctx context.Context, indexName string, query []byte, filterExpr []byte, from, size int, sort []string) (*executor.SearchResult, error)
}

// masterClientInterface defines the methods needed from master client
type masterClientInterface interface {
	GetShardRouting(ctx context.Context, indexName string) (map[int32]*pb.ShardRouting, error)
	GetIndexMetadata(ctx context.Context, indexName string) (*pb.IndexMetadataResponse, error)
}

// NewQueryService creates a new query service with the complete planner pipeline
func NewQueryService(
	queryExecutor queryExecutorInterface,
	masterClient masterClientInterface,
	logger *zap.Logger,
) *QueryService {
	return &QueryService{
		logger:           logger,
		queryParser:      parser.NewQueryParser(),
		converter:        planner.NewConverter(),
		optimizer:        planner.NewOptimizer(),
		costModel:        planner.NewDefaultCostModel(),
		physicalPlanner:  planner.NewPlanner(planner.NewDefaultCostModel()),
		queryExecutor:    queryExecutor,
		masterClient:     masterClient,
		queryCache:       cache.NewQueryCache(cache.DefaultQueryCacheConfig()),
		pipelineRegistry: nil, // Pipelines optional
		pipelineExecutor: nil,
	}
}

// NewQueryServiceWithCache creates a new query service with custom cache configuration
func NewQueryServiceWithCache(
	queryExecutor queryExecutorInterface,
	masterClient masterClientInterface,
	logger *zap.Logger,
	cacheConfig *cache.QueryCacheConfig,
) *QueryService {
	return &QueryService{
		logger:           logger,
		queryParser:      parser.NewQueryParser(),
		converter:        planner.NewConverter(),
		optimizer:        planner.NewOptimizer(),
		costModel:        planner.NewDefaultCostModel(),
		physicalPlanner:  planner.NewPlanner(planner.NewDefaultCostModel()),
		queryExecutor:    queryExecutor,
		masterClient:     masterClient,
		queryCache:       cache.NewQueryCache(cacheConfig),
		pipelineRegistry: nil, // Pipelines optional
		pipelineExecutor: nil,
	}
}

// SetPipelineComponents sets the pipeline registry and executor (optional)
func (qs *QueryService) SetPipelineComponents(registry *pipeline.Registry, executor *pipeline.Executor) {
	qs.pipelineRegistry = registry
	qs.pipelineExecutor = executor
}

// SearchResult represents a search result with all metadata
type SearchResult struct {
	TookMillis   int64
	TotalHits    int64
	MaxScore     float64
	Hits         []*SearchHit
	Aggregations map[string]*AggregationResult
	Shards       *ShardInfo
}

// SearchHit represents a single hit
type SearchHit struct {
	ID         string
	Score      float64
	Source     map[string]interface{}
	SortValues []float64
}

// AggregationResult represents an aggregation result
type AggregationResult struct {
	Type    string
	Buckets []*AggregationBucket

	// For single-value aggregations
	Value float64

	// For stats aggregations
	Count int64
	Min   float64
	Max   float64
	Avg   float64
	Sum   float64
}

// AggregationBucket represents a bucket in a bucket aggregation
type AggregationBucket struct {
	Key      interface{}
	DocCount int64
	SubAggs  map[string]*AggregationResult
}

// ShardInfo represents shard execution information
type ShardInfo struct {
	Total      int
	Successful int
	Skipped    int
	Failed     int
}

// ExecuteSearch executes a search query using the complete planner pipeline
func (qs *QueryService) ExecuteSearch(ctx context.Context, indexName string, requestBody []byte) (*SearchResult, error) {
	startTime := time.Now()

	qs.logger.Debug("QueryService.ExecuteSearch",
		zap.String("index", indexName),
		zap.Int("body_len", len(requestBody)))

	// Step 1: Parse query
	parseStart := time.Now()
	var searchReq *parser.SearchRequest
	var err error

	if len(requestBody) > 0 {
		searchReq, err = qs.queryParser.ParseSearchRequest(requestBody)
		if err != nil {
			qs.logger.Error("Failed to parse query", zap.Error(err))
			return nil, fmt.Errorf("failed to parse query: %w", err)
		}

		// Validate parsed query
		if searchReq.ParsedQuery != nil {
			if err := qs.queryParser.Validate(searchReq.ParsedQuery); err != nil {
				qs.logger.Error("Query validation failed", zap.Error(err))
				return nil, fmt.Errorf("query validation failed: %w", err)
			}
		}
	} else {
		// Empty body - match all query
		searchReq = &parser.SearchRequest{
			ParsedQuery: &parser.MatchAllQuery{},
			Size:        10,
		}
	}

	qs.logger.Debug("Query parsed",
		zap.String("index", indexName))

	queryPlanningTime.WithLabelValues(indexName, "parse").Observe(time.Since(parseStart).Seconds())

	// Step 1.5: Execute query pipeline if configured
	if qs.pipelineRegistry != nil && qs.pipelineExecutor != nil {
		queryPipelineStart := time.Now()
		modifiedReq, err := qs.executeQueryPipeline(ctx, indexName, searchReq)
		if err != nil {
			// Log warning but continue with original request (graceful degradation)
			qs.logger.Warn("Query pipeline failed, continuing with original request",
				zap.String("index", indexName),
				zap.Error(err))
		} else if modifiedReq != nil {
			searchReq = modifiedReq
			qs.logger.Info("Query pipeline executed successfully",
				zap.String("index", indexName),
				zap.Duration("duration", time.Since(queryPipelineStart)))
		}
		queryPlanningTime.WithLabelValues(indexName, "query_pipeline").Observe(time.Since(queryPipelineStart).Seconds())
	}

	// Aggregation fast path: for size=0 + aggs, bypass planner and forward
	// raw query JSON directly to data nodes. The data node has a native fast path
	// for aggregations that computes results in <1ms using C++ hash maps.
	if searchReq.Size == 0 && (len(searchReq.Aggregations) > 0 || len(searchReq.Aggs) > 0) {
		result, err := qs.executeAggFastPath(ctx, indexName, requestBody, startTime)
		if err == nil {
			return result, nil
		}
		qs.logger.Warn("Agg fast path failed, falling back to planner",
			zap.String("index", indexName),
			zap.Error(err))
	}

	// Direct search fast path: bypass the planner's expression conversion which
	// loses information for query types like query_string, bool filter, etc.
	// Forward the raw query JSON directly to data nodes, same as _count does.
	if len(searchReq.Aggregations) == 0 && len(searchReq.Aggs) == 0 {
		result, err := qs.executeDirectSearch(ctx, indexName, requestBody, searchReq, startTime)
		if err == nil {
			return result, nil
		}
		qs.logger.Warn("Direct search fast path failed, falling back to planner",
			zap.String("index", indexName),
			zap.Error(err))
	}

	// Step 2: Get shard routing for this index
	routing, err := qs.masterClient.GetShardRouting(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to get shard routing: %w", err)
	}

	// Extract shard IDs
	shardIDs := make([]int32, 0, len(routing))
	for shardID, shard := range routing {
		if shard.Allocation != nil && shard.Allocation.State == pb.ShardAllocation_SHARD_STATE_STARTED {
			shardIDs = append(shardIDs, shardID)
		}
	}

	if len(shardIDs) == 0 {
		return nil, fmt.Errorf("no active shards found for index %s", indexName)
	}

	// Step 3: Check logical plan cache or convert AST to Logical Plan
	convertStart := time.Now()
	var logicalPlan planner.LogicalPlan

	// Try to get from cache
	cachedLogicalPlan, found := qs.queryCache.GetLogicalPlan(indexName, searchReq, shardIDs)
	if found {
		logicalPlan = cachedLogicalPlan
		qs.logger.Debug("Logical plan retrieved from cache",
			zap.String("index", indexName),
			zap.String("plan", logicalPlan.String()))
	} else {
		// Convert AST to Logical Plan
		logicalPlan, err = qs.converter.ConvertSearchRequest(searchReq, indexName, shardIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to convert query to logical plan: %w", err)
		}
		qs.logger.Debug("Logical plan created",
			zap.String("index", indexName),
			zap.String("plan", logicalPlan.String()))
	}
	queryPlanningTime.WithLabelValues(indexName, "convert").Observe(time.Since(convertStart).Seconds())

	// Record logical plan complexity
	logicalPlanComplexity.WithLabelValues(indexName).Observe(float64(logicalPlan.Cardinality()))

	// Step 4: Optimize Logical Plan (if not from cache)
	var optimizedPlan planner.LogicalPlan
	if !found {
		// Plan was just created, so optimize it
		optimizeStart := time.Now()
		optimizedPlan, err = qs.optimizer.Optimize(logicalPlan)
		if err != nil {
			qs.logger.Warn("Optimization failed, using unoptimized plan",
				zap.String("index", indexName),
				zap.Error(err))
			optimizedPlan = logicalPlan
		}
		optimizeTime := time.Since(optimizeStart)
		queryOptimizationTime.WithLabelValues(indexName).Observe(optimizeTime.Seconds())

		// Record optimization passes (simplified - using 1 or 0)
		if optimizedPlan != logicalPlan {
			optimizationPassCount.WithLabelValues(indexName).Observe(1)
		} else {
			optimizationPassCount.WithLabelValues(indexName).Observe(0)
		}

		qs.logger.Debug("Logical plan optimized",
			zap.String("index", indexName),
			zap.String("optimized_plan", optimizedPlan.String()),
			zap.Duration("optimization_time", optimizeTime))

		// Cache the optimized logical plan
		qs.queryCache.PutLogicalPlan(indexName, searchReq, shardIDs, optimizedPlan)
	} else {
		// Plan from cache is already optimized
		optimizedPlan = logicalPlan
	}

	// Step 5: Check physical plan cache or create Physical Plan
	physicalStart := time.Now()
	var physicalPlan planner.PhysicalPlan

	// Try to get from cache
	cachedPhysicalPlan, foundPhysical := qs.queryCache.GetPhysicalPlan(indexName, optimizedPlan)
	if foundPhysical {
		physicalPlan = cachedPhysicalPlan
		qs.logger.Debug("Physical plan retrieved from cache",
			zap.String("index", indexName),
			zap.String("plan", physicalPlan.String()))
	} else {
		// Convert to Physical Plan
		physicalPlan, err = qs.physicalPlanner.Plan(optimizedPlan)
		if err != nil {
			return nil, fmt.Errorf("failed to create physical plan: %w", err)
		}

		qs.logger.Debug("Physical plan created",
			zap.String("index", indexName),
			zap.String("plan", physicalPlan.String()),
			zap.Float64("estimated_cost", physicalPlan.Cost().TotalCost))

		// Cache the physical plan
		qs.queryCache.PutPhysicalPlan(indexName, optimizedPlan, physicalPlan)
	}
	queryPlanningTime.WithLabelValues(indexName, "physical").Observe(time.Since(physicalStart).Seconds())

	// Step 6: Execute Physical Plan
	executeStart := time.Now()

	// Create execution context
	execCtx := &planner.ExecutionContext{
		QueryExecutor: qs.queryExecutor,
		Logger:        qs.logger,
	}
	ctxWithExec := planner.WithExecutionContext(ctx, execCtx)

	// Execute plan
	executionResult, err := physicalPlan.Execute(ctxWithExec)
	executeTime := time.Since(executeStart)

	if err != nil {
		queryExecutionTime.WithLabelValues(indexName, "error").Observe(executeTime.Seconds())
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	queryExecutionTime.WithLabelValues(indexName, "success").Observe(executeTime.Seconds())

	// Convert ExecutionResult to SearchResult
	totalTime := time.Since(startTime)
	result := qs.convertToSearchResult(executionResult, totalTime, len(shardIDs))

	// Step 7: Execute result pipeline if configured
	if qs.pipelineRegistry != nil && qs.pipelineExecutor != nil {
		resultPipelineStart := time.Now()
		modifiedResult, err := qs.executeResultPipeline(ctx, indexName, result, searchReq)
		if err != nil {
			// Log warning but continue with original results (graceful degradation)
			qs.logger.Warn("Result pipeline failed, continuing with original results",
				zap.String("index", indexName),
				zap.Error(err))
		} else if modifiedResult != nil {
			result = modifiedResult
			qs.logger.Info("Result pipeline executed successfully",
				zap.String("index", indexName),
				zap.Duration("duration", time.Since(resultPipelineStart)))
		}
		queryPlanningTime.WithLabelValues(indexName, "result_pipeline").Observe(time.Since(resultPipelineStart).Seconds())
	}

	qs.logger.Info("Query executed successfully",
		zap.String("index", indexName),
		zap.Int64("total_hits", result.TotalHits),
		zap.Int("hits_returned", len(result.Hits)),
		zap.Duration("total_time", totalTime),
		zap.Duration("execute_time", executeTime))

	return result, nil
}

// aggResponseCacheKey generates a cache key for aggregation responses
func aggResponseCacheKey(indexName string, rawQuery []byte) string {
	h := sha256.New()
	h.Write([]byte(indexName))
	h.Write(rawQuery)
	return hex.EncodeToString(h.Sum(nil))
}

// executeDirectSearch bypasses the planner and sends raw query JSON directly to
// data nodes. The planner's expression conversion loses information for certain
// query types (query_string, bool filter, must_not combos). This path mirrors
// how _count works — raw JSON goes straight to the data node's Diagon bridge.
func (qs *QueryService) executeDirectSearch(ctx context.Context, indexName string, rawBody []byte, searchReq *parser.SearchRequest, startTime time.Time) (*SearchResult, error) {
	qs.logger.Debug("Direct search fast path: bypassing planner",
		zap.String("index", indexName))

	// Extract just the query portion from the request body
	var bodyMap map[string]json.RawMessage
	if err := json.Unmarshal(rawBody, &bodyMap); err != nil {
		return nil, fmt.Errorf("failed to parse request body: %w", err)
	}

	queryBytes, hasQuery := bodyMap["query"]
	if !hasQuery {
		// No explicit query → match_all
		queryBytes = []byte(`{"match_all":{}}`)
	}

	// Extract sort if present
	var sort []string
	if sortRaw, ok := bodyMap["sort"]; ok {
		var sortArr []interface{}
		if err := json.Unmarshal(sortRaw, &sortArr); err == nil {
			for _, s := range sortArr {
				switch v := s.(type) {
				case string:
					sort = append(sort, v)
				case map[string]interface{}:
					for field, opts := range v {
						if optsMap, ok := opts.(map[string]interface{}); ok {
							if order, ok := optsMap["order"].(string); ok {
								sort = append(sort, field+":"+order)
							} else {
								sort = append(sort, field)
							}
						} else {
							sort = append(sort, field)
						}
					}
				}
			}
		}
	}

	size := searchReq.Size
	from := searchReq.From

	// Execute via QueryExecutor — fans out to all data node shards
	execResult, err := qs.queryExecutor.ExecuteSearch(ctx, indexName, queryBytes, nil, from, size, sort)
	if err != nil {
		return nil, err
	}

	totalTime := time.Since(startTime)
	result := &SearchResult{
		TookMillis:   totalTime.Milliseconds(),
		TotalHits:    execResult.TotalHits,
		MaxScore:     execResult.MaxScore,
		Hits:         make([]*SearchHit, len(execResult.Hits)),
		Aggregations: make(map[string]*AggregationResult),
		Shards:       &ShardInfo{Total: 1, Successful: 1},
	}

	for i, hit := range execResult.Hits {
		result.Hits[i] = &SearchHit{
			ID:         hit.ID,
			Score:      hit.Score,
			Source:     hit.Source,
			SortValues: hit.SortValues,
		}
	}

	qs.logger.Info("Direct search fast path succeeded",
		zap.String("index", indexName),
		zap.Int64("total_hits", execResult.TotalHits),
		zap.Int("hits", len(result.Hits)),
		zap.Duration("total_time", totalTime))

	return result, nil
}

// executeAggFastPath bypasses the planner and sends raw query JSON directly to
// data nodes for aggregation queries (size=0). The executor fans out to all shards,
// and each data node's native aggregation path computes results using C++ hash maps.
func (qs *QueryService) executeAggFastPath(ctx context.Context, indexName string, rawQuery []byte, startTime time.Time) (*SearchResult, error) {
	qs.logger.Debug("Agg fast path: bypassing planner for size=0 agg query",
		zap.String("index", indexName))

	// Check coordination-level agg response cache
	cacheKey := aggResponseCacheKey(indexName, rawQuery)
	if entry, ok := qs.aggResponseCache.Load(cacheKey); ok {
		cached := entry.(*aggResponseCacheEntry)
		if time.Since(cached.cachedAt) < aggResponseCacheTTL {
			totalTime := time.Since(startTime)
			// Return cached result with updated took time
			result := &SearchResult{
				TookMillis:   totalTime.Milliseconds(),
				TotalHits:    cached.result.TotalHits,
				MaxScore:     cached.result.MaxScore,
				Hits:         cached.result.Hits,
				Aggregations: cached.result.Aggregations,
				Shards:       cached.result.Shards,
			}
			qs.logger.Debug("Agg fast path: cache hit",
				zap.String("index", indexName),
				zap.Duration("total_time", totalTime))
			return result, nil
		}
		// Expired - delete
		qs.aggResponseCache.Delete(cacheKey)
	}

	// Cache miss - execute via gRPC
	// Use the existing executor fan-out — sends raw query bytes to all data nodes
	// with size=0. Data node detects embedded aggs and uses native C++ fast path.
	execResult, err := qs.queryExecutor.ExecuteSearch(ctx, indexName, rawQuery, nil, 0, 0, nil)
	if err != nil {
		return nil, err
	}

	totalTime := time.Since(startTime)
	result := &SearchResult{
		TookMillis:   totalTime.Milliseconds(),
		TotalHits:    execResult.TotalHits,
		MaxScore:     execResult.MaxScore,
		Hits:         []*SearchHit{},
		Aggregations: make(map[string]*AggregationResult),
		Shards:       &ShardInfo{Total: 1, Successful: 1},
	}

	// Convert executor aggregations to coordination SearchResult aggregations
	for name, agg := range execResult.Aggregations {
		result.Aggregations[name] = convertExecutorAggregation(agg)
	}

	// Store in coordination-level cache
	qs.aggResponseCache.Store(cacheKey, &aggResponseCacheEntry{
		result:   result,
		cachedAt: time.Now(),
	})

	qs.logger.Info("Agg fast path succeeded",
		zap.String("index", indexName),
		zap.Int64("total_hits", execResult.TotalHits),
		zap.Int("agg_results", len(result.Aggregations)),
		zap.Duration("total_time", totalTime))

	return result, nil
}

// convertExecutorAggregation converts executor.AggregationResult to coordination AggregationResult
func convertExecutorAggregation(agg *executor.AggregationResult) *AggregationResult {
	if agg == nil {
		return &AggregationResult{}
	}

	result := &AggregationResult{
		Type:    agg.Type,
		Buckets: make([]*AggregationBucket, len(agg.Buckets)),
		Value:   float64(agg.Value),
		Count:   agg.Count,
		Min:     agg.Min,
		Max:     agg.Max,
		Avg:     agg.Avg,
		Sum:     agg.Sum,
	}

	for i, b := range agg.Buckets {
		key := interface{}(b.Key)
		if b.Key == "" && b.NumericKey != 0 {
			key = b.NumericKey
		}
		bucket := &AggregationBucket{
			Key:      key,
			DocCount: b.DocCount,
		}
		if len(b.SubAggs) > 0 {
			bucket.SubAggs = make(map[string]*AggregationResult, len(b.SubAggs))
			for subName, subAgg := range b.SubAggs {
				bucket.SubAggs[subName] = convertExecutorAggregation(subAgg)
			}
		}
		result.Buckets[i] = bucket
	}

	return result
}

// convertToSearchResult converts ExecutionResult to SearchResult
func (qs *QueryService) convertToSearchResult(execResult *planner.ExecutionResult, totalTime time.Duration, totalShards int) *SearchResult {
	result := &SearchResult{
		TookMillis:   totalTime.Milliseconds(),
		TotalHits:    execResult.TotalHits,
		MaxScore:     execResult.MaxScore,
		Hits:         make([]*SearchHit, len(execResult.Rows)),
		Aggregations: make(map[string]*AggregationResult),
		Shards: &ShardInfo{
			Total:      totalShards,
			Successful: totalShards,
			Skipped:    0,
			Failed:     0,
		},
	}

	// Convert hits
	for i, row := range execResult.Rows {
		hit := &SearchHit{
			Source: make(map[string]interface{}),
		}

		// Extract _id and _score
		if id, ok := row["_id"].(string); ok {
			hit.ID = id
			delete(row, "_id")
		}
		if score, ok := row["_score"].(float64); ok {
			hit.Score = score
			delete(row, "_score")
		}
		// Extract sort values
		if rawSV, exists := row["_sort_values"]; exists {
			if sv, ok := rawSV.([]float64); ok {
				hit.SortValues = sv
			}
			delete(row, "_sort_values")
		}

		// Copy remaining fields to source
		for k, v := range row {
			hit.Source[k] = v
		}

		result.Hits[i] = hit
	}

	// Convert aggregations
	for name, agg := range execResult.Aggregations {
		result.Aggregations[name] = qs.convertAggregation(agg)
	}

	return result
}

// convertAggregation converts planner.AggregationResult to SearchResult.AggregationResult
func (qs *QueryService) convertAggregation(agg *planner.AggregationResult) *AggregationResult {
	result := &AggregationResult{
		Type:    string(agg.Type),
		Buckets: make([]*AggregationBucket, len(agg.Buckets)),
		Value:   agg.Value,
	}

	// Convert buckets
	for i, bucket := range agg.Buckets {
		b := &AggregationBucket{
			Key:      bucket.Key,
			DocCount: bucket.DocCount,
		}
		if len(bucket.SubAggs) > 0 {
			b.SubAggs = make(map[string]*AggregationResult, len(bucket.SubAggs))
			for subName, subAgg := range bucket.SubAggs {
				b.SubAggs[subName] = qs.convertAggregation(subAgg)
			}
		}
		result.Buckets[i] = b
	}

	// For stats aggregations
	if agg.Stats != nil {
		result.Count = agg.Stats.Count
		result.Min = agg.Stats.Min
		result.Max = agg.Stats.Max
		result.Avg = agg.Stats.Avg
		result.Sum = agg.Stats.Sum
	}

	return result
}

// executeQueryPipeline executes the query pipeline for an index if configured
func (qs *QueryService) executeQueryPipeline(ctx context.Context, indexName string, req *parser.SearchRequest) (*parser.SearchRequest, error) {
	// Get query pipeline for this index
	pipe, err := qs.pipelineRegistry.GetPipelineForIndex(indexName, pipeline.PipelineTypeQuery)
	if err != nil {
		// No pipeline configured - not an error
		return req, nil
	}

	qs.logger.Debug("Executing query pipeline",
		zap.String("index", indexName),
		zap.String("pipeline", pipe.Name()))

	// Convert SearchRequest to map for pipeline execution
	requestMap, err := qs.searchRequestToMap(req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert search request: %w", err)
	}

	// Execute pipeline
	output, err := pipe.Execute(ctx, requestMap)
	if err != nil {
		return nil, fmt.Errorf("pipeline execution failed: %w", err)
	}

	// Convert back to SearchRequest
	outputMap, ok := output.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("pipeline output is not a map, got %T", output)
	}

	modifiedReq, err := qs.mapToSearchRequest(outputMap)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pipeline output: %w", err)
	}

	return modifiedReq, nil
}

// executeResultPipeline executes the result pipeline for an index if configured
func (qs *QueryService) executeResultPipeline(ctx context.Context, indexName string, result *SearchResult, originalReq *parser.SearchRequest) (*SearchResult, error) {
	// Get result pipeline for this index
	pipe, err := qs.pipelineRegistry.GetPipelineForIndex(indexName, pipeline.PipelineTypeResult)
	if err != nil {
		// No pipeline configured - not an error
		return result, nil
	}

	qs.logger.Debug("Executing result pipeline",
		zap.String("index", indexName),
		zap.String("pipeline", pipe.Name()))

	// Convert SearchResult to map for pipeline execution
	// Include both results and original request for context
	pipelineInput := map[string]interface{}{
		"results": qs.searchResultToMap(result),
		"request": qs.searchRequestToMapSimple(originalReq),
	}

	// Execute pipeline
	output, err := pipe.Execute(ctx, pipelineInput)
	if err != nil {
		return nil, fmt.Errorf("pipeline execution failed: %w", err)
	}

	// Convert back to SearchResult
	outputMap, ok := output.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("pipeline output is not a map, got %T", output)
	}

	// Extract results from output
	resultsMap, ok := outputMap["results"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("pipeline output missing 'results' field")
	}

	modifiedResult, err := qs.mapToSearchResult(resultsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to convert pipeline output: %w", err)
	}

	return modifiedResult, nil
}

// searchRequestToMap converts SearchRequest to map for pipeline
func (qs *QueryService) searchRequestToMap(req *parser.SearchRequest) (map[string]interface{}, error) {
	// Marshal to JSON then unmarshal to map (simple conversion)
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// searchRequestToMapSimple converts SearchRequest to map (simplified version)
func (qs *QueryService) searchRequestToMapSimple(req *parser.SearchRequest) map[string]interface{} {
	return map[string]interface{}{
		"from": req.From,
		"size": req.Size,
		// Add other relevant fields as needed
	}
}

// mapToSearchRequest converts map back to SearchRequest
func (qs *QueryService) mapToSearchRequest(m map[string]interface{}) (*parser.SearchRequest, error) {
	// Marshal to JSON then unmarshal to SearchRequest
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	// Parse the modified JSON
	return qs.queryParser.ParseSearchRequest(data)
}

// searchResultToMap converts SearchResult to map for pipeline
func (qs *QueryService) searchResultToMap(result *SearchResult) map[string]interface{} {
	// Convert hits
	hits := make([]interface{}, len(result.Hits))
	for i, hit := range result.Hits {
		hits[i] = map[string]interface{}{
			"_id":     hit.ID,
			"_score":  hit.Score,
			"_source": hit.Source,
		}
	}

	return map[string]interface{}{
		"took":       result.TookMillis,
		"total_hits": result.TotalHits,
		"max_score":  result.MaxScore,
		"hits":       hits,
	}
}

// mapToSearchResult converts map back to SearchResult
func (qs *QueryService) mapToSearchResult(m map[string]interface{}) (*SearchResult, error) {
	result := &SearchResult{
		Aggregations: make(map[string]*AggregationResult),
		Hits:         []*SearchHit{},
		Shards: &ShardInfo{
			Total:      0,
			Successful: 0,
			Skipped:    0,
			Failed:     0,
		},
	}

	// Extract basic fields
	if took, ok := m["took"].(float64); ok {
		result.TookMillis = int64(took)
	}
	if took, ok := m["took"].(int64); ok {
		result.TookMillis = took
	}
	if totalHits, ok := m["total_hits"].(float64); ok {
		result.TotalHits = int64(totalHits)
	}
	if totalHits, ok := m["total_hits"].(int64); ok {
		result.TotalHits = totalHits
	}
	if maxScore, ok := m["max_score"].(float64); ok {
		result.MaxScore = maxScore
	}

	// Extract hits
	if hitsData, ok := m["hits"].([]interface{}); ok {
		result.Hits = make([]*SearchHit, 0, len(hitsData))
		for _, hitData := range hitsData {
			hitMap, ok := hitData.(map[string]interface{})
			if !ok {
				continue
			}

			hit := &SearchHit{
				Source: make(map[string]interface{}),
			}
			if id, ok := hitMap["_id"].(string); ok {
				hit.ID = id
			}
			if score, ok := hitMap["_score"].(float64); ok {
				hit.Score = score
			}
			if source, ok := hitMap["_source"].(map[string]interface{}); ok {
				hit.Source = source
			}

			result.Hits = append(result.Hits, hit)
		}
	}

	return result, nil
}
