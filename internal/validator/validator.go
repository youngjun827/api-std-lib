package validator

import (
	"errors"
	"net/mail"
	"regexp"
	"unicode"

	"github.com/youngjun827/api-std-lib/internal/database/models"
)

type Validator struct{}

func (v *Validator) ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (v *Validator) ValidatePassword(password string) bool {
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

func (v *Validator) ValidateUser(user models.User) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if len(user.Name) < 3 {
		return errors.New("name should be at least 3 characters long")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if !v.ValidateEmail(user.Email) {
		return errors.New("invalid email format")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	if !v.ValidatePassword(user.Password) {
		return errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, and one digit")
	}
	return nil
}
