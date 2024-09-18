package processor

import (
	"gurms/internal/infra/logging/core/idle"
	"gurms/internal/infra/logging/core/model/logrecord"
	"gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	mpsc "gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"sync"
	"sync/atomic"
)

type LogProcessor struct {
	active bool
	wait   *sync.WaitGroup
	count  int32
	queue  *mpscunboundedarrayqueue.MpscUnboundedArrayQueue[logrecord.LogRecord]
}

func NewLogProcessor(queue *mpscunboundedarrayqueue.MpscUnboundedArrayQueue[logrecord.LogRecord]) *LogProcessor {
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

func (lp *LogProcessor) drainLogsForever(recordQueue mpsc.MpscUnboundedArrayQueue[logrecord.LogRecord]) {
	idleStrategy := idle.NewBackoffIdleStrategy(128, 128, 1024000, 1024000)
	var logRecord logrecord.LogRecord
	var success bool
	for {
		for {
			logRecord, success = recordQueue.RelaxedPoll()
			if success == false {
				break
			}
			idleStrategy.Reset()
			appenders := logRecord.GetAppenders()
			// for appender: appenders {
			// 	appender.append(logRecord)
			// }
			logRecord.ClearData()
		}
		if !lp.active {
			break
		}
		idleStrategy.Idle()
	}
	//for
}
