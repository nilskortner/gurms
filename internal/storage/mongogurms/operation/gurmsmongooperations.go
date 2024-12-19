package operation

import (
	"context"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/storage/mongogurms"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var GURMSMONGOOPERATIONSLOGGER logger.Logger = factory.GetLogger("GurmsMongoOperationsLogger")

type GurmsMongoOperations struct {
	ctx *mongogurms.MongoContext
}

func (g *GurmsMongoOperations) Watch(name string) (*mongo.ChangeStream, error) {
	collection := g.ctx.Database.Collection(name)
	ctx := context.Background()
	stream, err := collection.Watch(ctx, nil)
	if err != nil {
		GURMSMONGOOPERATIONSLOGGER.FatalWithError("couldnt subscribe to stream", err)
		return nil, err
	}
	return stream, nil
}
func (g *GurmsMongoOperations) Insert(value any) error {
	return g.InsertWithSession(nil, value)
}

// TODO: mongoexception translator

func (g *GurmsMongoOperations) InsertWithSession(session *mongo.Session, value any) error {
	collection, err := g.ctx.GetCollection(value)
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

	return err
}

func (g *GurmsMongoOperations) FindById()
func (g *GurmsMongoOperations) FindOne()
func (g *GurmsMongoOperations) FindOneWithFilter()
func (g *GurmsMongoOperations) FindOneWithQueryOptions()
