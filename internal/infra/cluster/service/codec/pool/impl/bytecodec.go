package impl

import "gurms/internal/infra/cluster/service/codec/pool"

type ByteCodec struct {
}

func (b ByteCodec) GetCodecId() int {
	return pool.PRIMITVE_BYTE
}

func (b ByteCodec) InitialCapacity(data byte) int {
	return 1
}
