package rpcserv

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service/codec"
	"gurms/internal/infra/cluster/service/connection"
	"gurms/internal/infra/cluster/service/discovery"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/cluster"
)

var RPCLOGGER logger.Logger = factory.GetLogger("Rpc")

const (
	METRICS_NAME_GRPC_REQUEST          = "rpc.request"
	METRICS_TAG_REQUEST_NAME           = "name"
	METRICS_TAG_REQUEST_TARGET_NODE_ID = "node"
)

type RpcService struct {
	nodeType              nodetype.NodeType
	defaultRequestTimeout int
	codecService          *codec.CodecService
	connectionService     *connection.ConnectionService
	discoveryService      *discovery.DiscoveryService
}

func NewRpcService(rpcProperties *cluster.RpcProperties) *RpcService {
	return &RpcService{
		defaultRequestTimeout: rpcProperties.JobTimeoutMillis,
	}
}

// func LazyInit() {

// }
