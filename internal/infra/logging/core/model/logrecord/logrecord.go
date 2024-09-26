package logrecord

import (
	"bytes"
	"gurms/internal/infra/logging/core/model/loglevel"
)

type LogRecord struct {
	logger    any
	level     loglevel.LogLevel
	timestamp int64
	data      *bytes.Buffer
}

// logger needs to be of type *AsyncLogger
//
// *bytes.Buffer cant be nil
func NewLogRecord(logger any, level loglevel.LogLevel, timestamp int64, data *bytes.Buffer) LogRecord {
	if data == nil {
		panic("nil pointer in NewLogRecord")
	}
	return LogRecord{
		logger:    logger,
		level:     level,
		timestamp: timestamp,
		data:      data,
	}
}

func (l *LogRecord) Level() loglevel.LogLevel {
	return l.level
}

func (l *LogRecord) Timestamp() int64 {
	return l.timestamp
}

func (l *LogRecord) GetLogger() any {
	return l.logger
}

func (l *LogRecord) GetBuffer() *bytes.Buffer {
	return l.data
}

func (l *LogRecord) ClearData() {
	l.data = nil
}
