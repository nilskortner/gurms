package injection

import "gurms/internal/infra/cluster/service/connectionservice"

type Node interface {
	OpeningHandshakeRequestCall(*connectionservice.GurmsConnection) any
	KeepAliveRequestCall()
}
