package option

import "go.mongodb.org/mongo-driver/v2/bson"

type Update struct {
	Update bson.M
}

func NewUpdate(memberHost, adminApiAddress, wsAddress, tcpAddress, udpAddress any) *Update {
	updateMap := make(bson.M, 0)
	if memberHost != "" {
		updateMap["memberHost"] = memberHost
	}
	if adminApiAddress != "" {
		updateMap["adminApiAddress"] = adminApiAddress
	}
	if wsAddress != "" {
		updateMap["wsAddress"] = wsAddress
	}
	if tcpAddress != "" {
		updateMap["tcpAddress"] = tcpAddress
	}
	if udpAddress != "" {
		updateMap["udpAddress"] = udpAddress
	}
	return &Update{
		Update: updateMap,
	}
}
