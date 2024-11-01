package property

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var GPMLOGGER logger.Logger = factory.GetLogger("GurmsPropertiesManager")

type Node interface {
	Start()
	InitNodeId(id string) string
}

type GurmsPropertiesManager struct {
	Node                 *Node
	LocalGurmsProperties *GurmsProperties
}

func NewGurmsPropertiesManager(localGurmsProperties *GurmsProperties) *GurmsPropertiesManager {
	return &GurmsPropertiesManager{
		LocalGurmsProperties: localGurmsProperties,
	}
}

func (g *GurmsPropertiesManager) SetNode(node *Node) {
	g.Node = node
}
