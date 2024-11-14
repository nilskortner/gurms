package service

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/discovery"
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
	discoveryService      *discovery.DiscoveryService
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

func RequestResponseWithId(memberNodeId string, request dto.RpcRequest) {}

func RequestResponseWithDuration(memberNodeId string, request dto.RpcRequest, timeout int64) {

}

func RequestResponseWithGurmsConnection(memberNodeId string, request dto.RpcRequest, timeout int64, connection *connectionservice.GurmsConnection) {
}

func RequestResponseWithRpcEndpoint()

// internal implentations
func requestResponse0() {}

func OnConnectionOpened() {

}

// region MemberConnectionListener

type RpcMemberConnectionListener struct {
}

func (r *RpcMemberConnectionListener) OnConnectionOpened(connection *connectionservice.GurmsConnection)
func (r *RpcMemberConnectionListener) OnConnectionClosed()
func (r *RpcMemberConnectionListener) OnOpeningHandshakeCompleted(member *configdiscovery.Member)
func (r *RpcMemberConnectionListener) OnClosingHandshakeCompleted()
func (r *RpcMemberConnectionListener) OnDataReceived()
