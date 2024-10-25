package healthcheck

import "gurms/internal/infra/unit"

type MemoryHealthCheckProperties struct {
	maxAvailableMemoryPercentage           int
	maxAvailableDirectMemoryPercentage     int
	minFreeSystemMemoryBytes               int
	directMemoryWarningThresholdPercentage int
	heapMemoryWarningThresholdPercentage   int
	minMemoryWarningIntervalSeconds        int
	heapMemoryGcThresholdPercentage        int
	minHeapMemoryGcIntervalSeconds         int
}

func NewMemoryHealthCheckProperties() *MemoryHealthCheckProperties {
	return &MemoryHealthCheckProperties{
		maxAvailableMemoryPercentage:           95,
		maxAvailableDirectMemoryPercentage:     95,
		minFreeSystemMemoryBytes:               128 * unit.MB,
		directMemoryWarningThresholdPercentage: 50,
		heapMemoryWarningThresholdPercentage:   95,
		minMemoryWarningIntervalSeconds:        10,
		heapMemoryGcThresholdPercentage:        60,
		minHeapMemoryGcIntervalSeconds:         10,
	}
}
