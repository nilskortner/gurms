package configproperty

import "gurms/internal/infra/property/env/common/cluster"

type CommonProperties struct {
	Cluster        *cluster.ClusterProperties `bson:"clusterProperties"`
	FlightRecorder *FlightRecorderProperties  `bson:"flightRecorderProperties"`
	HealthCheck    *HealthCheckProperties     `bson:"healthCheckProperties"`
	Ip             *IpProperties              `bson:"ipProperties"`
	Location       *LocationProperties        `bson:"locationProperties"`
	Logging        *LoggingProperties         `bson:"loggingProperties"`
	Plugin         *PluginProperties          `bson:"pluginProperties"`
	Security       *SecurityProperties        `bson:"securityProperties"`
	Shutdown       *ShutdownProperties        `bson:"shutdownProperties"`
	UserStatus     *UserStatusProperties      `bson:"userStatusProperties"`
}
