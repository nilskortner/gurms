package node

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/cluster"
	"regexp"
	"strconv"
)

var nodeLogger logger.Logger = factory.GetLogger("Node")

var nodeId string
var nodeType nodetype.NodeType

func InitNodeId(id string) string {
	if nodeId != "" {
		return nodeId
	}
	if id == "" {

		nodeLogger.WarnWithArgs(
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
