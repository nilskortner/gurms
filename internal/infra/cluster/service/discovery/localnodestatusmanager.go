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
	LocalMember             *configdiscovery.Member
	IsLocalNodeRegistered   bool
	IsClosing               bool
	HeartbeatInterval       time.Time
	HeartbeatIntervalMIllis int64
	IsHealthStatusUpdating  atomic.Bool
}

//TODO
