package netutil

import "net"

const (
	IPV4_BYTE_LENGTH = 4
	IPV6_BYTE_LENGTH = 16
)

func IsIp(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}
