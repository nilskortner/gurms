package factory

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/appender/file"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/logger"
	"gurms/internal/infra/logging/core/model/logrecord"
	"gurms/internal/infra/property/env/common/logging"
	"gurms/internal/infra/system"
	"strings"
	"sync"

	"gurms/internal/supportpkgs/datastructures/copyonwriteslice"
	"gurms/internal/supportpkgs/datastructures/linkedlist"
	"gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
)

const (
	PROPERTY_NAME_TURMS_AI_SERVING_HOME = "GURMS_AI_SERVING_HOME"
	PROPERTY_NAME_TURMS_GATEWAY_HOME    = "GURMS_GATEWAY_HOME"
	PROPERTY_NAME_TURMS_SERVICE_HOME    = "GURMS_SERVICE_HOME"
	SERVER_TYPE_UNKNOWN                 = "unknown"
)

var once sync.Once

var loggerlayout *layout.GurmsTemplateLayout

var initialized bool

var ALL_APPENDERS = copyonwriteslice.NewCopyOnWriteSlice[appender.Appender]()
var DEFAULT_APPENDERS = make([]appender.Appender, 0, 2)
var Queue *mpscunboundedarrayqueue.MpscUnboundedArrayQueue[logrecord.LogRecord]
var UNINITIALIZED_LOGGERS linkedlist.LinkedList

var homeDir string
var serverTypeName string
var fileLoggingProperties logging.FileLoggingProperties
var defaultConsoleAppender appender.Appender

var logprocessor LogProcessor

func Loggerfactory(runWithTests bool,
	nodeId string,
	nodeType nodetype.NodeType,
	properties logging.LoggingProperties) {
	once.Do(func() {
		initialize(runWithTests, nodeId, nodeType, properties)
	})
}

func WaitClose(timeoutMillis int64) {
	logprocessor.waitClose(timeoutMillis)
}

func GetLogger(name string) logger.Logger {
	options := logger.NewLoggerOptions(name)
	return getLogger(options)
}

func initialize(
	runWithTests bool,
	nodeId string,
	nodeType nodetype.NodeType,
	properties logging.LoggingProperties) {
	if initialized {
		return
	}
	switch nodeType {
	case 0:
		homeDir = system.GetProperty("PROPERTY_NAME_GURMS_AI_SERVING_HOME")
	case 1:
		homeDir = system.GetProperty("PROPERTY_NAME_GURMS_GATEWAY_HOME")
	case 2:
		homeDir = system.GetProperty("PROPERTY_NAME_GURMS_SERVICE_HOME")
	}
	if homeDir == "" {
		homeDir = "."
	}
	serverTypeName = nodeType.GetId()
	consoleLoggingProperties := properties.GetConsole()
	fileLoggingProperties = properties.GetFile()
	if consoleLoggingProperties.IsEnabled() {
		var consoleAppender appender.Appender
		if runWithTests {
			consoleAppender = appender.NewSystemConsoleAppender(consoleLoggingProperties.Level())
		} else {
			consoleAppender = appender.NewChannelConsoleAppender(consoleLoggingProperties.Level())
		}
		defaultConsoleAppender = consoleAppender
		DEFAULT_APPENDERS = append(DEFAULT_APPENDERS, consoleAppender)
	}
	if fileLoggingProperties.IsEnabled() {
		fileAppender := file.NewRollingFileAppender(
			fileLoggingProperties.GetLevel(),
			getFilePath(fileLoggingProperties.GetFilePath()),
			fileLoggingProperties.GetMaxFiles(),
			int64(fileLoggingProperties.GetMaxFilesSizeMb()),
			fileLoggingProperties.GetCompression(),
		)
		DEFAULT_APPENDERS = append(DEFAULT_APPENDERS, fileAppender)
	}

	loggerlayout = layout.NewGurmsTemplateLayout(nodeType, nodeId)
	initialized = true

	processor := NewLogProcessor(Queue)
	processor.Start()
}

func getFilePath(path string) string {
	if path == "" {
		return "."
	}
	path = strings.Replace(path, "@HOME", homeDir, -1)
	path = strings.Replace(path, "@SERVICE_TYPE_NAME", serverTypeName, -1)
	return path
}

func getLogger(options *logger.LoggerOptions) logger.Logger {
	loggerName := options.GetName()
	filePath := options.GetPath()
	appenders := make([]appender.Appender, 2)
	if filePath != "" {
		filePath = getFilePath(filePath)
		level := options.GetLevel()
		if level == -1 {
			level = fileLoggingProperties.GetLevel()
		}
		appender := file.NewRollingFileAppender(
			level,
			filePath,
			fileLoggingProperties.GetMaxFiles(),
			int64(fileLoggingProperties.GetMaxFilesSizeMb()),
			fileLoggingProperties.GetCompression())
		appenders = append(appenders, appender)
		ALL_APPENDERS.Add(appender)
		if defaultConsoleAppender != nil {
			appenders = append(appenders, defaultConsoleAppender)
		}
	} else {
		appenders = DEFAULT_APPENDERS
	}
	return logger.NewAsyncLogger(loggerName, options.IsShouldParse(), appenders, loggerlayout, Queue)
}
