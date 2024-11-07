package connectionservice

import (
	"gurms/internal/infra/cluster/service/discovery"
	"time"

	grpc "google.golang.org/grpc"
)

type GurmsConnection struct {
	nodeId                 string
	connection             *grpc.ClientConn
	isLocalNodeClient      bool
	lastKeepaliveTimestamp int64
	listeners              []*discovery.MemberConnectionListener
	isClosing              bool
}

func NewGurmsConnection(
	nodeid string,
	connection *grpc.ClientConn,
	isLocalNodeClient bool,
	listeners []*discovery.MemberConnectionListener,
) *GurmsConnection {
	return &GurmsConnection{
		nodeId:                 nodeid,
		connection:             connection,
		isLocalNodeClient:      isLocalNodeClient,
		listeners:              listeners,
		lastKeepaliveTimestamp: time.Now().UnixMilli(),
		isClosing:              false,
	}
}
