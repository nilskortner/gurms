package healthcheck

type ServerStatusManager struct {
	//context            *GurmsApplicationContext
	node               Node
	healthCheckManager *HealthCheckManager
}

// TODO:
func NewServerStatusManager(
	//context *GurmsApplicationContext,
	node Node,
	healthCheckManager *HealthCheckManager) *ServerStatusManager {
	return &ServerStatusManager{
		context:            context,
		node:               node,
		healthCheckManager: healthCheckManager,
	}
}

func (s *ServerStatusManager) getServiceAvailability() ServiceAvailability {
	if context.isClosing() {
		return SHUTTING_DOWN
	}
	if !s.node.IsActive() {
		return INACTIVE
	}
	unhealthyReason := s.healthCheckManager.cpuHealthChecker.GetUnhealthyReason()
	if unhealthyReason != "" {
		return ServiceAvailability{
			status: STATUS_HIGH_CPU_USAGE,
			reason: unhealthyReason,
		}
	}
	unhealthyReason = s.healthCheckManager.memoryHealthChecker.GetUnhealthyReason()
	if unhealthyReason != "" {
		return ServiceAvailability{
			status: STATUS_INSUFFICIENT_MEMORY,
			reason: unhealthyReason,
		}
	}
	return AVAILABLE
}
