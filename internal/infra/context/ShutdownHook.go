package context

import "time"

type Shutdownhook interface {
	Run(timeout time.Duration) error
}
