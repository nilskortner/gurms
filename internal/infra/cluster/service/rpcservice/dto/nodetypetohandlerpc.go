package dto

type NodeTypeToHandleRpc int

const (
	BOTH NodeTypeToHandleRpc = iota
	GATEWAY
	SERVICE
)

func (n NodeTypeToHandleRpc) String() string {
	switch n {
	case BOTH:
		return "BOTH"
	case GATEWAY:
		return "GATEWAY"
	case SERVICE:
		return "SERVICE"
	default:
		return "UNKNOWN"
	}
}
