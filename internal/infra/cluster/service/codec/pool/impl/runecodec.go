package impl

import (
	"gurms/internal/infra/cluster/service/codec/codecio"
)

type RuneCodec struct {
}

func NewRuneCodec() *RuneCodec {
	return &RuneCodec{}
}

func (r *RuneCodec) GetCodecId() int {
	return PRIMITIVE_RUNE
}

func (r *RuneCodec) InitialCapacity(data ) int {
	return 4
}

func (r *RuneCodec) Write(output codecio.CodecStream, data rune ) {
	output.
}

func (r *RuneCodec) Read(input codecio.CodecStream) {
	input.
}
