package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/youngjun827/api-std-lib/db"
	"github.com/youngjun827/api-std-lib/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	idParam := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	rows, err := db.DB.Query(`SELECT id, name, email, password FROM users`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	defer r.Body.Close()

	sqlStatement := `UPDATE users SET name=$1, email=$2, password=$3 WHERE id=$4`
	_, err = db.DB.Exec(sqlStatement, user.Name, user.Email, user.Password, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
