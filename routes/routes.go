package routes

import (
	"net/http"

	"github.com/youngjun827/api-std-lib/api/handlers"
	"github.com/youngjun827/api-std-lib/db"
	"github.com/youngjun827/api-std-lib/middleware"
)

func SetupRoutes(userRepository db.UserRepository) *http.ServeMux {
	mux := http.NewServeMux()

	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListUsers(w, r, userRepository)
		case http.MethodPost:
			handlers.CreateUser(w, r, userRepository)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	singleUserHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetUser(w, r, userRepository)
		case http.MethodPut:
			handlers.UpdateUser(w, r, userRepository)
		case http.MethodDelete:
			handlers.DeleteUser(w, r, userRepository)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/user", middleware.RateLimiter(userHandler))
	mux.Handle("/user/", middleware.RateLimiter(singleUserHandler))

	return mux
}
