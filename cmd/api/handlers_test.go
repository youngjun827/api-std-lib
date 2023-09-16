package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/youngjun827/api-std-lib/internal/database/mock"
	"github.com/youngjun827/api-std-lib/internal/database/models"
)

func setupMockApp() *application {
	return &application{
		users: &mock.MockUserModel{},
	}
}

func generateUserData() models.User {
	return models.User{
		Name:     "Test",
		Email:    "test@example.com",
		Password: "Jooa005500!",
	}
}

func TestCreateUser(t *testing.T) {
	app := setupMockApp()

	userData := generateUserData()

	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest(http.MethodPost, "/user/", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()

	app.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response map[string]int
	_ = json.Unmarshal(rr.Body.Bytes(), &response)
	if (response["id"] + 1) != mock.MockUser.ID {
		t.Errorf("Unexpected user ID: got %v want %v", response["id"], mock.MockUser.ID)
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	app := setupMockApp()
	invalidJSON := "this is not a JSON string"

	req, _ := http.NewRequest(http.MethodPost, "/user/", bytes.NewBufferString(invalidJSON))
	rr := httptest.NewRecorder()

	app.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for invalid JSON: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCreateUser_InvalidUser(t *testing.T) {
	app := setupMockApp()

	userData := generateUserData()
	userData.Name = ""

	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest(http.MethodPost, "/user/", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()

	app.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for invalid user: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCreateUser_UserExists(t *testing.T) {
	app := setupMockApp()
	app.users = mock.NewMockUserModel(true) 

	userData := generateUserData()
	userData.Email = "exists@example.com"

	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest(http.MethodPost, "/user/", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()

	app.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for existing user: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCreateUser_UnexpectedError(t *testing.T) {
	app := setupMockApp()
	app.users = mock.NewMockUserModel(true) 

	userData := generateUserData()
	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest(http.MethodPost, "/user/", bytes.NewBuffer(jsonData))
	rr := httptest.NewRecorder()

	app.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code for unexpected error: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestGetUser(t *testing.T) {
    app := setupMockApp()

    userID := 1

    req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", userID), nil)
    rr := httptest.NewRecorder()

    app.GetUser(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedUser := mock.MockUser

    var responseUser models.User
    if err := json.Unmarshal(rr.Body.Bytes(), &responseUser); err != nil {
        t.Errorf("Error decoding JSON response: %v", err)
    }

    if responseUser.ID != expectedUser.ID ||
        responseUser.Name != expectedUser.Name ||
        responseUser.Email != expectedUser.Email ||
        responseUser.Password != expectedUser.Password {
        t.Errorf("Handler returned unexpected user data: got %+v, want %+v", responseUser, expectedUser)
    }
}

func TestGetUser_InvalidID(t *testing.T) {
    app := setupMockApp()

    req, _ := http.NewRequest(http.MethodGet, "/user/invalid", nil)
    rr := httptest.NewRecorder()

    app.GetUser(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("Handler returned wrong status code for invalid ID: got %v want %v", status, http.StatusNotFound)
    }
}

func TestGetUser_UserNotFound(t *testing.T) {
    app := setupMockApp()

    nonExistingUserID := 999

    req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", nonExistingUserID), nil)
    rr := httptest.NewRecorder()

    app.GetUser(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("Handler returned wrong status code for user not found: got %v want %v", status, http.StatusNotFound)
    }
}

func TestGetUser_InternalServerError(t *testing.T) {
    app := setupMockApp()

    app.users.(*mock.MockUserModel).ErrorMode = true

    userID := 1

    req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", userID), nil)
    rr := httptest.NewRecorder()

    app.GetUser(rr, req)

    if status := rr.Code; status != http.StatusInternalServerError {
        t.Errorf("Handler returned wrong status code for unexpected error: got %v want %v", status, http.StatusInternalServerError)
    }
}

func TestListUsers(t *testing.T) {
	app := setupMockApp()

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app.ListUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	var users []models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &users); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	expectedUser := mock.MockUser
	if len(users) != 1 || users[0].ID != expectedUser.ID || users[0].Name != expectedUser.Name {
		t.Errorf("Handler returned unexpected user data: got %v, want %v", users, []models.User{expectedUser})
	}
}
func TestListUsers_InternalServerError(t *testing.T) {
	
	app := setupMockApp()
	app.users = mock.NewMockUserModel(true)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app.ListUsers(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code for unexpected error: got %v, want %v", status, http.StatusInternalServerError)
	}
}

func TestUpdateUser_StatusOK(t *testing.T) {
	app := setupMockApp()

	user := models.User{
		ID: 	  1,
		Name:     "UpdatedName",
		Email:    "updated@example.com",
		Password: "UpdatedPassword123",
	}
	userJSON, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPut, "/user/1", bytes.NewBuffer(userJSON))
	rr := httptest.NewRecorder()

	app.UpdateUser(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestUpdateUser_InvalidInputParameter(t *testing.T) {
    app := setupMockApp()

    req, _ := http.NewRequest(http.MethodPut, "/user/invalid", nil)
    rr := httptest.NewRecorder()

    app.UpdateUser(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("Handler returned wrong status code for valid input: got %v want %v", status, http.StatusNotFound)
    }
}


func TestUpdateUser_InvalidJSON(t *testing.T) {
    app := setupMockApp()

    req, _ := http.NewRequest(http.MethodPut, "/user/1", bytes.NewBufferString("this is not a JSON string"))
    rr := httptest.NewRecorder()

    app.UpdateUser(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("Handler returned wrong status code for invalid JSON: got %v want %v", status, http.StatusBadRequest)
    }
}

func TestUpdateUser_InvalidUser(t *testing.T) {
    app := setupMockApp()

    userData := generateUserData()
    userData.Name = ""

    jsonData, _ := json.Marshal(userData)

    req, _ := http.NewRequest(http.MethodPut, "/user/1", bytes.NewBuffer(jsonData))
    rr := httptest.NewRecorder()

    app.UpdateUser(rr, req)

    // Check the response status code.
    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("Handler returned wrong status code for invalid user: got %v want %v", status, http.StatusBadRequest)
    }
}

func TestUpdateUser_UserNotFound(t *testing.T) {
    app := setupMockApp()

    nonExistingUserID := 999

    userData := generateUserData()

    jsonData, _ := json.Marshal(userData)

    req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/user/%d", nonExistingUserID), bytes.NewBuffer(jsonData))
    rr := httptest.NewRecorder()

    app.UpdateUser(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("Handler returned wrong status code for user not found: got %v want %v", status, http.StatusNotFound)
    }
}

func TestUpdateUser_UnexpectedError(t *testing.T) {
    app := setupMockApp()
    app.users = mock.NewMockUserModel(true)

    userData := generateUserData()

    jsonData, _ := json.Marshal(userData)

    req, _ := http.NewRequest(http.MethodPut, "/user/1", bytes.NewBuffer(jsonData))
    rr := httptest.NewRecorder()

    app.UpdateUser(rr, req)

    if status := rr.Code; status != http.StatusInternalServerError {
        t.Errorf("Handler returned wrong status code for unexpected error: got %v want %v", status, http.StatusInternalServerError)
    }
}

func TestDeleteUser_Success(t *testing.T) {
    app := setupMockApp()

    // Assume the user with ID 1 exists
    userID := 1

    req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d", userID), nil)
    rr := httptest.NewRecorder()

    app.DeleteUser(rr, req)

    if status := rr.Code; status != http.StatusNoContent {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
    }
}

func TestDeleteUser_InvalidInputParameter(t *testing.T) {
    app := setupMockApp()

    req, _ := http.NewRequest(http.MethodDelete, "/user/invalid", nil)
    rr := httptest.NewRecorder()

    app.DeleteUser(rr, req)

    if status := rr.Code; status != http.StatusBadRequest {
        t.Errorf("Handler returned wrong status code for valid input: got %v want %v", status, http.StatusNotFound)
    }
}

func TestDeleteUser_UserNotFound(t *testing.T) {
    app := setupMockApp()

    // Assume the user with ID 999 does not exist
    nonExistingUserID := 999

    req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d", nonExistingUserID), nil)
    rr := httptest.NewRecorder()

    app.DeleteUser(rr, req)

    if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("Handler returned wrong status code for user not found: got %v want %v", status, http.StatusNotFound)
    }
}

func TestDeleteUser_UnexpectedError(t *testing.T) {
    app := setupMockApp()
    app.users = mock.NewMockUserModel(true) // Simulate an unexpected error

    // Assume the user with ID 1 exists
    userID := 1

    req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%d", userID), nil)
    rr := httptest.NewRecorder()

    app.DeleteUser(rr, req)

    if status := rr.Code; status != http.StatusInternalServerError {
        t.Errorf("Handler returned wrong status code for unexpected error: got %v want %v", status, http.StatusInternalServerError)
    }
}

