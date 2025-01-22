package service

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var ISERVICELOGGER logger.Logger = factory.GetLogger("IdService")

const FLAKE_ID_GENERATORS_LENGTH = getSharedEnumConstants()

type IdService struct {
	discoveryService          *DiscoveryService
	idGenerators              []*idgen.SnowflakeIdGenerator
	previousLocalDataCenterId int
	previousLocalWorkerId     int
}

func NewIdService(discoveryService *DiscoveryService) *IdService {

	idGenerators := make()
	for i := 0; i < FLAKE_ID_GENERATORS_LENGTH; i++ {
		id
	}

	return &IdService{
		discoveryService: discoveryService,
	}
}
