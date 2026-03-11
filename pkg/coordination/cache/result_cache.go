package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/conjugate/conjugate/pkg/coordination/executor"
	"github.com/conjugate/conjugate/pkg/coordination/parser"
	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics for result cache
var (
	resultCacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "conjugate_result_cache_hits_total",
			Help: "Total number of result cache hits",
		},
		[]string{"index"},
	)

	resultCacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "conjugate_result_cache_misses_total",
			Help: "Total number of result cache misses",
		},
		[]string{"index"},
	)

	resultCacheEvictions = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "conjugate_result_cache_evictions_total",
			Help: "Total number of result cache evictions",
		},
	)

	resultCacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "conjugate_result_cache_entries",
			Help: "Current number of entries in result cache",
		},
	)

	resultCacheHitRate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "conjugate_result_cache_hit_rate",
			Help: "Result cache hit rate (hits / total requests)",
		},
	)

	resultCacheLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "conjugate_result_cache_lookup_seconds",
			Help:    "Result cache lookup latency in seconds",
			Buckets: prometheus.ExponentialBuckets(0.000001, 2, 12), // 1μs to ~4ms
		},
	)
)

// ResultCache provides high-performance caching of query results using Ristretto
type ResultCache struct {
	cache      *ristretto.Cache
	ttl        time.Duration
	enabled    bool
	hits       uint64
	misses     uint64
	maxCost    int64 // Maximum cache size in bytes
	defaultCost int64 // Default cost per entry if size cannot be estimated
}

// ResultCacheConfig configures the result cache
type ResultCacheConfig struct {
	// Cache size
	MaxCost     int64         // Max size in bytes (default: 1GB)
	NumCounters int64         // Number of keys to track frequency (default: 100M)
	MaxEntries  int64         // Max number of cache entries (default: 10000)

	// TTL
	TTL time.Duration // Time-to-live for cached results (default: 2 minutes)

	// Feature flag
	Enabled bool

	// Cost estimation
	DefaultCost int64 // Default cost per entry in bytes (default: 10KB)
}

// DefaultResultCacheConfig returns default result cache configuration
func DefaultResultCacheConfig() *ResultCacheConfig {
	return &ResultCacheConfig{
		MaxCost:     1024 * 1024 * 1024, // 1 GB
		NumCounters: 100_000_000,        // 100M counters
		MaxEntries:  10000,               // 10K entries max
		TTL:         2 * time.Minute,    // 2 minute TTL
		Enabled:     true,
		DefaultCost: 10 * 1024,          // 10 KB default
	}
}

// NewResultCache creates a new result cache using Ristretto
func NewResultCache(config *ResultCacheConfig) (*ResultCache, error) {
	if config == nil {
		config = DefaultResultCacheConfig()
	}

	if !config.Enabled {
		return &ResultCache{
			enabled: false,
		}, nil
	}

	// Create Ristretto cache
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: config.NumCounters,
		MaxCost:     config.MaxCost,
		BufferItems: 64,
		Metrics:     true, // Enable metrics for monitoring
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ristretto cache: %w", err)
	}

	rc := &ResultCache{
		cache:       ristrettoCache,
		ttl:         config.TTL,
		enabled:     true,
		maxCost:     config.MaxCost,
		defaultCost: config.DefaultCost,
	}

	// Start background metrics reporter
	go rc.reportMetrics()

	return rc, nil
}

// CachedResult wraps a search result with metadata
type CachedResult struct {
	Result    *executor.SearchResult
	CachedAt  time.Time
	IndexName string
	QueryHash string
}

// Get retrieves a cached result for a query
func (rc *ResultCache) Get(indexName string, searchReq *parser.SearchRequest) (*executor.SearchResult, bool) {
	if !rc.enabled {
		return nil, false
	}

	start := time.Now()
	defer func() {
		resultCacheLatency.Observe(time.Since(start).Seconds())
	}()

	key := rc.generateKey(indexName, searchReq)

	value, found := rc.cache.Get(key)
	if !found {
		rc.misses++
		resultCacheMisses.WithLabelValues(indexName).Inc()
		rc.updateHitRate()
		return nil, false
	}

	cached, ok := value.(*CachedResult)
	if !ok {
		rc.misses++
		resultCacheMisses.WithLabelValues(indexName).Inc()
		rc.updateHitRate()
		return nil, false
	}

	// Check if entry has expired (Ristretto doesn't support TTL natively)
	if time.Since(cached.CachedAt) > rc.ttl {
		rc.cache.Del(key)
		rc.misses++
		resultCacheMisses.WithLabelValues(indexName).Inc()
		rc.updateHitRate()
		return nil, false
	}

	rc.hits++
	resultCacheHits.WithLabelValues(indexName).Inc()
	rc.updateHitRate()

	return cached.Result, true
}

