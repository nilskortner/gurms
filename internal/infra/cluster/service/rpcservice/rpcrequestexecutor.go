package rpcservice

import (
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

func RunRpcRequest[T comparable](
	rpcRequest *dto.RpcRequest[T],
	connection *connectionservice.GurmsConnection,
	fromNodeId string,
) (T, error) {
	rpcRequest.Init(connection, fromNodeId)
	var result T

	if rpcRequest.IsAsync() {
		result = rpcRequest.CallAsync()
	} else {
		result = rpcRequest.Call()

	}
	return result, nil
}
