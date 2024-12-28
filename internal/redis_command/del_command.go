package redis_command

import (
	"redis-go-clone/internal/model"
	"strconv"
	"sync"
)

func Del(cmdArray []any, storedData map[string]model.StoredData, mu *sync.RWMutex) string {
	if len(cmdArray) < 2 {
		return "-ERR missing argument for DEL\r\n"
	}

	var deletedCount int

	mu.Lock()
	for _, arg := range cmdArray[1:] {
		key, ok := arg.(string)
		if !ok {
			mu.Unlock()
			return "-ERR invalid argument for DEL\r\n"
		}
		if _, exists := storedData[key]; exists {
			delete(storedData, key)
			deletedCount++
		}
	}
	mu.Unlock()

	return ":" + strconv.Itoa(deletedCount) + "\r\n"
}
