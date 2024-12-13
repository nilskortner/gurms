package cluster

import "gurms/internal/infra/property/env/common"

type DiscoveryProperties struct {
	HeartbeatTimeoutSeconds          int
	HeartbeatIntervalSeconds         int
	DelayToNotifyMemberChangeSeconds int
	Address                          *common.AddressProperties
}

func NewDiscoveryProperties() *DiscoveryProperties {
	return &DiscoveryProperties{
		HeartbeatTimeoutSeconds:          30,
		HeartbeatIntervalSeconds:         10,
		DelayToNotifyMemberChangeSeconds: 3,
		Address:                          common.NewAddressProperties(),
	}
}
