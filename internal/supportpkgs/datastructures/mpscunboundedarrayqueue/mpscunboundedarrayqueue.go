package mpscunboundedarrayqueue

import (
	"math"
)

type MpscUnboundedArrayQueue[T comparable] struct {
	bmlaq *BaseMpscLinkedArrayQueue[T]
}

func NewMpscUnboundedQueue[T comparable](chunkSize int) *MpscUnboundedArrayQueue[T] {
	return &MpscUnboundedArrayQueue[T]{
		bmlaq: NewBaseMpscLinkedArrayQueue[T](chunkSize),
	}
}

func availableInQueue(pIndex, cIndex int64) int64 {
	_, _ = pIndex, cIndex // unused for this implementation
	return math.MaxInt64
}

func getCurrentBufferCapacity(mask int64) int64 {
	return mask
}
