package operation

import (
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
}
