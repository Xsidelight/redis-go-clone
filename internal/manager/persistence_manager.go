package manager

import (
	"encoding/json"
	"log"
	"os"
	"redis-go-clone/internal/model"
	"sync"
)

const saveFile = "data.json"

func LoadData(storedData map[string]model.StoredData, mu *sync.RWMutex) {
	file, err := os.Open(saveFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No existing data file found. Starting with empty database.")
			return
		}
		log.Fatalf("Failed to open data file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	mu.Lock()
	defer mu.Unlock()
	err = decoder.Decode(&storedData)
	if err != nil {
		log.Fatalf("Failed to decode data file: %v", err)
	}

	log.Println("Database loaded successfully from disk.")
}
