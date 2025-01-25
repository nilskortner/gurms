package aiserving

const MB = 1024 * 1024

type AdminHttpProperties struct {
	host                      string
	port                      int
	connectTimeoutMillis      int
	idleTimeoutMillis         int
	requestReadTimeoutMillis  int
	maxRequestBodySizeRequest int
	ssl                       *SslProperties
}

func NewAdminHttpProperties() *AdminHttpProperties {
	return &AdminHttpProperties{
		host:                      "0.0.0.0",
		port:                      -1,
		connectTimeoutMillis:      30 * 1000,
		idleTimeoutMillis:         3 * 60 * 1000,
		requestReadTimeoutMillis:  3 * 60 * 1000,
		maxRequestBodySizeRequest: 10 * MB,
		ssl:                       NewSslProperties(),
	}
}
