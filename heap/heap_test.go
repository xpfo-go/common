package heap

import (
	"testing"
)

func TestHeap(t *testing.T) {
	type node struct {
		a int
	}

	h := NewHeap[*node](nil, func(a, b *node) bool {
		return a.a < b.a
	})

	t.Logf("pop: %v", h.Pop())
	t.Logf("size: %v", h.Size())
	t.Logf("empty: %v", h.Empty())

	h.Push(&node{a: 12})
	h.Push(&node{a: -231})
	h.Push(&node{a: 122})
	h.Push(&node{a: 1})
	t.Logf("pop: %v", h.Pop())
	t.Logf("pop: %v", h.Pop())
	t.Logf("size: %v", h.Size())
	t.Logf("empty: %v", h.Empty())
	t.Logf("pop: %v", h.Pop())
	t.Logf("pop: %v", h.Pop())
	t.Logf("size: %v", h.Size())
	t.Logf("empty: %v", h.Empty())
	t.Logf("pop: %v", h.Pop())
	t.Logf("pop: %v", h.Pop())
	t.Logf("size: %v", h.Size())
	t.Logf("empty: %v", h.Empty())
}
