package healthcheck

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var HCMLOGGER logger.Logger = factory.GetLogger("HealthCheckManager")

type HealthCheckManager struct {
	node                node.Node
	cpuHealthChecker    *CpuHealthChecker
	memoryHealthChecker *MemoryHealthChecker
	lastUpdateTimestamp int64
}

func NewHealthCheckManager(node node.Node, propertiesManager *GurmsPropertiesManager) {

}
