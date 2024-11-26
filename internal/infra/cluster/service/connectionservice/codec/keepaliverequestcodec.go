package codec

import (
	"gurms/internal/infra/cluster/service/codec/pool/impl"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
	"gurms/internal/infra/cluster/service/rpcservice/rpccodec"
)

const TRACE_ID_LENGTH = 8

type KeepaliveRequestCodec[T comparable] struct {
	*rpccodec.RpcRequestCodec[T]
}

func (r *KeepaliveRequestCodec[T]) GetCodecId() int {
	return impl.RPC_KEEPALIVE
}

func (r *KeepaliveRequestCodec[T]) InitialCapacity(data *dto.RpcRequest[T]) int {
	return TRACE_ID_LENGTH
}

// func (r *KeepaliveRequestCodec[T]) ReadRequestData(input *CodecStreamInput) *request.KeepaliveRequest[T] {
// 	return request.NewKeepAliveRequest[T]()
// }
