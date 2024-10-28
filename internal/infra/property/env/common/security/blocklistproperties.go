package security

type BlocklistProperties struct {
	ip     *IpBlocklistTypeProperties
	userId *UserIdBlocklistTypeProperties
}

func NewBlocklistProperties() *BlocklistProperties {
	return &BlocklistProperties{
		ip:     NewIpBlocklistTypeProperties(),
		userId: NewUserIdBlocklistTypeProperties(),
	}
}

type IpBlocklistTypeProperties struct {
	enabled                     bool
	syncBlocklistIntervalMillis int
	autoBlock                   *IpAutoBlocklistProperties
}

func NewIpBlocklistTypeProperties() *IpBlocklistTypeProperties {
	return &IpBlocklistTypeProperties{
		enabled:                     true,
		syncBlocklistIntervalMillis: 10 * 1000,
		autoBlock:                   NewIpAutoBlocklistProperties(),
	}
}

type UserIdBlocklistTypeProperties struct {
	enabled                     bool
	syncBlocklistIntervalMillis int
	autoBlock                   *UserIdAutoBlocklistProperties
}

func NewUserIdBlocklistTypeProperties() *UserIdBlocklistTypeProperties {
	return &UserIdBlocklistTypeProperties{
		enabled:                     true,
		syncBlocklistIntervalMillis: 10 * 1000,
		autoBlock:                   NewUserIdAutoBlocklistProperties(),
	}
}

type IpAutoBlocklistProperties struct {
	corruptedFrame   *AutoBlockItemProperties
	corruptedRequest *AutoBlockItemProperties
	frequentRequest  *AutoBlockItemProperties
}

func NewIpAutoBlocklistProperties() *IpAutoBlocklistProperties {
	return &IpAutoBlocklistProperties{
		corruptedFrame:   NewAutoBlockItemProperties(),
		corruptedRequest: NewAutoBlockItemProperties(),
		frequentRequest:  NewAutoBlockItemProperties(),
	}
}

type UserIdAutoBlocklistProperties struct {
	corruptedFrame   *AutoBlockItemProperties
	corruptedRequest *AutoBlockItemProperties
	frequentRequest  *AutoBlockItemProperties
}

func NewUserIdAutoBlocklistProperties() *UserIdAutoBlocklistProperties {
	return &UserIdAutoBlocklistProperties{
		corruptedFrame:   NewAutoBlockItemProperties(),
		corruptedRequest: NewAutoBlockItemProperties(),
		frequentRequest:  NewAutoBlockItemProperties(),
	}
}
