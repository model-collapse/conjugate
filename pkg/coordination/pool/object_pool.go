package pool

import (
	"bytes"
	"sync"

	pb "github.com/conjugate/conjugate/pkg/common/proto"
)

// ObjectPool provides pooling for frequently allocated objects
// to reduce GC pressure and improve performance
type ObjectPool struct {
	// Buffer pool for JSON/response building
	bufferPool sync.Pool

	// IO buffer pool for reading request bodies
	ioBufferPool sync.Pool

	// Search request context pool
	searchRequestPool sync.Pool

	// Search response pool
	searchResponsePool sync.Pool

	// Hit slice pool (for collecting results)
	hitSlicePool sync.Pool
}

// NewObjectPool creates a new object pool
func NewObjectPool() *ObjectPool {
	return &ObjectPool{
		bufferPool: sync.Pool{
			New: func() interface{} {
				// 4KB initial capacity - typical response size
				return bytes.NewBuffer(make([]byte, 0, 4096))
			},
		},
		ioBufferPool: sync.Pool{
			New: func() interface{} {
				// 16KB buffer for reading request bodies
				buf := make([]byte, 16384)
				return &buf
			},
		},
		searchRequestPool: sync.Pool{
			New: func() interface{} {
				return &SearchRequestContext{}
			},
		},
		searchResponsePool: sync.Pool{
			New: func() interface{} {
				return &SearchResponseContext{}
			},
		},
		hitSlicePool: sync.Pool{
			New: func() interface{} {
				// Pre-allocate capacity for typical result size (100 hits)
				hits := make([]*pb.SearchHit, 0, 100)
				return &hits
			},
		},
	}
}

// GetBuffer returns a buffer from the pool
func (p *ObjectPool) GetBuffer() *bytes.Buffer {
	buf := p.bufferPool.Get().(*bytes.Buffer)
	buf.Reset() // Clear any previous content
	return buf
}

// PutBuffer returns a buffer to the pool
func (p *ObjectPool) PutBuffer(buf *bytes.Buffer) {
	// Don't pool buffers that grew too large (>1MB)
	if buf.Cap() > 1024*1024 {
		return
	}
	p.bufferPool.Put(buf)
}

// GetIOBuffer returns an IO buffer from the pool for reading request bodies
func (p *ObjectPool) GetIOBuffer() *[]byte {
	return p.ioBufferPool.Get().(*[]byte)
}

// PutIOBuffer returns an IO buffer to the pool
func (p *ObjectPool) PutIOBuffer(buf *[]byte) {
	// Don't pool buffers that grew too large (>1MB)
	if cap(*buf) > 1024*1024 {
		return
	}
	// Reset to original size
	*buf = (*buf)[:cap(*buf)]
	p.ioBufferPool.Put(buf)
}

// GetSearchRequest returns a search request context from the pool
func (p *ObjectPool) GetSearchRequest() *SearchRequestContext {
	req := p.searchRequestPool.Get().(*SearchRequestContext)
	req.Reset()
	return req
}

// PutSearchRequest returns a search request context to the pool
func (p *ObjectPool) PutSearchRequest(req *SearchRequestContext) {
	p.searchRequestPool.Put(req)
}

// GetSearchResponse returns a search response context from the pool
func (p *ObjectPool) GetSearchResponse() *SearchResponseContext {
	resp := p.searchResponsePool.Get().(*SearchResponseContext)
	resp.Reset()
	return resp
}

// PutSearchResponse returns a search response context to the pool
func (p *ObjectPool) PutSearchResponse(resp *SearchResponseContext) {
	p.searchResponsePool.Put(resp)
}

// GetHitSlice returns a hit slice from the pool
func (p *ObjectPool) GetHitSlice() *[]*pb.SearchHit {
	slice := p.hitSlicePool.Get().(*[]*pb.SearchHit)
	*slice = (*slice)[:0] // Reset length to 0, keep capacity
	return slice
}

// PutHitSlice returns a hit slice to the pool
func (p *ObjectPool) PutHitSlice(slice *[]*pb.SearchHit) {
	// Don't pool slices that grew too large
	if cap(*slice) > 10000 {
		return
	}
	p.hitSlicePool.Put(slice)
}

// SearchRequestContext holds parsed search request data
type SearchRequestContext struct {
	IndexName        string
	Query            []byte
	FilterExpression []byte
	From             int
	Size             int
	Source           []string
	Sort             []string
	Aggregations     map[string]interface{}
}

// Reset clears the search request context for reuse
func (r *SearchRequestContext) Reset() {
	r.IndexName = ""
	r.Query = r.Query[:0]
	r.FilterExpression = r.FilterExpression[:0]
	r.From = 0
	r.Size = 10
	r.Source = r.Source[:0]
	r.Sort = r.Sort[:0]
	// Don't clear aggregations map, just set to nil
	// Maps are expensive to clear
	r.Aggregations = nil
}

// SearchResponseContext holds search response building data
type SearchResponseContext struct {
	TookMillis   int64
	TotalHits    int64
	MaxScore     float64
	Hits         []interface{} // Generic hit storage
	Aggregations map[string]interface{}
}

// Reset clears the search response context for reuse
func (r *SearchResponseContext) Reset() {
	r.TookMillis = 0
	r.TotalHits = 0
	r.MaxScore = 0
	r.Hits = r.Hits[:0]
	r.Aggregations = nil
}
