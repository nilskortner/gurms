package node

import (
	"gurms/internal/infra/address"
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/cluster/service"
	"gurms/internal/infra/cluster/service/connectionservice"
	"gurms/internal/infra/healthcheck"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/env/common/cluster"
	"gurms/internal/infra/property/env/common/cluster/connection"
	"regexp"
	"strconv"
)

var NODELOGGER logger.Logger = factory.GetLogger("Node")

var nodeId string
var nodeType nodetype.NodeType

type Node struct {
	gurmsProperties        *property.GurmsProperties
	clusterProperties      *cluster.ClusterProperties
	sharedConfigProperties *cluster.SharedConfigProperties
	nodeProperties         *cluster.NodeProperties
	connectionProperties   *connection.ConnectionProperties
	discoveryProperties    *cluster.DiscoveryProperties
	rpcProperties          *cluster.RpcProperties

	clusterId string
	zone      string
	name      string

	SharedConfigService   *service.SharedConfigService
	SharedPropertyService *service.SharedPropertyService
	CodecService          *service.CodecService
	ConnectionService     *service.ConnectionService
	DiscoveryService      *service.DiscoveryService
	RpcService            *service.RpcService
	IdService             *service.IdService
}

func NewNode(
	nType nodetype.NodeType,
	propertiesManager *property.GurmsPropertiesManager,
	baseServiceAddresssManager address.ServiceAddressManager,
	healthCheckManager *healthcheck.HealthCheckManager,
) *Node {
	nodeType = nType
	properties := propertiesManager.LocalGurmsProperties
	clusterProperties := properties.Cluster
	sharedConfigProperties := clusterProperties.SharedConfig
	nodeProperties := clusterProperties.Node
	connectionProperties := clusterProperties.Connection
	discoveryProperties := clusterProperties.Discovery
	rpcProperties := clusterProperties.Rpc

	version := gurmsContext.BuildProperties.Version
	nodeVersion := nodetype.Parse(version)
	NODELOGGER.InfoWithArgs("the local node version is: ", version)

	clusterId := clusterProperties.Id
	nodeId = InitNodeId(nodeProperties.Id)
	zone := nodeProperties.Zone
	name := nodeProperties.Name
	if name == "" {
		name = nodeId
	} else {
		if len(name) > cluster.NODE_NAME_MAX_LENGTH {
			panic("length of node id must be less than or equal to " + strconv.Itoa(cluster.NODE_NAME_MAX_LENGTH))
		}
		matched, err := regexp.MatchString("^[a-zA-Z_]\\w*$", name)
		if !matched {
			panic("The node ID must start with a letter or underscore, " +
				"and match zero or more of characters [a-zA-Z0-9_] after the beginning" +
				err.Error())
		}
	}

	node := &Node{
		gurmsProperties:        properties,
		clusterProperties:      clusterProperties,
		sharedConfigProperties: sharedConfigProperties,
		nodeProperties:         nodeProperties,
		connectionProperties:   connectionProperties,
		discoveryProperties:    discoveryProperties,
		rpcProperties:          rpcProperties,

		clusterId: clusterId,
		zone:      zone,
		name:      name,
	}

	node.CodecService = service.NewCodecService()
	node.ConnectionService = service.NewConnectionService(connectionProperties, node)
	node.RpcService = service.NewRpcService(nodeType, rpcProperties)
	node.SharedConfigService = service.NewSharedConfigService(sharedConfigProperties.Mongo)
	node.DiscoveryService = service.NewDiscoveryService(
		clusterId, nodeId, zone, name, nodeType, nodeVersion,
		nodeType == nodetype.SERVICE && nodeProperties.LeaderEligible,
		nodeProperties.Priority, nodeProperties.ActiveByDefault,
		healthCheckManager.isHealthy(),
		node.ConnectionService,
		discoveryProperties,
		baseServiceAddresssManager,
		node.SharedConfigService,
	)
	node.SharedPropertyService = service.NewSharedPropertyService(clusterId, nodeType, propertiesManager)
	node.IdService = service.NewIdService(node.DiscoveryService)

	// lazy init
	node.ConnectionService.LazyInit(node.DiscoveryService, node.RpcService)
	node.DiscoveryService.LazyInit(node.ConnectionService)
	node.RpcService.LazyInit(node.CodecService, node.ConnectionService, node.DiscoveryService)

	return node
}

func InitNodeId(id string) string {
	if nodeId != "" {
		return nodeId
	}
	if id == "" {
		id = randomstringutils.randomAlphabetic(8).toLowerCase()
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
	n.SharedPropertyService.Start()
	n.DiscoveryService.Start()
}

// for Node Injection

func (n *Node) OpeningHandshakeRequestCall(connection *connectionservice.GurmsConnection) any {
	return n.ConnectionService.HandleHandshakeRequest(connection, nodeId)
}

func (n *Node) KeepAliveRequestCall() {
	n.ConnectionService.Keepalive(nodeId)
}

func (n *Node) UpdateHealthStatus(isHealthy bool) {
	n.DiscoveryService.LocalNodeStatusManager.UpdateHealthStatus(isHealthy)
}
