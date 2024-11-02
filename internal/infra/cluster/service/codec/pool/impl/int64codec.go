package impl

import "gurms/internal/infra/io"

type Int64Codec struct {
}

func NewInt64Codec() *Int64Codec {
	return &Int64Codec{}
}

func (*Int64Codec) GetCodecId() int {
	return PRIMITIVE_INT64
}

func (*Int64Codec) InitialCapacity(data int64) int {
	return 8
}

func (*Int64Codec) Write(output io.Stream, data int64) {
	output.
}

func (*Int64Codec) Read(input io.Stream) {
	input.
}
