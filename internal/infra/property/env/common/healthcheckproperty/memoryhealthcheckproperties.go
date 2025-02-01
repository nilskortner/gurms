package healthcheckproperty

import "gurms/internal/infra/unit"

type MemoryHealthCheckProperties struct {
	MaxAvailableMemoryPercentage           int
	MaxAvailableDirectMemoryPercentage     int
	MinFreeSystemMemoryBytes               int
	DirectMemoryWarningThresholdPercentage int
	HeapMemoryWarningThresholdPercentage   int
	MinMemoryWarningIntervalSeconds        int
	HeapMemoryGcThresholdPercentage        int
	MinHeapMemoryGcIntervalSeconds         int
}

func NewMemoryHealthCheckProperties() *MemoryHealthCheckProperties {
	return &MemoryHealthCheckProperties{
		MaxAvailableMemoryPercentage:           95,
		MaxAvailableDirectMemoryPercentage:     95,
		MinFreeSystemMemoryBytes:               128 * unit.MB,
		DirectMemoryWarningThresholdPercentage: 50,
		HeapMemoryWarningThresholdPercentage:   95,
		MinMemoryWarningIntervalSeconds:        10,
		HeapMemoryGcThresholdPercentage:        60,
		MinHeapMemoryGcIntervalSeconds:         10,
	}
}
