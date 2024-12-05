package entity

type MongoEntityWrap interface {
	GetShardKeyBson()
}

type MongoEntity[T comparable] struct {
	constructor    func() T
	collectionName string
	jsonSchema     BsonDocument
	//shard key zone index
	shardKey        *ShardKey
	zone            *Zone
	compoundIndexes []*CompoundIndex
	indexes         []*Index
	//Field
	idFiledName string
	nameToField map[string]*EntityField
	fields      []*EntityField
}

func (m *MongoEntity[T]) GetShardKeyBson() {}
