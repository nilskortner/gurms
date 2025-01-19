package service

import (
	"context"
	"errors"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/property/env/common/mongoproperties"
	"gurms/internal/storage/mongogurms"
	"gurms/internal/storage/mongogurms/operation/option"

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

func (s *SharedConfigService) Find(name string, filter *option.Filter) (*mongo.Cursor, error) {
	return s.mongoClient.Operations.FindMany(name, filter)
}

func (s *SharedConfigService) FindOne(name string, filter *option.Filter) *mongo.SingleResult {
	return s.mongoClient.Operations.FindOneWithFilter(name, filter)
}

func (s *SharedConfigService) Insert(record any) error {
	return s.mongoClient.Operations.Insert(record)
}

// TODO: Check
func (s *SharedConfigService) UpdateOne(name string, filter *option.Filter,
	update *option.Update) (*mongo.UpdateResult, error) {
	return s.mongoClient.Operations.UpdateOne(name, filter, update)
}

func (s *SharedConfigService) UpdateMany(name string, filter *option.Filter,
	update *option.Update) (*mongo.UpdateResult, error) {
	return s.mongoClient.Operations.UpdateMany(name, filter, update)
}

func (s *SharedConfigService) Upsert(name string, filter *option.Filter,
	update *option.Update, value any) error {

	result, err := s.UpdateOne(name, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		err := s.mongoClient.Operations.Insert(value)
		if err != nil {
			var errWE mongo.WriteException
			if errors.As(err, &errWE) {
				if errWE.HasErrorCode(11000) {
					s.Upsert(name, filter, update, value)
				}
			}
			return err
		}
	}
	return nil
}

func (s *SharedConfigService) RemoveOne(name string, filter *option.Filter) (*mongo.DeleteResult, error) {
	return s.mongoClient.Operations.DeleteOne(name, filter)
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
