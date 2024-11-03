package impl

import (
	"gurms/internal/infra/io"
)

type Int16Codec struct {
}

func NewInt16Codec() *Int16Codec {
	return &Int16Codec{}
}

func (*Int16Codec) GetCodecId() int {
	return PRIMITIVE_INT16
}

func (*Int16Codec) InitialCapacity(data int16) int {
	return 2
}

func (*Int16Codec) Write(output *io.Stream, data int16) {
	//output.
}

func (*Int16Codec) Read(input *io.Stream) {
	//input.
}
