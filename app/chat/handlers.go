package chat

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Studiumz/studiumz-api/app"
	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/go-chi/chi/v5"
)

func createMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	userId := user.Id
	chatId := chi.URLParam(r, "chatId")

	var body createMessageReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	message, errs, err := createMessage(ctx, userId, chatId, body.Text)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch {
		case errors.As(err, &ErrInvalidChatId):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrChatDoesNotExist):
			app.WriteHttpError(w, http.StatusNotFound, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, message)
}

func getChatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	userId := user.Id

	chats, err := getChats(ctx, userId)
	if err != nil {
		switch err {
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, chats)
}

func createChatHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	firstUserId := user.Id

	var body createChatReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	chat, errs, err := createChat(ctx, firstUserId, body)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch {
		case errors.As(err, &ErrInvalidChatUserId):
			app.WriteHttpError(w, http.StatusBadRequest, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, chat)
}
