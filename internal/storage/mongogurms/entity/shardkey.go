package entity

import "go.mongodb.org/mongo-driver/v2/bson"

type ShardKey struct {
	document bson.M
}

type Path struct {
	isIdField bool
	fullPath  string
	path      []string
}
