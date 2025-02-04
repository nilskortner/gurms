package healthcheck

type HealthChecker interface {
	IsHealthy() bool
	GetUnhealthyReason() string
	UpdateHealthStatus()
}
