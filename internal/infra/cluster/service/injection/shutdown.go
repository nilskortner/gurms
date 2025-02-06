package injection

import "context"

type ShutDown interface {
	AddClosingContext(ctxClose context.CancelFunc)
	AddShutdownChannel(shutdown chan struct{})
}
