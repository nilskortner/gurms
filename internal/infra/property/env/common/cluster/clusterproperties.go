package cluster

import "gurms/internal/infra/property/env/common/cluster/connection"

type ClusterProperties struct {
	Id           string                           `bson:"id"`
	Node         *NodeProperties                  `bson:",inline"`
	Connection   *connection.ConnectionProperties `bson:",inline"`
	Discovery    *DiscoveryProperties             `bson:",inline"`
	SharedConfig *SharedConfigProperties          `bson:",inline"`
	Rpc          *RpcProperties                   `bson:",inline"`
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
