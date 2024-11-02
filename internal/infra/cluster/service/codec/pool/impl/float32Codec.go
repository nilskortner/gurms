package impl

import "gurms/internal/infra/io"

type Float32Codec struct {
}

func NewFloat32Codec() *Float32Codec {
	return &Float32Codec{}
}

func (*Float32Codec) GetCodecId() int {
	return PRIMITIVE_FLOAT32
}

func (*Float32Codec) InitialCapacity(data float32) int {
	return 4
}

func (*Float32Codec) Write(output io.Stream, data float32) {
	output.
}

func (*Float32Codec) Read(input io.Stream) {
	input.
}
