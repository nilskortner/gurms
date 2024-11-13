package service

import "gurms/internal/infra/cluster/service/connrpcservice"

type DiscoveryService struct {
	connectionService *ConnectionService
}

func (d *DiscoveryService) LazyInit(connectionService *ConnectionService) {
	d.connectionService. = append(, func() memberconnectionlistener.MemberConnectionListener{

	})
}
