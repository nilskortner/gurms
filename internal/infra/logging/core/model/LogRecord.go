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

func (l *LogRecord) Level() LogLevel {
	return l.level
}
