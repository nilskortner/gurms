package appender

import (
	"gurms/internal/infra/logging/core/model"
	"os"
)

type ChannelConsoleAppender struct {
	appender *ChannelAppender
	output   *os.File
}

func NewChannelConsoleAppender(level model.LogLevel) *ChannelConsoleAppender {
	return &ChannelConsoleAppender{
		appender: NewChannelAppender(level),
		output:   os.Stdout,
	}
}

func (c *ChannelConsoleAppender) Append() {}
