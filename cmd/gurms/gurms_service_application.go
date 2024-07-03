package main

import (
	"log"
)

// TurmsServiceApplication extends BaseApplication
type TurmsServiceApplication struct {
	*BaseApplication
}

// NewTurmsServiceApplication creates a new TurmsServiceApplication
func NewTurmsServiceApplication() *TurmsServiceApplication {
	baseApp := NewBaseApplication()
	return &TurmsServiceApplication{BaseApplication: baseApp}
}

// Run overrides the Run method to provide specific functionality
func (app *TurmsServiceApplication) Run() {
	app.BaseApplication.Run()
	log.Println("Running TurmsServiceApplication")
}

func main() {
	app := NewTurmsServiceApplication()
	app.Init()
	defer app.handleShutdown()

	app.Run()
}
