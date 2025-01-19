package operation

import (
	"gurms/internal/storage/mongogurms/operation/option"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoOperationsSupport interface {
	FindOne(name string) *mongo.SingleResult
	FindOneWithFilter(name string, filter *option.Filter) *mongo.SingleResult
	FindOneWithOptions(name string, filter *option.Filter,
		opts *options.FindOneOptionsBuilder) *mongo.SingleResult
	FindMany(name string, filter *option.Filter) (*mongo.Cursor, error)
	FindManyWithOptions(name string, filter *option.Filter,
		options *options.FindOptionsBuilder) (*mongo.Cursor, error)

	Insert(value any) error
	Watch(name string, opts *options.ChangeStreamOptionsBuilder) (*mongo.ChangeStream, error)
	UpdateOne(name string, filter *option.Filter, update *option.Update) (*mongo.UpdateResult, error)
	UpdateMany(name string, filter *option.Filter, update *option.Update) (*mongo.UpdateResult, error)
	DeleteOne(name string, filter *option.Filter) (*mongo.DeleteResult, error)
}
