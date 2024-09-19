package bufferpool

import (
	"bytes"
	"sync"
)

var BufferPool sync.Pool

func init() {
	BufferPool = sync.Pool{
		New: func() interface{} {
			var buffer bytes.Buffer
			return buffer
		},
	}
}
