package redis_command

import (
	"encoding/json"
	"os"
	"redis-go-clone/internal/model"
	"sync"
)

const saveFile = "data.json"

func Save(cmdArray []any, storedData map[string]model.StoredData, mu *sync.RWMutex) string {
	if len(cmdArray) != 1 {
		return "-ERR wrong number of arguments for 'SAVE' command\r\n"
	}

	mu.RLock()
	defer mu.RUnlock()

	data, err := json.MarshalIndent(storedData, "", "  ")
	if err != nil {
		return "-ERR error saving data\r\n"
	}

	err = os.WriteFile(saveFile, data, 0644)
	if err != nil {
		return "-ERR error saving data\r\n"
	}

	return "+OK\r\n"
}
