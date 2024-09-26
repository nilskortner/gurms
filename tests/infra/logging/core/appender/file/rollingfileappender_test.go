package file_test

import (
	"gurms/internal/infra/logging/core/appender/file"
	"testing"
)

func TestRollingFileAppender(t *testing.T) {
	rf := file.NewRollingFileAppender(5,
		"@HOME/log/.log",
		320,
		32,
		false)

	t.Log(rf)

	t.Error("rfa test")
}
