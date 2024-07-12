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

func (c *ConsoleLoggingProperties) IsEnabled() bool {
	return c.enabled
}

func (c *ConsoleLoggingProperties) Level() model.LogLevel {
	return c.level
}
