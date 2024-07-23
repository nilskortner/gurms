package timesupport

import "time"

type DateTimeFormatter struct {
	fileMiddle string
	zone       *time.Location
}

func NewDateTimeFormatter(fileMiddle string, zone *time.Location) *DateTimeFormatter {
	return &DateTimeFormatter{
		fileMiddle: fileMiddle,
		zone:       zone,
	}
}
