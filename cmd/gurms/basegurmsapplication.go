package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/cluster/node/nodetype"
	_ "gurms/internal/infra/collection"
	_ "gurms/internal/infra/lang"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/env/common/logging"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type BaseApplication struct {
	nodeType nodetype.NodeType
}

var logger zerolog.Logger

func NewBaseApplication(nodeType nodetype.NodeType) *BaseApplication {
	return &BaseApplication{
		nodeType: nodeType,
	}
}

// Init sets up the environment and logging
func (app *BaseApplication) init() {
	readConfigFile()

	app.setDefaults()
	app.setupLogging()
	app.initEnv()
}

func readConfigFile() {
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("configs")        // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

// setDefaults sets default environment variables
func (app *BaseApplication) setDefaults() {
	os.Setenv("TZ", "UTC")
	time.Local = time.UTC
	// Set other system properties as needed
}

// setupLogging initializes the logger
func (app *BaseApplication) setupLogging() {
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}

func (app *BaseApplication) initEnv() {
	// AsyncLoggerFactoryInit
	var nodeId string = node.InitNodeId(viper.GetString(property.GURMS_CLUSTER_NODE_ID))
	loggingProperties := configureContextForLogging()
	factory.Loggerfactory(false, nodeId, app.nodeType, loggingProperties)

	property.InitGurmsProperties()
}

// Run starts the application
func (app *BaseApplication) Run() {
	log.Println("Running BaseApplication")

	defer app.handlePanic()

	panic("test")

}

// handleShutdown ensures graceful shutdown
func (app *BaseApplication) handlePanic() {
	if r := recover(); r != nil {
		stack := make([]byte, 1024)
		runtime.Stack(stack, false)
		logger.Error().
			Str("time", time.Now().Format(time.RFC3339)).
			Interface("panic", r).
			Str("stack", string(stack)).
			Msg("Recovered from panic")
		logger.Info().Msg("Application shutting down gracefully")
		os.Exit(1)
	}
}

func configureContextForLogging() *logging.LoggingProperties {
	var enabled bool
	if viper.IsSet(property.GURMS_LOGGING_CONSOLE_ENABLED) {
		enabled = viper.GetBool(property.GURMS_LOGGING_CONSOLE_ENABLED)
	} else {
		enabled = logging.CONSOLE_DEFAULT_VALUE_ENABLED
	}
	var level loglevel.LogLevel
	if viper.IsSet(property.GURMS_LOGGING_CONSOLE_LEVEL) {
		level = loglevel.LogLevel(viper.GetInt(property.GURMS_LOGGING_CONSOLE_LEVEL))
	} else {
		level = logging.CONSOLE_DEFAULT_VALUE_LEVEL
	}
	consoleLoggingProperties := logging.NewConsoleLoggingProperties(enabled, level)

	if viper.IsSet(property.GURMS_LOGGING_FILE_ENABLED) {
		enabled = viper.GetBool(property.GURMS_LOGGING_FILE_ENABLED)
	} else {
		enabled = logging.FILE_DEFAULT_VALUE_ENABLED
	}
	if viper.IsSet(property.GURMS_LOGGING_FILE_LEVEL) {
		level = loglevel.LogLevel(viper.GetInt(property.GURMS_LOGGING_FILE_LEVEL))
	} else {
		level = logging.FILE_DEFAULT_VALUE_LEVEL
	}
	var path string
	if viper.IsSet(property.GURMS_LOGGING_FILE_PATH) {
		path = viper.GetString(property.GURMS_LOGGING_FILE_PATH)
	} else {
		path = logging.FILE_DEFAULT_VALUE_FILE_PATH
	}
	var maxFiles int
	if viper.IsSet(property.GURMS_LOGGING_FILE_MAX_FILES) {
		maxFiles = viper.GetInt(property.GURMS_LOGGING_FILE_MAX_FILES)
	} else {
		maxFiles = logging.FILE_DEFAULT_VALUE_MAX_FILES
	}
	var fileSizeMb int
	if viper.IsSet(property.GURMS_LOGGING_FILE_MAX_FILE_SIZE_MB) {
		fileSizeMb = viper.GetInt(property.GURMS_LOGGING_FILE_MAX_FILE_SIZE_MB)
	} else {
		fileSizeMb = logging.FILE_DEFAULT_VALUE_FILE_SIZE_MB
	}
	var compression bool
	if viper.IsSet(property.GURMS_LOGGING_FILE_COMPRESSION_ENABLED) {
		compression = viper.GetBool(property.GURMS_LOGGING_FILE_COMPRESSION_ENABLED)
	} else {
		compression = logging.FILE_DEFAULT_VALUE_COMPRESSION_ENABLED
	}
	fileLoggingProperties := logging.NewFileLoggingProperties(enabled, level, path, maxFiles, fileSizeMb, compression)

	return logging.NewLoggingProperties(consoleLoggingProperties, fileLoggingProperties)
}
