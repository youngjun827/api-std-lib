package routes

import (
	"net/http"

	"github.com/youngjun827/api-std-lib/controllers"
	"github.com/youngjun827/api-std-lib/utility"
)

func SetupRoutes() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		utility.RateLimiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				controllers.ListUsers(w, r)
			case "POST":
				controllers.CreateUser(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})).ServeHTTP(w, r)
	})

	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		utility.RateLimiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				controllers.GetUser(w, r)
			case "PUT":
				controllers.UpdateUser(w, r)
			case "DELETE":
				controllers.DeleteUser(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})).ServeHTTP(w, r)
	})
}
