package match

import "github.com/go-chi/chi/v5"

func Router() *chi.Mux {
	r := chi.NewRouter()

	// retrieve
	r.Get("/incoming", getIncoming)
	r.Get("/outgoing", getOutgoing)

	// for creating new match/skips
	r.Post("/skip/{matchee_id}", newSkip)
	r.Post("/connect/{matchee_id}", newMatch)

	// to take action on existing matches
	r.Post("/accept/{match_id}", acceptIncoming)
	r.Post("/reject/{match_id}", rejectIncoming)
	r.Post("/withdraw/{match_id}", withdrawOutgoing)

	return r
}
