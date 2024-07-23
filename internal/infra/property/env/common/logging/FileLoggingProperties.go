package logging

import "gurms/internal/infra/logging/core/model"

const DEFAULT_VALUE_ENABLED = true

var DEFAULT_VALUE_LEVEL model.LogLevel = "INFO"

const DEFAULT_VALUE_FILE_PATH = "@HOME/log/.log"
const DEFAULT_VALUE_MAX_FILES = 320
const DEFAULT_VALUE_FILE_SIZE_MB = 32

type FileLoggingProperties struct {
	enabled       bool
	level         model.LogLevel
	filePath      string
	maxFiles      int
	maxFileSizeMb int
}

func NewFileLoggingPropertiesDefault() *FileLoggingProperties {
	return &FileLoggingProperties{
		enabled:       DEFAULT_VALUE_ENABLED,
		level:         DEFAULT_VALUE_LEVEL,
		filePath:      DEFAULT_VALUE_FILE_PATH,
		maxFiles:      DEFAULT_VALUE_MAX_FILES,
		maxFileSizeMb: DEFAULT_VALUE_FILE_SIZE_MB,
	}
}

func (f *FileLoggingProperties) IsEnabled() bool {
	return f.enabled
}
