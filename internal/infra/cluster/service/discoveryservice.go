package service

import (
	"context"
	"fmt"
	"gurms/internal/infra/address"
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/cluster/service/discovery"
	"gurms/internal/infra/collection"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/cluster"
	"gurms/internal/storage/mongogurms/operation/option"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DISCOVERYSERVICELOGGER logger.Logger = factory.GetLogger("DiscoveryService")

const CRUD_TIMEOUT_DURATION = 1 * time.Minute
const (
	INSERT     = "insert"
	REPLACE    = "replace"
	UPDATE     = "update"
	DELETE     = "delete"
	INVALIDATE = "invalidate"
)

type DiscoveryService struct {
	notifyMembersChangeFuture *Future

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

	HeartbeatTimeoutMillis int64

	cancelLeaderChangeRoutine context.CancelFunc
	cancelMemberChangeRoutine context.CancelFunc

	mu sync.Mutex
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
	heartbeatTimeoutMillis := (discoveryProperties.HeartbeatTimeoutSeconds * 1000)

	discoveryService := &DiscoveryService{
		notifyMembersChangeFuture: &Future{},
		DiscoveryProperties:       discoveryProperties,
		SharedConfigService:       sharedConfigService,
		HeartbeatTimeoutMillis:    int64(heartbeatTimeoutMillis),
		// maps
		ActiveSortedAiServingMembers:         make([]*configdiscovery.Member, 0),
		ActiveSortedServiceMembers:           make([]*configdiscovery.Member, 0),
		ActiveSortedGatewayMembers:           make([]*configdiscovery.Member, 0),
		OtherActiveConnectedAiServingMembers: make([]*configdiscovery.Member, 0),
		OtherActiveConnectedGatewayMembers:   make([]*configdiscovery.Member, 0),
		OtherActiveConnectedServiceMembers:   make([]*configdiscovery.Member, 0),
		OtherActiveConnectedMembers:          make([]*configdiscovery.Member, 0),
		MembersChangeListeners:               make([]discovery.MembersChangeListener, 0),
	}

	localNodeStatusManager := discovery.NewLocalNodeStatusManager(
		discoveryService,
		sharedConfigService,
		localMember,
		discoveryProperties.HeartbeatIntervalSeconds,
	)
	discoveryService.LocalNodeStatusManager = localNodeStatusManager

	serviceAddressManager.AddOnNodeAddressInfoChangedListener(func(info *address.NodeAddressInfo) {

		update := option.NewUpdate(info.MemberHost,
			info.AdminApiAddress,
			info.WsAddress,
			info.TcpAddress,
			info.UdpAddress)
		err := localNodeStatusManager.UpsertLocalNodeInfo(update)
		if err != nil {
			DISCOVERYSERVICELOGGER.ErrorWithMessage("caught an error while upserting the local node info", err)
		}
	})

	return discoveryService
}

func compareMemberPriority(m1, m2 *configdiscovery.Member) int {
	m1Prio := m1.Priority
	m2Prio := m2.Priority
	if m1Prio == m2Prio {
		if m1.Key.NodeId < m2.Key.NodeId {
			return -1
		} else {
			return 1
		}
	}
	if m1Prio < m2Prio {
		return -1
	} else {
		return 1
	}
}

func (d *DiscoveryService) Start() {
	d.listenLeaderChangeEvent()

	//Members
	d.listenMembersChangeEvent()
	var memberList []*configdiscovery.Member
	memberList = d.queryMembers()
	time.Sleep(CRUD_TIMEOUT_DURATION)

	localMember := d.LocalNodeStatusManager.LocalMember
	var isSameId bool
	var isSameAddress bool
	for _, member := range memberList {
		isSameId = localMember.IsSameId(member)
		isSameAddress = localMember.IsSameAddress(member)
		if isSameId || isSameAddress {
			if !d.isAvailableMember(member, time.Now()) {
				var removedMemberIfInavailable bool
				if isSameId {
					d.removeMemberIfInavailable(member.Key.NodeId, "", 0)
				} else {
					d.removeMemberIfInavailable(member.Key.NodeId,
						member.MemberHost,
						member.MemberPort)
				}
				isConflictedNodeRemoved := removedMemberIfInavailable
				if isConflictedNodeRemoved {
					continue
				}
			}
			err := fmt.Errorf("Failed to bootstrap the local node because the local node has been registered. "+
				"Local Node: %s. Registered Node: %s", localMember, member)
			return err
		}
		d.onMemberAddedOrReplaced(member)
	}
	//
	d.onMemberAddedOrReplaced(localMember)
	d.updateActiveMembers(d.AllKnownMembers.Items())

	err := d.LocalNodeStatusManager.registerLocalNodeAsMember(false)
	if err != nil {
		fmt.Errorf("Caught an error while registering the local node as a member", err)
	}
	isLeader, err := d.LocalNodeStatusManager.tryBecomeFirstLeader()
	if err != nil {
		fmt.Errorf("Caught an error while trying to become the first leader", err)
	}
	if isLeader {
		DISCOVERYSERVICELOGGER.InfoWithArgs("the local node has become the first leader")
	}
	d.LocalNodeStatusManager.StartHeartbeat()
}

func (d *DiscoveryService) LazyInit(connectionService *ConnectionService) {
	d.ConnectionService = connectionService
	d.ConnectionService.addMemberConnectionListenerSupplier(func() connectionservice.MemberConnectionListener {
		var listener connectionservice.MemberConnectionListener
		listener = &DiscoveryMemberConnectionListener{
			discoveryService: d,
		}
		return listener
	})
}

func (d *DiscoveryService) queryMembers() []*configdiscovery.Member {
	clusterId := d.LocalNodeStatusManager.LocalMember.Key.ClusterId
	filter := eq(Member.ID_CLUSTER_ID, clusterId)
	return d.SharedConfigService.find(configdiscovery.MEMBERNAME, filter)
}

func (d *DiscoveryService) removeMemberIfInavailable(
	nodeId string,
	memberHost string,
	memberPort int) (bool, error) {

	clusterId := d.LocalNodeStatusManager.LocalMember.Key.ClusterId
	filter := option.Filter.newBuilder(3)
	if memberHost == "" {
		filter = filter.eq(ID_CLUSTER_ID, clusterId).eq(ID_NODE_ID, nodeId)
	} else {
		filter = filter.eq(a, memberHost).eq(a, memberPort)
	}
	return d.SharedConfigService.removeOne(configdiscovery.MEMBERNAME, filter)

}

func (d *DiscoveryService) listenLeaderChangeEvent() {
	go func() {
		opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
		stream, err := d.SharedConfigService.Subscribe("leader", opts)
		if err != nil {
			DISCOVERYSERVICELOGGER.FatalWithError("Error subscribing to change stream of collection:", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		d.cancelLeaderChangeRoutine = cancel
		for stream.Next(ctx) {
			var streamEvent bson.M
			if err := stream.Decode(&streamEvent); err != nil {
				DISCOVERYSERVICELOGGER.FatalWithError("Error decoding change stream event:", err)
				continue
			}
			var changedLeader *configdiscovery.Leader
			if err := stream.Decode(&changedLeader); err != nil {
				DISCOVERYSERVICELOGGER.FatalWithError("Error decoding change stream event:", err)
				continue
			}
			fullDoc, fullDocumentFound := streamEvent["fullDocument"].(bson.M)
			if !fullDocumentFound && changedLeader == nil {
				DISCOVERYSERVICELOGGER.Fatal("clusterId can not be obtained")
				continue
			}
			var clusterId string
			if changedLeader != nil {
				clusterId = changedLeader.ClusterId
			} else {
				clusterId = fullDoc["_id"].(string)
			}
			if clusterId != d.LocalNodeStatusManager.LocalMember.Key.ClusterId {
				continue
			}
			if operationType, ok := streamEvent["operationType"].(string); ok {
				switch operationType {
				case INSERT, REPLACE, UPDATE:
					if d.Leader == nil {
						str := fmt.Sprint("The leader has changed to: %s", changedLeader)
						DISCOVERYSERVICELOGGER.InfoWithArgs(str)
					} else if d.Leader.NodeId != changedLeader.NodeId || d.Leader.Generation != changedLeader.Generation {
						str := fmt.Sprint("The leader has changed from: (%s) to: %s", d.Leader, changedLeader)
						DISCOVERYSERVICELOGGER.InfoWithArgs(str)
					}
					d.Leader = changedLeader
				case DELETE:
					d.Leader = nil
					delay := int(5 * rand.Float64())
					str := fmt.Sprint("The leader has been deleted. Trying to be the first leader after %s seconds", delay)
					DISCOVERYSERVICELOGGER.InfoWithArgs(str)
					time.Sleep(time.Duration(delay) * time.Second)
					if d.Leader == nil {
						isLeader, err := d.LocalNodeStatusManager.TryBecomeFirstLeader()
						if err != nil {
							DISCOVERYSERVICELOGGER.ErrorWithMessage(
								"Caught an error while trying to become the first leader", err,
							)
						} else if isLeader {
							DISCOVERYSERVICELOGGER.InfoWithArgs("The local node has become the first leader")
						} else {
							DISCOVERYSERVICELOGGER.InfoWithArgs("Another node has become the first leader")
						}
					}
				case INVALIDATE:
					d.Leader = nil
				default:
					str := fmt.Sprint("Detected an illegal operation"+
						" on the collection leader in the change stream event: %s", streamEvent)
					DISCOVERYSERVICELOGGER.Fatal(str)
				}
			}
		}
	}()
}

func (d *DiscoveryService) listenMembersChangeEvent() {
	go func() {
		opts := options.ChangeStream().SetFullDocument(options.Default)
		stream, err := d.SharedConfigService.Subscribe("member", opts)
		if err != nil {
			DISCOVERYSERVICELOGGER.FatalWithError("Error subscribing to change stream of collection member:", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		d.cancelMemberChangeRoutine = cancel
		for stream.Next(ctx) {
			var streamEvent bson.M
			if err := stream.Decode(&streamEvent); err != nil {
				DISCOVERYSERVICELOGGER.FatalWithError("Error decoding change stream event of member:", err)
				continue
			}
			var changedMember *configdiscovery.Member
			if err := stream.Decode(&changedMember); err != nil {
				DISCOVERYSERVICELOGGER.FatalWithError("Error decoding change stream event to member:", err)
				continue
			}
			defaultDoc, defaultDocumentFound := streamEvent["documentKey"].(bson.M)
			if !defaultDocumentFound && changedMember == nil {
				DISCOVERYSERVICELOGGER.Fatal("clusterId can not be obtained")
				continue
			}
			var clusterId string
			if changedMember != nil {
				clusterId = changedMember.Key.ClusterId
			} else {
				clusterId = defaultDoc["cluster_id"].(string)
			}
			nodeId := defaultDoc["_id"].(string)
			if clusterId != d.LocalNodeStatusManager.LocalMember.Key.ClusterId {
				continue
			}

			if operationType, ok := streamEvent["operationType"].(string); ok {
				switch operationType {
				case INSERT, REPLACE:
					onMemberAddedOrReplaced(changedMember)
				case UPDATE:
					onMemberUpdated(nodeId, streamEvent[updatedescription])
				case DELETE:
					deletedMember, exists := d.AllKnownMembers.Get(nodeId)
					if exists {
						d.AllKnownMembers.Remove(nodeId)
					} else {
						continue
					}
					str := fmt.Sprint("a member has been deleted: %v", deletedMember)
					DISCOVERYSERVICELOGGER.InfoWithArgs(str)
					d.updateOtherActiveConnectedMemberList(false, deletedMember)

					if nodeId == d.LocalNodeStatusManager.LocalMember.Key.NodeId {
						if !d.LocalNodeStatusManager.IsClosing {
							d.LocalNodeStatusManager.registerLocalNodeAsMember(true)
						}
					}
				case INVALIDATE:
					for key, deletetedMember := range d.AllKnownMembers.Items() {
						str := fmt.Sprint("a member has been deleted: %v", deletetedMember)
						DISCOVERYSERVICELOGGER.InfoWithArgs(str)
						d.updateOtherActiveConnectedMemberList(false, deletetedMember)
						d.AllKnownMembers.Remove(key)
					}

				default:
					str := fmt.Sprint("detected an illegal operation on the collection \""+
						configdiscovery.MEMBERNAME+"\" in the change stream event: %v", streamEvent)
					DISCOVERYSERVICELOGGER.Fatal(str)
				}
			}
			d.updateActiveMembers(d.AllKnownMembers.Items())
			d.ConnectionService.updateHasConnectedToAllMembers(d.AllKnownMembers)
		}
	}()
}

func (d *DiscoveryService) onMemberAddedOrReplaced(newMember *configdiscovery.Member) {
	nodeId := newMember.Key.NodeId
	localMember := d.LocalNodeStatusManager.LocalMember
	isLocalNode := nodeId == localMember.Key.NodeId
	if d.AllKnownMembers.SetIfAbsent(nodeId, newMember) == true {
		str := fmt.Sprint("a new member has been added: ")
		DISCOVERYSERVICELOGGER.InfoWithArgs(str, newMember)
	}

	d.mu.Lock()
	if isLocalNode {
		d.LocalNodeStatusManager.UpdateInfo(newMember)
	}
	if newMember.Status.IsActive && d.ConnectionService.IsMemberConnected(nodeId) {
		d.updateOtherActiveConnectedMemberList(true, newMember)
		if d.notifyMembersChangeFuture.wait.Load() == true {
			d.notifyMembersChangeFuture.wait.Store(false)
		}
		d.notifyMembersChangeFuture.computeFuture(
			d.notifyMembersChangeListeners,
			d.DiscoveryProperties.DelayToNotifyMemberChangeSeconds)
	}
	d.mu.Unlock()

	shouldLocalNodeBeClient := compareMemberPriority(localMember, newMember) < 0
	if !isLocalNode && shouldLocalNodeBeClient {
		d.ConnectionService.connectMemberUntilSucceedOrRemoved(newMember)
	}
}

func (d *DiscoveryService) updateActiveMembers(allKnownMembers map[string]*configdiscovery.Member) {
	knownMembers
}

// region event

func (d *DiscoveryService) notifyMembersChangeListeners() {
	for _, listener := range d.MembersChangeListeners {
		listener.OnMembersChange()
	}
}

func (d *DiscoveryService) GetMember(nodeId string) *configdiscovery.Member {
	value, _ := d.AllKnownMembers.Get(nodeId)
	return value
}

func (d *DiscoveryService) IsKnownMember(nodeId string) bool {
	return d.AllKnownMembers.Has(nodeId)
}

func (d *DiscoveryService) isAvailableMember(knownMember *configdiscovery.Member, now time.Time) bool {
	memberHeartbeat := knownMember.Status.LastHeartbeatDate
	var t time.Time
	return memberHeartbeat != t && (now.UnixMilli()-memberHeartbeat.UnixMilli() < d.HeartbeatTimeoutMillis)
}

// end region

func (d *DiscoveryService) updateOtherActiveConnectedMemberList(isAdd bool, member *configdiscovery.Member) {
	d.mu.Lock()
	defer d.mu.Unlock()

	isLocalNode := member.IsSameNode(d.LocalNodeStatusManager.LocalMember)
	if isLocalNode {
		return
	}
	nodeType := member.NodeType
	var memberList []*configdiscovery.Member
	switch nodeType {
	case 0:
		memberList = d.OtherActiveConnectedAiServingMembers
	case 1:
		memberList = d.OtherActiveConnectedGatewayMembers
	case 2:
		memberList = d.OtherActiveConnectedServiceMembers
	default:
		memberList = make([]*configdiscovery.Member, 0)
	}
	var size int
	if isAdd {
		size = len(memberList) + 1
	} else {
		size = len(memberList)
	}
	tempOtherActiveConnectedMembers := make([]*configdiscovery.Member, size)
	copy(tempOtherActiveConnectedMembers, memberList)
	if isAdd {
		tempOtherActiveConnectedMembers = append(tempOtherActiveConnectedMembers, member)
	} else {
		tempOtherActiveConnectedMembers = collection.RemoveByValue(tempOtherActiveConnectedMembers, member)
	}
	switch nodeType {
	case 0:
		d.OtherActiveConnectedAiServingMembers = tempOtherActiveConnectedMembers
	case 1:
		d.OtherActiveConnectedGatewayMembers = tempOtherActiveConnectedMembers
	case 2:
		d.OtherActiveConnectedServiceMembers = tempOtherActiveConnectedMembers
	}
	d.OtherActiveConnectedMembers = collection.UnionThreeSlices(d.OtherActiveConnectedAiServingMembers,
		d.OtherActiveConnectedGatewayMembers, d.OtherActiveConnectedServiceMembers)
}

// region Leader

func (d *DiscoveryService) FindQualifiedMembersToBeLeader() []*configdiscovery.Member {
	members := make([]*configdiscovery.Member, len(d.ActiveSortedServiceMembers))
	highestPriority := math.MinInt
	for _, member := range d.ActiveSortedServiceMembers {
		if member.Priority < highestPriority {
			return members
		}
		if d.isQualifiedToBeLeader(member) {
			highestPriority = member.Priority
			members = append(members, member)
		}
	}
	return members
}

func (d *DiscoveryService) isQualifiedToBeLeader(member *configdiscovery.Member) bool {
	return member.NodeType == nodetype.SERVICE && member.IsLeaderEligible && member.Status.IsActive
}

// end region

// region MemberConnectionListener

// TODO: check where return value unneccesary
type DiscoveryMemberConnectionListener struct {
	discoveryService *DiscoveryService
	member           *configdiscovery.Member
}

func (d *DiscoveryMemberConnectionListener) GetName() string {
	return "DiscoveryMemberConnectionListener"
}

func (d *DiscoveryMemberConnectionListener) OnOpeningHandshakeCompleted(member *configdiscovery.Member) error {
	d.member = member
	d.discoveryService.updateOtherActiveConnectedMemberList(true, d.member)
	return nil
}

func (d *DiscoveryMemberConnectionListener) OnConnectionClosed() error {
	if d.member != nil {
		d.discoveryService.updateOtherActiveConnectedMemberList(false, d.member)
	}
	return nil
}

// not implemented
func (d *DiscoveryMemberConnectionListener) OnConnectionOpened(connection *connectionservice.GurmsConnection) error {
	return nil
}

// not implemented
func (d *DiscoveryMemberConnectionListener) OnClosingHandshakeCompleted() {
}

// not implemented
func (d *DiscoveryMemberConnectionListener) OnDataReceived(value any) error {
	return nil
}

// end region

// region future
type Future struct {
	wait atomic.Bool
}

func (f *Future) computeFuture(notifyMembersChangeListeners func(), delay int) {
	go func() {
		f.wait.Store(true)
		time.Sleep(time.Duration(delay) * time.Second)
		if f.wait.Load() == false {
			return
		}
		notifyMembersChangeListeners()
		f.wait.Store(false)
	}()
}
