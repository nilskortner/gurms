package connectionservice

import (
	"time"

	grpc "google.golang.org/grpc"
)

type GurmsConnection struct {
	NodeId                 string
	Connection             *ConnectionChannels
	IsLocalNodeClient      bool
	LastKeepaliveTimestamp int64
	Listeners              []MemberConnectionListener
	IsClosing              bool
}

func NewGurmsConnection(
	nodeid string,
	connection *grpc.ClientConn,
	isLocalNodeClient bool,
	listeners []MemberConnectionListener,
) *GurmsConnection {
	return &GurmsConnection{
		NodeId:                 nodeid,
		Connection:             NewConnectionChannels(),
		IsLocalNodeClient:      isLocalNodeClient,
		Listeners:              listeners,
		LastKeepaliveTimestamp: time.Now().UnixMilli(),
		IsClosing:              false,
	}
}
