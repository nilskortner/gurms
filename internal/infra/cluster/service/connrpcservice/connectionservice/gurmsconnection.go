package connectionservice

import (
	"gurms/internal/infra/cluster/service/discovery"
	"time"

	grpc "google.golang.org/grpc"
)

type GurmsConnection struct {
	NodeId                 string
	Connection             *grpc.ClientConn
	IsLocalNodeClient      bool
	LastKeepaliveTimestamp int64
	Listeners              []*discovery.MemberConnectionListener
	IsClosing              bool
}

func NewGurmsConnection(
	nodeid string,
	connection *grpc.ClientConn,
	isLocalNodeClient bool,
	listeners []*discovery.MemberConnectionListener,
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
