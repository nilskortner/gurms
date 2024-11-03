package impl

import (
	"gurms/internal/infra/io"
)

type StringCodec struct {
}

func NewStringCodec() *StringCodec {
	return &StringCodec{}
}

func (s *StringCodec) GetCodecId() int {
	return STRING
}

func (s *StringCodec) InitialCapacity(data string) int {
	return 2 + len(data) + 1
}

func (s *StringCodec) Write(output *io.Stream, data string) {
	//output.
}

func (s *StringCodec) Read(input *io.Stream) {
	//input.
}
