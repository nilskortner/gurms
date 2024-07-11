package copyonwriteslice

import (
	"gurms/internal/infra/logging/core/appender"
	"sync"
)

type CopyOnWriteSliceAppender struct {
	mu   sync.RWMutex
	data []appender.Appender
}

func NewCopyOnWriteSliceAppender() *CopyOnWriteSliceAppender {
	return &CopyOnWriteSliceAppender{}
}
