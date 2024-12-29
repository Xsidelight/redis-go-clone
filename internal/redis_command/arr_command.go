package redis_command

import (
	"redis-go-clone/internal/model"
	"sync"
)

func LPush(cmdArray []any, storedData map[string]model.StoredData, mu *sync.RWMutex) string {
	if len(cmdArray) < 3 {
		return "-ERR missing argument for LPUSH\r\n"
	}

	key, ok := cmdArray[1].(string)
	if !ok {
		return "-ERR invalid argument for LPUSH\r\n"
	}

	mu.Lock()
	defer mu.Unlock()

	value, found := storedData[key]
	if !found {
		value = model.StoredData{Value: []any{}}
	}

	list, ok := value.Value.([]any)
	if !ok {
		return "-ERR value is not type of list\r\n"
	}

	for i := 2; i < len(cmdArray); i++ {
		list = append([]any{cmdArray[i]}, list...)
	}

	storedData[key] = model.StoredData{Value: list}
	return "+OK\r\n"
}

func RPush(cmdArray []any, storedData map[string]model.StoredData, mu *sync.RWMutex) string {
	if len(cmdArray) < 3 {
		return "-ERR missing argument for RPUSH\r\n"
	}

	key, ok := cmdArray[1].(string)
	if !ok {
		return "-ERR invalid argument for RPUSH\r\n"
	}

	mu.Lock()
	defer mu.Unlock()

	value, found := storedData[key]
	if !found {
		value = model.StoredData{Value: []any{}}
	}

	list, ok := value.Value.([]any)
	if !ok {
		return "-ERR value is not type of list\r\n"
	}

	for i := 2; i < len(cmdArray); i++ {
		list = append(list, cmdArray[i])
	}

	storedData[key] = model.StoredData{Value: list}
	return "+OK\r\n"
}