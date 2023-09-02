package routes

import (
	"net/http"

	"github.com/youngjun827/api-std-lib/controllers"
)

func SetupRoutes() {
	http.HandleFunc("/user", controllers.CreateUser)
}
