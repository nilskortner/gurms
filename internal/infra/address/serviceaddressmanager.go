package address

import (
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/env/aiserving"
	"gurms/internal/infra/property/env/common"
)

type ServiceAddressManager interface {
	GetMemberHost() string
	GetAdminApiAddress() string
	GetWsAddress() string
	GetTcpAddress() string
	GetUdpAddress() string
	AddOnNodeAddressInfoChangedListener(func(*NodeAddressInfo))
	GetAdminAddressProperties(properties *property.GurmsProperties) *common.AddressProperties
	UpdateCustomAddresses(adminHttpProperties *aiserving.AdminHttpProperties,
		properties *property.GurmsProperties) (bool, error)
	UpdateAdminApiAddresses(adminHttpProperties *aiserving.AdminHttpProperties,
		newAdminApiAddressProperties *common.AddressProperties) error
	UpdateMemberHostIfChanged(newProperties *property.GurmsProperties) (bool, error)
	notifyOnNodeAddressInfoChangedListeners(addresses *NodeAddressInfo)
}
