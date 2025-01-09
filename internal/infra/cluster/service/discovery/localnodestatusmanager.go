package discovery

import (
	"fmt"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/collection"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/storage/mongogurms/operation/option"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
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
	heartbeatFuncActive     bool
	IsHealthStatusUpdating  atomic.Bool
}

// region injection

type DiscoveryService interface {
	FindQualifiedMembersToBeLeader() []*configdiscovery.Member
	RegisterMember(*configdiscovery.Member) error
	GetLeader() *configdiscovery.Leader
}

type SharedConfigService interface {
	Upsert(filter *option.Filter, update *option.Update, entity string) error
	Insert(value any) (bool, error)
	UpdateOne(name string, filter *option.Filter, update *option.Update) (*mongo.UpdateResult, error)
	RemoveOne(name string, filter *option.Filter) (*mongo.DeleteResult, error)
}

// end region

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

func (n *LocalNodeStatusManager) RegisterLocalNodeAsMember(suppressDuplicateMemberError bool) error {
	LOCALNODESTATUSMANAGERLOGGER.InfoWithArgs("Registering the local node as a member")
	return n.DiscoveryService.RegisterMember(n.LocalMember)
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

func (n *LocalNodeStatusManager) UnregisterLocalMemberLeadership() (bool, error) {
	query := option.NewFilter()
	query.Eq(configdiscovery.CLUSTERIDLEADER, n.LocalMember.Key.ClusterId)
	query.Eq(configdiscovery.NODEIDLEADER, n.LocalMember.Key.NodeId)
	result, err := n.SharedConfigService.RemoveOne("leader", query)
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

func (n *LocalNodeStatusManager) isLocalNodeLeader() bool {
	leader := n.DiscoveryService.GetLeader()
	return leader != nil &&
		(leader.NodeId == n.LocalMember.Key.NodeId) &&
		(leader.ClusterId == n.LocalMember.Key.ClusterId)
}

func (n *LocalNodeStatusManager) StartHeartbeat() {
	if n.heartbeatFuncActive == true {
		return
	}
	n.heartbeatFuncActive = true
	go func() {
		for {
			if n.IsClosing {
				return
			}
			timeout := n.HeartbeatInterval * time.Second
			now := time.Now()
			update := option.NewUpdate()
			update.Set(configdiscovery.LASTHEARTBEATDATE, now)
			err := n.UpsertLocalNodeInfo(update)
			if err != nil {
				LOCALNODESTATUSMANAGERLOGGER.FatalWithError("caught an error while upserting the local node information", err)
			}
			if n.isLocalNodeLeader() {
				isLeader, err := n.renewLocalNodeAsLeader(now)
				if err != nil {
					LOCALNODESTATUSMANAGERLOGGER.FatalWithError("caught an error while renewing the local node as the leader", err)
				}
				if isLeader {
					err = n.updateMemberStatus(now)
					if err != nil {
						LOCALNODESTATUSMANAGERLOGGER.FatalWithError("caught an error while updating the information "+
							"\"lastHeatbeatDate\" of the local node", err)
					}
				}
			}
			if time.Since(now) > timeout {
				err := fmt.Errorf("timeout while sending the heartbeat request")
				LOCALNODESTATUSMANAGERLOGGER.Error(err)
			} else {
				n.LocalMember.Status.LastHeartbeatDate = now
			}
		}
	}()
}

func (n *LocalNodeStatusManager) UpdateInfo(member *configdiscovery.Member) {
	isLeaderEligible := member.IsLeaderEligible
	wasLeaderEligible := n.LocalMember.IsLeaderEligible
	isLeaderEligibleChanged := isLeaderEligible != wasLeaderEligible
	n.LocalMember = member
	if isLeaderEligibleChanged {
		if isLeaderEligible {
			_, err := n.TryBecomeFirstLeader()
			if err != nil {
				LOCALNODESTATUSMANAGERLOGGER.ErrorWithMessage("caught an error while trying to become first leader", err)
			}
		} else {
			_, err := n.UnregisterLocalMemberLeadership()
			if err != nil {
				LOCALNODESTATUSMANAGERLOGGER.ErrorWithMessage(
					"caught an error while unregistering the leadership of the local node ", err)
			}
		}
	}
}

func (n *LocalNodeStatusManager) renewLocalNodeAsLeader(renewDate time.Time) (bool, error) {
	leader := n.DiscoveryService.GetLeader()
	if leader == nil {
		return false, nil
	}
	filter := option.NewFilter()
	filter.Eq(configdiscovery.CLUSTERID, n.LocalMember.Key.ClusterId)
	filter.Eq(configdiscovery.NODEID, n.LocalMember.Key.NodeId)
	filter.Eq(configdiscovery.GENERATIONLEADER, leader.Generation)
	update := option.NewUpdate()
	update.Set(configdiscovery.ISHEALTHY, n.LocalMember.Status.IsHealthy)

	result, err := n.SharedConfigService.UpdateOne(configdiscovery.LEADERNAME, filter, update)
	if err != nil {
		return false, err
	}
	if result.MatchedCount == 0 {
		_, err = n.SharedConfigService.Insert(leader)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (n *LocalNodeStatusManager) updateMemberStatus(lastHearthbeatDate time.Time) error {
	knownMembers := n.DiscoveryService.
}
