package service

import (
	"gurms/internal/infra/address"
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/cluster/service/discovery"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/cluster"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var DISCOVERYSERVICELOGGER logger.Logger = factory.GetLogger("DiscoveryService")

type DiscoveryService struct {
	DiscoveryProperties *cluster.DiscoveryProperties

	SharedConfigService *SharedConfigService

	ConnectionService      *ConnectionService
	LocalNodeStatusManager *discovery.LocalNodeStatusManager

	Leader *configdiscovery.Leader

	AllKnownMembers cmap.ConcurrentMap[string, *configdiscovery.Member]

	ActiveSortedAiServingMembers []*configdiscovery.Member
	ActiveSortedServiceMembers   []*configdiscovery.Member
	ActiveSortedGatewayMembers   []*configdiscovery.Member

	OtherActiveConnectedAiServingMembers []*configdiscovery.Member
	OtherActiveConnectedGatewayMembers   []*configdiscovery.Member
	OtherActiveConnectedServiceMembers   []*configdiscovery.Member
	OtherActiveConnectedMembers          []*configdiscovery.Member

	MembersChangeListeners []discovery.MembersChangeListener

	HeartbeatTimeoutMillis time.Time
}

func NewDiscoveryService(
	clusterId string,
	nodeId string,
	zone string,
	name string,
	nodeType nodetype.NodeType,
	nodeversion *nodetype.NodeVersion,
	isLeaderEligible bool,
	priority int,
	isActive bool,
	isHealthy bool,
	memberBindPort int,
	discoveryProperties *cluster.DiscoveryProperties,
	serviceAddressManager address.ServiceAddressManager,
	sharedConfigService *SharedConfigService,
) *DiscoveryService {
	now := time.Now()
	localMember := configdiscovery.NewMember(
		clusterId,
		nodeId,
		zone,
		name,
		nodeType,
		nodeversion,
		false,
		isLeaderEligible,
		now,
		priority,
		serviceAddressManager.GetMemberHost(),
		memberBindPort,
		serviceAddressManager.GetAdminApiAddress(),
		serviceAddressManager.GetWsAddress(),
		serviceAddressManager.GetTcpAddress(),
		serviceAddressManager.GetUdpAddress(),
		false,
		isActive,
		isHealthy,
	)
	heartbeatTimeoutMillis := mathutil.multiply(discovereyProperties.getHeartbeatTimeoutSeconds(), 1000)
	localNodeStatusManager := NewLocalNodeStatusManager(
		sharedConfigService,
		localMember,
		discoveryProperties.HeartbeatIntervalSeconds,
	)

	return &DiscoveryService{
		DiscoveryProperties:    discoveryProperties,
		SharedConfigService:    sharedConfigService,
		LocalNodeStatusManager: localNodeStatusManager,
		HeartbeatTimeoutMillis: heartbeatTimeoutMillis,
		// maps
	}
}

func (d *DiscoveryService) Start() {}

func (d *DiscoveryService) LazyInit(connectionService *ConnectionService) {
	//d.connectionService. = append(, func() memberconnectionlistener.MemberConnectionListener{

	//})
}

func (d *DiscoveryService) GetMember(nodeId string) *configdiscovery.Member {
	value, _ := d.AllKnownMembers.Get(nodeId)
	return value
}

func updateOtherActiveConnectedMemberList(isAdd bool, member *configdiscovery.Member) {

}

func (d *DiscoveryService) IsKnownMember(nodeId string) bool {
	return d.AllKnownMembers.Has(nodeId)
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
