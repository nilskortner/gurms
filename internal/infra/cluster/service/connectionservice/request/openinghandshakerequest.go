package request

import "gurms/internal/infra/cluster/service/rpcservice/dto"

const (
	OPENINGHANDSHAKEREQUESTNAME string = "openingHandshake"

	RESPONSE_CODE_SUCCESS                   byte = 1
	RESPONSE_CODE_CONNECTION_CLOSED         byte = 2
	RESPONSE_CODE_CONNECTION_ALREADY_EXISTS byte = 3
	RESPONSE_CODE_UNKNOWN_MEMBER            byte = 4
)

type OpeningHandshakeRequest[T comparable] struct {
	*dto.RpcBaseRequest
	node   Node
	nodeId string
}

func NewOpeningHandshakeRequest[T comparable](nodeId string) *OpeningHandshakeRequest[T] {
	return &OpeningHandshakeRequest[T]{
		nodeId: nodeId,
	}
}

// for rpcrequest interface

func (o *OpeningHandshakeRequest[T]) IsAsync() bool {
	return false
}
func (o *OpeningHandshakeRequest[T]) CallAsync() T {
	var zero T
	return zero
}
func (o *OpeningHandshakeRequest[T]) Call() T {
	// TODO
	return zero
}
func (o *OpeningHandshakeRequest[T]) NodeTypeToRequest() dto.NodeTypeToHandleRpc {
	return dto.BOTH
}
func (o *OpeningHandshakeRequest[T]) NodeTypeToRespond() dto.NodeTypeToHandleRpc {
	return dto.BOTH
}

func (o *OpeningHandshakeRequest[T]) Name() string {
	return NAME
}
