package rpcservice

import (
	"encoding/gob"
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
var pendingRequestMap *nonblockingmap.Map[int64, struct{}] = nonblockingmap.NewSized[int64, struct{}](uintptr(initSize))

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

// TODO: check this function and the response handling
func SendRequest[T comparable](endpoint *RpcEndpoint, request dto.RpcRequest[T]) (T, error) {
	conn := endpoint.Connection.Connection
	if conn == nil {
		var zero T
		err := fmt.Errorf("connection already closed")
		return zero, err
	}
	var response T
	for {
		requestId := generateId()
		_, ok := pendingRequestMap.GetOrInsert(requestId, struct{}{})
		if ok {
			continue
		}
		request.SetRequestId(requestId)

		buffer, err := channel.EncodeRequest(request)
		if err != nil {
			buffer = nil
			var zero T
			return zero, resolveRequest(requestId, response, err)
		}
		_, err = conn.Write(buffer.Bytes())
		if err != nil {
			var zero T
			return zero, resolveRequest(requestId, response, err)
		}
		decoder := gob.NewDecoder(conn)
		err = decoder.Decode(&response)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("error decoding response")
		}
		break
	}
	return response, nil
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
	case *dto.RpcResponse[byte]:
		HandleResponse[byte](value)
	default:
		RPCENDPOINTLOGGER.ErrorWithArgs("Couldnt resolve Instantiation of Response[Type?]: ", wrap)
	}
}

func HandleResponse[T comparable](response *dto.RpcResponse[T]) {
	resolveRequest(response.RequestId, response.Result, response.RpcError)
}

func resolveRequest[T comparable](requestId int64, response T, err error) error {
	// TODO: handle value emit
	_, ok := pendingRequestMap.Get(requestId)
	if !ok {
		message := fmt.Sprint("Could not find a pending request with the ID %s",
			requestId)
		RPCENDPOINTLOGGER.Warn(message)
		return nil
	}
	ok = pendingRequestMap.Del(requestId)
	if !ok {
		message := fmt.Sprint("Could not delete request with the ID %s",
			requestId)
		RPCENDPOINTLOGGER.Warn(message)
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
