package logging

import "gurms/internal/infra/logging/core/model/loglevel"

const CONSOLE_DEFAULT_VALUE_ENABLED = false
const CONSOLE_DEFAULT_VALUE_LEVEL loglevel.LogLevel = 1

type ConsoleLoggingProperties struct {
	enabled bool
	level   loglevel.LogLevel
}

func NewConsoleLoggingProperties(enabled bool, level loglevel.LogLevel) *ConsoleLoggingProperties {
	return &ConsoleLoggingProperties{
		enabled: enabled,
		level:   level,
	}
}

func (c *ConsoleLoggingProperties) IsEnabled() bool {
	return c.enabled
}

func (c *ConsoleLoggingProperties) Level() loglevel.LogLevel {
	return c.level
}
