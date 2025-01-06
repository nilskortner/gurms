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
	query.Eq()
	query.Eq()
	return n.SharedConfigService.RemoveOne("leader", query)
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
			err := n.UpsertLocalNodeInfo()
			if err != nil {
				LOCALNODESTATUSMANAGERLOGGER.FatalWithError("caught an error while upserting the local node information", err)
			}
			if n.isLocalNodeLeader() {
				err := n.renewLocalNodeAsLeader(now)
				if err != nil {
					LOCALNODESTATUSMANAGERLOGGER.FatalWithError("caught an error while renewing the local node as the leader", err)
				}
				if n.isLeader {
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

func (n *LocalNodeStatusManager) renewLocalNodeAsLeader(renewDate time.Time) bool {
	leader := n.DiscoveryService.GetLeader()
	if leader == nil {
		return false
	}
	filter := option.Filter{}
	update := option.NewUpdate()

	result, err := SharedConfigService.UpdateOne(configdiscovery.LEADERNAME, filter, update)
	if result.getMatchedCOunt() == 0 {
		n.SharedConfigService.Insert(leader)
	}
}
