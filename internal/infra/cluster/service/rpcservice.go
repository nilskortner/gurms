package service

import (
	"fmt"
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/rpcservice"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/cluster"

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
	defaultRequestTimeout int
	codecService          *CodecService
	connectionService     *ConnectionService
	discoveryService      *DiscoveryService
	nodeIdToEndpoint      cmap.ConcurrentMap[string, *rpcservice.RpcEndpoint]
}

func NewRpcService(nodeType nodetype.NodeType, rpcProperties *cluster.RpcProperties) *RpcService {
	return &RpcService{
		nodeType:              nodeType,
		defaultRequestTimeout: rpcProperties.JobTimeoutMillis,
		nodeIdToEndpoint:      cmap.New[*rpcservice.RpcEndpoint](),
	}
}

// func LazyInit() {

// }

func RequestResponse(request *dto.RpcRequest) {

}

func RequestResponseWithId[T comparable](r *RpcService, memberNodeId string, request dto.RpcRequest[T]) chan T {
	return RequestResponseWithGurmsConnection(r, memberNodeId, request, -1, nil)
}

func RequestResponseWithDuration[T comparable](r *RpcService, memberNodeId string, request dto.RpcRequest[T], timeout int64) chan T {
}

func RequestResponseWithGurmsConnection[T comparable](r *RpcService, memberNodeId string, request dto.RpcRequest[T], timeout int64, connection *connectionservice.GurmsConnection) (chan T, error) {
	if r.discoveryService.localMember.Key.NodeId == memberNodeId {
		return rpcservice.RunRpcRequest()
	}
	endpoint, err := r.getOrCreateEndpointWithConnection(memberNodeId, connection)
	if err != nil {
		request.Release()
		return nil, err
	}
	return requestResponse0(endpoint, request, timeout), nil
}

func RequestResponseWithRpcEndpoint()

// internal implentations
func requestResponse0[T comparable]() chan T {

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
