package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"gurms/internal/infra/address"
	"gurms/internal/infra/application"
	"gurms/internal/infra/cluster"
	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/cluster/node/nodetype"
	_ "gurms/internal/infra/collection"
	"gurms/internal/infra/healthcheck"
	_ "gurms/internal/infra/lang"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/property"
	"gurms/internal/infra/property/constants"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type BaseApplication struct {
	nodeType        nodetype.NodeType
	shutDownManager *application.ShutDownManager
	addressManager  address.ServiceAddressManager
}

var logger zerolog.Logger
var properties *property.GurmsProperties

func NewBaseApplication(nodeType nodetype.NodeType,
	addressManager address.ServiceAddressManager) *BaseApplication {
	return &BaseApplication{
		nodeType:       nodeType,
		addressManager: addressManager,
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

// setupLogging initializes the internal logger
func (app *BaseApplication) setupLogging() {
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}

func (app *BaseApplication) initEnv() {
	app.shutDownManager = application.NewShutDownManager()
	properties = property.NewGurmsProperties()

	var nodeId string = node.InitNodeId(viper.GetString(constants.GURMS_CLUSTER_NODE_ID))
	factory.Loggerfactory(false, nodeId, app.nodeType, properties.Logging)

	// iMongoCollectionInitializer
	propertiesManager := property.NewGurmsPropertiesManager(nil, properties)
	// TODO: check if lazy init and injections work properly
	healthcheckManager := healthcheck.NewHealthCheckManager(app.shutDownManager, nil, propertiesManager)
	node := cluster.NewNode(app.nodeType, app.shutDownManager,
		propertiesManager, app.addressManager, healthcheckManager)
	propertiesManager.SetNode(node)
	//metricsconfig

}

// Run starts the application
func (app *BaseApplication) Run() {
	log.Println("Running BaseApplication")
	defer app.handlePanic()
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
		if factory.IsInitialized() {
			factory.WaitClose(60 * 1000)
		}
		logger.Info().Msg("Application shutting down gracefully")
		os.Exit(1)
	}
}
