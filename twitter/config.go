package twitter

import (
	"sync"
)

var config struct {
	token string
	mutex sync.RWMutex
}

func getToken() string {
	config.mutex.RLock()
	token := config.token
	config.mutex.RUnlock()
	return token
}

func SetToken(token string) {
	config.mutex.Lock()
	config.token = token
	config.mutex.Unlock()
}
