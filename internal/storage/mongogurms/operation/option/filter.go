package option

import "go.mongodb.org/mongo-driver/v2/bson"

// TODO: make document bson.D for ordering of filtering?

type Filter struct {
	Document bson.M
}

func NewFilter() *Filter {
	return &Filter{
		Document: make(bson.M),
	}
}

func (f *Filter) Eq(key string, value any) {
	f.Document[key] = value
}

func (f *Filter) In(key string, value any) {
	f.Document[key] = bson.M{"$in": value}
}

func (f *Filter) Lt(key string, value any) {
	f.Document[key] = bson.M{"$lt": value}
}
