package treeset

import "gurms/internal/infra/logging/core/appender/file/logfile"

type Comparator func(x, y logfile.LogFile) int

func LogComparator(x, y logfile.LogFile) int {
	switch {
	case x.GetIndex() > y.GetIndex():
		return 1
	case x.GetIndex() < y.GetIndex():
		return -1
	default:
		return 0
	}
}
