package operation

import (
	"context"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var GURMSMONGOOPERATIONSLOGGER logger.Logger = factory.GetLogger("GurmsMongoOperationsLogger")

type GurmsMongoOperations struct {
	ctx MongoContextInjection
}

type MongoContextInjection interface {
	GetDatabaseCollection(name string) *mongo.Collection
	GetCollectionByValue(value any) (*mongo.Collection, error)
}

func NewGurmsMongoOperations(ctx MongoContextInjection) *GurmsMongoOperations {
	return &GurmsMongoOperations{
		ctx: ctx,
	}
}

func (g *GurmsMongoOperations) Watch(name string, opts *options.ChangeStreamOptionsBuilder) (*mongo.ChangeStream, error) {
	collection := g.ctx.GetDatabaseCollection(name)
	ctx := context.Background()
	stream, err := collection.Watch(ctx, mongo.Pipeline{}, opts)
	if err != nil {
		GURMSMONGOOPERATIONSLOGGER.FatalWithError("couldnt subscribe to stream", err)
		return nil, err
	}
	return stream, nil
}
func (g *GurmsMongoOperations) Insert(value any) error {
	return g.InsertWithSession(nil, value)
}

// TODO: do i need mongoexception translator?

func (g *GurmsMongoOperations) InsertWithSession(session *mongo.Session, value any) error {
	collection, err := g.ctx.GetCollectionByValue(value)
	if err != nil {
		return err
	}

	ctx := context.Background()
	insertFn := func(ctx context.Context) error {
		_, err = collection.InsertOne(ctx, value)
		return err
	}

	if session != nil {
		err = mongo.WithSession(ctx, session, insertFn)
	} else {
		_, err = collection.InsertOne(ctx, value)
	}
	if err != nil {
		return err
	}

	return nil
}

func (g *GurmsMongoOperations) FindById()
func (g *GurmsMongoOperations) FindOne()
func (g *GurmsMongoOperations) FindOneWithFilter()
func (g *GurmsMongoOperations) FindOneWithQueryOptions()
