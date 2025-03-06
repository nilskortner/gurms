package sharding

type ModShardingAlgorithm struct {
}

func (m *ModShardingAlgorithm) DoSharding(shardKey int64, serverCount int) int {
	return int(shardKey % int64(serverCount))
}
