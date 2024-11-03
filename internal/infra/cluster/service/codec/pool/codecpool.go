package pool

import (
	"fmt"
	"gurms/internal/infra/cluster/service/codec/pool/impl"
)

var ID_TO_CODEC map[int]Codec

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
	float32Codec := impl.NewFloat32Codec()
	register(float32Codec)
	intCodec := impl.NewIntCodec()
	register(intCodec)
	int64Codec := impl.NewInt64Codec()
	register(int64Codec)
	int16Codec := impl.NewInt16Codec()
	register(int16Codec)
	stringCodec := impl.NewStringCodec()
	register(stringCodec)

	// Datastructures
}

func register(codec Codec) error {
	id := codec.GetCodecId()
	_, exists := ID_TO_CODEC[id]
	if exists {
		return fmt.Errorf("the codec ID (%d) has already existed", id)
	} else {
		ID_TO_CODEC[id] = codec
		return nil
	}
}
