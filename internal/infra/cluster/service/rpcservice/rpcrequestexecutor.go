package rpcservice

import (
	"bytes"
	"encoding/gob"
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
}

// TODO: add other request types
func UnwrapRunRpcRequest(wrapRequest dto.RpcRequestWrap, connection *connectionservice.GurmsConnection, fromNodeId string) *bytes.Buffer {
	switch value := wrapRequest.(type) {
	case dto.RpcRequest[byte]:
		return EncodeRequest(RunRpcRequest[byte](value, connection, fromNodeId))
	default:
		RPCENDPOINTLOGGER.ErrorWithArgs("Couldnt resolve Instantiation of Request[Type?]: ", wrapRequest)
		return nil
	}
}

func EncodeRequest[T comparable](result T) *bytes.Buffer {
	buffer := bytes.NewBuffer(make([]byte, 0))
	encoder := gob.NewEncoder(buffer)
	encoder.Encode(result)

	return buffer
}
