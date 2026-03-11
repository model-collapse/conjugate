package data

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// AutoMergeScheduler manages automatic segment merging in the background
type AutoMergeScheduler struct {
	node   *DataNode
	logger *zap.Logger

	// Configuration
	enabled           bool
	checkInterval     time.Duration
	segmentThreshold  int32 // Trigger merge if segments > threshold
	targetSegments    int   // Target number of segments after merge
	minDocsForMerge   int64 // Minimum documents to consider merging
	maxConcurrentJobs int   // Maximum concurrent merge operations

	// State
	mu              sync.RWMutex
	running         bool
	stopCh          chan struct{}
	wg              sync.WaitGroup
	activeMerges    map[string]bool // shardKey -> is merging
	lastMergeTime   map[string]time.Time
	mergeMetrics    *MergeMetrics
}

// MergeMetrics tracks merge operation statistics
type MergeMetrics struct {
	mu                    sync.RWMutex
	TotalMerges           int64
	SuccessfulMerges      int64
	FailedMerges          int64
	TotalDuration         time.Duration
	TotalSegmentsReduced  int64
}

// AutoMergeConfig configures the auto-merge scheduler
type AutoMergeConfig struct {
	Enabled           bool
	CheckInterval     time.Duration
	SegmentThreshold  int32
	TargetSegments    int
	MinDocsForMerge   int64
	MaxConcurrentJobs int
}

// DefaultAutoMergeConfig returns default configuration
func DefaultAutoMergeConfig() *AutoMergeConfig {
	return &AutoMergeConfig{
		Enabled:           true,
		CheckInterval:     5 * time.Minute,  // Check every 5 minutes
		SegmentThreshold:  10,                // Merge if > 10 segments
		TargetSegments:    3,                 // Merge down to 3 segments
		MinDocsForMerge:   1000,              // Only merge shards with >= 1000 docs
		MaxConcurrentJobs: 2,                 // Max 2 concurrent merges
	}
}

// NewAutoMergeScheduler creates a new auto-merge scheduler
func NewAutoMergeScheduler(node *DataNode, config *AutoMergeConfig, logger *zap.Logger) *AutoMergeScheduler {
	if config == nil {
		config = DefaultAutoMergeConfig()
	}

	return &AutoMergeScheduler{
		node:              node,
		logger:            logger,
		enabled:           config.Enabled,
		checkInterval:     config.CheckInterval,
		segmentThreshold:  config.SegmentThreshold,
		targetSegments:    config.TargetSegments,
		minDocsForMerge:   config.MinDocsForMerge,
		maxConcurrentJobs: config.MaxConcurrentJobs,
		stopCh:            make(chan struct{}),
		activeMerges:      make(map[string]bool),
		lastMergeTime:     make(map[string]time.Time),
		mergeMetrics:      &MergeMetrics{},
	}
}

// Start starts the auto-merge scheduler
func (s *AutoMergeScheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		s.logger.Info("Auto-merge scheduler disabled")
		return nil
	}

	if s.running {
		return nil
	}

	s.running = true
	s.wg.Add(1)

	go s.run(ctx)

	s.logger.Info("Auto-merge scheduler started",
		zap.Duration("check_interval", s.checkInterval),
		zap.Int32("segment_threshold", s.segmentThreshold),
		zap.Int("target_segments", s.targetSegments),
		zap.Int64("min_docs", s.minDocsForMerge),
		zap.Int("max_concurrent", s.maxConcurrentJobs))

	return nil
}

// Stop stops the auto-merge scheduler
func (s *AutoMergeScheduler) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}

	s.logger.Info("Stopping auto-merge scheduler...")
	close(s.stopCh)
	s.mu.Unlock()

	// Wait for background goroutine to finish
	s.wg.Wait()

	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	s.logger.Info("Auto-merge scheduler stopped")
	return nil
}

// run is the main loop that periodically checks and merges shards
func (s *AutoMergeScheduler) run(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	s.logger.Info("Auto-merge scheduler running")

	for {
		select {
		case <-s.stopCh:
			s.logger.Info("Auto-merge scheduler received stop signal")
			return
		case <-ctx.Done():
			s.logger.Info("Auto-merge scheduler context cancelled")
			return
		case <-ticker.C:
			s.checkAndMergeShards(ctx)
		}
	}
}

// checkAndMergeShards checks all shards and triggers merge if needed
func (s *AutoMergeScheduler) checkAndMergeShards(ctx context.Context) {
	s.logger.Debug("Checking shards for auto-merge")

	// Get all shards
	shards := s.node.shards.List()

	var candidateCount int
	var mergedCount int

	for _, shard := range shards {
		// Skip if already merging
		shardKey := shardKey(shard.IndexName, shard.ShardID)

		s.mu.RLock()
		isMerging := s.activeMerges[shardKey]
		lastMerge := s.lastMergeTime[shardKey]
		activeMergeCount := len(s.activeMerges)
		s.mu.RUnlock()

		if isMerging {
			continue
		}

		// Rate limit: Don't merge same shard within 30 minutes
		if time.Since(lastMerge) < 30*time.Minute {
			continue
		}

		// Check if max concurrent jobs reached
		if activeMergeCount >= s.maxConcurrentJobs {
			s.logger.Debug("Max concurrent merge jobs reached",
				zap.Int("active", activeMergeCount),
				zap.Int("max", s.maxConcurrentJobs))
			break
		}

		// Get shard stats
		stats := shard.Stats()

		// Check if merge is needed
		needsMerge := s.shouldMerge(stats)
		if needsMerge {
			candidateCount++

			// Trigger merge in background
			s.mu.Lock()
			s.activeMerges[shardKey] = true
			s.mu.Unlock()

			mergedCount++

			go s.mergeShard(ctx, shard, shardKey, stats)
		}
	}

	if candidateCount > 0 {
		s.logger.Info("Auto-merge check completed",
			zap.Int("candidates", candidateCount),
			zap.Int("merged", mergedCount),
			zap.Int("total_shards", len(shards)))
	}
}

