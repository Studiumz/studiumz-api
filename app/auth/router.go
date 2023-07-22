package auth

import "github.com/go-chi/chi/v5"

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/sign-in/google", signInWithGoogleHandler)

	if env == "local" || env == "development" {
		r.Get("/test/firebase-id-token", getTestFirebaseIdTokenHandler)
	}

	return r
}
