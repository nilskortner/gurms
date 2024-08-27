package logger

import (
	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/appender/file"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/processor"
	"gurms/internal/infra/property/env/common/logging"
	"gurms/internal/infra/system"
	"strings"

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

var loggerlayout *layout.GurmsTemplateLayout

var initialized bool

var ALL_APPENDERS copyonwriteslice.CopyOnWriteSliceAppender
var DEFAULT_APPENDERS = make([]appender.Appender, 0, 2)
var Queue *mpscunboundedarrayqueue.MpscUnboundedArrayQueue
var UNINITIALIZED_LOGGERS linkedlist.LinkedList

var homeDir string
var serverTypeName string
var fileLoggingProperties logging.FileLoggingProperties
var defaultConsoleAppender appender.Appender

var logprocessor processor.LogProcessor

// use sync once?
func initialize(
	runWithTests bool,
	nodeType node.NodeType,
	nodeId string,
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
	fileLoggingProperties := properties.GetFile()
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
		DEFAULT_APPENDERS[0] = fileAppender
	}

	loggerlayout = layout.NewGurmsTemplateLayout(nodeType, nodeId)
	initialized = true

	InternalLogger.init()

	processor := processor.NewLogProcessor(Queue)
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
