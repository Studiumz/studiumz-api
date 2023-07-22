package subject

import (
	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(auth.UserAuthMiddleware)

	// retrieve all subjects
	r.Get("/", listSubjects)

	// create new
	r.Post("/", createSubjects)

	// create new subject_users
	r.Post("/self", addSubjectToSelf) // for manual use

	// retreive user's subjects
	r.Get("/self", listSelfSubjects)

	return r
}
