package channel

import (
	"bytes"
	"encoding/gob"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

func DecodeRequest(data []byte) (dto.RpcRequest, error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	buffer.Write(data)
	dec := gob.NewDecoder(buffer)
	var request dto.RpcRequest
	err := dec.Decode(request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeResponse(data []byte) (*dto.RpcResponse, error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	buffer.Write(data)
	dec := gob.NewDecoder(buffer)
	var response *dto.RpcResponse
	err := dec.Decode(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
