package common

import (
	"gurms/internal/infra/property/env/common/redisproperties"
)

type BaseRedisProperties struct {
	Session         *redisproperties.RedisProperties
	Location        *redisproperties.RedisProperties
	IpBlocklist     *redisproperties.SimpleRedisProperties
	UserIdBlocklist *redisproperties.SimpleRedisProperties
}
