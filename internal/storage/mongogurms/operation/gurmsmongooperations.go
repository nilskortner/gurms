package operation

import (
	"context"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/storage/mongogurms/operation/option"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var GURMSMONGOOPERATIONSLOGGER logger.Logger = factory.GetLogger("GurmsMongoOperationsLogger")

var FILTER_ALL *option.Filter = option.NewFilter()

// CRUD
var ()

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

// Query

func (g *GurmsMongoOperations) FindOne(name string) *mongo.SingleResult {
	return g.FindOneWithOptions(name, FILTER_ALL, nil)
}

func (g *GurmsMongoOperations) FindOneWithFilter(name string,
	filter *option.Filter) *mongo.SingleResult {
	return g.FindOneWithOptions(name, filter, nil)
}

func (g *GurmsMongoOperations) FindOneWithOptions(name string, filter *option.Filter,
	opts *options.FindOneOptionsBuilder) *mongo.SingleResult {

	collection := g.ctx.GetDatabaseCollection(name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	return collection.FindOne(ctx, filter.Document, opts)
}

func (g *GurmsMongoOperations) FindMany(name string, filter *option.Filter) (*mongo.Cursor, error) {
	return g.FindManyWithOptions(name, filter, nil)
}

func (g *GurmsMongoOperations) FindManyWithOptions(name string, filter *option.Filter,
	options *options.FindOptionsBuilder) (*mongo.Cursor, error) {
	collection := g.ctx.GetDatabaseCollection(name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	return collection.Find(ctx, filter.Document, options)
}

func (g *GurmsMongoOperations) Watch(name string, opts *options.ChangeStreamOptionsBuilder) (*mongo.ChangeStream, error) {
	collection := g.ctx.GetDatabaseCollection(name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
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

// Update

func (g *GurmsMongoOperations) UpdateOne(name string,
	filter *option.Filter, update *option.Update) (*mongo.UpdateResult, error) {
	return g.UpdateOneWithSession(nil, name, filter, update)
}

func (g *GurmsMongoOperations) UpdateOneWithSession(session *mongo.Session,
	name string, filter *option.Filter, update *option.Update) (*mongo.UpdateResult, error) {
	collection := g.ctx.GetDatabaseCollection(name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var source *mongo.UpdateResult
	var err error
	sessionFn := func(ctx context.Context) error {
		source, err = collection.UpdateOne(ctx, filter.Document, update.Update)
		if err != nil {
			return err
		}
		return nil
	}

	if session == nil {
		source, err = collection.UpdateOne(ctx, filter.Document, update.Update)
		if err != nil {
			return nil, err
		}
	} else {
		err := mongo.WithSession(ctx, session, sessionFn)
		if err != nil {
			return nil, err
		}
	}
	return source, nil
}

func (g *GurmsMongoOperations) UpdateMany(name string,
	filter *option.Filter, update *option.Update) (*mongo.UpdateResult, error) {
	return g.UpdateManyWithSession(nil, name, filter, update)
}

func (g *GurmsMongoOperations) UpdateManyWithSession(session *mongo.Session,
	name string, filter *option.Filter, update *option.Update) (*mongo.UpdateResult, error) {
	collection := g.ctx.GetDatabaseCollection(name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var source *mongo.UpdateResult
	var err error
	sessionFn := func(ctx context.Context) error {
		source, err = collection.UpdateMany(ctx, filter.Document, update.Update)
		if err != nil {
			return err
		}
		return nil
	}

	if session == nil {
		source, err = collection.UpdateMany(ctx, filter.Document, update.Update)
		if err != nil {
			return nil, err
		}
	} else {
		err := mongo.WithSession(ctx, session, sessionFn)
		if err != nil {
			return nil, err
		}
	}
	return source, nil
}

// Delete
func (g *GurmsMongoOperations) DeleteOneWithSession(session *mongo.Session, name string,
	filter *option.Filter) (*mongo.DeleteResult, error) {

	collection := g.ctx.GetDatabaseCollection(name)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	var source *mongo.DeleteResult
	var err error
	sessionFn := func(ctx context.Context) error {
		source, err = collection.DeleteOne(ctx, filter.Document)
		if err != nil {
			return err
		}
		return nil
	}

	if session == nil {
		source, err = collection.DeleteOne(ctx, filter.Document)
		if err != nil {
			return nil, err
		}
	} else {
		err := mongo.WithSession(ctx, session, sessionFn)
		if err != nil {
			return nil, err
		}
	}
	return source, nil
}

func (g *GurmsMongoOperations) DeleteOne(name string, filter *option.Filter) (*mongo.DeleteResult, error) {
	return g.DeleteOneWithSession(nil, name, filter)
}

// TODO: delete?

// func (g *GurmsMongoOperations) getPublisher(collection *mongo.Collection,
// 	filter bson.M, queryOptions *option.QueryOptions) func() (*mongo.SingleResult, error) {

// 	opts := []options.Lister[options.FindOneOptions]{queryOptions.Opts}

// 	publisher := func() (*mongo.SingleResult, error) {
// 		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
// 		defer cancel()
// 		return collection.FindOne(ctx, filter, opts...)
// 	}

// 	return publisher
// }
