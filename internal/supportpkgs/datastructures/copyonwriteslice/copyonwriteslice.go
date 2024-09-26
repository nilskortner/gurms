package copyonwriteslice

import (
	"sync"
)

type CopyOnWriteSlice[T any] struct {
	mu   sync.RWMutex
	data []T
}

func NewCopyOnWriteSlice[T any]() *CopyOnWriteSlice[T] {
	return &CopyOnWriteSlice[T]{}
}

func (c *CopyOnWriteSlice[T]) Add(value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	length := len(c.data)
	newData := make([]T, length+1)
	copy(newData, c.data)
	newData[length] = value
	c.data = newData
}

func (c *CopyOnWriteSlice[T]) List() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	newData := make([]T, len(c.data))
	copy(newData, c.data)
	return newData
}
