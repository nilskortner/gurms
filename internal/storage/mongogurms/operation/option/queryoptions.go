package option

import (
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type QueryOptions struct {
	Opts *options.FindOneOptionsBuilder
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Opts: options.FindOne(),
	}
}
