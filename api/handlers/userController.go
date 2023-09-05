package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/youngjun827/api-std-lib/api/models"
	"github.com/youngjun827/api-std-lib/db"
	"github.com/youngjun827/api-std-lib/middleware"
)

var userRepository db.UserRepository

func init() {
	userRepository = db.NewUserRepositorySQL(db.DB)
}

func CreateUser(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
	var user models.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		slog.Error("Failed to decode user credentials: %v", err)
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateUser(user); err != nil {
		slog.Error("Failed to validate user password: %v", err)
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	id, err := userRepository.CreateUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, fmt.Errorf("User already exists"), http.StatusBadRequest)
			return
		}
		slog.Error("Failed to create user: %v", err)
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
		slog.Error("Failed to parse user ID: %v", err)
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	user, err := userRepository.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, fmt.Errorf("User with ID %d not found", id), http.StatusNotFound)
			return
		}
		slog.Error("Failed to get user with ID")
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
		slog.Error("Failed to list users: %v", err)
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
		slog.Error("Failed to parse user ID: %v", err)
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		slog.Error("Failed to decode user data: %v", err)
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateUser(user); err != nil {
		slog.Error("Invalid user credentials: %v", err)
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	err = userRepository.UpdateUser(id, user)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, fmt.Errorf("User with ID %d not found", id), http.StatusNotFound)
			return
		}
		slog.Error("Failed to update user with ID")
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error("Failed to parse user ID")
		middleware.JSONError(w, err, http.StatusBadRequest)
		return
	}

	err = userRepository.DeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			middleware.JSONError(w, fmt.Errorf("User with ID %d not found", id), http.StatusNotFound)
			return
		}
		slog.Error("Failed to delete user with ID")
		middleware.JSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
