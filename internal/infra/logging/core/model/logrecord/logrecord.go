package logrecord

import (
	"bytes"
	"gurms/internal/infra/logging/core/model/loglevel"
)

type LogRecord struct {
	loggername string
	level      loglevel.LogLevel
	timestamp  int64
	data       *bytes.Buffer
}

func NewLogRecord(loggername string, level loglevel.LogLevel, timestamp int64, data *bytes.Buffer) LogRecord {
	return LogRecord{
		loggername: loggername,
		level:      level,
		timestamp:  timestamp,
		data:       data,
	}
}

func (l *LogRecord) Level() loglevel.LogLevel {
	return l.level
}

func (l *LogRecord) GetLogger() string {
	return l.loggername
}

func (l *LogRecord) ClearData() {
	l.data = nil
}
