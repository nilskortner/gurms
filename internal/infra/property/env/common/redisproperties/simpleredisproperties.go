package redisproperties

type SimpleRedisProperties struct {
	uri string
}

func NewSimpleRedisProperties() *SimpleRedisProperties {
	return &SimpleRedisProperties{
		uri: "redis://localhost",
	}
}
