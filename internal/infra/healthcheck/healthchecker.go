package healthcheck

type HealthChecker interface {
	IsHealthy() bool
	getUnhealthyReason() string
	updateHealthStatus()
}
