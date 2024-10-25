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

func NewClusterProperties(id string) *ClusterProperties {
	return &ClusterProperties{
		id:          id,
		node:        InitNodeProperties(),
		connection:  connection.InitConnectionProperties(),
		discovery:   NewDiscoveryProperties(),
		shardConfig: NewSharedConfigProperties(),
		grpc:        NewGrpcProperties(),
	}
}
