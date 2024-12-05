package operation

type MongoOperationsSupport interface {
	FindById()
	FindOne()
	FindOneWithFilter()
	FindOneWithQueryOptions()
}
