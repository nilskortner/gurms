package property

import (
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/cluster"
	"gurms/internal/infra/property/env/common/healthcheck"
)

const PROPERTIES_PREFIX = "gurms"

const SCHEMA_VERSION = 1

type GurmsProperties struct {
	// Common
	cluster        *cluster.ClusterProperties
	flightRecorder *common.FlightRecorderProperties
	healthCheck    *healthcheck.HealthCheckProperties
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
