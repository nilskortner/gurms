package application

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"time"
)

var SHUTDOWNMANAGERLOGGER logger.Logger = factory.GetLogger("ShutDownManager")

type ShutDownManager struct {
	isClosing            []*bool
	closeRoutineChannels []chan struct{}
}

func (s *ShutDownManager) ShutdownApp() {
	s.StartCloseRoutines()
	s.StartCloseRessources()
}

func (s *ShutDownManager) StartCloseRoutines() {
	SHUTDOWNMANAGERLOGGER.Warn("started shutting down ALL go routines")
	for _, setClose := range s.isClosing {
		*setClose = true
	}
	time.Sleep(5 * time.Second)
	SHUTDOWNMANAGERLOGGER.Warn("shutting down ALL go routines")
	s.closeRoutines()
	SHUTDOWNMANAGERLOGGER.InfoWithArgs("go routines have been shut down")
}

func (s *ShutDownManager) closeRoutines() {
	for _, channel := range s.closeRoutineChannels {
		channel <- struct{}{}
	}
}

// TODO:
func (s *ShutDownManager) StartCloseRessources() {

}
