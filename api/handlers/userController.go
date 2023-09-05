package handlers

import (
	"database/sql"
	"encoding/json"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateUser(user); err != nil {
		slog.Error("Failed to validate user password: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	id, err := userRepository.CreateUser(user)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Failed to create user: user already exists")
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}
		slog.Error("Failed to create user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := userRepository.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			// [Image of User not found icon]
			slog.Error("Failed to get user with the ID. User not found")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to get user with ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func ListUsers(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
	users, err := userRepository.ListUsers()
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Failed to list users: no users found")
			http.Error(w, "No users found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to list users: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		slog.Error("Failed to decode user data: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := middleware.ValidateUser(user); err != nil {
		slog.Error("Invalid user credentials: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = userRepository.UpdateUser(id, user)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Failed to update user with ID. User not found")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to update user with ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func DeleteUser(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		slog.Error("Failed to parse user ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = userRepository.DeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("User with the ID does not exist")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		slog.Error("Failed to delete user with ID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}