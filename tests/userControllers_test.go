package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/youngjun827/api-std-lib/controllers"
	"github.com/youngjun827/api-std-lib/db"
)

func TestUserController(t *testing.T) {
	db.InitDB()
	resetDB()

	t.Run("Create user", func(t *testing.T) {
		user := map[string]interface{}{
			"name":     "John",
			"email":    "john@example.com",
			"password": "Uppercasea005500",
		}
		userJSON, _ := json.Marshal(user)

		req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(controllers.CreateUser)
		handler.ServeHTTP(recorder, req)

		if status := recorder.Code; status != http.StatusCreated {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})

	t.Run("Get user", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/user/1", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(controllers.GetUser)
		handler.ServeHTTP(recorder, req)

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("List users", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(controllers.ListUsers)
		handler.ServeHTTP(recorder, req)

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("Update user", func(t *testing.T) {
		user := map[string]interface{}{
			"name":     "John Updated",
			"email":    "john_updated@example.com",
			"password": "UpdatedUppercasea005500",
		}
		userJSON, _ := json.Marshal(user)

		req, err := http.NewRequest("PUT", "/user/1", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(controllers.UpdateUser)
		handler.ServeHTTP(recorder, req)

		if status := recorder.Code; status != http.StatusNoContent {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}
	})

	t.Run("Delete user", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/user/1", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(controllers.DeleteUser)
		handler.ServeHTTP(recorder, req)

		if status := recorder.Code; status != http.StatusNoContent {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}
	})
}

func resetDB() {
	conn := db.DB

	_, err := conn.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		log.Fatalf("Could not truncate table: %v", err)
	}
}
