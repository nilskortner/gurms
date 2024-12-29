package configdiscovery

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"time"
)

const MEMBERNAME = "member"

type Member struct {
	Key              Key                   `bson:",inline"`
	Zone             string                `bson:"zone"`
	Name             string                `bson:"name"`
	NodeType         nodetype.NodeType     `bson:"nodetype"`
	Nodeversion      *nodetype.NodeVersion `bson:",inline"`
	IsSeed           bool                  `bson:"is_seed"`
	IsLeaderEligible bool                  `bson:"is_leader_eligible"`
	RegistrationDate time.Time             `bson:"registration_date"`
	Priority         int                   `bson:"priority"`
	MemberHost       string                `bson:"memberhost"`
	MemberPort       int                   `bson:"member_port"`
	AdminApiAddress  string                `bson:"admin_api_address"`
	WsAddress        string                `bson:"ws_address"`
	TcpAddress       string                `bson:"tcp_address"`
	UdpAddress       string                `bson:"udp_address"`
	Status           *MemberStatus         `bson:",inline"`
}

type Key struct {
	ClusterId string `bson:"cluster_id"`
	NodeId    string `bson:"_id"`
}

type MemberStatus struct {
	HasJoinedCluster  bool      `bson:"has_joined_cluster"`
	IsHealthy         bool      `bson:"is_healthy"`
	IsActive          bool      `bson:"is_active"`
	LastHeartbeatDate time.Time `bson:"last_heartbeat_date"`
}

func NewMember(
	clusterId string,
	nodeId string,
	zone string,
	name string,
	nodeType nodetype.NodeType,
	nodeversion *nodetype.NodeVersion,
	isSeed bool,
	isLeaderEligible bool,
	registrationDate time.Time,
	priority int,
	memberHost string,
	memberPort int,
	adminApiAddress string,
	wsAddress string,
	tcpAddress string,
	udpAddress string,
	hasJoinedCluster bool,
	isActive bool,
	isHealthy bool,
) *Member {
	key := NewKey(clusterId, nodeId)
	date := time.Now()
	memberStatus := NewMemberStatus(hasJoinedCluster, isHealthy, isActive, date)
	return &Member{
		Key:              key,
		Zone:             zone,
		Name:             name,
		NodeType:         nodeType,
		Nodeversion:      nodeversion,
		IsSeed:           isSeed,
		IsLeaderEligible: isLeaderEligible,
		RegistrationDate: registrationDate,
		Priority:         priority,
		MemberHost:       memberHost,
		MemberPort:       memberPort,
		AdminApiAddress:  adminApiAddress,
		WsAddress:        wsAddress,
		TcpAddress:       tcpAddress,
		UdpAddress:       udpAddress,
		Status:           memberStatus,
	}
}

func (m *Member) IsSameNode(member *Member) bool {
	return m.IsSameId(member) || m.IsSameAddress(member)
}

func (m *Member) IsSameId(member *Member) bool {
	return m.Key == member.Key
}

func (m *Member) IsSameAddress(member *Member) bool {
	return m.MemberHost == member.MemberHost && m.MemberPort == member.MemberPort
}

// TODO

func NewKey(clusterId string, nodeId string) Key {
	return Key{
		ClusterId: clusterId,
		NodeId:    nodeId,
	}
}

func NewMemberStatus(hasJoinedCluster bool,
	isHealthy bool,
	isActive bool,
	lastHearthbeatDate time.Time) *MemberStatus {

	return &MemberStatus{
		hasJoinedCluster:  hasJoinedCluster,
		isHealthy:         isHealthy,
		isActive:          isActive,
		lastHeartbeatDate: lastHearthbeatDate,
	}
}
