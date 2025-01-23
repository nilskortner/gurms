package service

import (
	"gurms/internal/infra/cluster/service/idgen"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
)

var ISERVICELOGGER logger.Logger = factory.GetLogger("IdService")

const FLAKE_ID_GENERATORS_LENGTH = int(idgen.NumServiceTypes)

type IdService struct {
	discoveryService          *DiscoveryService
	idGenerators              []*idgen.SnowflakeIdGenerator
	previousLocalDataCenterId int
	previousLocalWorkerId     int
}

func NewIdService(discoveryService *DiscoveryService) *IdService {
	idService := &IdService{
		discoveryService: discoveryService,
	}

	idGenerators := make([]*idgen.SnowflakeIdGenerator, FLAKE_ID_GENERATORS_LENGTH)
	for i := 0; i < FLAKE_ID_GENERATORS_LENGTH; i++ {
		generator, err := idgen.NewSnowflakeIdGenerator(0, 0)
		if err != nil {
			ISERVICELOGGER.FatalWithError("could not initialize SnowflakeIdGenerator", err)
			continue
		}
		idGenerators[i] = generator
	}
	idService.idGenerators = idGenerators
	discoveryService.addOnMembersChangeListener(func() {
		dataCenterId := idService.findNewDataCenterId()
	})

	return idService
}

func (i *IdService) nextIncreasingId() int64 {}

func (i *IdService) nextLargeGapId() int64 {}

func (i *IdService) findNewDataCenterId() int {}

func (i *IdService) findNewWorkerID() int {}
