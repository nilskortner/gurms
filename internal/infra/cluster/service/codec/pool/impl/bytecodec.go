package impl

import (
	"gurms/internal/infra/io"
)

type ByteCodec struct {
}

func NewByteCodec() *ByteCodec {
	return &ByteCodec{}
}

func (b *ByteCodec) GetCodecId() int {
	return PRIMITVE_BYTE
}

func (b *ByteCodec) InitialCapacity(data byte) int {
	return 1
}

func (b *ByteCodec) Write(output *io.Stream, data byte) {
	//
}

func (b *ByteCodec) Read(input *io.Stream) {
	//input.
}
