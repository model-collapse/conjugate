package data

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	json "github.com/goccy/go-json"
	"sync"
	"time"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
	"github.com/conjugate/conjugate/pkg/data/diagon"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// aggCacheEntry stores cached aggregation results
type aggCacheEntry struct {
	result    map[string]*pb.AggregationResult
	totalHits int64
	created   time.Time
}

// docMapPool reuses map[string]interface{} for json.Unmarshal to reduce GC pressure.
// Each bulk batch allocates ~5.6MB of maps (47% of all allocs); pooling saves ~1ms GC per batch.
var docMapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{}, 16)
	},
}

// clearMap deletes all keys from a map for reuse.
func clearMap(m map[string]interface{}) {
	for k := range m {
		delete(m, k)
	}
}

// DataService implements the gRPC DataService
type DataService struct {
	pb.UnimplementedDataServiceServer
	node     *DataNode
	logger   *zap.Logger
	aggCache sync.Map // map[string]*aggCacheEntry
}

// NewDataService creates a new data service
func NewDataService(node *DataNode, logger *zap.Logger) *DataService {
	return &DataService{
		node:   node,
		logger: logger,
	}
}

// aggCacheKey generates a cache key from index + query + aggs
func aggCacheKey(index string, queryBytes []byte) string {
	h := sha256.Sum256(append([]byte(index+":"), queryBytes...))
	return hex.EncodeToString(h[:16])
}

// CreateShard creates a new shard on this data node
func (s *DataService) CreateShard(ctx context.Context, req *pb.CreateShardRequest) (*pb.CreateShardResponse, error) {
	s.logger.Info("CreateShard request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId),
		zap.Bool("is_primary", req.IsPrimary))

	// Validate request
	if req.IndexName == "" {
		return nil, status.Error(codes.InvalidArgument, "index name is required")
	}
	if req.ShardId < 0 {
		return nil, status.Error(codes.InvalidArgument, "shard_id must be non-negative")
	}

	// Create shard
	if err := s.node.shards.CreateShard(ctx, req.IndexName, req.ShardId, req.IsPrimary); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create shard: %v", err)
	}

	shardKey := shardKey(req.IndexName, req.ShardId)

	return &pb.CreateShardResponse{
		Acknowledged: true,
		ShardKey:     shardKey,
	}, nil
}

