package connection

type ConnectionProperties struct {
	Client *ConnectionClientProperties
	Server *ConnectionServerProperties
}

func InitConnectionProperties() *ConnectionProperties {
	return &ConnectionProperties{
		Client: NewConnectionClientProperties(),
		Server: NewConnectionServerProperties(),
	}
}
