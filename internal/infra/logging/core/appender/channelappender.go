package appender

import (
	"fmt"
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/infra/logging/core/model/logrecord"
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

func (c *ChannelAppender) Close() {
	err := c.File.Close()
	if err != nil {
		fmt.Println("internal logger: channelappender.close() ", err)
	}
}

func (c *ChannelAppender) GetLevel() loglevel.LogLevel {
	return c.appender.GetLevel()
}

func (c *ChannelAppender) Append(record logrecord.LogRecord) int {
	if !record.Level().IsLoggable(c.appender.level) {
		return 0
	}
	buffer := record.GetBuffer()

	fmt.Println(buffer.String())

	return buffer.Len()
}
