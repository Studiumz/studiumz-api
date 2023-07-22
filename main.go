package main

import (
	"net/http"

	"github.com/Studiumz/studiumz-api/app"
	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/Studiumz/studiumz-api/app/recommendation"
	"github.com/Studiumz/studiumz-api/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func main() {
	c := app.LoadConfig()

	// App Configurations
	app.ConfigureLogger(c)
	app.ConfigureCors(c)
	co := app.ConfigureCohere(c)

	recommendation.InjectCohereClientAdapter(co)

	// Configure Adapters and Dependency Injection
	pool := db.CreateConnPool(c.DbDsn)

	auth.SetEnv(c.Env)
	auth.SetPool(pool)
	auth.ConfigureFirebaseAdminSdk(c.FirebaseAdminServiceAccountFile)
	auth.ConfigureJWTProperties(c.StudiumzJwtIssuer, c.StudiumzJwtAccessTokenSecret)

	r := chi.NewRouter()

	// Global middlewares
	r.Use(app.ReqLoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(app.CorsMiddleware)

	// Default route handlers
	r.NotFound(app.NotFound)
	r.MethodNotAllowed(app.MethodNotAllowed)
	r.Get("/", app.Heartbeat)

	// Normal routes
	r.Group(func(r chi.Router) {
		r.Mount("/auth", auth.Router())
	})

	log.Info().Msgf("Running server on port %s in %s mode...", c.Port, c.Env)
	http.ListenAndServe(":"+c.Port, r)
}
