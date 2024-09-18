package appender

import (
	"gurms/internal/infra/logging/core/model/loglevel"
	"os"
)

type ChannelAppender struct {
	appender *BaseAppender
	File     *os.File
}

func NewChannelAppender(level loglevel.LogLevel) *ChannelAppender {
	appender := NewAppender(level)
	return &ChannelAppender{
		appender: appender,
	}
}

func (c *ChannelAppender) GetLevel() loglevel.LogLevel {
	return c.appender.GetLevel()
}
