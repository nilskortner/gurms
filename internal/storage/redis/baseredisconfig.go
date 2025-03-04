package redis

type BaseRedisConfig struct {
	sessionRedisClientManager  *GurmsRedisClientManager
	locationRedisClientManager *GurmsRedisClientManager

	ipBlocklistRedisClient     *GurmsRedisClient
	userIdBlocklistRedisClient *GurmsRedisClient

	registeredClientManagers []GurmsRedisClientManager
	registeredClients        []GurmsRedisClient
}

func NewBaseRedisConfig(redisProperties *BaseRedisProperties,
	treatUserIdAndDeviceTypeAsUniqueUser bool) *BaseRedisConfig {

}
