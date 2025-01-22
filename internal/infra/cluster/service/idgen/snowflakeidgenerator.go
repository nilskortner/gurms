package idgen

import (
	"fmt"
	"sync/atomic"
)

const (
	EPOCH                int64 = 1602547200000
	TIMESTAMP_BITS       int   = 41
	DATA_CENTER_ID_BITS  int   = 4
	WORKER_ID_BITS       int   = 8
	SEQUENCE_NUMBER_BITS int   = 10

	TIMESTAMP_LEFT_SHIFT int64 = int64(SEQUENCE_NUMBER_BITS) + int64(WORKER_ID_BITS) + int64(DATA_CENTER_ID_BITS)
	DATA_CENTER_ID_SHIFT int64 = int64(SEQUENCE_NUMBER_BITS) + int64(WORKER_ID_BITS)
	WORKER_ID_SHIFT      int64 = int64(SEQUENCE_NUMBER_BITS)
	SEQUENCE_NUMBER_MASK int64 = (1 << SEQUENCE_NUMBER_BITS) - 1
	MAX_DATA_CENTER_ID   int   = 1 << DATA_CENTER_ID_BITS
	MAX_WORKER_ID        int   = 1 << WORKER_ID_BITS
)

type SnowflakeIdGenerator struct {
	lastTimestamp atomic.Int64
	// TODO: random init
	sequenceNumber atomic.Int64
	dataCenterId   int64
	workerId       int64
}

func NewSnowflakeIdGenerator(dataCenterId, workerId int) *SnowflakeIdGenerator {
	generator := &SnowflakeIdGenerator{}
	generator.updateNodeInfo(dataCenterId, workerId)
	return generator
}

// TODO: panic?
func (s *SnowflakeIdGenerator) updateNodeInfo(dataCenterId, workerId int) {
	if dataCenterId >= MAX_DATA_CENTER_ID {
		reason := fmt.Sprintf("the data center ID must be in the range: [0, %d], but got: %d", MAX_DATA_CENTER_ID, dataCenterId)
		panic(reason)
	}
	if workerId >= (1 << WORKER_ID_BITS) {
		reason := fmt.Sprintf("the worker ID must be in the range: [0, %d], but got: %d", 1<<WORKER_ID_BITS, workerId)
		panic(reason)
	}
	s.dataCenterId = int64(dataCenterId)
	s.workerId = int64(workerId)
}

func (s *SnowflakeIdGenerator) nextIncreasingId() int64 {
	sequenceNum := s.sequenceNumber.
}
