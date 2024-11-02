package codec

import (
	"bytes"
	"gurms/internal/infra/cluster/service/codec/pool"
	"gurms/internal/infra/cluster/service/codec/pool/impl"
)

type CodecService struct {
	pool *pool.Pool
}

func NewCodecService() *CodecService {
	newPool := pool.NewPool()
	pool.Init()

	return &CodecService{
		pool: newPool,
	}
}

func SerializeWithoutId(data any) *bytes.Buffer {
	var codecPicked pool.Codec
	switch data.(type) {
	case bool:
		codecPicked := impl.BoolCodec{}
		//protobuf codecPicked
	default:
		panic("No Codec for this Type")
	}
	buf := bytes.NewBuffer(make([]byte, 0, codecPicked.InitialCapacity()))

	initialCapacity := codecPicked.InitialCapacity
}
