package executor

import (
	"fmt"
	"testing"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
	"go.uber.org/zap"
)

// BenchmarkResultSorting benchmarks the heap-based top-K selection
func BenchmarkResultSorting(b *testing.B) {
	logger := zap.NewNop()
	qe := &QueryExecutor{logger: logger}

	// Create test data: 5 shards, each returning 100 hits (500 total)
	// Requesting top 10 (from=0, size=10)
	responses := make([]*pb.SearchResponse, 5)
	for i := 0; i < 5; i++ {
		hits := make([]*pb.SearchHit, 100)
		for j := 0; j < 100; j++ {
			hits[j] = &pb.SearchHit{
				Id:         fmt.Sprintf("doc-%d-%d", i, j),
				Score:      float64(1000 - i*100 - j), // Descending scores
				// Benchmark doesn't need actual source data
			}
		}
		responses[i] = &pb.SearchResponse{
			Hits: &pb.SearchHits{
				Total:    &pb.TotalHits{Value: 100, Relation: "eq"},
				MaxScore: float64(1000 - i*100),
				Hits:     hits,
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := qe.aggregateSearchResults(responses, 0, 10)
		if len(result.Hits) != 10 {
			b.Fatalf("Expected 10 hits, got %d", len(result.Hits))
		}
	}
}

// BenchmarkResultSortingLargeDataset benchmarks with many hits
func BenchmarkResultSortingLargeDataset(b *testing.B) {
	logger := zap.NewNop()
	qe := &QueryExecutor{logger: logger}

	// Create test data: 10 shards, each returning 1000 hits (10,000 total)
	// Requesting top 10 (from=0, size=10)
	responses := make([]*pb.SearchResponse, 10)
	for i := 0; i < 10; i++ {
		hits := make([]*pb.SearchHit, 1000)
		for j := 0; j < 1000; j++ {
			hits[j] = &pb.SearchHit{
				Id:         fmt.Sprintf("doc-%d-%d", i, j),
				Score:      float64(10000 - i*1000 - j),
				// Benchmark doesn't need actual source data
			}
		}
		responses[i] = &pb.SearchResponse{
			Hits: &pb.SearchHits{
				Total:    &pb.TotalHits{Value: 1000, Relation: "eq"},
				MaxScore: float64(10000 - i*1000),
				Hits:     hits,
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := qe.aggregateSearchResults(responses, 0, 10)
		if len(result.Hits) != 10 {
			b.Fatalf("Expected 10 hits, got %d", len(result.Hits))
		}
	}
}

// BenchmarkResultSortingDeepPagination benchmarks with deep pagination
func BenchmarkResultSortingDeepPagination(b *testing.B) {
	logger := zap.NewNop()
	qe := &QueryExecutor{logger: logger}

	// Create test data: 5 shards, each returning 100 hits
	// Requesting from=90, size=10 (deep pagination)
	responses := make([]*pb.SearchResponse, 5)
	for i := 0; i < 5; i++ {
		hits := make([]*pb.SearchHit, 100)
		for j := 0; j < 100; j++ {
			hits[j] = &pb.SearchHit{
				Id:         fmt.Sprintf("doc-%d-%d", i, j),
				Score:      float64(1000 - i*100 - j),
				// Benchmark doesn't need actual source data
			}
		}
		responses[i] = &pb.SearchResponse{
			Hits: &pb.SearchHits{
				Total:    &pb.TotalHits{Value: 100, Relation: "eq"},
				MaxScore: float64(1000 - i*100),
				Hits:     hits,
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := qe.aggregateSearchResults(responses, 90, 10)
		if len(result.Hits) != 10 {
			b.Fatalf("Expected 10 hits, got %d", len(result.Hits))
		}
	}
}

// BenchmarkTopKHeap benchmarks the heap operations directly
func BenchmarkTopKHeap(b *testing.B) {
	k := 10
	numHits := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := newTopKHeap(k)
		for j := 0; j < numHits; j++ {
			hit := &SearchHit{
				ID:    fmt.Sprintf("doc-%d", j),
				Score: float64(numHits - j),
			}
			h.Add(hit)
		}
		_ = h.ToSlice()
	}
}
