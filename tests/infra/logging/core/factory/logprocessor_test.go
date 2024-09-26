package factory_test

import (
	"errors"
	"fmt"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/logging/core/model/logrecord"
	"gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"testing"
	"time"
)

func TestLogProcessor(t *testing.T) {

	list := make([]appender.Appender, 0)

	layoutAL := layout.NewGurmsTemplateLayout(2, "TURMS_CLUSTER_NODE_ID")

	mpsc := mpscunboundedarrayqueue.NewMpscUnboundedQueue[logrecord.LogRecord](1024)

	//chanAppender := appender.NewChannelAppender(0)
	//sysAppender := appender.NewSystemConsoleAppender(5)
	sysAppender2 := appender.NewSystemConsoleAppender(3)
	//list = append(list, chanAppender)
	//list = append(list, sysAppender)
	list = append(list, sysAppender2)

	var asyncL logger.Logger = logger.NewAsyncLogger("AsyncLggerNr1", true, list, layoutAL, mpsc)

	for i := 0; i < 1; i++ {
		err := errors.New("fox")
		err = fmt.Errorf("failed to complete task %w", err)
		asyncL.Fatal("error in function")
		asyncL.FatalWithError("peter", err)
		asyncL.Warn("warning")
	}

	processor := factory.NewLogProcessor(mpsc)
	processor.Start()

	time.Sleep(5 * time.Second)

	err := errors.New("Internal error")
	wErr := fmt.Errorf("error in creator: %w", err)
	wErr2 := fmt.Errorf("error with function create(): %w", wErr)

	asyncL.Error(wErr2)

	t.Error("test2")
}
