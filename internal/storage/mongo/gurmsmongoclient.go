package mongo

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/mongoproperties"

	"go.mongodb.org/mongo-driver/v2/event"
)

var GURMSMONGOCLIENTLOGGER logger.Logger = factory.GetLogger("GurmsMongoClient")

type GurmsMongoClient struct {
	connect      chan struct{}
	names        map[string]struct{}
	descriptions []event.ServerDescription
	ctx          *MongoContext
	operations   *GurmsMongoOperations
}

func NewGurmsMongoClient(properties *mongoproperties.MongoProperties, name string) *GurmsMongoClient {
	emptySet := make(map[*event.TopologyDescription]struct{}, 0)
	return NewGurmsMongoClientWithSet(properties, name, emptySet)
}

func NewGurmsMongoClientWithSet(properties *mongoproperties.MongoProperties, name string, requiredClusterType map[*event.TopologyDescription]struct{}) *GurmsMongoClient {
	connect := make(chan (struct{}))
	client := gurmsMongoClient(properties, name, requiredClusterType, connect)
	return client
}

func gurmsMongoClient(properties *mongoproperties.MongoProperties,
	name string,
	requiredClusterType map[*event.TopologyDescription]struct{},
	connect chan struct{}) *GurmsMongoClient {

	operations := NewGurmsMongoOperations(ctx)
}

func verifyServers() {

}
