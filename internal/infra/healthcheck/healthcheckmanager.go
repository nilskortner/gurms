package healthcheck

import (
	"context"
	"fmt"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property"
	"time"
)

var HEALTHCHECKMANAGERLOGGER logger.Logger = factory.GetLogger("HealthCheckManager")

type Node interface {
	UpdateHealthStatus(bool)
	IsActive() bool
}

type HealthCheckManager struct {
	node                Node
	cpuHealthChecker    *CpuHealthChecker
	memoryHealthChecker *MemoryHealthChecker
	lastUpdateTimestamp int64
	stopChan            chan struct{}
}

type ShutDown interface {
	AddClosingContext(ctxClose context.CancelFunc)
	AddShutdownChannel(shutdown chan struct{})
}

func NewHealthCheckManager(
	shutdown ShutDown,
	node Node,
	propertiesManager *property.GurmsPropertiesManager) *HealthCheckManager {

	properties := propertiesManager.LocalGurmsProperties.HealthCheck
	healthCheckManager := &HealthCheckManager{
		node:                node,
		stopChan:            make(chan struct{}),
		cpuHealthChecker:    NewCpuHealthChecker(properties.Cpu),
		memoryHealthChecker: NewMemoryHealthChecker(properties.Memory),
	}
	shutdown.AddShutdownChannel(healthCheckManager.stopChan)
	healthCheckManager.startHealthCheck(properties.CheckIntervalSeconds)
	return healthCheckManager
}

func (h *HealthCheckManager) IsHealthy() bool {
	return h.cpuHealthChecker.isCpuHealthy && h.memoryHealthChecker.isMemoryHealthy
}

func (h *HealthCheckManager) startHealthCheck(intervalSeconds int) {
	intervalMillis := int64(intervalSeconds) * 1000
	interval := time.Duration(intervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-h.stopChan:
				ticker.Stop()
				return
			case <-ticker.C:
				h.cpuHealthChecker.UpdateHealthStatus()
				h.memoryHealthChecker.UpdateHealthStatus()
				h.node.UpdateHealthStatus(h.IsHealthy())
				now := time.Now()
				previousUpdateTimeStamp := h.lastUpdateTimestamp + intervalMillis
				h.lastUpdateTimestamp = now.UnixMilli()
				if previousUpdateTimeStamp >= now.UnixMilli() {
					HEALTHCHECKMANAGERLOGGER.WarnWithArgs(fmt.Sprintf(
						"the system time goes backwards. the time drift is %d millis",
						previousUpdateTimeStamp-now.UnixMilli()))
				}
			}
		}
	}()
}

func (h *HealthCheckManager) StopHealthCheckRoutine() {
	h.stopChan <- struct{}{}
}
