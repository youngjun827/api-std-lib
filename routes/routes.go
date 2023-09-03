package routes

import (
	"net/http"

	"github.com/youngjun827/api-std-lib/controllers"
	"github.com/youngjun827/api-std-lib/rate"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.ListUsers(w, r)
		case http.MethodPost:
			controllers.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	singleUserHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.GetUser(w, r)
		case http.MethodPut:
			controllers.UpdateUser(w, r)
		case http.MethodDelete:
			controllers.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/user", rate.RateLimiter(userHandler))
	mux.Handle("/user/", rate.RateLimiter(singleUserHandler))

	return mux
}
