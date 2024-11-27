package channel

import (
	"bytes"
	"encoding/gob"
	"gurms/internal/infra/cluster/service/codec/pool"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
)

type RpcFrameEncoder struct{}

var INSTANCE *RpcFrameEncoder = &RpcFrameEncoder{}

func EncodeRequest[T comparable](request dto.RpcRequest[T], requestBody *bytes.Buffer) (*bytes.Buffer, error) {
	buffer := bytes.NewBuffer(make([]byte, requestBody.Len()))
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(requestBody.Bytes())
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func getCodec[T any](value T) pool.Codec {
	return pool.GetCodec(value)
}

// func EncodeRequest[T comparable](request dto.RpcRequest[T], requestBody *bytes.Buffer) (*bytes.Buffer, error) {
// 	buffer := bytes.NewBuffer(make([]byte, requestBody.Len()+6))
// 	err := binary.Write(buffer, binary.BigEndian, int64(getCodec(request).GetCodecId()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = binary.Write(buffer, binary.BigEndian, request.GetRequestId())
// 	_, err = buffer.Write(requestBody.Bytes())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return buffer, nil
// }
