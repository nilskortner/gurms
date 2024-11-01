package codec

import "gurms/internal/infra/cluster/service/codec/pool"

type CodecService struct {
	pool *pool.Pool
}

func NewCodecService() *CodecService {
	pool := pool.NewPool()

	return &CodecService{
		pool: pool,
	}
}

func (c *CodecService) Register() {

}
