package healthcheck

import (
	"fmt"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/infra/property/env/common/healthcheckproperty"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

var MEMORYHEALTHCHECKERLOGGER logger.Logger = factory.GetLogger("MemoryHealthChecker")

const MB uint64 = 1024 * 1024

type MemoryHealthChecker struct {
	isMemoryHealthy                      bool
	unhealthyReason                      string
	maxAvailableMemory                   uint64
	maxHeapMemory                        uint64
	minFreeSystemMemory                  int
	totalPhysicalMemorySize              uint64
	usedAvailableMemory                  uint64
	usedHeapMemory                       uint64
	usedSystemMemory                     uint64
	heapMemoryWarningThresholdPercentage int
	minMemoryWarningIntervalMillis       int64
	lastHeapMemoryWarningTimestamp       int64
}

func NewMemoryHealthChecker(
	properties *healthcheckproperty.MemoryHealthCheckProperties) *MemoryHealthChecker {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	vmStats := getVirtualMemory()
	maxHeapMemory := memStats.HeapSys
	totalPhysicalMemorySize := vmStats.Total
	maxAvailableMemory := totalPhysicalMemorySize * uint64(properties.MaxAvailableMemoryPercentage/100.0)
	minAvailableMemory := 1000 * MB
	if maxAvailableMemory < minAvailableMemory {
		str := fmt.Sprintf("the max available memory is too small to run"+
			"expected: >= %dMB. actual: %dMB", minAvailableMemory/MB, maxAvailableMemory/MB)
		panic(str)
	}

	if maxAvailableMemory < maxHeapMemory {
		str := fmt.Sprintf("the max available memory %dMB"+
			" should not be less than the max heap memory %dMB ",
			maxAvailableMemory/MB, maxHeapMemory/MB)
		panic(str)
	}
	if maxAvailableMemory > memStats.Sys {
		str := fmt.Sprintf("the max available memory %d is larger than the os memory %d"+
			" which indicates that some memory will never be used by the server",
			maxAvailableMemory/MB, memStats.Sys/MB)
		MEMORYHEALTHCHECKERLOGGER.Warn(str)
	}

	return &MemoryHealthChecker{
		maxHeapMemory:                        maxHeapMemory,
		maxAvailableMemory:                   maxAvailableMemory,
		heapMemoryWarningThresholdPercentage: properties.HeapMemoryWarningThresholdPercentage,
		minMemoryWarningIntervalMillis:       int64(properties.MinMemoryWarningIntervalSeconds * 1000),
		minFreeSystemMemory:                  properties.MinFreeSystemMemoryBytes,
		totalPhysicalMemorySize:              totalPhysicalMemorySize,
	}
}

func (m *MemoryHealthChecker) IsHealthy() bool {
	return m.isMemoryHealthy
}

func (m *MemoryHealthChecker) GetUnhealthyReason() string {
	return m.unhealthyReason
}

func (m *MemoryHealthChecker) UpdateHealthStatus() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	usedHeapMemory := memStats.HeapAlloc
	usedStackMemory := memStats.StackInuse
	usedOtherMemory := memStats.Sys - memStats.HeapSys - memStats.StackSys
	localUsedAvailableMemory := usedHeapMemory + usedStackMemory + usedOtherMemory

	vmStats := getVirtualMemory()
	localFreeSystemMemory := vmStats.Available
	localUsedSystemMemory := vmStats.Total - localFreeSystemMemory
	m.usedSystemMemory = localUsedSystemMemory

	var localIsMemoryHealthy bool
	if localUsedAvailableMemory < m.maxAvailableMemory &&
		localFreeSystemMemory > uint64(m.minFreeSystemMemory) {
		m.unhealthyReason = ""
		localIsMemoryHealthy = true
	} else {
		m.unhealthyReason = fmt.Sprintf("the memory is insufficient. "+
			"the insufficient memory usage snapshot is: "+
			"used system memory: %d/%d; used available memory: %d/%d; "+
			"", localUsedSystemMemory/MB, m.totalPhysicalMemorySize/MB,
			localUsedAvailableMemory/MB, m.maxAvailableMemory/MB)
		localIsMemoryHealthy = false
	}
	m.isMemoryHealthy = localIsMemoryHealthy
	m.tryLog(localIsMemoryHealthy)
}

func (m *MemoryHealthChecker) tryLog(isHealthy bool) {
	var loglvl loglevel.LogLevel
	if isHealthy {
		loglvl = loglevel.DEBUG
	} else {
		loglvl = loglevel.WARN
	}
	if MEMORYHEALTHCHECKERLOGGER.IsEnabled(loglvl) {
		MEMORYHEALTHCHECKERLOGGER.Log(loglvl,
			fmt.Sprintf("used system memory: %d/%d; "+
				"used available memory: %d/%d; "+
				"used heap memory: %d/%d;",
				m.usedSystemMemory, m.totalPhysicalMemorySize,
				m.usedAvailableMemory, m.maxAvailableMemory,
				m.usedHeapMemory, m.maxHeapMemory))
	}
	now := time.Now()
	usedMemoryPercentage := 100.00 * float64(m.usedHeapMemory/m.maxHeapMemory)
	if m.heapMemoryWarningThresholdPercentage > 0 &&
		m.heapMemoryWarningThresholdPercentage < int(usedMemoryPercentage) &&
		m.minMemoryWarningIntervalMillis < (now.UnixMilli()-m.lastHeapMemoryWarningTimestamp) {
		m.lastHeapMemoryWarningTimestamp = int64(now.UnixMilli())
		MEMORYHEALTHCHECKERLOGGER.Warn(fmt.Sprintf(
			"the used heap mempory has exceeded the warning threshold: %d/%d/%d/%d",
			m.usedHeapMemory/MB, m.maxHeapMemory/MB, usedMemoryPercentage,
			m.heapMemoryWarningThresholdPercentage))
	}
}

func getVirtualMemory() *mem.VirtualMemoryStat {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		MEMORYHEALTHCHECKERLOGGER.ErrorWithMessage("error getting virtual memory", err)
	}
	return vmStat
}
