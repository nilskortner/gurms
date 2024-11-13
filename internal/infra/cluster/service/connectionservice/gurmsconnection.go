package connectionservice

import (
	"time"

	grpc "google.golang.org/grpc"
)

type MemberConnectionListener interface {
	OnConnectionClosed()
}

type GurmsConnection struct {
	NodeId                 string
	Connection             *grpc.ClientConn
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
		Connection:             connection,
		IsLocalNodeClient:      isLocalNodeClient,
		Listeners:              listeners,
		LastKeepaliveTimestamp: time.Now().UnixMilli(),
		IsClosing:              false,
	}
}
