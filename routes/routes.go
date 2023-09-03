package routes

import (
	"net/http"

	"github.com/youngjun827/api-std-lib/controllers"
	"github.com/youngjun827/api-std-lib/db"
	"github.com/youngjun827/api-std-lib/rate"
)

func SetupRoutes(userRepository db.UserRepository) *http.ServeMux {
	mux := http.NewServeMux()

	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListUsers(w, r, userRepository)
		case http.MethodPost:
			controllers.CreateUser(w, r, userRepository)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	singleUserHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetUser(w, r, userRepository)
		case http.MethodPut:
			controllers.UpdateUser(w, r, userRepository)
		case http.MethodDelete:
			controllers.DeleteUser(w, r, userRepository)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/user", rate.RateLimiter(userHandler))
	mux.Handle("/user/", rate.RateLimiter(singleUserHandler))

	return mux
}
