package datatype

type Set[T comparable] struct {
	Data map[T]bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{Data: make(map[T]bool)}
}

func (s *Set[T]) Add(element T) {
	s.Data[element] = true
}

func (s *Set[T]) Remove(element T) {
	delete(s.Data, element)
}

func (s *Set[T]) Contains(element T) bool {
	return s.Data[element]
}

func (s *Set[T]) Size() int {
	return len(s.Data)
}
func (s *Set[T]) IsEmpty() bool {
	return len(s.Data) == 0
}
func (s *Set[T]) Subtract(other *Set[T]) {
	for elem := range other.Data {
		if s.Contains(elem) {
			s.Remove(elem)
		}
	}
}
