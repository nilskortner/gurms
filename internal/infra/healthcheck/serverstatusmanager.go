package healthcheck

import "context"

type ServerStatusManager struct {
	ctx                context.Context
	node               Node
	healthCheckManager *HealthCheckManager
}

func NewServerStatusManager(
	shutdown ShutDown,
	node Node,
	healthCheckManager *HealthCheckManager) *ServerStatusManager {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown.AddShutdownFunction(cancel)
	return &ServerStatusManager{
		ctx:                ctx,
		node:               node,
		healthCheckManager: healthCheckManager,
	}
}

func (s *ServerStatusManager) getServiceAvailability() ServiceAvailability {
	select {
	case <-s.ctx.Done():
		return SHUTTING_DOWN
	default:
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
}
