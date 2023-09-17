package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/youngjun827/api-std-lib/internal/database/models"
	"github.com/youngjun827/api-std-lib/internal/middleware"
)

func decodeAndValidateUser(r *http.Request, w http.ResponseWriter) (models.User, error) {
    var user models.User
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&user); err != nil {
        return user, err
    }

    if err := middleware.ValidateUser(user); err != nil {
        return user, err
    }

    return user, nil
}

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