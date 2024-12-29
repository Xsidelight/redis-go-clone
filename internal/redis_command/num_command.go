package redis_command

import (
	"redis-go-clone/internal/model"
	"sync"
)

func Incr(cmdArray []any, storedData map[string]model.StoredData, mu *sync.RWMutex) string {
	if len(cmdArray) < 2 {
		return "-ERR missing argument for INCR\r\n"
	}

	key, ok := cmdArray[1].(string)
	if !ok {
		return "-ERR invalid argument for INCR\r\n"
	}

	mu.Lock()
	defer mu.Unlock()

	value, found := storedData[key]
	if !found {
		return "-ERR key does not exist\r\n"
	}

	if v, ok := value.Value.(int); ok {
		v++
		storedData[key] = model.StoredData{Value: v}
		return "+OK\r\n"
	} else {
		return "-ERR value is not type of int\r\n"
	}
}

func Decr(cmdArray []any, storedData map[string]model.StoredData, mu *sync.RWMutex) string {
	if len(cmdArray) < 2 {
		return "-ERR missing argument for DECR\r\n"
	}

	key, ok := cmdArray[1].(string)
	if !ok {
		return "-ERR invalid argument for DECR\r\n"
	}

	mu.Lock()
	defer mu.Unlock()

	value, found := storedData[key]
	if !found {
		return "-ERR key does not exist\r\n"
	}

	if v, ok := value.Value.(int); ok {
		v--
		storedData[key] = model.StoredData{Value: v}
		return "+OK\r\n"
	} else {
		return "-ERR value is not type of int\r\n"
	}
}
