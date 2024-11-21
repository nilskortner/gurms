package rpcservice

import (
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var RPCENDPOINTLOGGER logger.Logger = factory.GetLogger("RpcEndpoint")

const (
	EXPECTED_MAX_QPS           = 1000
	EXPECTED_AVERAGE_RTT       = 10
	INITAL_CAPACITY_PERCENTAGE = 10
)

var pendingRequestMap cmap.ConcurrentMap = cmap.New()[]

type RpcEndpoint struct {
	NodeId     string
	Connection *connectionservice.GurmsConnection
}

func NewRpcEndpoint(nodeId string, connection *connectionservice.GurmsConnection) *RpcEndpoint {
	return &RpcEndpoint{
		nodeId: nodeId,
	}
}

func (r *RpcEndpoint) HandleResponse(response *dto.RpcResponse) {
	resolveRequest(response.RequestId, response.Result, response.Rpcerror)
}

func resolveRequest(requestId int, response any, err error) {
	pendingRequestMap.remove(requestId)
}
