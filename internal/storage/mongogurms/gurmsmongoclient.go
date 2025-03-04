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

type GurmsMongoClient struct {
	names               map[string]struct{}
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

	names := make(map[string]struct{}, 8)
	names[name] = struct{}{}
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
	err = ctx.Client.Ping(pingCtx, nil)
	if err != nil {
		return nil, err
	}
	var operations operation.MongoOperationsSupport = operation.NewGurmsMongoOperations(ctx)

	return &GurmsMongoClient{
		names:               names,
		TopologyDescription: &topologyDescription,
		Ctx:                 ctx,
		Operations:          operations,
	}, nil
}

func (g *GurmsMongoClient) destroy(timeoutMillis int64) error {
	return g.Ctx.destroy(timeoutMillis)
}

func (g *GurmsMongoClient) verifyClusterTypes(name string, requiredClusterTypes map[string]struct{}) {
	if g.TopologyDescription == nil {
		GURMSMONGOCLIENTLOGGER.Fatal(
			"verification can only work after the mongo client has been initialized")
	}
	g.names[name] = struct{}{}
	if len(requiredClusterTypes) == 0 {
		return
	}
	for _, description := range g.TopologyDescription.Servers {
		_, ok := requiredClusterTypes[description.Kind]
		if !ok {
			GURMSMONGOCLIENTLOGGER.Fatal(fmt.Sprintf(
				"the cluster types for the mongo clients %v must be one of the types: %v",
				g.names, requiredClusterTypes))
		}
	}
}

func verifyServers(topologyDescription *event.TopologyDescription,
	name string,
	requiredClusterTypes map[string]struct{}) error {

	for _, description := range topologyDescription.Servers {
		if description.MaxWireVersion < 8 {
			return fmt.Errorf("the version of MongoDB server should be at least 4.2. ")
		}
		_, ok := requiredClusterTypes[description.Kind]
		if len(requiredClusterTypes) != 0 && !ok {
			return fmt.Errorf("the cluster types for the mongo client %s \" must be one of the types: %s",
				name, requiredClusterTypes)
		}
	}
	return nil
}
