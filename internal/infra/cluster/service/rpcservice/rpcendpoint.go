package rpcservice

import (
	"bytes"
	"fmt"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/rpcservice/channel"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"math"
	"math/rand"

	nonblockingmap "github.com/cornelk/hashmap"
)

var RPCENDPOINTLOGGER logger.Logger = factory.GetLogger("RpcEndpoint")

const (
	EXPECTED_MAX_QPS           = 1000
	EXPECTED_AVERAGE_RTT       = 10
	INITAL_CAPACITY_PERCENTAGE = 10
)

var initSize int = int(EXPECTED_MAX_QPS * EXPECTED_AVERAGE_RTT * (INITAL_CAPACITY_PERCENTAGE / 100.0))
var pendingRequestMap *nonblockingmap.Map[int64, int] = nonblockingmap.NewSized[int64, int](uintptr(initSize))

type RpcEndpoint struct {
	NodeId     string
	Connection *connectionservice.GurmsConnection
}

func NewRpcEndpoint(nodeId string, connection *connectionservice.GurmsConnection) *RpcEndpoint {
	return &RpcEndpoint{
		NodeId:     nodeId,
		Connection: connection,
	}
}

func SendRequest[T comparable](endpoint *RpcEndpoint, request *dto.RpcRequest[T], requestBody *bytes.Buffer) (T, error) {
	conn := endpoint.Connection.Connection
	if requestBody == nil {
		err := fmt.Errorf("the request body has been released")
		var zero T
		return zero, err
	}
	if conn == nil {
		err := fmt.Errorf("connection already closed")
		var zero T
		return zero, err
	}
	for {
		requestId := generateId()
		_, ok := pendingRequestMap.GetOrInsert(requestId, value)
		if ok {
			continue
		}
		request.RequestId = requestId
		var buffer *bytes.Buffer

		buffer = channel.INSTANCE.EncodeRequest(request, requestBody)

		_, err := conn.Write(buffer.Bytes())
		if err != nil {
			var zero T
			return zero, err
		}

		break
	}
}

func generateId() int64 {
	var id int64
	for {
		id = rand.Int63n(math.MaxInt64)
		_, ok := pendingRequestMap.Get(id)
		if !ok {
			break
		}
	}
	return id
}

func HandleResponse[T comparable](response *dto.RpcResponse[T]) {
	resolveRequest(response.RequestId, response.Result, response.Rpcerror)
}

func resolveRequest[T comparable](requestId int64, response T, err error) {
	ok := pendingRequestMap.Del(requestId)
	if !ok {
		message := fmt.Sprint("Could not find a pending request with the ID %s for the response: %#v",
			requestId, response)
		RPCENDPOINTLOGGER.Warn(message)
		return
	}
}
