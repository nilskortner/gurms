package mongogurms

import (
	"context"
	"fmt"
	"gurms/internal/infra/cluster/service/config/entity/configdiscovery"
	"gurms/internal/storage/mongogurms/entity"
	"net/url"
	"strings"
	"time"

	"github.com/cornelk/hashmap"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var nameToEntity *hashmap.Map[string, entity.MongoEntityWrap] = hashmap.New[string, entity.MongoEntityWrap]()
var nameToCollection *hashmap.Map[string, *mongo.Collection] = hashmap.New[string, *mongo.Collection]()

type MongoContext struct {
	Client         *mongo.Client
	Database       *mongo.Database
	AdminDatabase  *mongo.Database
	ConfigDatabase *mongo.Database
}

func NewMongoContext(connectionString string, onServerDescriptionChange func([]event.ServerDescription)) (*MongoContext, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("the connection string must not be empty")
	}
	settings := options.Client().ApplyURI(connectionString)

	clusterMonitor := &event.ServerMonitor{
		TopologyDescriptionChanged: func(e *event.TopologyDescriptionChangedEvent) {
			onServerDescriptionChange(e.NewDescription.Servers)
		},
	}
	settings.SetServerMonitor(clusterMonitor)
	settings.SetRegistry(CODEC_REGISTRY)

	//

	client, err := mongo.Connect(settings)
	if err != nil {
		return nil, fmt.Errorf("couldnt create new client", err)
	}

	databaseName, err := getDatabaseFromConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return &MongoContext{
		Client:         client,
		Database:       client.Database(databaseName),
		AdminDatabase:  client.Database("admin"),
		ConfigDatabase: client.Database("config"),
	}, nil
}

func (m *MongoContext) destroy(timeoutMillis int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMillis))
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func (m *MongoContext) GetCollection(value any) (*mongo.Collection, error) {
	switch value.(type) {
	case configdiscovery.Leader:
		return m.Database.Collection("leader"), nil
	case configdiscovery.Member:
		return m.Database.Collection("member"), nil
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

func getDatabaseFromConnectionString(connectionString string) (string, error) {
	// Parse the connection string using net/url
	parsedURL, err := url.Parse(connectionString)
	if err != nil {
		return "", fmt.Errorf("invalid connection string: %w", err)
	}

	// Extract the path (the part after the host)
	if parsedURL.Path == "" || parsedURL.Path == "/" {
		return "", fmt.Errorf("no database specified in connection string")
	}

	// Remove the leading "/" and return the database name
	database := strings.TrimPrefix(parsedURL.Path, "/")
	return database, nil
}

// region: injection functions
func (m *MongoContext) GetDatabaseCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

func (m *MongoContext) GetCollectionByValue(value any) (*mongo.Collection, error) {
	return m.GetCollection(value)
}
