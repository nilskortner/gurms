package common

type FlightRecorderProperties struct {
	closedRecordingRetentionPeriod int
}

func NewFlightRecorderProperties() *FlightRecorderProperties {
	return &FlightRecorderProperties{
		closedRecordingRetentionPeriod: 0,
	}
}
