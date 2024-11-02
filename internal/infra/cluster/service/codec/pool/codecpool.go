package pool

import (
	"fmt"
	"gurms/internal/infra/cluster/service/codec/pool/impl"
)

var ID_TO_CODEC map[int]Codec

type Pool struct {
}

func NewPool() *Pool {
	return &Pool{}
}

func Init() {
	ID_TO_CODEC = make(map[int]Codec, 32)

	if !(len(ID_TO_CODEC) == 0) {
		return
	}
	//Primitives
	boolCodec := impl.NewBoolCodec()
	register(boolCodec)
	byteCodec := impl.NewByteCodec()
	register(byteCodec)
	runeCodec := impl.NewRuneCodec()
	register(runeCodec)
	float64Codec := impl.NewFloat64Codec()
	register(float64Codec)

	// Datastructures
}

func register(codec Codec) error {
	id := codec.GetCodecId()
	_, exists := ID_TO_CODEC[id]
	if exists {
		return fmt.Errorf("The codec ID (%w) has already existed", id)
	} else {
		ID_TO_CODEC[id] = codec
		return nil
	}
}
