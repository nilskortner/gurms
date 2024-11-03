package connectionservice

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common"
)

var CONNECTIONSERVERLOGGER logger.Logger = factory.GetLogger("ConnectionServer")

type ConnectionServer struct {
	host              string
	proposedPort      int
	portAutoIncrement bool
	portCount         int
	ssl               *common.SslProperties
	connectionConsumer
	port   int
	server DisposableServer
}

func NewConnectionServer(
	host string,
	port int,
	portAutoIncrement bool,
	portCount int,
	ssl *common.SslProperties,
	connectionConsumer *Consumer,
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
