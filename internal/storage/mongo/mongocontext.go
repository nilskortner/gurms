package mongo

import (
	"gurms/internal/storage/mongo/entity"

	"github.com/cornelk/hashmap"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoContext struct {
	client         *GurmsMongoClient
	database       *mongo.Database
	adminDatabase  *mongo.Database
	configDatabase *mongo.Database

	nameToEntity     hashmap.Map[string, *entity.MongoEntityWrap]
	nameToCollection hashmap.Map[string, *mongo.Collection]
}

func NewMongoContext(connectionString string, onServerDescriptionChange func([]event.ServerDescription)) *MongoContext {

}
