package request

import (
	"gurms/internal/infra/cluster/service/injection"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

const KEEPALIVEREQUESTNAME = "keepalive"

type KeepaliveRequest[T comparable] struct {
	*dto.RpcBaseRequest
	node   injection.Node
	nodeId string
}

func NewKeepAliveRequest[T comparable](nodeId string) *KeepaliveRequest[T] {
	return &KeepaliveRequest[T]{}
}

// for rpcrequest interface

func (k *KeepaliveRequest[T]) IsAsync() bool {
	return false
}
func (k *KeepaliveRequest[T]) CallAsync() T {
	var zero T
	return zero
}
func (k *KeepaliveRequest[T]) Call() T {
	k.node.KeepAliveRequestCall()
	var zero T
	return zero
}
func (k *KeepaliveRequest[T]) NodeTypeToRequest() dto.NodeTypeToHandleRpc {
	return dto.BOTH
}
func (k *KeepaliveRequest[T]) NodeTypeToRespond() dto.NodeTypeToHandleRpc {
	return dto.BOTH
}

func (k *KeepaliveRequest[T]) Name() string {
	return KEEPALIVEREQUESTNAME
}
