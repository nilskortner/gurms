package connrpcservice

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service/codec"
	"gurms/internal/infra/cluster/service/discovery"
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
	codecService          *codec.CodecService
	connectionService     *ConnectionService
	discoveryService      *discovery.DiscoveryService
	nodeIdToEndpoint      *cmap.ConcurrentMap[string, *RpcEndpoint]
}

func NewRpcService(nodeType nodetype.NodeType, rpcProperties *cluster.RpcProperties) *RpcService {
	return &RpcService{
		nodeType:              nodeType,
		defaultRequestTimeout: rpcProperties.JobTimeoutMillis,
		nodeIdToEndpoint:      cmap.New[string, *RpcEndpoint](),
	}
}

// func LazyInit() {

// }

func RequestResponse(request RpcRequest) {

}

func RequestResponseWithId(memberNodeId string, request RpcRequest) {}

func RequestResponseWithDuration(memberNodeId string, request RpcRequest, timeout int64) {

}

func RequestResponseWithGurmsConnection(memberNodeId string, request RpcRequest, timeout int64, connection *GurmsConnection) {
}

func RequestResponseWithRpcEndpoint()

// internal implentations
func requestResponse0() {}
