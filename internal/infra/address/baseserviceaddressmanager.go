package address

import (
	"context"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/env/common"
	"time"
)

var BASESERVICEADDRESSMANAGERLOGGER logger.Logger = factory.GetLogger("BaseServiceAddressManager")

type BaseServiceAddressManager struct {
	IpDetector                        *IpDetector
	OnNodeAddressInfoChangedListeners []func(*NodeAddressInfo)
	MemberAddressProperties           *common.AddressProperties
	MemberBindHost                    string
	MemberHost                        string
	AdminApiAddressProperties         *common.AddressProperties
	AdminApiAddress                   string
}

func newBaseServiceAddressManager(
	adminHttpProperties *AdminHttpProperties,
	ipDetector *IpDetector,
	propertiesManager *property.GurmsPropertiesManager,
	serviceAddressManager ServiceAddressManager,
) *BaseServiceAddressManager {

	gurmsProperties := propertiesManager.LocalGurmsProperties
	memberBindHost := gurmsProperties.Cluster.Connection.Server.Host
	adminAddressProperties := serviceAddressManager.GetAdminAddressProperties(gurmsProperties)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60*time.Second))
	defer cancel()

	propertiesManager.AddLocalPropertiesChangeListener(func(properties *property.GurmsProperties) {
		newAdminApiDiscoveryProperties := serviceAddressManager.GetAdminAddressProperties(properties)
		areAdminAddressPropertiesChange := adminAddressProperties != newAdminApiDiscoveryProperties
		var updateAdminApiAddresses func() error
		if areAdminAddressPropertiesChange {
			updateAdminApiAddresses = func() error {
				serviceAddressManager.UpdateAdminApiAddresses(
					adminHttpProperties, newAdminApiDiscoveryProperties)
			}
		} else {
			updateAdminApiAddresses = func() error { return nil }
		}
		updateMemberHost := func() bool {
			return serviceAddressManager.UpdateMemberHostIfChanged(properties)
		}
		updateCustomAddresses := func() bool {
			return serviceAddressManager.UpdateCustomAddresses(
				adminHttpProperties, gurmsProperties)
		}
		go func() {
			err := updateAdminApiAddresses()
			if err != nil {
				BASESERVICEADDRESSMANAGERLOGGER.ErrorWithMessage(
					"caught an error while updating the node address information", err)
				return
			}
		}()
	})

	return &BaseServiceAddressManager{}
}
