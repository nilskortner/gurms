package logging

type LoggingProperties struct {
	console *ConsoleLoggingProperties
	file    *FileLoggingProperties
}

func NewLoggingProperties(console *ConsoleLoggingProperties, file *FileLoggingProperties) *LoggingProperties {
	return &LoggingProperties{
		console: console,
		file:    file,
	}
}

func (properties *LoggingProperties) GetConsole() *ConsoleLoggingProperties {
	return properties.console
}

func (properties *LoggingProperties) GetFile() *FileLoggingProperties {
	return properties.file
}
