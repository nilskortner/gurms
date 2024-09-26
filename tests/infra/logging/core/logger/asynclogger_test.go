package logger_test

import (
	"errors"
	"fmt"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/logging/core/model/logrecord"
	"gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"testing"
)

func TestAsyncLogger(t *testing.T) {
	list := make([]appender.Appender, 0)

	layoutAL := layout.NewGurmsTemplateLayout(2, "TURMS_CLUSTER_NODE_ID")

	mpsc := mpscunboundedarrayqueue.NewMpscUnboundedQueue[logrecord.LogRecord](1024)

	chanAppender := appender.NewChannelAppender(0)
	list = append(list, chanAppender)

	asyncL := logger.NewAsyncLogger("AsyncLggerNr1", true, list, layoutAL, mpsc)

	t.Log(mpsc.GetBuffer())

	asyncL.Log(5, "Frank Walter Steinmeier")

	value := asyncL.IsDebugEnabled()
	t.Logf("debug?: %v", value)

	value = asyncL.IsFatalEnabled()
	t.Logf("trace?: %v", value)

	buf := mpsc.GetBuffer()
	val := buf[0].Load()

	t.Log(val.GetBuffer())

	t.Error("test")
}

func TestAsyncLogger2(t *testing.T) {

	list := make([]appender.Appender, 0)

	layoutAL := layout.NewGurmsTemplateLayout(2, "TURMS_CLUSTER_NODE_ID")

	mpsc := mpscunboundedarrayqueue.NewMpscUnboundedQueue[logrecord.LogRecord](1024)

	chanAppender := appender.NewChannelAppender(0)
	list = append(list, chanAppender)

	var asyncL logger.Logger = logger.NewAsyncLogger("AsyncLggerNr1", true, list, layoutAL, mpsc)

	err := errors.New("Internal error")
	wErr := fmt.Errorf("error in creator: %w", err)
	wErr2 := fmt.Errorf("error with function create(): %w", wErr)

	asyncL.Error(wErr2)

	buf := mpsc.GetBuffer()
	val := buf[0].Load()

	t.Log(val.GetBuffer())

	t.Error("test2")
}
