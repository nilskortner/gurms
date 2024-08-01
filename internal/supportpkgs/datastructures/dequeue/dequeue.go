package dequeue

import (
	"gurms/internal/infra/logging/core/appender/file/logfile"
)

type Dequeue struct {
	data []logfile.LogFile
}

func NewDequeue(files []logfile.LogFile) *Dequeue {
	return &Dequeue{data: files}
}