// shouldMerge determines if a shard needs merging
func (s *AutoMergeScheduler) shouldMerge(stats *ShardStats) bool {
	// Must have minimum documents
	if stats.DocsCount < s.minDocsForMerge {
		return false
	}

	// Must exceed segment threshold
	if stats.SegmentCount <= s.segmentThreshold {
		return false
	}

	// Shard must be in started state
	if stats.State != ShardStateStarted {
		return false
	}

	return true
}

// mergeShard performs the merge operation for a single shard
func (s *AutoMergeScheduler) mergeShard(ctx context.Context, shard *Shard, shardKey string, stats *ShardStats) {
	defer func() {
		s.mu.Lock()
		delete(s.activeMerges, shardKey)
		s.lastMergeTime[shardKey] = time.Now()
		s.mu.Unlock()
	}()

	s.logger.Info("Starting auto-merge",
		zap.String("index", stats.IndexName),
		zap.Int32("shard_id", stats.ShardID),
		zap.Int32("segments", stats.SegmentCount),
		zap.Int("target", s.targetSegments),
		zap.Int64("docs", stats.DocsCount))

	startTime := time.Now()
	segmentsBefore := stats.SegmentCount

	// Perform flush before merge
	if err := shard.Flush(ctx); err != nil {
		s.logger.Warn("Failed to flush before merge",
			zap.String("index", stats.IndexName),
			zap.Int32("shard_id", stats.ShardID),
			zap.Error(err))
		s.recordMergeFailure()
		return
	}

	// Perform merge
	err := shard.ForceMerge(s.targetSegments)
	duration := time.Since(startTime)

	if err != nil {
		s.logger.Error("Auto-merge failed",
			zap.String("index", stats.IndexName),
			zap.Int32("shard_id", stats.ShardID),
			zap.Error(err),
			zap.Duration("duration", duration))
		s.recordMergeFailure()
		return
	}

	// Get stats after merge
	statsAfter := shard.Stats()
	segmentsAfter := statsAfter.SegmentCount
	segmentsReduced := segmentsBefore - segmentsAfter

	s.logger.Info("Auto-merge completed successfully",
		zap.String("index", stats.IndexName),
		zap.Int32("shard_id", stats.ShardID),
		zap.Int32("segments_before", segmentsBefore),
		zap.Int32("segments_after", segmentsAfter),
		zap.Int32("reduced", segmentsReduced),
		zap.Duration("duration", duration))

	s.recordMergeSuccess(duration, int64(segmentsReduced))
}

// recordMergeSuccess updates metrics for successful merge
func (s *AutoMergeScheduler) recordMergeSuccess(duration time.Duration, segmentsReduced int64) {
	s.mergeMetrics.mu.Lock()
	defer s.mergeMetrics.mu.Unlock()

	s.mergeMetrics.TotalMerges++
	s.mergeMetrics.SuccessfulMerges++
	s.mergeMetrics.TotalDuration += duration
	s.mergeMetrics.TotalSegmentsReduced += segmentsReduced
}

// recordMergeFailure updates metrics for failed merge
func (s *AutoMergeScheduler) recordMergeFailure() {
	s.mergeMetrics.mu.Lock()
	defer s.mergeMetrics.mu.Unlock()

	s.mergeMetrics.TotalMerges++
	s.mergeMetrics.FailedMerges++
}

// GetMetrics returns current merge metrics
func (s *AutoMergeScheduler) GetMetrics() map[string]interface{} {
	s.mergeMetrics.mu.RLock()
	defer s.mergeMetrics.mu.RUnlock()

	avgDuration := time.Duration(0)
	if s.mergeMetrics.SuccessfulMerges > 0 {
		avgDuration = s.mergeMetrics.TotalDuration / time.Duration(s.mergeMetrics.SuccessfulMerges)
	}

	s.mu.RLock()
	activeMerges := len(s.activeMerges)
	s.mu.RUnlock()

	return map[string]interface{}{
		"enabled":               s.enabled,
		"running":               s.running,
		"total_merges":          s.mergeMetrics.TotalMerges,
		"successful_merges":     s.mergeMetrics.SuccessfulMerges,
		"failed_merges":         s.mergeMetrics.FailedMerges,
		"total_segments_reduced": s.mergeMetrics.TotalSegmentsReduced,
		"avg_duration_ms":       avgDuration.Milliseconds(),
		"active_merges":         activeMerges,
		"check_interval_sec":    s.checkInterval.Seconds(),
		"segment_threshold":     s.segmentThreshold,
		"target_segments":       s.targetSegments,
	}
}

// UpdateConfig updates the scheduler configuration at runtime
func (s *AutoMergeScheduler) UpdateConfig(config *AutoMergeConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if config.CheckInterval > 0 {
		s.checkInterval = config.CheckInterval
	}
	if config.SegmentThreshold > 0 {
		s.segmentThreshold = config.SegmentThreshold
	}
	if config.TargetSegments > 0 {
		s.targetSegments = config.TargetSegments
	}
	if config.MinDocsForMerge > 0 {
		s.minDocsForMerge = config.MinDocsForMerge
	}
	if config.MaxConcurrentJobs > 0 {
		s.maxConcurrentJobs = config.MaxConcurrentJobs
	}

	s.logger.Info("Auto-merge configuration updated",
		zap.Duration("check_interval", s.checkInterval),
		zap.Int32("segment_threshold", s.segmentThreshold),
		zap.Int("target_segments", s.targetSegments))
}
