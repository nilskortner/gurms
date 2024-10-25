package cluster

type GrpcProperties struct {
	jobTimeoutMillis int64
}

func NewGrpcProperties() *GrpcProperties {
	return &GrpcProperties{
		jobTimeoutMillis: 120 * 1000,
	}
}
