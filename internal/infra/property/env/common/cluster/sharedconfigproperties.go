package cluster

import "gurms/internal/infra/property/env/common/mongo"

type SharedConfigProperties struct {
	mongo *mongo.MongoProperties
}

func NewSharedConfigProperties() *SharedConfigProperties {
	return &SharedConfigProperties{
		mongo: mongo.NewMongoProperties(),
	}
}
