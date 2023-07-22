package match

import "github.com/go-chi/chi/v5"

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/{matchee_id}", newMatch)
	r.Post("/accept/{match_id}", acceptIncoming)
	r.Post("/reject/{match_id}", rejectIncoming)
	r.Post("/withdraw/{match_id}", withdrawOutgoing)

	return r
}
