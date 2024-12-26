package configdiscovery

import "time"

type Leader struct {
	ClusterId  string    `bson:"_id"`
	NodeId     string    `bson:"node_id"`
	RenewDate  time.Time `bson:"renew_date"`
	Generation int       `bson:"generation"`
}

func NewLeader(clusterId, nodeId string, renewDate time.Time, generation int) *Leader {
	return &Leader{
		ClusterId:  clusterId,
		NodeId:     nodeId,
		RenewDate:  renewDate,
		Generation: generation,
	}
}
