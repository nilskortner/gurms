package logger

import (
	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/processor"
	"gurms/internal/infra/property/env/common/logging"
	"gurms/internal/infra/system"

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

var loggerlayout layout.GurmsTemplateLayout

var initialized bool

var ALL_APPENDERS copyonwriteslice.CopyOnWriteSliceAppender
var DEFAULT_APPENDERS [2]appender.Appender
var Queue mpscunboundedarrayqueue.MpscUnboundedArrayQueue
var UNINITIALIZED_LOGGERS linkedlist.LinkedList

var homeDir string
var serverTypeName string
var fileLoggingProperties logging.FileLoggingProperties
var defaultConsoleAppender appender.Appender

var logprocessor processor.LogProcessor

func initialize(
	runWithTests bool,
	nodeType node.NodeType,
	nodeId string,
	properties logging.Loggingproperties) {
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
}
