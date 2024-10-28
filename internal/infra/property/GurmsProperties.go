package property

import (
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/cluster"
	"gurms/internal/infra/property/env/common/healthcheck"
	"gurms/internal/infra/property/env/common/location"
	"gurms/internal/infra/property/env/common/logging"
	"gurms/internal/infra/property/env/common/security"
)

const PROPERTIES_PREFIX = "gurms"

const SCHEMA_VERSION = 1

type GurmsProperties struct {
	// Common
	cluster        *cluster.ClusterProperties
	flightRecorder *common.FlightRecorderProperties
	healthCheck    *healthcheck.HealthCheckProperties
	ip             *common.IpProperties
	location       *location.LocationProperties
	logging        *logging.LoggingProperties
	security       *security.SecurityProperties
	shutdown       *common.ShutdownProperties
	userStatus     *common.UserStatusProperties
	// AI Serving, Gateway and Service
	//aiServing *AiServingProperties
	//gateway   *GatewayProperties
	//service   *ServiceProperties
	//mockNode  *MockNodeProperties
}

func InitGurmsProperties() *GurmsProperties {
	return &GurmsProperties{
		cluster:        cluster.NewClusterProperties(),
		flightRecorder: common.NewFlightRecorderProperties(),
		healthCheck:    healthcheck.NewHealthCheckProperties(),
		ip:             common.NewIpProperties(),
		location:       location.NewLocationProperties(),
		logging:        logging.NewLoggingProperties(),
		security:       security.NewSecurityProperties(),
		shutdown:       common.NewShutdonwProperties(),
		userStatus:     common.NewUserStatusProperties(),
	}
}
