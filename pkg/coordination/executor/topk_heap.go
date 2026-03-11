package executor

import (
	"container/heap"
)

// topKHeap is a min-heap that maintains the top K search hits by score
// We use a min-heap because we want to quickly evict the lowest-scoring hit
// when we encounter a higher-scoring one
type topKHeap struct {
	hits []*SearchHit
	k    int // Maximum capacity
}

// Len implements heap.Interface
func (h *topKHeap) Len() int {
	return len(h.hits)
}

// Less implements heap.Interface
// Min-heap: parent has lower score than children
func (h *topKHeap) Less(i, j int) bool {
	return h.hits[i].Score < h.hits[j].Score
}

// Swap implements heap.Interface
func (h *topKHeap) Swap(i, j int) {
	h.hits[i], h.hits[j] = h.hits[j], h.hits[i]
}

// Push implements heap.Interface
func (h *topKHeap) Push(x interface{}) {
	h.hits = append(h.hits, x.(*SearchHit))
}

// Pop implements heap.Interface
func (h *topKHeap) Pop() interface{} {
	old := h.hits
	n := len(old)
	x := old[n-1]
	h.hits = old[0 : n-1]
	return x
}

// Add adds a hit to the heap, maintaining the top K constraint
func (h *topKHeap) Add(hit *SearchHit) {
	if len(h.hits) < h.k {
		// Heap not full yet, just add
		heap.Push(h, hit)
	} else if hit.Score > h.hits[0].Score {
		// Hit has higher score than minimum, replace minimum
		h.hits[0] = hit
		heap.Fix(h, 0)
	}
	// Otherwise, hit doesn't make the cut, discard it
}

// ToSlice extracts all hits from the heap in descending score order
func (h *topKHeap) ToSlice() []*SearchHit {
	// Extract elements from heap (they come out in arbitrary order)
	result := make([]*SearchHit, len(h.hits))
	copy(result, h.hits)

	// Sort in descending order by score
	// We use a simple bubble-sort-like approach since the heap is small (typically 10-100)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Score < result[j].Score {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// newTopKHeap creates a new top-K heap with the given capacity
func newTopKHeap(k int) *topKHeap {
	h := &topKHeap{
		hits: make([]*SearchHit, 0, k),
		k:    k,
	}
	heap.Init(h)
	return h
}
