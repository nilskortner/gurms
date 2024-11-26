package channel

import (
	"bytes"
	"encoding/binary"
	"gurms/internal/infra/cluster/service/codec/pool"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

type RpcFrameEncoder struct{}

var INSTANCE *RpcFrameEncoder = &RpcFrameEncoder{}

func EncodeRequest[T comparable](request *dto.RpcRequest[T], requestBody *bytes.Buffer) (*bytes.Buffer, error) {
	buffer := bytes.NewBuffer(make([]byte, requestBody.Len()+6))
	err := binary.Write(buffer, binary.BigEndian, int64(getCodec(request).GetCodecId()))
	err = binary.Write(buffer, binary.BigEndian, request.RequestId)
	_, err = buffer.Write(requestBody.Bytes())
	return buffer
}

func getCodec[T any](value T) pool.Codec {
	return pool.GetCodec(value)
}
