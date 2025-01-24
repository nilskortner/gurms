package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/connectionservice/request"
	"gurms/internal/infra/cluster/service/injection"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/infra/netutil"
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/cluster/connection"
	"gurms/internal/supportpkgs/mathsupport"
	"net"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var CONNECTIONLOGGER logger.Logger = factory.GetLogger("ConnectionService")

type ConnectionService struct {
	node injection.Node

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

	mu              sync.Mutex
	cancelKeepAlive context.CancelFunc
}

func NewConnectionService(connectionProperties *connection.ConnectionProperties, node injection.Node) *ConnectionService {
	clientProperties := connectionProperties.Client

	service := &ConnectionService{
		node:                              node,
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

func (c *ConnectionService) LazyInit(discoveryService *DiscoveryService, rpcService *RpcService) {
	c.discoveryService = discoveryService
	c.rpcService = rpcService
}

func (c *ConnectionService) IsMemberConnected(memberId string) bool {
	connection, ok := c.nodeIdToConnection.Get(memberId)
	if !ok || connection == nil {
		return false
	}
	if !connection.IsClosed.Load() && !connection.IsClosing {
		return true
	}
	return false
}

func (c *ConnectionService) updateHasConnectedToAllMembers(allMemberNodeIds cmap.ConcurrentMap[string, *configdiscovery.Member]) {
	connectedToAllMembers := true
	c.mu.Lock()
	defer c.mu.Unlock()
	for nodeId := range allMemberNodeIds.Items() {
		if !c.IsMemberConnected(nodeId) && !(c.discoveryService.LocalNodeStatusManager.LocalMember.Key.NodeId == nodeId) {
			connectedToAllMembers = false
			break
		}
	}
	c.hasConnectedToAllMembers = connectedToAllMembers
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

func (c *ConnectionService) initTcpConnection(host string, port int) (net.Conn, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	var conn net.Conn
	var err error

	tlsConfig := netutil.CreateTlsConfig(c.serverProperties.Tls, false)

	if c.serverProperties.Tls.Enabled {
		conn, err = tls.Dial("tcp", address, tlsConfig)
	} else {
		conn, err = net.Dial("tcp", address)
	}
	return conn, err
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
	conn, err := c.initTcpConnection(member.MemberHost, member.MemberPort)
	if err != nil {
		c.connectMemberUntilSucceedOrRemoved0Retry(nodeId, member.MemberHost, member.MemberPort, member, err)
		return
	}
	connection := connectionservice.NewGurmsConnection(
		nodeId, conn, true, c.newMemberConnectionListeners())
	c.OnMemberConnectionAdded(member, connection)
	localNodeId := c.discoveryService.LocalNodeStatusManager.LocalMember.Key.NodeId
	message = fmt.Sprint(
		"[Client] Sending an opening handshake request to the member: {id=%s}, host={%s}, port={%d}",
		nodeId, member.MemberHost, member.MemberPort)
	CONNECTIONLOGGER.InfoWithArgs(message)
	handshakeRequest := request.NewOpeningHandshakeRequest[byte](localNodeId, c.node)
	code, err := RequestResponseWithGurmsConnection[byte](c.rpcService, nodeId, handshakeRequest, -1, connection)
	if err != nil {
		connectMemberUntilSucceedOrRemoved0LogError(nodeId, member.MemberHost, member.MemberPort, connection, err)
		return
	}
	if code == request.RESPONSE_CODE_SUCCESS {
		c.OnMemberConnectionHandshakeCompleted(member, connection, true)
	} else {
		err = fmt.Errorf("failure code: " + string(code))
		c.connectMemberUntilSucceedOrRemoved0Retry(nodeId, member.MemberHost, member.MemberPort, member, err)
	}
}

func (c *ConnectionService) connectMemberUntilSucceedOrRemoved0Retry(
	nodeId string, host string, port int, member *configdiscovery.Member, err error) {
	if !c.discoveryService.IsKnownMember(nodeId) {
		return
	}
	retryTimes, _ := c.nodeIdToConnectionRetries.Get(nodeId)
	message := fmt.Sprint("[Client] Failed to connect to member: {id=%s, host=%s, port=%d}. Retry times: %d",
		nodeId, host, port, retryTimes)
	CONNECTIONLOGGER.ErrorWithMessage(message, err)
	retryTimes++
	c.nodeIdToConnectionRetries.Set(nodeId, retryTimes)
	time.Sleep(time.Duration(mathsupport.MinInt64(int64(retryTimes)*10, 60)) * time.Second)
	if !c.IsMemberConnected(nodeId) && c.discoveryService.IsKnownMember(nodeId) {
		c.connectMemberUntilSucceedOrRemoved0(member)
	} else {
		c.nodeIdToConnectionRetries.Remove(nodeId)
	}
}

func connectMemberUntilSucceedOrRemoved0LogError(nodeId string,
	host string, port int, connection *connectionservice.GurmsConnection, err error) {
	message := fmt.Sprint("[Client] Failed to complete the opening handshake with the the member:"+
		"{id=%s}, host={%s}, port={%s}. Closing Connection to reconnect",
		nodeId, host, port)
	CONNECTIONLOGGER.ErrorWithMessage(message, err)
	disconnectConnection(connection)
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

func (c *ConnectionService) Keepalive(nodeId string) {
	// TODO
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
	keepAliveRequest := request.NewKeepAliveRequest[int](connection.NodeId)
	_, err := RequestResponseWithId(c.rpcService, id, keepAliveRequest)
	if err != nil {
		CONNECTIONLOGGER.WarnWithArgs("failed to send a keepalive request to the member: "+id, err)
	} else {
		connection.LastKeepaliveTimestamp = time.Now().UnixMilli()
	}
}

// handshake

func (c *ConnectionService) HandleHandshakeRequest(connection *connectionservice.GurmsConnection, nodeId string) byte {
	//TODO
	return 0
}

func disconnectConnection(connection *connectionservice.GurmsConnection) {
	connection.IsClosing = true
	err := connection.Connection.Close()
	if err != nil {
		CONNECTIONLOGGER.ErrorWithMessage("error closing connection "+connection.NodeId, err)
	}
	// TODO retry here because dispose() is missing?
	connection.StopDecoderChan <- struct{}{}
	connection.StopListenerChan <- struct{}{}
	connection.IsClosed.Store(true)
}

// lifecycle listeners

func (c *ConnectionService) addMemberConnectionListenerSupplier(supplier func() connectionservice.MemberConnectionListener) {
	c.memberConnectionListenerSuppliers = append(c.memberConnectionListenerSuppliers, supplier)
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
			case <-connection.StopListenerChan:
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
		c.discoveryService.GetMember(nodeId)
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
	isKnownMember := nodeId != "" && c.discoveryService.IsKnownMember(nodeId)
	isClosing := c.discoveryService.LocalNodeStatusManager.IsClosing
	if isLocalNodeClient && isKnownMember && !isClosing {
		message := fmt.Sprint("[%s] Try to reconnect the member%s after %s millis",
			nodeType, memberIdAndAddress, c.reconnectInterval*1000)
		CONNECTIONLOGGER.InfoWithArgs(message)
		time.Sleep(time.Duration(c.reconnectInterval) * time.Second)
		memberToConnect, _ := c.discoveryService.AllKnownMembers.Get(nodeId)
		if memberToConnect == nil {
			message := fmt.Sprint("[%s] Stop to reconnect the member%s because it has been unregistered",
				nodeType, memberIdAndAddress)
			CONNECTIONLOGGER.InfoWithArgs(message)
		} else {
			c.connectMemberUntilSucceedOrRemoved(memberToConnect)
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

func (c *ConnectionService) OnMemberConnectionHandshakeCompleted(
	member *configdiscovery.Member,
	connection *connectionservice.GurmsConnection,
	isLocalNodeCLient bool,
) {
	nodeId := member.Key.NodeId
	var local string
	if isLocalNodeCLient {
		local = "Client"
	} else {
		local = "Server"
	}
	message := fmt.Sprint("[%s] Completed the opening handshake with the member :%s[%s:%s]",
		local, nodeId, member.MemberHost, member.MemberPort)
	CONNECTIONLOGGER.InfoWithArgs(message)
	c.nodeIdToConnection.Set(nodeId, connection)
	c.nodeIdToConnectionRetries.Remove(nodeId)
	c.connectingMembers.Remove(nodeId)
	c.updateHasConnectedToAllMembers(c.discoveryService.AllKnownMembers)
	for _, listener := range connection.Listeners {
		err := listener.OnOpeningHandshakeCompleted(member)
		if err != nil {
			message = fmt.Sprint(
				"Caught an error while notifiying the onOpeningHandshakeCompleted listener: %s",
				listener.GetName())
			CONNECTIONLOGGER.ErrorWithMessage(message, err)
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
