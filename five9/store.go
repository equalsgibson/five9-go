package five9

import "sync"

type store[T any] struct {
	data  T
	mutex *sync.Mutex
}

func newStore[T any](initialValue T) *store[T] {
	return &store[T]{
		data:  initialValue,
		mutex: &sync.Mutex{},
	}
}

func (s *store[T]) Reset() {
	s.data = T{}
}

func (s *store[T]) Write() {
	s.data = T{}
}
