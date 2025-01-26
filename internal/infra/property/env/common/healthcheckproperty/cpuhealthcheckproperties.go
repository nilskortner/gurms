package healthcheckproperty

type CpuHealthCheckProperties struct {
	UnhealthyLoadThresholdPercentage int
	Retries                          int
}

func NewCpuHealthCheckProperties() *CpuHealthCheckProperties {
	return &CpuHealthCheckProperties{
		UnhealthyLoadThresholdPercentage: 95,
		Retries:                          5,
	}
}
