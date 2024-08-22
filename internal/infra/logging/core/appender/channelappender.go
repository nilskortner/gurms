package appender

import (
	"gurms/internal/infra/logging/core/model"
	"os"
)

type ChannelAppender struct {
	appender *BaseAppender
	File     *os.File
}

func NewChannelAppender(level model.LogLevel) *ChannelAppender {
	appender := NewAppender(level)
	return &ChannelAppender{
		appender: appender,
	}
}
