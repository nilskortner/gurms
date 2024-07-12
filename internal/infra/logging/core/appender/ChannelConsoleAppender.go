package appender

import (
	"gurms/internal/infra/logging/core/model"
)

type ChannelConsoleAppender struct {
	appender Appender
}

func NewChannelConsoleAppender(level model.LogLevel) *ChannelConsoleAppender {
	return &ChannelConsoleAppender{
		appender: *NewAppender(level),
	}
}
