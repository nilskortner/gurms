package memberconnectionlistener

import (
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/infra/cluster/service/connectionservice"
)

type MemberConnectionListener interface {
	OnConnectionOpened(connection *connectionservice.GurmsConnection)
	OnConnectionClosed()
	OnOpeningHandshakeCompleted(member *configdiscovery.Member)
}
