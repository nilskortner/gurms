package sharding

const (
	SLOT_COUNT       = 1 << 16
	SLOT_MASK  int64 = SLOT_COUNT - 1
)

type ConsistentHashingAlgorithm struct {
}

func (c *ConsistentHashingAlgorithm) DoSharding(shardKey int64, serverCount int) int {
	slot := int(shardKey & SLOT_MASK)
	span := SLOT_COUNT / serverCount
	return slot / span
}
