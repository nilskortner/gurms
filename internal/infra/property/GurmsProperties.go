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
	Cluster        *cluster.ClusterProperties
	FlightRecorder *common.FlightRecorderProperties
	HealthCheck    *healthcheck.HealthCheckProperties
	Ip             *common.IpProperties
	Location       *location.LocationProperties
	Logging        *logging.LoggingProperties
	Security       *security.SecurityProperties
	Shutdown       *common.ShutdownProperties
	UserStatus     *common.UserStatusProperties
	// AI Serving, Gateway and Service
	//aiServing *AiServingProperties
	//gateway   *GatewayProperties
	//service   *ServiceProperties
	//mockNode  *MockNodeProperties
}

func InitGurmsProperties() *GurmsProperties {
	return &GurmsProperties{
		Cluster:        cluster.NewClusterProperties(),
		FlightRecorder: common.NewFlightRecorderProperties(),
		HealthCheck:    healthcheck.NewHealthCheckProperties(),
		Ip:             common.NewIpProperties(),
		Location:       location.NewLocationProperties(),
		Logging:        logging.NewLoggingProperties(),
		Security:       security.NewSecurityProperties(),
		Shutdown:       common.NewShutdonwProperties(),
		UserStatus:     common.NewUserStatusProperties(),
	}
}
