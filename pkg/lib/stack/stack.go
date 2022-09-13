package stack

type Stack[T any] []T

func New[T any]() *Stack[T] {
	s := make(Stack[T], 0)
	return &s
}

func (s *Stack[T]) Push(elem T) {
	*s = append(*s, elem)
}

func (s *Stack[T]) Peek() *T {
	if s.Len() == 0 {
		return nil
	}
	elem := (*s)[len(*s)-1]
	return &elem
}

func (s *Stack[T]) Pop() *T {
	if s.Len() == 0 {
		return nil
	}
	elem := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return &elem
}

func (s *Stack[T]) Len() int {
	return len(*s)
}

func (s *Stack[T]) IsEmpty() bool {
	return s.Len() == 0
}
