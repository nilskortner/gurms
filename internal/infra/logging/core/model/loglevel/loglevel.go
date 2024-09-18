package loglevel

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

func (level LogLevel) IsLoggable(enabledLevel LogLevel) bool {
	return enabledLevel <= level
}

func (level LogLevel) IsErrorOrFatal() bool {
	return level == ERROR || level == FATAL
}

func (level LogLevel) String() string {
	switch level {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
