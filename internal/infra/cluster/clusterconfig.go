package cluster

import (
	"gurms/internal/infra/address"
	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/healthcheck"
	"gurms/internal/infra/property"
)

func NewNode(
	nodeType nodetype.NodeType,
	propertiesManager *property.GurmsPropertiesManager,
	serviceAddressManager *address.BaseServiceAddressManager,
	healthCheckManager *healthcheck.HealthCheckManager) *node.Node {
	node := node.NewNode(
		nodeType,
		propertiesManager,
		serviceAddressManager,
		healthCheckManager)
	node.Start()
	return node
}
