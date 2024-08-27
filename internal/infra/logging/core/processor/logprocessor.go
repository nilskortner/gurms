package processor

import (
	"gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"sync"
	"sync/atomic"
)

type LogProcessor struct {
	active bool
	wait   *sync.WaitGroup
	count  int32
	queue  *mpscunboundedarrayqueue.MpscUnboundedArrayQueue
}

func NewLogProcessor(queue *mpscunboundedarrayqueue.MpscUnboundedArrayQueue) *LogProcessor {
	return &LogProcessor{
		active: true,
		wait:   &sync.WaitGroup{},
		queue:  queue,
	}
}

func (lp *LogProcessor) Start() {
	if atomic.LoadInt32(&lp.count) == 0 {
		lp.wait.Add(1)
		atomic.AddInt32(&lp.count, 1)
		go lp.drainLogsForever()
	}
}

func (lp *LogProcessor) waitClose(timeoutMillis int64) {
	lp.active = false
}

func (lp *LogProcessor) drainLogsForever() {

}
