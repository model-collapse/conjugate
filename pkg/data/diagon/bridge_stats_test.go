package diagon

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
)

func TestShardGetStats(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "diagon-stats-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	shardPath := filepath.Join(tmpDir, "shard0")
	if err := os.MkdirAll(shardPath, 0755); err != nil {
		t.Fatalf("Failed to create shard directory: %v", err)
	}

	// Create bridge
	logger, _ := zap.NewDevelopment()
	bridge, err := NewDiagonBridge(&Config{
		DataDir: tmpDir,
		Logger:  logger,
	})
	if err != nil {
		t.Fatalf("Failed to create bridge: %v", err)
	}

	// Start bridge
	if err := bridge.Start(); err != nil {
		t.Fatalf("Failed to start bridge: %v", err)
	}
	defer bridge.Stop()

	// Create shard
	shard, err := bridge.CreateShard(shardPath)
	if err != nil {
		t.Fatalf("Failed to create shard: %v", err)
	}
	defer shard.Close()

	// Index some documents
	docs := []map[string]interface{}{
		{"title": "Document 1", "count": int64(100)},
		{"title": "Document 2", "count": int64(200)},
		{"title": "Document 3", "count": int64(300)},
	}

	for i, doc := range docs {
		docID := string(rune('A' + i))
		if err := shard.IndexDocument(docID, doc); err != nil {
			t.Fatalf("Failed to index document %s: %v", docID, err)
		}
	}

	// Commit to disk
	if err := shard.Commit(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Get stats
	stats, err := shard.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	// Verify stats
	if stats.NumDocs != 3 {
		t.Errorf("Expected 3 documents, got %d", stats.NumDocs)
	}

	if stats.MaxDoc < 3 {
		t.Errorf("Expected max_doc >= 3, got %d", stats.MaxDoc)
	}

	if stats.SegmentCount < 1 {
		t.Errorf("Expected at least 1 segment, got %d", stats.SegmentCount)
	}

	if stats.SizeBytes < 0 {
		t.Errorf("Expected non-negative size, got %d", stats.SizeBytes)
	}

	t.Logf("Shard stats: %d docs, %d max_doc, %d segments, %d bytes",
		stats.NumDocs, stats.MaxDoc, stats.SegmentCount, stats.SizeBytes)

	// Note: SizeBytes might be 0 for new indexes depending on Diagon implementation
}

func TestShardForceMerge(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "diagon-merge-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	shardPath := filepath.Join(tmpDir, "shard0")
	if err := os.MkdirAll(shardPath, 0755); err != nil {
		t.Fatalf("Failed to create shard directory: %v", err)
	}

	// Create bridge
	logger, _ := zap.NewDevelopment()
	bridge, err := NewDiagonBridge(&Config{
		DataDir: tmpDir,
		Logger:  logger,
	})
	if err != nil {
		t.Fatalf("Failed to create bridge: %v", err)
	}

	// Start bridge
	if err := bridge.Start(); err != nil {
		t.Fatalf("Failed to start bridge: %v", err)
	}
	defer bridge.Stop()

	// Create shard
	shard, err := bridge.CreateShard(shardPath)
	if err != nil {
		t.Fatalf("Failed to create shard: %v", err)
	}
	defer shard.Close()

	// Index documents in batches to create multiple segments
	for batch := 0; batch < 5; batch++ {
		for i := 0; i < 100; i++ {
			docID := string(rune('A' + (batch*100 + i)))
			doc := map[string]interface{}{
				"title": "Document batch " + string(rune('0'+batch)),
				"count": int64(batch*100 + i),
			}
			if err := shard.IndexDocument(docID, doc); err != nil {
				t.Fatalf("Failed to index document: %v", err)
			}
		}
		// Flush each batch to create separate segments
		if err := shard.Flush(); err != nil {
			t.Fatalf("Failed to flush batch %d: %v", batch, err)
		}
	}

	// Commit all
	if err := shard.Commit(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Get stats before merge
	statsBefore, err := shard.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats before merge: %v", err)
	}

	t.Logf("Before merge: %d segments", statsBefore.SegmentCount)

	// Force merge to 1 segment
	if err := shard.ForceMerge(1); err != nil {
		t.Fatalf("Failed to force merge: %v", err)
	}

	// Get stats after merge
	statsAfter, err := shard.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats after merge: %v", err)
	}

	t.Logf("After merge: %d segments", statsAfter.SegmentCount)

	// Verify merge worked
	if statsAfter.SegmentCount > statsBefore.SegmentCount {
		t.Errorf("Segment count increased after merge: before=%d, after=%d",
			statsBefore.SegmentCount, statsAfter.SegmentCount)
	}

	// Document count should remain the same
	if statsAfter.NumDocs != statsBefore.NumDocs {
		t.Errorf("Document count changed after merge: before=%d, after=%d",
			statsBefore.NumDocs, statsAfter.NumDocs)
	}

	// Verify search still works after merge
	query := []byte(`{"match_all": {}}`)
	result, err := shard.Search(query, nil)
	if err != nil {
		t.Fatalf("Search failed after merge: %v", err)
	}

	if result.TotalHits != 500 {
		t.Errorf("Expected 500 documents, got %d", result.TotalHits)
	}
}

func TestShardGetStatsEmpty(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "diagon-stats-empty-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	shardPath := filepath.Join(tmpDir, "shard0")
	if err := os.MkdirAll(shardPath, 0755); err != nil {
		t.Fatalf("Failed to create shard directory: %v", err)
	}

	// Create bridge
	logger, _ := zap.NewDevelopment()
	bridge, err := NewDiagonBridge(&Config{
		DataDir: tmpDir,
		Logger:  logger,
	})
	if err != nil {
		t.Fatalf("Failed to create bridge: %v", err)
	}

	// Start bridge
	if err := bridge.Start(); err != nil {
		t.Fatalf("Failed to start bridge: %v", err)
	}
	defer bridge.Stop()

	// Create shard
	shard, err := bridge.CreateShard(shardPath)
	if err != nil {
		t.Fatalf("Failed to create shard: %v", err)
	}
	defer shard.Close()

	// Commit empty index
	if err := shard.Commit(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Get stats on empty index
	stats, err := shard.GetStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	// Verify empty stats
	if stats.NumDocs != 0 {
		t.Errorf("Expected 0 documents, got %d", stats.NumDocs)
	}

	t.Logf("Empty shard stats: %d docs, %d segments, %d bytes",
		stats.NumDocs, stats.SegmentCount, stats.SizeBytes)
}
