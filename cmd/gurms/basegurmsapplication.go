package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"gurms/internal/infra/cluster"
	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/cluster/node/nodetype"
	_ "gurms/internal/infra/collection"
	_ "gurms/internal/infra/lang"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/constants"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type BaseApplication struct {
	nodeType nodetype.NodeType
}

var logger zerolog.Logger
var properties *property.GurmsProperties

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
	properties = property.InitGurmsProperties()

	var nodeId string = node.InitNodeId(viper.GetString(constants.GURMS_CLUSTER_NODE_ID))
	factory.Loggerfactory(false, nodeId, app.nodeType, properties.Logging)

	// iMongoCollectionInitializer
	node := cluster.NewNode()
	propertiesmanager := property.NewGurmsPropertiesManager(node, properties)
	//metricsconfig

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
