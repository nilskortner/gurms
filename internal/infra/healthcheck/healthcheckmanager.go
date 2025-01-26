package healthcheck

import (
	"fmt"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property"
	"time"
)

var HEALTHCHECKMANAGERLOGGER logger.Logger = factory.GetLogger("HealthCheckManager")

type Node interface {
	UpdateHealthStatus(bool)
}

type HealthCheckManager struct {
	node                Node
	cpuHealthChecker    *CpuHealthChecker
	memoryHealthChecker *MemoryHealthChecker
	lastUpdateTimestamp int64
}

func NewHealthCheckManager(node Node,
	propertiesManager *property.GurmsPropertiesManager) *HealthCheckManager {

	properties := propertiesManager.LocalGurmsProperties.HealthCheck
	healthCheckManager := &HealthCheckManager{
		node:                node,
		cpuHealthChecker:    NewCpuHealtChecker(properties.Cpu),
		memoryHealthChecker: NewMemoryHealthChecker(properties.Memory),
	}
	healthCheckManager.startHealthCheck(properties.CheckIntervalSeconds)
	return healthCheckManager
}

func (h *HealthCheckManager) isHealthy() bool {
	return h.cpuHealthChecker.isCpuHealthy && h.memoryHealthChecker.isMemoryHealthy
}

func (h *HealthCheckManager) startHealthCheck(intervalSeconds int) {
	intervalMillis := int64(intervalSeconds) * 1000
	interval := time.Duration(intervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	go func() {
		for {
			h.cpuHealthChecker.updateHealthStatus()
			h.memoryHealthChecker.updateHealthStatus()
			h.node.UpdateHealthStatus(h.isHealthy())

			// TODO: correct schedueling
			now := time.Now()
			previousUpdateTimeStamp := h.lastUpdateTimestamp + intervalMillis
			h.lastUpdateTimestamp = now.UnixMilli()
			if previousUpdateTimeStamp >= now.UnixMilli() {
				HEALTHCHECKMANAGERLOGGER.WarnWithArgs(fmt.Sprintf(
					"the system time goes backwards. the time drift is %d millis",
					previousUpdateTimeStamp-now.UnixMilli()))
			}
		}
	}()
}
