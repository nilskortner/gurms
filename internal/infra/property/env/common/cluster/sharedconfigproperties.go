package cluster

import "gurms/internal/infra/property/env/common/mongoproperties"

type SharedConfigProperties struct {
	mongo *mongoproperties.MongoProperties
}

func NewSharedConfigProperties() *SharedConfigProperties {
	return &SharedConfigProperties{
		mongo: mongoproperties.NewMongoProperties(),
	}
}
