package gateway

type GatewayProperties struct {
	// API
	adminApi            *AdminApiProperties            `bson:",inline"`
	clientApi           *ClientApiProperties           `bson:",inline"`
	notificationLogging *NotificationLoggingProperties `bson:",inline"`

	// Business Behavior
	session           *SessionProperties           `bson:",inline"`
	simultaneousLogin *SimultaneousLoginProperties `bson:",inline"`

	// Cluster
	serviceDiscovery *DiscoveryProperties `bson:",inline"`

	// Data Store
	mongo *MongoGroupProperties `bson:",inline"`
	redis *TurmsRedisProperties `bson:",inline"`

	// Faking
	fake *FakeProperties `bson:",inline"`

	// Network Access Layer
	udp       *UdpProperties       `bson:",inline"`
	tcp       *TcpProperties       `bson:",inline"`
	websocket *WebSocketProperties `bson:",inline"`
}
