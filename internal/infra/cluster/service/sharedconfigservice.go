package service

import "gurms/internal/storage/mongo"

const EXPIRABLE_RECORD_TTL = 60

type SharedConfigService struct {
	mongoClient *mongo.GurmsMongoClient
}

func NewSharedConfigService(properties *MongoProperties) *SharedConfigService {
	mongoClient := mongo.NewGurmsMongoClient(properties, "shared-config")

	structs := make([]any, 3)
	mongoClient.RegisterEntitiesByStructs(structs)
	for entityStruct := range structs {
		mongoClient.CreateCollectionIfNotExists(entityStruct)
	}

	return &SharedConfigService{
		mongoClient: mongoClient,
	}
}
