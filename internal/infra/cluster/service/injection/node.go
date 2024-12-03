package injection

type Node interface {
	DoKeepalive()
	GetNodeId()
	OpeningHandshakeRequestCall() any
	KeepAliveRequestCall()
}
