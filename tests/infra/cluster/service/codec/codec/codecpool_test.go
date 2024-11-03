package codec_test

import (
	"gurms/internal/infra/cluster/service/codec/pool"
	"testing"
)

func TestCodecPoolInit(t *testing.T) {
	pool.Init()

	t.Log(pool.ID_TO_CODEC)
}
