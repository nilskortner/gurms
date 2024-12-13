package discovery

import (
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"sync/atomic"
	"time"
)

var LOCALNODESTATUSMANAGERLOGGER logger.Logger = factory.GetLogger("LocalNodeStatusManager")

type LocalNodeStatusManager struct {
	DiscoveryService        DiscoveryService
	SharedConfigService     SharedConfigService
	LocalMember             *configdiscovery.Member
	IsLocalNodeRegistered   bool
	IsClosing               bool
	HeartbeatInterval       time.Duration
	HeartbeatIntervalMillis int64
	IsHealthStatusUpdating  atomic.Bool
}

type DiscoveryService interface{}

type SharedConfigService interface{}

func NewLocalNodeStatusManager(
	discoveryService DiscoveryService,
	sharedConfigService SharedConfigService,
	localMember *configdiscovery.Member,
	heartbeatIntervalSeconds int,
) *LocalNodeStatusManager {
	heartbeatInterval := time.Duration(heartbeatIntervalSeconds) * time.Second
	heartbeatIntervalMillis := int64(heartbeatInterval.Milliseconds())
	return &LocalNodeStatusManager{
		DiscoveryService:        discoveryService,
		SharedConfigService:     sharedConfigService,
		LocalMember:             localMember,
		HeartbeatInterval:       heartbeatInterval,
		HeartbeatIntervalMillis: heartbeatIntervalMillis,
	}
}
