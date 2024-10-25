package connection

import "gurms/internal/infra/property/env/common"

type ConnectionClientProperties struct {
	keepAliveTimeoutSeconds  int
	keepAliveIntervalSeconds int
	reconnectIntervalSeconds int
	ssl                      *common.SslProperties
}

func NewConnectionClientProperties() *ConnectionClientProperties {
	return &ConnectionClientProperties{
		keepAliveTimeoutSeconds:  15,
		keepAliveIntervalSeconds: 5,
		reconnectIntervalSeconds: 15,
		ssl:                      &common.SslProperties{},
	}
}
