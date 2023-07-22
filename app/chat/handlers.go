package chat

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Studiumz/studiumz-api/app"
	"github.com/Studiumz/studiumz-api/app/auth"
)

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
