package service

import (
	"fmt"
	"gurms/internal/infra/cluster/service/idgen"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/supportpkgs/datastructures/treeset"
)

var IDSERVICELOGGER logger.Logger = factory.GetLogger("IdService")

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
			IDSERVICELOGGER.FatalWithError("could not initialize SnowflakeIdGenerator", err)
			continue
		}
		idGenerators[i] = generator
	}
	idService.idGenerators = idGenerators

	discoveryService.addOnMembersChangeListener(func() {
		dataCenterId := idService.findNewDataCenterId()
		localWorkerId := idService.findNewWorkerID()
		if idService.previousLocalDataCenterId != dataCenterId ||
			localWorkerId != idService.previousLocalWorkerId {
			for _, idGenerator := range idService.idGenerators {
				idGenerator.UpdateNodeInfo(dataCenterId, localWorkerId)
			}
			idService.previousLocalDataCenterId = dataCenterId
			idService.previousLocalWorkerId = localWorkerId
		}
	})

	return idService
}

func (i *IdService) nextIncreasingId(servicetype idgen.ServiceType) int64 {
	return i.idGenerators[servicetype].NextIncreasingId()
}

func (i *IdService) nextLargeGapId(servicetype idgen.ServiceType) int64 {
	return i.idGenerators[servicetype].NextLargeGapId()
}

func (i *IdService) findNewDataCenterId() int {
	zones := treeset.New[string](treeset.StringComparator)
	for _, member := range i.discoveryService.AllKnownMembers.Items() {
		zones.Put(member.Zone)
	}
	dataCenterId := zones.HeadSetSize(i.discoveryService.LocalNodeStatusManager.LocalMember.Zone)
	if dataCenterId >= idgen.MAX_DATA_CENTER_ID {
		fallbackDataCenterId := dataCenterId % idgen.MAX_DATA_CENTER_ID
		str := fmt.Sprintf("the data center ID %d is larger than %d, so the ID falls back to %d."+
			" it runs the risk of generating same IDs in the cluster",
			dataCenterId, idgen.MAX_DATA_CENTER_ID, fallbackDataCenterId)
		IDSERVICELOGGER.Warn(str)
		dataCenterId = fallbackDataCenterId
	}
	return dataCenterId
}

func (i *IdService) findNewWorkerID() int {
	localWorkerId := i.discoveryService.getLocalServiceMemberIndex()
	if localWorkerId == -1 {
		return i.previousLocalWorkerId
	}
	if localWorkerId >= idgen.MAX_WORKER_ID {
		fallbackWorkerId := localWorkerId % idgen.MAX_WORKER_ID
		str := fmt.Sprintf("the data center ID %d is larger than %d, so the ID falls back to %d."+
			" it runs the risk of generating same IDs in the cluster",
			localWorkerId, idgen.MAX_WORKER_ID, fallbackWorkerId)
		IDSERVICELOGGER.Warn(str)
		localWorkerId = fallbackWorkerId
	}
	return localWorkerId

}
