package discovery

import (
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/collection"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/storage/mongogurms/operation/option"
	"sync/atomic"
	"time"
)

var LOCALNODESTATUSMANAGERLOGGER logger.Logger = factory.GetLogger("LocalNodeStatusManager")

type LocalNodeStatusManager struct {
	DiscoveryService        DiscoveryService
	SharedConfigService     SharedConfigService
	LocalMember             *configdiscovery.Member
	IsLocalNodeRegistered   bool
	IsClosing               bool
	HeartbeatInterval       time.Duration
	HeartbeatIntervalMillis int64
	IsHealthStatusUpdating  atomic.Bool
}

type DiscoveryService interface {
	FindQualifiedMembersToBeLeader() []*configdiscovery.Member
}

type SharedConfigService interface {
	Upsert(filter *option.Filter, update *option.Update, entity string) error
	Insert(leader configdiscovery.Leader) (bool, error)
}

func NewLocalNodeStatusManager(
	discoveryService DiscoveryService,
	sharedConfigService SharedConfigService,
	localMember *configdiscovery.Member,
	heartbeatIntervalSeconds int,
) *LocalNodeStatusManager {
	heartbeatInterval := time.Duration(heartbeatIntervalSeconds) * time.Second
	heartbeatIntervalMillis := int64(heartbeatInterval.Milliseconds())
	return &LocalNodeStatusManager{
		DiscoveryService:        discoveryService,
		SharedConfigService:     sharedConfigService,
		LocalMember:             localMember,
		HeartbeatInterval:       heartbeatInterval,
		HeartbeatIntervalMillis: heartbeatIntervalMillis,
	}
}

func (n *LocalNodeStatusManager) UpsertLocalNodeInfo(update *option.Update) error {
	nodeId := n.LocalMember.Key.NodeId
	memberFilter := option.NewFilter()
	memberFilter.Eq("_id."+n.LocalMember.Key.ClusterId, n.LocalMember.Key.ClusterId)
	memberFilter.Eq("_id."+nodeId, nodeId)

	err := n.SharedConfigService.Upsert(memberFilter, update, n.LocalMember.Name)
	if err == nil {
		n.IsLocalNodeRegistered = true
		return nil
	}
	return err
}

func (n *LocalNodeStatusManager) TryBecomeFirstLeader() (bool, error) {
	qualifiedMembersToBeLeader := n.DiscoveryService.FindQualifiedMembersToBeLeader()
	if !collection.Contains(qualifiedMembersToBeLeader, n.LocalMember) {
		return false, nil
	}
	clusterId := n.LocalMember.Key.ClusterId
	localLeader := configdiscovery.NewLeader(clusterId, n.LocalMember.Key.NodeId, time.Now(), 1)
	return n.SharedConfigService.Insert(localLeader)
}
