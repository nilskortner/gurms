package appender

import (
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/infra/logging/core/model/logrecord"
)

type SystemConsoleAppender struct {
	appender BaseAppender
}

func NewSystemConsoleAppender(level loglevel.LogLevel) *SystemConsoleAppender {
	return &SystemConsoleAppender{
		appender: *NewAppender(level),
	}
}

func (s *SystemConsoleAppender) Append(logrecord.LogRecord) {
}

func (s *SystemConsoleAppender) GetLevel() loglevel.LogLevel {
	return s.appender.GetLevel()
}

// func append(record model.LogRecord) int {
// 	if !record.Level().IsLoggable(appenderlevel) {
// 		return 0
// 	}
// 	s := ByteBufUtil.GetString(record.data())

// }
