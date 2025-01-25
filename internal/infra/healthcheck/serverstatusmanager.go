package healthcheck

type ServerStatusManager struct {
	//context            *GurmsApplicationContext
	node               *Node
	healthCheckManager *HealthCheckManager
}

func NewServerStatusManager(
	//context *GurmsApplicationContext,
	node *Node,
	healthCheckManager *HealthCheckManager) *ServerStatusManager {
	return &ServerStatusManager{
		//context:            context,
		node:               node,
		healthCheckManager: healthCheckManager,
	}
}
