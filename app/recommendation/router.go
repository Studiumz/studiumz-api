package recommendation

import (
	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(auth.UserAuthMiddleware)

	// retrieve
	r.Get("/", showRecommendations)

	return r
}
