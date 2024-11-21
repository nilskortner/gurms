package connection

import (
	"gurms/internal/infra/property/env/common"
)

type ConnectionServerProperties struct {
	Host              string
	Port              int
	PortAutoIncrement bool
	PortCount         int
	Tls               *common.TlsProperties
}

func NewConnectionServerProperties() *ConnectionServerProperties {
	return &ConnectionServerProperties{
		Host:              "0.0.0.0",
		Port:              7510,
		PortAutoIncrement: false,
		PortCount:         100,
		Tls:               &common.TlsProperties{},
	}
}
