package impl

import (
	"gurms/internal/infra/io"
)

type BoolCodec struct {
}

func NewBoolCodec() *BoolCodec {
	return &BoolCodec{}
}

func (b BoolCodec) GetCodecId() int {
	return PRIMITIVE_BOOL
}

func (b BoolCodec) InitialCapacity(data bool) int {
	return 1
}

func (b BoolCodec) Write(output *io.Stream, data bool) {
	//output.WriteBoolean(data)
}

func (b BoolCodec) Read(input *io.Stream) {
	//input.ReadBoolean()
}
