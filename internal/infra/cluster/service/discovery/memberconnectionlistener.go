package discovery

import (
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
)

type MemberConnectionListener interface {
	OnConnectionClosed()
	OnOpeningHandshakeCompleted(member *configdiscovery.Member)
}
