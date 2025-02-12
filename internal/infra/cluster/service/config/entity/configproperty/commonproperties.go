package configproperty

import (
	"gurms/internal/infra/property/env/common"
	"gurms/internal/infra/property/env/common/cluster"
	"gurms/internal/infra/property/env/common/healthcheckproperty"
	"gurms/internal/infra/property/env/common/location"
	"gurms/internal/infra/property/env/common/logging"
	"gurms/internal/infra/property/env/common/plugin"
	"gurms/internal/infra/property/env/common/security"
)

type CommonProperties struct {
	Cluster        *cluster.ClusterProperties                 `bson:"clusterProperties"`
	FlightRecorder *common.FlightRecorderProperties           `bson:"flightRecorderProperties"`
	HealthCheck    *healthcheckproperty.HealthCheckProperties `bson:"healthCheckProperties"`
	Ip             *common.IpProperties                       `bson:"ipProperties"`
	Location       *location.LocationProperties               `bson:"locationProperties"`
	Logging        *logging.LoggingProperties                 `bson:"loggingProperties"`
	Plugin         *plugin.PluginProperties                   `bson:"pluginProperties"`
	Security       *security.SecurityProperties               `bson:"securityProperties"`
	Shutdown       *common.ShutdownProperties                 `bson:"shutdownProperties"`
	UserStatus     *common.UserStatusProperties               `bson:"userStatusProperties"`
}
