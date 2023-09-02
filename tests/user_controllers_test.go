package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/youngjun827/api-std-lib/controllers"
	"github.com/youngjun827/api-std-lib/db"
	"github.com/youngjun827/api-std-lib/models"
)

func TestCreateUser(t *testing.T) {
	db.InitDB()

	user := map[string]interface{}{
		"name":     "John",
		"email":    "john@example.com",
		"password": "password",
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

	var returnedID int
	err = json.NewDecoder(recorder.Body).Decode(&returnedID)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}
	if returnedID <= 0 {
		t.Errorf("Invalid ID returned: got %v", returnedID)
	}
}


func TestGetUser(t *testing.T) {
	// Initialize Database Connection
	db.InitDB()

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

	var user models.User
	err = json.NewDecoder(recorder.Body).Decode(&user)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if user.ID <= 0 {
		t.Errorf("Invalid ID returned: got %v", user.ID)
	}
}


func TestListUsers(t *testing.T) {
	db.InitDB()

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

	var users []models.User
	err = json.NewDecoder(recorder.Body).Decode(&users)
	if err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}

	if len(users) == 0 {
		t.Errorf("No users found: got %v", users)
	}
}
