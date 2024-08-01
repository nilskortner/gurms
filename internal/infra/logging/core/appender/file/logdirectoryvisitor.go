package file

import (
	"fmt"
	"gurms/internal/supportpkgs/datastructures/dequeue"
	"gurms/internal/supportpkgs/datastructures/treeset"
	"gurms/internal/supportpkgs/mathsupport"
	"gurms/internal/supportpkgs/timesupport"
	"os"
	"path/filepath"
	"strings"
)

type LogDirectoryVisitor struct {
	filePrefix            string
	fileSuffix            string
	fileMiddle            string
	fileDateTimeFormatter *timesupport.DateTimeFormatter
	maxFilesToKeep        int
	deleteExceedFiles     bool
	files                 *treeset.Tree
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
		files:                 treeset.New(treeset.LogComparator),
	}
}

func Visit(
	directory string,
	prefix string,
	suffix string,
	middle string,
	template *timesupport.DateTimeFormatter,
	maxFiles int) (*dequeue.Dequeue, error) {
	visitor := NewLogDirectoryVisitor(prefix, suffix, middle, template, maxFiles)
	maxDepth := 1
	err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to visit the directory: %v, %v", path, err)
		}
		depth := strings.Count(path, string(os.PathSeparator)) - strings.Count(directory, string(os.PathSeparator))

		if depth > maxDepth {
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dequeue.NewDequeue(visitor.files.Keys()), nil
}
