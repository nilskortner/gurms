package service

import (
	"context"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/connectionservice/request"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/cluster/connection"
	"net"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var CONNECTIONLOGGER logger.Logger = factory.GetLogger("ConnectionService")

type ConnectionService struct {
	clientTls               *common.TlsProperties
	keepaliveIntervalMillis int64
	keepaliveTimeoutMillis  int64
	reconnectInterval       int64

	nodeIdToConnection        cmap.ConcurrentMap[string, *connectionservice.GurmsConnection]
	nodeIdToConnectionRetries cmap.ConcurrentMap[string, int]
	connectingMembers         cmap.ConcurrentMap[string, struct{}]

	memberConnectionListenerSuppliers []func() *connectionservice.MemberConnectionListener
	discoveryService                  *DiscoveryService
	rpcService                        *RpcService
	hasConnectedToAllMembers          bool
	serverProperties                  *connection.ConnectionServerProperties
	server                            *connectionservice.ConnectionServer

	cancelKeepAlive context.CancelFunc
}

func NewConnectionService(connectionProperties *connection.ConnectionProperties) *ConnectionService {
	clientProperties := connectionProperties.Client

	service := &ConnectionService{
		memberConnectionListenerSuppliers: make([]func() *connectionservice.MemberConnectionListener, 0, 4),
		serverProperties:                  connectionProperties.Server,
		clientTls:                         clientProperties.Tls,
		keepaliveIntervalMillis:           int64(clientProperties.KeepAliveIntervalSeconds) * 1000,
		keepaliveTimeoutMillis:            int64(clientProperties.KeepAliveTimeoutSeconds) * 1000,
		reconnectInterval:                 int64(clientProperties.ReconnectIntervalSeconds),
	}

	var ctx context.Context
	ctx, service.cancelKeepAlive = context.WithCancel(context.Background())
	service.startSendKeepAliveToConnectionsForeverRoutine(ctx)

	service.server = service.setupServer()

	return service
}

func (c *ConnectionService) LazyInitConnectionService(discoveryService *DiscoveryService, rpcService *RpcService) {
	c.discoveryService = discoveryService
	c.rpcService = rpcService
}

func (c *ConnectionService) setupServer() *connectionservice.ConnectionServer {
	conn := func(conn net.Conn) {
		connection := connectionservice.NewGurmsConnection("", conn, false, c.newMemberConnectionListeners())
		c.OnMemberConnectionAdded(nil, connection)
	}

	server := connectionservice.NewConnectionServer(
		c.serverProperties.Host,
		c.serverProperties.Port,
		c.serverProperties.PortAutoIncrement,
		c.serverProperties.PortCount,
		c.serverProperties.Tls,
		conn,
	)
	server.BlockUntilConnect()
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
				for id, connection := range c.nodeIdToConnection.Items() {
					c.sendKeepAlive(id, connection)
				}
				time.Sleep(1000 * time.Millisecond)
			}
		}
	}()
}

func (c *ConnectionService) sendKeepAlive(id string, connection *connectionservice.GurmsConnection) {
	if !c.nodeIdToConnection.Has(id) {
		return
	}
	if !connection.IsLocalNodeClient {
		return
	}
	now := time.Now().UnixMilli()
	elapsedTime := now - connection.LastKeepaliveTimestamp
	if elapsedTime > c.keepaliveTimeoutMillis {
		CONNECTIONLOGGER.Warn("Reconnection to the member " + connection.NodeId + " due to keepalive timeout")
		disconnectConnection(connection)
		c.nodeIdToConnection.Remove(id)
		return
	}
	if elapsedTime < c.keepaliveIntervalMillis {
		return
	}
	// TODO: request
	keepAliveRequest := request.NewKeepAliveRequest()
	rrChan := RequestResponse(id, keepAliveRequest)
	subscribe(rrChan)
}

func disconnectConnection(connection *connectionservice.GurmsConnection) {
	connection.IsClosing = true
	err := connection.Connection.Close()
	if err != nil {
		CONNECTIONLOGGER.ErrorWithMessage("error closing connection "+connection.NodeId, err)
	}
}

func (c *ConnectionService) newMemberConnectionListeners() []*connectionservice.MemberConnectionListener {
	list := make([]*connectionservice.MemberConnectionListener, len(c.memberConnectionListenerSuppliers))
	for _, listener := range c.memberConnectionListenerSuppliers {
		list = append(list, listener())
	}
	return list
}

func (c *ConnectionService) OnMemberConnectionAdded(member *configdiscovery.Member, connection *connectionservice.GurmsConnection) {
	var endpointType string
	if connection.IsLocalNodeClient {
		endpointType = "Client"
	} else {
		endpointType = "Server"
	}
	memberIdAndAddress := getMemberIdAndAddress(connection.NodeId, member)
	CONNECTIONLOGGER.InfoWithArgs("[{}] Connected to the Member" + memberIdAndAddress)
	for _, listener := range connection.Listeners {
		err := listener.OnConnectionOpened(connection)
		if err != nil {
			CONNECTIONLOGGER.ErrorWithMessage("caught an error while notifiying the OnConnectionOpened listener: ", err)
		}
	}
	conn := connection.Connection
	for value := range conn.DataChan {
		for _, listener := range connection.Listeners {
			err := listener.OnDataReceived(value)
			if err != nil {
				CONNECTIONLOGGER.ErrorWithMessage("caught an error while notifiying the onDataReceived listener.", err)
			}
		}
	}
}

func getMemberIdAndAddress(nodeId string, member *configdiscovery.Member) string {
	if member == nil {
		return nodeId
	}
	return "{id=" +
		member.Key.NodeId +
		", host=" +
		member.MemberHost +
		", port=" +
		string(member.MemberPort) + "}"
}
