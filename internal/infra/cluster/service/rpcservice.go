package service

import (
	"bytes"
	"fmt"
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/rpcservice"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/cluster"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var RPCLOGGER logger.Logger = factory.GetLogger("RpcService")

const (
	METRICS_NAME_GRPC_REQUEST          = "rpc.request"
	METRICS_TAG_REQUEST_NAME           = "name"
	METRICS_TAG_REQUEST_TARGET_NODE_ID = "node"
)

type RpcService struct {
	nodeType              nodetype.NodeType
	defaultRequestTimeout time.Duration
	codecService          *CodecService
	connectionService     *ConnectionService
	discoveryService      *DiscoveryService
	nodeIdToEndpoint      cmap.ConcurrentMap[string, *rpcservice.RpcEndpoint]
}

func NewRpcService(nodeType nodetype.NodeType, rpcProperties *cluster.RpcProperties) *RpcService {
	return &RpcService{
		nodeType:              nodeType,
		defaultRequestTimeout: time.Duration(int64(rpcProperties.JobTimeoutMillis)) * time.Millisecond,
		nodeIdToEndpoint:      cmap.New[*rpcservice.RpcEndpoint](),
	}
}

// func LazyInit() {

// }

func RequestResponse(request dto.RpcRequest) {

}

func RequestResponseWithId[T comparable](r *RpcService, memberNodeId string, request dto.RpcRequest[T]) (T, error) {
	return RequestResponseWithGurmsConnection(r, memberNodeId, request, -1, nil)
}

func RequestResponseWithDuration[T comparable](r *RpcService, memberNodeId string, request dto.RpcRequest[T], timeout int64) chan T {
}

func RequestResponseWithGurmsConnection[T comparable](r *RpcService, memberNodeId string, request dto.RpcRequest[T], timeout int64, connection *connectionservice.GurmsConnection) (T, error) {
	if r.discoveryService.localMember.Key.NodeId == memberNodeId {
		value, err := rpcservice.RunRpcRequest[T](request, nil, memberNodeId)
		if err != nil {
			var zero T
			return zero, err
		} else {
			return value, nil
		}
	}
	endpoint, err := r.getOrCreateEndpointWithConnection(memberNodeId, connection)
	if err != nil {
		request.Release()
		var zero T
		return zero, err
	}
	return requestResponse0(r, endpoint, request, timeout)
}

func RequestResponseWithRpcEndpoint()

// internal implentations
func requestResponse0[T comparable](r *RpcService, endpoint *rpcservice.RpcEndpoint, request dto.RpcRequest[T], timeout int64) (T, error) {
	err := assertCurrentNodeIsAllowedToSend(r, request)
	if err != nil {
		request.Release()
		var zero T
		return zero, err
	}
	if timeout == -1 {
		timeout = r.defaultRequestTimeout.Milliseconds()
	}
	var requestBody *bytes.Buffer
	requestBody, err = Serialize(request)
	if err != nil {
		var zero T
		return zero, err
	}
	return rpcservice.SendRequest(endpoint, request, requestBody)
}

func OnConnectionOpened() {

}

// region MemberConnectionListener

type RpcMemberConnectionListener struct {
	rpcService *RpcService
	connection *connectionservice.GurmsConnection
	endpoint   *rpcservice.RpcEndpoint
}

func (r *RpcMemberConnectionListener) OnConnectionOpened(connection *connectionservice.GurmsConnection) error

func (r *RpcMemberConnectionListener) OnConnectionClosed()

func (r *RpcMemberConnectionListener) OnOpeningHandshakeCompleted(member *configdiscovery.Member)

func (r *RpcMemberConnectionListener) OnClosingHandshakeCompleted()

func (r *RpcMemberConnectionListener) OnDataReceived(value any) {
	switch value := value.(type) {
	case dto.RpcRequest:
		r.onRequestReceived(value)
	case dto.RpcResponse:
		r.onResponseReceived(value)
	default:
		RPCLOGGER.ErrorWithArgs("Received unkown data: ", value)
	}
}

func (r *RpcMemberConnectionListener) onRequestReceived(request dto.RpcRequest) {

}

func (r *RpcMemberConnectionListener) onResponseReceived(response dto.RpcResponse) {
	if r.endpoint == nil {
		r.endpoint = r.getOrCreateEndpoint(r.connection.NodeId, r.connection)
	}
	r.endpoint.HandleResponse(response)
}

func (r *RpcService) getOrCreateEndpoint() (*rpcservice.RpcEndpoint, error) {

}

func (r *RpcService) getOrCreateEndpointWithConnection(nodeId string, connection *connectionservice.GurmsConnection) (*rpcservice.RpcEndpoint, error) {
	if nodeId == r.discoveryService.localMember.Key.NodeId {
		return nil, fmt.Errorf("The target node ID of RPC endpoint cannot be the local node ID: " + nodeId)
	}
	endpoint, success := r.nodeIdToEndpoint.Get(nodeId)
	if success == true && (connection == nil || connection == endpoint.Connection) {
		return endpoint, nil
	}
	var err error
	endpoint, err = r.createEndpoint(nodeId, connection)
	if err != nil {
		return nil, err
	}
	r.nodeIdToEndpoint.SetIfAbsent(nodeId, endpoint)
	return endpoint, nil
}

func (r *RpcService) createEndpoint(nodeId string, connection *connectionservice.GurmsConnection) (*rpcservice.RpcEndpoint, error) {
	if connection == nil {
		var ok bool
		connection, ok = r.connectionService.nodeIdToConnection.Get(nodeId)
		if !ok {
			return nil, fmt.Errorf("the connection to the member " + nodeId + " does not exist")
		}
	}
	return rpcservice.NewRpcEndpoint(nodeId, connection), nil
}

func assertCurrentNodeIsAllowedToSend[T comparable](r *RpcService, request dto.RpcRequest[T]) error {
	typ := request.NodeTypeToRequest()
	var allowed bool
	switch typ {
	case dto.BOTH:
		allowed = true
	case dto.GATEWAY:
		allowed = r.nodeType == nodetype.GATEWAY
	case dto.SERVICE:
		allowed = r.nodeType == nodetype.SERVICE
	}
	if !allowed {
		return fmt.Errorf("the node type of the current server is: %s, which cannot send the request \"%s\" that requires the node type: %s",
			r.nodeType.GetDisplayName(),
			request.Name(),
			typ)
	}
	return nil
}
