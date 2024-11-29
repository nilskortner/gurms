package connectionservice

import (
	"encoding/gob"
	"net"
	"sync/atomic"
	"time"
)

type GurmsConnection struct {
	NodeId                 string
	Connection             net.Conn
	IsLocalNodeClient      bool
	LastKeepaliveTimestamp int64
	Listeners              []MemberConnectionListener
	IsClosing              bool
	IsClosed               atomic.Bool
	DataChan               chan any
	StopListenerChan       chan struct{}
	StopDecoderChan        chan struct{}
}

func NewGurmsConnection(
	nodeid string,
	connection net.Conn,
	isLocalNodeClient bool,
	listeners []MemberConnectionListener,
) *GurmsConnection {
	channel := make(chan any, 256)
	stopListener := make(chan struct{})
	stopDecoder := make(chan struct{})
	g := &GurmsConnection{
		NodeId:                 nodeid,
		Connection:             connection,
		IsLocalNodeClient:      isLocalNodeClient,
		Listeners:              listeners,
		LastKeepaliveTimestamp: time.Now().UnixMilli(),
		IsClosing:              false,
		StopListenerChan:       stopListener,
		StopDecoderChan:        stopDecoder,
		DataChan:               channel,
	}
	g.StartConnectionDecoderRoutine()
	return g
}

func (g *GurmsConnection) StartConnectionDecoderRoutine() {
	decoder := gob.NewDecoder(g.Connection)

	go func() {
		for {
			select {
			case <-g.StopDecoderChan:
				return
			default:
				var data any
				err := decoder.Decode(&data)
				if err == nil {
					g.DataChan <- data
				}
			}
		}
	}()
}
