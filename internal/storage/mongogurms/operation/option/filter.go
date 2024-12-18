package option

import "go.mongodb.org/mongo-driver/v2/bson"

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
