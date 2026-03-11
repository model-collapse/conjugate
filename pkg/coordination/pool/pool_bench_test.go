package pool

import (
	"bytes"
	"testing"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
)

// BenchmarkBufferPooling benchmarks buffer allocation with and without pooling
func BenchmarkBufferPooling(b *testing.B) {
	b.Run("WithoutPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := bytes.NewBuffer(make([]byte, 0, 4096))
			buf.WriteString("test data")
			_ = buf.Bytes()
		}
	})

	b.Run("WithPool", func(b *testing.B) {
		pool := NewObjectPool()
		for i := 0; i < b.N; i++ {
			buf := pool.GetBuffer()
			buf.WriteString("test data")
			_ = buf.Bytes()
			pool.PutBuffer(buf)
		}
	})
}

// BenchmarkIOBufferPooling benchmarks IO buffer allocation
func BenchmarkIOBufferPooling(b *testing.B) {
	b.Run("WithoutPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := make([]byte, 16384)
			_ = buf
		}
	})

	b.Run("WithPool", func(b *testing.B) {
		pool := NewObjectPool()
		for i := 0; i < b.N; i++ {
			buf := pool.GetIOBuffer()
			_ = buf
			pool.PutIOBuffer(buf)
		}
	})
}

// BenchmarkHitSlicePooling benchmarks hit slice allocation
func BenchmarkHitSlicePooling(b *testing.B) {
	b.Run("WithoutPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			hits := make([]*pb.SearchHit, 0, 100)
			for j := 0; j < 10; j++ {
				hits = append(hits, &pb.SearchHit{
					Id:    "doc1",
					Score: 1.0,
				})
			}
			_ = hits
		}
	})

	b.Run("WithPool", func(b *testing.B) {
		pool := NewObjectPool()
		for i := 0; i < b.N; i++ {
			hits := pool.GetHitSlice()
			for j := 0; j < 10; j++ {
				*hits = append(*hits, &pb.SearchHit{
					Id:    "doc1",
					Score: 1.0,
				})
			}
			pool.PutHitSlice(hits)
		}
	})
}

// BenchmarkSearchRequestContextPooling benchmarks request context allocation
func BenchmarkSearchRequestContextPooling(b *testing.B) {
	b.Run("WithoutPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := &SearchRequestContext{
				IndexName: "test",
				From:      0,
				Size:      10,
			}
			_ = req
		}
	})

	b.Run("WithPool", func(b *testing.B) {
		pool := NewObjectPool()
		for i := 0; i < b.N; i++ {
			req := pool.GetSearchRequest()
			req.IndexName = "test"
			req.From = 0
			req.Size = 10
			pool.PutSearchRequest(req)
		}
	})
}

// BenchmarkHighThroughput simulates high-throughput query handling
func BenchmarkHighThroughput(b *testing.B) {
	pool := NewObjectPool()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate a query request lifecycle
		reqCtx := pool.GetSearchRequest()
		reqCtx.IndexName = "logs"
		reqCtx.From = 0
		reqCtx.Size = 10

		// Simulate processing
		hits := pool.GetHitSlice()
		for j := 0; j < 10; j++ {
			*hits = append(*hits, &pb.SearchHit{
				Id:     "doc",
				Score:  1.0,
				SourceJson: nil, // Benchmark doesn't need actual source data
			})
		}

		// Simulate response building
		buf := pool.GetBuffer()
		buf.WriteString(`{"took": 5, "hits": {"total": 10, "hits": [`)
		buf.WriteString(`]}}`)

		// Return to pools
		pool.PutHitSlice(hits)
		pool.PutBuffer(buf)
		pool.PutSearchRequest(reqCtx)
	}
}
