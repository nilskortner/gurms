package mongo

import (
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/logger"
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
