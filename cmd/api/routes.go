package main

import (
	"fmt"
	"net/http"

	"github.com/youngjun827/api-std-lib/internal/middleware"
)

func (app *application) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	usersHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.ListUsers(w, r)
		default:
			middleware.JSONError(w, fmt.Errorf("Method not allowed. Only GET method is allowed."), http.StatusMethodNotAllowed)
		}
	})

	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			app.CreateUser(w, r)
		case http.MethodGet:
			app.GetUser(w, r)
		case http.MethodPut:
			app.UpdateUser(w, r)
		case http.MethodDelete:
			app.DeleteUser(w, r)
		default:
			middleware.JSONError(w, fmt.Errorf("Method not allowed. POST, GET, PUT, DELETE methods are allowed."), http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/user/", middleware.RateLimiter(userHandler))
	mux.Handle("/users", middleware.RateLimiter(usersHandler))

	return mux
}
