package file

import (
	"fmt"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/appender/file/logfile"
	"gurms/internal/infra/logging/core/compression"
	"gurms/internal/infra/logging/core/model"
	"gurms/internal/infra/time"
	"gurms/internal/supportpkgs/datastructures/dequeue"
	"gurms/internal/supportpkgs/mathsupport"
	"gurms/internal/supportpkgs/timesupport"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const FIELD_DELIMITER string = "_"
const FILE_MIDDLE = "20060102"

const ARCHIVE_FILE_SUFFIX = ".gz"

const COMPRESSION_LEVEL = 3

var READ_OPTIONS = []int{os.O_RDONLY}
var CREATE_NEW_OPTIONS = []int{os.O_CREATE | os.O_WRONLY | os.O_EXCL}

var APPEND_OPTIONS = []int{os.O_CREATE | os.O_WRONLY | os.O_APPEND}

type RollingFileAppender struct {
	channelAppender       *appender.ChannelAppender
	filePrefix            string
	fileSuffix            string
	fileDirectory         string
	fileDirectoryFile     fs.FileInfo
	fileDateTimeFormatter *timesupport.DateTimeFormatter
	maxFiles              int
	maxFilesBytes         int64
	minUsableSpaceBytes   int64
	files                 dequeue.Dequeue
	currentFile           logfile.LogFile
	nextFileBytes         int64
	nextIndex             int64
	nextDay               int64
	enableCompression     bool
	gzipOutputStream      *compression.FastGzipOutputStream
}

func NewRollingFileAppender(
	level model.LogLevel,
	file string,
	maxFiles int,
	maxFileMb int64,
	enableCompression bool,
) *RollingFileAppender {
	filePath, err := filepath.Abs(file)
	if err != nil {
		fmt.Println(err)
	}
	fileName := filepath.Base(filePath)
	index := strings.LastIndex(fileName, ".")

	var filePrefix string
	var fileSuffix string
	if index == -1 {
		filePrefix = fileName
		fileSuffix = ""
	} else {
		filePrefix = fileName[:index]
		fileSuffix = fileName[index:]
	}
	fileDirectory := filepath.Dir(filePath)
	fileDirectoryFile, err := os.Stat(fileDirectory)
	if err != nil {
		fmt.Println(err)
	}
	maxFiles = mathsupport.Max(maxFiles, 0)
	var maxFileBytes int64
	if maxFileMb > 0 {
		maxFileBytes = maxFileMb * 1024
	} else {
		maxFileBytes = 1024 * 1024 * 1024
	}
	minUsableSpaceBytes := int64(float64(maxFileBytes) * 2.5)

	var gzipOutputStream *compression.FastGzipOutputStream
	if enableCompression {
		gzipOutputStream, err = compression.NewFastGzipOutputStream("a", COMPRESSION_LEVEL, int(maxFileBytes/10))
		if err != nil {
			fmt.Println(err)
		}
	}

	err = os.MkdirAll(fileDirectory, 0755)
	if err != nil {
		fmt.Println("Failed to create the directory ("+fileDirectory+")for log files", err)
	}
	fileDateTimeFormatter := timesupport.NewDateTimeFormatter(FILE_MIDDLE, time.ZONE_ID)
	var files *dequeue.Dequeue
	files, err = Visit(fileDirectory, filePrefix, fileSuffix, FILE_MIDDLE, fileDateTimeFormatter, maxFiles)
	if err != nil {
		fmt.Println(err)
	}

	//logFile := files.peekLast()

	return &RollingFileAppender{
		channelAppender:       appender.NewChannelAppender(level),
		filePrefix:            filePrefix,
		fileSuffix:            fileSuffix,
		fileDirectory:         fileDirectory,
		fileDirectoryFile:     fileDirectoryFile,
		maxFiles:              maxFiles,
		maxFilesBytes:         maxFileBytes,
		minUsableSpaceBytes:   minUsableSpaceBytes,
		fileDateTimeFormatter: fileDateTimeFormatter,
		nextDay:               math.MinInt64,
		enableCompression:     enableCompression,
		gzipOutputStream:      gzipOutputStream,
	}
}
