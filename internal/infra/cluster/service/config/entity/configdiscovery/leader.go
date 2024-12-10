package configdiscovery

import "time"

type Leader struct {
	clusterId  string    `bson:"cluster_id"`
	nodeId     string    `bson:"node_id"`
	RemewDate  time.Time `bson:"renew_date"`
	generation int       `bson:"generation"`
}
