package connectionservice

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common"
)

var CONNECTIONSERVERLOGGER logger.Logger = factory.GetLogger("ConnectionServer")

type ConnectionServer struct {
	host               string
	proposedPort       int
	portCount          int
	portAutoIncrement  bool
	ssl                *common.SslProperties
	connectionConsumer *any
	port               int
	server             *ConnectionServer
}

func NewConnectionServer(
	host string,
	port int,
	portAutoIncrement bool,
	portCount int,
	ssl *common.SslProperties,
	connectionConsumer *any,
) *ConnectionServer {
	return &ConnectionServer{
		port:               -1,
		host:               host,
		proposedPort:       port,
		portAutoIncrement:  portAutoIncrement,
		portCount:          portCount,
		ssl:                ssl,
		connectionConsumer: connectionConsumer,
	}
}

func (c *ConnectionServer) blockUntilConnect() {
	if c.server == nil {
		return
	}
}
