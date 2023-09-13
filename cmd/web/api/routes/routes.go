package routes

import (
	"fmt"
	"net/http"

	"github.com/youngjun827/api-std-lib/cmd/web/api/handlers"
	"github.com/youngjun827/api-std-lib/internal/database"
	"github.com/youngjun827/api-std-lib/internal/middleware"
)

func SetupRoutes(userRepository database.UserRepository) *http.ServeMux {
	mux := http.NewServeMux()

	usersHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListUsers(w, r, userRepository)
		default:
			middleware.JSONError(w, fmt.Errorf("Method not allowed. Only GET method is allowed."), http.StatusMethodNotAllowed)
		}
	})

	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateUser(w, r, userRepository)
		case http.MethodGet:
			handlers.GetUser(w, r, userRepository)
		case http.MethodPut:
			handlers.UpdateUser(w, r, userRepository)
		case http.MethodDelete:
			handlers.DeleteUser(w, r, userRepository)
		default:
			middleware.JSONError(w, fmt.Errorf("Method not allowed. POST, GET, PUT, DELETE methods are allowed."), http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/user/", middleware.RateLimiter(userHandler))
	mux.Handle("/users", middleware.RateLimiter(usersHandler))

	return mux
}
