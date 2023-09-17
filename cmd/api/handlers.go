package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/youngjun827/api-std-lib/internal/middleware"
)

func (app *application) CreateUser(w http.ResponseWriter, r *http.Request) {
    user, err := decodeAndValidateUser(r, w)
    if err != nil {
        middleware.JSONError(w, err, http.StatusBadRequest)
        return
    }

    id, err := app.users.CreateUserQuery(user)
    handleUserQueryResult(w, err, id, http.StatusCreated)
}

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.users.GetUserByIDQuery(id)
	handleUserQueryResult(w, err, user, http.StatusOK)
}

func (app *application) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.users.ListUsersQuery()
	handleUserQueryResult(w, err, users, http.StatusOK)
}

func (app *application) UpdateUser(w http.ResponseWriter, r *http.Request) {
    idParam := strings.TrimPrefix(r.URL.Path, "/user/")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        middleware.JSONError(w, err, http.StatusBadRequest)
        return
    }

    user, err := decodeAndValidateUser(r, w)
    if err != nil {
        middleware.JSONError(w, err, http.StatusBadRequest)
        return
    }

    err = app.users.UpdateUserQuery(id, user)
    handleUserQueryResult(w, err, nil, http.StatusOK)
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	err = app.users.DeleteUserQuery(id)
	handleUserQueryResult(w, err, nil, http.StatusNoContent)
}
