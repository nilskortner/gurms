package cluster

import "gurms/internal/infra/property/env/common/cluster/connection"

type ClusterProperties struct {
	Id           string
	Node         *NodeProperties
	Connection   *connection.ConnectionProperties
	Discovery    *DiscoveryProperties
	SharedConfig *SharedConfigProperties
	Rpc          *RpcProperties
}

func NewClusterProperties() *ClusterProperties {
	return &ClusterProperties{
		Id:           "gurms",
		Node:         InitNodeProperties(),
		Connection:   connection.InitConnectionProperties(),
		Discovery:    NewDiscoveryProperties(),
		SharedConfig: NewSharedConfigProperties(),
		Rpc:          NewRpcProperties(),
	}
}
