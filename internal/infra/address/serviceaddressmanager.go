package address

type ServiceAddressManager interface {
	GetMemberHost() string
	GetAdminApiAddress() string
	GetWsAddress() string
	GetTcpAddress() string
	GetUdpAddress() string
}
