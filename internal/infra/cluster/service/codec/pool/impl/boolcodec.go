package impl

import "gurms/internal/infra/cluster/service/codec/pool"

type BoolCodec struct {
}

func (b BoolCodec) GetCodecId() int {
	return pool.PRIMITIVE_BOOL
}

func Encode(output CodecStreamOutput, data bool) {
	output.EncodeBoolean(data)
}

func (b BoolCodec) InitialCapacity(data byte) int {
	return 1
}
