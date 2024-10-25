package mongo

type MongoProperties struct {
	uri string
}

var uri string = "mongodb://localhost:27017/turms-config"

func NewMongoProperties() *MongoProperties {
	return &MongoProperties{
		uri: uri,
	}
}
