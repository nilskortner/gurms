package sharding

type ShardingAlgorithm interface {
	DoSharding(shardKey int64, serverSize int) int
}
