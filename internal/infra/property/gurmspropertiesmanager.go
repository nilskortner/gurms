package property

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var GPMLOGGER logger.Logger = factory.GetLogger("GurmsPropertiesManager")

var DEFAULT_PROPERTIES = NewGurmsProperties()
var DEFAULT_PROPERTIES_JSON_NODE = ""

type Node interface {
	Start()
	InitNodeId(id string) string
}

type GurmsPropertiesManager struct {
	LocalPropertiesChangeListeners []func(properties *GurmsProperties)
	LastestConfigFilePath          string
	Node                           *Node
	LocalGurmsProperties           *GurmsProperties
}

func NewGurmsPropertiesManager(node *Node, localGurmsProperties *GurmsProperties) *GurmsPropertiesManager {
	return &GurmsPropertiesManager{
		Node:                 node,
		LocalGurmsProperties: localGurmsProperties,
	}
}

func (g *GurmsPropertiesManager) SetNode(node *Node) {
	g.Node = node
}

// Listener

func (g *GurmsPropertiesManager) AddLocalPropertiesChangeListener(listener func(properties *GurmsProperties)) {
	g.LocalPropertiesChangeListeners = append(g.LocalPropertiesChangeListeners, listener)
}
