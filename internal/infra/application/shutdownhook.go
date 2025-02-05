package application

type ShutdownHook interface {
	Run(timeoutMillis int64) error
}
