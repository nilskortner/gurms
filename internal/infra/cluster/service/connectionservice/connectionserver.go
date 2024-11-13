package connectionservice

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common"

	"google.golang.org/grpc"
)

var CONNECTIONSERVERLOGGER logger.Logger = factory.GetLogger("ConnectionServer")

var server *ConnectionServer

type ConnectionServer struct {
	host              string
	proposedPort      int
	portCount         int
	portAutoIncrement bool
	ssl               *common.SslProperties
	connectionStream  func(conn *grpc.ClientConn)
	port              int
}

func NewConnectionServer(
	host string,
	port int,
	portAutoIncrement bool,
	portCount int,
	ssl *common.SslProperties,
) *ConnectionServer {
	return &ConnectionServer{
		port:              -1,
		host:              host,
		proposedPort:      port,
		portAutoIncrement: portAutoIncrement,
		portCount:         portCount,
		ssl:               ssl,
	}
}

func (c *ConnectionServer) BlockUntilConnect() {
	if server != nil {
		return
	}
}

func (c *ConnectionServer) Shutdown() {
}

func (c *ConnectionServer) SetStream(stream func(conn *grpc.ClientConn)) {
	c.connectionStream = stream
}