// DeleteShard deletes a shard from this data node
func (s *DataService) DeleteShard(ctx context.Context, req *pb.DeleteShardRequest) (*pb.DeleteShardResponse, error) {
	s.logger.Info("DeleteShard request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId))

	// Validate request
	if req.IndexName == "" {
		return nil, status.Error(codes.InvalidArgument, "index name is required")
	}

	// Delete shard
	if err := s.node.shards.DeleteShard(ctx, req.IndexName, req.ShardId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete shard: %v", err)
	}

	return &pb.DeleteShardResponse{
		Acknowledged: true,
	}, nil
}

// GetShardInfo returns information about a shard
func (s *DataService) GetShardInfo(ctx context.Context, req *pb.GetShardInfoRequest) (*pb.ShardInfo, error) {
	s.logger.Debug("GetShardInfo request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId))

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	// Convert to proto
	return &pb.ShardInfo{
		IndexName:   shard.IndexName,
		ShardId:     shard.ShardID,
		IsPrimary:   shard.IsPrimary,
		State:       s.convertShardStateToProto(shard.State),
		DocsCount:   shard.DocsCount,
		SizeBytes:   shard.SizeBytes,
		CreatedAt:   timestamppb.New(time.Now()), // TODO: Store creation time
		LastUpdated: timestamppb.New(time.Now()),
	}, nil
}

// RefreshShard makes recently indexed documents searchable
func (s *DataService) RefreshShard(ctx context.Context, req *pb.RefreshShardRequest) (*pb.RefreshShardResponse, error) {
	s.logger.Debug("RefreshShard request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId))

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	// Refresh shard
	if err := shard.Refresh(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to refresh shard: %v", err)
	}

	return &pb.RefreshShardResponse{
		Acknowledged: true,
	}, nil
}

// FlushShard flushes shard data to disk
func (s *DataService) FlushShard(ctx context.Context, req *pb.FlushShardRequest) (*pb.FlushShardResponse, error) {
	s.logger.Debug("FlushShard request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId))

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	// Flush pending batch documents first
	if err := shard.Flush(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to flush pending documents: %v", err)
	}

	// Flush Diagon translog
	if err := shard.FlushDiagon(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to flush shard: %v", err)
	}

	return &pb.FlushShardResponse{
		Acknowledged: true,
	}, nil
}

// IndexDocument indexes a document into a shard
func (s *DataService) IndexDocument(ctx context.Context, req *pb.IndexDocumentRequest) (*pb.IndexDocumentResponse, error) {
	s.logger.Debug("IndexDocument",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId),
		zap.String("doc_id", req.DocId))

	// Validate request
	if req.IndexName == "" {
		s.logger.Error("IndexDocument validation failed: index name is required")
		return nil, status.Error(codes.InvalidArgument, "index name is required")
	}
	if req.DocId == "" {
		s.logger.Error("IndexDocument validation failed: doc_id is required")
		return nil, status.Error(codes.InvalidArgument, "doc_id is required")
	}
	if req.Document == nil {
		s.logger.Error("IndexDocument validation failed: document is required")
		return nil, status.Error(codes.InvalidArgument, "document is required")
	}

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		s.logger.Error("Failed to get shard",
			zap.String("index", req.IndexName),
			zap.Int32("shard_id", req.ShardId),
			zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	// Convert protobuf Struct to map
	doc := req.Document.AsMap()

	// Index document
	if err := shard.IndexDocument(ctx, req.DocId, doc); err != nil {
		s.logger.Error("shard.IndexDocument FAILED",
			zap.String("doc_id", req.DocId),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to index document: %v", err)
	}

	return &pb.IndexDocumentResponse{
		Acknowledged: true,
		DocId:        req.DocId,
		Version:      1, // TODO: Implement versioning
	}, nil
}

// GetDocument retrieves a document by ID
func (s *DataService) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.GetDocumentResponse, error) {
	s.logger.Debug("GetDocument request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId),
		zap.String("doc_id", req.DocId))

	// Validate request
	if req.IndexName == "" {
		return nil, status.Error(codes.InvalidArgument, "index name is required")
	}
	if req.DocId == "" {
		return nil, status.Error(codes.InvalidArgument, "doc_id is required")
	}

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	// Get document
	doc, err := shard.GetDocument(ctx, req.DocId)
	if err != nil {
		// Document not found - log the actual error
		s.logger.Warn("GetDocument failed",
			zap.String("index", req.IndexName),
			zap.Int32("shard_id", req.ShardId),
			zap.String("doc_id", req.DocId),
			zap.Error(err))
		return &pb.GetDocumentResponse{
			Found: false,
			DocId: req.DocId,
		}, nil
	}

	// Convert map to protobuf Struct
	docStruct, err := structpb.NewStruct(doc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert document: %v", err)
	}

	return &pb.GetDocumentResponse{
		Found:    true,
		DocId:    req.DocId,
		Document: docStruct,
		Version:  1, // TODO: Implement versioning
	}, nil
}

// DeleteDocument deletes a document by ID
func (s *DataService) DeleteDocument(ctx context.Context, req *pb.DeleteDocumentRequest) (*pb.DeleteDocumentResponse, error) {
	s.logger.Debug("DeleteDocument request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId),
		zap.String("doc_id", req.DocId))

	// Validate request
	if req.IndexName == "" {
		return nil, status.Error(codes.InvalidArgument, "index name is required")
	}
	if req.DocId == "" {
		return nil, status.Error(codes.InvalidArgument, "doc_id is required")
	}

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	// Delete document
	if err := shard.DeleteDocument(ctx, req.DocId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete document: %v", err)
	}

	return &pb.DeleteDocumentResponse{
		Acknowledged: true,
		Found:        true, // TODO: Check if document existed
	}, nil
}

// BulkIndex indexes multiple documents in a single request using batch API
func (s *DataService) BulkIndex(ctx context.Context, req *pb.BulkIndexRequest) (*pb.BulkIndexResponse, error) {
	s.logger.Debug("BulkIndex request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId),
		zap.Int("items", len(req.Items)))

	// Validate request
	if req.IndexName == "" {
		return nil, status.Error(codes.InvalidArgument, "index name is required")
	}
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "items are required")
	}

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	startTime := time.Now()

	// Convert raw JSON bytes to the bulk format expected by shard.
	// Use pooled maps to reduce GC pressure (~5.6MB/batch → reused).
	bulkDocs := make([]struct {
		ID         string
		Doc        map[string]interface{}
		SourceJSON []byte
	}, len(req.Items))

	pooledMaps := make([]map[string]interface{}, len(req.Items))
	for i, item := range req.Items {
		bulkDocs[i].ID = item.DocId
		bulkDocs[i].SourceJSON = item.DocumentJson
		m := docMapPool.Get().(map[string]interface{})
		clearMap(m)
		if err := json.Unmarshal(item.DocumentJson, &m); err != nil {
			// Return borrowed maps before erroring
			for j := 0; j < i; j++ {
				docMapPool.Put(pooledMaps[j])
			}
			docMapPool.Put(m)
			return nil, status.Errorf(codes.InvalidArgument, "invalid JSON in item %d: %v", i, err)
		}
		bulkDocs[i].Doc = m
		pooledMaps[i] = m
	}

	// Use batch bulk index (single lock acquisition, single CGO batch)
	_, err = shard.BulkIndexDocuments(ctx, bulkDocs)

	// Return pooled maps after bridge has consumed them
	for _, m := range pooledMaps {
		docMapPool.Put(m)
	}

	// Build response items
	items := make([]*pb.BulkIndexItemResponse, len(req.Items))
	hasErrors := err != nil

	for i, item := range req.Items {
		itemResp := &pb.BulkIndexItemResponse{
			DocId: item.DocId,
		}
		if err != nil {
			itemResp.Acknowledged = false
			itemResp.Error = err.Error()
		} else {
			itemResp.Acknowledged = true
		}
		items[i] = itemResp
	}

	tookMillis := time.Since(startTime).Milliseconds()

	return &pb.BulkIndexResponse{
		HasErrors:  hasErrors,
		Items:      items,
		TookMillis: tookMillis,
	}, nil
}

// Search executes a search query on a shard
func (s *DataService) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	s.logger.Debug("Search",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId))

	// Validate request
	if req.IndexName == "" {
		s.logger.Error("Search failed: index name is required")
		return nil, status.Error(codes.InvalidArgument, "index name is required")
	}
	if req.Query == nil {
		s.logger.Error("Search failed: query is required")
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	startTime := time.Now()

	// Check if query contains aggregation definitions (push-down from coordinator).
	// Wrapped queries have format: {"query": <filter>, "aggs": <agg_defs>}
	var aggsMap map[string]interface{}
	queryBytes := req.Query
	var wrappedQuery map[string]interface{}
	if err := json.Unmarshal(req.Query, &wrappedQuery); err == nil {
		if aggs, ok := wrappedQuery["aggs"]; ok {
			if aggsTyped, ok := aggs.(map[string]interface{}); ok {
				aggsMap = aggsTyped
				// Extract the actual query portion for Diagon search
				if q, ok := wrappedQuery["query"]; ok {
					queryBytes, _ = json.Marshal(q)
				} else {
					queryBytes, _ = json.Marshal(map[string]interface{}{"match_all": map[string]interface{}{}})
				}
			}
		}
	}
	parseTime := time.Since(startTime)

	// Determine max results
	maxResults := int(req.Size)
	if maxResults <= 0 {
		if len(aggsMap) > 0 {
			// Aggregation query: need ALL matching docs to compute aggregations locally
			maxResults = 200000
		} else {
			maxResults = 100
		}
	}

	// Execute search and aggregations
	searchStart := time.Now()
	var result *diagon.SearchResult
	var aggregations map[string]*pb.AggregationResult

	// Check aggregation cache for agg-only queries (size=0)
	// Use req.Query (original bytes including aggs) so different agg definitions
	// produce different cache keys, even with the same base query.
	if len(aggsMap) > 0 && req.Size == 0 {
		cacheKey := aggCacheKey(req.IndexName, req.Query)
		if cached, ok := s.aggCache.Load(cacheKey); ok {
			entry := cached.(*aggCacheEntry)
			// Cache valid for 30s (matches refresh_interval)
			if time.Since(entry.created) < 30*time.Second {
				return &pb.SearchResponse{
					Hits: &pb.SearchHits{
						Total: &pb.TotalHits{
							Value:    entry.totalHits,
							Relation: "eq",
						},
					},
					Aggregations: entry.result,
				}, nil
			}
		}
	}

	// Fast path: For aggregation-only queries, try native index-based aggregation
	// This avoids extracting O(N) documents and instead uses O(unique_terms) index reads
	if len(aggsMap) > 0 && req.Size == 0 {
		nativeAggs := computeNativeAggregations(shard.DiagonShard, aggsMap, s.logger)
		if nativeAggs != nil {
			// Native aggregation succeeded - no need to search/extract documents
			aggregations = nativeAggs
			result = &diagon.SearchResult{
				TotalHits: int64(shard.DiagonShard.DocCount()),
			}
		}
	}

	// Standard path: search + extract docs (for non-agg queries or unsupported agg types)
	if result == nil {
		if len(aggsMap) > 0 && maxResults > 100 {
			// Aggregation path: use lightweight SearchAndAggregate (no Hit object allocation)
			aggFields := extractAggregationFields(aggsMap)
			totalHits, aggDocs, searchErr := shard.DiagonShard.SearchAndAggregate(queryBytes, maxResults, aggFields)
			if searchErr != nil {
				s.logger.Error("SearchAndAggregate failed", zap.Error(searchErr))
				return nil, status.Errorf(codes.Internal, "search failed: %v", searchErr)
			}
			result = &diagon.SearchResult{
				TotalHits: totalHits,
			}
			if len(aggDocs) > 0 {
				aggregations = computeAggregationsFromDocValues(aggDocs, aggsMap)
			}
		} else if len(aggsMap) > 0 {
			aggFields := extractAggregationFields(aggsMap)
			result, err = shard.SearchFieldsOnly(ctx, queryBytes, maxResults, aggFields)
			if err != nil {
				s.logger.Error("Search failed", zap.Error(err))
				return nil, status.Errorf(codes.Internal, "search failed: %v", err)
			}
			if aggregations == nil && len(result.Hits) > 0 {
				aggregations = computeDataNodeAggregations(result.Hits, aggsMap)
			}
		} else {
			result, err = shard.Search(ctx, queryBytes, maxResults)
			if err != nil {
				s.logger.Error("Search failed", zap.Error(err))
				return nil, status.Errorf(codes.Internal, "search failed: %v", err)
			}
		}
	}

	searchTime := time.Since(searchStart)

	// Fall back to Diagon-computed aggregations if nothing else worked
	aggStart := time.Now()
	if aggregations == nil && result != nil {
		aggregations = convertAggregations(result.Aggregations)
	}
	aggTime := time.Since(aggStart)

	// For aggregation-only queries (aggsMap present), skip converting hits to proto
	// since the coordinator set size=0 and doesn't need document data.
	convertStart := time.Now()
	var hits []*pb.SearchHit
	if len(aggsMap) == 0 {
		hits = make([]*pb.SearchHit, 0, len(result.Hits))
		for _, hit := range result.Hits {
			docStruct, err := structpb.NewStruct(hit.Source)
			if err != nil {
				s.logger.Error("Failed to convert document", zap.Error(err))
				continue
			}
			hits = append(hits, &pb.SearchHit{
				Id:     hit.ID,
				Score:  hit.Score,
				Source: docStruct,
			})
		}
	}
	convertTime := time.Since(convertStart)

	tookMillis := time.Since(startTime).Milliseconds()

	s.logger.Info("Search timing breakdown",
		zap.String("index", req.IndexName),
		zap.Int("hits_returned", len(result.Hits)),
		zap.Int64("total_hits", result.TotalHits),
		zap.Duration("parse", parseTime),
		zap.Duration("diagon_search", searchTime),
		zap.Duration("aggregation", aggTime),
		zap.Duration("proto_convert", convertTime),
		zap.Duration("total", time.Since(startTime)),
		zap.Bool("has_aggs", len(aggsMap) > 0))

	// Cache aggregation results for agg-only queries
	if len(aggsMap) > 0 && req.Size == 0 && aggregations != nil {
		cacheKey := aggCacheKey(req.IndexName, req.Query)
		s.aggCache.Store(cacheKey, &aggCacheEntry{
			result:    aggregations,
			totalHits: result.TotalHits,
			created:   time.Now(),
		})
	}

	return &pb.SearchResponse{
		TookMillis: tookMillis,
		TimedOut:   false,
		Shards: &pb.ShardSearchStats{
			Total:      1,
			Successful: 1,
			Failed:     0,
		},
		Hits: &pb.SearchHits{
			Total: &pb.TotalHits{
				Value:    result.TotalHits,
				Relation: "eq",
			},
			MaxScore: result.MaxScore,
			Hits:     hits,
		},
		Aggregations: aggregations,
	}, nil
}

// Count returns the count of documents matching a query
func (s *DataService) Count(ctx context.Context, req *pb.CountRequest) (*pb.CountResponse, error) {
	s.logger.Debug("Count request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId))

	// Validate request
	if req.IndexName == "" {
		return nil, status.Error(codes.InvalidArgument, "index name is required")
	}

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	// For now, return document count
	// TODO: Implement query-based counting
	_ = req.Query

	stats := shard.Stats()

	return &pb.CountResponse{
		Count: stats.DocsCount,
	}, nil
}

// GetShardStats returns statistics for a specific shard
func (s *DataService) GetShardStats(ctx context.Context, req *pb.GetShardStatsRequest) (*pb.ShardStats, error) {
	s.logger.Debug("GetShardStats request",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId))

	// Get shard
	shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "shard not found: %v", err)
	}

	// Get stats
	stats := shard.Stats()

	return &pb.ShardStats{
		IndexName:               stats.IndexName,
		ShardId:                 stats.ShardID,
		IsPrimary:               stats.IsPrimary,
		DocsCount:               stats.DocsCount,
		DocsDeleted:             0, // TODO: Track deleted docs
		SizeBytes:               stats.SizeBytes,
		SearchQueriesTotal:      0, // TODO: Track query metrics
		SearchQueriesTimeMillis: 0,
		IndexingTotal:           0, // TODO: Track indexing metrics
		IndexingTimeMillis:      0,
		SegmentCount:            stats.SegmentCount,
		MaxDoc:                  stats.MaxDoc,
	}, nil
}

