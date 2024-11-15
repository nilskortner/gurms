package connectionservice

type ConnectionChannels struct {
	DataChan chan any
}

func NewConnectionChannels() *ConnectionChannels {
	return &ConnectionChannels{
		DataChan: make(chan any, 256),
	}
}

func (c *ConnectionChannels) Close() {
	close(c.DataChan)
}
