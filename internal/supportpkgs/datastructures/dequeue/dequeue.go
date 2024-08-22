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

func (d *Dequeue) PeekLast() (logfile.LogFile, bool) {
	if len(d.data) > 0 {
		return d.data[len(d.data)-1], true
	}
	return logfile.LogFile{}, false
}

func (d *Dequeue) Add(file logfile.LogFile) {
	d.data = append(d.data, file)
}
