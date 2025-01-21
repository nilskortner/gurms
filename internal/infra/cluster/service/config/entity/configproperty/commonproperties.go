package configproperty

import "gurms/internal/infra/property/env/common/cluster"

type CommonProperties struct {
	Cluster        *cluster.ClusterProperties `bson:",inline"`
	FlightRecorder *FlightRecorderProperties  `bson:",inline"`
	HealthCheck    *HealthCheckProperties     `bson:",inline"`
	Ip             *IpProperties              `bson:",inline"`
	Location       *LocationProperties        `bson:",inline"`
	Logging        *LoggingProperties         `bson:",inline"`
	Plugin         *PluginProperties          `bson:",inline"`
	Security       *SecurityProperties        `bson:",inline"`
	Shutdown       *ShutdownProperties        `bson:",inline"`
	UserStatus     *UserStatusProperties      `bson:",inline"`
}
