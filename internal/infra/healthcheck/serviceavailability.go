package healthcheck

const (
	STATUS_AVAILABLE Status = iota
	STATUS_SHUTTING_DOWN
	STATUS_INACTIVE
	STATUS_INSUFFICIENT_MEMORY
	STATUS_HIGH_CPU_USAGE
)

type Status int

type ServiceAvailability struct {
	status Status
	reason string
}

var (
	AVAILABLE     ServiceAvailability = ServiceAvailability{status: STATUS_AVAILABLE, reason: ""}
	SHUTTING_DOWN ServiceAvailability = ServiceAvailability{status: STATUS_SHUTTING_DOWN,
		reason: "the server is shutting down"}
	INACTIVE ServiceAvailability = ServiceAvailability{status: STATUS_INACTIVE,
		reason: "the server node is inactive"}
)
