package impl

import "gurms/internal/infra/cluster/service/codec/pool"

type Float32Codec struct {
}

func (f Float32Codec) GetCodecId() int {
	return pool.PRIMITIVE_BOOL
}

func (f Float32Codec) InitialCapacity(data byte) int {
	return 1
}
