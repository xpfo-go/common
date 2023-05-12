package heap

import "container/heap"

// NewHeap 对外只开放这一个函数，防止使用不当造成的空指针异常
func NewHeap[T any](v []T, comparator func(a, b T) bool) *pHeap[T] {
	if comparator == nil {
		panic("comparator can not be nil")
	}

	th := &pHeap[T]{heap: &h[T]{comparator: comparator}}
	if v != nil && len(v) != 0 {
		th.heap.values = v
	}

	heap.Init(th.heap)
	return th
}

type pHeap[T any] struct {
	heap *h[T]
}

func (h *pHeap[T]) Push(v T) {
	heap.Push(h.heap, v)
}

func (h *pHeap[T]) Pop() T {
	if h.Empty() {
		return *new(T)
	}
	return heap.Pop(h.heap).(T)
}

func (h *pHeap[T]) Size() int {
	return h.heap.Len()
}

func (h *pHeap[T]) Empty() bool {
	return h.Size() == 0
}

// 实现 container/heap 中得 Interface 接口, 可以使用go内置得堆操作函数
type h[T any] struct {
	values     []T
	comparator func(a, b T) bool
}

func (h *h[T]) Len() int {
	return len(h.values)
}

func (h *h[T]) Less(i, j int) bool {
	return h.comparator(h.values[i], h.values[j])
}

func (h *h[T]) Swap(i, j int) {
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

func (h *h[T]) Push(x any) {
	h.values = append(h.values, x.(T))
}

func (h *h[T]) Pop() (v any) {
	h.values, v = h.values[:h.Len()-1], h.values[h.Len()-1]
	return
}
