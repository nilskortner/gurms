package connection

import "gurms/internal/infra/property/env/common"

type ConnectionClientProperties struct {
	KeepAliveTimeoutSeconds  int
	KeepAliveIntervalSeconds int
	ReconnectIntervalSeconds int
	Tls                      *common.TlsProperties
}

func NewConnectionClientProperties() *ConnectionClientProperties {
	return &ConnectionClientProperties{
		KeepAliveTimeoutSeconds:  15,
		KeepAliveIntervalSeconds: 5,
		ReconnectIntervalSeconds: 15,
		Tls:                      &common.TlsProperties{},
	}
}
