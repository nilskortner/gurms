package rpcservice

import (
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
var pendingRequestMap *nonblockingmap.Map[int64, chan any] = nonblockingmap.NewSized[int64, chan any](uintptr(initSize))

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

func SendRequest[T comparable](endpoint *RpcEndpoint, request dto.RpcRequest[T]) (chan any, error) {
	conn := endpoint.Connection.Connection
	if conn == nil {
		err := fmt.Errorf("connection already closed")
		return nil, err
	}
	sink := make(chan any)
	for {
		requestId := generateId()
		_, ok := pendingRequestMap.GetOrInsert(requestId, sink)
		if ok {
			continue
		}
		request.SetRequestId(requestId)

		buffer, err := channel.EncodeRequest(request)
		if err != nil {
			buffer = nil
			var zero T
			resolveRequest(requestId, zero, err)
			return nil, err
		}
		_, err = conn.Write(buffer.Bytes())
		if err != nil {
			var zero T
			resolveRequest(requestId, zero, err)
			return nil, err
		}
		break
	}
	return sink, nil
}

func sendRequest() {}

func sendRequestAsync() {}

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

// TODO: add all response types
func UnwrapResponse(wrap dto.RpcResponseWrap) {
	// Type switch
	switch value := wrap.(type) {
	case *dto.RpcResponse:
		HandleResponse(value)
	default:
		RPCENDPOINTLOGGER.ErrorWithArgs("Couldnt resolve Instantiation of Response[Type?]: ", wrap)
	}
}

func HandleResponse(response *dto.RpcResponse) {
	resolveRequest(response.RequestId, response.Result, response.RpcError)
}

func resolveRequest(requestId int64, response any, err error) {
	sink, ok := pendingRequestMap.Get(requestId)
	if !ok {
		message := fmt.Sprint("Could not find a pending request with the ID %s for the response: %s",
			requestId, response)
		RPCENDPOINTLOGGER.Warn(message)
		return
	}
	if err == nil {
		sink <- response
	} else {
		sink <- err
	}
}
