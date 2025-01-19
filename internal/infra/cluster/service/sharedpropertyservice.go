package service

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var SHAREDPROPERTYLOGGER logger.Logger = factory.GetLogger("SharedProperty")

type SharedPropertyService struct {
	clusterId string
	nodeType  nodetype.NodeType

	propertiesManager       *GurmsPropertiesManager
	sharedClusterProperties *sharedClusterProperties
	sharedConfigService     *SharedConfigService

	propertiesChangeListeners []func(*GurmsProperties)
}

func NewSharedPropertyService(clusterId string, nodeType nodetype.NodeType,
	propertiesManager *GurmsPropertiesManager) *SharedPropertyService {

	&NewSharedPropertyService{
		clusterId:                 clusterId,
		nodeType:                  nodeType,
		propertiesManager:         propertiesManager,
		propertiesChangeListeners: make([]func(*GurmsProperties), 0),
	}
}

func (s *SharedPropertyService) LazyInit(sharedConfigService *SharedConfigService) {
	s.sharedConfigService = sharedConfigService
}

func (s *SharedPropertyService) Start() {
	go func() {
		opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
		stream, err := s.sharedConfigService.Subscribe("sharedclusterproperties", opts)
		if err != nil {
			SHAREDPROPERTYLOGGER.FatalWithError("error subscribing to changestream", err)
		}

		stream.Next()
	}()
	err := s.initializeSharedProperties()
	if err != nil {
		SHAREDPROPERTYLOGGER.FatalWithError("failed to initialize the shared properties", err)
	}
}
