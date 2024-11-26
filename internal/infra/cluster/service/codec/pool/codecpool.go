package pool

import (
	"fmt"
	"gurms/internal/infra/cluster/service/codec/pool/impl"
	"gurms/internal/infra/cluster/service/rpcservice/dto"
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

func GetCodec(value any) Codec {

	switch value.(type) {
	case bool:
		return ID_TO_CODEC[impl.PRIMITIVE_BOOL]
	case byte:
		return ID_TO_CODEC[impl.PRIMITVE_BYTE]
	case rune:
		return ID_TO_CODEC[impl.PRIMITIVE_RUNE]
	case float64:
		return ID_TO_CODEC[impl.PRIMITIVE_FLOAT64]
	case float32:
		return ID_TO_CODEC[impl.PRIMITIVE_FLOAT32]
	case int:
		return ID_TO_CODEC[impl.PRIMITIVE_INT]
	case int64:
		return ID_TO_CODEC[impl.PRIMITIVE_INT64]
	case int16:
		return ID_TO_CODEC[impl.PRIMITIVE_INT16]
	case string:
		return ID_TO_CODEC[impl.STRING]
	case *dto.RpcRequest[T]:
		return ID_TO_CODEC[impl.RPC_KEEPALIVE]
	default:
		panic("no codec for this type")
	}
}
