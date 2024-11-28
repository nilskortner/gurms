package service

import (
	"context"
	"fmt"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/connectionservice/request"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/logging/core/model/loglevel"
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

	memberConnectionListenerSuppliers []func() connectionservice.MemberConnectionListener
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
		memberConnectionListenerSuppliers: make([]func() connectionservice.MemberConnectionListener, 0, 4),
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

func (c *ConnectionService) IsMemberConnected(memberId string) bool {
	connection, ok := c.nodeIdToConnection.Get(memberId)
	if !ok || connection == nil {
		return false
	}
	if !connection.IsClosed() && !connection.IsClosing {
		return true
	}
	return false
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

func (c *ConnectionService) connectMemberUntilSucceedOrRemoved(member *configdiscovery.Member) {
	nodeId := member.Key.NodeId
	_, ok := c.connectingMembers.Get(nodeId)
	if !member.IsSameNode(c.discoveryService.LocalNodeStatusManager.LocalMember) &&
		!c.IsMemberConnected(nodeId) &&
		!ok {
		c.connectingMembers.Set(nodeId, struct{}{})
		c.connectMemberUntilSucceedOrRemoved0(member)
	}
}

func (c *ConnectionService) connectMemberUntilSucceedOrRemoved0(member *configdiscovery.Member) {
	nodeId := member.Key.NodeId
	mappedId, _ := c.nodeIdToConnectionRetries.Get(nodeId)
	message := fmt.Sprint(
		"[Client] Connecting to the member: {id=%s}, host={%s}, port={%s}. Retry times: %s",
		nodeId, member.MemberHost, member.MemberPort, mappedId)
	CONNECTIONLOGGER.InfoWithArgs(message)
	conn := initTcpConnection(member.MemberHost, member.MemberPort)
	connection := connectionservice.NewGurmsConnection()
	c.OnMemberConnectionAdded(member, connection)
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
	keepAliveRequest := request.NewKeepAliveRequest[int]()
	_, err := RequestResponseWithId(c.rpcService, id, keepAliveRequest)
	if err != nil {
		CONNECTIONLOGGER.WarnWithArgs("failed to send a keepalive request to the member: "+id, err)
	} else {
		connection.LastKeepaliveTimestamp = time.Now().UnixMilli()
	}
}

func disconnectConnection(connection *connectionservice.GurmsConnection) {
	connection.IsClosing = true
	err := connection.Connection.Close()
	if err != nil {
		CONNECTIONLOGGER.ErrorWithMessage("error closing connection "+connection.NodeId, err)
	}
}

func (c *ConnectionService) newMemberConnectionListeners() []connectionservice.MemberConnectionListener {
	list := make([]connectionservice.MemberConnectionListener, len(c.memberConnectionListenerSuppliers))
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
	CONNECTIONLOGGER.InfoWithArgs("[{}] Connected to the Member"+memberIdAndAddress, endpointType)
	for _, listener := range connection.Listeners {
		err := listener.OnConnectionOpened(connection)
		if err != nil {
			CONNECTIONLOGGER.ErrorWithMessage("caught an error while notifiying the OnConnectionOpened listener: ", err)
		}
	}
	go func() {
		for {
			select {
			case <-connection.CloseContext.Done():
				c.onConnectionClosed(connection, nil)
				return
			case value, ok := <-connection.DataChan:
				if !ok && !connection.IsClosing {
					CONNECTIONLOGGER.Error(fmt.Errorf("failed to read from connection data channel"))
				} else {
					for _, listener := range connection.Listeners {
						err := listener.OnDataReceived(value)
						if err != nil {
							CONNECTIONLOGGER.ErrorWithMessage("caught an error while notifiying the onDataReceived listener.", err)
						}
					}
				}
			}
		}
	}()
}

func (c *ConnectionService) onConnectionClosed(connection *connectionservice.GurmsConnection, err error) {
	isLocalNodeClient := connection.IsLocalNodeClient
	var nodeType string
	if isLocalNodeClient {
		nodeType = "Client"
	} else {
		nodeType = "Server"
	}
	var member *configdiscovery.Member
	nodeId := connection.NodeId
	if nodeId == "" {
		member = nil
	} else {
		getMember(nodeId)
	}
	var memberIdAndAddress string
	if member == nil {
		memberIdAndAddress = "" + getMemberIdAndAddress(nodeId, member)
	} else {
		memberIdAndAddress = " " + getMemberIdAndAddress(nodeId, member)
	}
	closing := connection.IsClosing
	var level loglevel.LogLevel
	if closing {
		level = loglevel.INFO
	} else {
		level = loglevel.WARN
	}
	if member == nil {
		messageExpect := ""
		if closing {
			messageExpect = " unexpectedly"
		}
		message := fmt.Sprint("[%s] the connection to an unknown member has been closed%s",
			nodeType, messageExpect)
		CONNECTIONLOGGER.LogWithError(level, message, err)
	} else {
		c.connectingMembers.Remove(nodeId)
		messageExpect := ""
		if closing {
			messageExpect = " unexpectedly"
		}
		message := fmt.Sprint("[%s] the connection to the member %s has been closed%s",
			nodeType, memberIdAndAddress, messageExpect)
		CONNECTIONLOGGER.LogWithError(level, message, err)
	}
	for _, listener := range connection.Listeners {
		err := listener.OnConnectionClosed()
		if err != nil {
			CONNECTIONLOGGER.ErrorWithMessage(
				"caught an error while notifiyng the onConnectionClosed listener: "+listener.GetName(), err)
		}
	}
	isKnownMember := nodeId != "" && discoveryservice.IsKnownMember(nodeid)
	isClosing := c.discoveryService.localNodeStatusManager.isClosing
	if isLocalNodeClient && isKnownMember && !isClosing {
		message := fmt.Sprint("[%s] Try to reconnect the member%s after %s millis",
			nodeType, memberIdAndAddress, c.reconnectInterval*1000)
		CONNECTIONLOGGER.InfoWithArgs(message)
		time.Sleep(time.Duration(c.reconnectInterval) * time.Second)
		memberToConnect := c.discoveryService.GetAllKnownMembers().Get(nodeId)
		if memberToConnect == nil {
			message := fmt.Sprint("[%s] Stop to reconnect the member%s because it has been unregistered",
				nodeType, memberIdAndAddress)
			CONNECTIONLOGGER.InfoWithArgs(message)
		} else {
			connectMemberUntilSucceedOrRemoved(memberToConnect)
		}
	} else {
		var reason string
		if isLocalNodeClient {
			if isKnownMember {
				reason = "the local node is closing"
			} else {
				reason = "the member is unknown"
			}
		} else {
			reason = "the local node is server"
		}
		message := fmt.Sprint("[%s] Stop to connect the member%s because %s",
			nodeType, memberIdAndAddress, reason)
		CONNECTIONLOGGER.InfoWithArgs(message)
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
