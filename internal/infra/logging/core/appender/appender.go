package appender

import (
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/infra/logging/core/model/logrecord"
)

type BaseAppender struct {
	level loglevel.LogLevel
}

type Appender interface {
	Append(logrecord.LogRecord) int
	GetLevel() loglevel.LogLevel
}

func NewAppender(level loglevel.LogLevel) *BaseAppender {
	return &BaseAppender{
		level: level,
	}
}

func (a BaseAppender) GetLevel() loglevel.LogLevel {
	return a.level
}
