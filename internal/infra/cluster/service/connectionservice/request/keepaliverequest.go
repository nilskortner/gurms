package request

import (
	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

const NAME = "keepalive"

type KeepaliveRequest[T comparable] struct {
	*dto.RpcRequest[T]
	Node node.Node
}

func NewKeepAliveRequest[T comparable]() *KeepaliveRequest[T] {
	return &KeepaliveRequest[T]{}
}
