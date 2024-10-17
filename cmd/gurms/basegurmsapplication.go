package main

import (
	"log"
	"os"
	"runtime"
	"time"

	_ "gurms/internal/infra/collection"
	_ "gurms/internal/infra/lang"
	"gurms/internal/infra/logging/core/factory"
	"gurms/internal/infra/property"

	"github.com/rs/zerolog"
)

// BaseApplication holds common setup and initialization
type BaseApplication struct{}

var logger zerolog.Logger

// NewBaseApplication creates a new BaseApplication
func NewBaseApplication() *BaseApplication {
	return &BaseApplication{}
}

// Init sets up the environment and logging
func (app *BaseApplication) init() {
	app.setDefaults()
	app.initEnv()
	app.setupLogging()
}

// setDefaults sets default environment variables
func (app *BaseApplication) setDefaults() {
	os.Setenv("TZ", "UTC")
	time.Local = time.UTC
	// Set other system properties as needed
}

// validateEnv checks for required environment configurations
func (app *BaseApplication) initEnv() {
	var nodeId string = property.GURMS_CLUSTER_NODE_ID

	factory.Loggerfactory(false, nodeId, nodeType, loggingProperties)
}

// setupLogging initializes the logger
func (app *BaseApplication) setupLogging() {
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
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
		logger.Info().Msg("Application shutting down gracefully")
		os.Exit(1)
	}
}
