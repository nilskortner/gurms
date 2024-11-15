package rpcservice

import (
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var RPCENDPOINTLOGGER logger.Logger = factory.GetLogger("RpcEndpoint")

const (
	EXPECTED_MAX_QPS           = 1000
	EXPECTED_ABERAGE_RTT       = 10
	INITAL_CAPACITY_PERCENTAGE = 10
)

type RpcEndpoint struct {
	NodeId     string
	Connection *connectionservice.GurmsConnection
	pendingRequestMap
}

func NewRpcEndpoint(nodeId string, connection *connectionservice.GurmsConnection) *RpcEndpoint {
	return &RpcEndpoint{
		nodeId: nodeId,
	}
}

func (r *RpcEndpoint) HandleResponse(response *dto.RpcResponse) {

}
