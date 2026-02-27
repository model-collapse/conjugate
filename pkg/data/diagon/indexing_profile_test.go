package diagon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.uber.org/zap"
)

// TestBulkIndexBreakdown profiles the indexing path to identify bottlenecks.
// Run with: go test ./pkg/data/diagon/ -run TestBulkIndexBreakdown -v -count=1 -timeout 120s
func TestBulkIndexBreakdown(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "diagon-profile-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	shardPath := filepath.Join(tmpDir, "shard0")
	os.MkdirAll(shardPath, 0755)

	logger, _ := zap.NewDevelopment()
	bridge, err := NewDiagonBridge(&Config{DataDir: tmpDir, Logger: logger})
	if err != nil {
		t.Fatalf("Failed to create bridge: %v", err)
	}
	bridge.Start()
	defer bridge.Stop()

	shard, err := bridge.CreateShard(shardPath)
	if err != nil {
		t.Fatalf("Failed to create shard: %v", err)
	}
	defer shard.Close()

	// Realistic Big5 document (~20 fields, nested objects)
	sampleDoc := map[string]interface{}{
		"@timestamp": "2023-01-15T10:30:00.000Z",
		"message":    "Jan 15 10:30:00 ip-172-31-20-153 kernel: [1234567.890123] audit: type=1400 msg=audit(1234567890.123:456789): apparmor=\"STATUS\" operation=\"profile_load\" profile=\"unconfined\" name=\"/usr/bin/evince\" pid=12345 comm=\"apparmor_parser\"",
		"log": map[string]interface{}{
			"file": map[string]interface{}{"path": "/var/log/messages/birdknight"},
			"level": "info",
		},
		"cloud": map[string]interface{}{"region": "us-east-1"},
		"agent": map[string]interface{}{
			"ephemeral_id": "abc123-def456",
			"id":           "agent-001",
			"name":         "filebeat",
			"type":         "filebeat",
			"version":      "8.5.0",
		},
		"aws": map[string]interface{}{
			"cloudwatch": map[string]interface{}{
				"log_group":      "/aws/lambda/my-function",
				"log_stream":     "2023/01/15/[$LATEST]abc123",
				"ingestion_time": "2023-01-15T10:30:05.000Z",
			},
		},
		"data_stream": map[string]interface{}{
			"dataset":   "aws.cloudwatch_logs",
			"namespace": "default",
			"type":      "logs",
		},
		"host":  map[string]interface{}{"name": "ip-172-31-20-153"},
		"event": map[string]interface{}{"kind": "event", "category": "process", "type": "info"},
		"metrics": map[string]interface{}{
			"size": float64(2500),
			"tmin": float64(100),
		},
		"process": map[string]interface{}{"name": "kernel", "pid": float64(1234)},
		"tags":    "production",
	}

	batchSize := 500
	numBatches := 20
	totalDocs := batchSize * numBatches

	t.Logf("=== Indexing Profile: %d docs (%d batches x %d) ===", totalDocs, numBatches, batchSize)

	// Phase 1: JSON Marshal
	t.Run("1_JSONMarshal", func(t *testing.T) {
		start := time.Now()
		for i := 0; i < totalDocs; i++ {
			json.Marshal(sampleDoc)
		}
		elapsed := time.Since(start)
		t.Logf("JSON marshal: %.0f docs/sec (%.1f µs/doc)",
			float64(totalDocs)/elapsed.Seconds(),
			float64(elapsed.Microseconds())/float64(totalDocs))
	})

	// Phase 2: flattenMap
	t.Run("2_FlattenMap", func(t *testing.T) {
		start := time.Now()
		for i := 0; i < totalDocs; i++ {
			flattenMap("", sampleDoc)
		}
		elapsed := time.Since(start)
		t.Logf("flattenMap: %.0f docs/sec (%.1f µs/doc)",
			float64(totalDocs)/elapsed.Seconds(),
			float64(elapsed.Microseconds())/float64(totalDocs))
	})

	// Phase 3: Date parsing
	t.Run("3_DateParsing", func(t *testing.T) {
		dateStr := "2023-01-15T10:30:00.000Z"
		start := time.Now()
		for i := 0; i < totalDocs; i++ {
			tryParseDateToEpochMs(dateStr)
		}
		elapsed := time.Since(start)
		t.Logf("Date parsing: %.0f calls/sec (%.1f µs/call)",
			float64(totalDocs)/elapsed.Seconds(),
			float64(elapsed.Microseconds())/float64(totalDocs))
	})

	// Phase 4: isKeywordLike
	t.Run("4_IsKeywordLike", func(t *testing.T) {
		vals := []string{
			"/var/log/messages/birdknight",
			"us-east-1",
			"filebeat",
			"Jan 15 10:30:00 ip-172-31-20-153 kernel: long text message here with many words",
		}
		start := time.Now()
		for i := 0; i < totalDocs; i++ {
			for _, v := range vals {
				isKeywordLike(v)
			}
		}
		elapsed := time.Since(start)
		t.Logf("isKeywordLike: %.0f calls/sec (%.1f µs/call, 4 strings each)",
			float64(totalDocs)/elapsed.Seconds(),
			float64(elapsed.Microseconds())/float64(totalDocs))
	})

	// Phase 5: Full BulkIndex (CGO + Diagon C++)
	t.Run("5_FullBulkIndex", func(t *testing.T) {
		var indexTime time.Duration
		var commitTime time.Duration

		for b := 0; b < numBatches; b++ {
			batch := make([]struct {
				ID  string
				Doc map[string]interface{}
			}, batchSize)
			for i := 0; i < batchSize; i++ {
				batch[i].ID = fmt.Sprintf("doc-%d-%d", b, i)
				batch[i].Doc = sampleDoc
			}

			start := time.Now()
			err := shard.BulkIndexDocuments(batch)
			indexTime += time.Since(start)
			if err != nil {
				t.Fatalf("BulkIndex failed at batch %d: %v", b, err)
			}

			commitStart := time.Now()
			shard.Commit()
			commitTime += time.Since(commitStart)
		}

		totalTime := indexTime + commitTime
		t.Logf("BulkIndex (no commit): %.0f docs/sec (%.1f µs/doc)",
			float64(totalDocs)/indexTime.Seconds(),
			float64(indexTime.Microseconds())/float64(totalDocs))
		t.Logf("Commit: %d commits in %v (%.1f ms/commit)",
			numBatches, commitTime,
			float64(commitTime.Milliseconds())/float64(numBatches))
		t.Logf("Total (index+commit): %.0f docs/sec",
			float64(totalDocs)/totalTime.Seconds())
	})

	// Phase 6: BulkIndex WITHOUT commit (pure indexing throughput)
	t.Run("6_BulkIndexNoCommit", func(t *testing.T) {
		// Create a fresh shard to avoid segment accumulation effects
		shardPath2 := filepath.Join(tmpDir, "shard1")
		os.MkdirAll(shardPath2, 0755)
		shard2, err := bridge.CreateShard(shardPath2)
		if err != nil {
			t.Fatalf("Failed to create shard2: %v", err)
		}
		defer shard2.Close()

		start := time.Now()
		for b := 0; b < numBatches; b++ {
			batch := make([]struct {
				ID  string
				Doc map[string]interface{}
			}, batchSize)
			for i := 0; i < batchSize; i++ {
				batch[i].ID = fmt.Sprintf("nc-%d-%d", b, i)
				batch[i].Doc = sampleDoc
			}
			if err := shard2.BulkIndexDocuments(batch); err != nil {
				t.Fatalf("BulkIndex failed: %v", err)
			}
		}
		indexTime := time.Since(start)

		// Single commit at end
		commitStart := time.Now()
		shard2.Commit()
		commitTime := time.Since(commitStart)

		t.Logf("BulkIndex (no intermediate commits): %.0f docs/sec (%.1f µs/doc)",
			float64(totalDocs)/indexTime.Seconds(),
			float64(indexTime.Microseconds())/float64(totalDocs))
		t.Logf("Single final commit: %v", commitTime)
		t.Logf("Total: %.0f docs/sec",
			float64(totalDocs)/(indexTime+commitTime).Seconds())
	})

	// Summary
	t.Logf("\n=== SUMMARY ===")
	t.Logf("Document: ~%d fields (Big5 realistic)", len(flattenMap("", sampleDoc)))
	t.Logf("Batches: %d x %d docs = %d total", numBatches, batchSize, totalDocs)
}
