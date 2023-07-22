package auth

import "github.com/go-chi/chi/v5"

func SubjectRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(UserAuthMiddleware)

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
