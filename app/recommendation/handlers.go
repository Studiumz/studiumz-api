package recommendation

import (
	"errors"
	"net/http"

	"github.com/Studiumz/studiumz-api/app"
	"github.com/Studiumz/studiumz-api/app/auth"
)

func showRecommendations(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	recommendations, err := CreateRecommendation(user)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, errors.New("could not get your incoming requests"))
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]interface{}{"data": recommendations})
}
