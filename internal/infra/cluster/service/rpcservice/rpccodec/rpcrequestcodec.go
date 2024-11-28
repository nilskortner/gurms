package rpccodec

import (
	"gurms/internal/infra/io"
)

const TRACE_ID_LENGTH = 8

type RpcRequestCodec[T comparable] struct {
}

func NewRpcRequestCodec[T comparable]() *RpcRequestCodec[T] {
	return &RpcRequestCodec[T]{}
}

func (r *RpcRequestCodec[T]) GetCodecId() int {
	return 0
}

func (r *RpcRequestCodec[T]) Write(output *io.Stream, data bool) {
	output.WriteBoolean(data)
}

func (r *RpcRequestCodec[T]) Read(input *io.Stream) {
	//input.ReadBoolean()
}
