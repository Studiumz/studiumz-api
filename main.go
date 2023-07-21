package main

import (
	"github.com/Studiumz/studiumz-api/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	c := app.LoadConfig()

	// App Configurations
	app.ConfigureLogger(c)
	app.ConfigureCors(c)

	r := chi.NewRouter()

	// Global middlewares
	r.Use(app.ReqLoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(app.CorsMiddleware)
}
