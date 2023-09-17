package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/youngjun827/api-std-lib/internal/database/models"
	"github.com/youngjun827/api-std-lib/internal/middleware"
)

func handleUserQueryResult(w http.ResponseWriter, err error, data interface{}, successStatus int) {
    if err != nil {
        if errors.Is(err, models.ErrNoModels) {
            middleware.JSONError(w, fmt.Errorf("User not found"), http.StatusNotFound)
        } else {
            middleware.JSONError(w, err, http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if data != nil {
        w.WriteHeader(successStatus)
        json.NewEncoder(w).Encode(data)
    } else {
        w.WriteHeader(successStatus)
    }
}

func (app *application) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateUser(user); err != nil {
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

	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateUser(user); err != nil {
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
