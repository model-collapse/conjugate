// Copyright 2024 CONJUGATE Project
// Licensed under the Apache License, Version 2.0

package data

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/conjugate/conjugate/pkg/common/config"
	"github.com/conjugate/conjugate/pkg/data/diagon"
	"github.com/conjugate/conjugate/pkg/wasm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBatchIndexing_SeparateIntervals(t *testing.T) {
	// Create temp directory
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("conjugate-batch-test-%d", time.Now().UnixNano()))
	defer os.RemoveAll(tempDir)

	// Create logger
	logger, _ := zap.NewDevelopment()

	// Create Diagon bridge
	cfg := &config.DataNodeConfig{
		NodeID:    "test-node",
		DataDir:   tempDir,
		MaxShards: 10,
	}
	diagonBridge, err := diagon.NewDiagonBridge(&diagon.Config{
		DataDir:     tempDir,
		SIMDEnabled: false,
		Logger:      logger,
	})
	require.NoError(t, err)

	// Create shard manager
	shardManager := NewShardManager(cfg, logger, diagonBridge, &wasm.UDFRegistry{})
	err = shardManager.Start(context.Background())
	require.NoError(t, err)
	defer shardManager.Stop(context.Background())

	// Create shard
	err = shardManager.CreateShard(context.Background(), "batch-test", 0, true)
	require.NoError(t, err)

	// Get shard
	shard, err := shardManager.GetShard("batch-test", 0)
	require.NoError(t, err)

	t.Run("DefaultConfiguration", func(t *testing.T) {
		stats := shard.GetBatchStats()
		assert.Equal(t, int64(1000), stats["commit_batch_size"])
		assert.Equal(t, int64(1000), stats["commit_interval_ms"])
		assert.Equal(t, int64(1000), stats["refresh_interval_ms"])
	})

	t.Run("ConfigureSeparateIntervals", func(t *testing.T) {
		// Set different commit and refresh intervals
		shard.SetBatchConfig(500, 200*time.Millisecond, 1*time.Second)

		stats := shard.GetBatchStats()
		assert.Equal(t, int64(500), stats["commit_batch_size"])
		assert.Equal(t, int64(200), stats["commit_interval_ms"])
		assert.Equal(t, int64(1000), stats["refresh_interval_ms"])
	})

	t.Run("IndexAndCheckCommitRefreshSeparately", func(t *testing.T) {
		// Configure: commit 100ms, refresh 500ms
		shard.SetBatchConfig(100, 100*time.Millisecond, 500*time.Millisecond)

		// Index 50 documents (less than batch size)
		ctx := context.Background()
		for i := 0; i < 50; i++ {
			doc := map[string]interface{}{
				"id":    fmt.Sprintf("doc-%d", i),
				"field": "value",
			}
			err := shard.IndexDocument(ctx, fmt.Sprintf("doc-%d", i), doc)
			require.NoError(t, err)
		}

		// Wait 150ms - commit should have happened, but not refresh
		time.Sleep(150 * time.Millisecond)

		stats := shard.GetBatchStats()
		// Commit should have happened (pending_docs = 0)
		assert.Equal(t, int64(0), stats["pending_docs"])
		assert.False(t, stats["needs_commit"].(bool))

		// But refresh should be pending
		assert.True(t, stats["needs_refresh"].(bool))

		// Wait 400ms more - refresh should have happened
		time.Sleep(400 * time.Millisecond)

		stats = shard.GetBatchStats()
		assert.False(t, stats["needs_refresh"].(bool))
	})

	t.Run("FlushCommitsAndRefreshesBoth", func(t *testing.T) {
		// Configure: long intervals
		shard.SetBatchConfig(1000, 10*time.Second, 10*time.Second)

		// Index documents
		ctx := context.Background()
		for i := 0; i < 10; i++ {
			doc := map[string]interface{}{
				"id":    fmt.Sprintf("flush-doc-%d", i),
				"field": "value",
			}
			err := shard.IndexDocument(ctx, fmt.Sprintf("flush-doc-%d", i), doc)
			require.NoError(t, err)
		}

		// Check pending
		stats := shard.GetBatchStats()
		assert.Greater(t, stats["pending_docs"].(int), 0)
		assert.True(t, stats["needs_commit"].(bool))

		// Flush should commit and refresh immediately
		err = shard.Flush(ctx)
		require.NoError(t, err)

		// Check both cleared
		stats = shard.GetBatchStats()
		assert.Equal(t, int64(0), stats["pending_docs"])
		assert.False(t, stats["needs_commit"].(bool))
		assert.False(t, stats["needs_refresh"].(bool))
	})

	t.Run("HighThroughputProfile", func(t *testing.T) {
		// Configure for high throughput: commit 1s, refresh 5s
		shard.SetBatchConfig(5000, 1*time.Second, 5*time.Second)

		stats := shard.GetBatchStats()
		assert.Equal(t, int64(5000), stats["commit_batch_size"])
		assert.Equal(t, int64(1000), stats["commit_interval_ms"])
		assert.Equal(t, int64(5000), stats["refresh_interval_ms"])
	})

	t.Run("NearRealTimeProfile", func(t *testing.T) {
		// Configure for near real-time: commit 100ms, refresh 100ms
		shard.SetBatchConfig(100, 100*time.Millisecond, 100*time.Millisecond)

		stats := shard.GetBatchStats()
		assert.Equal(t, int64(100), stats["commit_batch_size"])
		assert.Equal(t, int64(100), stats["commit_interval_ms"])
		assert.Equal(t, int64(100), stats["refresh_interval_ms"])
	})
}

