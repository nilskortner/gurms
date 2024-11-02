package cluster

type RpcProperties struct {
	JobTimeoutMillis int
}

func NewRpcProperties() *RpcProperties {
	return &RpcProperties{
		JobTimeoutMillis: 120 * 1000,
	}
}
