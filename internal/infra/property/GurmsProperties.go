package property

import (
	c "gurms/internal/infra/property/env/common/cluster"
)

const PROPERTIES_PREFIX = "gurms"

const SCHEMA_VERSION = 1

type GurmsProperties struct {
	// Common
	cluster        *c.ClusterProperties
	flightRecorder *FlightRecorderProperties
	healthCheck    *HealthCheckProperties
	ip             *IpProperties
	location       *LocationProperties
	logging        *LoggingProperties
	plugin         *PluginProperties
	security       *SecurityProperties
	shutdown       *ShutdownProperties
	userStatus     *UserStatusProperties
	// AI Serving, Gateway and Service
	aiServing *AiServingProperties
	gateway   *GatewayProperties
	service   *ServiceProperties
	mockNode  *MockNodeProperties
}

func InitGurmsProperties() *GurmsProperties {

}
