package healthcheckproperty

type HealthCheckProperties struct {
	CheckIntervalSeconds int
	Cpu                  *CpuHealthCheckProperties
	Memory               *MemoryHealthCheckProperties
}

func NewHealthCheckProperties() *HealthCheckProperties {
	return &HealthCheckProperties{
		CheckIntervalSeconds: 3,
		Cpu:                  NewCpuHealthCheckProperties(),
		Memory:               NewMemoryHealthCheckProperties(),
	}
}
