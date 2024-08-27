package mpscunboundedarrayqueue

import (
	"gurms/internal/infra/logging/core/model"
	"sync"
)

type Node struct {
	data  []model.LogRecord
	next  *Node
	index uint64
}

type MpscUnboundedArrayQueue struct {
	head *Node
	tail *Node
	lock sync.Mutex
}

func NewMpscUnboundedQueue() *MpscUnboundedArrayQueue {
	node := &Node{
		data: make([]model.LogRecord, 64),
	}
	return &MpscUnboundedArrayQueue{
		head: node,
		tail: node,
	}
}
