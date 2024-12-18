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
