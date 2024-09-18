package logger

import (
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/model/logrecord"
	mpsc "gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"gurms/internal/supportpkgs/mathsupport"
	"math"
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
