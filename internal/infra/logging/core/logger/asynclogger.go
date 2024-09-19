package logger

import (
	"bytes"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/infra/logging/core/model/logrecord"
	mpsc "gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"gurms/internal/supportpkgs/mathsupport"
	"math"
	"time"
)

type AsyncLogger struct {
	name        string
	shouldParse bool
	appenders   []appender.Appender
	layout      layout.GurmsTemplateLayout
	queue       *mpsc.MpscUnboundedArrayQueue[logrecord.LogRecord]
	nameForLog  []byte
	level       int
}

func NewAsyncLogger(
	name string,
	shouldParse bool,
	appenders []appender.Appender,
	layoutAL layout.GurmsTemplateLayout,
	queue *mpsc.MpscUnboundedArrayQueue[logrecord.LogRecord]) *AsyncLogger {
	nameForLog := layout.FormatStructName(name)

	var level int
	if len(appenders) == 0 {
		level = math.MaxInt
	} else {
		level = -1
		for _, appender := range appenders {
			level = mathsupport.Max(level, int(appender.GetLevel()))
		}
	}

	return &AsyncLogger{
		name:        name,
		shouldParse: shouldParse,
		appenders:   appenders,
		layout:      layoutAL,
		queue:       queue,
		nameForLog:  nameForLog,
		level:       level,
	}
}

func (a *AsyncLogger) IsTraceEnabled() bool {
	return a.level <= int(loglevel.TRACE)
}

func (a *AsyncLogger) IsDebugEnabled() bool {
	return a.level <= int(loglevel.DEBUG)
}

func (a *AsyncLogger) IsInfoEnabled() bool {
	return a.level <= int(loglevel.INFO)
}

func (a *AsyncLogger) IsWarnEnabled() bool {
	return a.level <= int(loglevel.WARN)
}

func (a *AsyncLogger) IsErrorEnabled() bool {
	return a.level <= int(loglevel.ERROR)
}

func (a *AsyncLogger) IsFatalEnabled() bool {
	return a.level <= int(loglevel.FATAL)
}

func (a *AsyncLogger) IsEnabled(loglevel loglevel.LogLevel) bool {
	return a.level <= int(loglevel)
}

func (a *AsyncLogger) Log(level loglevel.LogLevel, message string) {
	if !a.IsEnabled(level) {
		return
	}
	doLog(level, message, 0, 0)
}

func (a *AsyncLogger) doLog(level loglevel.LogLevel, message string) {
	var buffer *bytes.Buffer

	buffer := layout.Format(a.shouldParse, a.nameForLog, a.level, message, a.layout)
	a.queue.Offer(logrecord.NewLogRecord(a.name, level, time.Now().UnixMilli(), buffer))

}
