package logrecord

import (
	"bytes"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/model/loglevel"
)

type LogRecord struct {
	appenders []appender.Appender
	level     loglevel.LogLevel
	timestamp int64
	data      *bytes.Buffer
}

func NewLogRecord(appenders []appender.Appender, level loglevel.LogLevel, timestamp int64, data bytes.Buffer) LogRecord {
	return LogRecord{
		appenders: appenders,
		level:     level,
		timestamp: timestamp,
		data:      &data,
	}
}

func (l *LogRecord) Level() loglevel.LogLevel {
	return l.level
}

func (l *LogRecord) GetAppender() []appender.Appender {
	return l.appenders
}

func (l *LogRecord) ClearData() {
	l.data = nil
}

// func equals
