package mongogurms

import (
	"gurms/internal/infra/property/env/common/mongoproperties"
	"sync"

	"golang.org/x/exp/maps"
)

// TODO: make shutdown for service and gateway mongoconfig

type BaseMongoConfig struct {
	uriToCLient map[string]*GurmsMongoClient
	mu          sync.Mutex
}

func NewBaseMongoConfig() *BaseMongoConfig {
	return &BaseMongoConfig{
		uriToCLient: make(map[string]*GurmsMongoClient, 8),
	}
}

func (b *BaseMongoConfig) destroy(timeoutMillis int64) error {
	clients := maps.Values(b.uriToCLient)
	size := len(clients)
	if size == 0 {
		return nil
	}
	errors := make([]error, size)
	for _, client := range clients {
		errors = append(errors, client.destroy(timeoutMillis))
	}
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BaseMongoConfig) getMongoClient(properties *mongoproperties.MongoProperties, name string,
	requiredClusterTypes map[string]struct{}) *GurmsMongoClient {

	b.mu.Lock()
	defer b.mu.Unlock()

	return func() *GurmsMongoClient {
		mongoClient := b.uriToCLient[properties.Uri]
		if mongoClient == nil {
			client, err := NewGurmsMongoClientWithSet(properties, name, requiredClusterTypes)
			if err != nil {
				GURMSMONGOCLIENTLOGGER.FatalWithError("failed to create the mongo client", err)
			}
			return client
		}
		mongoClient.verifyClusterTypes(name, requiredClusterTypes)
		return mongoClient
	}()
}
