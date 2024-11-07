package connrpcservice

import (
	"context"
	"gurms/internal/infra/cluster/service/connectionservice/request"
	"gurms/internal/infra/cluster/service/discovery"
	"gurms/internal/infra/cluster/service/rpcserv"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/cluster/connection"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"google.golang.org/grpc/connectivity"
)

var CONNECTIONLOGGER logger.Logger = factory.GetLogger("ConnectionService")

var nodeIdToConnection cmap.ConcurrentMap[string, *GurmsConnection] = cmap.New[*GurmsConnection]()

var nodeIdToConnectionRetries cmap.ConcurrentMap[string, int] = cmap.New[int]()
var connectingMembers cmap.ConcurrentMap[string, struct{}] = cmap.New[struct{}]()

type ConnectionService struct {
	clientSsl               *common.SslProperties
	keepaliveIntervalMillis int64
	keepaliveTimeoutMillis  int64
	reconnectInterval       int64

	nodeIdToConnection        cmap.ConcurrentMap[string, *GurmsConnection]
	nodeIdToConnectionRetries cmap.ConcurrentMap[string, int]
	connectingMembers         cmap.ConcurrentMap[string, struct{}]

	discoveryService         *discovery.DiscoveryService
	rpcService               *RpcService
	hasConnectedToAllMembers bool
	serverProperties         *connection.ConnectionServerProperties
	server                   *ConnectionServer
}

func NewConnectionService(connectionProperties *connection.ConnectionProperties) *ConnectionService {
	clientProperties := connectionProperties.Client

	service := &ConnectionService{
		serverProperties:        connectionProperties.Server,
		clientSsl:               clientProperties.Ssl,
		keepaliveIntervalMillis: int64(clientProperties.KeepAliveIntervalSeconds) * 1000,
		keepaliveTimeoutMillis:  int64(clientProperties.KeepAliveTimeoutSeconds) * 1000,
		reconnectInterval:       int64(clientProperties.ReconnectIntervalSeconds),
	}

	service.startSendKeepAliveToConnectionsForeverRoutine(context.Background())

	server := setupServer()

	return service
}

func (c *ConnectionService) LazyInitConnectionService(discoveryService *discovery.DiscoveryService, rpcService *rpcserv.RpcService) {
	c.discoveryService = discoveryService
	c.rpcService = rpcService
}

func (c *ConnectionService) setupServer() *ConnectionServer {
	server := NewConnectionServer()

	server.blockUntilConnect()
	return server
}

func (c *ConnectionService) startSendKeepAliveToConnectionsForeverRoutine(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				CONNECTIONLOGGER.Warn("SendKeepAliveToConnectionsForeverRoutine has been stopped")
				return
			default:
				for id, connection := range nodeIdToConnection.Items() {
					c.sendKeepAlive(id, connection)
				}
				time.Sleep(1000 * time.Millisecond)
			}
		}
	}()
}

func (c *ConnectionService) sendKeepAlive(id string, connection *GurmsConnection) {
	conn := connection.connection

	if conn.GetState() == connectivity.Shutdown {
		c.nodeIdToConnection.Remove(id)
		return
	}
	if !connection.isLocalNodeClient {
		return
	}
	now := time.Now().UnixMilli()
	elapsedTime := now - connection.lastKeepaliveTimestamp
	if elapsedTime > c.keepaliveIntervalMillis {
		CONNECTIONLOGGER.Warn("Reconnection to the member " + connection.nodeId + " due to keepalive timeout")
		disconnectConnection(connection)
		return
	}
	if elapsedTime < c.keepaliveIntervalMillis {
		return
	}
	rpcserv.RequestResponse(id, request.KeepaliveRequest{})
}

func disconnectConnection(connection *GurmsConnection) {
	connection.isClosing = true
	err := connection.connection.Close()
	if err != nil {
		CONNECTIONLOGGER.ErrorWithMessage("error closing connection "+connection.nodeId, err)
	}
}
