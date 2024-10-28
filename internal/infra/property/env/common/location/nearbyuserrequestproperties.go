package location

type NearbyUserRequestProperties struct {
	defautlMaxNearbyUserCount int16
	maxNearbyUserCount        int16
	defaultMaxDistanceMeters  int
	maxDistanceMeter          int
}

func NewNearbyUserRequestProperties() *NearbyUserRequestProperties {
	return &NearbyUserRequestProperties{
		defautlMaxNearbyUserCount: 20,
		maxNearbyUserCount:        100,
		defaultMaxDistanceMeters:  10_000,
		maxDistanceMeter:          10_000,
	}
}
