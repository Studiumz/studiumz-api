package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Studiumz/studiumz-api/app"
)

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserInfoCtx).(User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidAccessToken)
		return
	}

	app.WriteHttpBodyJson(w, http.StatusOK, user)
}

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

func onboarding(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserInfoCtx).(User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidAccessToken)
		return
	}

	var body finishOnboardingReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, errors.New("invalid body"))
		return
	}

	err = finishOnboarding(user, body)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, errors.New("failed to finish onboarding"))
		return
	}

	app.WriteHttpBodyJson(w, http.StatusCreated, map[string]any{"message": "onboarding complete"})
}

func listSubjects(w http.ResponseWriter, r *http.Request) {
	subjects, err := getAllSubjects()
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, err)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]interface{}{"data": subjects})
}

// ADMIN ONLY
func createSubjects(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(UserInfoCtx).(User)
	if !ok { // TODO
		app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidAccessToken)
		return
	}

	var body addSubjectsReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	err = AddSubjects(body.SubjectNames)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, err)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "Added subjects succesfully"})
}

func addSubjectToSelf(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserInfoCtx).(User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidAccessToken)
		return
	}

	var body addSubjectToSelfReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		app.WriteHttpError(w, http.StatusBadRequest, err)
		return
	}

	err = UserAddSubjects(body.SubjectNames, user)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, err)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]string{"message": "Added subjects succesfully"})
}

func listSelfSubjects(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserInfoCtx).(User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, ErrInvalidAccessToken)
		return
	}

	subjects, err := getSubjectsOfUser(user)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, err)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]interface{}{"data": subjects})
}
