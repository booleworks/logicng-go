package sat

import (
	"fmt"
	"strings"
)

type lngheap struct {
	solver  *CoreSolver
	heap    []int32
	indices []int
}

func newLngHeap(solver *CoreSolver) *lngheap {
	return &lngheap{
		solver:  solver,
		heap:    make([]int32, 0, 1000),
		indices: make([]int, 0, 1000),
	}
}

func (h *lngheap) left(pos int) int {
	return pos*2 + 1
}

func (h *lngheap) right(pos int) int {
	return (pos + 1) * 2
}

func (h *lngheap) parent(pos int) int {
	return (pos - 1) >> 1
}

func (h *lngheap) size() int {
	return len(h.heap)
}

func (h *lngheap) isEmpty() bool {
	return len(h.heap) == 0
}

func (h *lngheap) inHeap(n int32) bool {
	return int(n) < len(h.indices) && h.indices[n] >= 0
}

func (h *lngheap) get(index int32) int32 {
	return h.heap[index]
}

func (h *lngheap) decrease(n int32) {
	h.percolateUp(h.indices[n])
}

func (h *lngheap) insert(n int32) {
	h.growTo(int(n) + 1)
	h.indices[n] = len(h.heap)
	h.heap = append(h.heap, n)
	h.percolateUp(h.indices[n])
}

func (h *lngheap) removeMin() int32 {
	x := h.heap[0]
	h.heap[0] = h.heap[len(h.heap)-1]
	h.indices[h.heap[0]] = 0
	h.indices[x] = -1
	h.heap = h.heap[:len(h.heap)-1]
	if len(h.heap) > 1 {
		h.percolateDown(0)
	}
	return x
}

func (h *lngheap) remove(n int32) {
	kPos := h.indices[n]
	h.indices[n] = -1
	if kPos < len(h.heap)-1 {
		h.heap[kPos] = h.heap[len(h.heap)-1]
		h.indices[h.heap[kPos]] = kPos
		h.heap = h.heap[:len(h.heap)-1]
		h.percolateDown(kPos)
	} else {
		h.heap = h.heap[:len(h.heap)-1]
	}
}

func (h *lngheap) build(ns []int32) {
	for i := 0; i < len(h.heap); i++ {
		h.indices[h.heap[i]] = -1
	}
	h.heap = []int32{}
	for i := 0; i < len(ns); i++ {
		h.indices[ns[i]] = i
		h.heap = append(h.heap, ns[i])
	}
	for i := len(h.heap)/2 - 1; i >= 0; i-- {
		h.percolateDown(i)
	}
}

func (h *lngheap) clear() {
	for i := 0; i < len(h.heap); i++ {
		h.indices[h.heap[i]] = -1
	}
	h.heap = []int32{}
}

func (h *lngheap) percolateUp(pos int) {
	x := h.heap[pos]
	p := h.parent(pos)
	j := pos

	for j != 0 && h.solver.lt(x, h.heap[p]) {
		h.heap[j] = h.heap[p]
		h.indices[h.heap[p]] = j
		j = p
		p = h.parent(p)
	}

	h.heap[j] = x
	h.indices[x] = j
}

func (h *lngheap) percolateDown(pos int) {
	p := pos
	y := h.heap[p]
	for h.left(p) < len(h.heap) {
		var child int
		if h.right(p) < len(h.heap) && h.solver.lt(h.heap[h.right(p)], h.heap[h.left(p)]) {
			child = h.right(p)
		} else {
			child = h.left(p)
		}
		if !h.solver.lt(h.heap[child], y) {
			break
		}
		h.heap[p] = h.heap[child]
		h.indices[h.heap[p]] = p
		p = child
	}
	h.heap[p] = y
	h.indices[y] = p
}

func (h *lngheap) String() string {
	var sb strings.Builder
	sb.WriteString("LNGHeap{")
	for i := 0; i < len(h.heap); i++ {
		sb.WriteString(fmt.Sprintf("[%d, %d]", h.heap[i], h.indices[i]))
		if i != len(h.heap)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (h *lngheap) growTo(size int) {
	if len(h.indices) >= size {
		return
	}
	numberNew := size - len(h.indices)
	for i := 0; i < numberNew; i++ {
		h.indices = append(h.indices, -1)
	}
}
