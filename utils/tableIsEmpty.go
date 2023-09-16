package utils

import "github.com/youngjun827/api-std-lib/internal/database/models"

func IsEmptyResult(users []models.User) bool {
	return len(users) == 0
}
