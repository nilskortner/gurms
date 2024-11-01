package service

type ClusterService interface {
	Start()
	Stop(timoutMillis int64)
	LazyInit()
}
