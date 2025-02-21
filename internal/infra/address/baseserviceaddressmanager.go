package address

import (
	"context"
	"fmt"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/adminapi"
	"sync"
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
	adminHttpProperties *adminapi.AdminHttpProperties,
	ipDetector *IpDetector,
	propertiesManager *property.GurmsPropertiesManager,
	serviceAddressManager ServiceAddressManager,
) *BaseServiceAddressManager {

	onNodeAddressInfoChangedListeners := make([]func(*NodeAddressInfo), 0)
	gurmsProperties := propertiesManager.LocalGurmsProperties
	memberBindHost := gurmsProperties.Cluster.Connection.Server.Host
	adminAddressProperties := serviceAddressManager.GetAdminAddressProperties(gurmsProperties)
	manager := &BaseServiceAddressManager{
		IpDetector:                        ipDetector,
		OnNodeAddressInfoChangedListeners: onNodeAddressInfoChangedListeners,
		MemberBindHost:                    memberBindHost,
	}
	// parallel setup
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60*time.Second))
	defer cancel()
	done := make(chan struct{})

	go func() {
		defer wg.Done()
		serviceAddressManager.UpdateMemberHostIfChanged(gurmsProperties)
	}()
	go func() {
		defer wg.Done()
		serviceAddressManager.UpdateAdminApiAddresses(adminHttpProperties, adminAddressProperties)
	}()
	go func() {
		defer wg.Done()
		serviceAddressManager.UpdateCustomAddresses(adminHttpProperties, gurmsProperties)
	}()
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		BASESERVICEADDRESSMANAGERLOGGER.Fatal(
			"timeout while first updating the node address information")
	}

	propertiesManager.AddLocalPropertiesChangeListener(func(properties *property.GurmsProperties) {
		newAdminApiDiscoveryProperties := serviceAddressManager.GetAdminAddressProperties(properties)
		areAdminAddressPropertiesChange := adminAddressProperties != newAdminApiDiscoveryProperties
		var updateAdminApiAddresses func() error
		if areAdminAddressPropertiesChange {
			updateAdminApiAddresses = func() error {
				return serviceAddressManager.UpdateAdminApiAddresses(
					adminHttpProperties, newAdminApiDiscoveryProperties)
			}
		} else {
			updateAdminApiAddresses = func() error { return nil }
		}
		updateMemberHost := func() (bool, error) {
			return serviceAddressManager.UpdateMemberHostIfChanged(properties)
		}
		updateCustomAddresses := func() (bool, error) {
			return serviceAddressManager.UpdateCustomAddresses(
				adminHttpProperties, gurmsProperties)
		}

		go func() {
			// parallel setup
			var wg sync.WaitGroup
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(60*time.Second))
			defer cancel()
			done := make(chan struct{})
			isMemberHostChangedChan := make(chan bool)
			areCustomAddressesChangedChan := make(chan bool)
			errors := make(chan error, 3)
			wg.Add(3)

			go func() {
				defer wg.Done()
				err := updateAdminApiAddresses()
				errors <- err
			}()
			go func() {
				defer wg.Done()
				isMemberHostChanged, err := updateMemberHost()
				isMemberHostChangedChan <- isMemberHostChanged
				errors <- err
				close(isMemberHostChangedChan)
			}()
			go func() {
				defer wg.Done()
				areCustomAddressesChanged, err := updateCustomAddresses()
				areCustomAddressesChangedChan <- areCustomAddressesChanged
				errors <- err
				close(areCustomAddressesChangedChan)
			}()
			go func() {
				wg.Wait()
				close(done)
				close(errors)
			}()
			select {
			case <-done:
			case <-ctx.Done():
				BASESERVICEADDRESSMANAGERLOGGER.ErrorWithArgs(
					"timeout while updating the node address information")
				return
			}
			isErr := false
			for err := range errors {
				if err != nil {
					BASESERVICEADDRESSMANAGERLOGGER.ErrorWithMessage(
						"caught an error while updating the node address information", err)
					isErr = true
				}
			}
			if isErr {
				return
			}
			isMemberHostChanged := <-isMemberHostChangedChan
			areCustomAddressesChanged := <-areCustomAddressesChangedChan

			if areAdminAddressPropertiesChange || isMemberHostChanged || areCustomAddressesChanged {
				addressInfo := NewAddressInfo(serviceAddressManager.GetMemberHost(),
					serviceAddressManager.GetAdminApiAddress(), serviceAddressManager.GetWsAddress(),
					serviceAddressManager.GetTcpAddress(), serviceAddressManager.GetUdpAddress())
				serviceAddressManager.notifyOnNodeAddressInfoChangedListeners(addressInfo)
			}
		}()
	})
	return manager
}

func (b *BaseServiceAddressManager) notifyOnNodeAddressInfoChangedListeners(
	addresses *NodeAddressInfo) {

	for _, listener := range b.OnNodeAddressInfoChangedListeners {
		listener(addresses)
	}
}

func (b *BaseServiceAddressManager) AddOnNodeAddressInfoChangedListener(func(*NodeAddressInfo)) {

}

func (b *BaseServiceAddressManager) GetMemberHost() string {
	return ""
}
func (b *BaseServiceAddressManager) GetAdminApiAddress() string {
	return ""
}
func (b *BaseServiceAddressManager) GetWsAddress() string {
	return ""
}
func (b *BaseServiceAddressManager) GetTcpAddress() string {
	return ""
}
func (b *BaseServiceAddressManager) GetUdpAddress() string {
	return ""
}
func (b *BaseServiceAddressManager) GetAdminAddressProperties(properties *property.GurmsProperties) *common.AddressProperties {
	return &common.AddressProperties{}
}
func (b *BaseServiceAddressManager) UpdateAdminApiAddresses(adminHttpProperties *adminapi.AdminHttpProperties,
	newAdminApiAddressProperties *common.AddressProperties) error {
	return fmt.Errorf("")
}
func (b *BaseServiceAddressManager) UpdateCustomAddresses(adminHttpProperties *adminapi.AdminHttpProperties,
	properties *property.GurmsProperties) (bool, error) {
	return false, fmt.Errorf("")
}
func (b *BaseServiceAddressManager) UpdateMemberHostIfChanged(newProperties *property.GurmsProperties) (bool, error) {
	return false, fmt.Errorf("")
}
