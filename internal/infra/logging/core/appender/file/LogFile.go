package file

import "time"

type LogFile struct {
	path        string
	archivePath string
	dateTime    time.Time
	index       int64
}

func NewLogFile(path string, archivePath string, dateTime time.Time, index int64) *LogFile {
	return &LogFile{
		path:        path,
		archivePath: archivePath,
		dateTime:    dateTime,
		index:       index,
	}
}
