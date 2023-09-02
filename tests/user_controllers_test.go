package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/youngjun827/api-std-lib/controllers"
	"github.com/youngjun827/api-std-lib/db"
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