// Put stores a search result in cache
func (rc *ResultCache) Put(indexName string, searchReq *parser.SearchRequest, result *executor.SearchResult) {
	if !rc.enabled {
		return
	}

	key := rc.generateKey(indexName, searchReq)

	cached := &CachedResult{
		Result:    result,
		CachedAt:  time.Now(),
		IndexName: indexName,
		QueryHash: key,
	}

	// Estimate cost (size in bytes)
	cost := rc.estimateCost(result)

	// Store with cost-based eviction
	// Ristretto uses cost for admission control and eviction
	rc.cache.Set(key, cached, cost)
}

// InvalidateIndex removes all cached results for a specific index
func (rc *ResultCache) InvalidateIndex(indexName string) {
	if !rc.enabled {
		return
	}

	// Ristretto doesn't support prefix-based deletion,
	// so we'd need to track keys per index separately
	// For now, we clear the entire cache (simple but effective)
	rc.cache.Clear()
}

// Clear removes all cached entries
func (rc *ResultCache) Clear() {
	if rc.enabled {
		rc.cache.Clear()
	}
}

// Close closes the cache and releases resources
func (rc *ResultCache) Close() {
	if rc.enabled {
		rc.cache.Close()
	}
}

// Stats returns cache statistics
func (rc *ResultCache) Stats() map[string]interface{} {
	if !rc.enabled {
		return map[string]interface{}{
			"enabled": false,
		}
	}

	metrics := rc.cache.Metrics

	return map[string]interface{}{
		"enabled":      true,
		"hits":         rc.hits,
		"misses":       rc.misses,
		"hit_rate":     rc.calculateHitRate(),
		"ratio":        metrics.Ratio(),
		"keys_added":   metrics.KeysAdded(),
		"keys_evicted": metrics.KeysEvicted(),
		"cost_added":   metrics.CostAdded(),
		"cost_evicted": metrics.CostEvicted(),
		"max_cost":     rc.maxCost,
	}
}

// generateKey creates a cache key for a search request
func (rc *ResultCache) generateKey(indexName string, searchReq *parser.SearchRequest) string {
	// Create a normalized representation including all query parameters
	keyData := struct {
		Index        string
		Query        interface{}
		Aggregations interface{}
		Size         int
		From         int
		Sort         interface{}
		Source       interface{}
		Highlight    interface{}
	}{
		Index:        indexName,
		Query:        normalizeQuery(searchReq.ParsedQuery),
		Aggregations: searchReq.Aggregations,
		Size:         searchReq.Size,
		From:         searchReq.From,
		Sort:         searchReq.Sort,
		Source:       searchReq.Source,
		Highlight:    searchReq.Highlight,
	}

	// Serialize to JSON for consistent hashing
	jsonData, err := json.Marshal(keyData)
	if err != nil {
		// Fallback to timestamp-based key (won't cache, but won't crash)
		return fmt.Sprintf("result:%s:%d", indexName, time.Now().UnixNano())
	}

	// Hash the JSON data using SHA-256
	hash := sha256.Sum256(jsonData)
	return "result:" + hex.EncodeToString(hash[:])
}

// estimateCost estimates the memory cost of a search result in bytes
func (rc *ResultCache) estimateCost(result *executor.SearchResult) int64 {
	if result == nil {
		return rc.defaultCost
	}

	// Rough estimation:
	// - Base overhead: 1KB
	// - Per hit: ~1KB base + field data
	// - Per aggregation: ~500 bytes

	cost := int64(1024) // Base overhead

	// Estimate hits cost
	if result.Hits != nil {
		for _, hit := range result.Hits {
			cost += 1024 // Base per hit

			// Estimate source field data
			if hit.Source != nil {
				// Rough JSON size estimation
				jsonBytes, _ := json.Marshal(hit.Source)
				cost += int64(len(jsonBytes))
			}
		}
	}

	// Estimate aggregations cost (simplified)
	if result.Aggregations != nil {
		aggBytes, _ := json.Marshal(result.Aggregations)
		cost += int64(len(aggBytes))
	}

	// Cap cost at reasonable maximum
	if cost > rc.maxCost/100 { // No single entry should be >1% of cache
		cost = rc.maxCost / 100
	}

	return cost
}

// calculateHitRate calculates the current hit rate
func (rc *ResultCache) calculateHitRate() float64 {
	total := rc.hits + rc.misses
	if total == 0 {
		return 0.0
	}
	return float64(rc.hits) / float64(total)
}

// updateHitRate updates the Prometheus hit rate metric
func (rc *ResultCache) updateHitRate() {
	resultCacheHitRate.Set(rc.calculateHitRate())
}

// reportMetrics periodically reports cache metrics
func (rc *ResultCache) reportMetrics() {
	if !rc.enabled {
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics := rc.cache.Metrics

		// Update Prometheus metrics
		resultCacheEvictions.Add(float64(metrics.KeysEvicted()))

		// Update cache size metric
		// Note: Ristretto doesn't expose current key count directly
		// We use KeysAdded - KeysEvicted as approximation
		currentKeys := metrics.KeysAdded() - metrics.KeysEvicted()
		resultCacheSize.Set(float64(currentKeys))
	}
}
