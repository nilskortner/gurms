package appender

import "gurms/internal/infra/logging/core/model"

type Appender struct {
	level model.LogLevel
}

func NewAppender(level model.LogLevel) *Appender {
	return &Appender{
		level: level,
	}
}
