package configdiscovery

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"time"
)

type Member struct {
	Key              Key
	Zone             string
	Zame             string
	NodeType         nodetype.NodeType
	Nodeversion      *Nodeversion
	IsSeed           bool
	IsLeaderEligible bool
	RegistrationDate time.Time
	Priority         int
	MemberHost       string
	MemberPort       int
	AdminApiAddress  string
	WsAddress        string
	TcpAddress       string
	UdpAddress       string
	Status           MemberStatus
}

type Key struct {
	ClusterId string
	NodeId    string
}

type MemberStatus struct {
	hasJoinedCluster  bool
	isHealthy         bool
	isActive          bool
	lastHeartbeatDate time.Time
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
