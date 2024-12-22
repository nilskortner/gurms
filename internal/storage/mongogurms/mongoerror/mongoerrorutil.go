package mongoerror

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Translate(err error) error {
	var writeErr mongo.WriteException
	if errors.As(err, &writeErr) {
		for _, we := range writeErr.WriteErrors {
			if we.HasErrorCodeWithMessage()
		}

	}

	var bulkErr mongo.BulkWriteException
	if errors.As(err, &bulkErr) {

	}

	return err
}
