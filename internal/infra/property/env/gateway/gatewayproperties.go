package gateway

import (
	"gurms/internal/infra/property/env/gateway/network"
	"gurms/internal/infra/property/env/gateway/redis"
	"gurms/internal/infra/property/env/gateway/session"
)

type GatewayProperties struct {
	// API
	adminApi            *AdminApiProperties            `bson:"adminApiProperties"`
	clientApi           *ClientApiProperties           `bson:"clientApiProperties"`
	notificationLogging *NotificationLoggingProperties `bson:"notificationLoggingProperties"`

	// Business Behavior
	session           *session.SessionProperties   `bson:"sessionProperties"`
	simultaneousLogin *SimultaneousLoginProperties `bson:"simultaneousLoginProperties"`

	// Cluster
	serviceDiscovery *DiscoveryProperties `bson:"discoveryProperties"`

	// Data Store
	mongo *MongoGroupProperties       `bson:"mongoGroupProperties"`
	redis *redis.GurmsRedisProperties `bson:"gurmsRedisProperties"`

	// Faking
	fake *FakeProperties `bson:"fakeProperties"`

	// Network Access Layer
	udp       *network.UdpProperties       `bson:"udpProperties"`
	tcp       *network.TcpProperties       `bson:"tcpProperties"`
	websocket *network.WebSocketProperties `bson:"webSocketProperties"`
}

func NewGatewayProperties() *GatewayProperties {
	return &GatewayProperties{
		adminApi:            NewAdminApiProperties(),
		clientApi:           NewClientApiProperties(),
		notificationLogging: NewNotificationLoggingProperties(),
		session:             session.NewSessionProperties(),
		simultaneousLogin:   NewSimultaneousLoginProperties(),
		serviceDiscovery:    NewDiscoveryProperties(),
		mongo:               NewMongoGroupProperties(),
		redis:               redis.NewGurmsRedisProperties(),
		fake:                NewFakeProperties(),
		udp:                 network.NewUdpProperties(),
		tcp:                 network.NewTcpProperties(),
		websocket:           network.NewWebSocketProperties(),
	}
}
