package configdiscovery

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"time"
)

type Member struct {
	Key              Key               `bson:",inline"`
	Zone             string            `bson:"zone"`
	Name             string            `bson:"name"`
	NodeType         nodetype.NodeType `bson:"nodetype"`
	Nodeversion      *Nodeversion      `bson:",inline"`
	IsSeed           bool              `bson:"is_seed"`
	IsLeaderEligible bool              `bson:"is_leader_eligible"`
	RegistrationDate time.Time         `bson:"registration_date"`
	Priority         int               `bson:"priority"`
	MemberHost       string            `bson:"memberhost"`
	MemberPort       int               `bson:"member_port"`
	AdminApiAddress  string            `bson:"admin_api_address"`
	WsAddress        string            `bson:"ws_address"`
	TcpAddress       string            `bson:"tcp_address"`
	UdpAddress       string            `bson:"udp_address"`
	Status           MemberStatus      `bson:",inline"`
}

type Key struct {
	ClusterId string `bson:"cluster_id"`
	NodeId    string `bson:"_id"`
}

type MemberStatus struct {
	hasJoinedCluster  bool      `bson:"has_joined_cluster"`
	isHealthy         bool      `bson:"is_healthy"`
	isActive          bool      `bson:"is_active"`
	lastHeartbeatDate time.Time `bson:"last_heartbeat_date"`
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
