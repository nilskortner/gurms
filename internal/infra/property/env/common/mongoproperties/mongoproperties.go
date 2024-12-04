package mongoproperties

type MongoProperties struct {
	uri string
}

var uri string = "mongodb://localhost:27017/gurms-config"

func NewMongoProperties() *MongoProperties {
	return &MongoProperties{
		uri: uri,
	}
}
