package mongo

type MongoContext struct {
	client         *GurmsMongoClient
	database       *MongoDataBase
	adminDatabase  *MongoDataBase
	configDatabase *MongoDataBase
}
