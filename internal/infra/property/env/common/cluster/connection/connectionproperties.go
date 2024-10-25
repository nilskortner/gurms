package connection

type ConnectionProperties struct {
	client *ConnectionClientProperties
	server *ConnectionServerProperties
}

func InitConnectionProperties() *ConnectionProperties {
	return &ConnectionProperties{
		client: NewConnectionClientProperties(),
		server: NewConnectionServerProperties(),
	}
}
