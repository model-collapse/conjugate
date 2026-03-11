package coordination

import (
	"context"
	"sync"
	"time"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
)

// cachingMasterClient wraps a masterClientInterface with a local cache for
// GetShardRouting. This avoids repeated gRPC round-trips to the master node
// for the same index within a short window. Both QueryService and QueryExecutor
// share the same master client, so this eliminates duplicate routing lookups
// that previously added ~2-4ms per query.
type cachingMasterClient struct {
	inner masterClientInterface

	mu    sync.RWMutex
	cache map[string]*routingCacheEntry
	ttl   time.Duration
}

type routingCacheEntry struct {
	routing   map[int32]*pb.ShardRouting
	timestamp time.Time
}

// newCachingMasterClient wraps an existing master client with routing cache.
func newCachingMasterClient(inner masterClientInterface, ttl time.Duration) *cachingMasterClient {
	return &cachingMasterClient{
		inner: inner,
		cache: make(map[string]*routingCacheEntry),
		ttl:   ttl,
	}
}

func (c *cachingMasterClient) GetShardRouting(ctx context.Context, indexName string) (map[int32]*pb.ShardRouting, error) {
	now := time.Now()

	// Fast path: read lock
	c.mu.RLock()
	if entry, ok := c.cache[indexName]; ok && now.Sub(entry.timestamp) < c.ttl {
		c.mu.RUnlock()
		return entry.routing, nil
	}
	c.mu.RUnlock()

	// Slow path: fetch from master and cache
	routing, err := c.inner.GetShardRouting(ctx, indexName)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[indexName] = &routingCacheEntry{
		routing:   routing,
		timestamp: now,
	}
	c.mu.Unlock()

	return routing, nil
}

func (c *cachingMasterClient) GetIndexMetadata(ctx context.Context, indexName string) (*pb.IndexMetadataResponse, error) {
	return c.inner.GetIndexMetadata(ctx, indexName)
}
