package main

import (
	"github.com/go-martini/martini"
)

func main() {

	// Initializing all the Martini goodness

	app.M = martini.Classic()

	// Set up the environment (eg. load toml config, load database, etc.)

	if err := app.SetupEnvironment(); err != nil {
		panic(err)
	}

	// Pass true if you would like `SetupControllers` to also "run" martini for you.
	// If you need to do additional work after setting up controllers but prior to
	// starting the web server, you would set false
	app.SetupControllers(true)
}
