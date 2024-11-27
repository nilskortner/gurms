package dto

import (
	"bytes"
	"gurms/internal/infra/cluster/service/connectionservice"
	"time"
)

type RpcRequest[T comparable] interface {
	IsAsync() bool
	CallAsync() T
	Call() T
	NodeTypeToRequest() NodeTypeToHandleRpc
	Name() string
	Init(connection *connectionservice.GurmsConnection, fromNodeId string)
	Release()
	GetRequestId() int64
	SetRequestId(int64)
}

type RpcBaseRequest struct {
	Connection  *connectionservice.GurmsConnection
	FromNodeId  string
	RequestId   int64
	RequestTime int64
	BoundBuffer *bytes.Buffer
}

func (r *RpcBaseRequest) Init(connection *connectionservice.GurmsConnection, fromNodeId string) {
	r.Connection = connection
	r.FromNodeId = fromNodeId
	r.RequestId = -1
	r.RequestTime = time.Now().UnixMilli()
	r.BoundBuffer = bytes.NewBuffer(make([]byte, 0))
}

func (r *RpcBaseRequest) Release() {
	r.BoundBuffer = nil
}

func (r *RpcBaseRequest) GetRequestId() int64 {
	return r.RequestId
}

func (r *RpcBaseRequest) SetRequestId(requestId int64) {
	r.RequestId = requestId
}
