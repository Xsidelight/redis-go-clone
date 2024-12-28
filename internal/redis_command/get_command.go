package redis_command

import (
	"redis-go-clone/internal/model"
	"redis-go-clone/pkg/resp"
	"sync"
	"time"
)

func Get(cmdArray []any, storedData map[string]model.StoredData, mu *sync.RWMutex) string {
	if len(cmdArray) < 2 {
		return "-ERR missing argument for GET\r\n"
	}

	key, ok := cmdArray[1].(string)
	if !ok {
		return "-ERR invalid argument for GET\r\n"
	}

	mu.RLock()
	value, found := storedData[key]
	mu.RUnlock()

	if !found {
		return "$-1\r\n" // Key not found
	}

	if value.ExpiryDate > 0 {
		expiryTime := time.Unix(value.ExpiryDate, 0)
		if time.Now().After(expiryTime) {
			// Key expired; delete it
			mu.Lock()
			delete(storedData, key)
			mu.Unlock()
			return "$-1\r\n"
		}
	}

	respValue := resp.SerializeRESP(value.Value, true)
	if respValue == "" {
		return "-ERR serialization error\r\n"
	}

	return respValue
}
