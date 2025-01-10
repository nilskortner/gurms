package operation

import (
	"gurms/internal/storage/mongogurms/operation/option"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoOperationsSupport interface {
	FindById()
	FindOne()
	FindOneWithFilter()
	FindOneWithQueryOptions()

	Insert(value any) error
	Watch(name string, opts *options.ChangeStreamOptionsBuilder) (*mongo.ChangeStream, error)
	UpdateOne(name string, filter *option.Filter, update *option.Update) (*mongo.UpdateResult, error)
	UpdateMany(name string, filter *option.Filter, update *option.Update) (*mongo.UpdateResult, error)
	DeleteOne(name string, filter *option.Filter) (*mongo.DeleteResult, error)
}
