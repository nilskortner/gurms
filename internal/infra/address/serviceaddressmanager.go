package address

import (
	"gurms/internal/infra/property"
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
	UpdateCustomAddresses(adminHttpProperties *AdmingHttpProperties,
		properties *property.GurmsProperties) bool
	UpdateAdminApiAddresses(adminHttpProperties *AdmingHttpProperties,
		newAdminApiAddressProperties *AddressManager) error
	UpdateMemberHostIfChanged(newProperties *property.GurmsProperties) bool
}
