package healthcheckproperty

type HealthCheckProperties struct {
	checkIntervalSeconds int
	cpu                  *CpuHealthCheckProperties
	memory               *MemoryHealthCheckProperties
}

func NewHealthCheckProperties() *HealthCheckProperties {
	return &HealthCheckProperties{
		checkIntervalSeconds: 3,
		cpu:                  NewCpuHealthCheckProperties(),
		memory:               NewMemoryHealthCheckProperties(),
	}
}
