package appender

import "gurms/internal/infra/logging/core/model"

type BaseAppender struct {
	level model.LogLevel
}

type Appender interface {
	Append()
}

func NewAppender(level model.LogLevel) *BaseAppender {
	return &BaseAppender{
		level: level,
	}
}
