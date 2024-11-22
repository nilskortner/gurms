package service

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"time"
)

var DISCOVERYLOGGER logger.Logger = factory.GetLogger("Discovery")

type DiscoveryService struct {
	connectionService      *ConnectionService
	localMember            *configdiscovery.Member
	localNodeStatusManager *LocalNodeStatusManager
}

func NewDiscoveryService(
	clusterId string,
	nodeId string,
	zone string,
	name string,
	nodeType nodetype.NodeType,
	nodeversion nodetype.NodeTypeInfo,
	isLeaderEligible bool,
	priority int,
	isActive bool,
	isHealthy bool,
	memberBindPort int,
	discoveryProperties *DiscoveryProperties,
	serviceAddressManager *BaseServiceAddressManager,
	sharedConfigService *SharedConfigService,
) *DiscoveryService {
	now := time.Now()
	localMember := configdiscovery.NewMember()
	heartbeatTimeoutMillis := mathutil.multiply(discovereyProperties.getHeartbeatTimeoutSeconds(), 1000)
	localNodeStatusManager := newLocalNodeStatusManager(
		sharedConfigService,
		localMember,
		discoveryProperties.getHeartbeatIntervalSeconds(),
	)
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
