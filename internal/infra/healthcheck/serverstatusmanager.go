package healthcheck

import "context"

type ServerStatusManager struct {
	ctx                context.Context
	node               Node
	healthCheckManager *HealthCheckManager
}

// TODO:
func NewServerStatusManager(
	ctx context.Context,
	node Node,
	healthCheckManager *HealthCheckManager) *ServerStatusManager {
	return &ServerStatusManager{
		ctx:                ctx,
		node:               node,
		healthCheckManager: healthCheckManager,
	}
}

func (s *ServerStatusManager) getServiceAvailability() ServiceAvailability {
	if s.ctx.Done() {
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
