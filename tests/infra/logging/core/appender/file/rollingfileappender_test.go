package file_test

import (
	"bytes"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/appender/file"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/logging/core/model/logrecord"
	"gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"testing"
	"time"
)

func TestRollingFileAppender(t *testing.T) {
	rf := file.NewRollingFileAppender(0,
		"@HOME/log/.log",
		320,
		32,
		false)

	// define asynclogger
	list := make([]appender.Appender, 0)

	layoutAL := layout.NewGurmsTemplateLayout(2, "TURMS_CLUSTER_NODE_ID")

	mpsc := mpscunboundedarrayqueue.NewMpscUnboundedQueue[logrecord.LogRecord](1024)

	chanAppender := appender.NewChannelAppender(0)
	list = append(list, chanAppender)

	var asyncL logger.Logger = logger.NewAsyncLogger("AsyncLggerNr1", true, list, layoutAL, mpsc)
	// end

	for count := 0; count < 320; count++ {
		buffer := bytes.NewBuffer(make([]byte, 0))
		buffer.Write([]byte("testingString "))
		record := logrecord.NewLogRecord(asyncL, 0, time.Now().UnixMilli(), buffer)
		rf.Append(record)
	}
	t.Log(rf)

	t.Error("rfa test")
}
