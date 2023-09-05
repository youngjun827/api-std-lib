package middleware

import (
	"encoding/json"
	"errors"
	"net/mail"
	"regexp"
	"unicode"

	"log/slog"

	"github.com/youngjun827/api-std-lib/api/models"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func ValidatePassword(password string) bool {
	var (
		hasUpper, hasLower, hasDigit bool
		length                       int
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
		length++
	}
	return length >= 8 && hasUpper && hasLower && hasDigit
}

func ValidateUser(user models.User) error {
	if user.Name == "" {
		slog.Error("Invalid user: name is required")
		return errors.New("name is required")
	}
	if len(user.Name) < 3 {
		slog.Error("Invalid user: name should be at least 3 characters long")
		return errors.New("name should be at least 3 characters long")
	}
	if user.Email == "" {
		slog.Error("Invalid user: email is required")
		return errors.New("email is required")
	}
	if !ValidateEmail(user.Email) {
		slog.Error("Invalid user: invalid email format")
		return errors.New("invalid email format")
	}
	if user.Password == "" {
		slog.Error("Invalid user: password is required")
		return errors.New("password is required")
	}
	if !ValidatePassword(user.Password) {
		slog.Error("Invalid user: password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, and one digit")
		return errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, and one digit")
	}
	return nil
}

func ReturnJSONError(err error) error {
	slog.Error("Error validating user: %v", err)
	jsonErr := map[string]string{
		"error": err.Error(),
	}
	jsonBytes, err := json.Marshal(jsonErr)
	if err != nil {
		return err
	}
	return errors.New(string(jsonBytes))
}
