package connectionservice

import (
	"gurms/internal/infra/cluster/service/discovery"
	"gurms/internal/infra/cluster/service/rpcserv"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/cluster/connection"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var CONNECTIONLOGGER logger.Logger = factory.GetLogger("ConnectionService")

var nodeIdToConnection cmap.ConcurrentMap[string, *GurmsConnection] = cmap.New[*GurmsConnection]()

var nodeIdToConnectionRetries cmap.ConcurrentMap[string, int] = cmap.New[int]()
var connectingMembers cmap.ConcurrentMap[string, struct{}] = cmap.New[struct{}]()

type ConnectionService struct {
	clientSsl               *common.SslProperties
	keepaliveIntervalMillis int64
	keepaliveTimeoutMillis  int64
	reconnectInterval       int64

	nodeIdToConnection        cmap.ConcurrentMap[string, *GurmsConnection]
	nodeIdToConnectionRetries cmap.ConcurrentMap[string, int]
	connectingMembers         cmap.ConcurrentMap[string, struct{}]

	discoveryService         *discovery.DiscoveryService
	rpcService               *rpcserv.RpcService
	hasConnectedToAllMembers bool
	serverProperties         *connection.ConnectionServerProperties
	server                   *ConnectionServer
}

func NewConnectionService() *ConnectionService {

	return &ConnectionService{}
}

func (c *ConnectionService) LazyInit() {
	c
}
