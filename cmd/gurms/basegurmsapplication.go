package main

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	_ "gurms/internal/infra/collection"
	_ "gurms/internal/infra/lang"
)

var logger = logrus.New()

// BaseApplication holds common setup and initialization
type BaseApplication struct{}

// NewBaseApplication creates a new BaseApplication
func NewBaseApplication() *BaseApplication {
	return &BaseApplication{}
}

// Init sets up the environment and logging
func (app *BaseApplication) Init() {
	app.setDefaults()
	app.validateEnv()
	app.setupLogging()
}

// setDefaults sets default environment variables
func (app *BaseApplication) setDefaults() {
	os.Setenv("TZ", "UTC")
	time.Local = time.UTC
	// Set other system properties as needed
}

// validateEnv checks for required environment configurations
func (app *BaseApplication) validateEnv() {
}

// setupLogging initializes the logger
func (app *BaseApplication) setupLogging() {
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
}

// Run starts the application
func (app *BaseApplication) Run() {
	log.Println("Running BaseApplication")
}

// handleShutdown ensures graceful shutdown
func (app *BaseApplication) handleShutdown() {
	if r := recover(); r != nil {
		logger.Errorf("Application panicked: %v", r)
		os.Exit(1)
	}
	logger.Info("Application shutting down gracefully")
}
