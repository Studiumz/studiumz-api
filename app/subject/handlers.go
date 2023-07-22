package subject

import (
	"encoding/json"
	"net/http"

	"github.com/Studiumz/studiumz-api/app"
	"github.com/Studiumz/studiumz-api/app/auth"
)

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
	_, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok { // TODO
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
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
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
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
	user, ok := r.Context().Value(auth.UserInfoCtx).(auth.User)
	if !ok {
		app.WriteHttpError(w, http.StatusUnauthorized, auth.ErrInvalidAccessToken)
		return
	}

	subjects, err := getSubjectsOfUser(user)
	if err != nil {
		app.WriteHttpError(w, http.StatusInternalServerError, err)
		return
	}
	app.WriteHttpBodyJson(w, http.StatusOK, map[string]interface{}{"data": subjects})
}
