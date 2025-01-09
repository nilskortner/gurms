package option

import "go.mongodb.org/mongo-driver/v2/bson"

type Update struct {
	Update bson.M
}

func NewUpdate() *Update {
	updateMap := make(bson.M, 0)
	return &Update{
		Update: updateMap,
	}
}

func (u *Update) Set(name string, value any) {
	if name != "" && value != nil {
		u.Update[name] = value
	}
}
