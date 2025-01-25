package address

type NodeAddressInfo struct {
	MemberHost      string
	AdminApiAddress string
	WsAddress       string
	TcpAddress      string
	UdpAddress      string
}

func NewAddressInfo(memberHost, adminApiAddress, wsAddress, tcpAddress, udpAddress string) *NodeAddressInfo {
	return &NodeAddressInfo{
		MemberHost:      memberHost,
		AdminApiAddress: adminApiAddress,
		WsAddress:       wsAddress,
		TcpAddress:      tcpAddress,
		UdpAddress:      udpAddress,
	}
}
