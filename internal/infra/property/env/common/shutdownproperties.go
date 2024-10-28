package common

type ShutdownProperties struct {
	jobTimeoutMillis int
}

func NewShutdonwProperties() *ShutdownProperties {
	return &ShutdownProperties{
		jobTimeoutMillis: 120 * 1000,
	}
}
