package healthcheck

type CpuHealthCheckProperties struct {
	unhealthyLoadThresholdPercentage int
	retries                          int
}

func NewCpuHealthCheckProperties() *CpuHealthCheckProperties {
	return &CpuHealthCheckProperties{
		unhealthyLoadThresholdPercentage: 95,
		retries:                          5,
	}
}
