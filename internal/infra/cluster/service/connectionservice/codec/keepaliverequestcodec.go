package codec

import (
	"gurms/internal/infra/cluster/service/codec/pool/impl"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
	"gurms/internal/infra/io"
)

const TRACE_ID_LENGTH = 8

type KeepaliveRequestCodec struct {
}

func (r *KeepaliveRequestCodec) GetCodecId() int {
	return impl.RPC_KEEPALIVE
}

func (r *KeepaliveRequestCodec) InitialCapacity(data *dto.RpcBaseRequest) int {
	return TRACE_ID_LENGTH
}

func (k *KeepaliveRequestCodec) Write(output *io.Stream, data string) {
	//output.
}

func (k *KeepaliveRequestCodec) Read(input *io.Stream) {
	//input.
}
