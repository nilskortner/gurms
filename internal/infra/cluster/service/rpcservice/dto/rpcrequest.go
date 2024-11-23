package dto

import (
	"bytes"
	"gurms/internal/infra/cluster/service/connectionservice"
)

type RpcFunctions[T comparable] interface {
	IsAsync() bool
	CallAsync() T
	Call() T
	NodeTypeToRequest() NodeTypeToHandleRpc
	Name() string
}

type RpcRequest[T comparable] struct {
	RpcFunctions[T]
	Connection  *connectionservice.GurmsConnection
	FromNodeId  string
	RequestId   int
	RequestTime int64
	BoundBuffer *bytes.Buffer
}

func (r *RpcRequest[T]) Init(connection *connectionservice.GurmsConnection, fromNodeId string) {
	r.Connection = connection
	r.FromNodeId = fromNodeId
}

func (r *RpcRequest[T]) Release() {
	r.BoundBuffer = nil
}
