package channel

import (
	"bytes"
	"encoding/gob"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

func DecodeRequest[T comparable](data []byte) (*dto.RpcRequest[T], error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	buffer.Write(data)
	dec := gob.NewDecoder(buffer)
	var request *dto.RpcRequest[T]
	err := dec.Decode(request)
	if err != nil {
		return nil, err
	}
	return request, nil
}
