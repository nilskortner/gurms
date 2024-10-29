package property

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var GPMLOGGER logger.Logger = factory.GetLogger("GurmsPropertiesManager")

type Node interface {
	Start()
}

type GurmsPropertiesManager struct {
	node                 Node
	localGurmsProperties *GurmsProperties
}

func NewGurmsPropertiesManager(node Node, localGurmsProperties *GurmsProperties) *GurmsPropertiesManager {
	return &GurmsPropertiesManager{
		node:                 node,
		localGurmsProperties: localGurmsProperties,
	}
}