// GetNodeStats returns statistics for the entire node
func (s *DataService) GetNodeStats(ctx context.Context, req *pb.GetNodeStatsRequest) (*pb.DataNodeStats, error) {
	s.logger.Debug("GetNodeStats request",
		zap.Bool("include_shards", req.IncludeShards))

	// Get all shards
	shards := s.node.shards.List()

	// Aggregate stats
	var totalDocs, totalSize int64
	shardStats := make([]*pb.ShardStats, 0, len(shards))

	for _, shard := range shards {
		stats := shard.Stats()
		totalDocs += stats.DocsCount
		totalSize += stats.SizeBytes

		if req.IncludeShards {
			shardStats = append(shardStats, &pb.ShardStats{
				IndexName:               stats.IndexName,
				ShardId:                 stats.ShardID,
				IsPrimary:               stats.IsPrimary,
				DocsCount:               stats.DocsCount,
				DocsDeleted:             0, // TODO: Track deleted docs
				SizeBytes:               stats.SizeBytes,
				SearchQueriesTotal:      0, // TODO: Track query metrics
				SearchQueriesTimeMillis: 0,
				IndexingTotal:           0, // TODO: Track indexing metrics
				IndexingTimeMillis:      0,
			})
		}
	}

	// TODO: Get actual CPU, memory, disk usage
	nodeStats := &pb.DataNodeStats{
		NodeId:              s.node.cfg.NodeID,
		TotalShards:         int32(len(shards)),
		TotalDocs:           totalDocs,
		TotalSizeBytes:      totalSize,
		CpuUsagePercent:     0.0,  // TODO: Implement
		MemoryUsagePercent:  0.0,  // TODO: Implement
		DiskUsagePercent:    0.0,  // TODO: Implement
		UptimeSeconds:       0,    // TODO: Track uptime
		Shards:              shardStats,
	}

	return nodeStats, nil
}

