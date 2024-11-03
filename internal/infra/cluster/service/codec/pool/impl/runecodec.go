package impl

import (
	"gurms/internal/infra/io"
)

type RuneCodec struct {
}

func NewRuneCodec() *RuneCodec {
	return &RuneCodec{}
}

func (r *RuneCodec) GetCodecId() int {
	return PRIMITIVE_RUNE
}

func (r *RuneCodec) InitialCapacity(data rune) int {
	return 4
}

func (r *RuneCodec) Write(output *io.Stream, data rune) {
	//output.
}

func (r *RuneCodec) Read(input *io.Stream) {
	//input.
}
