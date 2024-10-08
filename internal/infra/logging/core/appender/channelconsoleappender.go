package appender

import (
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/infra/logging/core/model/logrecord"
)

type ChannelConsoleAppender struct {
	appender *ChannelAppender
}

func NewChannelConsoleAppender(level loglevel.LogLevel) *ChannelConsoleAppender {
	return &ChannelConsoleAppender{
		appender: NewChannelAppender(level),
	}
}

func (c *ChannelConsoleAppender) Append(record logrecord.LogRecord) int {
	return c.appender.Append(record)
}

func (c *ChannelConsoleAppender) GetLevel() loglevel.LogLevel {
	return c.appender.appender.GetLevel()
}
