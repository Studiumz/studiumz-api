package main

import (
	"github.com/Studiumz/studiumz-api/app"
)

func main() {
	c := app.LoadConfig()

	// App Configurations
	app.ConfigureLogger(c)
}
