package appender

import (
	"gurms/internal/infra/logging/core/model/loglevel"
	"os"
)

type ChannelConsoleAppender struct {
	appender *ChannelAppender
	output   *os.File
}

func NewChannelConsoleAppender(level loglevel.LogLevel) *ChannelConsoleAppender {
	return &ChannelConsoleAppender{
		appender: NewChannelAppender(level),
		output:   os.Stdout,
	}
}

func (c *ChannelConsoleAppender) Append() {}

func (c *ChannelConsoleAppender) GetLevel() loglevel.LogLevel {
	return c.appender.appender.GetLevel()
}
