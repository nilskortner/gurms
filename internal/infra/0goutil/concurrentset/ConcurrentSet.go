package concurrentset

import "sync"

type ConcurrentSet struct {
	m sync.Map
}

func NewConcurrentSet() *ConcurrentSet {
	return &ConcurrentSet{}
}

func (s *ConcurrentSet) Add(item any) {
	s.m.Store(item, struct{}{})
}

func (s *ConcurrentSet) Contains(item any) bool {
	_, ok := s.m.Load(item)
	return ok
}

func (s *ConcurrentSet) Size() int {
	size := 0
	s.m.Range(func(_, _ any) bool {
		size++
		return true
	})
	return size
}

func (s *ConcurrentSet) Items() []any {
	var items []any
	s.m.Range(func(key, _ any) bool {
		items = append(items, key)
		return true
	})
	return items
}
