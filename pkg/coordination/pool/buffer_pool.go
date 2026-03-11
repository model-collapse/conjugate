package pool

import (
	"io"
	"sync"
)

// BufferPool provides pooling for IO buffers
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates a new buffer pool
func NewBufferPool(initialSize int) *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, initialSize)
			},
		},
	}
}

// Get returns a buffer from the pool
func (bp *BufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

// Put returns a buffer to the pool
func (bp *BufferPool) Put(buf []byte) {
	// Don't pool buffers that grew too large (>1MB)
	if cap(buf) > 1024*1024 {
		return
	}
	// Reset length to 0 but keep capacity
	bp.pool.Put(buf[:0])
}

// ReadAll reads from reader using a pooled buffer
// This is more efficient than io.ReadAll for typical request sizes
func (bp *BufferPool) ReadAll(r io.Reader) ([]byte, error) {
	buf := bp.Get()
	defer bp.Put(buf)

	// Read in chunks
	var result []byte
	for {
		if len(buf) == 0 {
			buf = make([]byte, 4096)
		}

		n, err := r.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
		}

		if err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}
	}
}
