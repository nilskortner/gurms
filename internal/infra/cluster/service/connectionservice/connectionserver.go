package connectionservice

import (
	"context"
	"crypto/tls"
	"fmt"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/netutil"
	"gurms/internal/infra/property/env/common"
	"net"
)

var CONNECTIONSERVERLOGGER logger.Logger = factory.GetLogger("ConnectionServer")

var server *ConnectionServer

type ConnectionServer struct {
	host               string
	proposedPort       int
	portCount          int
	portAutoIncrement  bool
	tlsConfig          *common.TlsProperties
	listener           net.Listener
	connectionConsumer func(conn net.Conn)
	port               int
	ctx                context.Context
	cancel             context.CancelFunc
}

func NewConnectionServer(
	host string,
	port int,
	portAutoIncrement bool,
	portCount int,
	tlsConfig *common.TlsProperties,
	connection func(conn net.Conn),
) *ConnectionServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &ConnectionServer{
		port:               -1,
		host:               host,
		proposedPort:       port,
		portAutoIncrement:  portAutoIncrement,
		portCount:          portCount,
		tlsConfig:          tlsConfig,
		connectionConsumer: connection,
		ctx:                ctx,
		cancel:             cancel,
	}
}

func (c *ConnectionServer) BlockUntilConnect() {
	if server != nil {
		return
	}
	currentPort := c.proposedPort

	for {
		address := fmt.Sprintf("%s:%d", c.host, currentPort)
		var listener net.Listener
		var err error

		tlsConfig := netutil.CreateTlsConfig(c.tlsConfig, true)

		if c.tlsConfig.Enabled {
			listener, err = tls.Listen("tcp", address, tlsConfig)
		} else {
			listener, err = net.Listen("tcp", address)
		}
		if err != nil {
			if c.portAutoIncrement && currentPort <= c.proposedPort+c.portCount {
				message := fmt.Sprint("failed to bind on the port: %d. trying to bind on the next port: %d", currentPort, currentPort+1)
				CONNECTIONSERVERLOGGER.Warn(message)
				currentPort++
				continue
			}
			panic("failed to set up the local discovery server")
		}
		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					select {
					case <-c.ctx.Done():
						CONNECTIONSERVERLOGGER.Warn("connectionserver has been stopped.")
						return
					default:
						CONNECTIONSERVERLOGGER.Error(err)
					}
					continue
				}
				c.connectionConsumer(conn)
			}
		}()

		c.listener = listener
		c.port = currentPort
		CONNECTIONSERVERLOGGER.InfoWithArgs("the local node server started on:" + address)
		return
	}
}

func (c *ConnectionServer) Shutdown() {
	c.cancel()
	c.listener.Close()
}
