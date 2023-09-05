package cache

import (
	"testing"

	"github.com/youngjun827/api-std-lib/api/models"
)

func TestGetUserFromCache(t *testing.T) {
	// Add a user to the cache
	user := models.User{
		ID:   1,
		Name: "John Doe",
	}
	SetUserToCache(user.ID, user)

	// Get the user from the cache
	cachedUser, found := GetUserFromCache(user.ID)

	// Assert that the user was found in the cache
	if !found {
		t.Errorf("User not found in cache")
	}

	// Assert that the user is the same as the one that was added to the cache
	if cachedUser != user {
		t.Errorf("User in cache is not the same as the user that was added to the cache")
	}
}

func TestSetUserToCache(t *testing.T) {
	// Add a user to the cache
	user := models.User{
		ID:   1,
		Name: "John Doe",
	}
	SetUserToCache(user.ID, user)

	// Assert that the user is in the cache
	cachedUser, found := GetUserFromCache(user.ID)
	if !found {
		t.Errorf("User not found in cache")
	}

	// Assert that the user is the same as the one that was added to the cache
	if cachedUser != user {
		t.Errorf("User in cache is not the same as the user that was added to the cache")
	}
}

func TestDeleteUserFromCache(t *testing.T) {
	// Add a user to the cache
	user := models.User{
		ID:   1,
		Name: "John Doe",
	}
	SetUserToCache(user.ID, user)

	// Delete the user from the cache
	DeleteUserFromCache(user.ID)

	// Assert that the user is not in the cache
	_, found := GetUserFromCache(user.ID)
	if found {
		t.Errorf("User still found in cache")
	}
}
