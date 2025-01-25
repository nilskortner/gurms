package healthcheck

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var HCMLOGGER logger.Logger = factory.GetLogger("HealthCheckManager")

type Node interface {
}

type HealthCheckManager struct {
	node                Node
	cpuHealthChecker    *CpuHealthChecker
	memoryHealthChecker *MemoryHealthChecker
	lastUpdateTimestamp int64
}

func NewHealthCheckManager(node Node, propertiesManager *GurmsPropertiesManager) {

}
