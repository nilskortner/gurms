package file

import (
	"gurms/internal/supportpkgs/datastructures/dequeue"
	"gurms/internal/supportpkgs/mathsupport"
	"gurms/internal/supportpkgs/timesupport"
)

type LogDirectoryVisitor struct {
	filePrefix            string
	fileSuffix            string
	fileMiddle            string
	fileDateTimeFormatter *timesupport.DateTimeFormatter
	maxFilesToKeep        int
	deleteExceedFiles     bool
}

func NewLogDirectoryVisitor(
	filePrefix string,
	fileSuffix string,
	fileMiddle string,
	fileDateTimeFormatter *timesupport.DateTimeFormatter,
	maxFiles int) *LogDirectoryVisitor {
	return &LogDirectoryVisitor{
		filePrefix:            filePrefix,
		fileSuffix:            fileSuffix,
		fileMiddle:            fileMiddle,
		fileDateTimeFormatter: fileDateTimeFormatter,
		maxFilesToKeep:        mathsupport.Max(1, maxFiles),
		deleteExceedFiles:     maxFiles > 0,
	}
}

func Visit(
	directory string,
	prefix string,
	suffix string,
	middle string,
	template *timesupport.DateTimeFormatter,
	maxFiles int) *dequeue.Dequeue {
	visitor := NewLogDirectoryVisitor(prefix, suffix, middle, template, maxFiles)

}
