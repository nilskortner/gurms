package security

var DEFAULT_BLOCK_LEVELS [3]*BlockLevel

type AutoBlockItemProperties struct {
	enabled           bool
	blockTriggerTimes int
	blockLevels       []*BlockLevel
}

func NewAutoBlockItemProperties() *AutoBlockItemProperties {
	DEFAULT_BLOCK_LEVELS[0] = NewBlockLevelParam(10*60, 60*1000, 1)
	DEFAULT_BLOCK_LEVELS[1] = NewBlockLevelParam(30*60, 60*1000, 1)
	DEFAULT_BLOCK_LEVELS[2] = NewBlockLevelParam(60*60, 60*1000, 0)

	return &AutoBlockItemProperties{
		enabled:           false,
		blockTriggerTimes: 5,
		blockLevels:       DEFAULT_BLOCK_LEVELS[:],
	}
}

type BlockLevel struct {
	blockDurationsSeconds              int64
	reduceOneTriggerTimeIntervalMillis int
	goNextLevelTriggerTimes            int
}

func NewBlockLevelParam(block int64, reduce int, goNext int) *BlockLevel {
	return &BlockLevel{
		blockDurationsSeconds:              block,
		reduceOneTriggerTimeIntervalMillis: reduce,
		goNextLevelTriggerTimes:            goNext,
	}
}
