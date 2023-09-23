package main

import (
	"fmt"
	"net/http"
)

func (app *application) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	usersHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.ListUsers(w, r)
		default:
			app.JsonErrorResponse(w, fmt.Errorf("Method not allowed. Only GET method is allowed."), http.StatusMethodNotAllowed)
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
			app.JsonErrorResponse(w, fmt.Errorf("Method not allowed. POST, GET, PUT, DELETE methods are allowed."), http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/user/", app.RateLimiter(userHandler))
	mux.Handle("/users", app.RateLimiter(usersHandler))

	return mux
}
