package property

import (
	"gurms/internal/infra/property/env/aiserving"
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/cluster"
	"gurms/internal/infra/property/env/common/healthcheckproperty"
	"gurms/internal/infra/property/env/common/location"
	"gurms/internal/infra/property/env/common/logging"
	"gurms/internal/infra/property/env/common/plugin"
	"gurms/internal/infra/property/env/common/security"
	"gurms/internal/infra/property/env/gateway"
	"gurms/internal/infra/property/env/mocknode"
	"gurms/internal/infra/property/env/service"
)

const PROPERTIES_PREFIX = "gurms"

const SCHEMA_VERSION = 1

type GurmsProperties struct {
	// Common
	Cluster        *cluster.ClusterProperties
	FlightRecorder *common.FlightRecorderProperties
	HealthCheck    *healthcheckproperty.HealthCheckProperties
	Ip             *common.IpProperties
	Location       *location.LocationProperties
	Logging        *logging.LoggingProperties
	Plugin         *plugin.PluginProperties
	Security       *security.SecurityProperties
	Shutdown       *common.ShutdownProperties
	UserStatus     *common.UserStatusProperties
	// AI Serving, Gateway and Service
	AiServing *aiserving.AiServingProperties
	Gateway   *gateway.GatewayProperties
	Service   *service.ServiceProperties
	MockNode  *mocknode.MockNodeProperties
}

func NewGurmsProperties() *GurmsProperties {
	return &GurmsProperties{
		Cluster:        cluster.NewClusterProperties(),
		FlightRecorder: common.NewFlightRecorderProperties(),
		HealthCheck:    healthcheckproperty.NewHealthCheckProperties(),
		Ip:             common.NewIpProperties(),
		Location:       location.NewLocationProperties(),
		Logging:        logging.NewLoggingProperties(),
		Security:       security.NewSecurityProperties(),
		Shutdown:       common.NewShutdonwProperties(),
		UserStatus:     common.NewUserStatusProperties(),
	}
}
