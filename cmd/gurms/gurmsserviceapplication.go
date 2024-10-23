package main

import (
	"gurms/internal/infra/cluster/node/nodetype"
	"log"
)

type GurmsServiceApplication struct {
	*BaseApplication
}

// NewTurmsServiceApplication creates a new TurmsServiceApplication
func NewGurmsServiceApplication() *GurmsServiceApplication {
	baseApp := NewBaseApplication(nodetype.SERVICE)
	return &GurmsServiceApplication{BaseApplication: baseApp}
}

// Run overrides the Run method to provide specific functionality
func (app *GurmsServiceApplication) run() {
	app.BaseApplication.Run()
	log.Println("Running GurmsServiceApplication")
}

func main() {
	app := NewGurmsServiceApplication()

	app.run()
}
