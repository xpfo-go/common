package queue

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := Queue[int]{}
	t.Logf("queue empty: %v", q.Empty())
	q.Push(1)
	t.Logf("queue lens: %v", q.Size())
	q.Push(2)
	t.Logf("queue lens: %v", q.Size())
	q.Push(3)
	t.Logf("queue lens: %v", q.Size())

	t.Logf("item: %v", q.Pop())

	t.Logf("item: %v", q.Front())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Back())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("queue empty: %v", q.Empty())

}

func TestQueueInt(t *testing.T) {
	q := Queue[int]{}
	q.Push(1)
	t.Logf("queue lens: %v", q.Size())
	q.Push(2)
	t.Logf("queue lens: %v", q.Size())
	q.Push(3)
	t.Logf("queue lens: %v", q.Size())

	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())

	t.Logf("queue lens: %v", q.Size())
	q.Push(1)
	t.Logf("queue lens: %v", q.Size())
	q.Push(2)
	t.Logf("queue lens: %v", q.Size())
	q.Push(3)
	t.Logf("queue lens: %v", q.Size())

	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())
}

func TestQueueString(t *testing.T) {
	q := Queue[string]{}
	q.Push("1")
	t.Logf("queue lens: %v", q.Size())
	q.Push("2")
	t.Logf("queue lens: %v", q.Size())
	q.Push("3")
	t.Logf("queue lens: %v", q.Size())

	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
}

func TestQueuePointer(t *testing.T) {

	q := Queue[*map[int]int]{}
	q.Push(&map[int]int{1: 1})
	t.Logf("queue lens: %v", q.Size())
	q.Push(&map[int]int{2: 2})
	t.Logf("queue lens: %v", q.Size())
	q.Push(&map[int]int{3: 3})
	t.Logf("queue lens: %v", q.Size())

	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
}

func TestQueueStruct(t *testing.T) {
	type St struct {
		a int
		b int
	}

	q := Queue[St]{}
	q.Push(St{a: 1, b: 1})
	t.Logf("queue lens: %v", q.Size())
	q.Push(St{a: 2, b: 2})
	t.Logf("queue lens: %v", q.Size())
	q.Push(St{a: 3, b: 3})
	t.Logf("queue lens: %v", q.Size())

	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())
	t.Logf("item: %v", q.Pop())
	t.Logf("queue lens: %v", q.Size())
}
