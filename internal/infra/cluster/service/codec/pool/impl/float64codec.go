package impl

import (
	"gurms/internal/infra/io"
)

type Float64Codec struct {
}

func NewFloat64Codec() *Float64Codec {
	return &Float64Codec{}
}

func (*Float64Codec) GetCodecId() int {
	return PRIMITIVE_FLOAT64
}

func (*Float64Codec) InitialCapacity(data float64) int {
	return 8
}

func (*Float64Codec) Write(output *io.Stream, data float64) {
	//output.
}

func (*Float64Codec) Read(input *io.Stream) {
	//input.
}
