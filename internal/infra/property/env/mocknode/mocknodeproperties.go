package mocknode

type MockNodeProperties struct {
	// API
	adminApi *AdminApiProperties `bson:"adminApiProperties"`

	// Data Store
	mongo *MongoGroupProperties `bson:"mongoGroupProperties"`
}

func NewMockNodeProperties() *MockNodeProperties {
	return &MockNodeProperties{
		adminApi: NewAdminApiProperties(),
		mongo:    NewMongoGroupProperties(),
	}
}