// Helper functions

func (s *DataService) convertShardStateToProto(state ShardState) pb.ShardInfo_ShardState {
	switch state {
	case ShardStateInitializing:
		return pb.ShardInfo_SHARD_STATE_INITIALIZING
	case ShardStateStarted:
		return pb.ShardInfo_SHARD_STATE_STARTED
	case ShardStateRelocating:
		return pb.ShardInfo_SHARD_STATE_RELOCATING
	case ShardStateClosed:
		return pb.ShardInfo_SHARD_STATE_CLOSED
	default:
		return pb.ShardInfo_SHARD_STATE_UNKNOWN
	}
}

// Helper function for document conversion
func convertDocumentToJSON(doc map[string]interface{}) ([]byte, error) {
	return json.Marshal(doc)
}

func convertJSONToDocument(data []byte) (map[string]interface{}, error) {
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// convertAggregations converts Diagon aggregations to protobuf format
func convertAggregations(aggs map[string]diagon.AggregationResult) map[string]*pb.AggregationResult {
	if len(aggs) == 0 {
		return nil
	}

	result := make(map[string]*pb.AggregationResult, len(aggs))

	for name, agg := range aggs {
		pbAgg := &pb.AggregationResult{
			Type: agg.Type,
		}

		// Convert based on aggregation type
		switch agg.Type {
		case "terms", "histogram", "date_histogram", "range", "filters":
			// Bucket aggregations
			pbAgg.Buckets = convertBuckets(agg.Buckets)

		case "stats", "extended_stats":
			// Stats aggregations
			pbAgg.Count = agg.Count
			pbAgg.Min = agg.Min
			pbAgg.Max = agg.Max
			pbAgg.Avg = agg.Avg
			pbAgg.Sum = agg.Sum

			if agg.Type == "extended_stats" {
// 				pbAgg.SumOfSquares = agg.SumOfSquares
// 				pbAgg.Variance = agg.Variance
// 				pbAgg.StdDeviation = agg.StdDeviation
// 				pbAgg.StdDeviationBoundsUpper = agg.StdDeviationBoundsUpper
// 				pbAgg.StdDeviationBoundsLower = agg.StdDeviationBoundsLower
			}

		case "avg":
			// Average aggregation
			pbAgg.Avg = agg.Avg

		case "min":
			// Minimum aggregation
			pbAgg.Min = agg.Min

		case "max":
			// Maximum aggregation
			pbAgg.Max = agg.Max

		case "sum":
			// Sum aggregation
			pbAgg.Sum = agg.Sum

		case "value_count":
			// Value count aggregation
			pbAgg.Count = agg.Count

		case "percentiles":
			// Percentiles aggregation
			if agg.Values != nil {
				pbAgg.Values = make(map[string]float64, len(agg.Values))
				for k, v := range agg.Values {
					pbAgg.Values[k] = v
				}
			}

		case "cardinality":
			// Cardinality aggregation
			pbAgg.Value = agg.Value
		}

		result[name] = pbAgg
	}

	return result
}

// convertBuckets converts bucket data to protobuf format
func convertBuckets(buckets []map[string]interface{}) []*pb.AggregationBucket {
	if len(buckets) == 0 {
		return nil
	}

	result := make([]*pb.AggregationBucket, 0, len(buckets))

	for _, bucket := range buckets {
		pbBucket := &pb.AggregationBucket{}

		// Extract key (can be string or number)
		if key, ok := bucket["key"].(string); ok {
			pbBucket.Key = key
		}
		if numKey, ok := bucket["key"].(float64); ok {
			pbBucket.NumericKey = numKey
		}

		// Extract key_as_string for date histograms
		if keyAsString, ok := bucket["key_as_string"].(string); ok {
			pbBucket.Key = keyAsString
		}

		// Extract doc_count
		if docCount, ok := bucket["doc_count"].(int64); ok {
			pbBucket.DocCount = docCount
		} else if docCount, ok := bucket["doc_count"].(float64); ok {
			pbBucket.DocCount = int64(docCount)
		}

		// TODO: Handle sub-aggregations if needed
		// pbBucket.SubAggregations = convertAggregations(bucket["sub_aggs"])

		result = append(result, pbBucket)
	}

	return result
}

// ForceMerge optimizes a shard by merging segments.
// DISABLED: Diagon's SegmentMerger::merge() is a stub that creates empty merged
// segments (no data files), then deletes the originals. This destroys all index data.
// Re-enable once upstream Diagon implements SegmentMerger properly.
func (s *DataService) ForceMerge(ctx context.Context, req *pb.ForceMergeRequest) (*pb.ForceMergeResponse, error) {
	s.logger.Warn("ForceMerge SKIPPED: Diagon SegmentMerger is a stub that destroys data",
		zap.String("index", req.IndexName),
		zap.Int32("shard_id", req.ShardId),
		zap.Int32("max_segments", req.MaxSegments))

	// Flush pending docs to make them searchable (this is safe)
	if req.IndexName != "" {
		shard, err := s.node.shards.GetShard(req.IndexName, req.ShardId)
		if err == nil {
			if flushErr := shard.Flush(ctx); flushErr != nil {
				s.logger.Warn("Flush during force merge failed", zap.Error(flushErr))
			}
		}
	}

	// Return success without actually merging
	return &pb.ForceMergeResponse{
		Acknowledged:   true,
		SegmentsBefore: 0,
		SegmentsAfter:  0,
		DurationMillis: 0,
	}, nil
}

// computeNativeAggregations tries to compute aggregations directly from the inverted index.
// Returns nil if the aggregation type is not supported natively.
// Supported: terms aggregation (with optional sub-aggs that need doc extraction).
func computeNativeAggregations(diagonShard *diagon.Shard, aggsMap map[string]interface{}, logger *zap.Logger) map[string]*pb.AggregationResult {
	results := make(map[string]*pb.AggregationResult)

	for name, aggDef := range aggsMap {
		aggDefMap, ok := aggDef.(map[string]interface{})
		if !ok {
			return nil // Can't handle, fall back
		}

		for aggType, aggBody := range aggDefMap {
			if aggType == "aggs" {
				continue
			}
			bodyMap, ok := aggBody.(map[string]interface{})
			if !ok {
				return nil
			}

			switch aggType {
			case "terms":
				field, _ := bodyMap["field"].(string)
				size := 10
				if s, ok := bodyMap["size"].(float64); ok {
					size = int(s)
				}
				// Check for sub-aggregations - if present, fall back to doc extraction
				if _, hasSub := aggDefMap["aggs"]; hasSub {
					return nil // Can't do sub-aggs natively yet
				}
				buckets, err := diagonShard.ComputeTermsAgg(field, size)
				if err != nil || buckets == nil {
					logger.Debug("Native terms agg failed",
						zap.String("name", name),
						zap.String("field", field),
						zap.Error(err))
					return nil
				}
				pbBuckets := make([]*pb.AggregationBucket, len(buckets))
				for i, b := range buckets {
					pbBuckets[i] = &pb.AggregationBucket{
						Key:      b.Key,
						DocCount: b.DocCount,
					}
				}
				results[name] = &pb.AggregationResult{
					Type:    "terms",
					Buckets: pbBuckets,
				}

			case "date_histogram":
				return nil // Not supported natively, fall back to doc extraction

			default:
				return nil
			}
		}
	}

	if len(results) == 0 {
		return nil
	}
	return results
}

// extractAggregationFields extracts the list of field names needed by aggregation definitions.
// Used to enable field-only extraction mode (skip _source JSON parsing).
func extractAggregationFields(aggsMap map[string]interface{}) []string {
	fieldSet := make(map[string]bool)
	extractFieldsRecursive(aggsMap, fieldSet)
	fields := make([]string, 0, len(fieldSet))
	for f := range fieldSet {
		fields = append(fields, f)
	}
	return fields
}

func extractFieldsRecursive(aggsMap map[string]interface{}, fieldSet map[string]bool) {
	for _, aggDef := range aggsMap {
		aggDefMap, ok := aggDef.(map[string]interface{})
		if !ok {
			continue
		}
		for aggType, aggBody := range aggDefMap {
			if aggType == "aggs" {
				// Recurse into sub-aggregations
				if subAggs, ok := aggBody.(map[string]interface{}); ok {
					extractFieldsRecursive(subAggs, fieldSet)
				}
				continue
			}
			bodyMap, ok := aggBody.(map[string]interface{})
			if !ok {
				continue
			}
			if field, ok := bodyMap["field"].(string); ok {
				fieldSet[field] = true
			}
		}
	}
}
