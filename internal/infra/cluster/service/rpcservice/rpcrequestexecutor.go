package rpcservice

import (
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

func RunRpcRequest[T comparable](
	rpcRequest dto.RpcRequest[T],
	connection *connectionservice.GurmsConnection,
	fromNodeId string,
) T {
	rpcRequest.Init(connection, fromNodeId)
	var result T

	if rpcRequest.IsAsync() {
		result = rpcRequest.CallAsync()
	} else {
		result = rpcRequest.Call()
	}
	return result

	// TODO: error handling?
}

// TODO: add other request types
// func RunRpcRequest(wrapRequest dto.RpcRequestWrap, connection *connectionservice.GurmsConnection, fromNodeId string) (*bytes.Buffer, error) {
// 	switch value := wrapRequest.(type) {
// 	case dto.RpcRequest[byte]:
// 		return channel.EncodeRequest(RunRpcRequest[byte](value, connection, fromNodeId))
// 	default:
// 		RPCENDPOINTLOGGER.ErrorWithArgs("Couldnt resolve Instantiation of Request[Type?]: ", wrapRequest)
// 		return nil, error
// 	}
// }
