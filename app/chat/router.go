package chat

import (
	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(auth.UserAuthMiddleware)

	r.Get("/", getChatsHandler)
	r.Post("/create", createChatHandler)
	r.Get("/{chatId}", getChatDetailHandler)
	r.Post("/{chatId}/create", createMessageHandler)

	return r
}
