package mongogurms

import (
	"context"
	"fmt"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/mongoproperties"
	"gurms/internal/storage/mongogurms/operation"
	"time"

	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/description"
)

var GURMSMONGOCLIENTLOGGER logger.Logger = factory.GetLogger("GurmsMongoClient")

var names map[string]struct{} = make(map[string]struct{}, 8)

type GurmsMongoClient struct {
	TopologyDescription *event.TopologyDescription
	Ctx                 *MongoContext
	Operations          operation.MongoOperationsSupport
}

func NewGurmsMongoClient(properties *mongoproperties.MongoProperties, name string) (*GurmsMongoClient, error) {
	emptySet := make(map[string]struct{}, 0)
	return NewGurmsMongoClientWithSet(properties, name, emptySet)
}

func NewGurmsMongoClientWithSet(properties *mongoproperties.MongoProperties,
	name string,
	requiredClusterTypes map[string]struct{}) (*GurmsMongoClient, error) {

	return gurmsMongoClient(properties, name, requiredClusterTypes)
}

func gurmsMongoClient(properties *mongoproperties.MongoProperties,
	name string,
	requiredClusterTypes map[string]struct{}) (*GurmsMongoClient, error) {

	var topologyDescription event.TopologyDescription
	ctx, err := NewMongoContext(properties.Uri, func(sd []event.ServerDescription) {
		for _, server := range sd {
			if server.Kind == description.UnknownStr {
				return
			}
		}
		topologyDescription.Servers = sd
	})
	if err != nil {
		return nil, err
	}

	err = verifyServers(&topologyDescription, name, requiredClusterTypes)
	if err != nil {
		return nil, err
	}
	pingCtx, cancelPing := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelPing()
	ctx.Client.Ping(pingCtx, nil)
	if err != nil {
		return nil, err
	}
	var operations operation.MongoOperationsSupport = operation.NewGurmsMongoOperations(ctx)

	return &GurmsMongoClient{
		TopologyDescription: &topologyDescription,
		Ctx:                 ctx,
		Operations:          operations,
	}, nil
}

func verifyServers(topologyDescription *event.TopologyDescription,
	name string,
	requiredClusterTypes map[string]struct{}) error {

	for _, description := range topologyDescription.Servers {
		if description.MaxWireVersion < 8 {
			return fmt.Errorf("the version of MongoDB server should be at least 4.2.")
		}
		_, ok := requiredClusterTypes[description.Kind]
		if len(requiredClusterTypes) != 0 && !ok {
			return fmt.Errorf("the cluster types for the mongo client %s \" must be one of the types: %s",
				name, requiredClusterTypes)
		}
	}
	return nil
}
