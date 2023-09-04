package models

import (
	"errors"
	"net/mail"
	"regexp"
	"unicode"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate performs validation on the user data.
func (u *User) Validate() error {
    if u.Name == "" {
        return errors.New("name is required")
    }
    if len(u.Name) < 3 {
        return errors.New("name should be at least 3 characters long")
    }
    if u.Email == "" {
        return errors.New("email is required")
    }
    if !u.ValidateEmail() {
        return errors.New("invalid email format")
    }
    if u.Password == "" {
        return errors.New("password is required")
    }
    if !u.ValidatePassword() {
        return errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, and one digit")
    }
    return nil
}

// ValidateEmail checks if the email is in a valid format.
func (u *User) ValidateEmail() bool {
    _, err := mail.ParseAddress(u.Email)
    if err != nil {
        return false
    }
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(u.Email)
}

// ValidatePassword checks if the password meets the criteria.
func (u *User) ValidatePassword() bool {
    var (
        hasUpper, hasLower, hasDigit bool
        length                       int
    )
    for _, char := range u.Password {
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
