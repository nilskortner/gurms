package request

import "gurms/internal/infra/cluster/node"

const NAME = "keepalive"

type KeepaliveRequest struct {
	node node.Node
}

func NewKeepAliveRequest() *KeepaliveRequest {
	return &KeepaliveRequest{}
}
