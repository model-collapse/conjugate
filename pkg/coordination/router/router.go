package router

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
	"go.uber.org/zap"
)

// DataNodeClient interface for communication with data nodes
type DataNodeClient interface {
	IndexDocument(ctx context.Context, indexName string, shardID int32, docID string, document map[string]interface{}) (*pb.IndexDocumentResponse, error)
	GetDocument(ctx context.Context, indexName string, shardID int32, docID string) (*pb.GetDocumentResponse, error)
	DeleteDocument(ctx context.Context, indexName string, shardID int32, docID string) (*pb.DeleteDocumentResponse, error)
	IsConnected() bool
	Connect(ctx context.Context) error
	NodeID() string
}

// MasterClient interface for getting cluster state
type MasterClient interface {
	GetShardRouting(ctx context.Context, indexName string) (map[int32]*pb.ShardRouting, error)
	GetIndexMetadata(ctx context.Context, indexName string) (*pb.IndexMetadataResponse, error)
}

// DocumentRouter routes document operations to the appropriate shards
type DocumentRouter struct {
	logger      *zap.Logger
	masterClient MasterClient
	dataClients map[string]DataNodeClient // nodeID -> client
}

// NewDocumentRouter creates a new document router
func NewDocumentRouter(masterClient MasterClient, dataClients map[string]DataNodeClient, logger *zap.Logger) *DocumentRouter {
	return &DocumentRouter{
		logger:      logger,
		masterClient: masterClient,
		dataClients: dataClients,
	}
}

