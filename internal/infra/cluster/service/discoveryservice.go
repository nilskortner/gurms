package service

import "gurms/internal/infra/cluster/service/config/entity/configdiscovery"

type DiscoveryService struct {
	connectionService *ConnectionService
}

func (d *DiscoveryService) LazyInit(connectionService *ConnectionService) {
	//d.connectionService. = append(, func() memberconnectionlistener.MemberConnectionListener{

	//})
}

func updateOtherActiveConnectedMemberList(isAdd bool, member *configdiscovery.Member) {

}

// region MemberConnectionListener

type DiscoveryMemberConnectionListener struct {
	member *configdiscovery.Member
}

func (d *DiscoveryMemberConnectionListener) OnOpeningHandshakeCompleted(member *configdiscovery.Member) {
	d.member = member
	updateOtherActiveConnectedMemberList(true, d.member)
}

func (d *DiscoveryMemberConnectionListener) OnConnectionClosed() {
	if d.member != nil {
		updateOtherActiveConnectedMemberList(false, d.member)
	}
}
