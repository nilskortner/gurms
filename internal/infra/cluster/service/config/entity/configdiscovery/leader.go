package configdiscovery

import "time"

type Leader struct {
	ClusterId  string    `bson:"cluster_id"`
	NodeId     string    `bson:"node_id"`
	RemewDate  time.Time `bson:"renew_date"`
	Generation int       `bson:"generation"`
}
