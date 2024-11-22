package connectionservice

import (
	"net"
	"time"
)

type GurmsConnection struct {
	NodeId                 string
	Connection             net.Conn
	IsLocalNodeClient      bool
	LastKeepaliveTimestamp int64
	Listeners              []*MemberConnectionListener
	IsClosing              bool
}

func NewGurmsConnection(
	nodeid string,
	connection net.Conn,
	isLocalNodeClient bool,
	listeners []*MemberConnectionListener,
) *GurmsConnection {
	return &GurmsConnection{
		NodeId:                 nodeid,
		Connection:             connection,
		IsLocalNodeClient:      isLocalNodeClient,
		Listeners:              listeners,
		LastKeepaliveTimestamp: time.Now().UnixMilli(),
		IsClosing:              false,
	}
}
