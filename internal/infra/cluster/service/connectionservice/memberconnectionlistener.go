package connectionservice

import (
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
)

type MemberConnectionListener interface {
	OnConnectionOpened(connection *GurmsConnection) error
	OnConnectionClosed()
	OnOpeningHandshakeCompleted(member *configdiscovery.Member)
	OnClosingHandshakeCompleted()
	OnDataReceived(value any)
}
