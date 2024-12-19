package operation

type MongoOperationsSupport interface {
	FindById()
	FindOne()
	FindOneWithFilter()
	FindOneWithQueryOptions()

	Insert(value any) error
}
