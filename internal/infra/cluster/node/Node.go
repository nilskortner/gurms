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

	SharedConfigService   *SharedConfigService
	SharedPropertyService *SharedPropertyService
	CodecService          *codec.CodecService
	ConnectionService     *service.ConnectionService
	DiscoveryService      *DiscoveryService
	RpcService            *rpcserv.RpcService
	IdService             *IdService
}

func NewNode(
	nType nodetype.NodeType,
	propertiesManager *property.GurmsPropertiesManager,
	baseServiceAddresssManager *address.BaseServiceAddressManager,
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
	codecService := codec.NewCodecService()
	connectionService := connectionservice.NewConnectionService()
	rpcService := rpcserv.NewRpcService(nodeType, rpcProperties)

	return &Node{
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
	a
}

// for interface Node

func GetNodeID()
