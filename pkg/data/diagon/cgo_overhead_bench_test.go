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

// TestCGOOverheadBreakdown measures each phase of the indexing pipeline separately:
//   - Phase A: Go→C conversion (flattenMap + CString alloc + diagon_create_*_field)
//   - Phase B: Pure C++ indexing (diagon_add_documents only)
//   - Phase C: Document cleanup (free)
//   - Full bridge reference (BulkIndexDocuments)
//
// Run: go test ./pkg/data/diagon/ -run TestCGOOverheadBreakdown -v -count=1 -timeout 120s
func TestCGOOverheadBreakdown(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cgo-overhead-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logger, _ := zap.NewDevelopment()
	bridge, err := NewDiagonBridge(&Config{DataDir: tmpDir, Logger: logger})
	if err != nil {
		t.Fatalf("Failed to create bridge: %v", err)
	}
	bridge.Start()
	defer bridge.Stop()

	// Big5 realistic document (~25 flattened fields)
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

	sampleJSON, _ := json.Marshal(sampleDoc)
	flatSample := flattenMap("", sampleDoc)
	t.Logf("Document: %d flattened fields, %d bytes JSON", len(flatSample), len(sampleJSON))

	batchSize := 500
	numBatches := 40
	totalDocs := batchSize * numBatches

	t.Logf("=== CGO Overhead Breakdown: %d docs (%d batches x %d) ===\n", totalDocs, numBatches, batchSize)

	// ========================================================
	// Phase measurement via BulkIndexWithTiming
	// ========================================================
	t.Run("PhaseSplit_PrepVsCpp", func(t *testing.T) {
		shardPath := filepath.Join(tmpDir, "shard-timed")
		if err := os.MkdirAll(shardPath, 0755); err != nil {
			t.Fatalf("Failed to create shard dir: %v", err)
		}
		shard, err := bridge.CreateShard(shardPath)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}
		defer shard.Close()

		var totalPrep, totalAdd, totalFree time.Duration

		for b := 0; b < numBatches; b++ {
			batch := make([]struct {
				ID         string
				Doc        map[string]interface{}
				SourceJSON []byte
			}, batchSize)
			for i := 0; i < batchSize; i++ {
				batch[i].ID = fmt.Sprintf("timed-%d-%d", b, i)
				batch[i].Doc = sampleDoc
				batch[i].SourceJSON = sampleJSON
			}

			timing, err := shard.BulkIndexWithTiming(batch)
			if err != nil {
				t.Fatalf("BulkIndexWithTiming failed at batch %d: %v", b, err)
			}

			totalPrep += timing.PrepTime
			totalAdd += timing.AddTime
			totalFree += timing.FreeTime
		}

		totalBridge := totalPrep + totalAdd + totalFree

		// Commit
		commitStart := time.Now()
		shard.Commit()
		commitTime := time.Since(commitStart)

		totalWithCommit := totalBridge + commitTime

		t.Logf("\n=== PER-PHASE RESULTS (%d docs) ===", totalDocs)
		t.Logf("Go→C prep (flatten+CString+fields): %v  (%.1f µs/doc, %.0f docs/sec)",
			totalPrep,
			float64(totalPrep.Microseconds())/float64(totalDocs),
			float64(totalDocs)/totalPrep.Seconds())
		t.Logf("Pure C++ add_documents:              %v  (%.1f µs/doc, %.0f docs/sec)",
			totalAdd,
			float64(totalAdd.Microseconds())/float64(totalDocs),
			float64(totalDocs)/totalAdd.Seconds())
		t.Logf("Doc free (cleanup):                  %v  (%.1f µs/doc)",
			totalFree,
			float64(totalFree.Microseconds())/float64(totalDocs))
		t.Logf("Bridge total (prep+add+free):        %v  (%.1f µs/doc, %.0f docs/sec)",
			totalBridge,
			float64(totalBridge.Microseconds())/float64(totalDocs),
			float64(totalDocs)/totalBridge.Seconds())
		t.Logf("Commit:                              %v", commitTime)
		t.Logf("Bridge + commit:                     %v  (%.0f docs/sec)",
			totalWithCommit,
			float64(totalDocs)/totalWithCommit.Seconds())

		t.Logf("\n=== PERCENTAGE BREAKDOWN (bridge time only) ===")
		t.Logf("Go→C prep:        %5.1f%%", float64(totalPrep)/float64(totalBridge)*100)
		t.Logf("C++ add_documents: %5.1f%%", float64(totalAdd)/float64(totalBridge)*100)
		t.Logf("Doc free:          %5.1f%%", float64(totalFree)/float64(totalBridge)*100)

		t.Logf("\n=== PERCENTAGE BREAKDOWN (bridge + commit) ===")
		t.Logf("Go→C prep:         %5.1f%%", float64(totalPrep)/float64(totalWithCommit)*100)
		t.Logf("C++ add_documents:  %5.1f%%", float64(totalAdd)/float64(totalWithCommit)*100)
		t.Logf("Doc free:           %5.1f%%", float64(totalFree)/float64(totalWithCommit)*100)
		t.Logf("Commit:             %5.1f%%", float64(commitTime)/float64(totalWithCommit)*100)

		// Estimate pipeline overhead from end-to-end numbers.
		// Reference: Big5 HTTP bulk indexing benchmark (Mar 2026, commit c170a30).
		// Update this value when re-running end-to-end benchmarks.
		e2eRate := 22141.0
		bridgeRate := float64(totalDocs) / totalWithCommit.Seconds()
		pipelineOverhead := (1.0/e2eRate - 1.0/bridgeRate) * 1e6 // µs per doc

		t.Logf("\n=== FULL PIPELINE OVERHEAD ESTIMATE ===")
		t.Logf("Raw bridge + commit: %.0f docs/sec (%.1f µs/doc)", bridgeRate, 1e6/bridgeRate)
		t.Logf("End-to-end HTTP:     %.0f docs/sec (%.1f µs/doc)", e2eRate, 1e6/e2eRate)
		t.Logf("HTTP/gRPC overhead:  %.1f µs/doc", pipelineOverhead)
		t.Logf("OpenSearch ref:      20824 docs/sec (48.0 µs/doc)")
	})

	// ========================================================
	// Reference: Full BulkIndexDocuments via normal API
	// ========================================================
	t.Run("Reference_FullBridge", func(t *testing.T) {
		shardPath := filepath.Join(tmpDir, "shard-ref")
		if err := os.MkdirAll(shardPath, 0755); err != nil {
			t.Fatalf("Failed to create shard dir: %v", err)
		}
		shard, err := bridge.CreateShard(shardPath)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}
		defer shard.Close()

		start := time.Now()
		for b := 0; b < numBatches; b++ {
			batch := make([]struct {
				ID         string
				Doc        map[string]interface{}
				SourceJSON []byte
			}, batchSize)
			for i := 0; i < batchSize; i++ {
				batch[i].ID = fmt.Sprintf("ref-%d-%d", b, i)
				batch[i].Doc = sampleDoc
				batch[i].SourceJSON = sampleJSON
			}
			if err := shard.BulkIndexDocuments(batch); err != nil {
				t.Fatalf("BulkIndex failed at batch %d: %v", b, err)
			}
		}
		indexTime := time.Since(start)

		commitStart := time.Now()
		shard.Commit()
		commitTime := time.Since(commitStart)

		t.Logf("BulkIndexDocuments (no commit): %.0f docs/sec (%.1f µs/doc)",
			float64(totalDocs)/indexTime.Seconds(),
			float64(indexTime.Microseconds())/float64(totalDocs))
		t.Logf("Commit: %v", commitTime)
		t.Logf("Total (index+commit): %.0f docs/sec",
			float64(totalDocs)/(indexTime+commitTime).Seconds())
	})

	// ========================================================
	// Larger batch test: 2000 docs per batch (amortize CGO overhead)
	// ========================================================
	t.Run("LargeBatch_2000", func(t *testing.T) {
		shardPath := filepath.Join(tmpDir, "shard-large")
		if err := os.MkdirAll(shardPath, 0755); err != nil {
			t.Fatalf("Failed to create shard dir: %v", err)
		}
		shard, err := bridge.CreateShard(shardPath)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}
		defer shard.Close()

		largeBatch := 2000
		numLargeBatches := 10
		largeTotalDocs := largeBatch * numLargeBatches

		var totalPrep, totalAdd, totalFree time.Duration

		for b := 0; b < numLargeBatches; b++ {
			batch := make([]struct {
				ID         string
				Doc        map[string]interface{}
				SourceJSON []byte
			}, largeBatch)
			for i := 0; i < largeBatch; i++ {
				batch[i].ID = fmt.Sprintf("lg-%d-%d", b, i)
				batch[i].Doc = sampleDoc
				batch[i].SourceJSON = sampleJSON
			}

			timing, err := shard.BulkIndexWithTiming(batch)
			if err != nil {
				t.Fatalf("BulkIndexWithTiming failed: %v", err)
			}
			totalPrep += timing.PrepTime
			totalAdd += timing.AddTime
			totalFree += timing.FreeTime
		}

		totalBridge := totalPrep + totalAdd + totalFree
		commitStart := time.Now()
		shard.Commit()
		commitTime := time.Since(commitStart)

		t.Logf("\n=== LARGE BATCH (2000 docs/batch, %d total) ===", largeTotalDocs)
		t.Logf("Go→C prep:           %.1f µs/doc (%.0f docs/sec)",
			float64(totalPrep.Microseconds())/float64(largeTotalDocs),
			float64(largeTotalDocs)/totalPrep.Seconds())
		t.Logf("C++ add_documents:    %.1f µs/doc (%.0f docs/sec)",
			float64(totalAdd.Microseconds())/float64(largeTotalDocs),
			float64(largeTotalDocs)/totalAdd.Seconds())
		t.Logf("Bridge total:         %.1f µs/doc (%.0f docs/sec)",
			float64(totalBridge.Microseconds())/float64(largeTotalDocs),
			float64(largeTotalDocs)/totalBridge.Seconds())
		t.Logf("Commit: %v", commitTime)
		t.Logf("C++ pct of bridge: %.1f%%", float64(totalAdd)/float64(totalBridge)*100)
	})
}
