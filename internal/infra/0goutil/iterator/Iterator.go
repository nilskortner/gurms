package goiterator

type Iterator interface {
	hasNext() bool
	Next() interface{}
}

type IntIterator struct {
	data []interface{}
	pos  int
}

func NewIntIterator(data ...interface{}) *IntIterator {
	return &IntIterator{data: data, pos: 0}
}

func (s *IntIterator) hasNext() bool {
	return s.pos < len(s.data)
}

func (s *IntIterator) next() interface{} {
	if s.hasNext() {
		value := s.data[s.pos]
		s.pos++
		return value
	}
	return nil
}
