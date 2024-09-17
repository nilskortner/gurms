package model

import (
	"bytes"
	"log"
)

type LogRecord struct {
	logger    *logger.AsyncLogger
	level     LogLevel
	timestamp int64
	data      *bytes.Buffer
}

func NewLogRecord(logger *logger.AsyncLogger, level LogLevel, timestamp int64, data bytes.Buffer) LogRecord {
	return LogRecord{
		logger:    logger,
		level:     level,
		timestamp: timestamp,
		data:      &data,
	}
}

func (l *LogRecord) Level() LogLevel {
	return l.level
}

func (l *LogRecord) GetLogger() *log.Logger {
	return l.logger
}

func (l *LogRecord) ClearData() {
	l.data = nil
}

// func equals
