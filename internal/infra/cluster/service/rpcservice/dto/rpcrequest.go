package dto

import (
	"bytes"
	"gurms/internal/infra/cluster/service/connectionservice"
)

type RpcRequest struct {
	connection  *connectionservice.GurmsConnection
	fromNodeId  string
	requestId   int
	requestTime int64
	boundBuffer *bytes.Buffer
}

func (r *RpcRequest) Init(connection *connectionservice.GurmsConnection, fromNodeId string) {

}
