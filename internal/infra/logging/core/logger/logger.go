package logger

import model "gurms/internal/infra/logging/core/model/loglevel"

type Logger interface {
	isTraceEnabled() bool

	isDebugEnabled() bool

	isInfoEnabled() bool

	isWarnEnabled() bool

	isErrorEnabled() bool

	isFatalEnabled() bool

	IsEnabled() bool

	Log(level model.LogLevel, message string)
	Logf(level model.LogLevel, format string, args ...interface{})
	LogError(level model.LogLevel, message string, err error)

	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	InfoBytes(data []byte)
	Warn(message string, args ...interface{})
	Error(err error)
	Errorf(message string, err error, args ...interface{})
	Fatal(message string, args ...interface{})
	FatalError(message string, err error)
}
