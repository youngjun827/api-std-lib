package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	db "github.com/youngjun827/api-std-lib/internal/database"
	"github.com/youngjun827/api-std-lib/internal/database/models"
	"github.com/youngjun827/api-std-lib/internal/middleware"
)

var userRepository db.UserRepository

func init() {
	userRepository = db.NewUserRepositorySQL(db.DB)
}

func CreateUser(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
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

	id, err := userRepository.CreateUser(user)
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

func GetUser(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	user, err := userRepository.GetUserByID(id)
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

func ListUsers(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
	users, err := userRepository.ListUsers()
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, fmt.Errorf("No users found"), http.StatusNotFound)
			return
		}
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
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

	err = userRepository.UpdateUser(id, user)
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

func DeleteUser(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	err = userRepository.DeleteUser(id)
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
