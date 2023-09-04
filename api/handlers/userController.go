package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/youngjun827/api-std-lib/api/models"
	"github.com/youngjun827/api-std-lib/cache"
	"github.com/youngjun827/api-std-lib/db"
	"github.com/youngjun827/api-std-lib/logger"
)

var userRepository db.UserRepository

func init() {
	logger.InitLogger()
	userRepository = db.NewUserRepositorySQL(db.DB)
}

func CreateUser(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
    var user models.User
    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&user); err != nil {
        logger.Error.Println("Failed to decode request body:", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := user.Validate(); err != nil {
        logger.Error.Println("Validation error:", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	defer r.Body.Close()

	id, err := userRepository.CreateUser(user)
	if err != nil {
		logger.Error.Println("Failed to create user:", err)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user, found := cache.GetUserFromCache(id); found {
		logger.Info.Println("Cache hit")
		json.NewEncoder(w).Encode(user)
		return
	}

	user, err := userRepository.GetUserByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error.Println("User not found:", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		logger.Error.Println("Failed to get user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cache.SetUserToCache(id, user)

	logger.Info.Printf("User fetched with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func ListUsers(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
	users, err := userRepository.ListUsers()
	if err != nil {
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
		logger.Error.Println("Invalid user ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    var user models.User
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&user); err != nil {
        logger.Error.Println("Failed to decode request body:", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := user.Validate(); err != nil {
        logger.Error.Println("Validation error:", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	defer r.Body.Close()

	err = userRepository.UpdateUser(id, user)
	if err != nil {
		logger.Error.Println("Failed to update user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info.Printf("User updated with ID: %d", id)

	w.WriteHeader(http.StatusNoContent)

	cache.SetUserToCache(id, user)

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, userRepository db.UserRepository) {
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logger.Error.Println("Invalid user ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = userRepository.DeleteUser(id)
	if err != nil {
		logger.Error.Println("Failed to delete user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info.Printf("User deleted with ID: %d", id)

	cache.DeleteUserFromCache(id)

	w.WriteHeader(http.StatusNoContent)
}
