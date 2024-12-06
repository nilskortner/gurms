package mongoproperties

type MongoProperties struct {
	Uri string
}

var uri string = "mongodb://localhost:27017/gurms-config"

func NewMongoProperties() *MongoProperties {
	return &MongoProperties{
		Uri: uri,
	}
}
