package healthcheck

import (
	"fmt"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/healthcheckproperty"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

var CPUHEALTHMANAGERLOGGER logger.Logger = factory.GetLogger("CpuHealthManager")

type CpuHealthChecker struct {
	isCpuHealthCheckAvailable bool
	cpuCheckRetries           int
	unhealthyLoadThreshold    float64
	isCpuHealthy              bool
	currentUnhealthyTimes     int
	unhealthyReason           string
}

func NewCpuHealthChecker(properties *healthcheckproperty.CpuHealthCheckProperties) *CpuHealthChecker {
	manager := &CpuHealthChecker{
		cpuCheckRetries:        properties.Retries,
		unhealthyLoadThreshold: float64(properties.UnhealthyLoadThresholdPercentage),
	}
	cpuLoad := getCpuLoad()
	manager.isCpuHealthCheckAvailable = true
	if cpuLoad < 0 {
		str := fmt.Sprintf("the cpu health checker can not work because the recent cpu usage" +
			" for the whole operating environment is unavailable")
		CPUHEALTHMANAGERLOGGER.Warn(str)
		manager.isCpuHealthCheckAvailable = false
		manager.isCpuHealthy = true

	}
	manager.updateHealthStatus()
	return manager
}

func (c *CpuHealthChecker) updateHealthStatus() {
	if !c.isCpuHealthCheckAvailable {
		return
	}
	wasCpuHealthy := c.isCpuHealthy
	var cpuLoad float64 = getCpuLoad()
	if cpuLoad > c.unhealthyLoadThreshold {
		c.currentUnhealthyTimes++
		if c.currentUnhealthyTimes > c.cpuCheckRetries {
			c.unhealthyReason = fmt.Sprintf("the CPU usage is too high,"+
				" and it has failed the CPU health check %d times, which exceeds max retries %d",
				c.currentUnhealthyTimes, c.cpuCheckRetries)
			c.isCpuHealthy = false
		}
	} else {
		c.unhealthyReason = ""
		c.currentUnhealthyTimes = 0
		c.isCpuHealthy = true
	}

	// log
	str := fmt.Sprintf("CPU load is: %d", cpuLoad)
	CPUHEALTHMANAGERLOGGER.Debug(str)
	if wasCpuHealthy != c.isCpuHealthy {
		if c.isCpuHealthy {
			str = fmt.Sprintf("the CPU has become healthy. the current CPU load is: %d", cpuLoad)
			CPUHEALTHMANAGERLOGGER.InfoWithArgs(str)
		} else {
			str = fmt.Sprintf("the CPU has become unhealthy. the current CPU load is: %d", cpuLoad)
			CPUHEALTHMANAGERLOGGER.InfoWithArgs(str)
		}
	}
}

// cpu load is in 0.00 to 100.00 format
func getCpuLoad() float64 {
	load, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		CPUHEALTHMANAGERLOGGER.ErrorWithMessage("error getting cpu load", err)
		return -1
	}
	return load[0]
}
