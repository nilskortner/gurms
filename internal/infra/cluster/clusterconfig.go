package cluster

import (
	"gurms/internal/infra/address"
	"gurms/internal/infra/application"
	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/healthcheck"
	"gurms/internal/infra/property"
)

func NewNode(
	nodeType nodetype.NodeType,
	shutDown *application.ShutDownManager,
	propertiesManager *property.GurmsPropertiesManager,
	serviceAddressManager *address.BaseServiceAddressManager,
	healthCheckManager *healthcheck.HealthCheckManager) *node.Node {
	node := node.NewNode(
		nodeType,
		shutDown,
		propertiesManager,
		serviceAddressManager,
		healthCheckManager)
	node.Start()
	return node
}
