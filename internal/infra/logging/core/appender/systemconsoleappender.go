package appender

import (
	"gurms/internal/infra/logging/core/model"
)

type SystemConsoleAppender struct {
	appender BaseAppender
}

func NewSystemConsoleAppender(level model.LogLevel) *SystemConsoleAppender {
	return &SystemConsoleAppender{
		appender: *NewAppender(level),
	}
}

func (s *SystemConsoleAppender) Append() {
}

// func append(record model.LogRecord) int {
// 	if !record.Level().IsLoggable(appenderlevel) {
// 		return 0
// 	}
// 	s := ByteBufUtil.GetString(record.data())

// }
