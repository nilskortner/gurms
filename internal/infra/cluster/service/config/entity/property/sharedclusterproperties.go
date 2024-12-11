package property

import "time"

// TODO: Make bson tags

type SharedClusterProperties struct {
	ClusterId           string               `bson:"_id"`
	SchemaVersion       int                  `bson:"schema_version"`
	CommonProperties    *CommonProperties    `bson:",inline"`
	AiServingProperties *AiServingProperties `bson:",inline"`
	GatewayProperties   *GatewayProperties   `bson:",inline"`
	ServiceProperties   *ServiceProperties   `bson:",inline"`
	MockNodeProperties  *MockNodeProperties  `bson:",inline"`
	LastUpdateTime      *time.Time           `bson:"last_updated_time"`
	GurmsProperties     *GurmsProperties     `bson:",inline"`
}
