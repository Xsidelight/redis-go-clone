package redis_command

import (
	"redis-go-clone/internal/model"
	"sync"
)

func Exist(cmdArray []any, storedData map[string]model.StoredData, mu *sync.RWMutex) string {
	if len(cmdArray) < 2 {
		return "-ERR missing argument for EXIST\r\n"
	}

	key, ok := cmdArray[1].(string)
	if !ok {
		return "-ERR invalid argument for EXIST\r\n"
	}

	mu.RLock()
	_, found := storedData[key]
	mu.RUnlock()

	if found {
		return ":1\r\n"
	}

	return ":0\r\n"

}