func TestBatchIndexing_BackgroundWorkers(t *testing.T) {
	// Create temp directory
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("conjugate-worker-test-%d", time.Now().UnixNano()))
	defer os.RemoveAll(tempDir)

	// Create components
	logger, _ := zap.NewDevelopment()
	cfg := &config.DataNodeConfig{
		NodeID:    "test-node",
		DataDir:   tempDir,
		MaxShards: 10,
	}
	diagonBridge, err := diagon.NewDiagonBridge(&diagon.Config{
		DataDir:     tempDir,
		SIMDEnabled: false,
		Logger:      logger,
	})
	require.NoError(t, err)

	shardManager := NewShardManager(cfg, logger, diagonBridge, &wasm.UDFRegistry{})
	err = shardManager.Start(context.Background())
	require.NoError(t, err)
	defer shardManager.Stop(context.Background())

	// Create shard
	err = shardManager.CreateShard(context.Background(), "worker-test", 0, true)
	require.NoError(t, err)

	shard, err := shardManager.GetShard("worker-test", 0)
	require.NoError(t, err)

	t.Run("BackgroundCommitterWorks", func(t *testing.T) {
		// Configure: commit 200ms, refresh 5s
		shard.SetBatchConfig(1000, 200*time.Millisecond, 5*time.Second)

		// Index 10 documents (less than batch size)
		ctx := context.Background()
		for i := 0; i < 10; i++ {
			doc := map[string]interface{}{
				"id":    fmt.Sprintf("bgc-doc-%d", i),
				"field": "value",
			}
			err := shard.IndexDocument(ctx, fmt.Sprintf("bgc-doc-%d", i), doc)
			require.NoError(t, err)
		}

		// Initially should have pending docs
		stats := shard.GetBatchStats()
		assert.Greater(t, stats["pending_docs"].(int), 0)

		// Wait for background committer (300ms > 200ms interval)
		time.Sleep(300 * time.Millisecond)

		// Should be committed now
		stats = shard.GetBatchStats()
		assert.Equal(t, int64(0), stats["pending_docs"])
		assert.False(t, stats["needs_commit"].(bool))
		assert.True(t, stats["needs_refresh"].(bool)) // But not refreshed yet
	})

	t.Run("BackgroundRefresherWorks", func(t *testing.T) {
		// Configure: commit 100ms, refresh 300ms
		shard.SetBatchConfig(1000, 100*time.Millisecond, 300*time.Millisecond)

		// Index documents
		ctx := context.Background()
		for i := 0; i < 10; i++ {
			doc := map[string]interface{}{
				"id":    fmt.Sprintf("bgr-doc-%d", i),
				"field": "value",
			}
			err := shard.IndexDocument(ctx, fmt.Sprintf("bgr-doc-%d", i), doc)
			require.NoError(t, err)
		}

		// Wait for commit (150ms > 100ms)
		time.Sleep(150 * time.Millisecond)

		stats := shard.GetBatchStats()
		assert.True(t, stats["needs_refresh"].(bool))

		// Wait for refresh (400ms total > 300ms interval)
		time.Sleep(250 * time.Millisecond)

		stats = shard.GetBatchStats()
		assert.False(t, stats["needs_refresh"].(bool))
	})
}
