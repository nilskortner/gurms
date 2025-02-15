package address

import (
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/adminapi"
)

type ServiceAddressManager interface {
	GetMemberHost() string
	GetAdminApiAddress() string
	GetWsAddress() string
	GetTcpAddress() string
	GetUdpAddress() string
	AddOnNodeAddressInfoChangedListener(func(*NodeAddressInfo))
	GetAdminAddressProperties(properties *property.GurmsProperties) *common.AddressProperties
	UpdateCustomAddresses(adminHttpProperties *adminapi.AdminHttpProperties,
		properties *property.GurmsProperties) (bool, error)
	UpdateAdminApiAddresses(adminHttpProperties *adminapi.AdminHttpProperties,
		newAdminApiAddressProperties *common.AddressProperties) error
	UpdateMemberHostIfChanged(newProperties *property.GurmsProperties) (bool, error)
	notifyOnNodeAddressInfoChangedListeners(addresses *NodeAddressInfo)
}
