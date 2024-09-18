package logging

import "gurms/internal/infra/logging/core/model/loglevel"

type ConsoleLoggingProperties struct {
	enabled bool
	level   loglevel.LogLevel
}

func NewConsoleLoggingProperties() *ConsoleLoggingProperties {
	return &ConsoleLoggingProperties{
		enabled: false,
		level:   2,
	}
}

func (c *ConsoleLoggingProperties) IsEnabled() bool {
	return c.enabled
}

func (c *ConsoleLoggingProperties) Level() loglevel.LogLevel {
	return c.level
}
