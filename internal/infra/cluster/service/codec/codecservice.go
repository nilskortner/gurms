package codec

import (
	"bytes"
	"gurms/internal/infra/cluster/service/codec/pool"
	"gurms/internal/infra/cluster/service/codec/pool/impl"
	"gurms/internal/infra/io"
)

type CodecService struct {
}

func NewCodecService() *CodecService {
	pool.Init()

	return &CodecService{}
}

func Serialize(data any) *bytes.Buffer {
	var buf *bytes.Buffer
	var stream *io.Stream

	switch data := data.(type) {
	case bool:
		codecPicked := impl.NewBoolCodec()
		buf = bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity(data)))
		stream = io.NewStream(buf)
		codecPicked.Write(stream, data)
	case byte:
		codecPicked := impl.NewByteCodec()
		buf = bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity(data)))
		stream = io.NewStream(buf)
		codecPicked.Write(stream, data)
	case rune:
		codecPicked := impl.NewRuneCodec()
		buf = bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity(data)))
		stream = io.NewStream(buf)
		codecPicked.Write(stream, data)
	case float64:
		codecPicked := impl.NewFloat64Codec()
		buf = bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity(data)))
		stream = io.NewStream(buf)
		codecPicked.Write(stream, data)
	case float32:
		codecPicked := impl.NewFloat32Codec()
		buf = bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity(data)))
		stream = io.NewStream(buf)
		codecPicked.Write(stream, data)
	case int:
		codecPicked := impl.NewIntCodec()
		buf = bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity(data)))
		stream = io.NewStream(buf)
		codecPicked.Write(stream, data)
	case int64:
		codecPicked := impl.NewInt64Codec()
		buf = bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity(data)))
		stream = io.NewStream(buf)
		codecPicked.Write(stream, data)
	case int16:
		codecPicked := impl.NewInt16Codec()
		buf = bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity(data)))
		stream = io.NewStream(buf)
		codecPicked.Write(stream, data)
	case string:
		codecPicked := impl.NewStringCodec()
		buf = bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity(data)))
		stream = io.NewStream(buf)
		codecPicked.Write(stream, data)
	default:
		panic("No Codec for this Type")
	}
	return buf
}
