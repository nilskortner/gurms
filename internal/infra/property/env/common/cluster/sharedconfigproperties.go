package cluster

import "gurms/internal/infra/property/env/common/mongoproperties"

type SharedConfigProperties struct {
	Mongo *mongoproperties.MongoProperties
}

func NewSharedConfigProperties() *SharedConfigProperties {
	return &SharedConfigProperties{
		Mongo: mongoproperties.NewMongoProperties(),
	}
}
