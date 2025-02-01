package healthcheck

import (
	"fmt"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/healthcheckproperty"
	"runtime"

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
	usedAvailableMemory                  uint64
	heapMemoryWarningThresholdPercentage int
	minMemoryWarningIntervalSeconds      int
	lastHeapMemoryWarningTimestamp       uint64
}

func NewMemoryHealthChecker(
	properties *healthcheckproperty.MemoryHealthCheckProperties) *MemoryHealthChecker {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	vmStats := getVirtualMemory()
	maxHeapMemory := memStats.HeapSys
	maxAvailableMemory := vmStats.Total * uint64(properties.MaxAvailableMemoryPercentage/100.0)
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
			maxAvailableMemory/100, memStats.Sys/100)
		MEMORYHEALTHCHECKERLOGGER.Warn(str)
	}

	return &MemoryHealthChecker{
		maxHeapMemory:                        maxHeapMemory,
		maxAvailableMemory:                   maxAvailableMemory,
		heapMemoryWarningThresholdPercentage: properties.HeapMemoryWarningThresholdPercentage,
		minMemoryWarningIntervalSeconds:      properties.MinMemoryWarningIntervalSeconds * 1000,
		minFreeSystemMemory:                  properties.MinFreeSystemMemoryBytes,
	}
}

func (m *MemoryHealthChecker) updateHealthStatus() {
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
		localFreeSystemMemory > m.minFreeSystemMemory {
		m.unhealthyReason = ""
		localIsMemoryHealthy = true
	} else {
		m.unhealthyReason = fmt.Sprintf("the memory is insufficient. " +
			"the insufficient memory usage snapshot is:")
		localIsMemoryHealthy = false
	}
	m.isMemoryHealthy = localIsMemoryHealthy
	m.tryLog(localIsMemoryHealthy)
}

func (m *MemoryHealthChecker) tryLog(isHealthy bool) {

}

func getVirtualMemory() *mem.VirtualMemoryStat {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		MEMORYHEALTHCHECKERLOGGER.ErrorWithMessage("error getting virtual memory", err)
	}
	return vmStat
}
