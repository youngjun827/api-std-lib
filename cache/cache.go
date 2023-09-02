package cache

import (
	"sync"

	"github.com/youngjun827/api-std-lib/models"
)

var userCache = make(map[int]models.User)
var cacheMutex = &sync.Mutex{}

func GetUserFromCache(id int) (models.User, bool) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	user, found := userCache[id]
	return user, found
}

func SetUserToCache(id int, user models.User) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	userCache[id] = user
}

func DeleteUserFromCache(id int) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	delete(userCache, id)
}
