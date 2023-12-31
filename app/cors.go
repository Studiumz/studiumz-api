package app

import (
	"net/http"

	"github.com/rs/cors"
)

var allowedOrigins []string
var allowedMethods []string
var allowedHeaders []string

func ConfigureCors(c Config) {
	allowedOrigins = []string{c.ClientAppUrl}
	allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	if c.Env == "local" || c.Env == "development" {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000")
		allowedMethods = append(allowedMethods, "HEAD", "TRACE")
	}

	allowedHeaders = []string{"Accept", "Authorization", "X-Forwarded-Authorization", "Content-Type", "X-Google-Id-Token", "X-Studiumz-Api-Key"}
}

func CorsMiddleware(h http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: allowedMethods,
		AllowedHeaders: allowedHeaders,
	}).Handler(h)
}
