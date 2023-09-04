package middleware

import (
	"testing"

	"github.com/youngjun827/api-std-lib/api/models"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{email: "john.doe@example.com", valid: true},
		{email: "john.doe@example.com.au", valid: true},
		{email: "john.doe@123.com", valid: true},
		{email: "john.doe@example", valid: false},
		{email: "john.doe@example.com.", valid: false},
		{email: "john.doe@example.com.au.", valid: false},
		{email: "john.doe", valid: false},
		{email: "@example.com", valid: false},
	}

	for _, test := range tests {
		if actual := ValidateEmail(test.email); actual != test.valid {
			t.Errorf("ValidateEmail(%q) should return %t, but got %t", test.email, test.valid, actual)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		valid bool
	}{
		{password: "Password123", valid: true},
		{password: "Password123", valid: true},
		{password: "Pa$$word123", valid: true},
		{password: "123456", valid: false},
		{password: "password", valid: false},
		{password: "", valid: false},
	}

	for _, test := range tests {
		if actual := ValidatePassword(test.password); actual != test.valid {
			t.Errorf("ValidatePassword(%q) should return %t, but got %t", test.password, test.valid, actual)
		}
	}
}

func TestValidateUser(t *testing.T) {
	tests := []struct {
		user   models.User
		valid  bool
		errMsg string
	}{
		{user: models.User{Name: "", Email: "john.doe@example.com", Password: "password123"}, valid: false, errMsg: "name is required"},
		{user: models.User{Name: "John Doe", Email: "", Password: "password123"}, valid: false, errMsg: "email is required"},
		{user: models.User{Name: "John", Email: "invalidemail", Password: "password123"}, valid: false, errMsg: "invalid email format"},
		{user: models.User{Name: "John Doe", Email: "john.doe@example.com", Password: "123"}, valid: false, errMsg: "password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, and one digit"},
		{user: models.User{Name: "J", Email: "john.doe@example.com", Password: "Password123"}, valid: false, errMsg: "name should be at least 3 characters long"},
		{user: models.User{Name: "John Doe", Email: "john.doe@example.com", Password: ""}, valid: false, errMsg: "password is required"},
		{user: models.User{Name: "John Doe", Email: "john.doe@example.com", Password: "Password123"}, valid: true},
	}

	for _, test := range tests {
		err := ValidateUser(test.user)
		if test.valid {
			if err != nil {
				t.Errorf("ValidateUser(%v) should not return an error, but got %v", test.user, err)
			}
		} else {
			if err == nil {
				t.Errorf("ValidateUser(%v) should return an error with message %q, but got no error", test.user, test.errMsg)
			} else if err.Error() != test.errMsg {
				t.Errorf("ValidateUser(%v) should return an error with message %q, but got %q", test.user, test.errMsg, err.Error())
			}
		}
	}
}
