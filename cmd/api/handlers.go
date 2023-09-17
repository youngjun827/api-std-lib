package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/youngjun827/api-std-lib/internal/database/models"
)

func (app *application) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		app.JsonErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := app.ValidateUser(user); err != nil {
		app.JsonErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	id, err := app.users.CreateUserQuery(user)
	if err != nil {
		if errors.Is(err, models.ErrNoModels) {
			app.JsonErrorResponse(w, fmt.Errorf("User already exists"), http.StatusNotFound)
		}
		app.JsonErrorResponse(w, err, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.JsonErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.users.GetUserByIDQuery(id)
	if err != nil {
		if errors.Is(err, models.ErrNoModels) {
			app.JsonErrorResponse(w, fmt.Errorf("User with ID %d not found", id), http.StatusNotFound)
			return
		}
		app.JsonErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (app *application) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.users.ListUsersQuery()
	if err != nil {
		app.JsonErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (app *application) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.JsonErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var user models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		app.JsonErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := app.ValidateUser(user); err != nil {
		app.JsonErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = app.users.UpdateUserQuery(id, user)
	if err != nil {
		if errors.Is(err, models.ErrNoModels) {
			app.JsonErrorResponse(w, fmt.Errorf("User with ID %d not found", id), http.StatusNotFound)
			return
		}
		app.JsonErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.JsonErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	err = app.users.DeleteUserQuery(id)
	if err != nil {
		if errors.Is(err, models.ErrNoModels) {
			app.JsonErrorResponse(w, fmt.Errorf("User with ID %d not found", id), http.StatusNotFound)
			return
		}
		app.JsonErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}