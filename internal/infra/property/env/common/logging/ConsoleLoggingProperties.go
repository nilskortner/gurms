package logging

import "gurms/internal/infra/logging/core/model"

type ConsoleLoggingProperties struct {
	enabled bool
	level   model.LogLevel
}

func NewConsoleLoggingProperties() *ConsoleLoggingProperties {
	return &ConsoleLoggingProperties{
		enabled: false,
		level:   "INFO",
	}
}
