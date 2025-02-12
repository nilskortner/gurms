package injection

type ShutDown interface {
	AddShutdownChannel(shutdown chan struct{})
	AddShutdownFunction(shutdown func())
}
