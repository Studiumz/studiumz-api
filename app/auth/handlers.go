package auth

import (
	"net/http"

	"github.com/Studiumz/studiumz-api/app"
)

func getTestFirebaseIdTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idToken, err := createTestFirebaseIdToken(ctx)
	if err != nil {
		app.WriteHttpInternalServerError(w)
		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"id_token": idToken})
}

func signInWithGoogleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idToken := r.Header.Get("X-Google-Id-Token")

	signIn, errs, err := signInWithGoogle(ctx, idToken)
	if errs != nil {
		app.WriteHttpErrors(w, http.StatusBadRequest, errs)
		return
	}
	if err != nil {
		switch err {
		case ErrInvalidFirebaseIdToken:
			app.WriteHttpError(w, http.StatusUnauthorized, err)
		default:
			app.WriteHttpInternalServerError(w)
		}

		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, signIn)
}
