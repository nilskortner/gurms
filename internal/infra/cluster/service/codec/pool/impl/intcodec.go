package impl

import (
	"gurms/internal/infra/io"
	"unsafe"
)

const THISSYSTEMINT = int(unsafe.Sizeof(int(0) * 8))

type IntCodec struct {
}

func NewIntCodec() *IntCodec {
	return &IntCodec{}
}

func (i *IntCodec) GetCodecId() int {
	return 
}

func (i *IntCodec) InitialCapacity(data int) int {
	return THISSYSTEMINT
}

func (i *IntCodec) Write(output io.Stream, data int) {
	output.
}

func (i *IntCodec) Read(input io.Stream) {
	input.
}
