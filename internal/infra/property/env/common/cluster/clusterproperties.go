package cluster

import "gurms/internal/infra/property/env/common/cluster/connection"

type ClusterProperties struct {
	id          string
	node        *NodeProperties
	connection  *connection.ConnectionProperties
	discovery   *DiscoveryProperties
	shardConfig *SharedConfigProperties
	grpc        *GrpcProperties
}

func NewClusterProperties() *ClusterProperties {
	return &ClusterProperties{
		id:          "gurms",
		node:        InitNodeProperties(),
		connection:  connection.InitConnectionProperties(),
		discovery:   NewDiscoveryProperties(),
		shardConfig: NewSharedConfigProperties(),
		grpc:        NewGrpcProperties(),
	}
}