// RouteIndexDocument routes an index document operation to the correct shard
func (dr *DocumentRouter) RouteIndexDocument(ctx context.Context, indexName, docID string, document map[string]interface{}) (*pb.IndexDocumentResponse, error) {
	// Get index metadata to determine number of shards
	metadata, err := dr.masterClient.GetIndexMetadata(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to get index metadata: %w", err)
	}

	numShards := metadata.Metadata.Settings.NumberOfShards
	if numShards == 0 {
		return nil, fmt.Errorf("index has no shards configured")
	}

	// Calculate which shard this document belongs to
	shardID := dr.calculateShardID(docID, numShards)

	// Get shard routing information
	routing, err := dr.masterClient.GetShardRouting(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to get shard routing: %w", err)
	}

	shard, exists := routing[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %d not found for index %s", shardID, indexName)
	}

	// Find primary shard for writes
	if shard.Allocation == nil || shard.Allocation.State != pb.ShardAllocation_SHARD_STATE_STARTED {
		return nil, fmt.Errorf("shard %d is not available (state: %v)", shardID, shard.Allocation.State)
	}

	// Only write to primary shard
	if !shard.IsPrimary {
		return nil, fmt.Errorf("shard %d is not a primary shard", shardID)
	}

	nodeID := shard.Allocation.NodeId
	if nodeID == "" {
		return nil, fmt.Errorf("shard %d has no node assignment", shardID)
	}

	// Get data node client
	client, exists := dr.dataClients[nodeID]
	if !exists {
		return nil, fmt.Errorf("data node %s not found", nodeID)
	}

	// Ensure client is connected
	if !client.IsConnected() {
		if err := client.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect to node %s: %w", nodeID, err)
		}
	}

	// Route to data node
	dr.logger.Debug("Routing index document",
		zap.String("index", indexName),
		zap.Int32("shard_id", shardID),
		zap.String("node_id", nodeID))

	resp, err := client.IndexDocument(ctx, indexName, shardID, docID, document)
	if err != nil {
		dr.logger.Error("IndexDocument call failed", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

// BulkDocItem represents a document to be indexed in a bulk operation
type BulkDocItem struct {
	DocID        string
	Document     map[string]interface{}
	DocumentJSON []byte // Raw JSON bytes for zero-copy pass-through
}

// BulkResultItem represents the result of indexing a single document
type BulkResultItem struct {
	DocID   string
	Success bool
	Error   string
	Version int64
}

// BulkIndexClient is the interface for data node clients that support bulk indexing
type BulkIndexClient interface {
	DataNodeClient
	BulkIndex(ctx context.Context, indexName string, shardID int32, items []*pb.BulkIndexItem) (*pb.BulkIndexResponse, error)
}

// RouteBulkIndex routes a batch of documents to the appropriate shards using a single
// metadata lookup and batched gRPC calls. This is dramatically faster than calling
// RouteIndexDocument for each document individually.
func (dr *DocumentRouter) RouteBulkIndex(ctx context.Context, indexName string, docs []BulkDocItem) ([]BulkResultItem, error) {
	results := make([]BulkResultItem, len(docs))

	// Get index metadata ONCE for the entire batch
	metadata, err := dr.masterClient.GetIndexMetadata(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to get index metadata: %w", err)
	}

	numShards := metadata.Metadata.Settings.NumberOfShards
	if numShards == 0 {
		return nil, fmt.Errorf("index has no shards configured")
	}

	// Get shard routing ONCE for the entire batch
	routing, err := dr.masterClient.GetShardRouting(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to get shard routing: %w", err)
	}

	// Group documents by shard
	type shardBatch struct {
		shardID  int32
		nodeID   string
		indices  []int // original indices in docs slice
		items    []*pb.BulkIndexItem
	}

	shardBatches := make(map[int32]*shardBatch)

	for i, doc := range docs {
		results[i].DocID = doc.DocID

		shardID := dr.calculateShardID(doc.DocID, numShards)

		shard, exists := routing[shardID]
		if !exists {
			results[i].Success = false
			results[i].Error = fmt.Sprintf("shard %d not found for index %s", shardID, indexName)
			continue
		}

		if shard.Allocation == nil || shard.Allocation.State != pb.ShardAllocation_SHARD_STATE_STARTED {
			results[i].Success = false
			results[i].Error = fmt.Sprintf("shard %d is not available", shardID)
			continue
		}

		if !shard.IsPrimary {
			results[i].Success = false
			results[i].Error = fmt.Sprintf("shard %d is not a primary shard", shardID)
			continue
		}

		nodeID := shard.Allocation.NodeId

		batch, exists := shardBatches[shardID]
		if !exists {
			batch = &shardBatch{
				shardID: shardID,
				nodeID:  nodeID,
			}
			shardBatches[shardID] = batch
		}

		// Use raw JSON bytes for zero-copy pass-through (skip structpb conversion)
		var jsonBytes []byte
		if len(doc.DocumentJSON) > 0 {
			jsonBytes = doc.DocumentJSON
		} else {
			// Fallback: marshal from map if raw bytes not available
			jsonBytes, err = json.Marshal(doc.Document)
			if err != nil {
				results[i].Success = false
				results[i].Error = fmt.Sprintf("failed to marshal document: %v", err)
				continue
			}
		}

		batch.indices = append(batch.indices, i)
		batch.items = append(batch.items, &pb.BulkIndexItem{
			DocId:        doc.DocID,
			DocumentJson: jsonBytes,
		})
	}

	// Send one BulkIndex RPC per shard
	for _, batch := range shardBatches {
		client, exists := dr.dataClients[batch.nodeID]
		if !exists {
			for _, idx := range batch.indices {
				results[idx].Success = false
				results[idx].Error = fmt.Sprintf("data node %s not found", batch.nodeID)
			}
			continue
		}

		// Ensure client is connected
		if !client.IsConnected() {
			if err := client.Connect(ctx); err != nil {
				for _, idx := range batch.indices {
					results[idx].Success = false
					results[idx].Error = fmt.Sprintf("failed to connect to node %s: %v", batch.nodeID, err)
				}
				continue
			}
		}

		// Try BulkIndex if client supports it, otherwise fall back to individual calls
		bulkClient, ok := client.(BulkIndexClient)
		if ok {
			resp, err := bulkClient.BulkIndex(ctx, indexName, batch.shardID, batch.items)
			if err != nil {
				for _, idx := range batch.indices {
					results[idx].Success = false
					results[idx].Error = err.Error()
				}
				continue
			}

			// Map responses back to original indices
			for j, itemResp := range resp.Items {
				if j < len(batch.indices) {
					idx := batch.indices[j]
					results[idx].Success = itemResp.Acknowledged
					results[idx].DocID = itemResp.DocId
					if itemResp.Error != "" {
						results[idx].Error = itemResp.Error
					}
				}
			}
		} else {
			// Fallback: individual IndexDocument calls
			for j, item := range batch.items {
				idx := batch.indices[j]
				docMap := docs[idx].Document
				resp, err := client.IndexDocument(ctx, indexName, batch.shardID, item.DocId, docMap)
				if err != nil {
					results[idx].Success = false
					results[idx].Error = err.Error()
				} else {
					results[idx].Success = true
					results[idx].Version = resp.Version
				}
			}
		}
	}

	dr.logger.Debug("Bulk index routed",
		zap.String("index", indexName),
		zap.Int("total_docs", len(docs)),
		zap.Int("shard_count", len(shardBatches)))

	return results, nil
}

// RouteGetDocument routes a get document operation to the correct shard
func (dr *DocumentRouter) RouteGetDocument(ctx context.Context, indexName, docID string) (*pb.GetDocumentResponse, error) {
	// Get index metadata to determine number of shards
	metadata, err := dr.masterClient.GetIndexMetadata(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to get index metadata: %w", err)
	}

	numShards := metadata.Metadata.Settings.NumberOfShards
	if numShards == 0 {
		return nil, fmt.Errorf("index has no shards configured")
	}

	// Calculate which shard this document belongs to
	shardID := dr.calculateShardID(docID, numShards)

	// Get shard routing information
	routing, err := dr.masterClient.GetShardRouting(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to get shard routing: %w", err)
	}

	shard, exists := routing[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %d not found for index %s", shardID, indexName)
	}

	// For reads, we can use primary or replica
	if shard.Allocation == nil || shard.Allocation.State != pb.ShardAllocation_SHARD_STATE_STARTED {
		return nil, fmt.Errorf("shard %d is not available", shardID)
	}

	nodeID := shard.Allocation.NodeId
	if nodeID == "" {
		return nil, fmt.Errorf("shard %d has no node assignment", shardID)
	}

	// Get data node client
	client, exists := dr.dataClients[nodeID]
	if !exists {
		return nil, fmt.Errorf("data node %s not found", nodeID)
	}

	// Ensure client is connected
	if !client.IsConnected() {
		if err := client.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect to node %s: %w", nodeID, err)
		}
	}

	// Route to data node
	dr.logger.Debug("Routing get document",
		zap.String("index", indexName),
		zap.String("doc_id", docID),
		zap.Int32("shard_id", shardID),
		zap.String("node_id", nodeID))

	return client.GetDocument(ctx, indexName, shardID, docID)
}

// RouteDeleteDocument routes a delete document operation to the correct shard
func (dr *DocumentRouter) RouteDeleteDocument(ctx context.Context, indexName, docID string) (*pb.DeleteDocumentResponse, error) {
	// Get index metadata to determine number of shards
	metadata, err := dr.masterClient.GetIndexMetadata(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to get index metadata: %w", err)
	}

	numShards := metadata.Metadata.Settings.NumberOfShards
	if numShards == 0 {
		return nil, fmt.Errorf("index has no shards configured")
	}

	// Calculate which shard this document belongs to
	shardID := dr.calculateShardID(docID, numShards)

	// Get shard routing information
	routing, err := dr.masterClient.GetShardRouting(ctx, indexName)
	if err != nil {
		return nil, fmt.Errorf("failed to get shard routing: %w", err)
	}

	shard, exists := routing[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %d not found for index %s", shardID, indexName)
	}

	// Only delete from primary shard
	if shard.Allocation == nil || shard.Allocation.State != pb.ShardAllocation_SHARD_STATE_STARTED {
		return nil, fmt.Errorf("shard %d is not available", shardID)
	}

	if !shard.IsPrimary {
		return nil, fmt.Errorf("shard %d is not a primary shard", shardID)
	}

	nodeID := shard.Allocation.NodeId
	if nodeID == "" {
		return nil, fmt.Errorf("shard %d has no node assignment", shardID)
	}

	// Get data node client
	client, exists := dr.dataClients[nodeID]
	if !exists {
		return nil, fmt.Errorf("data node %s not found", nodeID)
	}

	// Ensure client is connected
	if !client.IsConnected() {
		if err := client.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect to node %s: %w", nodeID, err)
		}
	}

	// Route to data node
	dr.logger.Debug("Routing delete document",
		zap.String("index", indexName),
		zap.String("doc_id", docID),
		zap.Int32("shard_id", shardID),
		zap.String("node_id", nodeID))

	return client.DeleteDocument(ctx, indexName, shardID, docID)
}

// calculateShardID uses consistent hashing to determine which shard a document belongs to
func (dr *DocumentRouter) calculateShardID(docID string, numShards int32) int32 {
	// Use FNV-1a hash (fast, good distribution)
	h := fnv.New32a()
	h.Write([]byte(docID))
	hash := h.Sum32()

	// Modulo to get shard ID
	return int32(hash % uint32(numShards))
}

// SetDataClients updates the data node clients
func (dr *DocumentRouter) SetDataClients(clients map[string]DataNodeClient) {
	dr.dataClients = clients
}
