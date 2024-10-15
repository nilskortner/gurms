package main

import (
	"log"
)

// TurmsServiceApplication extends BaseApplication
type GurmsServiceApplication struct {
	*BaseApplication
}

// NewTurmsServiceApplication creates a new TurmsServiceApplication
func NewGurmsServiceApplication() *GurmsServiceApplication {
	baseApp := NewBaseApplication()
	return &GurmsServiceApplication{BaseApplication: baseApp}
}

// Run overrides the Run method to provide specific functionality
func (app *GurmsServiceApplication) Run() {
	app.BaseApplication.Run()
	log.Println("Running TurmsServiceApplication")
}

func main() {
	app := NewGurmsServiceApplication()
	app.Init()
	defer app.handleShutdown()

	app.Run()
}
