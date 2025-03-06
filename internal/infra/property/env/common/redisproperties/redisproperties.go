package redisproperties

import "gurms/internal/storage/redis/sharding"

type RedisProperties struct {
	uriList           []string
	shardingAlgorithm sharding.ShardingAlgorithm
}

func NewRedisProperties() *RedisProperties {
	list := make([]string, 1)
	list[0] = "redis:://localhost"
	return &RedisProperties{
		uriList:           list,
		shardingAlgorithm: &sharding.ConsistentHashingAlgorithm{},
	}
}
