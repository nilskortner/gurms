package connectionservice

import (
	"context"
	"net"
	"time"
)

type GurmsConnection struct {
	NodeId                 string
	Connection             net.Conn
	IsLocalNodeClient      bool
	LastKeepaliveTimestamp int64
	Listeners              []MemberConnectionListener
	IsClosing              bool
	Cancel                 context.CancelFunc
	CloseContext           context.Context
	DataChan               chan any
}

func NewGurmsConnection(
	nodeid string,
	connection net.Conn,
	isLocalNodeClient bool,
	listeners []MemberConnectionListener,
) *GurmsConnection {
	ctx, cancel := context.WithCancel(context.Background())
	channel := make(chan any, 256)
	return &GurmsConnection{
		NodeId:                 nodeid,
		Connection:             connection,
		IsLocalNodeClient:      isLocalNodeClient,
		Listeners:              listeners,
		LastKeepaliveTimestamp: time.Now().UnixMilli(),
		IsClosing:              false,
		Cancel:                 cancel,
		CloseContext:           ctx,
		DataChan:               channel,
	}
}

func (g *GurmsConnection) IsClosed() bool {
	select {
	case <-g.CloseContext.Done():
		return true
	default:
		return false
	}
}
