package node

import (
	"gurms/internal/infra/address"
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/healthcheck"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/env/common/cluster"
	"regexp"
	"strconv"
)

var NODELOGGER logger.Logger = factory.GetLogger("Node")

var nodeId string
var nodeType nodetype.NodeType

type Node struct {
	sharedConfigService   *SharedConfigService
	sharedPropertyService *SharedPropertyService
	codecService          *CodecService
	connectionService     *ConnectionService
	discoveryService      *DiscoveryService
	grpcService           *GrocService
	idService             *IdService
}

func NewNode(
	nType nodetype.NodeType,
	propertiesManager *property.GurmsPropertiesManager,
	baseServiceAddresssManager *address.BaseServiceAddressManager,
	healthCheckManager *healthcheck.HealthCheckManager,
) *Node {
	nodeType = nType
}

func InitNodeId(id string) string {
	if nodeId != "" {
		return nodeId
	}
	if id == "" {

		NODELOGGER.WarnWithArgs(
			"A random node ID ({}) has been used. You should better set a node ID manually in production",
			id)
	} else {
		if len(id) > cluster.NODE_ID_MAX_LENGTH {
			panic("length of node id must be less than or equal to " + strconv.Itoa(cluster.NODE_ID_MAX_LENGTH))
		}
		matched, err := regexp.MatchString("^[a-zA-Z_]\\w*$", id)
		if !matched {
			panic("The node ID must start with a letter or underscore, " +
				"and match zero or more of characters [a-zA-Z0-9_] after the beginning" +
				err.Error())
		}
	}
	nodeId = id
	return nodeId
}

func (n *Node) Start() {

}
