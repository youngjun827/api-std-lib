package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/youngjun827/api-std-lib/db"
	"github.com/youngjun827/api-std-lib/logger"
	"github.com/youngjun827/api-std-lib/routes"
)

func main() {
	logger.InitLogger()
	db.InitDB()

	routes.SetupRoutes()

	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}