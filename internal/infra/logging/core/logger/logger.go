package logger

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

type Logger interface {
	isTraceEnabled() bool

	isDebugEnabled() bool

	isInfoEnabled() bool

	isWarnEnabled() bool

	isErrorEnabled() bool

	isFatalEnabled() bool

	IsEnabled() bool

	Log(level LogLevel, message string)
	Logf(level LogLevel, format string, args ...interface{})
	LogError(level LogLevel, message string, err error)

	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	InfoBytes(data []byte)
	Warn(message string, args ...interface{})
	Error(err error)
	Errorf(message string, err error, args ...interface{})
	Fatal(message string, args ...interface{})
	FatalError(message string, err error)
}
