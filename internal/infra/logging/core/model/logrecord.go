package model

import (
	"bytes"
	"log"
)

type LogRecord struct {
	logger    *log.Logger
	level     LogLevel
	timestamp int64
	data      bytes.Buffer
}

func NewLogRecord(logger *log.Logger, level LogLevel, timestamp int64, data bytes.Buffer) LogRecord {
	return LogRecord{
		logger:    logger,
		level:     level,
		timestamp: timestamp,
		data:      data,
	}
}

func (l *LogRecord) Level() LogLevel {
	return l.level
}
