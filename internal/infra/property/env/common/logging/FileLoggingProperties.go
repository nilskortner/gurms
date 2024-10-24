package logging

import (
	"gurms/internal/infra/logging/core/model/loglevel"
)

const FILE_DEFAULT_VALUE_ENABLED = true
const FILE_DEFAULT_VALUE_LEVEL loglevel.LogLevel = 1

const FILE_DEFAULT_VALUE_FILE_PATH = "@HOME/log/.log"
const FILE_DEFAULT_VALUE_MAX_FILES = 320
const FILE_DEFAULT_VALUE_FILE_SIZE_MB = 32

type FileLoggingProperties struct {
	enabled       bool
	level         loglevel.LogLevel
	filePath      string
	maxFiles      int
	maxFileSizeMb int
	compression   bool
}

func NewFileLoggingProperties(
	enabled bool,
	level loglevel.LogLevel,
	filePath string,
	maxFiles int,
	maxFileSizeMb int,
	compression bool) *FileLoggingProperties {
	return &FileLoggingProperties{
		enabled:       enabled,
		level:         level,
		filePath:      filePath,
		maxFiles:      maxFiles,
		maxFileSizeMb: maxFileSizeMb,
		compression:   compression,
	}
}

func (f *FileLoggingProperties) IsEnabled() bool {
	return f.enabled
}

func (f *FileLoggingProperties) GetLevel() loglevel.LogLevel {
	return f.level
}

func (f *FileLoggingProperties) GetFilePath() string {
	return f.filePath
}

func (f *FileLoggingProperties) GetMaxFiles() int {
	return f.maxFiles
}

func (f *FileLoggingProperties) GetMaxFilesSizeMb() int {
	return f.maxFileSizeMb
}

func (f *FileLoggingProperties) GetCompression() bool {
	return f.compression
}
