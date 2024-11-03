package connectionservice

import "gurms/internal/infra/cluster/service/discovery"

type GurmsConnection struct {
	nodeId                 string
	connection             ChannelOperations
	isClosing              bool
	isLocalNodeClient      bool
	lastKeepaliveTimestamp int64
	listeners              []*discovery.MemberConnectionListener
}

func NewGurmsConnection(
	nodeid string,
	connection,
	isLocalNodeClient bool,
	listeners []*discovery.MemberConnectionListener,
) *GurmsConnection {
	return &GurmsConnection{
		nodeId:            nodeid,
		connection:        connection,
		isLocalNodeClient: isLocalNodeClient,
	}
}
