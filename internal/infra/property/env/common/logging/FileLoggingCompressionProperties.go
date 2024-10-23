package logging

type FileLoggingCompressionProperties struct {
	enabled bool
}

func NewFileLoggingCompressionProperties() *FileLoggingCompressionProperties {
	return &FileLoggingCompressionProperties{
		enabled: FILE_DEFAULT_VALUE_ENABLED,
	}
}

func (f FileLoggingCompressionProperties) IsEnabled() bool {
	return f.enabled
}
