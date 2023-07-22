package match

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Studiumz/studiumz-api/app"
	"github.com/Studiumz/studiumz-api/app/auth"
	"github.com/go-chi/chi/v5"
)

func getIncoming(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	matches, err := getUserIncoming(user.Id)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, errors.New("could not get your incoming requests"))
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]interface{}{"data": matches})
}

func getOutgoing(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	matches, err := getUserOutgoing(user.Id)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, errors.New("could not get your incoming requests"))
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]interface{}{"data": matches})
}

func newMatch(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	var body CreateMatchReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, errors.New("invalid body"))
		return
	}

	matcheeId := chi.URLParam(r, "matchee_id")
	err = createNewMatch(matcheeId, body, user)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, ErrFailToUpdateMatch)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "Match request accepted"})
}

func newSkip(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	matcheeId := chi.URLParam(r, "matchee_id")
	err := createNewSkip(matcheeId, user)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, ErrFailToUpdateMatch)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "skip success"})
}

func acceptIncoming(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	matchId := chi.URLParam(r, "match_id")
	m, err := GetMatchById(matchId)
	if err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, ErrMatchNotFound)
		return
	}

	if user.Id != m.MatcheeId {
		app.WriteHttpError(w, http.StatusForbidden, auth.ErrInvalidUserId)
		return
	}

	err = acceptMatch(m)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, ErrFailToUpdateMatch)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "Match request accepted"})
}

func rejectIncoming(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	matchId := chi.URLParam(r, "match_id")
	m, err := GetMatchById(matchId)
	if err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, ErrMatchNotFound)
		return
	}

	if user.Id != m.MatcheeId {
		app.WriteHttpError(w, http.StatusForbidden, auth.ErrInvalidUserId)
		return
	}

	err = rejectMatch(m)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, ErrFailToUpdateMatch)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "Match request rejected"})
}

func withdrawOutgoing(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	matchId := chi.URLParam(r, "match_id")
	m, err := GetMatchById(matchId)
	if err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, ErrMatchNotFound)
		return
	}

	if user.Id != m.MatcherId {
		app.WriteHttpError(w, http.StatusForbidden, auth.ErrInvalidUserId)
		return
	}

	err = withdrawMatch(m)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, ErrFailToUpdateMatch)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "Match request withdrawn"})
}
