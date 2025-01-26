package healthcheck

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var MEMORYHEALTHCHECKER logger.Logger = factory.GetLogger("MemoryHealthChecker")

type MemoryHealthChecker struct {
	isMemoryHealthy bool
	unhealthyReason string
}

func NewMemoryHealthChecker() *MemoryHealthChecker {

}
