package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/youngjun827/api-std-lib/cache"
	"github.com/youngjun827/api-std-lib/db"
	"github.com/youngjun827/api-std-lib/logger"
	"github.com/youngjun827/api-std-lib/models"
	"github.com/youngjun827/api-std-lib/utility"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validator_err := utility.ValidateUser(user)
    if validator_err != nil {
        http.Error(w, validator_err.Error(), http.StatusBadRequest)
        return
    }

	defer r.Body.Close()

	sqlStatement := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	id := 0
	err := db.DB.QueryRow(sqlStatement, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("GetUser function started")

	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	  // Try to fetch user from cache first
	if user, found := cache.GetUserFromCache(id); found {
		logger.Info.Println("Cache hit")
		json.NewEncoder(w).Encode(user)
		return
	}

	var user models.User
	sqlStatement := `SELECT id, name, email, password FROM users WHERE id=$1`
	row := db.DB.QueryRow(sqlStatement, id)

	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
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

func ListUsers(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("ListUsers function started")

    userChannel := make(chan []models.User, 1)
    errorChannel := make(chan error, 1)

    go func() {
        fmt.Println("Goroutine for fetching users started")

        var users []models.User
        rows, err := db.DB.Query(`SELECT id, name, email, password FROM users`)
        if err != nil {
            errorChannel <- err
            return
        }
        defer rows.Close()

        for rows.Next() {
            var user models.User
            if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
                errorChannel <- err
                return
            }
            users = append(users, user)
        }

        if err := rows.Err(); err != nil {
            fmt.Println("Row error:", err)
            errorChannel <- err
            return
        }
        userChannel <- users
    }()

    select {
    case users := <-userChannel:
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    case err := <-errorChannel:
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

	logger.Info.Println("ListUsers function completed")
}


func UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("UpdateUser function started")

	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validator_err := utility.ValidateUser(user)
    if validator_err != nil {
        http.Error(w, validator_err.Error(), http.StatusBadRequest)
        return
    }

	defer r.Body.Close()

	sqlStatement := `UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4`
	_, err = db.DB.Exec(sqlStatement, user.Name, user.Email, user.Password, id)
	if err != nil {
		logger.Error.Println("Failed to update user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info.Printf("User updated with ID: %d", id)

	w.WriteHeader(http.StatusNoContent)

	// Update the cache
	cache.SetUserToCache(id, user)

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("DeleteUser function started")
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `DELETE FROM users WHERE id=$1`
	_, err = db.DB.Exec(sqlStatement, id)
	if err != nil {
		logger.Error.Println("Failed to delete user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info.Printf("User deleted with ID: %d", id)

	cache.DeleteUserFromCache(id)

    w.WriteHeader(http.StatusNoContent)
}
