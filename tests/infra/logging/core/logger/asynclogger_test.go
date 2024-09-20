package logger_test

import (
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/logging/core/model/logrecord"
	"gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"testing"
)

func TestAsyncLogger(t *testing.T) {
	list := make([]appender.Appender, 0)
	//var level loglevel.LogLevel = 5
	//app := appender.NewChannelAppender(level)
	//list = append(list, app)

	layoutAL := layout.NewGurmsTemplateLayout(0, "TURMS_CLUSTER_NODE_ID")

	mpsc := mpscunboundedarrayqueue.NewMpscUnboundedQueue[logrecord.LogRecord](1024)

	asyncL := logger.NewAsyncLogger("AsyncLggerNr1", true, list, layoutAL, mpsc)

	value := asyncL.IsDebugEnabled()
	t.Logf("debug?: %v", value)

	value = asyncL.IsTraceEnabled()
	t.Logf("trace?: %v", value)

	t.Error("test")
}
