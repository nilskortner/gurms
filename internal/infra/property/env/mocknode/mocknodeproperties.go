package mocknode

type MockNodeProperties struct {
	// API
	adminApi *AdminApiProperties `bson:",inline"`

	// Data Store
	mongo *MongoGroupProperties `bson:",inline"`
}
