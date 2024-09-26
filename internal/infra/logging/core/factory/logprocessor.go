package factory

import (
	"gurms/internal/infra/logging/core/idle"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/logging/core/model/logrecord"
	mpsc "gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"sync/atomic"
)

type LogProcessor struct {
	active bool
	wait   atomic.Int64
	queue  *mpsc.MpscUnboundedArrayQueue[logrecord.LogRecord]
}

func NewLogProcessor(queue *mpsc.MpscUnboundedArrayQueue[logrecord.LogRecord]) *LogProcessor {
	return &LogProcessor{
		active: true,
		queue:  queue,
	}
}

func (lp *LogProcessor) Start() {
	if lp.wait.Load() == 0 {
		lp.wait.Add(1)
		go lp.drainLogsForever(*lp.queue)
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
			if !success {
				break
			}
			idleStrategy.Reset()
			logger := logRecord.GetLogger().(*logger.AsyncLogger)
			appenders := logger.GetAppenders()
			for _, appender := range appenders {
				appender.Append(logRecord)
			}
			logRecord.ClearData()
		}
		if !lp.active {
			break
		}
		idleStrategy.Idle()
	}
}
