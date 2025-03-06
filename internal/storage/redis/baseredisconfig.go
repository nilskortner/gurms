package redis

import "gurms/internal/infra/property/env/common"

type BaseRedisConfig struct {
	sessionRedisClientManager  *GurmsRedisClientManager
	locationRedisClientManager *GurmsRedisClientManager

	ipBlocklistRedisClient     *GurmsRedisClient
	userIdBlocklistRedisClient *GurmsRedisClient

	registeredClientManagers []GurmsRedisClientManager
	registeredClients        []GurmsRedisClient
}

func NewBaseRedisConfig(redisProperties *common.BaseRedisProperties,
	treatUserIdAndDeviceTypeAsUniqueUser bool) *BaseRedisConfig {

	sessionRedisClientManager := NewSessionRedisClientManager(redisProperties.Session)
	locationRedisClientManager := NewLocationRedisClientManager(redisProperties.Location)
	ipBlocklistRedisClient := NewIpBlocklistRedisClient(redisProperties.IpBlocklist.Uri)
	userIdBlocklistRedisClient := NewUserIdBlocklistRedisClient(redisProperties.UserIdBlocklist.Uri)

	return &BaseRedisConfig{}
}
