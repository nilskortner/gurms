package mongogurms

import (
	"context"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var GURMSOPERATIONSLOGGER logger.Logger = factory.GetLogger("GurmsOperations")

// TODO: implement functions and variables + implemented interface named support

type GurmsMongoOperations struct {
	ctx *MongoContext
}

func NewGurmsMongoOperations(ctx *MongoContext) *GurmsMongoOperations {
	return &GurmsMongoOperations{
		ctx: ctx,
	}
}

func (g *GurmsMongoOperations) Watch(name string) (*mongo.ChangeStream, error) {
	collection := g.ctx.Database.Collection(name)
	ctx := context.Background()
	stream, err := collection.Watch(ctx, nil)
	if err != nil {
		GURMSOPERATIONSLOGGER.FatalWithError("couldnt subscribe to stream", err)
		return nil, err
	}
	return stream, nil
}
