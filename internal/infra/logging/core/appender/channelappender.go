package appender

import "gurms/internal/infra/logging/core/model"

type ChannelAppender struct {
	appender *BaseAppender
}

func NewChannelAppender(level model.LogLevel) *ChannelAppender {
	appender := NewAppender(level)
	return &ChannelAppender{
		appender: appender,
	}
}
