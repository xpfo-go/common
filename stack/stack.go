package stack

import "sync"

type node[T any] struct {
	value T
	next  *node[T]
}

type Stack[T any] struct {
	top    *node[T]
	bottom *node[T]

	count int
	mu    sync.Mutex
}

func (s *Stack[T]) Push(value T) {
	item := &node[T]{value: value, next: nil}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Empty() {
		s.top = item
		s.bottom = item
	} else {
		item.next = s.top
		s.top = item
	}

	s.count++
}

func (s *Stack[T]) Pop() T {
	if s.Empty() {
		return *new(T)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	res := s.top

	if res.next == nil {
		s.top = nil
		s.bottom = nil
	} else {
		s.top = res.next
	}

	s.count--

	return res.value
}

func (s *Stack[T]) Size() int {
	return s.count
}

func (s *Stack[T]) Empty() bool {
	return s.count == 0
}

func (s *Stack[T]) Top() T {
	if s.Empty() {
		return *new(T)
	}
	return s.top.value
}

func (s *Stack[T]) Bottom() T {
	if s.Empty() {
		return *new(T)
	}
	return s.bottom.value
}
