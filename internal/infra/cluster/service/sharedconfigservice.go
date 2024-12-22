package service

import (
	"context"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/mongoproperties"
	"gurms/internal/storage/mongogurms"
	"gurms/internal/storage/mongogurms/operation/option"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var SHAREDCONFIGLOGGER logger.Logger = factory.GetLogger("SharedConfig")

const EXPIRABLE_RECORD_TTL = 60

type SharedConfigService struct {
	mongoClient *mongogurms.GurmsMongoClient
}

// TODO: make bson tags for structs

// TODO: look into sharding functions

func NewSharedConfigService(properties *mongoproperties.MongoProperties) *SharedConfigService {
	mongoClient, err := mongogurms.NewGurmsMongoClient(properties, "shared-config")
	if err != nil {
		SHAREDCONFIGLOGGER.FatalWithError("failed to create the shared config service", err)
	}
	ctx := context.Background()

	setIndexes(ctx, mongoClient.Ctx.Database)

	return &SharedConfigService{
		mongoClient: mongoClient,
	}
}

func setIndexes(ctx context.Context, database *mongo.Database) {
	ensureSharedClusterPropertiesIndexes(database)
	ensureLeaderIndexes(database)
	ensureMemberIndexes(database)
}

func (s *SharedConfigService) Subscribe(name string, opts *options.ChangeStreamOptionsBuilder) (*mongo.ChangeStream, error) {
	return s.mongoClient.Operations.Watch(name, opts)
}

func (s *SharedConfigService) Insert(record any) error {
	return s.mongoClient.Operations.Insert(record)
}

func (s *SharedConfigService) updateOne(filter *option.Filter, update *option.Update, entity string, upsert bool) error {
	collection := s.mongoClient.Ctx.Database.Collection(entity)
	option := options.Update().SetUpsert(true)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.UpdateOne(ctx, filter.Document, update.Update, option)
	if err != nil {
		return err
	}
	return nil
}

func (s *SharedConfigService) Upsert(filter *option.Filter, update *option.Update, entity string) error {
	return s.updateOne(filter, update, entity, true)
}

// region Indexation

func ensureSharedClusterPropertiesIndexes(database *mongo.Database) {
	sharedclusterproperties := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	ctx := context.Background()
	name := "sharedclusterproperties"
	database.CreateCollection(ctx, name)
	collection := database.Collection(name)
	_, err := collection.Indexes().CreateOne(ctx, sharedclusterproperties)
	if err != nil {
		SHAREDCONFIGLOGGER.FatalWithError("couldnt create index for sharedclusterproperties: ", err)
	}
}

func ensureLeaderIndexes(database *mongo.Database) {
	renew := mongo.IndexModel{
		Keys:    bson.D{{Key: "renew_date", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(60),
	}
	clusterId := mongo.IndexModel{
		Keys:    bson.D{{Key: "cluster_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	ctx := context.Background()
	name := "leader"
	database.CreateCollection(ctx, name)
	collection := database.Collection(name)
	_, err := collection.Indexes().CreateOne(ctx, clusterId)
	if err != nil {
		SHAREDCONFIGLOGGER.FatalWithError("couldnt create clusterId index for leader: ", err)
	}
	_, err = collection.Indexes().CreateOne(ctx, renew)
	if err != nil {
		SHAREDCONFIGLOGGER.FatalWithError("couldnt create renew index for leader: ", err)
	}
}
func ensureMemberIndexes(database *mongo.Database) {
	member := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	ctx := context.Background()
	name := "member"
	database.CreateCollection(ctx, name)
	collection := database.Collection(name)
	_, err := collection.Indexes().CreateOne(ctx, member)
	if err != nil {
		SHAREDCONFIGLOGGER.FatalWithError("couldnt create index for member: ", err)
	}
}

// end region
