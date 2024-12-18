package mongogurms

import (
	"fmt"
	"gurms/internal/storage/mongogurms/entity"
	"net/url"
	"strings"

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
