package idgen

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"
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

func NewSnowflakeIdGenerator(dataCenterId, workerId int) (*SnowflakeIdGenerator, error) {
	generator := &SnowflakeIdGenerator{}
	err := generator.updateNodeInfo(dataCenterId, workerId)
	if err != nil {
		return generator, err
	}
	return generator, nil
}

func (s *SnowflakeIdGenerator) updateNodeInfo(dataCenterId, workerId int) error {
	if dataCenterId >= MAX_DATA_CENTER_ID {
		reason := fmt.Sprintf("the data center ID must be in the range: [0, %d], but got: %d", MAX_DATA_CENTER_ID, dataCenterId)
		return fmt.Errorf(reason)
	}
	if workerId >= (1 << WORKER_ID_BITS) {
		reason := fmt.Sprintf("the worker ID must be in the range: [0, %d], but got: %d", 1<<WORKER_ID_BITS, workerId)
		return fmt.Errorf(reason)
	}
	s.dataCenterId = int64(dataCenterId)
	s.workerId = int64(workerId)
	return nil
}

func (s *SnowflakeIdGenerator) nextIncreasingId() int64 {
	// prepare each ID part
	sequenceNum := s.sequenceNumber.Add(1) & SEQUENCE_NUMBER_MASK
	updateAndGet := func() int64 {
		for {
			stamp := s.lastTimestamp.Load()
			nonBackwardsTimestamp := int64(math.Max(float64(stamp), float64(time.Now().UnixMilli())))
			if sequenceNum == 0 {
				nonBackwardsTimestamp++
			}
			if s.lastTimestamp.CompareAndSwap(stamp, nonBackwardsTimestamp) {
				return nonBackwardsTimestamp
			}
		}
	}
	timestamp := updateAndGet() - EPOCH

	// Get ID
	return (timestamp << TIMESTAMP_LEFT_SHIFT) | (s.dataCenterId << DATA_CENTER_ID_SHIFT) |
		(s.workerId << WORKER_ID_SHIFT) | sequenceNum
}

func (s *SnowflakeIdGenerator) nextLargeGapId() int64 {
	// prepare each ID part
	sequenceNum := s.sequenceNumber.Add(1) & SEQUENCE_NUMBER_MASK
	updateAndGet := func() int64 {
		for {
			stamp := s.lastTimestamp.Load()
			nonBackwardsTimestamp := int64(math.Max(float64(stamp), float64(time.Now().UnixMilli())))
			if sequenceNum == 0 {
				nonBackwardsTimestamp++
			}
			if s.lastTimestamp.CompareAndSwap(stamp, nonBackwardsTimestamp) {
				return nonBackwardsTimestamp
			}
		}
	}
	timestamp := updateAndGet() - EPOCH

	// Get ID
	return (sequenceNum << (int64(TIMESTAMP_BITS + DATA_CENTER_ID_BITS + WORKER_ID_BITS))) |
		(timestamp << (int64(DATA_CENTER_ID_BITS + WORKER_ID_BITS))) |
		s.dataCenterId<<int64(WORKER_ID_BITS) | s.workerId
}
