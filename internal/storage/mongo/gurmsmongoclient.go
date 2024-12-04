package mongo

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/mongoproperties"
)

var GURMSMONGOCLIENTLOGGER logger.Logger = factory.GetLogger("GurmsMongoClient")

type GurmsMongoClient struct {
	connect chan struct{}
}

func NewGurmsMongoClient(properties *mongoproperties.MongoProperties, name string) *GurmsMongoClient {
	emptySet := make(map[*ClusterType]struct{}, 0)
	return NewGurmsMongoClientWithSet(properties, name, emptySet)
}

func NewGurmsMongoClientWithSet(properties *mongoproperties.MongoProperties, name string, requiredClusterType map[*ClusterType]struct{}) *GurmsMongoClient {
	connect := make(chan (struct{}))
	client := gurmsMongoClient(properties, name, requiredClusterType, connect)
	return client
}

func gurmsMongoClient(properties *mongoproperties.MongoProperties,
	name string,
	requiredClusterType map[*ClusterType]struct{},
	connect chan struct{}) *GurmsMongoClient {

	operations := NewGurmsMongoOperations(ctx)
}
