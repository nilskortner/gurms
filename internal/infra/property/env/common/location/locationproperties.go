package location

type LocationProperties struct {
	enabled                              bool
	treatUserIpAndDeviceTypeAsUniqueUser bool
	nearbyUserRequest                    *NearbyUserRequestProperties
}

func NewLocationProperties() *LocationProperties {
	return &LocationProperties{
		enabled:                              true,
		treatUserIpAndDeviceTypeAsUniqueUser: false,
		nearbyUserRequest:                    NewNearbyUserRequestProperties(),
	}
}
