package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/youngjun827/api-std-lib/internal/database/models"
	"github.com/youngjun827/api-std-lib/internal/middleware"
)


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

	id, err := app.users.CreateUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, fmt.Errorf("User already exists"), http.StatusBadRequest)
			return
		}
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.users.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, fmt.Errorf("User with ID %d not found", id), http.StatusNotFound)
			return
		}
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (app *application) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.users.ListUsers()
	if err != nil {
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
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

	err = app.users.UpdateUser(id, user)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, fmt.Errorf("User with ID %d not found", id), http.StatusNotFound)
			return
		}
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	err = app.users.DeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, fmt.Errorf("User with ID %d not found", id), http.StatusNotFound)
			return
		}
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
