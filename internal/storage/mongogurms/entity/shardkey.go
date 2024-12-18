package entity

type ShardKey struct {
	document BsonDocument
}

type Path struct {
	isIdField bool
	fullPath  string
	path      []string
}
