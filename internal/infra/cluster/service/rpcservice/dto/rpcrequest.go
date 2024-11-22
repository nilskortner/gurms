package dto

import (
	"bytes"
	"gurms/internal/infra/cluster/service/connectionservice"
)

type RpcRequest[T comparable] struct {
	connection  *connectionservice.GurmsConnection
	fromNodeId  string
	requestId   int
	requestTime int64
	boundBuffer *bytes.Buffer
}

func (r *RpcRequest[T]) Init(connection *connectionservice.GurmsConnection, fromNodeId string) {

}

func (r *RpcRequest[T]) Release() {
	r.boundBuffer = nil
}
