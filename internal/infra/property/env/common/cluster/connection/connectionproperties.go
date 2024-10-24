package connection

type ConnectionProperties struct {
	client *ConnectionCLientProperties
	server *ConnectionServerProperties
}

func InitConnectionProperties() *ConnectionProperties {
	return &ConnectionProperties{}
}
