package configproperty

import (
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/env/aiserving"
	"gurms/internal/infra/property/env/gateway"
	"gurms/internal/infra/property/env/mocknode"
	"gurms/internal/infra/property/env/service"
	"time"
)

const SHAREDCLUSTERPROPERTIESNAME = "sharedclusterproperties"

const (
	ID string = "_id"
)

type SharedClusterProperties struct {
	ClusterId           string                         `bson:"_id"`
	SchemaVersion       int                            `bson:"schema_version"`
	CommonProperties    *CommonProperties              `bson:",inline"`
	AiServingProperties *aiserving.AiServingProperties `bson:",inline"`
	GatewayProperties   *gateway.GatewayProperties     `bson:",inline"`
	ServiceProperties   *service.ServiceProperties     `bson:",inline"`
	MockNodeProperties  *mocknode.MockNodeProperties   `bson:",inline"`
	LastUpdateTime      time.Time                      `bson:"last_updated_time"`
	GurmsProperties     *property.GurmsProperties      // transient
}

func (s *SharedClusterProperties) tryInitGurmsProperties() {
	localCommonProperties := s.CommonProperties
	if localCommonProperties == nil {
		return
	}
	s.GurmsProperties = &property.GurmsProperties{
		Cluster:        localCommonProperties.Cluster,
		FlightRecorder: localCommonProperties.FlightRecorder,
		HealthCheck:    localCommonProperties.HealthCheck,
		Ip:             localCommonProperties.Ip,
		Location:       localCommonProperties.Location,
		Logging:        localCommonProperties.Logging,
		Plugin:         localCommonProperties.Plugin,
		Security:       localCommonProperties.Security,
		Shutdown:       localCommonProperties.Shutdown,
		UserStatus:     localCommonProperties.UserStatus,
		AiServing:      s.AiServingProperties,
		Gateway:        s.GatewayProperties,
		Service:        s.ServiceProperties,
		MockNode:       s.MockNodeProperties,
	}
}

func (s *SharedClusterProperties) SetCommonProperties(commonProperties *CommonProperties) {
	s.CommonProperties = commonProperties
	s.tryInitGurmsProperties()
}

func (s *SharedClusterProperties) SetAiServingProperties(aiServingProperties *aiserving.AiServingProperties) {
	s.AiServingProperties = aiServingProperties
	s.tryInitGurmsProperties()
}

func (s *SharedClusterProperties) SetGatewayProperties(gatewayProperties *gateway.GatewayProperties) {
	s.GatewayProperties = gatewayProperties
	s.tryInitGurmsProperties()
}

func (s *SharedClusterProperties) SetServiceProperties(serviceProperties *service.ServiceProperties) {
	s.ServiceProperties = serviceProperties
	s.tryInitGurmsProperties()
}

func (s *SharedClusterProperties) SetMocknodeProperties(mockNodeProperties *mocknode.MockNodeProperties) {
	s.MockNodeProperties = mockNodeProperties
	s.tryInitGurmsProperties()
}
