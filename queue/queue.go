package queue

import (
	"sync"
)

type node[T any] struct {
	value T
	next  *node[T]
}

type Queue[T any] struct {
	head *node[T]
	end  *node[T]

	count int
	mu    sync.Mutex
}

func (q *Queue[T]) Push(value T) {
	item := &node[T]{value: value, next: nil}

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.Empty() {
		q.head = item
		q.end = item
	} else {
		q.end.next = item
		q.end = item
	}

	q.count++
}

func (q *Queue[T]) Size() int {
	return q.count
}

func (q *Queue[T]) Empty() bool {
	return q.count == 0
}

func (q *Queue[T]) Pop() T {
	if q.Empty() {
		return *new(T)
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	res := q.head

	if q.head.next != nil {
		q.head = q.head.next
	} else {
		q.head = nil
		q.end = nil
	}

	q.count--

	return res.value
}

func (q *Queue[T]) Front() T {
	if q.Empty() {
		return *new(T)
	}

	return q.head.value
}

func (q *Queue[T]) Back() T {
	if q.Empty() {
		return *new(T)
	}

	return q.end.value
}
