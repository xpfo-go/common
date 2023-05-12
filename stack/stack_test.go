package stack

import "testing"

func TestStack(t *testing.T) {
	s := Stack[*map[int]int]{}

	t.Logf("empty: %v", s.Empty())

	s.Push(&map[int]int{1: 1})
	s.Push(&map[int]int{2: 2})
	s.Push(&map[int]int{3: 3})

	t.Logf("top : %v", s.Top())
	t.Logf("bottom : %v", s.Bottom())

	t.Logf("empty: %v", s.Empty())

	t.Logf("size: %v", s.Size())
	t.Logf("pop : %v", s.Pop())
	t.Logf("size: %v", s.Size())
	t.Logf("pop : %v", s.Pop())
	t.Logf("size: %v", s.Size())
	t.Logf("pop : %v", s.Pop())
	t.Logf("size: %v", s.Size())
	t.Logf("pop : %v", s.Pop())
	t.Logf("size: %v", s.Size())
	t.Logf("pop : %v", s.Pop())
	t.Logf("size: %v", s.Size())

}

func TestStack_Push(t *testing.T) {
	s := Stack[int]{}
	s.Push(1)
	t.Log(s.Size())
	s.Push(2)
	t.Log(s.Size())
}

func TestStack_Pop(t *testing.T) {
	s := Stack[int]{}
	s.Push(1)
	s.Push(2)
	t.Log(s.Pop())
	t.Log(s.Pop())
	t.Log(s.Pop())
	t.Log(s.Pop())
}

func TestStack_Top(t *testing.T) {
	s := Stack[int]{}
	s.Push(1)
	s.Push(2)
	t.Log(s.Top())
	s.Pop()
	t.Log(s.Top())
	s.Pop()
	s.Pop()
	t.Log(s.Top())
}

func TestStack_Empty(t *testing.T) {
	s := Stack[int]{}
	t.Log(s.Empty())
	s.Push(1)
	t.Log(s.Empty())
	s.Pop()
	t.Log(s.Empty())
}

func TestStack_Bottom(t *testing.T) {
	s := Stack[int]{}
	s.Push(1)
	t.Log(s.Bottom())
	s.Push(2)
	t.Log(s.Bottom())
	s.Pop()
	t.Log(s.Bottom())
	s.Pop()
	t.Log(s.Bottom())
}

func TestStack_Size(t *testing.T) {
	s := Stack[int]{}
	s.Push(1)
	t.Log(s.Size())
	s.Push(2)
	t.Log(s.Size())
	s.Pop()
	t.Log(s.Size())
	s.Pop()
	t.Log(s.Size())
	s.Pop()
	t.Log(s.Size())
}
