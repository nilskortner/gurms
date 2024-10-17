package node

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/cluster"
)

var nodeLogger logger.Logger = factory.GetLogger("Node")

var nodeId string
var nodeType nodetype.NodeType

func initNodeId(id string) string {
	if nodeId != "" {
		return nodeId
	}
	if id == "" {

		nodeLogger.WarnWithArgs(
			"A random node ID ({}) has been used. You should better set a node ID manually in production",
			id)
	} else {
		if len(id) > cluster.NODE_ID_MAX_LENGTH {

		}
		if !id.matches() {

		}
	}
	nodeId = id
	return nodeId
}
