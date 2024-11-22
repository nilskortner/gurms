package rpcservice

import (
	"context"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

func RunRpcRequest[T comparable](
	ctx context.Context,
	rpcRequest dto.RpcRequest[T],
	connection *connectionservice.GurmsConnection,
	fromNodeId string,
) chan T {

}
