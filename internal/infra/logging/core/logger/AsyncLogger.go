package logger

import (
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/model"
	mpsc "gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"gurms/internal/supportpkgs/mathsupport"
)

type AsyncLogger struct {
	name        string
	shouldParse bool
	appenders   []appender.Appender
	layout      layout.GurmsTemplateLayout
	queue       mpsc.MpscUnboundedArrayQueue[model.LogRecord]
	nameForLog  []byte
	level       int
}

func NewAsyncLogger(
	name string,
	shouldParse bool,
	appenders []appender.Appender,
	layoutAL layout.GurmsTemplateLayout,
	queue mpsc.MpscUnboundedArrayQueue[model.LogRecord]) *AsyncLogger {
	nameForLog := layout.FormatStructName(name)

	var level int
	if len(appenders) == 0 {
		level = 
	} else {
		level =- 1

	}

	return &AsyncLogger{
		name:        name,
		shouldParse: shouldParse,
		appenders:   appenders,
		layout:      layoutAL,
		queue:       queue,
		level: level,
	}
}
