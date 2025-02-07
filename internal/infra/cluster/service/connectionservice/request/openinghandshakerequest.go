package request

import (
	"gurms/internal/infra/cluster/service/injection"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

const (
	OPENINGHANDSHAKEREQUESTNAME string = "openingHandshake"

	RESPONSE_CODE_SUCCESS                   byte = 1
	RESPONSE_CODE_CONNECTION_CLOSED         byte = 2
	RESPONSE_CODE_CONNECTION_ALREADY_EXISTS byte = 3
	RESPONSE_CODE_UNKNOWN_MEMBER            byte = 4
)

type OpeningHandshakeRequest[T comparable] struct {
	*dto.RpcBaseRequest
	node   injection.Node
	nodeId string
}

// only initialize as boolean
func NewOpeningHandshakeRequest[T comparable](nodeId string, node injection.Node) *OpeningHandshakeRequest[T] {
	return &OpeningHandshakeRequest[T]{
		nodeId: nodeId,
		node:   node,
	}
}

// for rpcrequest interface

func (o *OpeningHandshakeRequest[T]) IsAsync() bool {
	return false
}
func (o *OpeningHandshakeRequest[T]) CallAsync() any {
	return nil
}
func (o *OpeningHandshakeRequest[T]) Call() any {
	return o.node.OpeningHandshakeRequestCall(o.Connection)
}
func (o *OpeningHandshakeRequest[T]) NodeTypeToRequest() dto.NodeTypeToHandleRpc {
	return dto.BOTH
}
func (o *OpeningHandshakeRequest[T]) NodeTypeToRespond() dto.NodeTypeToHandleRpc {
	return dto.BOTH
}

func (o *OpeningHandshakeRequest[T]) Name() string {
	return OPENINGHANDSHAKEREQUESTNAME
}
