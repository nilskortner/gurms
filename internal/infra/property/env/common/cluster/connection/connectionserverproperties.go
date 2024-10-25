package connection

import "gurms/internal/infra/property/env/common"

type ConnectionServerProperties struct {
	host              string
	port              int
	portAutoIncrement bool
	portCount         int
	ssl               *common.SslProperties
}

func NewConnectionServerProperties() *ConnectionServerProperties {
	return &ConnectionServerProperties{
		host:              "0.0.0.0",
		port:              7510,
		portAutoIncrement: false,
		portCount:         100,
		ssl:               &common.SslProperties{},
	}
}
