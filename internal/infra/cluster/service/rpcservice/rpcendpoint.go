package rpcservice

import (
	"bytes"
	"fmt"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var RPCENDPOINTLOGGER logger.Logger = factory.GetLogger("RpcEndpoint")

const (
	EXPECTED_MAX_QPS           = 1000
	EXPECTED_AVERAGE_RTT       = 10
	INITAL_CAPACITY_PERCENTAGE = 10
)

//var pendingRequestMap cmap.ConcurrentMap = cmap.New()[]

type RpcEndpoint struct {
	NodeId     string
	Connection *connectionservice.GurmsConnection
}

func NewRpcEndpoint(nodeId string, connection *connectionservice.GurmsConnection) *RpcEndpoint {
	return &RpcEndpoint{
		NodeId:     nodeId,
		Connection: connection,
	}
}

func SendRequest[T comparable](endpoint *RpcEndpoint, request *dto.RpcRequest[T], requestBody *bytes.Buffer) (T, error) {
	conn := endpoint.Connection.Connection
	if requestBody == nil {
		err := fmt.Errorf("the request body has been released")
		var zero T
		return zero, err
	}
	if conn.
}

func (r *RpcEndpoint) HandleResponse(response *dto.RpcResponse) {
	resolveRequest(response.RequestId, response.Result, response.Rpcerror)
}

func resolveRequest(requestId int, response any, err error) {
	pendingRequestMap.remove(requestId)
}
