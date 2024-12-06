package mongo

import "go.mongodb.org/mongo-driver/v2/bson"

var CODEC_REGISTRY *bson.Registry = bson.NewRegistry()

// TODO: implement codecs for mongo bson
