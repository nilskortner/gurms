package cluster

import "gurms/internal/infra/property/env/common"

type DiscoveryProperties struct {
	heartbeatTimeoutSeconds          int
	heartbeatIntervalSeconds         int
	delayToNotifyMemberChangeSeconds int
	address                          *common.AddressProperties
}

func NewDiscoveryProperties() *DiscoveryProperties {
	return &DiscoveryProperties{
		heartbeatTimeoutSeconds:          30,
		heartbeatIntervalSeconds:         10,
		delayToNotifyMemberChangeSeconds: 3,
		address:                          common.NewAddressProperties(),
	}
}
