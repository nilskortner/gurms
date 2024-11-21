package connectionservice

type ConnectionChannels struct {
	DataChan  chan any
	CloseChan chan struct{}
}

func NewConnectionChannels() *ConnectionChannels {
	return &ConnectionChannels{
		DataChan:  make(chan any, 256),
		CloseChan: make(chan struct{}),
	}
}

func (c *ConnectionChannels) Close() {
	close(c.DataChan)
}
