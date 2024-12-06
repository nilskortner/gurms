package mongo

import (
	"context"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/mongoproperties"
	"time"

	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/description"
)

var GURMSMONGOCLIENTLOGGER logger.Logger = factory.GetLogger("GurmsMongoClient")

var names map[string]struct{} = make(map[string]struct{}, 8)

type GurmsMongoClient struct {
	descriptions []event.ServerDescription
	ctx          *MongoContext
	operations   *GurmsMongoOperations
}

func NewGurmsMongoClient(properties *mongoproperties.MongoProperties, name string) (*GurmsMongoClient, error) {
	emptySet := make(map[*event.TopologyDescription]struct{}, 0)
	return NewGurmsMongoClientWithSet(properties, name, emptySet)
}

func NewGurmsMongoClientWithSet(properties *mongoproperties.MongoProperties,
	name string,
	requiredClusterTypes map[*event.TopologyDescription]struct{}) (*GurmsMongoClient, error) {

	return gurmsMongoClient(properties, name, requiredClusterTypes)
}

func gurmsMongoClient(properties *mongoproperties.MongoProperties,
	name string,
	requiredClusterTypes map[*event.TopologyDescription]struct{}) (*GurmsMongoClient, error) {

	var descriptions []event.ServerDescription
	ctx, err := NewMongoContext(properties.Uri, func(sd []event.ServerDescription) {
		for _, server := range sd {
			if server.Kind == description.UnknownStr {
				return
			}
		}
		descriptions = sd
	})
	if err != nil {
		return nil, err
	}

	verifyServers(descriptions, name, requiredClusterTypes)
	pingCtx, cancelPing := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelPing()
	ctx.client.Ping(pingCtx, nil)
	if err != nil {
		return nil, err
	}
	operations := NewGurmsMongoOperations(ctx)

	return &GurmsMongoClient{
		descriptions: descriptions,
		ctx:          ctx,
		operations:   operations,
	}, nil
}

func verifyServers() {

}
