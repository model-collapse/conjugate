package executor

import (
	"testing"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
	"go.uber.org/zap"
)

// BenchmarkAggregationMerging benchmarks parallel aggregation merging
func BenchmarkAggregationMerging(b *testing.B) {
	logger := zap.NewNop()
	qe := &QueryExecutor{logger: logger}

	// Create test data: 5 shards, each with 5 aggregations (typical query)
	responses := make([]*pb.SearchResponse, 5)
	for i := 0; i < 5; i++ {
		responses[i] = &pb.SearchResponse{
			Aggregations: map[string]*pb.AggregationResult{
				"agg1_terms": {
					Type: "terms",
					Buckets: []*pb.AggregationBucket{
						{Key: "bucket1", DocCount: 100},
						{Key: "bucket2", DocCount: 200},
						{Key: "bucket3", DocCount: 150},
					},
				},
				"agg2_stats": {
					Type:  "stats",
					Count: 1000,
					Min:   1.0,
					Max:   100.0,
					Sum:   50000.0,
					Avg:   50.0,
				},
				"agg3_histogram": {
					Type: "histogram",
					Buckets: []*pb.AggregationBucket{
						{NumericKey: 0.0, DocCount: 50},
						{NumericKey: 10.0, DocCount: 75},
						{NumericKey: 20.0, DocCount: 100},
					},
				},
				"agg4_cardinality": {
					Type:  "cardinality",
					Value: 500,
				},
				"agg5_avg": {
					Type: "avg",
					Avg:  42.5,
				},
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := qe.mergeAggregations(responses)
		if len(result) != 5 {
			b.Fatalf("Expected 5 aggregations, got %d", len(result))
		}
	}
}

// BenchmarkAggregationMergingManyAggs benchmarks with many aggregations
func BenchmarkAggregationMergingManyAggs(b *testing.B) {
	logger := zap.NewNop()
	qe := &QueryExecutor{logger: logger}

	// Create test data: 3 shards, each with 20 aggregations
	responses := make([]*pb.SearchResponse, 3)
	for i := 0; i < 3; i++ {
		aggs := make(map[string]*pb.AggregationResult)
		for j := 0; j < 20; j++ {
			aggName := string(rune('a' + j))
			aggs[aggName] = &pb.AggregationResult{
				Type: "terms",
				Buckets: []*pb.AggregationBucket{
					{Key: "bucket1", DocCount: 100},
					{Key: "bucket2", DocCount: 200},
				},
			}
		}
		responses[i] = &pb.SearchResponse{Aggregations: aggs}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := qe.mergeAggregations(responses)
		if len(result) != 20 {
			b.Fatalf("Expected 20 aggregations, got %d", len(result))
		}
	}
}

// BenchmarkAggregationMergingSingleAgg benchmarks with single aggregation (worst case for parallelism)
func BenchmarkAggregationMergingSingleAgg(b *testing.B) {
	logger := zap.NewNop()
	qe := &QueryExecutor{logger: logger}

	// Create test data: 5 shards, 1 aggregation (no parallelism benefit)
	responses := make([]*pb.SearchResponse, 5)
	for i := 0; i < 5; i++ {
		responses[i] = &pb.SearchResponse{
			Aggregations: map[string]*pb.AggregationResult{
				"agg1": {
					Type: "terms",
					Buckets: []*pb.AggregationBucket{
						{Key: "bucket1", DocCount: 100},
						{Key: "bucket2", DocCount: 200},
					},
				},
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := qe.mergeAggregations(responses)
		if len(result) != 1 {
			b.Fatalf("Expected 1 aggregation, got %d", len(result))
		}
	}
}
