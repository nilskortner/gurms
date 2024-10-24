package cluster

type ClusterProperties struct {
	id          string
	node        *NodeProperties
	connection  *CpnnectionsProperties
	discovery   *DiscoveryProperties
	shardConfig *SharedConfigProperties
	rpc         *RpcProperties
}
